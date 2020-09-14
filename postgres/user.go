package postgres

import (
	"database/sql"

	"github.com/dustinpianalto/goff"
)

type UserService struct {
	DB *sql.DB
}

func (s *UserService) User(id string) (*goff.User, error) {
	var u goff.User
	queryString := `SELECT id, banned, logging, steam_id, is_active, is_staff, is_admin 
						FROM users WHERE id=$1`
	row := s.DB.QueryRow(queryString, id)
	if err := row.Scan(&u.ID, &u.Banned, &u.Logging, &u.SteamID, &u.IsActive, &u.IsStaff, &u.IsAdmin); err != nil {
		return nil, err
	}
	var guilds []string
	queryString = `SELECT guild_id from x_users_guilds WHERE user_id=$1`
	rows, err := s.DB.Query(queryString, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var guildID string
		err = rows.Scan(&guildID)
		if err != nil {
			return nil, err
		}
		guilds = append(guilds, guildID)
	}
	u.Guilds = guilds
	return &u, nil
}
