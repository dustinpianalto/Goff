package exts

import (
	"github.com/dustinpianalto/disgoman"
)

func AddCommandHandlers(h *disgoman.CommandManager) {
	// Arguments:
	// name - command name - string
	// desc - command description - string
	// owneronly - only allow owners to run - bool
	// hidden - hide command from non-owners - bool
	// perms - permissisions required - anpan.Permission (int)
	// type - command type, sets where the command is available
	// run - function to run - func(anpan.Context, []string) / CommandRunFunc
	_ = h.AddCommand(&disgoman.Command{
		Name:                "ping",
		Aliases:             nil,
		Description:         "Check the bot's ping",
		OwnerOnly:           false,
		Hidden:              false,
		RequiredPermissions: 0,
		Invoke:              pingCommand,
	})
	_ = h.AddCommand(&disgoman.Command{
		Name:                "say",
		Aliases:             nil,
		Description:         "Repeat a message",
		OwnerOnly:           false,
		Hidden:              false,
		RequiredPermissions: 0,
		Invoke:              sayCommand,
	})
	_ = h.AddCommand(&disgoman.Command{
		Name:                "user",
		Aliases:             nil,
		Description:         "Get user info",
		OwnerOnly:           false,
		Hidden:              false,
		RequiredPermissions: 0,
		Invoke:              userCommand,
	})
	_ = h.AddCommand(&disgoman.Command{
		Name:                "git",
		Aliases:             nil,
		Description:         "Show my github link",
		OwnerOnly:           false,
		Hidden:              false,
		RequiredPermissions: 0,
		Invoke:              gitCommand,
	})
	_ = h.AddCommand(&disgoman.Command{
		Name:                "tag",
		Aliases:             nil,
		Description:         "Get a tag",
		OwnerOnly:           false,
		Hidden:              false,
		RequiredPermissions: 0,
		Invoke:              tagCommand,
	})
	_ = h.AddCommand(&disgoman.Command{
		Name:                "addtag",
		Aliases:             nil,
		Description:         "Add a tag",
		OwnerOnly:           false,
		Hidden:              false,
		RequiredPermissions: 0,
		Invoke:              addTagCommand,
	})
	_ = h.AddCommand(&disgoman.Command{
		Name:                "invite",
		Aliases:             nil,
		Description:         "Get the invite link for this bot or others",
		OwnerOnly:           false,
		Hidden:              false,
		RequiredPermissions: 0,
		Invoke:              inviteCommand,
	})
	_ = h.AddCommand(&disgoman.Command{
		Name:                "P",
		Aliases:             nil,
		Description:         "Interpret a P\" program and return the results",
		OwnerOnly:           false,
		Hidden:              false,
		RequiredPermissions: 0,
		Invoke:              pCommand,
	})
	_ = h.AddCommand(&disgoman.Command{
		Name:                "set-logging-channel",
		Aliases:             []string{"slc"},
		Description:         "Set the channel logging messages will be sent to.",
		OwnerOnly:           false,
		Hidden:              false,
		RequiredPermissions: disgoman.PermissionManageServer,
		Invoke:              loggingChannel,
	})
	_ = h.AddCommand(&disgoman.Command{
		Name:                "get-logging-channel",
		Aliases:             []string{"glc"},
		Description:         "Gets the channel logging messages will be sent to.",
		OwnerOnly:           false,
		Hidden:              false,
		RequiredPermissions: disgoman.PermissionManageServer,
		Invoke:              getLoggingChannel,
	})
	_ = h.AddCommand(&disgoman.Command{
		Name:                "set-welcome-channel",
		Aliases:             []string{"swc"},
		Description:         "Set the channel welcome messages will be sent to.",
		OwnerOnly:           false,
		Hidden:              false,
		RequiredPermissions: disgoman.PermissionManageServer,
		Invoke:              welcomeChannel,
	})
	_ = h.AddCommand(&disgoman.Command{
		Name:                "get-welcome-channel",
		Aliases:             []string{"gwc"},
		Description:         "Gets the channel welcome messages will be sent to.",
		OwnerOnly:           false,
		Hidden:              false,
		RequiredPermissions: disgoman.PermissionManageServer,
		Invoke:              getWelcomeChannel,
	})
	_ = h.AddCommand(&disgoman.Command{
		Name:                "kick",
		Aliases:             []string{""},
		Description:         "Kicks the given user with the given reason",
		OwnerOnly:           false,
		Hidden:              false,
		RequiredPermissions: disgoman.PermissionKickMembers,
		Invoke:              kickUser,
	})
}
