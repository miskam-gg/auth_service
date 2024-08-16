package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"database/sql"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your_secret_key")

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func handleToken(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	ipAddress := r.RemoteAddr

	accessToken, err := generateAccessToken(userID, ipAddress)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	refreshToken, hashedRefreshToken, err := generateRefreshToken()
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	err = SaveRefreshToken(db, userID, hashedRefreshToken, ipAddress)
	if err != nil {
		http.Error(w, "Failed to save refresh token", http.StatusInternalServerError)
		return
	}

	response := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleRefresh(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userID := r.URL.Query().Get("user_id")
	refreshToken := r.URL.Query().Get("refresh_token")

	if userID == "" || refreshToken == "" {
		http.Error(w, "user_id and refresh_token are required", http.StatusBadRequest)
		return
	}

	tokenDetails, err := GetRefreshToken(db, userID)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(tokenDetails.TokenHash), []byte(refreshToken))
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	ipAddress := r.RemoteAddr
	if tokenDetails.IPAddress != ipAddress {
		fmt.Println("Warning: IP address has changed!")
	}

	accessToken, err := generateAccessToken(userID, ipAddress)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	response := TokenResponse{
		AccessToken: accessToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func generateAccessToken(userID, ipAddress string) (string, error) {
	claims := &jwt.StandardClaims{
		Subject:   userID,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		Audience:  ipAddress,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(jwtSecret)
}

func generateRefreshToken() (string, string, error) {
	refreshToken := uuid.New().String()
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	return refreshToken, string(hashedToken), nil
}
