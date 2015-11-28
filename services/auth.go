package services

import (
	"crypto/sha256"
	"encoding/base64"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/pbkdf2"

	"github.com/dgrijalva/jwt-go"
	"github.com/grounded042/capacious/entities"
	"github.com/grounded042/capacious/utils"
)

// LoginUser is used to represent login details for a user. It's an object to
// hold user data for auth.
type LoginUser struct {
	UserID   string
	Email    string
	Password string
	Token    string
}

type authGateway interface {
	// GetUserLoginFromEmail gets a user login object from the db that has the
	// passed in email address
	GetUserLoginFromEmail(string) (entities.UserLogin, error)
}

type authService struct {
	da authGateway
}

func newAuthService(newDa authGateway) authService {
	return authService{
		da: newDa,
	}
}

// Login will authenticate login credentials from the lUser object
func (as authService) Login(lUser LoginUser) (string, utils.Error) {
	lUser.Email = strings.ToLower(lUser.Email)

	// get the userlogin object based on the email
	dbUser, err := as.da.GetUserLoginFromEmail(lUser.Email)
	if err != nil {
		return "", utils.NewApiError(401, "Could not find user.")
	}

	// see if the user login creds are valid
	success := as.authenticate(lUser, dbUser)
	if !success {
		return "", utils.NewApiError(401, "Authentication failed.")
	}

	// generate a token
	token, err := as.GenerateToken(dbUser.FkUserID)
	if err != nil {
		return "", utils.NewApiError(500, err.Error())
	}

	return token, nil
}

// authenticate will hash and then compare the password from the authUser with
// the users hashed password from the database. If they match, the correct
// passwrod for the user has been provided.
func (as authService) authenticate(authUser LoginUser, dbUser entities.UserLogin) bool {
	return as.hashPasswordWithSalt(authUser.Password, dbUser.Salt) == dbUser.Password
}

// hashPasswordWithSalt will hash the passed in password with the passed in salt
// and return a string of the hashed value
func (as authService) hashPasswordWithSalt(password string, salt string) string {
	hashed := pbkdf2.Key([]byte(salt+password), []byte(salt), 4096, sha256.Size, sha256.New)
	return string(base64.StdEncoding.EncodeToString(hashed))
}

// GenerateToken will generate a new token for the provided user id
func (as authService) GenerateToken(userID string) (string, utils.Error) {
	token := jwt.New(jwt.SigningMethodHS512)
	token.Claims["exp"] = time.Now().Add(time.Hour * time.Duration(72)).Unix()
	token.Claims["iat"] = time.Now().Unix()
	token.Claims["sub"] = userID
	tokenString, err := token.SignedString([]byte(os.Getenv("GO_JWT_MIDDLEWARE_KEY")))
	if err != nil {
		return "", utils.NewApiError(500, err.Error())
	}

	return tokenString, nil
}
