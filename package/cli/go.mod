module github.com/skygenesisenterprise/aether-vault/package/cli

go 1.25.5

require (
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/spf13/cobra v1.8.0
	golang.org/x/crypto v0.46.0
	golang.org/x/term v0.38.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	golang.org/x/sys v0.40.0 // indirect
)

replace github.com/skygenesisenterprise/aether-vault/package/cli/server => ./server
