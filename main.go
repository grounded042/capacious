package main

import (
	"flag"

	"github.com/grounded042/capacious/controllers"
	"github.com/grounded042/capacious/dal"
	"github.com/grounded042/capacious/middleware"
	"github.com/grounded042/capacious/routes"
	"github.com/grounded042/capacious/services"
	"github.com/zenazn/goji"
)

type appContext struct {
	Controllers controllers.List
}

func main() {
	var prefix = flag.String("prefix", "/api/v1", "The prefix for all calls.")

	flag.Parse()

	capaciousAPIServer := goji.DefaultMux
	ac := getAppContext()

	// apply the middleware
	goji.Use(middleware.ContentTypeHeader)
	goji.Use(middleware.CORS)

	routes.BuildRoutes(capaciousAPIServer, routes.EventRoutes(ac.Controllers), *prefix)
	routes.BuildRoutes(capaciousAPIServer, routes.InviteeRoutes(ac.Controllers), *prefix)

	goji.Serve()
}

func getAppContext() appContext {
	da := dal.NewDal()
	co := services.NewCoordinator(da)
	cl := controllers.NewControllersList(co)

	return appContext{
		Controllers: cl,
	}
}
