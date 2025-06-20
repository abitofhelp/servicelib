// Copyright (c) 2025 A Bit of Help, Inc.

// Example usage of the Auth module for a quick start
package main

import (
	"context"
	"net/http"

	"github.com/abitofhelp/servicelib/auth"
	"go.uber.org/zap"
)

func main() {
	// Create a logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Create a context
	ctx := context.Background()

	// Create a configuration
	config := auth.DefaultConfig()
	config.JWT.SecretKey = "your-secret-key"

	// Create an auth instance
	authInstance, err := auth.New(ctx, config, logger)
	if err != nil {
		logger.Fatal("Failed to create auth instance", zap.Error(err))
	}

	// Create an HTTP handler
	http.Handle("/", authInstance.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the user is authorized to perform an operation
		authorized, err := authInstance.IsAuthorized(r.Context(), "read:resource")
		if err != nil {
			http.Error(w, "Authorization error", http.StatusInternalServerError)
			return
		}

		if !authorized {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Get the user ID
		userID, err := authInstance.GetUserID(r.Context())
		if err != nil {
			http.Error(w, "User ID not found", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Hello, " + userID))
	})))

	// Start the server
	http.ListenAndServe(":8080", nil)
}
