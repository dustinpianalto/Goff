package goff

import "database/sql"

type Guild struct {
	ID             string
	WelcomeMessage string
	GoodbyeMessage string
	LoggingChannel string
	WelcomeChannel string
	PuzzleChannel  string
	PuzzleRole     sql.NullString
	Prefixes       []string
}

type GuildService interface {
	Guild(id string) (*Guild, error)
	CreateGuild(g *Guild) error
	DeleteGuild(g *Guild) error
	GuildUsers(g *Guild) ([]*User, error)
	UpdateGuild(g *Guild) error
}
