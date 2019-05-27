package example

import (
	"github.com/compty001/go-admin/context"
	"github.com/compty001/go-admin/modules/config"
	"github.com/compty001/go-admin/plugins"
)

type Example struct {
	app *context.App
}

var Plug = new(Example)

var Config config.Config

func (example *Example) InitPlugin() {
	Config = config.Get()
	Config.PREFIX = "/" + Config.PREFIX
	Plug.app = InitRouter(Config.PREFIX)
}

func NewExample() *Example {
	return Plug
}

func (example *Example) GetRequest() []context.Path {
	return example.app.Requests
}

func (example *Example) GetHandler(url, method string) context.Handler {
	return plugins.GetHandler(url, method, example.app)
}
