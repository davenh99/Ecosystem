package handlers

import (
	"apps/ecosystem/tools"
	"apps/ecosystem/tools/types"
	"net/http"

	"github.com/go-chi/chi"
)

type RecordHandler struct {
	recordStore types.RecordStore
}

func NewRecordHandler(recordStore types.RecordStore) *RecordHandler {
	return &RecordHandler{recordStore}
}

func (h *RecordHandler) RegisterRoutes(router *chi.Mux) {
	recordsRouter := chi.NewRouter()
	router.Mount("/tables/{table}/records", recordsRouter)

	recordsRouter.Get("/{id}", h.handleGetOne)
	recordsRouter.Get("/", h.handleGetList)
	recordsRouter.Post("/", h.handleCreate)
	recordsRouter.Put("/{id}", h.handleUpdate)
	recordsRouter.Delete("/{id}", h.handleDelete)
}

// TODO figure out better way than setting table name every time
func (h *RecordHandler) handleGetOne(w http.ResponseWriter, r *http.Request) {
	h.recordStore.SetTableName(chi.URLParam(r, "table"))

	record, err := h.recordStore.GetByID(chi.URLParam(r, "id"))
	if err != nil {
		tools.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	tools.WriteJSON(w, http.StatusOK, record)
}

func (h *RecordHandler) handleGetList(w http.ResponseWriter, r *http.Request) {
	h.recordStore.SetTableName(chi.URLParam(r, "table"))

	records, err := h.recordStore.GetList(r.Context())
	if err != nil {
		tools.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	tools.WriteJSON(w, http.StatusOK, records)
}

func (h *RecordHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	h.recordStore.SetTableName(chi.URLParam(r, "table"))

	var payload map[string]any

	// parse the payload
	if err := tools.ParseJSON(r, &payload); err != nil {
		tools.WriteError(w, http.StatusBadRequest, err)
		return
	}

	id, err := h.recordStore.Create(payload)
	if err != nil {
		tools.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	tools.WriteJSON(w, http.StatusOK, id)
}

func (h *RecordHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	h.recordStore.SetTableName(chi.URLParam(r, "table"))

	var payload map[string]any

	if err := tools.ParseJSON(r, &payload); err != nil {
		tools.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err := h.recordStore.Update(chi.URLParam(r, "id"), payload)
	if err != nil {
		tools.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	tools.WriteJSON(w, http.StatusOK, nil)
}

func (h *RecordHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	h.recordStore.SetTableName(chi.URLParam(r, "table"))

	err := h.recordStore.Delete(chi.URLParam(r, "id"))
	if err != nil {
		tools.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	tools.WriteJSON(w, http.StatusOK, nil)
}
