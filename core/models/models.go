package models

import (
	"time"
)

// TODO probs add method on tablemodel to convert fields to json?

type FieldType string

const (
	FieldBoolean FieldType = "BOOLEAN"
	FieldChar FieldType = "CHAR"
	FieldVarchar FieldType = "VARCHAR"
	FieldText FieldType = "TEXT"
	FieldDatetime FieldType = "DATETIME"
	FieldTimestamp FieldType = "TIMESTAMP"
	FieldInt FieldType = "INT"
	FieldJson FieldType = "JSON"
	FieldFloat FieldType = "FLOAT"
)

type ForeignKey struct {
	Table  string `json:"table"`
	Column string `json:"column"`
}

type Field struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Type  FieldType `json:"type"`
	Size int `json:"size"`
	Nullable bool `json:"nullable"`
	Primary bool `json:"primary"`
	Default string `json:"default"`
	ForeignKey *ForeignKey `json:"foreignKey"`
	Unique bool `json:"unique"`
	Index bool `json:"index"`
	AutoIncrement bool `json:"autoIncrement"`
}

type TableModel struct {
	Id string `json:"id"`
	System bool `json:"system"`
	Module string `json:"module"`
	Name string `json:"name"`
	Fields []Field `json:"fields"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type RoleModel struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type UserRoleModel struct {
	UserId string `json:"userId"`
	RoleId string `json:"roleId"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type AuthModel struct {
	Id string `json:"id"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Email string `json:"email"`
	Password string `json:"-"`
}

type UserModel struct {
	Id string `json:"id"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Email string `json:"email"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type ViewModel struct {
	Id string `json:"id"`
	Name string `json:"name"`
	View string `json:"view"`
	ViewJSON string `json:"viewJSON"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type ScriptModel struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Script string `json:"script"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}
