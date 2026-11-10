//go:build wasm

package main

import app_components "exocomp-app/components"
import app_controllers "exocomp-app/controllers"
import app_types "exocomp-app/types"
import app_views "exocomp-app/views"
import "github.com/cookiengineer/gooey/components/app"
import "github.com/cookiengineer/gooey/components/content"
import "github.com/cookiengineer/gooey/components/layout"
import "github.com/cookiengineer/gooey/components/ui"
import "time"

func main() {

	main := app.NewMain()

	content.RegisterTo(main.Document)
	layout.RegisterTo(main.Document)
	ui.RegisterTo(main.Document)
	app_components.RegisterTo(main.Document)

	go func() {

		config, err := app_types.BootstrapConfig(main.Client, "")

		if err != nil || config == nil {
			config = app_types.NewConfig()
		}

		app_views.RegisterTo(main)
		app_controllers.RegisterTo(main, config)

		main.Mount()
		main.Render()

	}()

	for true {
		time.Sleep(1 * time.Second)
	}

}
