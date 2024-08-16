package main

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestGenerateAccessToken(t *testing.T) {
	userID := "test_user"
	ipAddress := "192.168.1.1"
	tokenString, err := generateAccessToken(userID, ipAddress)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	assert.NoError(t, err)
	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		assert.Equal(t, userID, claims.Subject)
		assert.WithinDuration(t, time.Now(), time.Unix(claims.IssuedAt, 0), time.Minute)
		assert.WithinDuration(t, time.Now().Add(time.Minute*15), time.Unix(claims.ExpiresAt, 0), time.Minute)
		assert.Equal(t, ipAddress, claims.Audience)
	} else {
		t.Fail()
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	refreshToken, hashedToken, err := generateRefreshToken()
	assert.NoError(t, err)
	assert.NotEmpty(t, refreshToken)
	assert.NotEmpty(t, hashedToken)
	err = bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(refreshToken))
	assert.NoError(t, err)
}

func TestSaveAndGetRefreshToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	userID := "test_user"
	refreshToken := "test_refresh_token"
	ipAddress := "192.168.1.1"
	hashedToken, _ := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	mock.ExpectExec("INSERT INTO refresh_tokens").
		WithArgs(sqlmock.AnyArg(), userID, string(hashedToken), sqlmock.AnyArg(), ipAddress).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = SaveRefreshToken(db, userID, string(hashedToken), ipAddress)
	assert.NoError(t, err)
	rows := sqlmock.NewRows([]string{"id", "user_id", "token_hash", "created_at", "ip_address"}).
		AddRow("test_id", userID, string(hashedToken), time.Now(), ipAddress)
	mock.ExpectQuery("SELECT id, user_id, token_hash, created_at, ip_address FROM refresh_tokens WHERE user_id = ?").
		WithArgs(userID).
		WillReturnRows(rows)
	tokenDetails, err := GetRefreshToken(db, userID)
	assert.NoError(t, err)
	assert.Equal(t, string(hashedToken), tokenDetails.TokenHash)
	assert.Equal(t, ipAddress, tokenDetails.IPAddress)
}
