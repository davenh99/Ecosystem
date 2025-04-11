package stores

import (
	"apps/ecosystem/core/models"
	"apps/ecosystem/tools/db"
	"context"
	"database/sql"
	"fmt"

	"github.com/stephenafamo/bob/dialect/mysql"
	"github.com/stephenafamo/bob/dialect/mysql/sm"

	"github.com/google/uuid"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db}
}

// TODO add request auth params? (api rules from pocketbase)
func (s *UserStore) GetList(ctx context.Context) ([]models.UserModel, error) {
	// TODO add search params, filters, limit etc
	q, args, err := mysql.Select(sm.From("_users")).Build(ctx)
	if err != nil {
		return nil, err
	}

	// rows, err := db.NewQueryBuilder("_users").Select().Query(s.db)
	rows, err := s.db.Query(q, args...)
	if err != nil {
		return nil, err
	}

	users := make([]models.UserModel, 0)
	for rows.Next() {
		u, err := db.ScanRowIntoModelUser(rows)
		if err != nil {
			return nil, err
		}

		users = append(users, *u)
	}

	return users, nil
}

func (s *UserStore) GetAuthByEmail(email string) (*models.AuthModel, error) {
	rows, err := db.NewQueryBuilder("_users").Select("id", "password").Where("email", email).Query(s.db)
	if err != nil {
		return nil, err
	}

	u := new(models.AuthModel)
	for rows.Next() {
		u, err = db.ScanRowIntoModelAuth(rows)
		if err != nil {
			return nil, err
		}
	}

	// TODO is below correct???
	if u.Id == "" {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func (s *UserStore) GetByID(id string) (*models.UserModel, error) {
	rows, err := db.NewQueryBuilder("_users").Select().Where("id", id).Query(s.db)
	if err != nil {
		return nil, err
	}

	u := new(models.UserModel)
	for rows.Next() {
		// TODO don't scan password here
		u, err = db.ScanRowIntoModelUser(rows)
		if err != nil {
			return nil, err
		}
	}

	// TODO below is huh? why are we checking only if it is an empty string? (was checking if 0 before when id was int)
	if u.Id == "" {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func (s *UserStore) Create(user models.AuthModel) (string, error) {
	id := uuid.New()
	// TODO below can probs be extracted when new struct/reflect stuff being used for orm stuff?
	cols := []string{"id", "firstName", "lastName", "email", "password"}
	vals := []any{id.String(), user.FirstName, user.LastName, user.Email, user.Password}

	_, err := db.NewQueryBuilder("_users").Insert(cols, vals).Exec(s.db)
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

// TODO if keeping as below use PATCH instead of PUT???? or vice versa. huh?
func (s *UserStore) Update(id string, payload map[string]any) error {
	cols := make([]string, 0)
	vals := make([]any, 0)

	for col, val := range payload {
		cols = append(cols, col)
		vals = append(vals, val)
	}
	_, err := db.NewQueryBuilder("_users").Update(cols, vals).Where("id", id).Exec(s.db)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) Delete(id string) error {
	_, err := db.NewQueryBuilder("_users").Delete().Where("id", id).Exec(s.db)
	if err != nil {
		return err
	}
	return nil
}
