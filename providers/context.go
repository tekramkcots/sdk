package providers

import "github.com/tekramkcots/sdk/app"

func InjectDefaultAppContext(appCtx *app.Context) *app.Context {
	appCtx.Logger = app.NewLogger()
	return appCtx
}

func GetAppContext() *app.Context {
	return InjectDefaultAppContext(&app.Context{})
}
