//go:build wasm

package controllers

import app_types "exocomp-app/types"
import "github.com/cookiengineer/gooey/components/app"
import "github.com/cookiengineer/gooey/components/interfaces"

func RegisterTo(main *app.Main, config *app_types.Config) {

	main.RegisterController("chat", app.WrapController(func(main *app.Main, view interfaces.View) *Chat {
		return NewChat(main, view, config)
	}))

	main.RegisterController("bugs", app.WrapController(func(main *app.Main, view interfaces.View) *Bugs {
		return NewBugs(main, view)
	}))

	main.RegisterController("changelog", app.WrapController(func(main *app.Main, view interfaces.View) *Changelog {
		return NewChangelog(main, view)
	}))

	main.RegisterController("requirements", app.WrapController(func(main *app.Main, view interfaces.View) *Requirements {
		return NewRequirements(main, view)
	}))

}
