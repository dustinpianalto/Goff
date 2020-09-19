package postgres

import (
	"database/sql"

	"github.com/dustinpianalto/goff"
)

type GuildService struct {
	DB *sql.DB
}

func (s *GuildService) Guild(id string) (*goff.Guild, error) {
	var g goff.Guild
	queryString := `SELECT id, welcome_message, goodbye_message, 
							logging_channel, welcome_channel, puzzle_channel, puzzle_role 
						FROM guilds 
						WHERE id = $1`
	row := s.DB.QueryRow(queryString, id)
	err := row.Scan(
		&g.ID,
		&g.WelcomeMessage,
		&g.GoodbyeMessage,
		&g.LoggingChannel,
		&g.WelcomeChannel,
		&g.PuzzleChannel,
		&g.PuzzleRole,
	)
	if err != nil {
		return nil, err
	}
	var prefixes []string
	queryString = `SELECT p.prefix 
					FROM prefixes p, x_guilds_prefixes xgp 
					WHERE p.id = xgp.prefix_id AND xgp.guild_id = $1`
	rows, err := s.DB.Query(queryString, id)
	if err == nil {
		for rows.Next() {
			var prefix string
			err = rows.Scan(&prefix)
			if err != nil {
				continue
			}
			prefixes = append(prefixes, prefix)
		}
	}
	g.Prefixes = prefixes
	return &g, nil
}

func (s *GuildService) CreateGuild(g *goff.Guild) error {
	return nil
}

func (s *GuildService) DeleteGuild(g *goff.Guild) error {
	return nil
}

func (s *GuildService) GuildUsers(g *goff.Guild) ([]*goff.User, error) {
	return []*goff.User{}, nil
}

func (s *GuildService) UpdateGuild(g *goff.Guild) error {
	return nil
}
