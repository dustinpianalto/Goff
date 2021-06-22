package roles

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dustinpianalto/disgoman"
	"github.com/dustinpianalto/goff"
	"github.com/dustinpianalto/goff/internal/postgres"
	"github.com/dustinpianalto/goff/internal/services"
)

var MakeRoleSelfAssignableCommand = &disgoman.Command{
	Name:                "make-role-self-assignable",
	Aliases:             []string{"makesar"},
	Description:         "Makes the passed in role self assignable by anyone",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: disgoman.PermissionManageServer,
	Invoke:              makeRoleSelfAssignableCommandFunc,
}

func makeRoleSelfAssignableCommandFunc(ctx disgoman.Context, args []string) {
	var roleString string
	var roleID string
	if len(args) > 0 {
		roleString = strings.Join(args, " ")
		if strings.HasPrefix(roleString, "<@&") && strings.HasSuffix(roleString, ">") {
			roleID = roleString[3 : len(roleString)-1]
		} else if _, err := strconv.Atoi(roleString); err == nil {
			roleID = roleString
		} else {
			for _, role := range ctx.Guild.Roles {
				if roleString == role.Name {
					roleID = role.ID
				}
			}
		}
	}
	fmt.Println(roleID)
	var role *discordgo.Role
	var err error
	if role, err = ctx.Session.State.Role(ctx.Guild.ID, roleID); err != nil {
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Can't find that Role.",
			Error:   err,
		}
		return
	}
	_, err = postgres.DB.Exec("INSERT INTO roles (id, guild_id, self_assignable) VALUES ($1, $2, true) ON CONFLICT (id) DO UPDATE SET self_assignable=true", role.ID, ctx.Guild.ID)
	if err != nil {
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Error Updating Database",
			Error:   err,
		}
		return
	}
	_, _ = ctx.Send(fmt.Sprintf("%s is now self assignable", role.Name))
}

var RemoveSelfAssignableCommand = &disgoman.Command{
	Name:                "remove-self-assignable-role",
	Aliases:             []string{"removesar"},
	Description:         "Makes a role that was previously self assignable not so",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: disgoman.PermissionManageServer,
	Invoke:              removeSelfAssignableRoleCommandFunc,
}

func removeSelfAssignableRoleCommandFunc(ctx disgoman.Context, args []string) {
	var roleString string
	var roleID string
	if len(args) > 0 {
		roleString = strings.Join(args, " ")
		if strings.HasPrefix(roleString, "<@&") && strings.HasSuffix(roleString, ">") {
			roleID = roleString[3 : len(roleString)-1]
		} else if _, err := strconv.Atoi(roleString); err == nil {
			roleID = roleString
		} else {
			for _, role := range ctx.Guild.Roles {
				if roleString == role.Name {
					roleID = role.ID
				}
			}
		}
	}
	fmt.Println(roleID)
	var err error
	var role *discordgo.Role
	if role, err = ctx.Session.State.Role(ctx.Guild.ID, roleID); err != nil {
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Can't find that Role.",
			Error:   err,
		}
		return
	}
	_, err = postgres.DB.Exec("INSERT INTO roles (id, guild_id, self_assignable) VALUES ($1, $2, false) ON CONFLICT (id) DO UPDATE SET self_assignable=false", role.ID, ctx.Guild.ID)
	if err != nil {
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Error Updating Database",
			Error:   err,
		}
		return
	}
	_, _ = ctx.Send(fmt.Sprintf("%s's self assignability has been removed.", role.Name))
}

var SelfAssignRoleCommand = &disgoman.Command{
	Name:                "giverole",
	Aliases:             []string{"iwant", "givetome"},
	Description:         "Assigns a person the passed in role if it is self assignable",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: 0,
	Invoke:              selfAssignRoleCommandFunc,
}

func selfAssignRoleCommandFunc(ctx disgoman.Context, args []string) {
	var roleString string
	var roleID string
	if len(args) > 0 {
		roleString = strings.Join(args, " ")
		if strings.HasPrefix(roleString, "<@&") && strings.HasSuffix(roleString, ">") {
			roleID = roleString[3 : len(roleString)-1]
		} else if _, err := strconv.Atoi(roleString); err == nil {
			roleID = roleString
		} else {
			for _, role := range ctx.Guild.Roles {
				if roleString == role.Name {
					roleID = role.ID
				}
			}
		}
	}
	fmt.Println(roleID)
	var role *discordgo.Role
	var err error
	if role, err = ctx.Session.State.Role(ctx.Guild.ID, roleID); err != nil {
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Can't find that Role.",
			Error:   err,
		}
		return
	}
	if memberHasRole(ctx.Member, role.ID) {
		_, _ = ctx.Send(fmt.Sprintf("You already have the %s role silly...", role.Name))
		return
	}
	var selfAssignable bool
	err = postgres.DB.QueryRow("SELECT self_assignable FROM roles where id=$1", role.ID).Scan(&selfAssignable)
	if err != nil {
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Error Updating Database",
			Error:   err,
		}
		return
	}
	if !selfAssignable {
		_, _ = ctx.Send(fmt.Sprintf("You aren't allowed to assign yourself the %s role", role.Name))
		return
	}
	err = ctx.Session.GuildMemberRoleAdd(ctx.Guild.ID, ctx.User.ID, role.ID)
	if err != nil {
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "There was a problem adding that role to you.",
			Error:   err,
		}
		return
	}
	_, _ = ctx.Send(fmt.Sprintf("Congratulations! The %s role has been added to your... Ummm... Thing.", role.Name))
}

var UnAssignRoleCommand = &disgoman.Command{
	Name:                "removerole",
	Aliases:             []string{"idon'twant"},
	Description:         "Removes a role from a person if the role is self assignable",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: 0,
	Invoke:              unAssignRoleCommandFunc,
}

func unAssignRoleCommandFunc(ctx disgoman.Context, args []string) {
	var roleString string
	var roleID string
	if len(args) > 0 {
		roleString = strings.Join(args, " ")
		if strings.HasPrefix(roleString, "<@&") && strings.HasSuffix(roleString, ">") {
			roleID = roleString[3 : len(roleString)-1]
		} else if _, err := strconv.Atoi(roleString); err == nil {
			roleID = roleString
		} else {
			for _, role := range ctx.Guild.Roles {
				if roleString == role.Name {
					roleID = role.ID
				}
			}
		}
	}
	fmt.Println(roleID)
	var role *discordgo.Role
	var err error
	if role, err = ctx.Session.State.Role(ctx.Guild.ID, roleID); err != nil {
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Can't find that Role.",
			Error:   err,
		}
		return
	}
	if !memberHasRole(ctx.Member, role.ID) {
		_, _ = ctx.Send(fmt.Sprintf("I can't remove the %s role from you because you don't have it...", role.Name))
		return
	}
	var selfAssignable bool
	err = postgres.DB.QueryRow("SELECT self_assignable FROM roles where id=$1", role.ID).Scan(&selfAssignable)
	if err != nil {
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Error Updating Database",
			Error:   err,
		}
		return
	}
	if !selfAssignable {
		_, _ = ctx.Send(fmt.Sprintf("You aren't allowed to remove the %s role from yourself", role.Name))
		return
	}
	err = ctx.Session.GuildMemberRoleRemove(ctx.Guild.ID, ctx.User.ID, role.ID)
	if err != nil {
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "There was a problem removing that role from you.",
			Error:   err,
		}
		return
	}
	_, _ = ctx.Send(fmt.Sprintf("Sad to see you go... but the %s role has been removed.", role.Name))
}

func memberHasRole(m *discordgo.Member, id string) bool {
	for _, r := range m.Roles {
		if r == id {
			return true
		}
	}
	return false
}

var MakeAutoRoleCommand = &disgoman.Command{
	Name:                "make-auto-role",
	Aliases:             []string{"makear"},
	Description:         "Marks the passed in role to auto add to new members.",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: disgoman.PermissionManageServer,
	Invoke:              makeAutoRoleCommandFunc,
}

func makeAutoRoleCommandFunc(ctx disgoman.Context, args []string) {
	var roleString string
	var roleID string
	if len(args) > 0 {
		roleString = strings.Join(args, " ")
		if strings.HasPrefix(roleString, "<@&") && strings.HasSuffix(roleString, ">") {
			roleID = roleString[3 : len(roleString)-1]
		} else if _, err := strconv.Atoi(roleString); err == nil {
			roleID = roleString
		} else {
			for _, role := range ctx.Guild.Roles {
				if roleString == role.Name {
					roleID = role.ID
				}
			}
		}
	}
	fmt.Println(roleID)
	var r *discordgo.Role
	var err error
	if r, err = ctx.Session.State.Role(ctx.Guild.ID, roleID); err != nil {
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Can't find that Role.",
			Error:   err,
		}
		return
	}
	role, err := services.RoleService.Role(r.ID)
	if err != nil {
		role = &goff.Role{
			ID:             r.ID,
			IsModerator:    false,
			IsAdmin:        false,
			SelfAssignable: false,
			AutoRole:       true,
			Guild:          ctx.Guild.ID,
		}
		role, err = services.RoleService.AddRole(role)
		if err != nil {
			ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
				Context: ctx,
				Message: "Error making an Auto Role",
				Error:   err,
			}
			return
		}
	} else {
		err = services.RoleService.MakeAutoRole(role)
		if err != nil {
			ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
				Context: ctx,
				Message: "Error making an Auto Role",
				Error:   err,
			}
			return
		}
	}

}

var RemoveAutoRoleCommand = &disgoman.Command{
	Name:                "remove-auto-role",
	Aliases:             []string{"removear"},
	Description:         "Remove role from those added to new members.",
	OwnerOnly:           false,
	Hidden:              false,
	RequiredPermissions: disgoman.PermissionManageServer,
	Invoke:              removeAutoRoleCommandFunc,
}

func removeAutoRoleCommandFunc(ctx disgoman.Context, args []string) {
	var roleString string
	var roleID string
	if len(args) > 0 {
		roleString = strings.Join(args, " ")
		if strings.HasPrefix(roleString, "<@&") && strings.HasSuffix(roleString, ">") {
			roleID = roleString[3 : len(roleString)-1]
		} else if _, err := strconv.Atoi(roleString); err == nil {
			roleID = roleString
		} else {
			for _, role := range ctx.Guild.Roles {
				if roleString == role.Name {
					roleID = role.ID
				}
			}
		}
	}
	fmt.Println(roleID)
	var r *discordgo.Role
	var err error
	if r, err = ctx.Session.State.Role(ctx.Guild.ID, roleID); err != nil {
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Can't find that Role.",
			Error:   err,
		}
		return
	}
	role, err := services.RoleService.Role(r.ID)
	if err != nil {
		ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
			Context: ctx,
			Message: "Role was not configured as an Auto Role",
			Error:   err,
		}
		return
	} else {
		err = services.RoleService.RemoveAutoRole(role)
		if err != nil {
			ctx.CommandManager.ErrorChannel <- disgoman.CommandError{
				Context: ctx,
				Message: "Error removing Auto Role",
				Error:   err,
			}
			return
		}
	}

}
