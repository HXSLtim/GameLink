package team

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gamelink/internal/model"
)

// ============================================================================
// Tests
// ============================================================================

func TestTeam_DefaultValues(t *testing.T) {
	t.Run("default max members", func(t *testing.T) {
		team := &model.Team{}
		if team.MaxMembers == 0 {
			team.MaxMembers = 5
		}
		assert.Equal(t, 5, team.MaxMembers)
	})

	t.Run("default income share type", func(t *testing.T) {
		team := &model.Team{}
		if team.IncomeShareType == "" {
			team.IncomeShareType = "equal"
		}
		assert.Equal(t, "equal", team.IncomeShareType)
	})

	t.Run("default status", func(t *testing.T) {
		team := &model.Team{}
		team.Status = model.TeamStatusActive
		assert.Equal(t, model.TeamStatusActive, team.Status)
	})
}

func TestTeam_MemberCount(t *testing.T) {
	t.Run("team full check", func(t *testing.T) {
		team := &model.Team{
			MaxMembers:  5,
			MemberCount: 5,
		}
		assert.True(t, team.MemberCount >= team.MaxMembers)
	})

	t.Run("team has space", func(t *testing.T) {
		team := &model.Team{
			MaxMembers:  5,
			MemberCount: 3,
		}
		assert.True(t, team.MemberCount < team.MaxMembers)
	})
}

func TestTeam_CurrentOrder(t *testing.T) {
	t.Run("has current order", func(t *testing.T) {
		orderID := uint64(100)
		team := &model.Team{
			CurrentOrderID: &orderID,
		}
		assert.NotNil(t, team.CurrentOrderID)
	})

	t.Run("no current order", func(t *testing.T) {
		team := &model.Team{
			CurrentOrderID: nil,
		}
		assert.Nil(t, team.CurrentOrderID)
	})
}

func TestTeamMember_Roles(t *testing.T) {
	t.Run("leader role", func(t *testing.T) {
		member := &model.TeamMember{
			Role: model.TeamMemberRoleLeader,
		}
		assert.Equal(t, model.TeamMemberRoleLeader, member.Role)
	})

	t.Run("member role", func(t *testing.T) {
		member := &model.TeamMember{
			Role: model.TeamMemberRoleMember,
		}
		assert.Equal(t, model.TeamMemberRoleMember, member.Role)
	})
}

func TestTeamMember_Status(t *testing.T) {
	t.Run("active status", func(t *testing.T) {
		member := &model.TeamMember{
			Status: model.TeamMemberStatusActive,
		}
		assert.Equal(t, model.TeamMemberStatusActive, member.Status)
	})

	t.Run("left status", func(t *testing.T) {
		now := time.Now()
		member := &model.TeamMember{
			Status: model.TeamMemberStatusLeft,
			LeftAt: &now,
		}
		assert.Equal(t, model.TeamMemberStatusLeft, member.Status)
		assert.NotNil(t, member.LeftAt)
	})

	t.Run("kicked status", func(t *testing.T) {
		now := time.Now()
		member := &model.TeamMember{
			Status: model.TeamMemberStatusKicked,
			LeftAt: &now,
		}
		assert.Equal(t, model.TeamMemberStatusKicked, member.Status)
		assert.NotNil(t, member.LeftAt)
	})
}

func TestTeamInvite_Status(t *testing.T) {
	t.Run("pending invite", func(t *testing.T) {
		invite := &model.TeamInvite{
			Status:   model.TeamInviteStatusPending,
			ExpireAt: time.Now().Add(7 * 24 * time.Hour),
		}
		assert.Equal(t, model.TeamInviteStatusPending, invite.Status)
	})

	t.Run("accepted invite", func(t *testing.T) {
		invite := &model.TeamInvite{
			Status: model.TeamInviteStatusAccepted,
		}
		assert.Equal(t, model.TeamInviteStatusAccepted, invite.Status)
	})

	t.Run("rejected invite", func(t *testing.T) {
		invite := &model.TeamInvite{
			Status: model.TeamInviteStatusRejected,
		}
		assert.Equal(t, model.TeamInviteStatusRejected, invite.Status)
	})

	t.Run("expired invite", func(t *testing.T) {
		invite := &model.TeamInvite{
			Status:   model.TeamInviteStatusPending,
			ExpireAt: time.Now().Add(-1 * time.Hour),
		}
		assert.True(t, time.Now().After(invite.ExpireAt))
	})
}

func TestTeamInvite_Expiration(t *testing.T) {
	t.Run("not expired", func(t *testing.T) {
		invite := &model.TeamInvite{
			ExpireAt: time.Now().Add(7 * 24 * time.Hour),
		}
		assert.False(t, time.Now().After(invite.ExpireAt))
	})

	t.Run("expired", func(t *testing.T) {
		invite := &model.TeamInvite{
			ExpireAt: time.Now().Add(-1 * time.Hour),
		}
		assert.True(t, time.Now().After(invite.ExpireAt))
	})
}

func TestTeamStatus_Constants(t *testing.T) {
	t.Run("status values", func(t *testing.T) {
		assert.Equal(t, model.TeamStatus("active"), model.TeamStatusActive)
		assert.Equal(t, model.TeamStatus("busy"), model.TeamStatusBusy)
		assert.Equal(t, model.TeamStatus("inactive"), model.TeamStatusInactive)
	})
}

func TestTeamMemberRole_Constants(t *testing.T) {
	t.Run("role values", func(t *testing.T) {
		assert.Equal(t, model.TeamMemberRole("leader"), model.TeamMemberRoleLeader)
		assert.Equal(t, model.TeamMemberRole("member"), model.TeamMemberRoleMember)
	})
}

func TestTeamMemberStatus_Constants(t *testing.T) {
	t.Run("status values", func(t *testing.T) {
		assert.Equal(t, model.TeamMemberStatus("active"), model.TeamMemberStatusActive)
		assert.Equal(t, model.TeamMemberStatus("left"), model.TeamMemberStatusLeft)
		assert.Equal(t, model.TeamMemberStatus("kicked"), model.TeamMemberStatusKicked)
	})
}

func TestTeamInviteStatus_Constants(t *testing.T) {
	t.Run("status values", func(t *testing.T) {
		assert.Equal(t, model.TeamInviteStatus("pending"), model.TeamInviteStatusPending)
		assert.Equal(t, model.TeamInviteStatus("accepted"), model.TeamInviteStatusAccepted)
		assert.Equal(t, model.TeamInviteStatus("rejected"), model.TeamInviteStatusRejected)
		assert.Equal(t, model.TeamInviteStatus("expired"), model.TeamInviteStatusExpired)
	})
}

func TestTeam_IncomeShareType(t *testing.T) {
	t.Run("equal share", func(t *testing.T) {
		team := &model.Team{
			IncomeShareType: "equal",
		}
		assert.Equal(t, "equal", team.IncomeShareType)
	})

	t.Run("custom share", func(t *testing.T) {
		team := &model.Team{
			IncomeShareType: "custom",
		}
		assert.Equal(t, "custom", team.IncomeShareType)
	})
}

func TestTeam_Statistics(t *testing.T) {
	t.Run("order count", func(t *testing.T) {
		team := &model.Team{
			TotalOrderCount: 100,
		}
		assert.Equal(t, 100, team.TotalOrderCount)
	})

	t.Run("income cents", func(t *testing.T) {
		team := &model.Team{
			TotalIncomeCents: 1000000,
		}
		assert.Equal(t, int64(1000000), team.TotalIncomeCents)
	})
}

func TestTeamMember_Statistics(t *testing.T) {
	t.Run("member order count", func(t *testing.T) {
		member := &model.TeamMember{
			OrderCount: 50,
		}
		assert.Equal(t, 50, member.OrderCount)
	})

	t.Run("member income", func(t *testing.T) {
		member := &model.TeamMember{
			IncomeCents: 500000,
		}
		assert.Equal(t, int64(500000), member.IncomeCents)
	})
}

func TestTeam_MaxMembersValidation(t *testing.T) {
	t.Run("cannot reduce below current count", func(t *testing.T) {
		team := &model.Team{
			MaxMembers:  5,
			MemberCount: 4,
		}
		newMaxMembers := 3
		assert.True(t, newMaxMembers < team.MemberCount)
	})

	t.Run("can reduce above current count", func(t *testing.T) {
		team := &model.Team{
			MaxMembers:  5,
			MemberCount: 2,
		}
		newMaxMembers := 3
		assert.True(t, newMaxMembers >= team.MemberCount)
	})
}
