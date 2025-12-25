package permission

import (
	"gamelink/pkg/apierr"
)

// Permission-related errors
var (
	// ErrPermissionNotFound indicates the requested permission does not exist.
	ErrPermissionNotFound = apierr.NotFound("权限不存在")

	// ErrPermissionCodeExists indicates the permission code already exists.
	ErrPermissionCodeExists = apierr.Conflict("权限码已存在")

	// ErrPermissionCodeInvalid indicates the permission code format is invalid.
	// Valid format: module.resource.action (three dot-separated lowercase segments)
	ErrPermissionCodeInvalid = apierr.BadRequest("权限码格式无效，应为 module.resource.action")

	// ErrPermissionCodeImmutable indicates the permission code cannot be modified after creation.
	ErrPermissionCodeImmutable = apierr.BadRequest("权限码创建后不可修改")

	// ErrPermissionInUse indicates the permission is referenced by roles and cannot be deleted.
	ErrPermissionInUse = apierr.BadRequest("权限被角色引用，无法删除")

	// ErrPermissionIsSystem indicates the permission is a system permission and cannot be deleted.
	ErrPermissionIsSystem = apierr.BadRequest("系统权限不可删除")
)
