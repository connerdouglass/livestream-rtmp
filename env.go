package main

import (
	"log"
	"os"
)

func RequireEnv(key string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	log.Fatalf("Missing environment variable: \"%s\"\n", key)
	return ""
}

func EnvOrDefault(key, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultValue
}
