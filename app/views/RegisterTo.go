//go:build wasm

package views

import "github.com/cookiengineer/gooey/components/app"

func RegisterTo(main *app.Main) {

	main.RegisterView("chat", app.WrapView(ToChat))

}
