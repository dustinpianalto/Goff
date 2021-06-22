package postgres

import (
	"database/sql"
	"log"

	"github.com/dustinpianalto/goff"
)

type RoleService struct {
	DB *sql.DB
}

func (s *RoleService) Role(id string) (*goff.Role, error) {
	var r goff.Role
	queryString := `SELECT id, moderator, admin, self_assignable, guild_id FROM roles WHERE id = $1`
	row := s.DB.QueryRow(queryString, id)
	if err := row.Scan(&r.ID, &r.IsModerator, &r.IsAdmin, &r.SelfAssignable, &r.Guild); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *RoleService) AddRole(r *goff.Role) (*goff.Role, error) {
	queryString := `INSERT INTO roles (id, moderator, admin, self_assignable, auto_role, guild)
						VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING`
	_, err := s.DB.Exec(queryString, r.ID, r.IsModerator, r.IsAdmin, r.SelfAssignable, r.AutoRole, r.Guild)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (s *RoleService) DeleteRole(r *goff.Role) error {
	queryString := `DELETE FROM roles WHERE id = $1`
	_, err := s.DB.Exec(queryString, r.ID)
	return err
}

func (s *RoleService) MakeSelfAssignable(r *goff.Role) error {
	queryString := `UPDATE roles SET self_assignable = true WHERE id = $1`
	_, err := s.DB.Exec(queryString, r.ID)
	return err
}

func (s *RoleService) RemoveSelfAssignable(r *goff.Role) error {
	queryString := `UPDATE roles SET self_assignable = false WHERE id = $1`
	_, err := s.DB.Exec(queryString, r.ID)
	return err
}

func (s *RoleService) MakeAutoRole(r *goff.Role) error {
	queryString := `UPDATE roles SET auto_role = true WHERE id = $1`
	_, err := s.DB.Exec(queryString, r.ID)
	return err
}

func (s *RoleService) RemoveAutoRole(r *goff.Role) error {
	queryString := `UPDATE roles SET auto_role = false WHERE id = $1`
	_, err := s.DB.Exec(queryString, r.ID)
	return err
}

func (s *RoleService) GetAutoRoles(gID string) ([]*goff.Role, error) {
	var roles []*goff.Role
	queryString := `SELECT id FROM roles WHERE guild_id = $1`
	rows, err := s.DB.Query(queryString, gID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			log.Println(err)
			continue
		}
		role, err := s.Role(id)
		if err != nil {
			log.Println(err)
			continue
		}
		roles = append(roles, role)
	}
	return roles, nil
}
