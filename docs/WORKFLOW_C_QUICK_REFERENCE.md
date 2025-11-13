# Workflow C: Quick Reference Guide

## Quick Start

### For Backend Developers

#### Running Tests
```bash
cd backend
go test -v ./internal/service/assignment
go test -v ./internal/service/assignment -run TestInitiateDisputeFullFlow
```

#### Key Files
- **Models**: `backend/internal/model/dispute.go`
- **Repository**: `backend/internal/repository/dispute/repository.go`
- **Service**: `backend/internal/service/assignment/service.go`
- **Handlers**: 
  - Admin: `backend/internal/handler/admin/dispute.go`
  - User: `backend/internal/handler/user/dispute.go`

#### Service Usage Example
```go
// Initialize service
svc := assignment.NewAssignmentService(
    disputeRepo,
    orderRepo,
    userRepo,
    opLogRepo,
    notifRepo,
    paymentRepo,
)

// Initiate dispute
resp, err := svc.InitiateDispute(ctx, assignment.InitiateDisputeRequest{
    OrderID:      123,
    UserID:       456,
    Reason:       "Service not provided",
    Description:  "Player did not show up",
    EvidenceURLs: []string{"https://example.com/screenshot.jpg"},
})

// Assign dispute
err := svc.AssignDispute(ctx, assignment.AssignDisputeRequest{
    DisputeID:        resp.DisputeID,
    AssignedToUserID: 789,
    Source:           model.AssignmentSourceSystem,
    ActorUserID:      999,
})

// Resolve dispute
err := svc.ResolveDispute(ctx, assignment.ResolveDisputeRequest{
    DisputeID:        resp.DisputeID,
    Resolution:       model.ResolutionRefund,
    ResolutionAmount: 10000,
    ResolutionNotes:  "Full refund approved",
    ActorUserID:      999,
})
```

### For Frontend Developers

#### API Endpoints

**User Endpoints:**
```
POST   /api/v1/user/orders/{id}/dispute
GET    /api/v1/user/orders/{id}/disputes
```

**Admin Endpoints:**
```
GET    /api/v1/admin/orders/pending-assign
GET    /api/v1/admin/orders/{id}/disputes
POST   /api/v1/admin/orders/{id}/assign
POST   /api/v1/admin/orders/{id}/assign/cancel
POST   /api/v1/admin/orders/{id}/mediate
```

#### React Query Hooks Example
```typescript
// Fetch pending disputes
const { data, isLoading } = useQuery({
  queryKey: ['disputes', 'pending'],
  queryFn: () => api.get('/admin/orders/pending-assign'),
});

// Assign dispute
const assignMutation = useMutation({
  mutationFn: (data) => api.post(`/admin/orders/${data.id}/assign`, data),
  onSuccess: () => queryClient.invalidateQueries(['disputes']),
});

// Resolve dispute
const resolveMutation = useMutation({
  mutationFn: (data) => api.post(`/admin/orders/${data.id}/mediate`, data),
  onSuccess: () => queryClient.invalidateQueries(['disputes']),
});
```

### For DevOps

#### Database Migration
```sql
-- Run migration script
-- Location: backend/db/migrations/XXX_create_order_disputes_table.sql
```

#### Monitoring Setup
```yaml
# Prometheus rules location
ops/alerts/dispute_sla_rules.yml

# Loki queries
{job="gamelink-backend"} | json | entity_type="dispute"
```

## Data Model Reference

### OrderDispute Fields
```go
type OrderDispute struct {
    ID                uint64            // Primary key
    OrderID           uint64            // Reference to order
    UserID            uint64            // User who initiated
    Status            DisputeStatus     // Current status
    Reason            string            // Dispute reason
    Description       string            // Detailed description
    EvidenceURLs      []string          // Screenshot URLs (max 9)
    AssignedToUserID  *uint64           // CS representative
    AssignmentSource  AssignmentSource  // system/manual/team
    AssignedAt        *time.Time        // Assignment timestamp
    SLADeadline       *time.Time        // CS response deadline
    SLABreached       bool              // SLA violation flag
    Resolution        DisputeResolution // Final decision
    ResolutionAmount  int64             // Refund amount (cents)
    ResolutionNotes   string            // Resolution notes
    ResolvedAt        *time.Time        // Resolution timestamp
    TraceID           string            // Audit trace ID
    CreatedAt         time.Time         // Creation timestamp
    UpdatedAt         time.Time         // Last update
}
```

### Status Lifecycle
```
pending → assigned → mediating → resolved
                  ↘ rejected
                  ↘ canceled
```

### Enums
```go
// DisputeStatus
"pending"    // Awaiting assignment
"assigned"   // Assigned to CS rep
"mediating"  // Being mediated
"resolved"   // Resolved
"rejected"   // Rejected
"canceled"   // Canceled

// DisputeResolution
"refund"     // Full refund
"partial"    // Partial refund
"reassign"   // Reassign order
"reject"     // Reject dispute
"pending"    // Not yet decided

// AssignmentSource
"system"     // System recommendation
"manual"     // Manual assignment
"team"       // Team assignment
```

## Common Operations

### Check SLA Status
```go
// Check if dispute exceeded SLA
if dispute.IsOverSLA() {
    // Handle SLA breach
}

// Get remaining time
remaining := dispute.GetSLARemaining() // seconds
```

### Validate Dispute Eligibility
```go
// Check if dispute can be initiated
if model.CanInitiateDispute(order) {
    // Can initiate dispute
}
```

### List Disputes with Filters
```go
disputes, total, err := svc.ListDisputesByStatus(ctx, 
    []model.DisputeStatus{model.DisputeStatusPending},
    1,  // page
    20, // pageSize
)
```

## Error Handling

### Common Errors
```go
assignment.ErrNotFound           // Dispute not found
assignment.ErrValidation         // Validation failed
assignment.ErrInvalidStatus      // Invalid status transition
assignment.ErrUnauthorized       // User not authorized
assignment.ErrCannotInitiateDispute // Outside 24h window
assignment.ErrDisputeExists      // Duplicate dispute
assignment.ErrOrderNotFound      // Order not found
```

### Error Response Example
```json
{
  "success": false,
  "code": 400,
  "message": "validation failed: reason is required"
}
```

## Testing

### Run All Tests
```bash
cd backend
make test
```

### Run Specific Test
```bash
go test -v ./internal/service/assignment -run TestInitiateDisputeFullFlow
```

### Run with Coverage
```bash
go test -v -coverprofile=coverage.out ./internal/service/assignment
go tool cover -html=coverage.out
```

## Logging

### Operation Log Format
```
entity_type: "dispute"
entity_id: <dispute_id>
action: "initiate_dispute" | "assign_dispute" | "resolve_dispute" | etc.
reason: <description>
trace_id: <uuid>
actor_user_id: <user_id>
timestamp: <created_at>
```

### Trace ID Usage
All operations include a trace ID for end-to-end tracing:
```
trace_id: "550e8400-e29b-41d4-a716-446655440000"
```

## Performance Tips

1. **Use pagination** for list queries (default 20 items)
2. **Index queries** on status, SLA deadline, trace ID
3. **Batch SLA checks** with `CheckAndMarkSLABreaches()`
4. **Cache user data** to avoid repeated lookups
5. **Use trace IDs** for debugging and monitoring

## Troubleshooting

### Dispute Not Found
- Verify dispute ID exists
- Check if dispute was soft-deleted
- Confirm user has access

### Cannot Initiate Dispute
- Check order status (must be in_progress or completed)
- Verify within 24h of completion
- Check for existing dispute

### SLA Not Triggering
- Verify SLA deadline is in the past
- Check if dispute status allows SLA check
- Run `CheckAndMarkSLABreaches()` manually

### Refund Not Processing
- Verify resolution is "refund"
- Check order has sufficient funds
- Verify payment integration

## Related Documentation

- [Full Implementation Guide](./WORKFLOW_C_IMPLEMENTATION_GUIDE.md)
- [Implementation Summary](./WORKFLOW_C_IMPLEMENTATION_SUMMARY.md)
- [GameLink Architecture](./ARCHITECTURE.md)
- [API Design Standards](./api/api-design-standards.md)

## Support

For issues or questions:
1. Check the implementation guide
2. Review integration tests for examples
3. Check operation logs for trace ID
4. Contact the development team
