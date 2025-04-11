package handlers

import (
	"apps/ecosystem/core/models"
	"apps/ecosystem/tools/types"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
)

func TestUserServiceHandlers(t *testing.T) {
	userStore := &mockUserStore{}
	roleStore := &mockRoleStore{}
	handler := NewUserHandler(userStore, roleStore)
	validPayload := types.UserRegisterPayload{
		FirstName: "fn",
		LastName: "ln",
		Email: "valid@email.com",
		Password: "12345678",
	}
	// var createdAdminId int

	t.Run("should fail if the user payload is invalid", func (t *testing.T) {
		payload := types.UserRegisterPayload{
			FirstName: "fn",
			LastName: "validlastname",
			Email: "invalid",
			Password: "12345678",
		}
		marshalled, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/user/register", bytes.NewBuffer(marshalled))

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := chi.NewRouter()
		router.HandleFunc("/user/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should correctly register the user", func (t *testing.T) {
		marshalled, _ := json.Marshal(validPayload)
		req, err := http.NewRequest(http.MethodPost, "/user/register", bytes.NewBuffer(marshalled))

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := chi.NewRouter()
		router.HandleFunc("/user/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d. %s", http.StatusCreated, rr.Code, rr.Body)
		}
	})

	// t.Run("should login the admin", func (t *testing.T) {

	// })

	// t.Run("should update the admin", func (t *testing.T) {

	// })

	// t.Run("should fail to delete the admin with un-authenticated request", func (t *testing.T) {

	// })

	// t.Run("should delete the admin with authenticated request", func (t *testing.T) {

	// })
}

type mockUserStore struct {}

func (m *mockUserStore) GetAuthByEmail(email string) (*models.AuthModel, error) {
	return nil, fmt.Errorf("user not found")
}
func (m *mockUserStore) GetByID(id string) (*models.UserModel, error) {
	return new(models.UserModel), nil
}
func (m *mockUserStore) Create(models.AuthModel) (string, error) {
	return "", nil
}
func (m *mockUserStore) GetList(ctx context.Context) ([]models.UserModel, error) {
	return []models.UserModel{}, nil
}
func (m *mockUserStore) Update(string, map[string]any) error {
	return nil
}
func (m *mockUserStore) Delete(string) error {
	return nil
}

type mockRoleStore struct {}

func (m *mockRoleStore) GetByID(id string) (*models.RoleModel, error) {
	return new(models.RoleModel), nil
}
func (m *mockRoleStore) Create(models.RoleModel) (string, error) {
	return "", nil
}
func (m *mockRoleStore) GetList(ctx context.Context) ([]models.RoleModel, error) {
	return []models.RoleModel{}, nil
}
func (m *mockRoleStore) Update(string, map[string]any) error {
	return nil
}
func (m *mockRoleStore) Delete(string) error {
	return nil
}
func (m *mockRoleStore) AssignRoleToUser(string, string) error {
	return nil
}
