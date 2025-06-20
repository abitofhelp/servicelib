// Copyright (c) 2025 A Bit of Help, Inc.

// Example usage of the Auth module configuration
package main

import (
    "fmt"
    "time"

    "github.com/abitofhelp/servicelib/auth"
)

func main() {
    // Create a configuration
    config := auth.Config{}

    // JWT configuration
    config.JWT.SecretKey = "your-secret-key"
    config.JWT.TokenDuration = 24 * time.Hour
    config.JWT.Issuer = "your-issuer"

    // JWT Remote validation configuration
    config.JWT.Remote.Enabled = true
    config.JWT.Remote.ValidationURL = "https://your-auth-server.com/validate"
    config.JWT.Remote.ClientID = "your-client-id"
    config.JWT.Remote.ClientSecret = "your-client-secret"
    config.JWT.Remote.Timeout = 5 * time.Second

    // OIDC configuration
    config.OIDC.IssuerURL = "https://your-oidc-provider.com"
    config.OIDC.ClientID = "your-client-id"
    config.OIDC.ClientSecret = "your-client-secret"
    config.OIDC.RedirectURL = "https://your-app.com/callback"
    config.OIDC.Scopes = []string{"openid", "profile", "email"}
    config.OIDC.Timeout = 10 * time.Second

    // Middleware configuration
    config.Middleware.SkipPaths = []string{"/public", "/health"}
    config.Middleware.RequireAuth = true

    // Print the configuration
    fmt.Printf("JWT Secret Key: %s\n", config.JWT.SecretKey)
    fmt.Printf("JWT Token Duration: %v\n", config.JWT.TokenDuration)
    fmt.Printf("JWT Issuer: %s\n", config.JWT.Issuer)
    fmt.Printf("OIDC Issuer URL: %s\n", config.OIDC.IssuerURL)
    fmt.Printf("OIDC Client ID: %s\n", config.OIDC.ClientID)
    fmt.Printf("Middleware Skip Paths: %v\n", config.Middleware.SkipPaths)
}