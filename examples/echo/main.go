package main

import (
	_ "github.com/compty001/go-admin/adapter/echo"
	"github.com/compty001/go-admin/engine"
	"github.com/compty001/go-admin/examples/datamodel"
	"github.com/compty001/go-admin/modules/config"
	"github.com/compty001/go-admin/modules/db"
	"github.com/compty001/go-admin/plugins/admin"
	"github.com/compty001/go-admin/plugins/example"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	eng := engine.Default()

	cfg := config.Config{
		DATABASE: []config.Database{
			{
				HOST:         "127.0.0.1",
				PORT:         "3306",
				USER:         "root",
				PWD:          "root",
				NAME:         "godmin",
				MAX_IDLE_CON: 50,
				MAX_OPEN_CON: 150,
				DRIVER:       db.DRIVER_MYSQL,
			},
		},
		DOMAIN: "localhost",
		PREFIX: "admin",
		INDEX:  "/",
		DEBUG:  true,
	}

	adminPlugin := admin.NewAdmin(datamodel.Generators)
	examplePlugin := example.NewExample()

	if err := eng.AddConfig(cfg).AddPlugins(adminPlugin, examplePlugin).Use(e); err != nil {
		panic(err)
	}

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
