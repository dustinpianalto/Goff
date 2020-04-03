package exts

import "github.com/MikeModder/anpan"

func AddCommandHandlers(h *anpan.CommandHandler) {
	// Arguments:
	// name - command name - string
	// desc - command description - string
	// owneronly - only allow owners to run - bool
	// hidden - hide command from non-owners - bool
	// perms - permissisions required - anpan.Permission (int)
	// type - command type, sets where the command is available
	// run - function to run - func(anpan.Context, []string) / CommandRunFunc
	h.AddCommand("ping", "Check the bot's ping", false, false, 0, anpan.CommandTypeEverywhere, pingCommand)
	h.AddCommand("say", "Repeat a message", false, false, 0, anpan.CommandTypeEverywhere, sayCommand)
	h.AddCommand("user", "Show info about a user", false, false, 0, anpan.CommandTypeEverywhere, userCommand)
}
