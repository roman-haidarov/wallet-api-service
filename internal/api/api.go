package api

import (
	"context"
	"encoding/json"
	// "fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"wallet-api-service/internal/config"
	"wallet-api-service/internal/db"
	"wallet-api-service/internal/kafka"
	"wallet-api-service/internal/types"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const(
	topic = "my-topic"
)

type API struct {
	cfg    config.Config
	db     db.DB
	kafka  *kafka.Producer
	server *http.Server
}

func New(cfg config.Config, database db.DB, kafkaClient *kafka.Producer) *API {
	return &API{
		cfg:   cfg,
		db:    database,
		kafka: kafkaClient,
	}
}

func (a *API) Serve(ctx context.Context) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/wallets", a.walletsHandler)
	mux.HandleFunc("/wallets/", a.walletsHandler)

	a.server = &http.Server{
		Addr:    ":" + strconv.Itoa(a.cfg.Server.ListenPort),
		Handler: mux,
	}

	log.Info().Msgf("Starting HTTP server on port %d", a.cfg.Server.ListenPort)
	return a.server.ListenAndServe()
}

func (a *API) walletsHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/wallets" {
		path = "/wallets/"
	}

	path = strings.TrimPrefix(path, "/wallets/")
	
	log.Debug().
			Str("path", path).
			Str("original_path", r.URL.Path).
			Str("method", r.Method).
			Msg("Request received")

	if path == "" {
		if r.Method == http.MethodPost {
			a.createWalletHandler(w, r)
			return
		}
		log.Debug().Msg("Method not allowed for wallet creation")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(path, "/")
	
	if len(parts) < 1 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	
	walletIDStr := parts[0]
	
	if len(parts) == 1 {
		if r.Method == http.MethodGet {
			a.getWalletHandler(w, r, walletIDStr)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if len(parts) == 2 && parts[1] == "topup" {
		if r.Method == http.MethodPost {
			a.topUpWalletHandler(w, r, walletIDStr)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	http.Error(w, "Not found", http.StatusNotFound)
}

func (a *API) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}

func (a *API) getWalletHandler(w http.ResponseWriter, r *http.Request, idStr string) {
	walletID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid wallet ID", http.StatusBadRequest)
		return
	}

	wallet, err := a.db.GetWallet(r.Context(), walletID)
	if err != nil {
		http.Error(w, "Wallet not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wallet)
}

func (a *API) topUpWalletHandler(w http.ResponseWriter, r *http.Request, idStr string) {
	walletID, err := uuid.Parse(idStr)
	// _, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid wallet ID", http.StatusBadRequest)
		return
	}

	var req types.Wallet
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "Amount must be positive", http.StatusBadRequest)
		return
	}

	req.WalletID = walletID
	jsonData, err := json.Marshal(&req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal message to JSON")
		http.Error(w, "Failed to process top-up", http.StatusInternalServerError)
		return
	}

	// fmt.Println(req)
	if err := a.kafka.Produce(string(jsonData), topic); err != nil {
		log.Error().Err(err).Msg("Failed to publish to Kafka")
		http.Error(w, "Failed to process top-up", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"message": "Top-up request sent successfully",
	})
}

func (a *API) createWalletHandler(w http.ResponseWriter, r *http.Request) {
	walletID := uuid.New()
	
	wallet := types.Wallet{
		WalletID:  walletID,
		Amount:    0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := a.db.CreateWallet(r.Context(), wallet); err != nil {
		http.Error(w, "Failed to create wallet", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(wallet)
}
