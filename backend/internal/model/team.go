package model

import "time"

// TeamRole represents a member's role within a team.
type TeamRole string

const (
	TeamRoleLeader TeamRole = "leader"
	TeamRoleMember TeamRole = "member"
)

// TeamMemberStatus indicates a member's current status.
type TeamMemberStatus string

const (
	TeamMemberStatusActive   TeamMemberStatus = "active"
	TeamMemberStatusInactive TeamMemberStatus = "inactive"
)

// Team models a cooperative group that can take assignments.
type Team struct {
	Base
	Name        string
	LeaderID    uint64
	Description string
	Status      string
}

// TeamMember represents a user in a team.
type TeamMember struct {
	Base
	TeamID    uint64
	UserID    uint64
	Role      TeamRole
	Status    TeamMemberStatus
	JoinedAt  *time.Time
	UpdatedAt time.Time
}

// TeamAssignmentStatus tracks assignment lifecycle.
type TeamAssignmentStatus string

const (
	TeamAssignmentStatusPending     TeamAssignmentStatus = "pending"
	TeamAssignmentStatusDispatching TeamAssignmentStatus = "dispatching"
	TeamAssignmentStatusReleased    TeamAssignmentStatus = "released"
	TeamAssignmentStatusCompleted   TeamAssignmentStatus = "completed"
)

// TeamOrderAssignment records a team's engagement with an order.
type TeamOrderAssignment struct {
	Base
	OrderID          uint64
	TeamID           uint64
	Status           TeamAssignmentStatus
	DispatchDeadline time.Time
	LockedAt         *time.Time
	ReleasedAt       *time.Time
}

// TeamAssignmentMemberState indicates invitation state.
type TeamAssignmentMemberState string

const (
	TeamAssignmentMemberStatePending   TeamAssignmentMemberState = "pending"
	TeamAssignmentMemberStateConfirmed TeamAssignmentMemberState = "confirmed"
	TeamAssignmentMemberStateDeclined  TeamAssignmentMemberState = "declined"
)

// TeamAssignmentMember represents a member participating in an assignment.
type TeamAssignmentMember struct {
	Base
	AssignmentID uint64
	MemberID     uint64
	State        TeamAssignmentMemberState
	ConfirmedAt  *time.Time
}

// TeamProfitMode describes how payouts are distributed.
type TeamProfitMode string

const (
	TeamProfitModeDefault TeamProfitMode = "default"
	TeamProfitModeCustom  TeamProfitMode = "custom"
)

// TeamPayoutShare indicates an individual member's profit share.
type TeamPayoutShare struct {
	MemberID uint64
	Percent  int
}

// TeamPayoutPlan stores the configured payout strategy for an assignment.
type TeamPayoutPlan struct {
	Base
	AssignmentID uint64
	ProfitMode   TeamProfitMode
	Shares       []TeamPayoutShare `gorm:"-"`
	UpdatedAt    time.Time
}
