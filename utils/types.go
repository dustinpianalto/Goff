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
