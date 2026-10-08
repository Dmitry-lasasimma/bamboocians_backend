package models

import "time"

// Role constants
const (
	RoleOrganizer = "organizer"
	RoleTalent    = "talent"
	RoleVenue     = "venue"
	RoleAdmin     = "admin"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Role      string    `json:"role" gorm:"not null"` // organizer | talent | venue | admin
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrganizerProfile struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	UserID      uint   `json:"user_id" gorm:"unique;not null"`
	User        User   `json:"user" gorm:"foreignKey:UserID"`
	CompanyName string `json:"company_name"`
	Bio         string `json:"bio"`
	Location    string `json:"location"`
	Phone       string `json:"phone"`
	AvatarURL   string `json:"avatar_url"`
}

type TalentProfile struct {
	ID         uint    `json:"id" gorm:"primaryKey"`
	UserID     uint    `json:"user_id" gorm:"unique;not null"`
	User       User    `json:"user" gorm:"foreignKey:UserID"`
	StageName  string  `json:"stage_name"`
	Genre      string  `json:"genre"`
	Bio        string  `json:"bio"`
	HourlyRate float64 `json:"hourly_rate"`
	Location   string  `json:"location"`
	Skills     string  `json:"skills"` // comma-separated
	AvatarURL  string  `json:"avatar_url"`
}

type VenueProfile struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	UserID      uint    `json:"user_id" gorm:"unique;not null"`
	User        User    `json:"user" gorm:"foreignKey:UserID"`
	VenueName   string  `json:"venue_name"`
	Description string  `json:"description"`
	Capacity    int     `json:"capacity"`
	Location    string  `json:"location"`
	Amenities   string  `json:"amenities"` // comma-separated
	HourlyRate  float64 `json:"hourly_rate"`
	AvatarURL   string  `json:"avatar_url"`
}
