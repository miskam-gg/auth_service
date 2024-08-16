// server.go

package main

import (
	"database/sql"
	"log"
	"net/http"
)

func setupRoutes(db *sql.DB) {
	http.HandleFunc("/auth/token", func(w http.ResponseWriter, r *http.Request) {
		handleToken(w, r, db)
	})
	http.HandleFunc("/auth/refresh", func(w http.ResponseWriter, r *http.Request) {
		handleRefresh(w, r, db)
	})

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
