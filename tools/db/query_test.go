package db

import (
	"apps/ecosystem/core/models"
	"errors"
	"testing"
)

func TestQueryBuilder(t *testing.T) {
	// TODO add tests: multiple foreign keys on create table, maybe properly test insert table?
	// TODO update tests: create/insert table, delete/drop table
	type result struct {
		query string
		args []any
		err error
	}

	type test struct {
		result result
		expectQuery string
		expectArgs []any
		expectError error
	}

	tests := map[string]test{}

	query, args, err := NewQueryBuilder("table").Select().Where("id", 123).Build()
	tests["build select all query"] = test{
		result: result{query, args, err},
		expectQuery: "SELECT * FROM table WHERE id=?",
		expectArgs: []any{123},
		expectError: nil,
	}

	columns := []string{"col1", "col2", "test"}
	query, args, err = NewQueryBuilder("table").Select(columns...).Where("not_id", "an_id").Build()
	tests["build select query with columns"] = test{
		result: result{query, args, err},
		expectQuery: "SELECT col1, col2, test FROM table WHERE not_id=?",
		expectArgs: []any{"an_id"},
		expectError: nil,
	}

	columns = []string{"col1", "col2"}
	vals := []any{0, false}
	query, args, err = NewQueryBuilder("foot").Update(columns, vals).Where("name", "test_id").Build()
	tests["build update query with 2 columns"] = test{
		result: result{query, args, err},
		expectQuery: "UPDATE foot SET col1=?, col2=? WHERE name=?",
		expectArgs: []any{0, false, "test_id"},
		expectError: nil,
	}

	columns = []string{"col1", "col2", "test"}
	vals = []any{12, 34, "str val"}
	query, args, err = NewQueryBuilder("foo").Update(columns, vals).Where("val", 98).Build()
	tests["build update query with 3 columns"] = test{
		result: result{query, args, err},
		expectQuery: "UPDATE foo SET col1=?, col2=?, test=? WHERE val=?",
		expectArgs: []any{12, 34, "str val", 98},
		expectError: nil,
	}
	
	query, args, err = NewQueryBuilder("foo").Delete().Where("val", 98).Build()
	tests["build delete query foo"] = test{
		result: result{query, args, err},
		expectQuery: "DELETE FROM foo WHERE val=?",
		expectArgs: []any{98},
		expectError: nil,
	}
	
	query, args, err = NewQueryBuilder("bar").Delete().Where("id", "good_id").Build()
	tests["build delete query bar"] = test{
		result: result{query, args, err},
		expectQuery: "DELETE FROM bar WHERE id=?",
		expectArgs: []any{"good_id"},
		expectError: nil,
	}

	columns = []string{"name", "in"}
	vals = []any{true, "help"}
	query, args, err = NewQueryBuilder("collection").Insert(columns, vals).Build()
	tests["build insert query with 2 columns"] = test{
		result: result{query, args, err},
		expectQuery: "INSERT INTO collection (name, in) VALUES (?, ?)",
		expectArgs: []any{true, "help"},
		expectError: nil,
	}

	columns = []string{"col1", "col2", "test", "an_col"}
	vals = []any{12, 34, "str val", false}
	query, args, err = NewQueryBuilder("table").Insert(columns, vals).Build()
	tests["build insert query with 4 columns"] = test{
		result: result{query, args, err},
		expectQuery: "INSERT INTO table (col1, col2, test, an_col) VALUES (?, ?, ?, ?)",
		expectArgs: []any{12, 34, "str val", false},
		expectError: nil,
	}

	// query, args, err = NewQueryBuilder("table_1").DropTable().Build()
	// tests["build drop table query"] = test{
	// 	result: result{query, args, err},
	// 	expectQuery: "DROP TABLE table_1",
	// 	expectArgs: []any{},
	// 	expectError: nil,
	// }

	// query, args, err = NewQueryBuilder("table 1").DropTable().Build()
	// tests["fail to build drop table query"] = test{
	// 	result: result{query, args, err},
	// 	expectQuery: "",
	// 	expectArgs: []any{},
	// 	expectError: fmt.Errorf("whitespace in table name"),
	// }

	// table := models.TableModel{
	// 	Id: "123",
	// 	Name: "new_table_1",
	// 	Fields: []models.Field{
	// 		{Name: "id", Primary: true, Nullable: false, Type: models.FieldInt, AutoIncrement: true},
	// 		{Name: "name", Type: models.FieldText, Nullable: true},
	// 		{Name: "column", Type: models.FieldBoolean, Default: "TRUE"},
	// 	},
	// }
	// query, args, err = NewQueryBuilder("notthetablename").CreateTable(&table).Build()
	// tests["build simple create table query"] = test{
	// 	result: result{query, args, err},
	// 	expectQuery: "CREATE TABLE new_table_1 (id INT NOT NULL AUTO_INCREMENT, name TEXT, column BOOLEAN NOT NULL DEFAULT " +
	// 				 "TRUE, CONSTRAINT pk_new_table_1 PRIMARY KEY (id))",
	// 	expectArgs: []any{},
	// 	expectError: nil,
	// }

	// table = models.TableModel{
	// 	Id: "456_hg",
	// 	Name: "good_table",
	// 	Fields: []models.Field{
	// 		{Name: "var", Primary: true, Nullable: false, Type: models.FieldInt},
	// 		{Name: "email", Type: models.FieldText, Unique: true},
	// 		{Name: "foo", Type: models.FieldChar, Default: "'1_3_5_7_555j'", Size: 12},
	// 		{Name: "bar", Type: models.FieldDatetime, Default: "CURRENT_TIMESTAMP"},
	// 	},
	// }
	// query, args, err = NewQueryBuilder("notthetablename").CreateTable(&table).Build()
	// tests["build average create table query"] = test{
	// 	result: result{query, args, err},
	// 	expectQuery: "CREATE TABLE good_table (var INT NOT NULL, email TEXT NOT NULL, foo CHAR(12) NOT NULL DEFAULT '1_3_5_7_555j'," +
	// 				" bar DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, CONSTRAINT uc_good_table UNIQUE KEY (email), CONSTRAINT" +
	// 				" pk_good_table PRIMARY KEY (var))",
	// 	expectArgs: []any{},
	// 	expectError: nil,
	// }

	// // fk, multiple pk, multiple unique
	// fk := models.ForeignKey{Table: "other_tbl", Column: "other_col"}
	// table = models.TableModel{
	// 	Id: "456_hg",
	// 	Name: "wow",
	// 	Fields: []models.Field{
	// 		{Name: "var", Primary: true, Nullable: false, Type: models.FieldInt, AutoIncrement: true},
	// 		{Name: "email", Type: models.FieldText, Primary: true},
	// 		{Name: "foo", Type: models.FieldText, ForeignKey: &fk},
	// 		{Name: "bar", Type: models.FieldInt, Unique: true},
	// 		{Name: "bar2", Type: models.FieldFloat, Unique: true, Size: 15, Nullable: true},
	// 	},
	// }
	// query, args, err = NewQueryBuilder("notthetablename").CreateTable(&table).Build()
	// tests["build complex create table query"] = test{
	// 	result: result{query, args, err},
	// 	expectQuery: "CREATE TABLE wow (var INT NOT NULL AUTO_INCREMENT, email TEXT NOT NULL, foo TEXT NOT NULL," +
	// 				" bar INT NOT NULL, bar2 FLOAT(15), CONSTRAINT uc_wow UNIQUE KEY (bar, bar2), CONSTRAINT pk_wow" +
	// 				" PRIMARY KEY (var, email), CONSTRAINT fk_wow_other_tbl FOREIGN KEY (foo) REFERENCES other_tbl(other_col))",
	// 	expectArgs: []any{},
	// 	expectError: nil,
	// }

	// table = models.TableModel{
	// 	Name: "bad 1",
	// 	Fields: []models.Field{
	// 		{Name: "var", Primary: true, Nullable: false, Type: models.FieldInt, AutoIncrement: true},
	// 		{Name: "email", Type: models.FieldText, Primary: true},
	// 	},
	// }
	// query, args, err = NewQueryBuilder("").CreateTable(&table).Build()
	// tests["fail to build create table query cos whitespace name"] = test{
	// 	result: result{query, args, err},
	// 	expectQuery: "",
	// 	expectArgs: []any{},
	// 	expectError: fmt.Errorf("whitespace in new table name"),
	// }

	// table = models.TableModel{
	// 	Name: "bad_2",
	// 	Fields: []models.Field{
	// 		{Name: "var", Primary: true, Nullable: true, Type: models.FieldInt, AutoIncrement: true},
	// 		{Name: "email", Type: models.FieldText, Primary: true},
	// 	},
	// }
	// query, args, err = NewQueryBuilder("").CreateTable(&table).Build()
	// tests["fail to build create table query cos nullable primary"] = test{
	// 	result: result{query, args, err},
	// 	expectQuery: "",
	// 	expectArgs: []any{},
	// 	expectError: fmt.Errorf("nullable primary key forbidden"),
	// }

	// table = models.TableModel{
	// 	Name: "bad_3",
	// 	Fields: []models.Field{
	// 		{Name: "var", Primary: true, Nullable: true, Type: models.FieldInt, AutoIncrement: true},
	// 		{Name: "var", Type: models.FieldText, Primary: true},
	// 	},
	// }
	// query, args, err = NewQueryBuilder("").CreateTable(&table).Build()
	// tests["fail to build create table query cos same col names"] = test{
	// 	result: result{query, args, err},
	// 	expectQuery: "",
	// 	expectArgs: []any{},
	// 	expectError: fmt.Errorf("duplicate column name var"),
	// }

	// query, args, err = NewQueryBuilder("yt").DeleteTableMeta().Build()
	// tests["build deletetable query"] = test{
	// 	result: result{query, args, err},
	// 	expectQuery: "DELETE FROM _tables WHERE name='yt'",
	// 	expectArgs: []any{},
	// 	expectError: nil,
	// }

	// table := models.TableModel{Name: "new_table_3"}
	// query, args, err = NewQueryBuilder("").InsertTableMeta(&table).Build()
	// tests["fail to build inserttable query"] = test{
	// 	result: result{query, args, err},
	// 	expectQuery: "",
	// 	expectArgs: []any{},
	// 	expectError: fmt.Errorf("id or name or fields not included"),
	// }

	fieldNew := models.Field{Name: "field1", Type: models.FieldText, Default: "'def'"}
	query, args, err = NewQueryBuilder("_1_").AlterTable().AddColumn(&fieldNew).Build()
	tests["build add column query"] = test{
		result: result{query, args, err},
		expectQuery: "ALTER TABLE _1_ ADD COLUMN field1 TEXT NOT NULL DEFAULT 'def'",
		expectArgs: []any{},
		expectError: nil,
	}

	fieldOld := models.Field{Name: "field1"}
	fieldNew = models.Field{Name: "field1_new"}
	query, args, err = NewQueryBuilder("_1_").AlterTable().AlterColumnName(&fieldOld, &fieldNew).Build()
	tests["build alter column name query"] = test{
		result: result{query, args, err},
		expectQuery: "ALTER TABLE _1_ RENAME COLUMN field1 TO field1_new",
		expectArgs: []any{},
		expectError: nil,
	}

	fieldOld = models.Field{Name: "field1"}
	fieldNew = models.Field{Name: "field1", Type: models.FieldFloat, Size: 5, Nullable: true}
	query, args, err = NewQueryBuilder("_1_").AlterTable().AlterColumnDefinition(&fieldOld, &fieldNew).Build()
	tests["build alter column definition query"] = test{
		result: result{query, args, err},
		expectQuery: "ALTER TABLE _1_ MODIFY COLUMN field1 FLOAT(5)",
		expectArgs: []any{},
		expectError: nil,
	}

	fieldOld = models.Field{Name: "help"}
	query, args, err = NewQueryBuilder("_2_").AlterTable().DropColumn(&fieldOld).Build()
	tests["build drop column query"] = test{
		result: result{query, args, err},
		expectQuery: "ALTER TABLE _2_ DROP COLUMN help",
		expectArgs: []any{},
		expectError: nil,
	}

	fieldNew = models.Field{Name: "help"}
	query, args, err = NewQueryBuilder("_tbl_").AlterTable().CreateIndex(&fieldNew).Build()
	tests["build add index query"] = test{
		result: result{query, args, err},
		expectQuery: "ALTER TABLE _tbl_ ADD INDEX idx__tbl__help (help)",
		expectArgs: []any{},
		expectError: nil,
	}

	fieldNew = models.Field{Name: "new", ForeignKey: &models.ForeignKey{Table: "huh", Column: "o_col"}}
	query, args, err = NewQueryBuilder("at_last").AlterTable().CreateForeignKey(&fieldNew).Build()
	tests["build add foreign key query"] = test{
		result: result{query, args, err},
		expectQuery: "ALTER TABLE at_last ADD CONSTRAINT fk_new_huh FOREIGN KEY (new) REFERENCES huh(o_col)",
		expectArgs: []any{},
		expectError: nil,
	}

	query, args, err = NewQueryBuilder("at_last").AlterTable().DropConstraint("fk_1_sdg").Build()
	tests["build drop constraint query"] = test{
		result: result{query, args, err},
		expectQuery: "ALTER TABLE at_last DROP CONSTRAINT fk_1_sdg",
		expectArgs: []any{},
		expectError: nil,
	}

	query, args, err = NewQueryBuilder("at_last").RenameTable("new_table").Build()
	tests["build rename table query"] = test{
		result: result{query, args, err},
		expectQuery: "RENAME TABLE at_last TO new_table",
		expectArgs: []any{},
		expectError: nil,
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.result.err != nil && tc.expectError != nil && errors.Is(tc.result.err, tc.expectError) {
				t.Fatalf("expected err: %v, got err: %v", tc.expectError, tc.result.err)
			}
			if !(tc.result.query == tc.expectQuery) {
				t.Fatalf("expected query:\n%v,\ngot query:\n%v\nerr: %v", tc.expectQuery, tc.result.query, tc.result.err)
			}
			if len(tc.expectArgs) != len(tc.result.args) {
				t.Fatalf("expected %d args, got %d",len(tc.expectArgs), len(tc.result.args))
			}
			for i, arg := range tc.result.args {
				if arg != tc.expectArgs[i] {
					t.Fatalf("expected arg at position %d: %v, got %v", i, tc.expectArgs[i], arg)
				}
			}
		})
	}
}
