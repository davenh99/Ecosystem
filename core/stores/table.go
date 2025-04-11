package stores

import (
	"apps/ecosystem/core/models"
	"apps/ecosystem/tools/db"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
)

type TableStore struct {
	tableName string
	db *sql.DB
}

func NewTableStore(db *sql.DB) *TableStore {
	return &TableStore{db: db}
}

func (s *TableStore) SetTableName(name string) {
	s.tableName = name
}

// TODO add request auth params? (api rules from pocketbase)
// TODO add search filters and params
// TODO do this soon, when table is big this will return a lot of data...
func (s *TableStore) GetList() ([]models.TableModel, error) {
	// get table metadata from tables table
	rows, err := db.NewQueryBuilder("_tables").Select().Query(s.db)
	if err != nil {
		return nil, fmt.Errorf("problem getting tables metadata: %w", err)
	}

	tables := make([]models.TableModel, 0)
	for rows.Next() {
		t, err := db.ScanRowIntoModelTable(rows)
		if err != nil {
			// return nil, err
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		tables = append(tables, *t)
	}

	// TODO find out if I should be doing this?
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// TODO check if no rows and return error, see if i need to do this

	return tables, nil
}

func (s *TableStore) GetByName(tableName string) (*models.TableModel, error) {
	// get table metadata from _tables table
	rows, err := db.NewQueryBuilder("_tables").Select().Where("name", tableName).Query(s.db)
	if err != nil {
		return nil, err
	}

	t := new(models.TableModel)
	for rows.Next() {
		t, err = db.ScanRowIntoModelTable(rows)
		if err != nil {
			return nil, err
			// return nil, fmt.Errorf("failed to scan row: %w", err)
		}
	}

	// TODO find out if I should be doing this?
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// TODO check if no rows and return error

	return t, nil
}

// TODO would below ever be used?
func (s *TableStore) GetByID(id string) (*models.TableModel, error) {
	// get table metadata from _tables table
	rows, err := db.NewQueryBuilder("_tables").Select().Where("id", id).Query(s.db)
	if err != nil {
		return nil, err
	}

	t := new(models.TableModel)
	for rows.Next() {
		t, err = db.ScanRowIntoModelTable(rows)
		if err != nil {
			return nil, err
			// return nil, fmt.Errorf("failed to scan row: %w", err)
		}
	}

	// TODO find out if I should be doing this?
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// TODO check if no rows and return error

	return t, nil
}

func (s *TableStore) Create(payload models.TableModel) (string, error) {
	// TODO wrap below in transaction
	// first create table
	_, err := db.NewQueryBuilder(s.tableName).CreateTable(&payload).Exec(s.db)
	if err != nil {
		return "", err
	}
	// adding indexes to tables using alter table
	for _, field := range payload.Fields {
		if field.Index {
			_, err := db.NewQueryBuilder(s.tableName).AlterTable().CreateIndex(&field).Exec(s.db)
			if err != nil {
				return "", err
			}
		}
	}
	
	// then, update tables table
	id := uuid.New()
	// payload.Id = id.String()
	// _, err = db.NewQueryBuilder(s.tableName).InsertTableMeta(&payload).Exec(s.db)
	// if err != nil {
	// 	return "", err
	// }

	// make a migration? or leave up to modules?

	return id.String(), nil
}

func (s *TableStore) Update(id string, tableNew models.TableModel) error {
	// TODO wrap below in transaction
	cols := make([]string, 0)
	vals := make([]any, 0)

	// get old table to compare against
	rows, err := db.NewQueryBuilder("_tables").Select().Where("id", id).Query(s.db)
	if err != nil {
		return err
	}
	tableOld := new(models.TableModel)
	for rows.Next() {
		tableOld, err = db.ScanRowIntoModelTable(rows)
		if err != nil {
			return err
		}
	}

	addFields, updateFields, dropFields := tableNew.GetUpdatedFields(*tableOld)
	// add fields
	for _, field := range dropFields {
		_, err := db.NewQueryBuilder(s.tableName).AlterTable().DropColumn(&field).Exec(s.db)
		if err != nil {
			return err
		}
	}
	// update fields
	for _, fields := range updateFields {
		fieldOld, fieldNew := fields[0], fields[1]
		if fieldOld.Name != fieldNew.Name {
			_, err := db.NewQueryBuilder(s.tableName).AlterTable().AlterColumnName(&fieldOld, &fieldNew).Exec(s.db)
			if err != nil {
				return err
			}
			fieldOld.Name = fieldNew.Name
			if !fieldOld.Equals(&fieldNew) {
				_, err := db.NewQueryBuilder(s.tableName).AlterTable().AlterColumnDefinition(&fieldOld, &fieldNew).Exec(s.db)
				if err != nil {
					return err
				}
			}
		} else {
			_, err := db.NewQueryBuilder(s.tableName).AlterTable().AlterColumnDefinition(&fieldOld, &fieldNew).Exec(s.db)
			if err != nil {
				return err
			}
		}
	}
	// drop fields
	for _, field := range dropFields {
		_, err := db.NewQueryBuilder(s.tableName).AlterTable().DropColumn(&field).Exec(s.db)
		if err != nil {
			return err
		}
	}

	// add updated fields to table update args
	if len(addFields) < 1 && len(updateFields) < 1 && len(dropFields) < 1 {
		cols = append(cols, "fields")
		b, err := json.Marshal(tableNew.Fields)
        if err != nil {
            return err
        }
        // need to escape the single quotes for mysql
		vals = append(vals, "'" + strings.ReplaceAll(string(b), "'", "\\'") + "'")
	}

	// check if we need to rename the table
	if tableOld.Name != tableNew.Name {
		_, err := db.NewQueryBuilder(tableOld.Name).RenameTable(tableNew.Name).Exec(s.db)
		if err != nil {
			return err
		}
		cols = append(cols, "name")
		vals = append(vals, tableNew.Name)
	}
	
	// finally, update tables table
	_, err = db.NewQueryBuilder("_tables").Update(cols, vals).Where("id", id).Exec(s.db)
	if err != nil {
		return err
	}

	// make a migration? or leave up to modules?

	return nil
}

func (s *TableStore) Delete() error {
	// TODO wrap below in transaction

	// first drop table
	_, err := db.NewQueryBuilder(s.tableName).DropTable().Exec(s.db)
	if err != nil {
		return err
	}
	
	// then, delete from tables table
	// _, err = db.NewQueryBuilder(s.tableName).DeleteTableMeta().Exec(s.db)
	// if err != nil {
	// 	return err
	// }

	// make a migration? or leave up to modules?

	return nil
}
