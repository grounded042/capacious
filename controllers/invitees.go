package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/grounded042/capacious/services"
	"github.com/zenazn/goji/web"
)

type InviteesController struct {
	sl services.List
}

func NewInviteesController(newSL services.List) InviteesController {
	return InviteesController{
		sl: newSL,
	}
}

func (ec InviteesController) GetInviteesForEvent(c web.C, w http.ResponseWriter, r *http.Request) {
	invitees, err := ec.sl.Invitees.GetInviteesForEvent(c.URLParams["id"])

	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(invitees)
}
