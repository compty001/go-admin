package main

import (
	_ "github.com/compty001/go-admin/adapter/http"
	"github.com/compty001/go-admin/engine"
	"github.com/compty001/go-admin/examples/datamodel"
	"github.com/compty001/go-admin/modules/config"
	"github.com/compty001/go-admin/modules/db"
	"github.com/compty001/go-admin/plugins/admin"
	"github.com/compty001/go-admin/plugins/example"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

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

	examplePlugin := example.NewExample()

	if err := eng.AddConfig(cfg).AddPlugins(admin.NewAdmin(datamodel.Generators), examplePlugin).
		Use(mux); err != nil {
		panic(err)
	}

	http.ListenAndServe(":9002", mux)
}
