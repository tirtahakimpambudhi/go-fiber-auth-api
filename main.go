package main

import (
	"github.com/tirtahakimpambudhi/restful_api/internal/configs"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs/bootstrap"
	"github.com/tirtahakimpambudhi/restful_api/internal/delivery/http"
	"github.com/tirtahakimpambudhi/restful_api/internal/delivery/http/route"
	"log"
)

func main() {
	configs.ConfigFile = ".env.dev"
	app, err := bootstrap.New()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	usersController, authController, err := http.NewController(app)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	routes, errRoute := route.NewRoute(app, usersController, authController)
	if errRoute != nil {
		log.Fatal(errRoute.Error())
		return
	}

	if errInit := routes.Init(app.FiberServer.App); errInit != nil {
		log.Fatal(errInit.Error())
		return
	}
	err = app.FiberServer.Serve()
	if err != nil {
		log.Fatal(err.Error())
	}
	//if errServer := app.FiberServer.App.Listen(":3000"); errServer != nil {
	//	log.Fatal(err.Error())
	//	return
	//}

}
