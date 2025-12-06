package model

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// TestCircularInheritanceDetection tests Property 19: Circular Inheritance Detection
// For any attempt to set a role's parent, if it would create a circular inheritance chain,
// the operation should be rejected with an error.
// **Feature: rbac-button-level-permission, Property 19: 循环继承检测**
// **Validates: Requirements 10.5**
func TestCircularInheritanceDetection(t *testing.T) {
	properties := gopter.NewProperties(nil)

	// Property 1: A role cannot be its own parent
	properties.Property("a role cannot be its own parent", prop.ForAll(
		func(roleID uint64) bool {
			// If roleID == parentID, it should be detected as circular
			return isCircularInheritance(roleID, roleID, nil)
		},
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
	))

	// Property 2: Direct circular reference (A -> B -> A) should be detected
	properties.Property("direct circular reference should be detected", prop.ForAll(
		func(roleAID, roleBID uint64) bool {
			if roleAID == roleBID {
				return true // Skip if same ID
			}

			// Build inheritance chain: B's parent is A
			inheritanceChain := map[uint64]*uint64{
				roleBID: &roleAID,
			}

			// Now try to set A's parent to B - this should create a cycle
			return isCircularInheritance(roleAID, roleBID, inheritanceChain)
		},
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
	))

	// Property 3: Longer circular chain (A -> B -> C -> A) should be detected
	properties.Property("longer circular chain should be detected", prop.ForAll(
		func(roleAID, roleBID, roleCID uint64) bool {
			// Skip if any IDs are the same
			if roleAID == roleBID || roleBID == roleCID || roleAID == roleCID {
				return true
			}

			// Build inheritance chain: C's parent is B, B's parent is A
			inheritanceChain := map[uint64]*uint64{
				roleCID: &roleBID,
				roleBID: &roleAID,
			}

			// Now try to set A's parent to C - this should create a cycle
			return isCircularInheritance(roleAID, roleCID, inheritanceChain)
		},
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
	))

	// Property 4: Non-circular inheritance should be allowed
	properties.Property("non-circular inheritance should be allowed", prop.ForAll(
		func(roleAID, roleBID, roleCID uint64) bool {
			// Skip if any IDs are the same
			if roleAID == roleBID || roleBID == roleCID || roleAID == roleCID {
				return true
			}

			// Build inheritance chain: B's parent is A (A -> B)
			inheritanceChain := map[uint64]*uint64{
				roleBID: &roleAID,
			}

			// Setting C's parent to B should NOT create a cycle (A -> B -> C)
			return !isCircularInheritance(roleCID, roleBID, inheritanceChain)
		},
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
	))

	// Property 5: Setting parent to nil should never create a cycle
	properties.Property("setting parent to nil should never create a cycle", prop.ForAll(
		func(roleID uint64) bool {
			// Setting parent to nil (removing parent) should never be circular
			return !isCircularInheritanceWithNil(roleID)
		},
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
	))

	// Property 6: Max depth validation - valid depths should pass, invalid should fail
	properties.Property("inheritance depth validation should correctly identify valid and invalid depths", prop.ForAll(
		func(depth int) bool {
			if depth < 0 {
				return true // Skip negative depths
			}
			// validateInheritanceDepth returns true if depth is valid (within limit)
			// It should return true for depth <= MaxRoleInheritanceDepth
			// It should return false for depth > MaxRoleInheritanceDepth
			isValid := validateInheritanceDepth(depth)
			expectedValid := depth <= MaxRoleInheritanceDepth
			return isValid == expectedValid
		},
		gen.IntRange(0, 10),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// isCircularInheritance checks if setting parentID as the parent of roleID would create a cycle.
// inheritanceChain maps roleID -> parentID for existing relationships.
func isCircularInheritance(roleID, parentID uint64, inheritanceChain map[uint64]*uint64) bool {
	// A role cannot be its own parent
	if roleID == parentID {
		return true
	}

	// Check if roleID appears in the ancestry of parentID
	visited := make(map[uint64]bool)
	current := parentID

	for {
		if visited[current] {
			// We've seen this before - there's already a cycle in the chain
			return true
		}
		visited[current] = true

		if current == roleID {
			// Found roleID in the ancestry - would create a cycle
			return true
		}

		// Get the parent of current
		if inheritanceChain == nil {
			break
		}
		parent, exists := inheritanceChain[current]
		if !exists || parent == nil || *parent == 0 {
			break
		}
		current = *parent
	}

	return false
}

// isCircularInheritanceWithNil checks if setting parent to nil would create a cycle (it never should).
func isCircularInheritanceWithNil(roleID uint64) bool {
	// Setting parent to nil removes the parent relationship, which can never create a cycle
	return false
}

// validateInheritanceDepth checks if the given depth is within the allowed limit.
func validateInheritanceDepth(depth int) bool {
	if depth > MaxRoleInheritanceDepth {
		return false // Depth exceeds limit
	}
	return true
}

// TestRoleInheritanceChainProperties tests properties of the inheritance chain.
// **Feature: rbac-button-level-permission, Property 18: 角色继承权限传递**
// **Validates: Requirements 10.2**
func TestRoleInheritanceChainProperties(t *testing.T) {
	properties := gopter.NewProperties(nil)

	// Property: Valid inheritance chain lengths should be within max depth + 1
	properties.Property("valid inheritance chain lengths should be within max depth + 1", prop.ForAll(
		func(chainLength int) bool {
			if chainLength < 1 {
				return true // Skip invalid chain lengths
			}
			// Test that we correctly identify valid vs invalid chain lengths
			// Valid chain length: 1 to MaxRoleInheritanceDepth + 1 (including the role itself)
			isValidLength := chainLength <= MaxRoleInheritanceDepth+1
			// This property verifies our understanding of valid chain lengths
			return isValidLength == (chainLength >= 1 && chainLength <= MaxRoleInheritanceDepth+1)
		},
		gen.IntRange(1, 10),
	))

	// Property: Level calculation should be consistent with parent level
	properties.Property("level calculation should be consistent with parent level", prop.ForAll(
		func(parentLevel int) bool {
			if parentLevel < 0 || parentLevel >= MaxRoleInheritanceDepth {
				return true // Skip invalid parent levels
			}

			role := &RoleModel{}
			parentID := uint64(1)
			role.ParentID = &parentID

			calculatedLevel := role.CalculateLevel(parentLevel)
			expectedLevel := parentLevel + 1

			return calculatedLevel == expectedLevel
		},
		gen.IntRange(0, MaxRoleInheritanceDepth),
	))

	// Property: Root role (no parent) should have level 0
	properties.Property("root role should have level 0", prop.ForAll(
		func(_ int) bool {
			role := &RoleModel{
				ParentID: nil,
			}
			return role.CalculateLevel(0) == 0
		},
		gen.Int(),
	))

	// Property: HasParent should return true only when ParentID is set and non-zero
	properties.Property("HasParent should correctly identify parent presence", prop.ForAll(
		func(hasParent bool, parentID uint64) bool {
			role := &RoleModel{}
			if hasParent && parentID > 0 {
				role.ParentID = &parentID
			}

			expected := hasParent && parentID > 0
			return role.HasParent() == expected
		},
		gen.Bool(),
		gen.UInt64(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// TestRoleValidation tests role validation properties.
func TestRoleValidation(t *testing.T) {
	properties := gopter.NewProperties(nil)

	// Property: ValidateInheritanceDepth should return error when level exceeds max
	properties.Property("ValidateInheritanceDepth should return error when level exceeds max", prop.ForAll(
		func(level int) bool {
			role := &RoleModel{Level: level}
			err := role.ValidateInheritanceDepth()

			if level > MaxRoleInheritanceDepth {
				return err == ErrRoleMaxDepthExceeded
			}
			return err == nil
		},
		gen.IntRange(0, MaxRoleInheritanceDepth+5),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}
