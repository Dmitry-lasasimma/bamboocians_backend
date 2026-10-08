package models

// GuestStatus constants
const (
	GuestStatusInvited   = "invited"
	GuestStatusConfirmed = "confirmed"
	GuestStatusDeclined  = "declined"
)

type Guest struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	EventID uint   `json:"event_id" gorm:"not null"`
	Event   Event  `json:"event" gorm:"foreignKey:EventID"`
	Name    string `json:"name" gorm:"not null"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Status  string `json:"status" gorm:"default:'invited'"`
	Ticket  string `json:"ticket"` // unique ticket reference code
}

// TicketType constants
const (
	TicketTypeGeneral  = "general"
	TicketTypeVIP      = "vip"
	TicketTypeBackstage = "backstage"
)

type Ticket struct {
	ID       uint    `json:"id" gorm:"primaryKey"`
	EventID  uint    `json:"event_id" gorm:"not null"`
	Event    Event   `json:"event" gorm:"foreignKey:EventID"`
	Type     string  `json:"type" gorm:"default:'general'"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Sold     int     `json:"sold" gorm:"default:0"`
}
