package types

import (
	"apps/ecosystem/core/models"
	"context"
)

type BaseStore[T any] interface {
	GetByID(id string) (*T, error)
	GetList(ctx context.Context) ([]T, error)
	Update(id string, payload map[string]any) error
	Delete(id string) error
}

type TableStore interface {
	SetTableName(name string)
	GetByID(id string) (*models.TableModel, error)
	GetList() ([]models.TableModel, error)
	Create(payload models.TableModel) (string, error)
	Update(id string, payload models.TableModel) error
	Delete() error
	GetByName(tableName string) (*models.TableModel, error)
}

type RecordStore interface {
	BaseStore[map[string]any]
	Create(payload map[string]any) (string, error)
	SetTableName(name string)
}

type UserStore interface {
	BaseStore[models.UserModel]
	Create(payload models.AuthModel) (string, error)
	GetAuthByEmail(email string) (*models.AuthModel, error)
}

type RoleStore interface {
	BaseStore[models.RoleModel]
	Create(payload models.RoleModel) (string, error)
	AssignRoleToUser(userId string, roleId string) error
}
