# Workflow C: Order Dispute & Customer Service Assignment Implementation Guide

## Overview

This document provides a comprehensive guide for implementing Workflow C: Order Dispute & Customer Service Assignment Enhancement for the GameLink platform.

## Architecture

### Data Models

#### OrderDispute Entity
Location: `backend/internal/model/dispute.go`

**Key Fields:**
- `OrderID` - Reference to the order
- `UserID` - User who initiated the dispute
- `Status` - Dispute lifecycle state (pending, assigned, mediating, resolved, rejected, canceled)
- `Reason` - Primary reason for dispute
- `Description` - Detailed description
- `EvidenceURLs` - Array of screenshot URLs (max 9)
- `AssignedToUserID` - Customer service representative assigned
- `AssignmentSource` - Where assignment came from (system/manual/team)
- `SLADeadline` - Deadline for CS response (default 30 minutes)
- `SLABreached` - Flag indicating SLA violation
- `Resolution` - Final decision (refund/partial/reassign/reject)
- `ResolutionAmount` - Refund amount in cents
- `TraceID` - Unique trace ID for logging

**Status Lifecycle:**
```
pending → assigned → mediating → resolved
                  ↘ rejected
                  ↘ canceled
```

#### AssignmentSource Enum
```go
const (
    AssignmentSourceSystem AssignmentSource = "system"   // System recommendation
    AssignmentSourceManual AssignmentSource = "manual"   // Manual assignment
    AssignmentSourceTeam   AssignmentSource = "team"     // Team assignment
)
```

### Repository Layer

Location: `backend/internal/repository/dispute/repository.go`

**Interface Methods:**
- `Create(ctx, dispute)` - Create new dispute
- `Get(ctx, id)` - Get dispute by ID
- `GetByOrderID(ctx, orderID)` - Get dispute for an order
- `Update(ctx, dispute)` - Update dispute
- `List(ctx, opts)` - List disputes with filters
- `ListPendingAssignment(ctx, page, pageSize)` - List unassigned disputes
- `ListSLABreached(ctx)` - List disputes exceeding SLA
- `MarkSLABreached(ctx, disputeID)` - Mark SLA as breached
- `CountByStatus(ctx, status)` - Count disputes by status
- `GetPendingCount(ctx)` - Get pending dispute count

### Service Layer

Location: `backend/internal/service/assignment/service.go`

**Key Methods:**

1. **InitiateDispute**
   - Validates order ownership
   - Checks if dispute can be initiated (24h window)
   - Creates dispute with SLA deadline
   - Logs operation with trace ID
   - Returns dispute ID and SLA deadline

2. **AssignDispute**
   - Assigns dispute to customer service rep
   - Validates assignment source
   - Updates dispute status to "assigned"
   - Sends notification to assigned user
   - Logs operation

3. **ResolveDispute**
   - Processes dispute resolution
   - Handles refund if needed
   - Updates order status
   - Sends notification to user
   - Logs operation

4. **RollbackAssignment**
   - Reverts assignment back to pending
   - Allows reassignment
   - Logs rollback reason
   - Maintains audit trail

5. **CheckAndMarkSLABreaches**
   - Scheduled task to check SLA violations
   - Marks breached disputes
   - Sends alerts to assigned users
   - Logs breaches

### Handler Layer

#### Admin Handlers
Location: `backend/internal/handler/admin/dispute.go`

**Endpoints:**
- `GET /admin/orders/pending-assign` - List pending disputes
- `GET /admin/orders/{id}/disputes` - Get dispute detail
- `POST /admin/orders/{id}/assign` - Assign dispute
- `POST /admin/orders/{id}/assign/cancel` - Rollback assignment
- `POST /admin/orders/{id}/mediate` - Resolve dispute

#### User Handlers
Location: `backend/internal/handler/user/dispute.go`

**Endpoints:**
- `POST /user/orders/{id}/dispute` - Initiate dispute
- `GET /user/orders/{id}/disputes` - Get dispute detail

## API Specifications

### User Endpoints

#### POST /api/v1/user/orders/{id}/dispute
**Request:**
```json
{
  "orderId": 123,
  "reason": "Service not provided",
  "description": "Player did not show up for the scheduled service",
  "evidenceUrls": [
    "https://example.com/screenshot1.jpg",
    "https://example.com/screenshot2.jpg"
  ]
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "code": 201,
  "data": {
    "disputeId": 456,
    "traceId": "uuid-string",
    "slaDeadline": "2025-01-13T10:30:00Z"
  }
}
```

**Error Cases:**
- 400: Invalid request (missing fields, too many evidence URLs)
- 404: Order not found
- 409: Cannot initiate dispute (outside 24h window) or dispute already exists
- 403: Unauthorized (not order owner)

#### GET /api/v1/user/orders/{id}/disputes
**Response (200 OK):**
```json
{
  "success": true,
  "code": 200,
  "data": {
    "id": 456,
    "orderId": 123,
    "userId": 1,
    "status": "assigned",
    "reason": "Service not provided",
    "description": "...",
    "evidenceUrls": ["..."],
    "assignedToUserId": 789,
    "assignmentSource": "system",
    "assignedAt": "2025-01-13T09:00:00Z",
    "slaDeadline": "2025-01-13T10:30:00Z",
    "slaBreached": false,
    "resolution": "pending",
    "traceId": "uuid-string",
    "createdAt": "2025-01-13T09:00:00Z"
  }
}
```

### Admin Endpoints

#### GET /api/v1/admin/orders/pending-assign
**Query Parameters:**
- `page` (optional, default: 1)
- `pageSize` (optional, default: 20, max: 100)

**Response (200 OK):**
```json
{
  "success": true,
  "code": 200,
  "data": {
    "disputes": [
      {
        "id": 456,
        "orderId": 123,
        "userId": 1,
        "status": "pending",
        "reason": "Service not provided",
        "slaDeadline": "2025-01-13T10:30:00Z",
        "slaBreached": false,
        "createdAt": "2025-01-13T09:00:00Z"
      }
    ],
    "total": 42,
    "page": 1,
    "pageSize": 20
  }
}
```

#### POST /api/v1/admin/orders/{id}/assign
**Request:**
```json
{
  "assignedToUserId": 789,
  "source": "system"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "code": 200,
  "data": {
    "message": "Dispute assigned successfully"
  }
}
```

#### POST /api/v1/admin/orders/{id}/assign/cancel
**Request:**
```json
{
  "rollbackReason": "Assigned user unavailable"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "code": 200,
  "data": {
    "message": "Assignment rolled back successfully"
  }
}
```

#### POST /api/v1/admin/orders/{id}/mediate
**Request:**
```json
{
  "resolution": "refund",
  "resolutionAmount": 10000,
  "resolutionNotes": "Full refund approved due to service not provided"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "code": 200,
  "data": {
    "message": "Dispute resolved with refund decision"
  }
}
```

## Database Migration

Create migration file: `backend/db/migrations/XXX_create_order_disputes_table.sql`

```sql
CREATE TABLE order_disputes (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'pending' INDEX,
    reason TEXT NOT NULL,
    description TEXT,
    evidence_urls JSON,
    assigned_to_user_id BIGINT,
    assignment_source VARCHAR(32),
    assigned_at TIMESTAMP NULL,
    sla_deadline TIMESTAMP NULL INDEX,
    sla_breached BOOLEAN DEFAULT FALSE,
    sla_breached_at TIMESTAMP NULL,
    resolution VARCHAR(32) DEFAULT 'pending',
    resolution_amount BIGINT DEFAULT 0,
    resolution_notes TEXT,
    resolved_at TIMESTAMP NULL,
    resolved_by_user_id BIGINT,
    rolled_back_at TIMESTAMP NULL,
    rolled_back_by_user_id BIGINT,
    rollback_reason TEXT,
    trace_id VARCHAR(64) INDEX,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP INDEX,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL INDEX,
    
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (assigned_to_user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (resolved_by_user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (rolled_back_by_user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Add dispute fields to orders table
ALTER TABLE orders ADD COLUMN has_dispute BOOLEAN DEFAULT FALSE INDEX;
ALTER TABLE orders ADD COLUMN dispute_id BIGINT INDEX;
ALTER TABLE orders ADD FOREIGN KEY (dispute_id) REFERENCES order_disputes(id) ON DELETE SET NULL;
```

## Integration Tests

Location: `backend/internal/service/assignment/integration_test.go`

**Test Scenarios:**

1. **Full Dispute Flow**
   - User initiates dispute
   - System assigns to CS rep
   - CS rep resolves with refund
   - Verify order status updated
   - Verify notifications sent
   - Verify trace ID logged

2. **SLA Breach Handling**
   - Create dispute
   - Wait past SLA deadline
   - Run SLA check
   - Verify SLA marked as breached
   - Verify alert sent

3. **Assignment Rollback**
   - Assign dispute
   - Rollback assignment
   - Verify status back to pending
   - Verify can reassign

4. **Validation Tests**
   - Cannot initiate outside 24h window
   - Cannot initiate duplicate dispute
   - Cannot resolve non-assigned dispute
   - Cannot assign to non-existent user

## Frontend Implementation

### Admin Assignment Workbench

Location: `frontend/src/pages/admin/AssignmentWorkbench.tsx`

**Features:**
- List of pending disputes
- Real-time SLA countdown
- Timeline view of dispute history
- Quick assign/rollback actions
- Dispute detail modal
- Resolution form

### Components

1. **DisputeList** - Paginated list of disputes
2. **DisputeDetail** - Full dispute information
3. **Timeline** - Chronological view of actions
4. **SLACountdown** - Real-time SLA timer
5. **AssignmentForm** - Assign to CS rep
6. **ResolutionForm** - Resolve dispute

### State Management

Use React Query for:
- Fetching pending disputes
- Polling for SLA updates
- Mutation for assign/resolve/rollback

## Monitoring & Alerts

### Prometheus Rules

Location: `ops/alerts/dispute_sla_rules.yml`

```yaml
groups:
  - name: dispute_sla
    rules:
      - alert: DisputeSLABreached
        expr: order_dispute_sla_breached{status="assigned"} > 0
        for: 1m
        annotations:
          summary: "Dispute {{ $labels.dispute_id }} SLA breached"
          description: "Dispute has exceeded SLA deadline"
      
      - alert: HighPendingDisputeCount
        expr: order_dispute_pending_count > 10
        for: 5m
        annotations:
          summary: "High number of pending disputes"
          description: "{{ $value }} disputes pending assignment"
```

### Loki Queries

```logql
# Find disputes exceeding SLA
{job="gamelink-backend"} | json | action="update_status" AND reason="SLA breached"

# Trace dispute lifecycle
{job="gamelink-backend"} | json | trace_id="<trace-id>"

# Count disputes by status
{job="gamelink-backend"} | json | entity_type="dispute" | stats count() by status
```

## Operation Logging

All dispute operations logged to `operation_log` table:

**Actions:**
- `initiate_dispute` - User creates dispute
- `assign_dispute` - CS rep assigned
- `mediate_dispute` - Dispute being mediated
- `resolve_dispute` - Dispute resolved
- `rollback_dispute` - Assignment rolled back
- `reject_dispute` - Dispute rejected

**Metadata:**
- `trace_id` - Unique request trace
- `actor_user_id` - Who performed action
- `reason` - Why action taken
- `timestamp` - When action occurred

## Implementation Checklist

### Backend
- [x] Data models (OrderDispute, AssignmentSource)
- [x] Repository layer with CRUD operations
- [x] Service layer with business logic
- [x] Admin handlers (list, assign, rollback, resolve)
- [x] User handlers (initiate, get detail)
- [ ] Database migration
- [ ] Integration tests
- [ ] Unit tests for service layer
- [ ] Error handling and validation

### Frontend
- [ ] Assignment workbench page
- [ ] Dispute list component
- [ ] Dispute detail modal
- [ ] Timeline component
- [ ] SLA countdown timer
- [ ] Assignment form
- [ ] Resolution form
- [ ] API integration with React Query
- [ ] Error handling and notifications

### DevOps
- [ ] Prometheus alert rules
- [ ] Loki query templates
- [ ] Deployment configuration
- [ ] Database migration scripts

### Documentation
- [ ] README API section update
- [ ] Swagger/OpenAPI specs
- [ ] Postman collection
- [ ] Architecture diagram
- [ ] Runbook for operations

## Testing Commands

```bash
# Run backend tests
cd backend
make test

# Run specific test
go test -v ./internal/service/assignment -run TestInitiateDispute

# Run with coverage
go test -v -coverprofile=coverage.out ./internal/service/assignment
go tool cover -html=coverage.out

# Run integration tests
go test -v -tags=integration ./internal/service/assignment

# Frontend tests
cd frontend
npm run test

# Lint checks
make lint
npm run lint
```

## Deployment Considerations

1. **Database Migration** - Run migration before deploying new code
2. **Backward Compatibility** - Existing orders unaffected
3. **Feature Flags** - Consider feature flag for gradual rollout
4. **Monitoring** - Enable alerts before going live
5. **Runbook** - Document SLA check scheduler setup

## Future Enhancements

1. **AI-Powered Recommendations** - ML model for assignment suggestions
2. **Bulk Operations** - Batch assign/resolve disputes
3. **Custom SLA Rules** - Per-category SLA times
4. **Escalation Workflow** - Multi-tier escalation
5. **Analytics Dashboard** - Dispute metrics and trends
6. **Webhook Integration** - External system notifications
7. **Mobile App Support** - CS rep mobile app

## References

- [GameLink Architecture](./ARCHITECTURE.md)
- [API Design Standards](./api/api-design-standards.md)
- [Go Coding Standards](./api/go-coding-standards.md)
- [Workflow Distribution Plan](./WORK_DISTRIBUTION_PLAN.md)
