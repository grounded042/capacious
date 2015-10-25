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
			Pattern: "/invitees/:invitee_id/relationships/guests/:guest_id",
			Handler: cl.Invitees.EditInviteeGuest,
		},
		Route{
			Method:  "post",
			Pattern: "/invitees/:invitee_id/relationships/guests",
			Handler: cl.Invitees.CreateInviteeGuest,
		},
	}
}
