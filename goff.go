package goff

import "github.com/dustinpianalto/disgoman"

type CommandManager struct {
	UserService  UserService
	GuildService GuildService
	disgoman.CommandManager
}
