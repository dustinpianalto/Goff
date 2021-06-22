package services

import "github.com/dustinpianalto/goff"

var UserService goff.UserService
var GuildService goff.GuildService
var RoleService goff.RoleService

func InitalizeServices(us goff.UserService, gs goff.GuildService, rs goff.RoleService) {
	UserService = us
	GuildService = gs
	RoleService = rs
}
