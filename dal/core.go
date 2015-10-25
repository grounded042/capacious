package dal

import (
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

	return dh.addInviteeGuestsToInvitees(invitees)
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

func (dh DataHandler) addInviteeSelfToInvitee(invitee entities.Invitee) (entities.Invitee, error) {
	var err error

	invitee.Self, err = dh.getGuestFromId(invitee.FkGuestId)

	if err != nil {
		return entities.Invitee{}, err
	}

	return invitee, nil
}

func (dh DataHandler) addInviteeGuestsToInvitees(list []entities.Invitee) ([]entities.Invitee, error) {
	for key, value := range list {
		inviteeGuests, err := dh.GetInviteeGuestsFromInviteeId(value.InviteeId)

		if err != nil {
			return []entities.Invitee{}, err
		}

		list[key].Guests = inviteeGuests
	}

	return list, nil
}

func (dh DataHandler) addInviteeGuestsToInvitee(invitee entities.Invitee) (entities.Invitee, error) {
	inviteeGuests, err := dh.GetInviteeGuestsFromInviteeId(invitee.InviteeId)

	if err != nil {
		return entities.Invitee{}, err
	}

	invitee.Guests = inviteeGuests

	return invitee, nil
}

func (dh DataHandler) GetInviteeGuestsFromInviteeId(id string) ([]entities.InviteeGuest, error) {
	var inviteeGuests []entities.InviteeGuest

	db := dh.conn.Debug().Table("invitee_guests").Where("fk_invitee_id = ?", id).Find(&inviteeGuests)

	if db.Error != nil {
		return []entities.InviteeGuest{}, db.Error
	}

	for key, value := range inviteeGuests {
		inviteeGuests[key].Self, db.Error = dh.getGuestFromId(value.FkGuestId)

		if db.Error != nil {
			return []entities.InviteeGuest{}, db.Error
		}
	}

	return inviteeGuests, nil
}

func (dh DataHandler) CreateInvitee(createMe *entities.Invitee) error {

	// TODO: check and make sure email doesn't exist yet
	// TODO: create invitee and then the self of the invitee
	db := dh.conn.Debug().Create(&createMe)

	if db.Error != nil {
		return db.Error
	}

	for _, value := range createMe.Guests {
		// TODO: use a method to create the inviteeguest and their self entry
		value.FkInviteeId = createMe.InviteeId
		db.Create(&value)

		if db.Error != nil {
			return db.Error
		}
	}

	return db.Error
}

func (dh DataHandler) GetInviteeFromId(id string) (entities.Invitee, error) {
	var invitee entities.Invitee

	db := dh.conn.Debug().Where("invitee_id = ?", id).First(&invitee)

	if db.Error != nil {
		return entities.Invitee{}, db.Error
	}

	invitee.Self, db.Error = dh.getGuestFromId(invitee.FkGuestId)

	if db.Error != nil {
		return entities.Invitee{}, db.Error
	}

	return dh.addInviteeGuestsToInvitee(invitee)
}

func (dh DataHandler) getGuestFromId(id string) (entities.Guest, error) {
	var guest entities.Guest

	db := dh.conn.Debug().Where("guest_id = ?", id).First(&guest)

	return guest, db.Error
}

func (dh DataHandler) UpdateInvitee(updateMe entities.Invitee) error {
	// TODO: update all guests too
	invitee, err := dh.GetInviteeFromId(updateMe.InviteeId)

	if err != nil {
		return err
	}

	updateMe.FkEventId = invitee.FkEventId

	db := dh.conn.Debug().Save(updateMe)

	return db.Error
}
