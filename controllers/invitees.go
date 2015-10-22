package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/grounded042/capacious/entities"
	"github.com/zenazn/goji/web"
)

type InviteeStub interface {
	GetInviteesForEvent(string) ([]entities.Invitee, error)
	CreateInviteeForEvent(*entities.Invitee, entities.Event) error
	GetInviteeFromId(string) (entities.Invitee, error)
	EditInvitee(entities.Invitee) error
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
		fmt.Println(err)
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
