package bo

// UserRoleBo
type UserRoleBo struct {
	UserId   int64  `db:"user_id,omitempty" json:"userId"`
	RoleId   int64  `db:"role_id,omitempty" json:"roleId"`
	RoleName string `db:"role_name,omitempty" json:"roleName"`
}

type RoleUserCount struct {
	RoleID int64 `json:"roleId" db:"role_id"`
	Count uint64 `json:"count" db:"count"`
}
