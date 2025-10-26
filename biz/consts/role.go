package consts

// UserRole 用户角色类型
type UserRole string

const (
	// RoleNormal 普通用户
	RoleNormal UserRole = "normal"

	// RoleAdmin 管理员
	RoleAdmin UserRole = "admin"
)

// String 返回角色的字符串表示
func (r UserRole) String() string {
	return string(r)
}

// IsValid 检查角色是否有效
func (r UserRole) IsValid() bool {
	switch r {
	case RoleNormal, RoleAdmin:
		return true
	default:
		return false
	}
}

// IsAdmin 判断是否是管理员
func (r UserRole) IsAdmin() bool {
	return r == RoleAdmin
}

// IsNormal 判断是否是普通用户
func (r UserRole) IsNormal() bool {
	return r == RoleNormal
}

// AllRoles 返回所有角色列表
func AllRoles() []UserRole {
	return []UserRole{RoleNormal, RoleAdmin}
}
