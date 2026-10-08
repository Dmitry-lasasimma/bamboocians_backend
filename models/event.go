package models

import "time"

// EventStatus constants
const (
	EventStatusDraft     = "draft"
	EventStatusPublished = "published"
	EventStatusCompleted = "completed"
	EventStatusCancelled = "cancelled"
)

// EventType constants
const (
	EventTypeWedding   = "wedding"
	EventTypeCorporate = "corporate"
	EventTypeFestival  = "festival"
	EventTypeBirthday  = "birthday"
	EventTypeConcert   = "concert"
	EventTypeOther     = "other"
)

type Event struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	OrganizerID uint      `json:"organizer_id" gorm:"not null"`
	Organizer   User      `json:"organizer" gorm:"foreignKey:OrganizerID"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	EventType   string    `json:"event_type"`
	Date        time.Time `json:"date"`
	Location    string    `json:"location"`
	Budget      float64   `json:"budget"`
	Status      string    `json:"status" gorm:"default:'draft'"`
	Capacity    int       `json:"capacity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
