package model

import "time"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	ID           int64     `json:"user_id"`
	Username     string    `json:"username"`
	DisplayName  string    `json:"display_name,omitempty"`
	PINHash      string    `json:"-"`
	Role         Role      `json:"role"`
	TokenVersion int       `json:"token_version"`
	AvatarIcon   string    `json:"avatar_icon"`
	AvatarPath   string    `json:"avatar_path,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserPublic struct {
	ID          int64  `json:"user_id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name,omitempty"`
	Role        Role   `json:"role"`
	AvatarIcon  string `json:"avatar_icon"`
	AvatarPath  string `json:"avatar_path,omitempty"`
}

func (u *User) ToPublic() UserPublic {
	return UserPublic{
		ID:          u.ID,
		Username:    u.Username,
		DisplayName: u.DisplayName,
		Role:        u.Role,
		AvatarIcon:  u.AvatarIcon,
		AvatarPath:  u.AvatarPath,
	}
}
