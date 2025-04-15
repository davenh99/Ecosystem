package db

import (
	"apps/ecosystem/core/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
)

// TODO maybe below is not performant, also figure out what the hell is going on...
func ScanRowIntoMap(rows *sql.Rows) (*map[string]any, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	record := make(map[string]any)

	values := make([]any, len(columns))
	valuePtrs := make([]any, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	err = rows.Scan(valuePtrs...)
	if err != nil {
		return nil, err
		// return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	for i, col := range columns {
		val := values[i]
		if b, ok := val.([]byte); ok {
			record[col] = string(b)
		} else {
			record[col] = val
		}
	}

	return &record, nil
}

func ScanRowIntoModel[T any](rows *sql.Rows) (*T, error) {
	model := new(T)
	
	v := reflect.ValueOf(model)
	values := make([]any, v.NumField())

	for i := range v.NumField() {
		values[i] = v.Field(i).Interface()
	}

	// TODO obvs below is wrong, the values don't end up back in 'model'
	err := rows.Scan(values...)
	
	if err != nil {
		return nil, err
	}

	return model, nil
}

func ScanRowIntoModelTable(rows *sql.Rows) (*models.TableModel, error) {
	table := new(models.TableModel)
	fieldsJson := new([]byte)
	module := new(sql.NullString)

	err := rows.Scan(
		&table.Id,
		&table.System,
		module,
		&table.Name,
		fieldsJson,
		&table.Created,
		&table.Updated,
	)
	if err != nil {
		return nil, err
	}

	if module.Valid {
		table.Module = module.String
	}

	err = json.Unmarshal(*fieldsJson, &table.Fields)
	if err != nil {
		return nil, err
	}

	return table, nil
}

func ScanRowIntoModelUser(rows *sql.Rows) (*models.UserModel, error) {
	user := new(models.UserModel)

	err := rows.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		// &user.Password,
		&user.Created,
		&user.Updated,
	)
	
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func ScanRowIntoModelAuth(rows *sql.Rows) (*models.AuthModel, error) {
	user := new(models.AuthModel)

	// TODO below will fail without proper query
	err := rows.Scan(
		&user.Id,
		// &user.Email,
		// &user.FirstName,
		// &user.LastName,
		&user.Password,
	)
	
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func ScanRowIntoModelRole(rows *sql.Rows) (*models.RoleModel, error) {
	role := new(models.RoleModel)

	err := rows.Scan(
		&role.Id,
		&role.Name,
		&role.Created,
		&role.Updated,
	)
	
	if err != nil {
		return nil, err
	}
	
	return role, nil
}