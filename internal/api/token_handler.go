package api

import (
	"encoding/json"
	"github.com/pistolricks/m-api/internal/middleware"
	"log"
	"net/http"
	"time"

	"github.com/pistolricks/m-api/internal/store"
	"github.com/pistolricks/m-api/internal/tokens"
	"github.com/pistolricks/m-api/internal/utils"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type createTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type deleteTokensRequest struct {
	Username string `json:"username"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func (h *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		h.logger.Printf("ERROR: createTokenRequest: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	// lets get the user
	user, err := h.userStore.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		h.logger.Printf("ERROR: GetUserByUsername: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	passwordsDoMatch, err := user.PasswordHash.Matches(req.Password)
	if err != nil {
		h.logger.Printf("ERORR: PasswordHash.Mathes %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if !passwordsDoMatch {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
		return
	}

	token, err := h.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		h.logger.Printf("ERORR: Creating Token %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return

	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"auth_token": token, "user": user})

}

func (h *TokenHandler) HandleDeleteAllUserTokens(w http.ResponseWriter, r *http.Request) {

	// lets get the user
	user := middleware.GetUser(r)

	err := h.tokenStore.DeleteAllTokensForUser(user.ID, tokens.ScopeAuth)
	if err != nil {
		h.logger.Printf("ERORR: Creating Token %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return

	}

	utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"user": store.AnonymousUser, "message": "Logged Out"})

}
