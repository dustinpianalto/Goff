package user_management

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/dustinpianalto/goff"
	"github.com/dustinpianalto/goff/internal/services"
)

func OnGuildMemberAdd(s *discordgo.Session, member *discordgo.GuildMemberAdd) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic in OnGuildMemberAdd", r)
		}
	}()
	user, err := services.UserService.User(member.User.ID)
	if err != nil {
		log.Println("Error getting user from database: ", err)
		user = &goff.User{
			ID:       member.User.ID,
			Banned:   false,
			Logging:  true,
			IsActive: true,
			IsStaff:  false,
			IsAdmin:  false,
		}
		err := services.UserService.CreateUser(user, member.GuildID)
		if err != nil {
			log.Println("Error adding user to database: ", err)
		}
	}
	if !user.IsActive {
		user.IsActive = true
		err = services.UserService.UpdateUser(user)
		if err != nil {
			log.Println("Error marking user as active: ", err)
		}
	}
	err = services.UserService.AddUserToGuild(user, &goff.Guild{ID: member.GuildID})
	if err != nil {
		log.Println("Error adding user to guild: ", err)
	}
	roles, err := services.RoleService.GetAutoRoles(member.GuildID)
	if err != nil {
		log.Println("Error getting Auto Join Roles: ", err)
	}
	log.Println(roles)
	for _, r := range roles {
		role, err := s.State.Role(member.GuildID, r.ID)
		if err != nil {
			log.Println("Error getting role: ", err)
			continue
		}
		err = s.GuildMemberRoleAdd(member.GuildID, member.User.ID, role.ID)
		if err != nil {
			log.Println("Error adding Role to member: ", err)
			continue
		}
	}
}

func OnGuildMemberRemove(s *discordgo.Session, member *discordgo.GuildMemberRemove) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic in OnGuildMemberRemove: ", r)
		}
	}()
	user, err := services.UserService.User(member.User.ID)
	if err != nil {
		log.Println("Error getting user from database: ", err)
		return
	}
	err = services.UserService.RemoveUserFromGuild(user, &goff.Guild{ID: member.GuildID})
	if err != nil {
		log.Println("Error removing user from guild: ", err)
	}
	for i, g := range user.Guilds {
		if g == member.GuildID {
			user.Guilds[len(user.Guilds)-1], user.Guilds[i] = user.Guilds[i], user.Guilds[len(user.Guilds)-1]
			user.Guilds = user.Guilds[:len(user.Guilds)-1]
		}
	}
	if len(user.Guilds) == 0 {
		err = services.UserService.MarkUserInactive(user)
		if err != nil {
			log.Println("Error marking user as inactive: ", err)
		}
	}
}
