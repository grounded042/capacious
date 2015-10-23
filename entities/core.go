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

type Guest struct {
	GuestId   string    `gorm:"primary_key" sql:"DEFAULT:uuid_generate_v1mc()" json:"guest_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Attending bool      `json:"attending"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Invitee struct {
	InviteeId string         `gorm:"primary_key" sql:"DEFAULT:uuid_generate_v1mc()" json:"invitee_id"`
	FkEventId string         `json:"-"`
	FkGuestId string         `json:"-"`
	Email     string         `json:"email"`
	Self      Guest          `json:"self"`
	Guests    []InviteeGuest `json:"guests"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type InviteeGuest struct {
	InviteeGuestId string    `gorm:"primary_key" sql:"DEFAULT:uuid_generate_v1mc()" json:"invitee_guest_id"`
	FkInviteeId    string    `json:"-"`
	FkGuestId      string    `json:"-"`
	Self           Guest     `json:"self"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
