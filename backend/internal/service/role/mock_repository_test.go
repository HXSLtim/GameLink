package role

import (
    "context"
    "errors"
    "testing"

    "gamelink/internal/model"
    "gamelink/internal/repository"

    "github.com/stretchr/testify/require"
)

func TestMockRoleRepository_List_Success(t *testing.T) {
    ctx := context.Background()
    m := &MockRoleRepository{}
    roles := []model.RoleModel{{Slug: "admin"}}
    m.On("List", ctx).Return(roles, nil)
    got, err := m.List(ctx)
    require.NoError(t, err)
    require.Equal(t, 1, len(got))
}

func TestMockRoleRepository_List_Error(t *testing.T) {
    ctx := context.Background()
    m := &MockRoleRepository{}
    m.On("List", ctx).Return(nil, errors.New("x"))
    got, err := m.List(ctx)
    require.Error(t, err)
    require.Nil(t, got)
}

func TestMockRoleRepository_Get_Success(t *testing.T) {
    ctx := context.Background()
    m := &MockRoleRepository{}
    r := &model.RoleModel{Slug: "admin"}
    m.On("Get", ctx, uint64(1)).Return(r, nil)
    got, err := m.Get(ctx, 1)
    require.NoError(t, err)
    require.NotNil(t, got)
}

func TestMockRoleRepository_Get_Error(t *testing.T) {
    ctx := context.Background()
    m := &MockRoleRepository{}
    m.On("Get", ctx, uint64(1)).Return(nil, repository.ErrNotFound)
    got, err := m.Get(ctx, 1)
    require.Error(t, err)
    require.Nil(t, got)
}

func TestMockRoleRepository_Create_Update_Delete(t *testing.T) {
    ctx := context.Background()
    m := &MockRoleRepository{}
    r := &model.RoleModel{Slug: "user"}
    m.On("Create", ctx, r).Return(nil)
    m.On("Update", ctx, r).Return(nil)
    m.On("Delete", ctx, uint64(2)).Return(nil)
    require.NoError(t, m.Create(ctx, r))
    require.NoError(t, m.Update(ctx, r))
    require.NoError(t, m.Delete(ctx, 2))
}

func TestMockRoleRepository_AssignAddRemovePermissions(t *testing.T) {
    ctx := context.Background()
    m := &MockRoleRepository{}
    ids := []uint64{1, 2}
    m.On("AssignPermissions", ctx, uint64(3), ids).Return(nil)
    m.On("AddPermissions", ctx, uint64(3), ids).Return(nil)
    m.On("RemovePermissions", ctx, uint64(3), ids).Return(nil)
    require.NoError(t, m.AssignPermissions(ctx, 3, ids))
    require.NoError(t, m.AddPermissions(ctx, 3, ids))
    require.NoError(t, m.RemovePermissions(ctx, 3, ids))
}

func TestMockRoleRepository_AssignRemoveUserRoles(t *testing.T) {
    ctx := context.Background()
    m := &MockRoleRepository{}
    ids := []uint64{1}
    m.On("AssignToUser", ctx, uint64(9), ids).Return(nil)
    m.On("RemoveFromUser", ctx, uint64(9), ids).Return(nil)
    require.NoError(t, m.AssignToUser(ctx, 9, ids))
    require.NoError(t, m.RemoveFromUser(ctx, 9, ids))
}

func TestMockRoleRepository_ListByUserID_And_CheckUserHasRole(t *testing.T) {
    ctx := context.Background()
    m := &MockRoleRepository{}
    roles := []model.RoleModel{{Slug: "user"}}
    m.On("ListByUserID", ctx, uint64(7)).Return(roles, nil)
    m.On("CheckUserHasRole", ctx, uint64(7), "admin").Return(true, nil)
    got, err := m.ListByUserID(ctx, 7)
    require.NoError(t, err)
    require.Equal(t, 1, len(got))
    ok, err := m.CheckUserHasRole(ctx, 7, "admin")
    require.NoError(t, err)
    require.True(t, ok)
}
