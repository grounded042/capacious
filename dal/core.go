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

	return dh.addDatesToInvitees(invitees)
}

func (dh DataHandler) addDatesToInvitees(list []entities.Invitee) ([]entities.Invitee, error) {
	for key, value := range list {
		date, err := dh.GetDateForInviteeId(value.InviteeId)

		if err != nil {
			return []entities.Invitee{}, err
		}

		list[key].Date = date
	}

	return list, nil
}

func (dh DataHandler) GetDateForInviteeId(id string) (entities.Date, error) {
	var date = entities.Date{FkInviteeId: id}

	db := dh.conn.Find(&date)

	fmt.Println(date)

	return date, db.Error
}
