package types

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID                 uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Username           string    `gorm:"not null"`
	Password           string    `gorm:"not null"`
	Email              string    `gorm:"not null"`
	TeamPermissions    []TeamPermission
	ProjectPermissions []ProjectPermission
	RefreshToken       *string
}

type Team struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name     string
	Users    []User `gorm:"many2many:team_permissions;"`
	color    string
	Projects []Project `gorm:"foreignKey:TeamID"`
}

type TeamPermission struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Permission string
	TeamID     uuid.UUID
	UserID     uuid.UUID
}

type ProjectPermission struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Permission string
	ProjectID  uuid.UUID
	UserID     uuid.UUID
}

type Project struct {
	gorm.Model
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	TeamID uuid.UUID
	Name   string
	Color  string
}

type DirectMessage struct {
	gorm.Model
	ID         uuid.UUID
	Message    string
	SenderID   uuid.UUID
	ReceiverID uuid.UUID
}

type TeamMessages struct {
	gorm.Model
	ID       uuid.UUID
	TeamID   string
	Message  string
	SenderID string
}

type Task struct {
	gorm.Model
	ID           uuid.UUID
	Content      string
	Dependencies []Task `gorm:"many2many:task_dependencies;"`
}
