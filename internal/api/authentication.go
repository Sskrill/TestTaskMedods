package api

import (
	"encoding/json"
	"fmt"
	"github.com/Sskrill/TestTaskMedods/internal/domain"
	"io/ioutil"
	"net"
	"net/http"
)

func (h *Handler) refreshTokens(w http.ResponseWriter, r *http.Request) {
	cookieToken, err := r.Cookie("refresh-token")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(domain.CustomErrorResponse{Message: err.Error()})
		w.Write(resp)
		return
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	aToken, rToken, err := h.service.RefreshTokens(r.Context(), cookieToken.Value, ip)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(domain.CustomErrorResponse{Message: err.Error()})
		w.Write(resp)
		return
	}
	resp, err := json.Marshal(map[string]string{"token": aToken})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(domain.CustomErrorResponse{Message: err.Error()})
		w.Write(resp)
		return
	}
	cookie := &http.Cookie{
		Name:     "refresh-token",
		Value:    fmt.Sprintf("'%s'", rToken),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	dataReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		resp, _ := json.Marshal(domain.CustomErrorResponse{Message: err.Error()})
		w.Write(resp)
		return
	}

	signInParam := domain.SignInInput{}
	if err := json.Unmarshal(dataReq, &signInParam); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(domain.CustomErrorResponse{Message: err.Error()})
		w.Write(resp)
		return
	}
	if err = signInParam.IsValid(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(domain.CustomErrorResponse{Message: err.Error()})
		w.Write(resp)
		return
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	aToken, rToken, err := h.service.SignIn(r.Context(), signInParam, ip)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(domain.CustomErrorResponse{Message: err.Error()})
		w.Write(resp)
		return
	}
	resp, err := json.Marshal(map[string]string{"token": aToken})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(domain.CustomErrorResponse{Message: err.Error()})
		w.Write(resp)
		return
	}
	cookie := &http.Cookie{
		Name:     "refresh-token",
		Value:    fmt.Sprintf("'%s'", rToken),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)

}

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	dataReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(domain.CustomErrorResponse{Message: err.Error()})

		w.Write(resp)
		return
	}
	signUpParam := domain.SignUpInput{}
	if err := json.Unmarshal(dataReq, &signUpParam); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(domain.CustomErrorResponse{Message: err.Error()})

		w.Write(resp)
		return
	}
	if err = signUpParam.IsValid(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(domain.CustomErrorResponse{Message: err.Error()})

		w.Write(resp)
		return
	}
	err = h.service.SignUp(r.Context(), signUpParam)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(domain.CustomErrorResponse{Message: err.Error()})

		w.Write(resp)
		return
	}
	w.WriteHeader(http.StatusOK)
}
