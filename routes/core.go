package routes

import (
	"errors"
	"fmt"
	"strings"

	"github.com/zenazn/goji/web"
)

// the struct all routes need to use
type Route struct {
	Method      string
	Pattern     string
	Description string
	Handler     web.HandlerType
}

// apply the prefix to each route in the routes array and add the
// route to the Mux.
func BuildRoutes(wm *web.Mux, routes []Route, prefix string) {
	for _, r := range routes {
		r.Pattern = prefix + r.Pattern
		err := getHandler(wm, &r)

		if err != nil {
			fmt.Println(err)
		}
	}
}

// attach r.Handler and r.Method to the correct verb function
func getHandler(wm *web.Mux, r *Route) error {
	switch strings.ToLower(r.Method) {
	case "get":
		wm.Get(r.Pattern, r.Handler)
	case "post":
		wm.Post(r.Pattern, r.Handler)
	case "put":
		wm.Put(r.Pattern, r.Handler)
	case "patch":
		wm.Patch(r.Pattern, r.Handler)
	case "delete":
		wm.Delete(r.Pattern, r.Handler)
	default:
		return errors.New("unsupported method: " + r.Method)
	}

	return nil
}
