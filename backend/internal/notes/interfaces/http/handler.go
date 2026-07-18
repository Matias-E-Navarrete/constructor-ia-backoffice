package http

import (
	nethttp "net/http"

	"github.com/go-chi/chi/v5"

	"rimu/backend/internal/notes/application"
	"rimu/backend/internal/notes/domain"
	"rimu/backend/internal/platform/httpkit"
	"rimu/backend/internal/platform/httpmiddleware"
)

type Handler struct {
	Create     *application.CreateNote
	Update     *application.UpdateNote
	Get        *application.GetNote
	List       *application.ListNotes
	Delete     *application.DeleteNote
	GetGraph   *application.GetVaultGraph
	ExportVault *application.ExportVault
}

func (h *Handler) HandleCreate(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())

	var req createNoteRequest
	if err := httpkit.DecodeJSON(r, &req); err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_body", "malformed JSON body")
		return
	}

	n, err := h.Create.Execute(r.Context(), application.CreateNoteInput{
		UserID: user.ID, Title: req.Title, Body: req.Body, Kind: domain.Kind(req.Kind),
	})
	if err == domain.ErrSlugTaken {
		httpkit.WriteError(w, nethttp.StatusConflict, "slug_taken", err.Error())
		return
	}
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusCreated, toNoteResponse(*n))
}

func (h *Handler) HandleList(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	notes, err := h.List.Execute(r.Context(), user.ID)
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	resp := make([]noteResponse, 0, len(notes))
	for _, n := range notes {
		resp = append(resp, toNoteResponse(n))
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, resp)
}

func (h *Handler) HandleGet(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	slug := chi.URLParam(r, "slug")

	n, err := h.Get.Execute(r.Context(), user.ID, slug)
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusNotFound, "not_found", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, toNoteResponse(*n))
}

func (h *Handler) HandleUpdate(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	slug := chi.URLParam(r, "slug")

	var req updateNoteRequest
	if err := httpkit.DecodeJSON(r, &req); err != nil {
		httpkit.WriteError(w, nethttp.StatusBadRequest, "invalid_body", "malformed JSON body")
		return
	}

	n, err := h.Update.Execute(r.Context(), application.UpdateNoteInput{
		UserID: user.ID, Slug: slug, Title: req.Title, Body: req.Body,
	})
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusNotFound, "not_found", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, toNoteResponse(*n))
}

func (h *Handler) HandleDelete(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	slug := chi.URLParam(r, "slug")

	if err := h.Delete.Execute(r.Context(), user.ID, slug); err != nil {
		httpkit.WriteError(w, nethttp.StatusNotFound, "not_found", err.Error())
		return
	}
	w.WriteHeader(nethttp.StatusNoContent)
}

func (h *Handler) HandleGraph(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	graph, err := h.GetGraph.Execute(r.Context(), user.ID)
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	httpkit.WriteJSON(w, nethttp.StatusOK, toGraphResponse(graph))
}

func (h *Handler) HandleExport(w nethttp.ResponseWriter, r *nethttp.Request) {
	user, _ := httpmiddleware.UserFromContext(r.Context())
	zipBytes, err := h.ExportVault.Execute(r.Context(), user.ID)
	if err != nil {
		httpkit.WriteError(w, nethttp.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=vault.zip")
	w.WriteHeader(nethttp.StatusOK)
	w.Write(zipBytes)
}
