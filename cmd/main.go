package main

import (
	"log"
	"os"

	// Blank-import the function package so the init() runs.
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	_ "github.com/tablelandnetwork/basin-storage"
	"gopkg.in/yaml.v2"
)

type uploaderVars struct {
	W3SToken string `yaml:"WEB3STORAGE_TOKEN"`
	CrdbConn string `yaml:"CRDB_CONN_STRING"`
}

type statusCheckerVars struct {
	W3SToken   string `yaml:"WEB3STORAGE_TOKEN"`
	CrdbConn   string `yaml:"CRDB_CONN_STRING"`
	PrivateKey string `yaml:"PRIVATE_KEY"`
	ChainID    string `yaml:"CHAIN_ID"`
}

func main() {
	// Use PORT environment variable, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	targetFn := os.Getenv("FUNCTION_TARGET")

	if targetFn == "Uploader" {
		data, err := os.ReadFile("uploader.env.yml")
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		vars := uploaderVars{}
		if err = yaml.Unmarshal(data, &vars); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err = os.Setenv("WEB3STORAGE_TOKEN", vars.W3SToken); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err = os.Setenv("CRDB_CONN_STRING", vars.CrdbConn); err != nil {
			log.Fatalf("error: %v", err)
		}
	}

	if targetFn == "StatusChecker" {
		data, err := os.ReadFile("checker.env.yml")
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		vars := statusCheckerVars{}
		if err = yaml.Unmarshal(data, &vars); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err = os.Setenv("WEB3STORAGE_TOKEN", vars.W3SToken); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err = os.Setenv("CRDB_CONN_STRING", vars.CrdbConn); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err = os.Setenv("PRIVATE_KEY", vars.PrivateKey); err != nil {
			log.Fatalf("error: %v", err)
		}
		if err = os.Setenv("CHAIN_ID", vars.ChainID); err != nil {
			log.Fatalf("error: %v", err)
		}
	}

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
