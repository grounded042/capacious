package dal

import (
	"errors"
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

	return dh.addInviteeFriendsToInvitees(invitees)
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

func (dh DataHandler) addInviteeFriendsToInvitees(list []entities.Invitee) ([]entities.Invitee, error) {
	for key, value := range list {
		inviteeFriends, err := dh.GetInviteeFriendsFromInviteeId(value.InviteeId)

		if err != nil {
			return []entities.Invitee{}, err
		}

		list[key].Friends = inviteeFriends
	}

	return list, nil
}

func (dh DataHandler) addInviteeFriendsToInvitee(invitee entities.Invitee) (entities.Invitee, error) {
	inviteeFriends, err := dh.GetInviteeFriendsFromInviteeId(invitee.InviteeId)

	if err != nil {
		return entities.Invitee{}, err
	}

	invitee.Friends = inviteeFriends

	return invitee, nil
}

func (dh DataHandler) GetInviteeFriendsFromInviteeId(id string) ([]entities.InviteeFriend, error) {
	var inviteeFriends []entities.InviteeFriend

	db := dh.conn.Debug().Table("invitee_friends").Where("fk_invitee_id = ?", id).Find(&inviteeFriends)

	if db.Error != nil {
		return []entities.InviteeFriend{}, db.Error
	}

	for key, value := range inviteeFriends {
		inviteeFriends[key].Self, db.Error = dh.getGuestFromId(value.FkGuestId)

		if db.Error != nil {
			return []entities.InviteeFriend{}, db.Error
		}
	}

	return inviteeFriends, nil
}

func (dh DataHandler) CreateInvitee(createMe *entities.Invitee) error {

	// TODO: check and make sure email doesn't exist yet

	// create the invitee self
	cErr := dh.createGuest(&createMe.Self)

	if cErr != nil {
		return cErr
	}

	// assign the id of self to the foreign key entry
	createMe.FkGuestId = createMe.Self.GuestId

	db := dh.conn.Debug().Create(&createMe)

	if db.Error != nil {
		return db.Error
	}

	for key, value := range createMe.Friends {
		value.FkInviteeId = createMe.InviteeId

		cigErr := dh.CreateInviteeFriend(&value)

		if cigErr != nil {
			return cigErr
		}

		// assign the value so we can get the ids on the obj
		createMe.Friends[key] = value
	}

	return db.Error
}

func (dh DataHandler) createGuest(createMe *entities.Guest) error {
	db := dh.conn.Debug().Create(&createMe)

	return db.Error
}

func (dh DataHandler) CreateInviteeFriend(createMe *entities.InviteeFriend) error {
	// create the invitee friend self
	cErr := dh.createGuest(&createMe.Self)

	if cErr != nil {
		return cErr
	}

	// assign the id of self to the foreign key entry
	createMe.FkGuestId = createMe.Self.GuestId

	db := dh.conn.Debug().Create(&createMe)

	return db.Error
}

func (dh DataHandler) GetInviteeFromId(id string) (entities.Invitee, error) {
	var invitee entities.Invitee

	db := dh.conn.Debug().Where("invitee_id = ?", id).First(&invitee)

	if db.Error != nil {
		return entities.Invitee{}, db.Error
	}

	invitee, db.Error = dh.addInviteeSelfToInvitee(invitee)

	if db.Error != nil {
		return entities.Invitee{}, db.Error
	}

	return dh.addInviteeFriendsToInvitee(invitee)
}

func (dh DataHandler) getGuestFromId(id string) (entities.Guest, error) {
	var guest entities.Guest

	db := dh.conn.Debug().Where("guest_id = ?", id).First(&guest)

	return guest, db.Error
}

func (dh DataHandler) UpdateInvitee(updateMe entities.Invitee) error {
	// get the current invitee to diff
	curInvitee, err := dh.GetInviteeFromId(updateMe.InviteeId)

	if err != nil {
		return err
	}

	// update info from db
	updateMe.FkEventId = curInvitee.FkEventId
	updateMe.FkGuestId = curInvitee.FkGuestId

	// check and make sure the self id is not different
	if updateMe.FkGuestId != updateMe.Self.GuestId ||
		updateMe.Self.GuestId != curInvitee.Self.GuestId {
		return errors.New("bad invitee self id")
	}

	// update the invitee self
	err = dh.updateGuest(updateMe.Self)

	if err != nil {
		return err
	}

	// lastly, update the invitee obj
	db := dh.conn.Debug().Save(updateMe)

	return db.Error
}

func (dh DataHandler) updateGuest(updateMe entities.Guest) error {
	return dh.conn.Debug().Save(updateMe).Error
}

func (dh DataHandler) UpdateInviteeFriend(updateMe entities.InviteeFriend) error {
	curIG, err := dh.GetInviteeFriendFromId(updateMe.InviteeFriendId)

	if err != nil {
		return err
	}

	// update with info from db
	updateMe.FkGuestId = curIG.FkGuestId
	updateMe.FkInviteeId = curIG.FkInviteeId

	// check and make sure the self id is not different
	if updateMe.FkGuestId != updateMe.Self.GuestId {
		return errors.New("bad invitee friend self id")
	}

	// update the invitee friend self
	err = dh.updateGuest(updateMe.Self)

	if err != nil {
		return err
	}

	// lastly, update the invitee obj
	db := dh.conn.Debug().Save(updateMe)

	return db.Error

}

func (dh DataHandler) GetInviteeFriendFromId(id string) (entities.InviteeFriend, error) {
	var friend entities.InviteeFriend

	db := dh.conn.Debug().Where("invitee_friend_id = ?", id).First(&friend)

	if db.Error != nil {
		return entities.InviteeFriend{}, db.Error
	}

	friend.Self, db.Error = dh.getGuestFromId(friend.FkGuestId)

	return friend, db.Error
}
