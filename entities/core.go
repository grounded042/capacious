package entities

import "time"

type Event struct {
	EventId       string    `gorm:"primary_key" sql:"DEFAULT:uuid_generate_v1mc()" json:"event_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Location      string    `json:"location"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	AllowedGuests int       `json:"allowed_guests"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Guest struct {
	GuestId     string       `gorm:"primary_key" sql:"DEFAULT:uuid_generate_v1mc()" json:"guest_id"`
	FirstName   string       `json:"first_name"`
	LastName    string       `json:"last_name"`
	Attending   bool         `json:"attending"`
	MenuChoices []MenuChoice `json:"menu_choices"`
	MenuNote    string       `sql:"-" json:"menu_note"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type Invitee struct {
	InviteeId       string                  `gorm:"primary_key" sql:"DEFAULT:uuid_generate_v1mc()" json:"invitee_id"`
	FkEventId       string                  `json:"-"`
	FkGuestId       string                  `json:"-"`
	Email           string                  `json:"email"`
	Self            Guest                   `json:"self"`
	Friends         []InviteeFriend         `json:"friends"`
	SeatingRequests []InviteeSeatingRequest `json:"seating_request"`
	CreatedAt       time.Time               `json:"created_at"`
	UpdatedAt       time.Time               `json:"updated_at"`
}

type InviteeFriend struct {
	InviteeFriendId string    `gorm:"primary_key" sql:"DEFAULT:uuid_generate_v1mc()" json:"invitee_friend_id"`
	FkInviteeId     string    `json:"-"`
	FkGuestId       string    `json:"-"`
	Self            Guest     `json:"self"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type MenuItem struct {
	MenuItemId string           `gorm:"primary_key" sql:"DEFAULT:uuid_generate_v1mc()" json:"menu_item_id"`
	FkEventId  string           `json:"-"`
	ItemOrder  int              `json:"item_order"`
	Name       string           `json:"name"`
	NumChoices int              `json:"num_choices"`
	Options    []MenuItemOption `json:"options"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

type MenuItemOption struct {
	MenuItemOptionId string    `gorm:"primary_key" sql:"DEFAULT:uuid_generate_v1mc()" json:"menu_item_option_id"`
	FkMenuItemId     string    `json:"-"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type MenuChoice struct {
	MenuChoiceId       string    `gorm:"primary_key" sql:"DEFAULT:uuid_generate_v1mc()" json:"menu_choice_id"`
	FkGuestId          string    `json:"-"`
	FkMenuItemId       string    `json:"menu_item_id"`
	FkMenuItemOptionId string    `json:"menu_item_option_id"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type MenuNote struct {
	MenuNoteId string    `gorm:"primary_key" sql:"DEFAULT:uuid_generate_v1mc()" json:"menu_note_id"`
	FkGuestId  string    `json:"-"`
	NoteBody   string    `json:"note_body"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type InviteeSeatingRequest struct {
	InviteeSeatingRequestId string    `gorm:"primary_key" sql:"DEFAULT:uuid_generate_v1mc()" json:"invitee_seating_request_id"`
	FkInviteeId             string    `json:"-"`
	FkInviteeRequestId      string    `json:"invitee_request_id"`
	FirstName               string    `sql:"-" json:"first_name"`
	LastName                string    `sql:"-" json:"last_name"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

type SeatingRequestChoice struct {
	FkInviteeRequestId string `json:"invitee_request_id"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
}
