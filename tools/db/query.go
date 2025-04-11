package db

import (
	"apps/ecosystem/core/models"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
)

// TODO maybe I need to backtick the crap out of everything?
// TODO consider better way to do table create, atm can grab table from q.tableNew.Name or q.table, should only be one place
// TODO fk naming is kinda wierd atm...
// TODO maybe eventually replace with bob query builder?

type QueryType int

const (
    QuerySelect QueryType = iota
    QueryInsert
    QueryUpdate
    QueryDelete
    QueryCreateTable
    // QueryInsertTableMeta
    // QueryDeleteTableMeta
    QueryAlterTable
    QueryRenameTable
    QueryAlterTableAddColumn
    QueryAlterTableModifyColumnDefinition
    QueryAlterTableModifyColumnName
    QueryAlterTableDropColumn
    QueryAlterTableAddIndex
    QueryAlterTableAddFK
    QueryAlterTableDropConstraint
    QueryDropTable
)

type QueryBuilder struct {
    queryType QueryType
    alterTableQueryType QueryType
	table string
	columns []string
	wheres []string
    tableModel *models.TableModel
    fieldOld *models.Field
    fieldNew *models.Field
    dropConstraint string
	// orderBy  string
	// limit    int
	// offset   int
    args []any
    err error
}

func NewQueryBuilder(table string) *QueryBuilder {
	return &QueryBuilder{table: table, queryType: -1}
}

func (q *QueryBuilder) Select(columns ...string) *QueryBuilder {
    q.queryType = QuerySelect
    q.columns = columns
    return q
}

func (q *QueryBuilder) Insert(columns []string, vals []any) *QueryBuilder {
    q.queryType = QueryInsert
    if len(columns) != len(vals) {
        q.err = fmt.Errorf("length of columns not same as length of vals")
    }
    q.columns = columns
    q.args = append(q.args, vals...)
    return q
}

func (q *QueryBuilder) Update(columns []string, vals []any) *QueryBuilder {
    q.queryType = QueryUpdate
    if len(columns) != len(vals) {
        q.err = fmt.Errorf("length of columns not same as length of vals")
    }
    q.columns = columns
    q.args = append(q.args, vals...)
    return q
}

func (q *QueryBuilder) Delete() *QueryBuilder {
    q.queryType = QueryDelete
    return q
}

func (q *QueryBuilder) CreateTable(tableNew *models.TableModel) *QueryBuilder {
    q.queryType = QueryCreateTable
    q.tableModel = tableNew
    return q
}

// func (q *QueryBuilder) InsertTableMeta(t *models.TableModel) *QueryBuilder {
//     q.queryType = QueryInsertTableMeta
//     q.tableModel = t
//     return q
// }

// func (q *QueryBuilder) DeleteTableMeta() *QueryBuilder {
//     q.queryType = QueryDeleteTableMeta
//     return q
// }

func (q *QueryBuilder) DropTable() *QueryBuilder {
    q.queryType = QueryDropTable
    return q
}

func (q *QueryBuilder) AlterTable() *QueryBuilder {
    q.queryType = QueryAlterTable
    return q
}

func (q *QueryBuilder) RenameTable(name string) *QueryBuilder {
    q.queryType = QueryRenameTable
    q.tableModel = &models.TableModel{Name: name}
    return q
}

func (q *QueryBuilder) AddColumn(fieldNew *models.Field) *QueryBuilder {
    q.alterTableQueryType = QueryAlterTableAddColumn
    q.fieldNew = fieldNew
    return q
}

func (q *QueryBuilder) AlterColumnDefinition(fieldOld *models.Field, fieldNew *models.Field) *QueryBuilder {
    q.alterTableQueryType = QueryAlterTableModifyColumnDefinition
    q.fieldOld = fieldOld
    q.fieldNew = fieldNew
    return q
}

func (q *QueryBuilder) AlterColumnName(fieldOld *models.Field, fieldNew *models.Field) *QueryBuilder {
    q.alterTableQueryType = QueryAlterTableModifyColumnName
    q.fieldOld = fieldOld
    q.fieldNew = fieldNew
    return q
}

func (q *QueryBuilder) DropColumn(fieldOld *models.Field) *QueryBuilder {
    q.alterTableQueryType = QueryAlterTableDropColumn
    q.fieldOld = fieldOld
    return q
}

func (q *QueryBuilder) CreateIndex(fieldNew *models.Field) *QueryBuilder {
    q.alterTableQueryType = QueryAlterTableAddIndex
    q.fieldNew = fieldNew
    return q
}

func (q *QueryBuilder) CreateForeignKey(fieldNew *models.Field) *QueryBuilder {
    q.alterTableQueryType = QueryAlterTableAddFK
    q.fieldNew = fieldNew
    return q
}

func (q *QueryBuilder) DropConstraint(constraint string) *QueryBuilder {
    q.alterTableQueryType = QueryAlterTableDropConstraint
    q.dropConstraint = constraint
    return q
}

func (q *QueryBuilder) Where(column string, value any) *QueryBuilder {
    if q.queryType != QueryDelete && q.queryType != QuerySelect && q.queryType != QueryUpdate {
        q.err = fmt.Errorf("where used on invalid query type")
    }
    q.wheres = append(q.wheres, column)
    q.args = append(q.args, value)
    return q
}

func (q *QueryBuilder) Build() (string, []any, error) {
    return q.build()
}

func (q *QueryBuilder) Exec(db *sql.DB) (sql.Result, error) {
    query, args, err := q.build()
    if err != nil {
		return nil, fmt.Errorf("could not build exec query:, %v", err)
	}

    return db.Exec(query, args...)
}

func (q *QueryBuilder) Query(db *sql.DB) (*sql.Rows, error) {
    query, args, err := q.build()
    if err != nil {
		return nil, fmt.Errorf("could not build query:, %v", err)
	}

    return db.Query(query, args...)
}

// TODO could be possible to DRY this out...
func (q *QueryBuilder) build() (string, []any, error) {
    if q.err != nil {
        return "", nil, q.err
    }
    if q.queryType < 0 {
        return "", nil, fmt.Errorf("no query type specified")
    }
    if regexp.MustCompile(`\s`).MatchString(q.table) {
        return "", nil, fmt.Errorf("whitespace in table name")
    }
    if q.tableModel != nil && regexp.MustCompile(`\s`).MatchString(q.tableModel.Name) {
        return "", nil, fmt.Errorf("whitespace in new table name")
    }

    query := ""
    countedArgs := 0

    buildWheres := func() {
        if len(q.wheres) > 0 {
            query += " WHERE " + q.wheres[0] + "=?"
            countedArgs += 1
            for _, where := range q.wheres[1:] {
                query += " AND " + where + "=?"
                countedArgs += 1
            }
        }
    }

    switch q.queryType {
    case QuerySelect:
        query += "SELECT "

        if len(q.columns) > 0 {
            query += q.columns[0]
            for _, col := range q.columns[1:] {
                query += ", " + col
            }
        } else {
            query += "*"
        }
        
        query += " FROM " + q.table
        buildWheres()
    case QueryInsert:
        query += "INSERT INTO " + q.table + " ("

        if len(q.columns) > 0 {
            query += q.columns[0]
            countedArgs += 1
            for _, col := range q.columns[1:] {
                query += ", " + col
                countedArgs += 1
            }
            query += ") VALUES (?"
            for range q.columns[1:] {
                query += ", ?"
            }
            query += ")"
        } else {
            return "", nil, fmt.Errorf("nothing to insert")
        }
    case QueryUpdate:
        query += "UPDATE " + q.table + " SET "

        if len(q.columns) > 0 {
            query += q.columns[0] + "=?"
            countedArgs += 1
            for _, col := range q.columns[1:] {
                query += ", " + col + "=?"
                countedArgs += 1
            }
        } else {
            return "", nil, fmt.Errorf("no columns to update")
        }
        buildWheres()
    case QueryDelete:
        query += "DELETE FROM " + q.table
        buildWheres()
    case QueryCreateTable:
        colInstances := map[string]int{}

        query = "START TRANSACTION;"

        for _, field := range q.tableModel.Fields {
            if colInstances[field.Name] > 1 {
                return "", nil, fmt.Errorf("duplicate column name %s", field.Name)
            }
            colInstances[field.Name] += 1
        }
        query += "CREATE TABLE " + q.tableModel.Name + " ("

        primaryKeys := make([]models.Field, 0)
        foreignKeys := make([]models.Field, 0)
        uniques := make([]models.Field, 0)

        for i, field := range q.tableModel.Fields {
            if i >= 1 {
                query += ", "
            }
            query += field.Name + " " + string(field.Type)
            if field.Size > 0 {
                query += "(" + strconv.Itoa(field.Size) + ")"
            }
            if !field.Nullable {
                query += " NOT NULL"
            }
            if field.Default != "" {
                query += " DEFAULT " + field.Default
            }
            if field.AutoIncrement {
                query += " AUTO_INCREMENT"
            }
            if field.Primary {
                if field.Nullable {
                    return "", nil, fmt.Errorf("nullable primary key forbidden")
                }
                primaryKeys = append(primaryKeys, field)
            }
            if field.ForeignKey != nil {
                foreignKeys = append(foreignKeys, field)
            }
            if field.Unique {
                uniques = append(uniques, field)
            }
        }
        if len(uniques) > 0 {
            query += ", CONSTRAINT uc_" + q.tableModel.Name + " UNIQUE KEY (" + uniques[0].Name
            for _, field := range uniques[1:] {
                query += ", " + field.Name
            }
            query += ")"
        }
        if len(primaryKeys) > 0 {
            query += ", CONSTRAINT pk_" + q.tableModel.Name + " PRIMARY KEY (" + primaryKeys[0].Name

            for _, field := range primaryKeys[1:] {
                query += ", " + field.Name
            }
            query += ")"
        }
        for _, field := range foreignKeys {
            query += ", CONSTRAINT fk_" + q.tableModel.Name + "_" + field.ForeignKey.Table + " FOREIGN KEY (" + field.Name
            query += ") REFERENCES " + field.ForeignKey.Table + "(" + field.ForeignKey.Column + ")"
        }

        // TODO checking if Module.String is probably not enough? check if valid? (sql.NullString.Valid)
        // TODO check if fields can be zero length, probably can?
        if q.tableModel.Id == "" || q.tableModel.Name == "" || len(q.tableModel.Fields) == 0 {
            return "", []any{}, fmt.Errorf("id or name or fields not included")
        }
        query := "INSERT INTO _tables (id, name, system"
        if q.tableModel.Module != "" {
            query += ", module"
        }
        query += ") VALUES ('"
        query += q.tableModel.Id + "', '" + q.tableModel.Name +  "', "
        if q.tableModel.System {
            query += "TRUE"
        } else {
            query += "FALSE"
        }

        // b, err := json.Marshal(q.tableModel.Fields)
        // if err != nil {
        //     return "", []any{}, err
        // }
        // need to escape the single quotes for mysql
        // query += ", '" + strings.ReplaceAll(string(b), "'", "\\'") + "'"

        if q.tableModel.Module != "" {
            query += ", '" + q.tableModel.Module + "'"
        }
        query += ")"

        // TODO add fields meta
        // columns := ?
        // vals := ?
        // query += "INSERT INTO " + q.table + " ("
        query += "INSERT INTO _fields (id, name, table, system, module, type, size, nullable, primary, default, foreignKey, unique, index, autoIncrement)"

        // query += columns[0]
        // countedArgs += 1
        // for _, col := range columns[1:] {
        //     query += ", " + col
        //     countedArgs += 1
        // }
        // query += ") VALUES (?"
        // for range columns[1:] {
        //     query += ", ?"
        // }
        
        query += ") COMMIT"
        return query, []any{}, nil
    case QueryAlterTable:
        query += "ALTER TABLE " + q.table
        switch q.alterTableQueryType {
        case QueryAlterTableAddColumn:
            query += " ADD COLUMN " + q.fieldNew.Name + " " + string(q.fieldNew.Type)
            if q.fieldNew.Size > 0 {
                query += "(" + strconv.Itoa(q.fieldNew.Size) + ")"
            }
            if !q.fieldNew.Nullable {
                query += " NOT NULL"
            }
            if q.fieldNew.Default != "" {
                query += " DEFAULT " + q.fieldNew.Default
            }
            if q.fieldNew.AutoIncrement {
                query += " AUTO_INCREMENT"
            }
        case QueryAlterTableModifyColumnDefinition:
            query += " MODIFY COLUMN " + q.fieldOld.Name + " " + string(q.fieldNew.Type)
            if q.fieldNew.Size > 0 {
                query += "(" + strconv.Itoa(q.fieldNew.Size) + ")"
            }
            if !q.fieldNew.Nullable {
                query += " NOT NULL"
            }
            if q.fieldNew.Default != "" {
                query += " DEFAULT " + q.fieldNew.Default
            }
            if q.fieldNew.AutoIncrement {
                query += " AUTO_INCREMENT"
            }
        case QueryAlterTableModifyColumnName:
            query += " RENAME COLUMN " + q.fieldOld.Name + " TO " + q.fieldNew.Name
        case QueryAlterTableDropColumn:
            query += " DROP COLUMN " + q.fieldOld.Name
        case QueryAlterTableAddIndex:
            query += " ADD INDEX idx_" + q.table + "_" + q.fieldNew.Name + " (" + q.fieldNew.Name + ")"
        case QueryAlterTableAddFK:
            query += " ADD CONSTRAINT fk_" + q.fieldNew.Name + "_" + q.fieldNew.ForeignKey.Table + " FOREIGN KEY ("
            query += q.fieldNew.Name + ") REFERENCES " + q.fieldNew.ForeignKey.Table + "(" + q.fieldNew.ForeignKey.Column + ")"
        case QueryAlterTableDropConstraint:
            query += " DROP CONSTRAINT " + q.dropConstraint
        default:
            return "", nil, fmt.Errorf("unknown query type for alter table")
        }
    case QueryRenameTable:
        query += "RENAME TABLE " + q.table + " TO " + q.tableModel.Name
    case QueryDropTable:
        query += "DROP TABLE IF EXISTS" + q.table
    // case QueryInsertTableMeta:
    //     // TODO checking if Module.String is probably not enough? check if valid? (sql.NullString.Valid)
    //     // TODO check if fields can be zero length, probably can?
    //     if q.tableModel.Id == "" || q.tableModel.Name == "" || len(q.tableModel.Fields) == 0 {
    //         return "", []any{}, fmt.Errorf("id or name or fields not included")
    //     }
    //     query := "INSERT INTO _tables (id, name, system, fields"
    //     if q.tableModel.Module != "" {
    //         query += ", module"
    //     }
    //     query += ") VALUES ('"
    //     query += q.tableModel.Id + "', '" + q.tableModel.Name +  "', "
    //     if q.tableModel.System {
    //         query += "TRUE"
    //     } else {
    //         query += "FALSE"
    //     }

    //     b, err := json.Marshal(q.tableModel.Fields)
    //     if err != nil {
    //         return "", []any{}, err
    //     }
    //     // need to escape the single quotes for mysql
    //     query += ", '" + strings.ReplaceAll(string(b), "'", "\\'") + "'"

    //     if q.tableModel.Module != "" {
    //         query += ", '" + q.tableModel.Module + "'"
    //     }
    //     query += ")"

    //     return query, []any{}, nil
    // case QueryDeleteTableMeta:
    //     return "DELETE FROM _tables WHERE name='" + q.table + "'", []any{}, nil
    default:
        return "", []any{}, fmt.Errorf("unknown query type")
    }

    if len(q.args) != countedArgs {
        return "", []any{}, fmt.Errorf("incorrect number of args, query: %s, num args: %d", query, len(q.args))
    }

    return query, q.args, nil
}
