package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/grounded042/capacious/services"
	"github.com/zenazn/goji/web"
)

type List struct {
	Events   EventsController
	Invitees InviteesController
	Auth     AuthController
}

func NewControllersList(coord services.Coordinator) List {
	return List{
		Events:   NewEventsController(coord),
		Invitees: NewInviteesController(coord),
		Auth:     NewAuthController(coord),
	}
}

// checkForAndHandleUserIDInContext checks for the UserID variable in the
// context c. If there is a user id in the context, it returns the user id along
// with true. If there is not, it writes the status 401, the provided mesaage to
// the http.ResponseWriter w, returns an empty string, and false.
func checkForAndHandleUserIDInContext(c web.C, w http.ResponseWriter, message string) (string, bool) {
	userID, ok := c.Env["UserID"].(string)

	if !ok || userID == "" {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(message)
		return "", false
	}

	return userID, true
}
