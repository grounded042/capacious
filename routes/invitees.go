package routes

import "github.com/grounded042/capacious/controllers"

func InviteeRoutes(cl controllers.List) []Route {
	return []Route{
		Route{
			Method:  "get",
			Pattern: "/invitees/:id",
			Handler: cl.Invitees.GetInvitee,
		},
		Route{
			Method:  "patch",
			Pattern: "/invitees/:id",
			Handler: cl.Invitees.EditInvitee,
		},
		Route{
			Method:  "patch",
			Pattern: "/invitees/:invitee_id/relationships/friends/:friend_id",
			Handler: cl.Invitees.EditInviteeFriend,
		},
		Route{
			Method:  "post",
			Pattern: "/invitees/:invitee_id/relationships/friends",
			Handler: cl.Invitees.CreateInviteeFriend,
		},
		Route{
			Method:  "post",
			Pattern: "/invitees/:invitee_id/relationships/menu_choices",
			Handler: cl.Invitees.SetInviteeMenuChoices,
		},
		Route{
			Method:  "post",
			Pattern: "/invitees/:invitee_id/relationships/menu_note",
			Handler: cl.Invitees.SetInviteeMenuNote,
		},
		// TODO: this might need to be moved into a better controller
		// maybe a guest controller
		Route{
			Method:  "post",
			Pattern: "/invitees/:invitee_id/relationships/friends/:friend_id/relationships/menu_choices",
			Handler: cl.Invitees.SetGuestMenuChoices,
		},
		Route{
			Method:  "post",
			Pattern: "/invitees/:invitee_id/relationships/friends/:friend_id/relationships/menu_note",
			Handler: cl.Invitees.SetInviteeFriendMenuNote,
		},
	}
}
