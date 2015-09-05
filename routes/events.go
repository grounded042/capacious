package routes

import "github.com/grounded042/capacious/controllers"

func EventRoutes(cl controllers.List) []Route {
	return []Route{
		Route{
			Method:  "get",
			Pattern: "/events",
			Handler: cl.Events.GetEvents,
		},
	}
}
