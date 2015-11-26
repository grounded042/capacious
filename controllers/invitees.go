package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/grounded042/capacious/entities"
	"github.com/grounded042/capacious/utils"
	"github.com/zenazn/goji/web"
)

type InviteeStub interface {
	GetInviteesForEvent(string) ([]entities.Invitee, utils.Error)
	CreateInviteeForEvent(*entities.Invitee, entities.Event) utils.Error
	GetInviteeFromID(string) (entities.Invitee, utils.Error)
	EditInvitee(entities.Invitee) utils.Error
	EditInviteeFriend(entities.InviteeFriend) utils.Error
	CreateInviteeFriend(*entities.InviteeFriend) utils.Error
	SetInviteeMenuChoices(string, []entities.MenuChoice) ([]entities.MenuChoice, utils.Error)
	SetInviteeFriendMenuChoices(string, []entities.MenuChoice) ([]entities.MenuChoice, utils.Error)
	SetInviteeMenuNote(string, entities.MenuNote) (entities.MenuNote, utils.Error)
	SetInviteeFriendMenuNote(string, entities.MenuNote) (entities.MenuNote, utils.Error)
	SetInviteeSeatingRequests(string, []entities.InviteeSeatingRequest) ([]entities.InviteeSeatingRequest, utils.Error)
}

type InviteesController struct {
	is InviteeStub
}

func NewInviteesController(newIs InviteeStub) InviteesController {
	return InviteesController{
		is: newIs,
	}
}

func (ec InviteesController) GetInviteesForEvent(c web.C, w http.ResponseWriter, r *http.Request) {
	invitees, err := ec.is.GetInviteesForEvent(c.URLParams["id"])

	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err.Error())
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(invitees)
}

func (ec InviteesController) CreateInviteeForEvent(c web.C, w http.ResponseWriter, r *http.Request) {
	var invitee entities.Invitee

	event := entities.Event{EventID: c.URLParams["id"]}

	rBody, ioErr := ioutil.ReadAll(r.Body)

	if ioErr != nil {
		w.WriteHeader(500)
		fmt.Println(ioErr)
		return
	}

	if err := json.Unmarshal(rBody, &invitee); err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
		return
	}

	if err := ec.is.CreateInviteeForEvent(&invitee, event); err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(invitee)
	}
}

func (ic InviteesController) GetInvitee(c web.C, w http.ResponseWriter, r *http.Request) {
	invitee, err := ic.is.GetInviteeFromID(c.URLParams["id"])

	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(invitee)
	}
}

func (ic InviteesController) EditInvitee(c web.C, w http.ResponseWriter, r *http.Request) {
	invitee := entities.Invitee{InviteeID: c.URLParams["id"]}

	rBody, ioErr := ioutil.ReadAll(r.Body)

	if ioErr != nil {
		w.WriteHeader(500)
		fmt.Println(ioErr)
		return
	}

	if err := json.Unmarshal(rBody, &invitee); err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
		return
	}

	err := ic.is.EditInvitee(invitee)

	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(invitee)
	}
}

func (ic InviteesController) EditInviteeFriend(c web.C, w http.ResponseWriter, r *http.Request) {
	iGuest := entities.InviteeFriend{InviteeFriendID: c.URLParams["friend_id"], FkInviteeID: c.URLParams["invitee_id"]}

	rBody, ioErr := ioutil.ReadAll(r.Body)

	if ioErr != nil {
		w.WriteHeader(500)
		fmt.Println(ioErr)
		return
	}

	if err := json.Unmarshal(rBody, &iGuest); err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
		return
	}

	err := ic.is.EditInviteeFriend(iGuest)

	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(iGuest)
	}

}

func (ec InviteesController) CreateInviteeFriend(c web.C, w http.ResponseWriter, r *http.Request) {
	iGuest := entities.InviteeFriend{FkInviteeID: c.URLParams["invitee_id"]}

	rBody, ioErr := ioutil.ReadAll(r.Body)

	if ioErr != nil {
		w.WriteHeader(500)
		fmt.Println(ioErr)
		return
	}

	if err := json.Unmarshal(rBody, &iGuest); err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
		return
	}

	if err := ec.is.CreateInviteeFriend(&iGuest); err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(iGuest)
	}
}

func (ec InviteesController) SetInviteeMenuChoices(c web.C, w http.ResponseWriter, r *http.Request) {
	inviteeID := c.URLParams["invitee_id"]
	var choices []entities.MenuChoice

	rBody, ioErr := ioutil.ReadAll(r.Body)

	if ioErr != nil {
		w.WriteHeader(500)
		fmt.Println(ioErr)
		return
	}

	if err := json.Unmarshal(rBody, &choices); err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
		return
	}

	updatedChoices, err := ec.is.SetInviteeMenuChoices(inviteeID, choices)

	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(updatedChoices)
	}
}

func (ec InviteesController) SetGuestMenuChoices(c web.C, w http.ResponseWriter, r *http.Request) {
	guestID := c.URLParams["friend_id"]
	var choices []entities.MenuChoice

	rBody, ioErr := ioutil.ReadAll(r.Body)

	if ioErr != nil {
		w.WriteHeader(500)
		fmt.Println(ioErr)
		return
	}

	if err := json.Unmarshal(rBody, &choices); err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
		return
	}

	updatedChoices, err := ec.is.SetInviteeFriendMenuChoices(guestID, choices)

	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(updatedChoices)
	}
}

func (ec InviteesController) SetInviteeMenuNote(c web.C, w http.ResponseWriter, r *http.Request) {
	inviteeID := c.URLParams["invitee_id"]
	var note entities.MenuNote

	rBody, ioErr := ioutil.ReadAll(r.Body)

	if ioErr != nil {
		w.WriteHeader(500)
		fmt.Println(ioErr)
		return
	}

	if err := json.Unmarshal(rBody, &note); err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
		return
	}

	updatedNote, err := ec.is.SetInviteeMenuNote(inviteeID, note)

	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(updatedNote)
	}
}

func (ec InviteesController) SetInviteeFriendMenuNote(c web.C, w http.ResponseWriter, r *http.Request) {
	friendID := c.URLParams["friend_id"]
	var note entities.MenuNote

	rBody, ioErr := ioutil.ReadAll(r.Body)

	if ioErr != nil {
		w.WriteHeader(500)
		fmt.Println(ioErr)
		return
	}

	if err := json.Unmarshal(rBody, &note); err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
		return
	}

	updatedNote, err := ec.is.SetInviteeFriendMenuNote(friendID, note)

	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(updatedNote)
	}
}

func (ec InviteesController) SetInviteeSeatingRequests(c web.C, w http.ResponseWriter, r *http.Request) {
	inviteeID := c.URLParams["invitee_id"]
	var requests []entities.InviteeSeatingRequest

	rBody, ioErr := ioutil.ReadAll(r.Body)

	if ioErr != nil {
		w.WriteHeader(500)
		fmt.Println(ioErr)
		return
	}

	if err := json.Unmarshal(rBody, &requests); err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
		return
	}

	updatedRequests, err := ec.is.SetInviteeSeatingRequests(inviteeID, requests)

	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(updatedRequests)
	}
}
