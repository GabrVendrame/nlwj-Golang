package main

// go:generate tern migrate --migrations ./internal/pgstore/migrations --config ./internal/pgstore/migrations/tern.conf go generate ./... | run go generate

// go:generate sqlc generate -f ./internal/pgstore/sqlc.yaml

