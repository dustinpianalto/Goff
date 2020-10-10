package postgres

import (
	"database/sql"
	"log"

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
	queryString := `INSERT INTO guilds (id, 
					welcome_message, 
					goodbye_message, 
					logging_channel, 
					welcome_channel, 
					puzzle_channel, 
					puzzle_role) 
					VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := s.DB.Exec(queryString,
		g.ID,
		g.WelcomeMessage,
		g.GoodbyeMessage,
		g.LoggingChannel,
		g.WelcomeChannel,
		g.PuzzleChannel,
		g.PuzzleRole,
	)
	return err
}

func (s *GuildService) DeleteGuild(g *goff.Guild) error {
	queryString := `DELETE CASCADE FROM guilds WHERE id = $1`
	_, err := s.DB.Exec(queryString, g.ID)
	return err
}

func (s *GuildService) GuildUsers(g *goff.Guild) ([]*goff.User, error) {
	var users []*goff.User
	queryString := `SELECT u.id, u.banned, u.logging, u.steam_id, u.is_active, u.is_staff, u.is_admin
					FROM users u, x_users_guilds xug
					WHERE xug.guild_id = $1
					AND xug.user_id = u.id`
	rows, err := s.DB.Query(queryString, g.ID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user goff.User
		err := rows.Scan(&user.ID,
			&user.Banned,
			&user.Logging,
			&user.SteamID,
			&user.IsActive,
			&user.IsStaff,
			&user.IsAdmin,
		)
		if err != nil {
			log.Println(err)
			continue
		}
		users = append(users, &user)
	}
	return users, nil
}

func (s *GuildService) UpdateGuild(g *goff.Guild) error {
	queryString := `UPDATE guilds SET 
					welcome_message = $1,
					goodbye_message = $2,
					logging_channel = $3,
					welcome_channel = $4,
					puzzle_channel = $5,
					puzzle_role = $6
					WHERE id = $7`
	_, err := s.DB.Exec(queryString,
		g.WelcomeMessage,
		g.GoodbyeMessage,
		g.LoggingChannel,
		g.WelcomeChannel,
		g.PuzzleChannel,
		g.PuzzleRole,
		g.ID,
	)
	return err
}
