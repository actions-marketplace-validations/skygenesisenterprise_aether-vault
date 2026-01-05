package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/skygenesisenterprise/aether-vault"
)

func main() {
	// Configuration example
	cfg := vault.NewConfig("http://localhost:8080", "your-token-here")
	cfg.Timeout = 10 * time.Second
	cfg.RetryCount = 3
	cfg.Debug = true

	// Create vault client
	client, err := vault.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create vault client: %v", err)
	}
	defer client.Close()

	// Check vault health
	ctx := context.Background()
	if err := client.Health(ctx); err != nil {
		log.Printf("Vault health check failed: %v", err)
	} else {
		fmt.Println("Vault is healthy")
	}

	// Get version
	version, err := client.Version(ctx)
	if err != nil {
		log.Printf("Failed to get vault version: %v", err)
	} else {
		fmt.Printf("Vault version: %s\n", version)
	}

	// Example: Create a secret
	secret, err := client.Secrets.Create(ctx, &vault.CreateSecretRequest{
		Name:        "example/database/password",
		Value:       "super-secret-password",
		Description: "Database password for example service",
		Tags: map[string]string{
			"environment": "dev",
			"service":     "example",
		},
	})
	if err != nil {
		log.Printf("Failed to create secret: %v", err)
	} else {
		fmt.Printf("Created secret: %s (ID: %s)\n", secret.Name, secret.ID)
	}

	// Example: Get a secret
	retrievedSecret, err := client.Secrets.Get(ctx, "example/database/password")
	if err != nil {
		log.Printf("Failed to get secret: %v", err)
	} else {
		fmt.Printf("Retrieved secret: %s = %s\n", retrievedSecret.Name, retrievedSecret.Value)
	}

	// Example: List secrets
	secretsList, err := client.Secrets.List(ctx, &vault.ListSecretsRequest{
		Prefix: "example/",
		Limit:  10,
	})
	if err != nil {
		log.Printf("Failed to list secrets: %v", err)
	} else {
		fmt.Printf("Found %d secrets:\n", secretsList.Total)
		for _, s := range secretsList.Secrets {
			fmt.Printf("  - %s (v%d)\n", s.Name, s.Version)
		}
	}

	// Example: Generate TOTP
	totpResp, err := client.TOTP.Generate(ctx, &vault.GenerateRequest{
		AccountName: "user@example.com",
		Issuer:      "Aether Vault",
	})
	if err != nil {
		log.Printf("Failed to generate TOTP: %v", err)
	} else {
		fmt.Printf("Generated TOTP secret for %s\n", totpResp.AccountName)
		fmt.Printf("QR Code URL: %s\n", totpResp.QRCode)
	}

	// Example: Authentication flow
	authResp, err := client.Auth.Login(ctx, "testuser", "testpass", "")
	if err != nil {
		log.Printf("Failed to login: %v", err)
	} else {
		fmt.Printf("Login successful! Token expires at: %s\n", authResp.ExpiresAt)

		// Verify token
		verifyResp, err := client.Auth.Verify(ctx, authResp.Token)
		if err != nil {
			log.Printf("Failed to verify token: %v", err)
		} else if verifyResp.Valid {
			fmt.Printf("Token is valid for user: %s\n", verifyResp.Username)
		}
	}

	fmt.Println("Aether Vault Go SDK demo completed!")
}
