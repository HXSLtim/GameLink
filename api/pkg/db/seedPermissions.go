package db

import (
	"fmt"
	"hash/fnv"
	"log"
	"strings"

	"gorm.io/gorm"

	"gamelink/internal/model"
)

// seedSystemPermissions creates or updates all system permissions from GetAllPermissionDefinitions.
// This ensures all permissions defined in permission_codes.go are seeded with IsSystem = true.
// The function is idempotent - it can be run multiple times safely.
func seedSystemPermissions(tx *gorm.DB) error {
	definitions := model.GetAllPermissionDefinitions()

	for _, def := range definitions {
		if err := upsertSystemPermissionByMethodPath(tx, def); err != nil {
			return err
		}
	}

	log.Printf("Seeded %d system permissions\n", len(definitions))
	return nil
}

// upsertSystemPermissionByMethodPath upserts system permission with method+path as primary key.
// Code conflicts are automatically resolved with deterministic fallback codes.
func upsertSystemPermissionByMethodPath(tx *gorm.DB, def model.PermissionDefinition) error {
	var existing model.Permission
	err := tx.Where("method = ? AND path = ?", def.Method, def.Path).First(&existing).Error
	if err == nil {
		updates := map[string]interface{}{
			"is_system":   def.IsSystem,
			"description": def.Description,
			"group":       def.Group,
		}

		targetCode, codeChanged, codeErr := resolvePermissionCode(tx, def, existing.ID, existing.Code)
		if codeErr != nil {
			return codeErr
		}
		if codeChanged {
			updates["code"] = targetCode
		}

		if err := tx.Model(&existing).Updates(updates).Error; err != nil {
			return err
		}
		return nil
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	targetCode, _, codeErr := resolvePermissionCode(tx, def, 0, "")
	if codeErr != nil {
		return codeErr
	}

	perm := model.Permission{
		Code:        targetCode,
		Method:      def.Method,
		Path:        def.Path,
		Group:       def.Group,
		Description: def.Description,
		IsSystem:    def.IsSystem,
	}
	if err := tx.Create(&perm).Error; err != nil {
		// 极端情况下（并发/竞态）再次兜底：若 method+path 已存在，转为更新
		if isDuplicateKeyError(err) {
			var concurrent model.Permission
			if getErr := tx.Where("method = ? AND path = ?", def.Method, def.Path).First(&concurrent).Error; getErr == nil {
				return tx.Model(&concurrent).Updates(map[string]interface{}{
					"is_system":   def.IsSystem,
					"description": def.Description,
					"group":       def.Group,
				}).Error
			}
		}
		return err
	}
	return nil
}

// resolvePermissionCode returns a code that is safe to use under unique(code) constraint.
// changed=false means keep currentCode unchanged.
func resolvePermissionCode(tx *gorm.DB, def model.PermissionDefinition, currentID uint64, currentCode string) (code string, changed bool, err error) {
	trimmedCurrentCode := strings.TrimSpace(currentCode)
	preferredCode := strings.TrimSpace(def.Code)

	// 当前 code 可用则优先保持，避免无意义 churn
	if trimmedCurrentCode != "" {
		available, checkErr := isPermissionCodeAvailable(tx, trimmedCurrentCode, currentID)
		if checkErr != nil {
			return "", false, checkErr
		}
		if available {
			return trimmedCurrentCode, false, nil
		}
	}

	// 尝试使用定义中的 code
	if preferredCode != "" {
		available, checkErr := isPermissionCodeAvailable(tx, preferredCode, currentID)
		if checkErr != nil {
			return "", false, checkErr
		}
		if available {
			if preferredCode == trimmedCurrentCode {
				return preferredCode, false, nil
			}
			return preferredCode, true, nil
		}
	}

	// code 冲突时，按 method+path 生成稳定 fallback code
	for salt := 0; salt < 8; salt++ {
		fallback := buildFallbackPermissionCode(def, salt)
		available, checkErr := isPermissionCodeAvailable(tx, fallback, currentID)
		if checkErr != nil {
			return "", false, checkErr
		}
		if available {
			if fallback != preferredCode {
				log.Printf("permission code conflict on [%s] %s: prefer=%s, fallback=%s", def.Method, def.Path, preferredCode, fallback)
			}
			if fallback == trimmedCurrentCode {
				return fallback, false, nil
			}
			return fallback, true, nil
		}
	}

	return "", false, fmt.Errorf("failed to resolve unique permission code for [%s] %s", def.Method, def.Path)
}

func isPermissionCodeAvailable(tx *gorm.DB, code string, currentID uint64) (bool, error) {
	if strings.TrimSpace(code) == "" {
		return false, nil
	}

	query := tx.Model(&model.Permission{}).Where("code = ?", code)
	if currentID > 0 {
		query = query.Where("id <> ?", currentID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}

// buildFallbackPermissionCode creates a deterministic unique code for duplicated semantic codes.
// Format: module.resource.action (action carries hash).
func buildFallbackPermissionCode(def model.PermissionDefinition, salt int) string {
	module, resource := extractPathSegments(def.Path)
	action := shortActionName(def.Method)

	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(fmt.Sprintf("%s:%s:%d", def.Method, def.Path, salt)))
	hashPart := fmt.Sprintf("h%08x", hasher.Sum32())

	return fmt.Sprintf("%s.%s.%s%s", module, resource, action, hashPart)
}

func extractPathSegments(path string) (string, string) {
	clean := strings.TrimSpace(path)
	clean = strings.TrimPrefix(clean, "/")
	parts := strings.Split(clean, "/")
	segments := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || strings.HasPrefix(part, ":") {
			continue
		}
		if part == "api" || (strings.HasPrefix(part, "v") && len(part) >= 2 && part[1] >= '0' && part[1] <= '9') {
			continue
		}
		segments = append(segments, sanitizePermissionSegment(part))
	}

	module := "sys"
	resource := "route"
	if len(segments) > 0 && segments[0] != "" {
		module = segments[0]
	}
	if len(segments) > 1 && segments[1] != "" {
		resource = segments[1]
	}
	return module, resource
}

func shortActionName(method model.HTTPMethod) string {
	switch method {
	case model.HTTPMethodGET:
		return "g"
	case model.HTTPMethodPOST:
		return "p"
	case model.HTTPMethodPUT:
		return "u"
	case model.HTTPMethodPATCH:
		return "a"
	case model.HTTPMethodDELETE:
		return "d"
	default:
		return "x"
	}
}

func sanitizePermissionSegment(raw string) string {
	raw = strings.ToLower(raw)
	builder := strings.Builder{}
	builder.Grow(len(raw))
	for _, ch := range raw {
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') {
			builder.WriteRune(ch)
		}
	}
	segment := builder.String()
	if segment == "" {
		return "x"
	}
	if segment[0] >= '0' && segment[0] <= '9' {
		return "x" + segment
	}
	return segment
}

// markExistingPermissionsAsSystem marks all permissions that match system permission codes as IsSystem = true.
// This is useful for updating existing permissions that were created before IsSystem was added.
func markExistingPermissionsAsSystem(tx *gorm.DB) error {
	definitions := model.GetAllPermissionDefinitions()

	// Build a set of system permission codes
	systemCodes := make(map[string]bool, len(definitions))
	for _, def := range definitions {
		systemCodes[def.Code] = true
	}

	// Update all permissions with matching codes
	var permissions []model.Permission
	if err := tx.Find(&permissions).Error; err != nil {
		return err
	}

	updatedCount := 0
	for _, perm := range permissions {
		if systemCodes[perm.Code] && !perm.IsSystem {
			if err := tx.Model(&perm).Update("is_system", true).Error; err != nil {
				return err
			}
			updatedCount++
		}
	}

	if updatedCount > 0 {
		log.Printf("Marked %d existing permissions as system permissions\n", updatedCount)
	}

	return nil
}
