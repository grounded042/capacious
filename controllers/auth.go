package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/grounded042/capacious/services"
	"github.com/grounded042/capacious/utils"
	"github.com/zenazn/goji/web"
)

type jsonTokenObj struct {
	Token string `json:"token"`
}

type AuthStub interface {
	Login(services.LoginUser) (string, utils.Error)
	GenerateToken(string) (string, utils.Error)
}

type AuthController struct {
	as AuthStub
}

func NewAuthController(newAS AuthStub) AuthController {
	return AuthController{
		as: newAS,
	}
}

func (ac AuthController) Login(c web.C, w http.ResponseWriter, r *http.Request) {
	user := new(services.LoginUser)
	decoder := json.NewDecoder(r.Body)

	if dErr := decoder.Decode(&user); dErr != nil {
		w.WriteHeader(400)
		fmt.Println(dErr)
		return
	}

	if token, err := ac.as.Login(*user); err != nil {
		w.WriteHeader(err.Code())
		fmt.Println(err)
	} else if jsonToken, mErr := json.Marshal(jsonTokenObj{Token: token}); mErr != nil {
		w.WriteHeader(500)
		fmt.Println(mErr)
	} else {
		w.Write(jsonToken)
	}
}

func (ac AuthController) RefreshToken(c web.C, w http.ResponseWriter, r *http.Request) {
	userId, ok := c.Env["UserID"].(string)
	if !ok || userId == "" {
		w.WriteHeader(401)
		w.Write([]byte("You need a valid user id to refresh your token!"))
		return
	}

	if token, err := ac.as.GenerateToken(userId); err != nil {
		w.WriteHeader(err.Code())
		fmt.Println(err)
		return
	} else if jsonToken, mErr := json.Marshal(jsonTokenObj{Token: token}); mErr != nil {
		w.WriteHeader(500)
		fmt.Println(mErr)
	} else {
		w.Write(jsonToken)
	}
}

func (ac AuthController) Logout(c web.C, w http.ResponseWriter, r *http.Request) {
	// TODO: implement this
}
