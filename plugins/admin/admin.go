package admin

import (
	"github.com/compty001/go-admin/context"
	"github.com/compty001/go-admin/modules/config"
	"github.com/compty001/go-admin/modules/db"
	"github.com/compty001/go-admin/plugins"
	"github.com/compty001/go-admin/plugins/admin/controller"
	"github.com/compty001/go-admin/plugins/admin/models"
)

type Admin struct {
	app      *context.App
	tableCfg map[string]models.TableGenerator
}

func (admin *Admin) InitPlugin() {

	cfg := config.Get()

	// Init database
	for _, databaseCfg := range cfg.DATABASE {
		db.GetConnectionByDriver(databaseCfg.DRIVER).InitDB(map[string]config.Database{
			"default": databaseCfg,
		})
	}

	// Init router
	App.app = InitRouter("/" + cfg.PREFIX)

	models.SetGenerators(map[string]models.TableGenerator{
		"manager":    models.GetManagerTable,
		"permission": models.GetPermissionTable,
		"roles":      models.GetRolesTable,
		"op":         models.GetOpTable,
		"menu":       models.GetMenuTable,
	})
	models.SetGenerators(admin.tableCfg)
	models.InitTableList()

	cfg.PREFIX = "/" + cfg.PREFIX
	controller.SetConfig(cfg)

}

var App = new(Admin)

func NewAdmin(tableCfg map[string]models.TableGenerator) *Admin {
	App.tableCfg = tableCfg
	return App
}

func (admin *Admin) GetRequest() []context.Path {
	return admin.app.Requests
}

func (admin *Admin) GetHandler(url, method string) context.Handler {
	return plugins.GetHandler(url, method, admin.app)
}
