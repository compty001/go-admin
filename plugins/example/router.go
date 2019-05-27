package example

import (
	"github.com/compty001/go-admin/context"
	"github.com/compty001/go-admin/modules/auth"
	"github.com/compty001/go-admin/plugins/admin/controller"
)

func InitRouter(prefix string) *context.App {
	app := context.NewApp()

	authenticator := auth.SetPrefix(prefix).SetAuthFailCallback(func(ctx *context.Context) {
		ctx.Write(302, map[string]string{
			"Location": prefix + "/login",
		}, ``)
	}).SetPermissionDenyCallback(func(ctx *context.Context) {
		controller.ShowErrorPage(ctx, "permission denied")
	})

	app.GET(prefix+"/example", authenticator.Middleware(TestHandler))

	if prefix == "" {
		app.GET(prefix+"/", authenticator.Middleware(TestHandler))
	} else {
		app.GET(prefix, authenticator.Middleware(TestHandler))
	}

	return app
}
