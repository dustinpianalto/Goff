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
		SanitizeEveryone:    true,
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
		SanitizeEveryone:    true,
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
		Aliases:             nil,
		Description:         "Kicks the given user with the given reason",
		OwnerOnly:           false,
		Hidden:              false,
		RequiredPermissions: disgoman.PermissionKickMembers,
		Invoke:              kickUserCommand,
	})
	_ = h.AddCommand(&disgoman.Command{
		Name:                "addGuild",
		Aliases:             nil,
		Description:         "Adds the current guild to the database",
		OwnerOnly:           true,
		Hidden:              false,
		RequiredPermissions: disgoman.PermissionManageServer,
		Invoke:              addGuildCommand,
	})
	_ = h.AddCommand(&disgoman.Command{
		Name:                "ban",
		Aliases:             []string{"ban-no-delete"},
		Description:         "Bans the given user with the given reason",
		OwnerOnly:           false,
		Hidden:              false,
		RequiredPermissions: disgoman.PermissionBanMembers,
		Invoke:              banUserCommand,
	})
	_ = h.AddCommand(&disgoman.Command{
		Name:                "unban",
		Aliases:             nil,
		Description:         "Unbans the given user",
		OwnerOnly:           false,
		Hidden:              false,
		RequiredPermissions: disgoman.PermissionBanMembers,
		Invoke:              unbanUserCommand,
	})
	_ = h.AddCommand(&disgoman.Command{
		Name:                "remind",
		Aliases:             nil,
		Description:         "Remind me at a later time",
		OwnerOnly:           false,
		Hidden:              false,
		RequiredPermissions: 0,
		Invoke:              addReminderCommand,
	})
	_ = h.AddCommand(&disgoman.Command{
		Name:                "encode",
		Aliases:             []string{"e"},
		Description:         "Encode 2 numbers",
		OwnerOnly:           false,
		Hidden:              false,
		RequiredPermissions: 0,
		Invoke:              interleave,
	})
	_ = h.AddCommand(&disgoman.Command{
		Name:                "decode",
		Aliases:             []string{"d"},
		Description:         "Decode 1 number into 2",
		OwnerOnly:           false,
		Hidden:              false,
		RequiredPermissions: 0,
		Invoke:              deinterleave,
	})
	_ = h.AddCommand(&disgoman.Command{
		Name:                "RPN",
		Aliases:             []string{"rpn"},
		Description:         "Convert infix to rpn",
		OwnerOnly:           false,
		Hidden:              false,
		RequiredPermissions: 0,
		Invoke:              generateRPNCommand,
	})
	_ = h.AddCommand(&disgoman.Command{
		Name:                "ParseRPN",
		Aliases:             []string{'PRPN'},
		Description:         "Parse RPN string and return the result",
		OwnerOnly:           false,
		Hidden:              false,
		RequiredPermissions: 0,
		Invoke:              parseRPNCommand,
	})
	_ = h.AddCommand(&disgoman.Command{
		Name:                "solve",
		Aliases:             []string{"math"},
		Description:         "Solve infix equation and return the result",
		OwnerOnly:           false,
		Hidden:              false,
		RequiredPermissions: 0,
		Invoke:              solveCommand,
	})

}
