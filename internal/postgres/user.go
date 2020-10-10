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

func (s *UserService) CreateUser(u *goff.User) error {
	queryString := `INSERT INTO users (id, banned, logging, steam_id, is_active, is_staff, is_admin)
						VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := s.DB.Exec(queryString, u.ID, u.Banned, u.Logging, u.SteamID, u.IsActive, u.IsStaff, u.IsAdmin)
	return err
}

func (s *UserService) DeleteUser(u *goff.User) error {
	queryString := `DELETE CASCADE FROM x_users_guilds WHERE user_id = $1; DELETE FROM users WHERE id = $1`
	_, err := s.DB.Exec(queryString, u.ID)
	return err
}

func (s *UserService) MarkUserInactive(u *goff.User) error {
	queryString := `UPDATE users SET is_active = false WHERE id = $1`
	_, err := s.DB.Exec(queryString, u.ID)
	if err == nil {
		u.IsActive = false
	}
	return err
}

func (s *UserService) AddUserToGuild(u *goff.User, g *goff.Guild) error {
	queryString := `INSERT INTO x_users_guilds (user_id, guild_id) VALUES ($1, $2)`
	_, err := s.DB.Exec(queryString, u.ID, g.ID)
	return err
}

func (s *UserService) RemoveUserFromGuild(u *goff.User, g *goff.Guild) error {
	queryString := `DELETE FROM x_users_guilds WHERE user_id = $1 AND guild_id = $2`
	_, err := s.DB.Exec(queryString, u.ID, g.ID)
	return err
}

func (s *UserService) UpdateUser(u *goff.User) error {
	queryString := `UPDATE users SET 
						banned = $1,
						logging = $2, 
						steam_id = $3, 
						is_active = $4, 
						is_staff = $5,
						is_admin = $6
						WHERE id = $7`
	_, err := s.DB.Exec(queryString, u.Banned, u.Logging, u.SteamID, u.IsActive, u.IsStaff, u.IsAdmin)
	return err
}
