package services

import "github.com/dustinpianalto/goff"

var UserService goff.UserService
var GuildService goff.GuildService

func InitalizeServices(us goff.UserService, gs goff.GuildService) {
	UserService = us
	GuildService = gs
}
