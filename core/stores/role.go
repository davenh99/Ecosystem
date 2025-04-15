package stores

import (
	"apps/ecosystem/core/models"
	"apps/ecosystem/tools/db"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type RoleStore struct {
	db *sql.DB
}

func NewRoleStore(db *sql.DB) *RoleStore {
	return &RoleStore{db}
}

func (s *RoleStore) GetList(ctx context.Context) ([]models.RoleModel, error) {
	rows, err := db.NewQueryBuilder("_roles").Select().Query(s.db)
	if err != nil {
		return nil, err
	}

	users := make([]models.RoleModel, 0)
	for rows.Next() {
		r, err := db.ScanRowIntoModelRole(rows)
		if err != nil {
			return nil, err
		}

		users = append(users, models.RoleModel{Name: r.Name})
	}

	return users, nil
}

func (s *RoleStore) GetByID(id string) (*models.RoleModel, error) {
	rows, err := db.NewQueryBuilder("_roles").Select().Where("id", id).Query(s.db)
	if err != nil {
		return nil, err
	}

	r := new(models.RoleModel)
	for rows.Next() {
		r, err = db.ScanRowIntoModelRole(rows)
		if err != nil {
			return nil, err
		}
	}

	// TODO .... same as user thing
	if r.Id == "" {
		return nil, fmt.Errorf("role not found")
	}

	return r, nil
}

func (s *RoleStore) Create(role models.RoleModel) (string, error) {
	id := uuid.New()

	_, err := db.NewQueryBuilder("_roles").Insert(
		[]string{"id", "name"},
		[]any{id.String(), role.Name},
		).Exec(s.db)
	if err != nil {
		return "", err
	}
	
	return id.String(), nil
}

// TODO if keeping as below use PATCH instead of PUT???? or vice versa. huh?
func (s *RoleStore) Update(id string, payload map[string]any) error {
	cols := make([]string, 0)
	vals := make([]any, 0)

	for col, val := range payload {
		cols = append(cols, col)
		vals = append(vals, val)
	}
	_, err := db.NewQueryBuilder("_roles").Update(cols, vals).Where("id", id).Exec(s.db)
	if err != nil {
		return err
	}
	return nil
}

func (s *RoleStore) Delete(id string) error {
	_, err := db.NewQueryBuilder("_roles").Delete().Where("id", id).Exec(s.db)
	if err != nil {
		return err
	}
	
	return nil
}

func (s *RoleStore) AssignRoleToUser(userId string, roleId string) error {
	_, err := db.NewQueryBuilder("_user_roles").Insert(
		[]string{"user", "role"},
		[]any{userId, roleId},
		).Exec(s.db)
	if err != nil {
		return err
	}
	return nil
}
