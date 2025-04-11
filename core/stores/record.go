package stores

import (
	"apps/ecosystem/tools/db"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type RecordStore struct {
	tableName string
	db *sql.DB
}

func NewRecordStore(db *sql.DB) *RecordStore {
	return &RecordStore{db: db}
}

func (s *RecordStore) SetTableName(name string) {
	s.tableName = name
}

// TODO add request auth params? (api rules from pocketbase)
// TODO add search filters and params
// TODO do this soon, when table is big this will return a lot of data...
func (s *RecordStore) GetList(ctx context.Context) ([]map[string]any, error) {
	rows, err := db.NewQueryBuilder(s.tableName).Select().Query(s.db)
	if err != nil {
		return nil, err
	}

	records := make([]map[string]any, 0)

	for rows.Next() {
		record, err := db.ScanRowIntoMap(rows)
		if err != nil {
			return nil, err
			// return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		records = append(records, *record)
	}

	// TODO check if no rows and return error

	return records, nil
}

func (s *RecordStore) GetByID(id string) (*map[string]any, error) {
	rows, err := db.NewQueryBuilder(s.tableName).Select().Where("id", id).Query(s.db)
	if err != nil {
		return nil, err
	}

	record := new(map[string]any)

	for rows.Next() {
		record, err = db.ScanRowIntoMap(rows)
		if err != nil {
			return nil, err
			// return nil, fmt.Errorf("failed to scan row: %w", err)
		}
	}

	// TODO check if no rows and return error

	return record, nil
}

// TODO in create and update, figure out whether i should check to make sure payload matches columns?
func (s *RecordStore) Create(payload map[string]any) (string, error) {
	id := uuid.New()

	// // confirm that columns match payload?
	// // then update by looping through payload?

	cols := []string{"id"}
	vals := []any{id.String()}

	for col, val := range payload {
		cols = append(cols, col)
		vals = append(vals, val)
	}
	_, err := db.NewQueryBuilder(s.tableName).Insert(cols, vals).Exec(s.db)

	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (s *RecordStore) Update(id string, payload map[string]any) error {
	cols := make([]string, 0)
	vals := make([]any, 0)

	for col, val := range payload {
		cols = append(cols, col)
		vals = append(vals, val)
	}
	_, err := db.NewQueryBuilder(s.tableName).Update(cols, vals).Where("id", id).Exec(s.db)
	if err != nil {
		return err
	}
	return nil
}

func (s *RecordStore) Delete(id string) error {
	_, err := db.NewQueryBuilder(s.tableName).Delete().Where("id", id).Exec(s.db)
	if err != nil {
		return err
	}
	return nil
}
