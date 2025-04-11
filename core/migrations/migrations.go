package migrations

import (
	"apps/ecosystem/core/models"
	"apps/ecosystem/tools/db"
	_ "embed"

	"github.com/golang-migrate/migrate/v4/source"
	"github.com/google/uuid"
)

//go:embed init_up.sql
var db_init_up string

//go:embed init_down.sql
var db_init_down string

func GetMigrations() (*source.Migrations, error) {
	m := source.NewMigrations()
	var curVersion uint = 1

	m.Append(&source.Migration{Version: curVersion, Identifier: "v1 create _logs, _fields, _tables, tables", Direction: source.Up, Raw: db_init_up})
	m.Append(&source.Migration{Version: curVersion, Identifier: "v1 drop _logs, _fields, _tables, tables", Direction: source.Down, Raw: db_init_down})

	curVersion += 1
	// TODO is this a bad programmatic style, passing in m and curVersion and modifying it?
	if err := addModulesTableMigrations(m, &curVersion); err != nil {
		return nil, err
	}

	// if err := addUsersTableMigrations(m, &curVersion); err != nil {
	// 	return nil, err
	// }

	// if err := addScriptsTableMigrations(m, &curVersion); err != nil {
	// 	return nil, err
	// }

	// if err := addViewsTableMigrations(m, &curVersion); err != nil {
	// 	return nil, err
	// }

	// if err := addRolesTableMigrations(m, &curVersion); err != nil {
	// 	return nil, err
	// }

	// // TODO make below a standard method for creating junction tables
	// // TODO add same (add roles to) tables, fields, views, scripts
	// if err := addUsersRolesTableMigrations(m, &curVersion); err != nil {
	// 	return nil, err
	// }

	return m, nil
}

// TODO think of a nice way to DRY this all up, same thing repeated multiple times here
func addModulesTableMigrations(m *source.Migrations, v *uint) error  {
	id := uuid.New()
	table := models.TableModel{
		Name: "_modules",
		Id: id.String(),
		System: true,
		Fields: []models.Field{
			{Name: "id", Type: models.FieldChar, Size: 36, Primary: true},
			{Name: "name", Type: models.FieldVarchar, Size: 255},
			{Name: "created", Type: models.FieldTimestamp, Default: "CURRENT_TIMESTAMP"},
			{Name: "updated", Type: models.FieldTimestamp, Default: "CURRENT_TIMESTAMP"},
		},
	}
	query, _, err := db.NewQueryBuilder("").CreateTable(&table).Build()
	if err != nil {return err}
	m.Append(&source.Migration{Version: *v, Identifier: "v1 create _modules table", Direction: source.Up, Raw: query})

	query, _, err = db.NewQueryBuilder(table.Name).DropTable().Build()
	if err != nil {return err}
	m.Append(&source.Migration{Version: *v, Identifier: "v1 drop _modules table", Direction: source.Down, Raw: query})
	// TODO this seems kinda crazy...
	*v += 1

	// query, _, err = db.NewQueryBuilder("").InsertTableMeta(&table).Build()
	// if err != nil {return err}
	// m.Append(&source.Migration{Version: *v, Identifier: "v1 add _modules table metadata", Direction: source.Up, Raw: query})

	// query, _, err = db.NewQueryBuilder(table.Name).DeleteTableMeta().Build()
	// if err != nil {return err}
	// m.Append(&source.Migration{Version: *v, Identifier: "v1 remove _modules table metadata", Direction: source.Down, Raw: query})
	// *v += 1

	// query, _, err = db.NewQueryBuilder("_tables").AlterTable().CreateForeignKey(
	// 	&models.Field{Name: "module", ForeignKey: &models.ForeignKey{Table: "_modules", Column: "id"}},
	// 	).Build()
	// if err != nil {return err}
	// m.Append(&source.Migration{Version: *v, Identifier: "v1 add module fk to _tables", Direction: source.Up, Raw: query})

	// query, _, err = db.NewQueryBuilder("_tables").AlterTable().DropConstraint("fk_module__modules").Build()
	// if err != nil {return err}
	// m.Append(&source.Migration{Version: *v, Identifier: "v1 drop module fk to _tables", Direction: source.Down, Raw: query})
	// *v += 1

	// TODO add foreign key to fields table as well

	return nil
}

func addUsersTableMigrations(m *source.Migrations, v *uint) error  {
	id := uuid.New()
	table := models.TableModel{
		Name: "_users",
		Id: id.String(),
		System: true,
		Fields: []models.Field{
			{Name: "id", Type: models.FieldChar, Size: 36, Primary: true},
			{Name: "firstName", Type: models.FieldVarchar, Size: 255},
			{Name: "lastName", Type: models.FieldVarchar, Size: 255},
			{Name: "email", Type: models.FieldVarchar, Size: 255, Unique: true},
			{Name: "password", Type: models.FieldVarchar, Size: 255},
			{Name: "created", Type: models.FieldTimestamp, Default: "CURRENT_TIMESTAMP"},
			{Name: "updated", Type: models.FieldTimestamp, Default: "CURRENT_TIMESTAMP"},
		},
	}
	query, _, err := db.NewQueryBuilder("").CreateTable(&table).Build()
	if err != nil {return err}
	m.Append(&source.Migration{Version: *v, Identifier: "v1 create _users table", Direction: source.Up, Raw: query})

	query, _, err = db.NewQueryBuilder(table.Name).DropTable().Build()
	if err != nil {return err}
	m.Append(&source.Migration{Version: *v, Identifier: "v1 drop _users table", Direction: source.Down, Raw: query})
	*v += 1

	// query, _, err = db.NewQueryBuilder("").InsertTableMeta(&table).Build()
	// if err != nil {return err}
	// m.Append(&source.Migration{Version: *v, Identifier: "v1 add _users table metadata", Direction: source.Up, Raw: query})

	// query, _, err = db.NewQueryBuilder(table.Name).DeleteTableMeta().Build()
	// if err != nil {return err}
	// m.Append(&source.Migration{Version: *v, Identifier: "v1 remove _users table metadata", Direction: source.Down, Raw: query})
	// *v += 1

	return nil
}

// func addRolesTableMigrations(m *source.Migrations, v *uint) error  {
// 	id := uuid.New()
// 	table := models.TableModel{
// 		Name: "_roles",
// 		Id: id.String(),
// 		System: true,
// 		Fields: []models.Field{
// 			{Name: "id", Type: models.FieldChar, Size: 36, Primary: true},
// 			{Name: "name", Type: models.FieldVarchar, Size: 255},
// 			{Name: "created", Type: models.FieldTimestamp, Default: "CURRENT_TIMESTAMP"},
// 			{Name: "updated", Type: models.FieldTimestamp, Default: "CURRENT_TIMESTAMP"},
// 		},
// 	}
// 	query, _, err := db.NewQueryBuilder("").CreateTable(&table).Build()
// 	if err != nil {return err}
// 	m.Append(&source.Migration{Version: *v, Identifier: "v1 create _roles table", Direction: source.Up, Raw: query})

// 	query, _, err = db.NewQueryBuilder(table.Name).DropTable().Build()
// 	if err != nil {return err}
// 	m.Append(&source.Migration{Version: *v, Identifier: "v1 drop _roles table", Direction: source.Down, Raw: query})
// 	*v += 1

// 	query, _, err = db.NewQueryBuilder("").InsertTableMeta(&table).Build()
// 	if err != nil {return err}
// 	m.Append(&source.Migration{Version: *v, Identifier: "v1 add _roles table metadata", Direction: source.Up, Raw: query})

// 	query, _, err = db.NewQueryBuilder(table.Name).DeleteTableMeta().Build()
// 	if err != nil {return err}
// 	m.Append(&source.Migration{Version: *v, Identifier: "v1 remove _roles table metadata", Direction: source.Down, Raw: query})
// 	*v += 1

// 	return nil
// }

// func addUsersRolesTableMigrations(m *source.Migrations, v *uint) error  {
// 	id := uuid.New()
// 	table := models.TableModel{
// 		Name: "_user_roles",
// 		Id: id.String(),
// 		System: true,
// 		Fields: []models.Field{
// 			{Name: "role", Type: models.FieldChar, Size: 36, Primary: true, ForeignKey: &models.ForeignKey{
// 				Table: "_roles", Column: "id"}},
// 			{Name: "user", Type: models.FieldChar, Size: 36, Primary: true, ForeignKey: &models.ForeignKey{
// 				Table: "_users", Column: "id"}},
// 			{Name: "created", Type: models.FieldTimestamp, Default: "CURRENT_TIMESTAMP"},
// 			{Name: "updated", Type: models.FieldTimestamp, Default: "CURRENT_TIMESTAMP"},
// 		},
// 	}
// 	query, _, err := db.NewQueryBuilder("").CreateTable(&table).Build()
// 	if err != nil {return err}
// 	m.Append(&source.Migration{Version: *v, Identifier: "v1 create _user_roles table", Direction: source.Up, Raw: query})

// 	query, _, err = db.NewQueryBuilder(table.Name).DropTable().Build()
// 	if err != nil {return err}
// 	m.Append(&source.Migration{Version: *v, Identifier: "v1 drop _user_roles table", Direction: source.Down, Raw: query})
// 	*v += 1

// 	query, _, err = db.NewQueryBuilder("").InsertTableMeta(&table).Build()
// 	if err != nil {return err}
// 	m.Append(&source.Migration{Version: *v, Identifier: "v1 add _user_roles table metadata", Direction: source.Up, Raw: query})

// 	query, _, err = db.NewQueryBuilder(table.Name).DeleteTableMeta().Build()
// 	if err != nil {return err}
// 	m.Append(&source.Migration{Version: *v, Identifier: "v1 remove _user_roles table metadata", Direction: source.Down, Raw: query})
// 	*v += 1

// 	return nil
// }

func addScriptsTableMigrations(m *source.Migrations, v *uint) error  {
	id := uuid.New()
	table := models.TableModel{
		Name: "_scripts",
		Id: id.String(),
		System: true,
		Fields: []models.Field{
			{Name: "id", Type: models.FieldChar, Size: 36, Primary: true},
			{Name: "name", Type: models.FieldVarchar, Size: 255, Unique: true},
			{Name: "script", Type: models.FieldText, Nullable: true, Default: "''"},
			{Name: "created", Type: models.FieldTimestamp, Default: "CURRENT_TIMESTAMP"},
			{Name: "updated", Type: models.FieldTimestamp, Default: "CURRENT_TIMESTAMP"},
		},
	}
	query, _, err := db.NewQueryBuilder("").CreateTable(&table).Build()
	if err != nil {return err}
	m.Append(&source.Migration{Version: *v, Identifier: "v1 create _scripts table", Direction: source.Up, Raw: query})

	query, _, err = db.NewQueryBuilder(table.Name).DropTable().Build()
	if err != nil {return err}
	m.Append(&source.Migration{Version: *v, Identifier: "v1 drop _scripts table", Direction: source.Down, Raw: query})
	*v += 1

	// query, _, err = db.NewQueryBuilder("").InsertTableMeta(&table).Build()
	// if err != nil {return err}
	// m.Append(&source.Migration{Version: *v, Identifier: "v1 add _scripts table metadata", Direction: source.Up, Raw: query})

	// query, _, err = db.NewQueryBuilder(table.Name).DeleteTableMeta().Build()
	// if err != nil {return err}
	// m.Append(&source.Migration{Version: *v, Identifier: "v1 remove _scripts table metadata", Direction: source.Down, Raw: query})
	// *v += 1

	return nil
}

func addViewsTableMigrations(m *source.Migrations, v *uint) error  {
	id := uuid.New()
	table := models.TableModel{
		Name: "_views",
		Id: id.String(),
		System: true,
		Fields: []models.Field{
			{Name: "id", Type: models.FieldChar, Size: 36, Primary: true},
			{Name: "name", Type: models.FieldVarchar, Size: 255, Unique: true},
			{Name: "view", Type: models.FieldText, Nullable: true, Default: "''"},
			{Name: "viewJSON", Type: models.FieldJson, Nullable: true, Default: "'{}'"},
			{Name: "created", Type: models.FieldTimestamp, Default: "CURRENT_TIMESTAMP"},
			{Name: "updated", Type: models.FieldTimestamp, Default: "CURRENT_TIMESTAMP"},
		},
	}
	query, _, err := db.NewQueryBuilder("").CreateTable(&table).Build()
	if err != nil {return err}
	m.Append(&source.Migration{Version: *v, Identifier: "v1 create _views table", Direction: source.Up, Raw: query})

	query, _, err = db.NewQueryBuilder(table.Name).DropTable().Build()
	if err != nil {return err}
	m.Append(&source.Migration{Version: *v, Identifier: "v1 drop _views table", Direction: source.Down, Raw: query})
	*v += 1

	// query, _, err = db.NewQueryBuilder("").InsertTableMeta(&table).Build()
	// if err != nil {return err}
	// m.Append(&source.Migration{Version: *v, Identifier: "v1 add _views table metadata", Direction: source.Up, Raw: query})

	// query, _, err = db.NewQueryBuilder(table.Name).DeleteTableMeta().Build()
	// if err != nil {return err}
	// m.Append(&source.Migration{Version: *v, Identifier: "v1 remove _views table metadata", Direction: source.Down, Raw: query})
	// *v += 1

	return nil
}
