package vault

// Package vault provides a complete Go SDK for Aether Vault, a sovereign, open-source, self-hostable secret management system.
//
// This SDK abstracts all HTTP calls to /api/v1/* and provides an idiomatic, typed, and secure Go API.
//
// Basic usage:
//
//	cfg := vault.NewConfig("http://localhost:8080", "your-token")
//	client, err := vault.New(cfg)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Get a secret
//	secret, err := client.Secrets.Get(context.Background(), "db/password")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Printf("Secret: %s\n", secret.Value)
//
// The SDK supports all major Vault features:
// - Authentication & token management
// - Secret management with versioning
// - TOTP (Time-based One-Time Password)
// - Identity & access management
// - Policy enforcement
// - Audit logging
// - Transport security (TLS/mTLS)
// - Retry mechanisms and middleware
//
// For more detailed examples, see the examples/ directory.
package vault