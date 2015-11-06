package routes

import "github.com/grounded042/capacious/controllers"

func EventRoutes(cl controllers.List) []Route {
	return []Route{
		Route{
			Method:  "get",
			Pattern: "/events",
			Handler: cl.Events.GetEvents,
		},
		Route{
			Method:  "post",
			Pattern: "/events",
			Handler: cl.Events.CreateEvent,
		},
		Route{
			Method:  "get",
			Pattern: "/events/:id",
			Handler: cl.Events.GetEventInfo,
		},
		Route{
			Method:  "get",
			Pattern: "/events/:id/relationships/invitees",
			Handler: cl.Invitees.GetInviteesForEvent,
		},
		Route{
			Method:  "post",
			Pattern: "/events/:id/relationships/invitees",
			Handler: cl.Invitees.CreateInviteeForEvent,
		},
	}
}
