package entities

import "time"

// Event represents an object that contains details about a specific event.
type Event struct {
	EventID        string    `gorm:"primary_key" sql:"DEFAULT:uuid_generate_v1mc()" json:"event_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Location       string    `json:"location"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	RespondBy      time.Time `json:"respond_by"`
	AllowedFriends int       `json:"allowed_friends"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Guest represents an object that contains details about a specific guest.
// A guest is always referenced from the self attribute of an Invitee or
// InviteeFriend
type Guest struct {
	GuestID     string       `gorm:"primary_key" sql:"DEFAULT:uuid_generate_v1mc()" json:"guest_id"`
	FirstName   string       `json:"first_name"`
	LastName    string       `json:"last_name"`
	Attending   bool         `json:"attending"`
	MenuChoices []MenuChoice `json:"menu_choices"`
	MenuNote    string       `sql:"-" json:"menu_note"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// Invitee represents an object that contains details about a specific invitee.
// An invitee is the person inviteed to an event and contains information
// about the specific guest it maps to, any friends the invitee will be
// bringing, and other invitees the invitee would like to be seated near.
// The Invitee object also holds db keys for the event it relates to as well as
// the guest it relates to.
type Invitee struct {
	InviteeID       string                  `gorm:"primary_key" sql:"DEFAULT:uuid_generate_v1mc()" json:"invitee_id"`
	FkEventID       string                  `json:"-"`
	FkGuestID       string                  `json:"-"`
	Email           string                  `json:"email"`
	Self            Guest                   `json:"self"`
	Friends         []InviteeFriend         `json:"friends"`
	SeatingRequests []InviteeSeatingRequest `json:"seating_request"`
	CreatedAt       time.Time               `json:"created_at"`
	UpdatedAt       time.Time               `json:"updated_at"`
}

// InviteeFriend represents an object that contains details about a specific
// InviteeFriend.
// An InviteeFriend simply holds a guest object for info about the friend.
// The InviteeFriend object also holds db keys for the invitee it relates to as
// well as the guest it relates to.
type InviteeFriend struct {
	InviteeFriendID string    `gorm:"primary_key" sql:"DEFAULT:uuid_generate_v1mc()" json:"invitee_friend_id"`
	FkInviteeID     string    `json:"-"`
	FkGuestID       string    `json:"-"`
	Self            Guest     `json:"self"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// MenuItem represents an object that contains details about a specific menu
// item.
// A MenuItem holds the order in which the item should appear on a list, the
// name of the item, the number of choices that are allowed to be choosen from
// the options, and the options for the menu item.
// The MenuItem object holds a db key for the event it relates to.
type MenuItem struct {
	MenuItemID string           `gorm:"primary_key" sql:"DEFAULT:uuid_generate_v1mc()" json:"menu_item_id"`
	FkEventID  string           `json:"-"`
	ItemOrder  int              `json:"item_order"`
	Name       string           `json:"name"`
	NumChoices int              `json:"num_choices"`
	Options    []MenuItemOption `json:"options"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

// MenuItemOption represents an object that contains details about a specific
// menu item option.
// A MenuItemOption holds the name and the description of the option.
// The MenuITemOption object holds a db kye for the MenuItem it relates to.
type MenuItemOption struct {
	MenuItemOptionID string    `gorm:"primary_key" sql:"DEFAULT:uuid_generate_v1mc()" json:"menu_item_option_id"`
	FkMenuItemID     string    `json:"-"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// MenuChoice represents an object that contains details about a specific menu
// choice choseen by a guest.
// The MenuChoice object is only used as a refernece object as it contains db
// keys for the guest it relates to, the menu item from the event it relates to
// and the choosen option for the aformentioned menu item.
type MenuChoice struct {
	MenuChoiceID       string    `gorm:"primary_key" sql:"DEFAULT:uuid_generate_v1mc()" json:"menu_choice_id"`
	FkGuestID          string    `json:"-"`
	FkMenuItemID       string    `json:"menu_item_id"`
	FkMenuItemOptionID string    `json:"menu_item_option_id"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// MenuNote represents an object that contains details about a specific menu
// note.
// The MenuNote object simply holds the body of a note in regards to menu
// selection at an event.
// The MenuNote object contains a db key for the guest it relates to.
type MenuNote struct {
	MenuNoteID string    `gorm:"primary_key" sql:"DEFAULT:uuid_generate_v1mc()" json:"menu_note_id"`
	FkGuestID  string    `json:"-"`
	NoteBody   string    `json:"note_body"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// InviteeSeatingRequest represents an object that contains details about a
// specific seating request from an invitee.
// For convenience the InviteeSeatingRequest object contains the first name
// and last name of the invitee whom the request is of.
// The InviteeSeatingRequest object contains db keys relating to the invitee
// (FkInviteeID) who is requesting to sit next to another invitee
// (FkInviteeRequestID).
type InviteeSeatingRequest struct {
	InviteeSeatingRequestID string    `gorm:"primary_key" sql:"DEFAULT:uuid_generate_v1mc()" json:"invitee_seating_request_id"`
	FkInviteeID             string    `json:"-"`
	FkInviteeRequestID      string    `json:"invitee_request_id"`
	FirstName               string    `sql:"-" json:"first_name"`
	LastName                string    `sql:"-" json:"last_name"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

// SeatingRequestChoice represents an object that contains details about a
// seating request an invitee can make.
// This is purely used for front-end display and does not relate to a specific
// table in the database.
type SeatingRequestChoice struct {
	FkInviteeRequestID string `json:"invitee_request_id"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
}
