// Package controller
package controller

import (
	"AvitoTech/internal/service"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

type APIController struct {
	l *zap.Logger

	auth service.Auth
	info service.Info
	coin service.Coin
}

func (a APIController) Register(r chi.Router) {
	r.Post("/api/auth", a.apiAuth)
	r.Get("/api/buy/{item}", a.apiBuyItem)
	r.Get("/api/info", a.apiInfo)
	r.Post("/api/sendCoin", a.apiSendCoin)
}

func (a APIController) apiAuth(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.writeError(w, http.StatusBadRequest, "Invalid request: missing username or password")
		return
	}

	var req AuthRequest
	err = json.Unmarshal(body, &req)
	if err != nil || req.Username == "" || req.Password == "" {
		a.writeError(w, http.StatusBadRequest, "Invalid request: missing username or password")
		return
	}

	token, err := a.auth.Authenticate(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrUnauthorized) {
			a.writeError(w, http.StatusUnauthorized, "User unauthorized")
			return
		}
		a.writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	resp := AuthResponse{Token: &token}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		a.writeError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResp)
	if err != nil {
		a.l.Error("Failed to write response", zap.Error(err))
		return
	}
}

func (a APIController) apiBuyItem(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		a.writeError(w, http.StatusBadRequest, "Missing token")
		return
	}
	id, err := a.auth.VerifyJWT(token)
	if err != nil {
		a.writeError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	item := chi.URLParam(r, "item")
	if item == "" {
		a.writeError(w, http.StatusBadRequest, "Item can't be empty")
	}

	err = a.coin.BuyItem(id, item)
	if err != nil {
		a.writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
}

func (a APIController) apiInfo(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		a.writeError(w, http.StatusBadRequest, "Missing token")
		return
	}
	id, err := a.auth.VerifyJWT(token)
	if err != nil {
		a.writeError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	info, err := a.info.GetInfo(id)
	if err != nil {
		a.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sent := make([]SendRecord, len(info.Sent))
	for i, item := range info.Sent {
		sent[i].ToUser = &item.ToUser
		sent[i].Amount = &item.Amount
	}

	received := make([]ReceiveRecord, len(info.Received))
	for i, item := range info.Received {
		received[i].FromUser = &item.FromUser
		received[i].Amount = &item.Amount
	}

	inventory := make([]InventoryRecord, 0, len(info.Inventory))
	for key, value := range info.Inventory {
		inventory = append(inventory, InventoryRecord{Quantity: &value, Type: &key})
	}
	resp := InfoResponse{
		CoinHistory: &History{Sent: &sent, Received: &received},
		Coins:       &info.Coins,
		Inventory:   &inventory,
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		a.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResp)
	if err != nil {
		a.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (a APIController) apiSendCoin(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		a.writeError(w, http.StatusBadRequest, "Missing token")
		return
	}
	id, err := a.auth.VerifyJWT(token)
	if err != nil {
		a.writeError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.writeError(w, http.StatusBadRequest, "Invalid request: missing username or password")
		return
	}

	var req SendCoinRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		a.writeError(w, http.StatusBadRequest, "Invalid request: missing username or password")
		return
	}

	if req.ToUser == "" {
		a.writeError(w, http.StatusBadRequest, "User is empty")
		return
	}

	err = a.coin.SendCoin(id, req.ToUser, req.Amount)
	if err != nil {
		a.writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
}

func (a APIController) writeError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	errorMessage := message
	errResp := ErrorResponse{Errors: &errorMessage}
	jsonResp, _ := json.Marshal(errResp)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(jsonResp)
	if err != nil {
		a.l.Error("Failed to write response", zap.Error(err))
		return
	}
}

func NewAPIController(
	l *zap.Logger,
	a service.Auth,
	i service.Info,
	c service.Coin,
) *APIController {
	return &APIController{
		l:    l,
		auth: a,
		info: i,
		coin: c,
	}
}
