package api

import (
	"encoding/json"
	"github.com/Sskrill/TestTaskMedods/internal/domain"
	"github.com/gorilla/mux"
	"net/http"
)

func (h *Handler) getTokens(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid := vars["guid"]

	aT, rT, err := h.service.GetTokensByGUID(r.Context(), guid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(domain.CustomErrorResponse{Message: err.Error()})
		w.Write(resp)
		return
	}

	resp, err := json.Marshal(domain.CooupleOfTokens{AToken: aT, RToken: rT})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ = json.Marshal(domain.CustomErrorResponse{Message: err.Error()})
		w.Write(resp)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)
}
