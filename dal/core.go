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

	return dh.addGuestsToInvitees(invitees)
}

func (dh DataHandler) addGuestsToInvitees(list []entities.Invitee) ([]entities.Invitee, error) {
	for key, value := range list {
		guests, err := dh.GetGuestsFromInviteeId(value.InviteeId)

		if err != nil {
			return []entities.Invitee{}, err
		}

		list[key].Guests = guests
	}

	return list, nil
}

func (dh DataHandler) GetGuestsFromInviteeId(id string) ([]entities.Guest, error) {
	var guests []entities.Guest

	db := dh.conn.Debug().Table("guests").Where("fk_invitee_id = ?", id).Find(&guests)

	return guests, db.Error
}

func (dh DataHandler) CreateInvitee(createMe *entities.Invitee) error {

	// TODO: check and make sure email doesn't exist yet
	db := dh.conn.Debug().Create(&createMe)

	if db.Error != nil {
		return db.Error
	}

	for _, value := range createMe.Guests {
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

	return invitee, db.Error
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
