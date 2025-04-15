package handlers

import (
	"apps/ecosystem/core/models"
	"apps/ecosystem/tools"
	"apps/ecosystem/tools/types"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

type TableHandler struct {
	tableStore types.TableStore
}

func NewTableHandler(tableStore types.TableStore) *TableHandler {
	return &TableHandler{tableStore}
}

func (h *TableHandler) RegisterRoutes(router *chi.Mux) {
	tablesRouter := chi.NewRouter()
	router.Mount("/tables", tablesRouter)

	tablesRouter.Get("/", h.handleGetList)
	tablesRouter.Post("/", h.handleCreate)
	tablesRouter.Get("/{table}", h.handleGetOne)
	tablesRouter.Put("/{table}", h.handleUpdate)
	tablesRouter.Delete("/{table}", h.handleDelete)
}

func (h *TableHandler) handleGetList(w http.ResponseWriter, r *http.Request) {
	tables, err := h.tableStore.GetList()
	if err != nil {
		tools.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	tools.WriteJSON(w, http.StatusOK, tables)
}

func (h *TableHandler) handleGetOne(w http.ResponseWriter, r *http.Request) {
	t, err := h.tableStore.GetByName(chi.URLParam(r, "table"))
	if err != nil {
		tools.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	tools.WriteJSON(w, http.StatusOK, &t)
}

func (h *TableHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	h.tableStore.SetTableName(chi.URLParam(r, "table"))

	var payload models.TableModel

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

	// check if table exists first (maybe here is not the best place to do this, should I just let
	// us get a sql error?)
	_, err := h.tableStore.GetByName(chi.URLParam(r, "table"))
	if err == nil {
		tools.WriteError(w, http.StatusBadRequest, fmt.Errorf("table %s already exists", payload.Name))
		return
	}

	id, err := h.tableStore.Create(payload)
	if err != nil {
		tools.WriteError(w, http.StatusInternalServerError, err)
	}

	// maybe return the whole table here?
	tools.WriteJSON(w, http.StatusCreated, id)
}

func (h *TableHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	h.tableStore.SetTableName(chi.URLParam(r, "table"))

	var payload models.TableModel

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

	err := h.tableStore.Update(payload.Id, payload)
	if err != nil {
		tools.WriteError(w, http.StatusInternalServerError, err)
	}

	tools.WriteJSON(w, http.StatusOK, nil)
}

func (h *TableHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	// TODO below seems wrong, maybe pass in the table name to delete instead?
	// consider for all other stuff using this 'SetTableName' as well
	
	h.tableStore.SetTableName(chi.URLParam(r, "table"))

	err := h.tableStore.Delete()
	if err != nil {
		tools.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	tools.WriteJSON(w, http.StatusOK, nil)
}