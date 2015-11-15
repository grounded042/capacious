package dal

import (
	"errors"
	"fmt"
	"os"

	"github.com/grounded042/capacious/entities"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type DataHandler struct {
	conn *gorm.DB
}

func NewDal() DataHandler {
	psqlURL := "postgres://" + os.Getenv("PSQL_USERNAME") + ":" + os.Getenv("PSQL_SECRET") + "@" + os.Getenv("PSQL_HOSTNAME") + ":" + os.Getenv("PSQL_PORT") + "/" + os.Getenv("PSQL_DB_NAME") + "?sslmode=disable"

	db, err := gorm.Open("postgres", fmt.Sprint(psqlURL))

	if err != nil {
		fmt.Println(err)
	}

	pErr := db.DB().Ping()

	if pErr != nil {
		fmt.Println(pErr)
	}

	return DataHandler{conn: &db}
}

func (dh DataHandler) GetAllEvents() ([]entities.Event, error) {
	var events = []entities.Event{}

	db := dh.conn.Find(&events)

	return events, db.Error
}

func (dh DataHandler) GetEventInfo(eventId string) (entities.Event, error) {
	var event = entities.Event{EventId: eventId}

	db := dh.conn.Find(&event)

	return event, db.Error
}

func (dh DataHandler) CreateEvent(createMe *entities.Event) error {
	db := dh.conn.Create(&createMe)

	return db.Error
}

func (dh DataHandler) GetAllInviteesForEvent(eventId string) ([]entities.Invitee, error) {
	var invitees = []entities.Invitee{}

	db := dh.conn.Where("fk_event_id = ?", eventId).Find(&invitees)

	if db.Error != nil {
		return []entities.Invitee{}, db.Error
	}

	invitees, db.Error = dh.addInviteeSelfToInvitees(invitees)

	if db.Error != nil {
		return []entities.Invitee{}, db.Error
	}

	invitees, db.Error = dh.addInviteeSeatingRequestsToInvitees(invitees)

	if db.Error != nil {
		return []entities.Invitee{}, db.Error
	}

	return dh.addInviteeFriendsToInvitees(invitees)
}

func (dh DataHandler) addInviteeSelfToInvitees(list []entities.Invitee) ([]entities.Invitee, error) {
	for key, value := range list {
		invitee, err := dh.addInviteeSelfToInvitee(value)

		if err != nil {
			return []entities.Invitee{}, err
		}

		list[key] = invitee
	}

	return list, nil
}

func (dh DataHandler) addInviteeSeatingRequestsToInvitees(list []entities.Invitee) ([]entities.Invitee, error) {
	for key, value := range list {
		invitee, err := dh.addInviteeSeatingRequestsToInvitee(value)

		if err != nil {
			return []entities.Invitee{}, err
		}

		list[key] = invitee
	}

	return list, nil
}

func (dh DataHandler) addInviteeSelfToInvitee(invitee entities.Invitee) (entities.Invitee, error) {
	var err error

	invitee.Self, err = dh.getGuestFromId(invitee.FkGuestId)

	if err != nil {
		return entities.Invitee{}, err
	}

	return invitee, nil
}

func (dh DataHandler) addInviteeSeatingRequestsToInvitee(invitee entities.Invitee) (entities.Invitee, error) {
	var err error

	invitee.SeatingRequests, err = dh.getInviteeSeatingRequestsForInviteeID(invitee.InviteeId)

	if err != nil {
		return entities.Invitee{}, err
	}

	// add the first and last names
	for key, value := range invitee.SeatingRequests {
		firstName, lastName, err := dh.getInviteeFirstNameAndLastNameFromId(value.FkInviteeRequestId)

		if err != nil {
			return entities.Invitee{}, err
		}

		invitee.SeatingRequests[key].FirstName = firstName
		invitee.SeatingRequests[key].LastName = lastName
	}

	return invitee, nil
}

func (dh DataHandler) addInviteeFriendsToInvitees(list []entities.Invitee) ([]entities.Invitee, error) {
	for key, value := range list {
		inviteeFriends, err := dh.GetInviteeFriendsFromInviteeId(value.InviteeId)

		if err != nil {
			return []entities.Invitee{}, err
		}

		list[key].Friends = inviteeFriends
	}

	return list, nil
}

func (dh DataHandler) addInviteeFriendsToInvitee(invitee entities.Invitee) (entities.Invitee, error) {
	inviteeFriends, err := dh.GetInviteeFriendsFromInviteeId(invitee.InviteeId)

	if err != nil {
		return entities.Invitee{}, err
	}

	invitee.Friends = inviteeFriends

	return invitee, nil
}

func (dh DataHandler) GetInviteeFriendsFromInviteeId(id string) ([]entities.InviteeFriend, error) {
	var inviteeFriends []entities.InviteeFriend
	var count int

	db := dh.conn.Debug().Table("invitee_friends").Where("fk_invitee_id = ?", id).Find(&inviteeFriends).Count(&count)

	if count == 0 {
		return []entities.InviteeFriend{}, nil
	} else if db.Error != nil {
		return []entities.InviteeFriend{}, db.Error
	}

	for key, value := range inviteeFriends {
		inviteeFriends[key].Self, db.Error = dh.getGuestFromId(value.FkGuestId)

		if db.Error != nil {
			return []entities.InviteeFriend{}, db.Error
		}
	}

	return inviteeFriends, nil
}

func (dh DataHandler) CreateInvitee(createMe *entities.Invitee) error {

	// TODO: check and make sure email doesn't exist yet

	// create the invitee self
	cErr := dh.createGuest(&createMe.Self)

	if cErr != nil {
		return cErr
	}

	// assign the id of self to the foreign key entry
	createMe.FkGuestId = createMe.Self.GuestId

	db := dh.conn.Debug().Create(&createMe)

	if db.Error != nil {
		return db.Error
	}

	for key, value := range createMe.Friends {
		value.FkInviteeId = createMe.InviteeId

		cigErr := dh.CreateInviteeFriend(&value)

		if cigErr != nil {
			return cigErr
		}

		// assign the value so we can get the ids on the obj
		createMe.Friends[key] = value
	}

	return db.Error
}

func (dh DataHandler) createGuest(createMe *entities.Guest) error {
	db := dh.conn.Debug().Create(&createMe)

	return db.Error
}

func (dh DataHandler) CreateInviteeFriend(createMe *entities.InviteeFriend) error {
	// create the invitee friend self
	cErr := dh.createGuest(&createMe.Self)

	if cErr != nil {
		return cErr
	}

	// assign the id of self to the foreign key entry
	createMe.FkGuestId = createMe.Self.GuestId

	db := dh.conn.Debug().Create(&createMe)

	return db.Error
}

func (dh DataHandler) GetInviteeFromId(id string) (entities.Invitee, error) {
	var invitee entities.Invitee

	db := dh.conn.Debug().Where("invitee_id = ?", id).First(&invitee)

	if db.Error != nil {
		return entities.Invitee{}, db.Error
	}

	invitee, db.Error = dh.addInviteeSelfToInvitee(invitee)

	if db.Error != nil {
		return entities.Invitee{}, db.Error
	}

	invitee, db.Error = dh.addInviteeSeatingRequestsToInvitee(invitee)

	if db.Error != nil {
		return entities.Invitee{}, db.Error
	}

	return dh.addInviteeFriendsToInvitee(invitee)
}

func (dh DataHandler) getInviteeFirstNameAndLastNameFromId(id string) (string, string, error) {
	var invitee entities.Invitee

	db := dh.conn.Debug().Where("invitee_id = ?", id).First(&invitee)

	if db.Error != nil {
		return "", "", db.Error
	}

	invitee, db.Error = dh.addInviteeSelfToInvitee(invitee)

	if db.Error != nil {
		return "", "", db.Error
	}

	return invitee.Self.FirstName, invitee.Self.LastName, nil
}

func (dh DataHandler) getGuestFromId(id string) (entities.Guest, error) {
	var guest entities.Guest
	var err error

	db := dh.conn.Debug().Where("guest_id = ?", id).First(&guest)

	// get the guests menu options
	if db.Error != nil {
		return entities.Guest{}, db.Error
	}

	guest.MenuChoices, err = dh.getMenuChoicesForGuestID(guest.GuestId)

	if err != nil {
		return entities.Guest{}, err
	}

	note, err := dh.getMenuNoteForGuestID(guest.GuestId)

	if err != nil {
		return entities.Guest{}, err
	}

	guest.MenuNote = note.NoteBody

	return guest, err
}

// getMenuChoicesForGuestID gets the menu choices associated with the supplied guestID.
// It returns a slice of entities.MenuChoice objs and any error that occured.
func (dh DataHandler) getMenuChoicesForGuestID(guestID string) ([]entities.MenuChoice, error) {
	var choices []entities.MenuChoice
	var count int

	db := dh.conn.Debug().Where("fk_guest_id = ?", guestID).Find(&choices).Count(&count)

	if count == 0 {
		return []entities.MenuChoice{}, nil
	}

	return choices, db.Error
}

func (dh DataHandler) UpdateInvitee(updateMe entities.Invitee) error {
	// get the current invitee to diff
	curInvitee, err := dh.GetInviteeFromId(updateMe.InviteeId)

	if err != nil {
		return err
	}

	// update info from db
	updateMe.FkEventId = curInvitee.FkEventId
	updateMe.FkGuestId = curInvitee.FkGuestId

	// check and make sure the self id is not different
	if updateMe.FkGuestId != updateMe.Self.GuestId ||
		updateMe.Self.GuestId != curInvitee.Self.GuestId {
		return errors.New("bad invitee self id")
	}

	// update the invitee self
	err = dh.updateGuest(updateMe.Self)

	if err != nil {
		return err
	}

	// lastly, update the invitee obj
	db := dh.conn.Debug().Save(updateMe)

	return db.Error
}

func (dh DataHandler) updateGuest(updateMe entities.Guest) error {
	return dh.conn.Debug().Save(updateMe).Error
}

func (dh DataHandler) UpdateInviteeFriend(updateMe entities.InviteeFriend) error {
	curIG, err := dh.GetInviteeFriendFromId(updateMe.InviteeFriendId)

	if err != nil {
		return err
	}

	// update with info from db
	updateMe.FkGuestId = curIG.FkGuestId
	updateMe.FkInviteeId = curIG.FkInviteeId

	// check and make sure the self id is not different
	if updateMe.FkGuestId != updateMe.Self.GuestId {
		return errors.New("bad invitee friend self id")
	}

	// update the invitee friend self
	err = dh.updateGuest(updateMe.Self)

	if err != nil {
		return err
	}

	// lastly, update the invitee obj
	db := dh.conn.Debug().Save(updateMe)

	return db.Error

}

func (dh DataHandler) GetInviteeFriendFromId(id string) (entities.InviteeFriend, error) {
	var friend entities.InviteeFriend
	var count int

	db := dh.conn.Debug().Where("invitee_friend_id = ?", id).First(&friend).Count(&count)

	if count == 0 {
		return entities.InviteeFriend{}, nil
	} else if db.Error != nil {
		return entities.InviteeFriend{}, db.Error
	}

	friend.Self, db.Error = dh.getGuestFromId(friend.FkGuestId)

	return friend, db.Error
}

// GetMenuItemsForEvent gets the menu items for an event bases on the
// event id eventID. It returns a slice of menu items and any errors that
// occured.
func (dh DataHandler) GetMenuItemsForEvent(eventID string) ([]entities.MenuItem, error) {
	var items = []entities.MenuItem{}
	var count int

	db := dh.conn.Where("fk_event_id = ?", eventID).Find(&items).Count(&count)

	if count == 0 {
		return []entities.MenuItem{}, nil
	} else if db.Error != nil {
		return []entities.MenuItem{}, db.Error
	}

	return dh.addMenuItemOptionsToMenuItems(items)
}

// addMenuItemOptionsToMenuItems adds all of the possible options for a
// menu item to that item object in the supplied entities.MenuItem slice.
// It returns a slice of items with the options added and any error that
// occured.
func (dh DataHandler) addMenuItemOptionsToMenuItems(items []entities.MenuItem) ([]entities.MenuItem, error) {
	for key, value := range items {
		newItem, err := dh.addMenuItemOptionToMenuItem(value)

		if err != nil {
			return items, err
		}

		items[key] = newItem
	}

	return items, nil
}

// addMenuItemOptionToMenuItem adds the possible options for a menu item to
// the supplied entities.MenuItem object. It returns the item with the options
// added and any error that occured.
func (dh DataHandler) addMenuItemOptionToMenuItem(item entities.MenuItem) (entities.MenuItem, error) {
	opts, err := dh.getMenuItemOptionsForMenuItemID(item.MenuItemId)

	if err != nil {
		return item, err
	}

	item.Options = opts

	return item, nil
}

// getMenuItemOptionsForMenuItemID gets the menu item options associated with
// the supplied menuItemID. It returns a slice of the entities.MenuItemOptions
// and any error that occured.
func (dh DataHandler) getMenuItemOptionsForMenuItemID(menuItemID string) ([]entities.MenuItemOption, error) {
	var opts []entities.MenuItemOption
	var count int

	db := dh.conn.Debug().Table("menu_item_options").Where("fk_menu_item_id = ?", menuItemID).Find(&opts).Count(&count)

	if count == 0 {
		return []entities.MenuItemOption{}, nil
	}

	return opts, db.Error
}

func (dh DataHandler) SetGuestMenuChoices(guestID string, choices []entities.MenuChoice) ([]entities.MenuChoice, error) {
	// delete all the current choices
	//  get all the current choices
	oldChoices, err := dh.getMenuChoicesForGuestID(guestID)

	if err != nil {
		return []entities.MenuChoice{}, err
	}

	if len(oldChoices) > 0 {
		for _, value := range oldChoices {
			db := dh.conn.Debug().Delete(value)

			if db.Error != nil {
				return []entities.MenuChoice{}, db.Error
			}
		}
	}

	// add the new choices
	for key, value := range choices {
		db := dh.conn.Debug().Create(&value)

		if db.Error != nil {
			return []entities.MenuChoice{}, db.Error
		}

		choices[key] = value
	}

	return choices, nil
}

func (dh DataHandler) SetGuestMenuNote(guestID string, note entities.MenuNote) (entities.MenuNote, error) {
	// delete the current note
	oldNote, err := dh.getMenuNoteForGuestID(guestID)

	if err != nil {
		return entities.MenuNote{}, err
	}

	if oldNote.MenuNoteId != "" {
		db := dh.conn.Debug().Delete(oldNote)

		if db.Error != nil {
			return entities.MenuNote{}, db.Error
		}
	}

	// add the new note
	db := dh.conn.Debug().Create(&note)

	if db.Error != nil {
		return entities.MenuNote{}, db.Error
	}

	return note, nil
}

func (dh DataHandler) getMenuNoteForGuestID(guestID string) (entities.MenuNote, error) {
	var note entities.MenuNote
	var count int

	db := dh.conn.Debug().Where("fk_guest_id = ?", guestID).Find(&note).Count(&count)

	if count == 0 {
		return entities.MenuNote{}, nil
	}

	return note, db.Error
}

func (dh DataHandler) getInviteeSeatingRequestsForInviteeID(inviteeID string) ([]entities.InviteeSeatingRequest, error) {
	var requests []entities.InviteeSeatingRequest
	var count int

	db := dh.conn.Debug().Where("fk_invitee_id = ?", inviteeID).Find(&requests).Count(&count)

	if count == 0 {
		return []entities.InviteeSeatingRequest{}, nil
	}

	return requests, db.Error
}

func (dh DataHandler) SetInviteeSeatingRequests(inviteeID string, requests []entities.InviteeSeatingRequest) ([]entities.InviteeSeatingRequest, error) {
	// delete all the current requests
	oldRequests, err := dh.getInviteeSeatingRequestsForInviteeID(inviteeID)

	if err != nil {
		return []entities.InviteeSeatingRequest{}, err
	}

	if len(oldRequests) > 0 {
		for _, value := range oldRequests {
			db := dh.conn.Debug().Delete(value)

			if db.Error != nil {
				return []entities.InviteeSeatingRequest{}, db.Error
			}
		}
	}

	// add the new requests
	for key, value := range requests {
		db := dh.conn.Debug().Create(&value)

		if db.Error != nil {
			return []entities.InviteeSeatingRequest{}, db.Error
		}

		requests[key] = value
	}

	return requests, nil
}

type getInviteesForRequest struct {
	InviteeID string
	FirstName string
	LastName  string
}

func (dh DataHandler) GetSeatingRequestInviteesForEvent(eventID string) ([]entities.Invitee, error) {
	var getStuff []getInviteesForRequest
	invitees := []entities.Invitee{}

	db := dh.conn.Debug().Table("invitees").Select("invitees.invitee_id, guests.first_name, guests.last_name").Joins("left join guests on guests.guest_id = invitees.fk_guest_id").Where("invitees.fk_event_id = ?", eventID).Scan(&getStuff)

	if db.Error != nil {
		return []entities.Invitee{}, db.Error
	}

	for _, value := range getStuff {
		invitee := entities.Invitee{
			InviteeId: value.InviteeID,
			Self: entities.Guest{
				FirstName: value.FirstName,
				LastName:  value.LastName,
			},
		}

		invitees = append(invitees, invitee)
	}

	return invitees, nil
}
