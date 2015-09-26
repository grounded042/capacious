package entities

import "time"

type Event struct {
	EventId       string    `gorm:"primary_key" sql:"DEFAULT:uuid_generate_v1mc()" json:"event_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	AllowedGuests int       `json:"allowed_guests"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Invitee struct {
	InviteeId string    `gorm:"primary_key" sql:"DEFAULT:uuid_generate_v1mc()" json:"invitee_id"`
	FkEventId string    `json:"-"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Guests    []Guest   `json:"guests"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Guest struct {
	GuestId     string    `gorm:"primary_key" sql:"DEFAULT:uuid_generate_v1mc()" json:"date_id"`
	FkInviteeId string    `json:"-"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Attending   bool      `json:"attending"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
