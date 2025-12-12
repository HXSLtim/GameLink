package routingrule_test

import (
	"context"
	"encoding/json"
	"testing"

	"gamelink/internal/model"
	"gamelink/internal/repository/collectionentity"
	"gamelink/internal/repository/routingrule"
	svc "gamelink/internal/service/routingrule"
	"gamelink/pkg/testutil"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/require"
)

// **Feature: payment-finance-module, Property 15: 收款分流规则优先级**
// **Validates: Requirements 16.2**
//
// Property 15: Routing Rule Priority
// *For any* payment request, when multiple routing rules match, the system must select
// the rule with the highest priority (lowest numeric value) to determine the collection entity.

// TestProperty15_RoutingRulePriority tests that rules are matched by priority order
func TestProperty15_RoutingRulePriority(t *testing.T) {
	// Setup test database
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.CollectionEntity{},
		&model.PaymentChannelConfig{},
		&model.RoutingRule{},
		&model.RoutingRuleHistory{},
		&model.RoutingLog{},
		&model.User{},
	)
	defer testutil.CleanDB(t, db)

	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	service := svc.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create test user
	testUser := &model.User{
		Name:         "test_admin",
		Email:        "admin@test.com",
		PasswordHash: "hashed_password",
		Role:         model.RoleAdmin,
		Status:       model.UserStatusActive,
	}
	require.NoError(t, db.Create(testUser).Error)

	// Create multiple collection entities
	entities := make([]*model.CollectionEntity, 5)
	for i := 0; i < 5; i++ {
		entity := &model.CollectionEntity{
			Name:       genEntityName(i),
			CreditCode: genValidCreditCode(i),
			Status:     model.EntityStatusActive,
			IsDefault:  i == 0, // First entity is default
			CreatedBy:  testUser.ID,
		}
		require.NoError(t, db.Create(entity).Error)
		entities[i] = entity
	}

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property 15.1: When multiple rules match, the one with lowest priority number wins
	properties.Property("lowest priority number rule should be selected when multiple rules match", prop.ForAll(
		func(priorities []int) bool {
			// Clean up existing rules
			db.Exec("DELETE FROM routing_rules")

			// Create rules with different priorities, all matching the same condition
			gameType := "LOL"
			for i, priority := range priorities {
				entityIndex := i % len(entities)
				conditions := []model.RoutingCondition{
					{
						Field:    model.ConditionFieldGameType,
						Operator: model.ConditionOperatorEquals,
						Value:    mustMarshal(gameType),
					},
				}

				_, err := service.CreateRule(ctx, &model.CreateRoutingRuleRequest{
					Name:           genRuleName(i),
					Priority:       priority,
					Conditions:     conditions,
					TargetEntityID: entities[entityIndex].ID,
					Description:    "Priority test rule",
				}, testUser.ID)
				if err != nil {
					return false
				}
			}

			// Test routing with matching game type
			result, err := service.MatchCollectionEntity(ctx, &model.RoutingTestRequest{
				GameType: gameType,
			})
			if err != nil {
				return false
			}

			// Find the minimum priority
			minPriority := priorities[0]
			minPriorityIndex := 0
			for i, p := range priorities {
				if p < minPriority {
					minPriority = p
					minPriorityIndex = i
				}
			}

			// The matched entity should be the one with the lowest priority number
			expectedEntityID := entities[minPriorityIndex%len(entities)].ID
			return result.CollectionEntityID == expectedEntityID && !result.IsDefault
		},
		gen.SliceOfN(3, gen.IntRange(1, 100)),
	))

	// Property 15.2: Rules are evaluated in priority order (ascending)
	properties.Property("rules should be evaluated in ascending priority order", prop.ForAll(
		func(seed int) bool {
			// Clean up existing rules
			db.Exec("DELETE FROM routing_rules")

			// Create rules with specific priorities
			priorities := []int{10, 5, 15, 3, 8}
			gameTypes := []string{"LOL", "DOTA2", "CSGO", "VALORANT", "APEX"}

			for i, priority := range priorities {
				conditions := []model.RoutingCondition{
					{
						Field:    model.ConditionFieldGameType,
						Operator: model.ConditionOperatorEquals,
						Value:    mustMarshal(gameTypes[i]),
					},
				}

				_, err := service.CreateRule(ctx, &model.CreateRoutingRuleRequest{
					Name:           genRuleName(i),
					Priority:       priority,
					Conditions:     conditions,
					TargetEntityID: entities[i%len(entities)].ID,
					Description:    "Order test rule",
				}, testUser.ID)
				if err != nil {
					return false
				}
			}

			// Get active rules by priority
			rules, err := service.ListActiveRulesByPriority(ctx)
			if err != nil {
				return false
			}

			// Verify rules are sorted by priority (ascending)
			for i := 1; i < len(rules); i++ {
				if rules[i].Priority < rules[i-1].Priority {
					return false
				}
			}

			return true
		},
		gen.IntRange(1, 100),
	))

	properties.TestingRun(t)
}

// **Feature: payment-finance-module, Property 16: 收款分流默认主体回退**
// **Validates: Requirements 16.3, 17.4**
//
// Property 16: Default Entity Fallback
// *For any* payment request, when no routing rule matches, the system must use
// the entity marked as default to process the payment.

// TestProperty16_DefaultEntityFallback tests that default entity is used when no rules match
func TestProperty16_DefaultEntityFallback(t *testing.T) {
	// Setup test database
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.CollectionEntity{},
		&model.PaymentChannelConfig{},
		&model.RoutingRule{},
		&model.RoutingRuleHistory{},
		&model.RoutingLog{},
		&model.User{},
	)
	defer testutil.CleanDB(t, db)

	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	service := svc.NewRoutingRuleService(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create test user
	testUser := &model.User{
		Name:         "test_admin_fallback",
		Email:        "admin_fallback@test.com",
		PasswordHash: "hashed_password",
		Role:         model.RoleAdmin,
		Status:       model.UserStatusActive,
	}
	require.NoError(t, db.Create(testUser).Error)

	// Create default collection entity
	defaultEntity := &model.CollectionEntity{
		Name:       "Default Entity",
		CreditCode: "91110000300000000A",
		Status:     model.EntityStatusActive,
		IsDefault:  true,
		CreatedBy:  testUser.ID,
	}
	require.NoError(t, db.Create(defaultEntity).Error)

	// Create non-default entities
	for i := 0; i < 3; i++ {
		entity := &model.CollectionEntity{
			Name:       genEntityName(i + 10),
			CreditCode: genValidCreditCodeForFallback(i),
			Status:     model.EntityStatusActive,
			IsDefault:  false,
			CreatedBy:  testUser.ID,
		}
		require.NoError(t, db.Create(entity).Error)
	}

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property 16.1: When no rules exist, default entity should be used
	properties.Property("default entity should be used when no rules exist", prop.ForAll(
		func(gameType string) bool {
			// Clean up existing rules
			db.Exec("DELETE FROM routing_rules")

			// Test routing with any game type (no rules to match)
			result, err := service.MatchCollectionEntity(ctx, &model.RoutingTestRequest{
				GameType: gameType,
			})
			if err != nil {
				return false
			}

			// Should use default entity
			return result.CollectionEntityID == defaultEntity.ID && result.IsDefault
		},
		gen.AnyString(),
	))

	// Property 16.2: When rules exist but none match, default entity should be used
	properties.Property("default entity should be used when no rules match", prop.ForAll(
		func(ruleGameType string, requestGameType string) bool {
			// Skip if game types are the same (rule would match)
			if ruleGameType == requestGameType {
				return true
			}

			// Clean up existing rules
			db.Exec("DELETE FROM routing_rules")

			// Create a rule that won't match
			conditions := []model.RoutingCondition{
				{
					Field:    model.ConditionFieldGameType,
					Operator: model.ConditionOperatorEquals,
					Value:    mustMarshal(ruleGameType),
				},
			}

			_, err := service.CreateRule(ctx, &model.CreateRoutingRuleRequest{
				Name:           "Non-matching rule",
				Priority:       1,
				Conditions:     conditions,
				TargetEntityID: defaultEntity.ID,
				Description:    "Rule that won't match",
			}, testUser.ID)
			if err != nil {
				return false
			}

			// Test routing with different game type
			result, err := service.MatchCollectionEntity(ctx, &model.RoutingTestRequest{
				GameType: requestGameType,
			})
			if err != nil {
				return false
			}

			// Should use default entity since rule doesn't match
			return result.CollectionEntityID == defaultEntity.ID && result.IsDefault
		},
		gen.AlphaString(),
		gen.AlphaString(),
	))

	// Property 16.3: Default entity should be active
	properties.Property("default entity must be active to be used", prop.ForAll(
		func(seed int) bool {
			// Get default entity
			entity, err := service.GetDefaultEntity(ctx)
			if err != nil {
				return false
			}

			// Default entity should be active
			return entity.Status == model.EntityStatusActive
		},
		gen.IntRange(1, 100),
	))

	properties.TestingRun(t)
}

// Helper functions

func genEntityName(index int) string {
	return "Test Entity " + string(rune('A'+index%26))
}

func genValidCreditCode(index int) string {
	codes := []string{
		"91110000100000000A",
		"91110000100000000B",
		"91110000100000000C",
		"91110000100000000D",
		"91110000100000000E",
	}
	return codes[index%len(codes)]
}

func genValidCreditCodeForFallback(index int) string {
	codes := []string{
		"91110000400000000A",
		"91110000400000000B",
		"91110000400000000C",
	}
	return codes[index%len(codes)]
}

func genRuleName(index int) string {
	return "Test Rule " + string(rune('A'+index%26))
}

func mustMarshal(v interface{}) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
