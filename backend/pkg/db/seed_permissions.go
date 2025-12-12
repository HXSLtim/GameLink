package db

import (
	"log"

	"gorm.io/gorm"

	"gamelink/internal/model"
)

// seedSystemPermissions creates or updates all system permissions from GetAllPermissionDefinitions.
// This ensures all permissions defined in permission_codes.go are seeded with IsSystem = true.
// The function is idempotent - it can be run multiple times safely.
func seedSystemPermissions(tx *gorm.DB) error {
	definitions := model.GetAllPermissionDefinitions()

	for _, def := range definitions {
		var existing model.Permission

		// First try to find by code (primary identifier)
		err := tx.Where("code = ?", def.Code).First(&existing).Error
		if err == nil {
			// Permission exists by code, update IsSystem flag and other fields
			updates := map[string]interface{}{
				"is_system":   true,
				"description": def.Description,
				"group":       def.Group,
			}
			if err := tx.Model(&existing).Updates(updates).Error; err != nil {
				return err
			}
			continue
		} else if err != gorm.ErrRecordNotFound {
			return err
		}

		// Try to find by method+path (for backward compatibility)
		err = tx.Where("method = ? AND path = ?", def.Method, def.Path).First(&existing).Error
		if err == nil {
			// Permission exists by method+path, update code and IsSystem
			updates := map[string]interface{}{
				"code":        def.Code,
				"is_system":   true,
				"description": def.Description,
				"group":       def.Group,
			}
			if err := tx.Model(&existing).Updates(updates).Error; err != nil {
				return err
			}
			continue
		} else if err != gorm.ErrRecordNotFound {
			return err
		}

		// Permission doesn't exist, create it
		perm := model.Permission{
			Code:        def.Code,
			Method:      def.Method,
			Path:        def.Path,
			Group:       def.Group,
			Description: def.Description,
			IsSystem:    true,
		}
		if err := tx.Create(&perm).Error; err != nil {
			// Skip duplicates (race condition protection)
			if isDuplicateKeyError(err) {
				log.Printf("Permission %s already exists, skipping\n", def.Code)
				continue
			}
			return err
		}
	}

	log.Printf("Seeded %d system permissions\n", len(definitions))
	return nil
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
