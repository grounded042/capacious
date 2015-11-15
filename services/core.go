package services

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"os"

	"github.com/grounded042/capacious/dal"
	"github.com/grounded042/capacious/entities"
	"github.com/grounded042/capacious/utils"
)

// commonIV for data encryption
var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
var key = os.Getenv("ENC_KEY")

// the Coordinator coordinates interactions between different services.
// it applies any overarching business logic relates to more than one service.
// this allows each service to only care about it's objects and eliminates
// the possibility of circular dependencies in the services

type Coordinator struct {
	events   eventsService
	invitees inviteeService
}

func NewCoordinator(newDa dal.DataHandler) Coordinator {
	return Coordinator{
		events:   newEventsService(newDa),
		invitees: newInviteeService(newDa),
	}
}

// events coordination

func (c Coordinator) GetEvents() ([]entities.Event, utils.Error) {
	return c.events.GetEvents()
}

func (c Coordinator) GetEventInfo(eventId string) (entities.Event, utils.Error) {
	return c.events.GetEventInfo(eventId)
}

func (c Coordinator) CreateEvent(event *entities.Event) utils.Error {
	return c.events.CreateEvent(event)
}

func (c Coordinator) GetMenuItemsForEvent(eventID string) ([]entities.MenuItem, utils.Error) {
	return c.events.GetMenuItemsForEvent(eventID)
}

func (c Coordinator) GetListOfSeatingRequestChoices(eventID string) ([]entities.SeatingRequestChoice, utils.Error) {
	iList, err := c.invitees.GetSeatingRequestInviteesForEvent(eventID)

	if err != nil {
		return []entities.SeatingRequestChoice{}, utils.NewApiError(500, err.Error())
	}

	return c.encryptInviteesToSeatingRequestChoiceList(iList)
}

func (c Coordinator) encryptInviteesToSeatingRequestChoiceList(iList []entities.Invitee) ([]entities.SeatingRequestChoice, utils.Error) {
	var srcl []entities.SeatingRequestChoice

	for _, value := range iList {
		src := entities.SeatingRequestChoice{
			FkInviteeRequestId: value.InviteeId,
			FirstName:          value.Self.FirstName,
			LastName:           value.Self.LastName,
		}

		eSrc, err := c.encryptSeatingRequestChoice(src)

		if err != nil {
			return []entities.SeatingRequestChoice{}, utils.NewApiError(500, err.Error())
		}

		srcl = append(srcl, eSrc)
	}

	return srcl, nil
}

func (c Coordinator) encryptSeatingRequestChoice(choice entities.SeatingRequestChoice) (entities.SeatingRequestChoice, utils.Error) {
	var err error

	choice.FkInviteeRequestId, err = c.encryptFkInviteeRequestId(choice.FkInviteeRequestId)

	if err != nil {
		return entities.SeatingRequestChoice{}, utils.NewApiError(500, err.Error())
	}

	return choice, nil
}

func (c Coordinator) encryptFkInviteeRequestId(toEncrypt string) (string, utils.Error) {
	teByte := []byte(toEncrypt)

	newCipher, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", utils.NewApiError(500, err.Error())
	}

	cfb := cipher.NewCFBEncrypter(newCipher, commonIV)
	ciphertext := make([]byte, len(teByte))
	cfb.XORKeyStream(ciphertext, teByte)

	base64Text := make([]byte, base64.StdEncoding.EncodedLen(len(ciphertext)))
	base64.StdEncoding.Encode(base64Text, []byte(ciphertext))

	return string(base64Text), nil
}

// end events coordination

// invitee coordination

func (c Coordinator) GetInviteesForEvent(eventId string) ([]entities.Invitee, utils.Error) {
	return c.invitees.GetInviteesForEvent(eventId)
}

func (c Coordinator) CreateInviteeForEvent(invitee *entities.Invitee, event entities.Event) utils.Error {
	return c.invitees.CreateInviteeForEvent(invitee, event)
}

func (c Coordinator) GetInviteeFromId(id string) (entities.Invitee, utils.Error) {
	return c.invitees.GetInviteeFromId(id)
}

func (c Coordinator) EditInvitee(updateMe entities.Invitee) utils.Error {
	return c.invitees.EditInvitee(updateMe)
}

func (c Coordinator) CreateInviteeFriend(updateMe *entities.InviteeFriend) utils.Error {
	// TODO: make sure to constrain the number of friends here
	return c.invitees.CreateInviteeFriend(updateMe)
}

func (c Coordinator) EditInviteeFriend(updateMe entities.InviteeFriend) utils.Error {
	// TODO: make sure to constrain the number of friends here
	return c.invitees.EditInviteeFriend(updateMe)
}

func (c Coordinator) SetInviteeMenuChoices(inviteeID string, choices []entities.MenuChoice) ([]entities.MenuChoice, utils.Error) {
	invitee, err := c.invitees.GetInviteeFromId(inviteeID)

	if err != nil {
		return []entities.MenuChoice{}, err
	}

	return c.SetGuestMenuChoices(invitee.FkEventId, invitee.Self.GuestId, choices)
}

func (c Coordinator) SetInviteeFriendMenuChoices(iFriendID string, choices []entities.MenuChoice) ([]entities.MenuChoice, utils.Error) {
	iFriend, err := c.invitees.GetInviteeFriendFromId(iFriendID)

	if err != nil {
		return []entities.MenuChoice{}, err
	}

	invitee, err := c.invitees.GetInviteeFromId(iFriend.FkInviteeId)

	if err != nil {
		return []entities.MenuChoice{}, err
	}

	return c.SetGuestMenuChoices(invitee.FkEventId, iFriend.FkGuestId, choices)
}

func (c Coordinator) SetGuestMenuChoices(eventID string, guestID string, choices []entities.MenuChoice) ([]entities.MenuChoice, utils.Error) {
	items, err := c.events.GetMenuItemsForEvent(eventID)

	if err != nil {
		return []entities.MenuChoice{}, err
	}

	if !c.validateMenuChoicesWithMenuItems(choices, items) {
		return []entities.MenuChoice{}, utils.NewApiError(400, "So yeah...you had an error in the list of menu choices you sent in. Rejected!")
	}

	return c.invitees.SetGuestMenuChoices(guestID, choices)
}

func (c Coordinator) SetInviteeMenuNote(inviteeID string, note entities.MenuNote) (entities.MenuNote, utils.Error) {
	invitee, err := c.invitees.GetInviteeFromId(inviteeID)

	if err != nil {
		return entities.MenuNote{}, err
	}

	return c.SetGuestMenuNote(invitee.FkGuestId, note)
}

func (c Coordinator) SetInviteeFriendMenuNote(iFriendID string, note entities.MenuNote) (entities.MenuNote, utils.Error) {
	iFriend, err := c.invitees.GetInviteeFriendFromId(iFriendID)

	if err != nil {
		return entities.MenuNote{}, err
	}

	return c.SetGuestMenuNote(iFriend.FkGuestId, note)
}

func (c Coordinator) SetGuestMenuNote(guestID string, note entities.MenuNote) (entities.MenuNote, utils.Error) {
	return c.invitees.SetGuestMenuNote(guestID, note)
}

func (c Coordinator) SetInviteeSeatingRequests(inviteeID string, requests []entities.InviteeSeatingRequest) ([]entities.InviteeSeatingRequest, utils.Error) {
	requests, err := c.decryptInviteeSeatingRequests(requests)

	if err != nil {
		return []entities.InviteeSeatingRequest{}, err
	}

	requests, err = c.invitees.SetInviteeSeatingRequests(inviteeID, requests)

	if err != nil {
		return []entities.InviteeSeatingRequest{}, err
	}

	for key, value := range requests {
		requests[key].FkInviteeRequestId, err = c.encryptFkInviteeRequestId(value.FkInviteeRequestId)

		if err != nil {
			return []entities.InviteeSeatingRequest{}, err
		}
	}

	return requests, nil
}

// validateMenuChoicesWithMenuItems validates that the supplied choices match
// up with the supplied menu items. It returns a bool regarding the validity.
// TODO: unit test this sucker
func (c Coordinator) validateMenuChoicesWithMenuItems(choices []entities.MenuChoice, items []entities.MenuItem) bool {
	itemUsage := make(map[string]int)

	for _, cValue := range choices {
		itemMatch := false
		// what item does the current choice align to?
		for _, iValue := range items {
			for _, oValue := range iValue.Options {
				if oValue.MenuItemOptionId == cValue.FkMenuItemOptionId {
					// we found an item option that matches the choice id,
					// so we have validated that the choice is valid
					itemMatch = true

					// up the count for the number of times a choice has been used for an item
					itemUsage[iValue.MenuItemId] = itemUsage[iValue.MenuItemId] + 1

					// is the number of choices higher than the number of choices allowed?
					if itemUsage[iValue.MenuItemId] > iValue.NumChoices {
						return false
					}

				}
			}
		}

		// check and see if the choice matched up with an option in an item
		if !itemMatch {
			return false
		}

	}

	// if we made it this far, we are good!
	return true
}

func (c Coordinator) decryptInviteeSeatingRequests(requests []entities.InviteeSeatingRequest) ([]entities.InviteeSeatingRequest, utils.Error) {
	for key, value := range requests {
		newValue, err := c.decryptInviteeSeatingRequest(value)

		if err != nil {
			return []entities.InviteeSeatingRequest{}, err
		}

		requests[key] = newValue
	}

	return requests, nil
}

func (c Coordinator) decryptInviteeSeatingRequest(request entities.InviteeSeatingRequest) (entities.InviteeSeatingRequest, utils.Error) {
	dbuf := make([]byte, base64.StdEncoding.DecodedLen(len(request.FkInviteeRequestId)))
	base64.StdEncoding.Decode(dbuf, []byte(request.FkInviteeRequestId))
	toDecrypt := []byte(dbuf)

	newCipher, err := aes.NewCipher([]byte(key))
	if err != nil {
		return entities.InviteeSeatingRequest{}, utils.NewApiError(500, err.Error())
	}

	cfbdec := cipher.NewCFBDecrypter(newCipher, commonIV)
	decrypted := make([]byte, len(toDecrypt))
	cfbdec.XORKeyStream(decrypted, toDecrypt)

	request.FkInviteeRequestId = string(decrypted)

	return request, nil
}

//end invitee coordination
