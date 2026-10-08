package models

import "time"

// BookingStatus constants
const (
	BookingStatusPending   = "pending"
	BookingStatusConfirmed = "confirmed"
	BookingStatusDeclined  = "declined"
	BookingStatusCancelled = "cancelled"
)

// BookingType constants
const (
	BookingTypeTalent = "talent"
	BookingTypeVenue  = "venue"
)

type Booking struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	EventID     uint      `json:"event_id" gorm:"not null"`
	Event       Event     `json:"event" gorm:"foreignKey:EventID"`
	BookedByID  uint      `json:"booked_by_id" gorm:"not null"` // organizer user ID
	BookedBy    User      `json:"booked_by" gorm:"foreignKey:BookedByID"`
	BookedToID  uint      `json:"booked_to_id" gorm:"not null"` // talent or venue user ID
	BookedTo    User      `json:"booked_to" gorm:"foreignKey:BookedToID"`
	BookingType string    `json:"booking_type" gorm:"not null"` // talent | venue
	Price       float64   `json:"price"`
	Status      string    `json:"status" gorm:"default:'pending'"`
	Notes       string    `json:"notes"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ContractStatus constants
const (
	ContractStatusDraft     = "draft"
	ContractStatusPending   = "pending"
	ContractStatusActive    = "active"
	ContractStatusCompleted = "completed"
)

type Contract struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	BookingID         uint      `json:"booking_id" gorm:"unique;not null"`
	Booking           Booking   `json:"booking" gorm:"foreignKey:BookingID"`
	Content           string    `json:"content"`
	OrganizerSigned   bool      `json:"organizer_signed" gorm:"default:false"`
	TalentVenueSigned bool      `json:"talent_venue_signed" gorm:"default:false"`
	Status            string    `json:"status" gorm:"default:'draft'"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
