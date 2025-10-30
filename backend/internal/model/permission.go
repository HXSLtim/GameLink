package model

// HTTPMethod defines standard HTTP methods for permission control.
type HTTPMethod string

// HTTPMethod values define standard HTTP methods.
const (
	HTTPMethodGET    HTTPMethod = "GET"
	HTTPMethodPOST   HTTPMethod = "POST"
	HTTPMethodPUT    HTTPMethod = "PUT"
	HTTPMethodPATCH  HTTPMethod = "PATCH"
	HTTPMethodDELETE HTTPMethod = "DELETE"
)

// Permission represents a backend API resource that can be accessed.
// It records method + path as a unique identifier for fine-grained authorization.
type Permission struct {
	Base
	Method      HTTPMethod `json:"method" gorm:"size:16;not null;uniqueIndex:idx_method_path"`
	Path        string     `json:"path" gorm:"size:255;not null;uniqueIndex:idx_method_path"`
	Code        string     `json:"code" gorm:"size:128;uniqueIndex;comment:语义化标识，如 admin.games.read"`
	Group       string     `json:"group" gorm:"size:64;index;comment:API 分组，如 /admin/games"`
	Description string     `json:"description" gorm:"size:255"`
}

// TableName specifies the table name for Permission.
func (Permission) TableName() string {
	return "permissions"
}



