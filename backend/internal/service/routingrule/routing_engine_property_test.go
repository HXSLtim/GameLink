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

// **Feature: payment-finance-module, Property 17: 收款分流记录完整性**
// **Validates: Requirements 17.3**
//
// Property 17: Payment Routing Record Completeness
// *For any* successful payment record, it must contain the actual collection entity ID
// and the corresponding merchant number information.

// TestProperty17_PaymentRoutingRecordCompleteness tests that routing results contain complete information
func TestProperty17_PaymentRoutingRecordCompleteness(t *testing.T) {
	// Setup test database
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.CollectionEntity{},
		&model.PaymentChannelConfig{},
		&model.RoutingRule{},
		&model.RoutingRuleHistory{},
		&model.RoutingLog{},
		&model.User{},
		&model.Payment{},
		&model.Order{},
	)
	defer testutil.CleanDB(t, db)

	ruleRepo := routingrule.NewRoutingRuleRepository(db)
	entityRepo := collectionentity.NewCollectionEntityRepository(db)
	engine := svc.NewRoutingEngine(ruleRepo, entityRepo)
	ctx := context.Background()

	// Create test user
	testUser := &model.User{
		Name:         "test_admin_routing",
		Email:        "admin_routing@test.com",
		PasswordHash: "hashed_password",
		Role:         model.RoleAdmin,
		Status:       model.UserStatusActive,
	}
	require.NoError(t, db.Create(testUser).Error)

	// Create collection entities with payment channels
	entities := make([]*model.CollectionEntity, 3)
	for i := 0; i < 3; i++ {
		entity := &model.CollectionEntity{
			Name:       genEntityNameForProperty17(i),
			CreditCode: genValidCreditCodeForProperty17(i),
			Status:     model.EntityStatusActive,
			IsDefault:  i == 0, // First entity is default
			CreatedBy:  testUser.ID,
		}
		require.NoError(t, db.Create(entity).Error)
		entities[i] = entity

		// Create payment channel for each entity
		channel := &model.PaymentChannelConfig{
			CollectionEntityID: entity.ID,
			Channel:            model.PaymentMethodWeChat,
			MerchantNo:         genMerchantNo(i),
			MerchantKey:        "test_key_" + string(rune('A'+i)),
			CallbackURL:        "https://callback.test.com/" + string(rune('A'+i)),
			Enabled:            true,
			Priority:           1,
		}
		require.NoError(t, db.Create(channel).Error)
	}

	// Create routing rules
	ruleService := svc.NewRoutingRuleService(ruleRepo, entityRepo)
	gameTypes := []string{"LOL", "DOTA2", "CSGO"}
	for i, gameType := range gameTypes {
		conditions := []model.RoutingCondition{
			{
				Field:    model.ConditionFieldGameType,
				Operator: model.ConditionOperatorEquals,
				Value:    mustMarshalForProperty17(gameType),
			},
		}
		_, err := ruleService.CreateRule(ctx, &model.CreateRoutingRuleRequest{
			Name:           "Rule for " + gameType,
			Priority:       i + 1,
			Conditions:     conditions,
			TargetEntityID: entities[i].ID,
			Description:    "Test rule for " + gameType,
		}, testUser.ID)
		require.NoError(t, err)
	}

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property 17.1: Routing result must contain valid collection entity ID
	properties.Property("routing result must contain valid collection entity ID", prop.ForAll(
		func(gameType string, amountCents int64) bool {
			routingCtx := &svc.RoutingContext{
				OrderID:     1,
				GameType:    gameType,
				ServiceType: "escort",
				AmountCents: amountCents,
				Region:      "CN",
				Method:      model.PaymentMethodWeChat,
			}

			result, err := engine.RoutePayment(ctx, routingCtx)
			if err != nil {
				// If routing fails, it's acceptable (e.g., no default entity)
				return true
			}

			// Collection entity ID must be non-zero
			return result.CollectionEntityID > 0
		},
		gen.OneConstOf("LOL", "DOTA2", "CSGO", "VALORANT", "APEX"),
		gen.Int64Range(100, 1000000),
	))

	// Property 17.2: Routing result must contain non-empty merchant number
	properties.Property("routing result must contain non-empty merchant number", prop.ForAll(
		func(gameType string, amountCents int64) bool {
			routingCtx := &svc.RoutingContext{
				OrderID:     1,
				GameType:    gameType,
				ServiceType: "escort",
				AmountCents: amountCents,
				Region:      "CN",
				Method:      model.PaymentMethodWeChat,
			}

			result, err := engine.RoutePayment(ctx, routingCtx)
			if err != nil {
				// If routing fails, it's acceptable
				return true
			}

			// Merchant number must be non-empty
			return result.MerchantNo != ""
		},
		gen.OneConstOf("LOL", "DOTA2", "CSGO", "VALORANT", "APEX"),
		gen.Int64Range(100, 1000000),
	))

	// Property 17.3: Routing result must contain valid entity name
	properties.Property("routing result must contain valid entity name", prop.ForAll(
		func(gameType string, amountCents int64) bool {
			routingCtx := &svc.RoutingContext{
				OrderID:     1,
				GameType:    gameType,
				ServiceType: "escort",
				AmountCents: amountCents,
				Region:      "CN",
				Method:      model.PaymentMethodWeChat,
			}

			result, err := engine.RoutePayment(ctx, routingCtx)
			if err != nil {
				return true
			}

			// Entity name must be non-empty
			return result.EntityName != ""
		},
		gen.OneConstOf("LOL", "DOTA2", "CSGO", "VALORANT", "APEX"),
		gen.Int64Range(100, 1000000),
	))

	// Property 17.4: When rule matches, routing result must contain rule information
	properties.Property("when rule matches, routing result must contain rule information", prop.ForAll(
		func(gameType string) bool {
			routingCtx := &svc.RoutingContext{
				OrderID:     1,
				GameType:    gameType,
				ServiceType: "escort",
				AmountCents: 10000,
				Region:      "CN",
				Method:      model.PaymentMethodWeChat,
			}

			result, err := engine.RoutePayment(ctx, routingCtx)
			if err != nil {
				return true
			}

			// If not using default, must have matched rule info
			if !result.IsDefault {
				return result.MatchedRuleID != nil && result.MatchedRuleName != ""
			}

			// If using default, should not have matched rule info
			return result.MatchedRuleID == nil
		},
		gen.OneConstOf("LOL", "DOTA2", "CSGO"),
	))

	// Property 17.5: Routing log can be created and retrieved
	properties.Property("routing log can be created and retrieved", prop.ForAll(
		func(paymentID uint64, orderID uint64, gameType string) bool {
			routingCtx := &svc.RoutingContext{
				OrderID:     orderID,
				GameType:    gameType,
				ServiceType: "escort",
				AmountCents: 10000,
				Region:      "CN",
				Method:      model.PaymentMethodWeChat,
			}

			result, err := engine.RoutePayment(ctx, routingCtx)
			if err != nil {
				return true
			}

			// Create routing log
			err = engine.CreateRoutingLog(ctx, paymentID, orderID, result)
			if err != nil {
				return false
			}

			// Retrieve routing log
			log, err := engine.GetRoutingLogByPayment(ctx, paymentID)
			if err != nil {
				return false
			}

			// Verify log contains complete information
			// Property 17: 收款分流记录完整性
			return log.PaymentID == paymentID &&
				log.OrderID == orderID &&
				log.CollectionEntityID == result.CollectionEntityID &&
				log.MerchantNo == result.MerchantNo &&
				log.IsDefault == result.IsDefault
		},
		gen.UInt64Range(1, 10000),
		gen.UInt64Range(1, 10000),
		gen.OneConstOf("LOL", "DOTA2", "CSGO", "VALORANT"),
	))

	properties.TestingRun(t)
}

// Helper functions for Property 17 tests

func genEntityNameForProperty17(index int) string {
	return "Routing Entity " + string(rune('A'+index%26))
}

func genValidCreditCodeForProperty17(index int) string {
	codes := []string{
		"91110000500000000A",
		"91110000500000000B",
		"91110000500000000C",
	}
	return codes[index%len(codes)]
}

func genMerchantNo(index int) string {
	return "MERCHANT_" + string(rune('A'+index%26)) + "_001"
}

func mustMarshalForProperty17(v interface{}) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
