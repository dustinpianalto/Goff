package info

import (
	"fmt"
	"github.com/dustinpianalto/disgoman"
)

var HelpCommand = &disgoman.Command{
	Name: "info",
	Aliases: []string{"?","h"},
	Description: "",
	OwnerOnly: false,
	Hidden: true,
	RequiredPermissions: 0,
	Invoke: HelpFunc,
}

func HelpFunc(ctx disgoman.Context, args []string) {
	var message = ">>> Command | Description\n"
	for name, command := range ctx.CommandManager.Commands {
		if command.OwnerOnly || command.Hidden {
			continue
		}

		message += fmt.Sprintf("%s | %s\n", name, command.Description)
	}
	_, _ = ctx.Send(message)
}