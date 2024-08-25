package entity

import "gorm.io/plugin/soft_delete"

// Users represents table users in database
type Users struct {
	ID        string                `gorm:"primary_key;column:id"`
	Username  string                `gorm:"column:username"`
	Email     string                `gorm:"column:email;unique"`
	Password  string                `gorm:"column:password;"`
	CreatedAt int64                 `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64                 `gorm:"column:updated_at;autoUpdateTime:milli"`
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;default:0"`
}

// Used for implement model gorm
func (u Users) TableName() string {
	return "users"
}
