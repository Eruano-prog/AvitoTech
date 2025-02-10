package controller

import (
	"AvitoTech/internal/service"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

type ApiController struct {
	l *zap.Logger

	auth *service.AuthService
	info *service.InfoService
	coin *service.CoinService
}

func (a ApiController) Register(r chi.Router) {
	r.Post("/api/auth", a.ApiAuth)
	r.Get("/api/buy/{item}", a.ApiBuyItem)
	r.Get("/api/info", a.ApiInfo)
	r.Post("/api/sendCoin", a.ApiSendCoin)
}

func (a ApiController) ApiAuth(w http.ResponseWriter, r *http.Request) {
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
		if errors.Is(err, service.UnauthorizedError) {
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
	w.Write(jsonResp)
}

func (a ApiController) ApiBuyItem(w http.ResponseWriter, r *http.Request) {
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

	// Invoke the callback with all the unmarshaled arguments
	err = a.coin.BuyItem(id, item)
	if err != nil {
		a.writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
}

func (a ApiController) ApiInfo(w http.ResponseWriter, r *http.Request) {
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

	// Invoke the callback with all the unmarshaled arguments
	info, err := a.info.GetInfo(id)
	if err != nil {
		a.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	//TODO: move it to service layer
	sent := make([]SendRecord, len(info.Sent))
	for i, item := range info.Sent {
		sent[i].ToUser = &item.User
		sent[i].Amount = &item.Amount
	}

	received := make([]ReceiveRecord, len(info.Received))
	for i, item := range info.Received {
		received[i].FromUser = &item.User
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

func (a ApiController) ApiSendCoin(w http.ResponseWriter, r *http.Request) {
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

func (a ApiController) writeError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	errorMessage := message
	errResp := ErrorResponse{Errors: &errorMessage}
	jsonResp, _ := json.Marshal(errResp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func NewApiController(
	l *zap.Logger,
	a *service.AuthService,
	i *service.InfoService,
	c *service.CoinService,
) *ApiController {
	return &ApiController{
		l:    l,
		auth: a,
		info: i,
		coin: c,
	}
}
