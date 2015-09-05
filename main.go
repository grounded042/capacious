package main

import (
	"flag"

	"github.com/grounded042/capacious/controllers"
	"github.com/grounded042/capacious/dal"
	"github.com/grounded042/capacious/routes"
	"github.com/grounded042/capacious/services"
	"github.com/zenazn/goji"
)

type appContext struct {
	DA          dal.DataHandler
	Controllers controllers.List
}

func main() {
	var prefix = flag.String("prefix", "/api/v1", "The prefix for all calls.")

	flag.Parse()

	capaciousAPIServer := goji.DefaultMux
	ac := getAppContext()

	routes.BuildRoutes(capaciousAPIServer, routes.EventRoutes(ac.Controllers), *prefix)

	goji.Serve()
}

func getAppContext() appContext {
	da := dal.NewDal()
	sl := services.NewServicesList(da)
	cl := controllers.NewControllersList(sl)

	return appContext{
		Controllers: cl,
	}
}
