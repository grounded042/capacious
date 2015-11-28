package routes

import "github.com/grounded042/capacious/controllers"

func AuthRoutes(cl controllers.List) []Route {
	return []Route{
		Route{
			Method:  "get",
			Pattern: "/token",
			Handler: cl.Auth.RefreshToken,
		},
		Route{
			Method:  "post",
			Pattern: "/token",
			Handler: cl.Auth.Login,
		},
		Route{
			Method:  "delete",
			Pattern: "/token",
			Handler: cl.Auth.Logout,
		},
	}
}
