//go:build wasm

package controllers

import app_types "exocomp-app/types"
import "github.com/cookiengineer/gooey/components/app"
import "github.com/cookiengineer/gooey/components/interfaces"

func RegisterTo(main *app.Main, config *app_types.Config) {

	main.RegisterController("chat", app.WrapController(func(main *app.Main, view interfaces.View) *Chat {
		return NewChat(main, view, config)
	}))

}
