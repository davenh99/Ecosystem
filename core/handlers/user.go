package handlers

import (
	"apps/ecosystem/core/models"
	"apps/ecosystem/tools"
	"apps/ecosystem/tools/auth"
	"apps/ecosystem/tools/types"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userStore types.UserStore
	roleStore types.RoleStore
}

func NewUserHandler(userStore types.UserStore, roleStore types.RoleStore) *UserHandler {
	return &UserHandler {userStore, roleStore}
}

func (h *UserHandler) RegisterRoutes(router *chi.Mux) {
	router.Get("/user/", h.handleGetList)
	router.Post("/user/authenticate", h.handleLogin)
	// should below be get or post? we are sending the authorization and cookie?
	router.Get("/user/{id}", auth.WithJWTAuth(h.handleGetOne, h.userStore))
	router.Post("/user/register", h.handleRegister)
	// router.Post("/user/refresh", h.handleRefresh)
	router.Put("/user/{id}", h.handleUpdate)
	router.Delete("/user/{id}", h.handleDelete)
}

func (h *UserHandler) handleGetList(w http.ResponseWriter, r *http.Request) {
	users, err := h.userStore.GetList(r.Context())
	if err != nil {
		tools.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	tools.WriteJSON(w, http.StatusOK, users)
}

func (h *UserHandler) handleGetOne(w http.ResponseWriter, r *http.Request) {
	// if auth failed with access token, try to issue new one and refresh token.
	userId := auth.GetUserIdFromContext(r.Context())
	tokenAccess := auth.GetJWTAccessFromContext(r.Context())

	user, err := h.userStore.GetByID(userId)
	if err != nil {
		tools.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	
	// TODO is this the new token that was generated from withjwtauth?
	// return the user view model
	tools.WriteJSON(w, http.StatusOK, types.UserAuthResponse{Token: tokenAccess, User: *user})
}

func (h *UserHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.UserLoginPayload

	// parse the payload
	if err := tools.ParseJSON(r, &payload); err != nil {
		tools.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate the payload
	if err := tools.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		// TODO need better error handling here, maybe just better message
		tools.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	u, err := h.userStore.GetAuthByEmail(payload.Email)
	if err != nil {
		tools.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to authenticate, email or password is incorrect"))
		return
	}

	if !auth.CheckPassword(u.Password, payload.Password) {
		tools.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to authenticate, email or password is incorrect"))
		return
	}

	// TODO consider abstracting out the below to do with JWTs and cookies

	tokenAccess, err := auth.NewUserAccessJWT(u.Id)
	if err != nil {
		tools.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	tokenRefresh, err := auth.NewUserRefreshJWT(u.Id)
	if err != nil {
		tools.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	user, err := h.userStore.GetByID(u.Id)
	if err != nil {
		tools.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	tools.SetCookie(w, tokenRefresh)

	tools.WriteJSON(w, http.StatusOK, types.UserAuthResponse{Token: tokenAccess, User: *user})
}

func (h *UserHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.UserRegisterPayload

	// get JSON payload
	// check if user exists
	// if it doesn't, create new one

	// parse the payload
	if err := tools.ParseJSON(r, &payload); err != nil {
		tools.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate the payload
	if err := tools.Validate.Struct(payload); err != nil {
		tools.WriteValidationError(w, http.StatusBadRequest, err.(validator.ValidationErrors))
		return
	}
	
	// check if user exists
	_, err := h.userStore.GetAuthByEmail(payload.Email)
	if err == nil {
		tools.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}
	
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		tools.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// get list of users to check if we need to init the thingy
	users, err := h.userStore.GetList(r.Context())
	if err != nil {
		tools.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// if no user, create one
	id, err := h.userStore.Create(models.AuthModel{
		FirstName: payload.FirstName,
		LastName: payload.LastName,
		Email: payload.Email,
		Password: hashedPassword,
	})
	
	if err != nil {
		tools.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	
	// if user is first user, assign role as owner and init roles
	if len(users) < 1 {
		// TODO extract this to some sort of init function for db, pb calls it 'installer'
		roleId, err := h.roleStore.Create(models.RoleModel{Name: "owner"})
		if err != nil {
			tools.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		_, err = h.roleStore.Create(models.RoleModel{Name: "admin"})
		if err != nil {
			tools.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		err = h.roleStore.AssignRoleToUser(id, roleId)
		if err != nil {
			tools.WriteError(w, http.StatusInternalServerError, err)
			return
		}
	}

	tokenAccess, err := auth.NewUserAccessJWT(id)
	if err != nil {
		tools.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	tokenRefresh, err := auth.NewUserRefreshJWT(id)
	if err != nil {
		tools.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// TODO is it bad performance to query a second time? is there something better?
	user, err := h.userStore.GetByID(id)
	if err != nil {
		tools.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	tools.SetCookie(w, tokenRefresh)

	// TODO try to send the 'Token' and 'User' as lowercase - low priority
	tools.WriteJSON(w, http.StatusCreated, types.UserAuthResponse{Token: tokenAccess, User: *user})
}

// func (h *Handler) handleRefresh(w http.ResponseWriter, r *http.Request) {
// 	refreshToken, err := tools.GetCookie(r)
// 	// TODO this error handling looks nice! consider implementing it elsewhere? throught the app?
// 	if err != nil {
// 		switch {
// 		case goErrors.Is(err, http.ErrNoCookie):
// 			tools.WriteError(w, http.StatusBadRequest, fmt.Errorf("cookie not found"))
// 		default:
// 			tools.WriteError(w, http.StatusInternalServerError, err)
// 			log.Println(err)
// 		}
// 	}

// 	// validate the JWT
// 	token, err := auth.ValidateToken(refreshToken)
// 	if err != nil {
// 		log.Printf("failed to validate token: %v", err)
// 		errors.PermissionDenied(w)
// 		return
// 	}

// 	if !token.Valid {
// 		log.Printf("invalid token")
// 		errors.PermissionDenied(w)
// 		return
// 	}

// 	claims := token.Claims.(jwt.MapClaims)
// 	str := claims["userId"].(string)

// 	userId, _ := strconv.Atoi(str)
// 	u, err := h.userStore.GetByID(userId)
// 	if err != nil {
// 		log.Printf("failed to get user by id: %v", err)
// 		errors.PermissionDenied(w)
// 		return
// 	}

// 	tokenAccess, err := auth.NewUserAccessJWT(u.Id)
// 	if err != nil {
// 		tools.WriteError(w, http.StatusInternalServerError, err)
// 		return
// 	}

// 	tokenRefresh, err := auth.NewUserRefreshJWT(u.Id)
// 	if err != nil {
// 		tools.WriteError(w, http.StatusInternalServerError, err)
// 		return
// 	}

// 	tools.SetCookie(w, tokenRefresh)

// 	tools.WriteJSON(w, http.StatusOK, types.UserAuthResponse{
// 		Token: tokenAccess, User: types.UserViewModel{
// 			Email: u.Email,
// 			Username: u.Username,
// 			Created: u.Created,
// 			Updated: u.Updated,
// 		},
// 	})
// }

func (h *UserHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	var payload map[string]any

	userId := auth.GetUserIdFromContext(r.Context())

	// get JSON payload
	// check if user exists
	// if they doesn't, create new one

	// parse the payload, not sure how this works for map[string]any?
	if err := tools.ParseJSON(r, &payload); err != nil {
		tools.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate the payload
	// if err := tools.Validate.Struct(payload); err != nil {
	// 	errors := err.(validator.ValidationErrors)
	// 	tools.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
	// 	return
	// }
	
	// TODO use jwt auth for this!!!
	err := h.userStore.Update(userId, payload)
	if err != nil {
		tools.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	tools.WriteJSON(w, http.StatusOK, nil)
}

func (h *UserHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	err := h.userStore.Delete(chi.URLParam(r, "id"))
	if err != nil {
		tools.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	tools.WriteJSON(w, http.StatusOK, nil)
}
