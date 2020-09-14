package utils

import "database/sql"

type Guild struct {
	ID             string
	WelcomeMessage string
	GoodbyeMessage string
	LoggingChannel string
	WelcomeChannel string
	PuzzleChannel  string
	PuzzleRole     sql.NullString
}

type User struct {
	ID       string
	Banned   bool
	Logging  bool
	SteamID  string
	IsActive bool
	IsStaff  bool
	IsAdmin  bool
	Guilds   []string
}
