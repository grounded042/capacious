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
	GetInviteeFromId(string) (entities.Invitee, utils.Error)
	EditInvitee(entities.Invitee) utils.Error
	EditInviteeFriend(entities.InviteeFriend) utils.Error
	CreateInviteeFriend(*entities.InviteeFriend) utils.Error
	SetInviteeMenuChoices(string, []entities.MenuChoice) ([]entities.MenuChoice, utils.Error)
	SetInviteeFriendMenuChoices(string, []entities.MenuChoice) ([]entities.MenuChoice, utils.Error)
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

	event := entities.Event{EventId: c.URLParams["id"]}

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
	invitee, err := ic.is.GetInviteeFromId(c.URLParams["id"])

	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(invitee)
	}
}

func (ic InviteesController) EditInvitee(c web.C, w http.ResponseWriter, r *http.Request) {
	invitee := entities.Invitee{InviteeId: c.URLParams["id"]}

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
	iGuest := entities.InviteeFriend{InviteeFriendId: c.URLParams["friend_id"], FkInviteeId: c.URLParams["invitee_id"]}

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
	iGuest := entities.InviteeFriend{FkInviteeId: c.URLParams["invitee_id"]}

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
	guestID := c.URLParams["guest_id"]
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
