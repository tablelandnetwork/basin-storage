package main

import (
	"log"
	"os"

	// Blank-import the function package so the init() runs
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	_ "github.com/tablelandnetwork/basin-storage"
	"gopkg.in/yaml.v2"
)

type UploaderVars struct {
	W3SToken string `yaml:"WEB3STORAGE_TOKEN"`
	CrdbConn string `yaml:"CRDB_CONN_STRING"`
}

type StatusCheckerVars struct {
	W3SToken   string `yaml:"WEB3STORAGE_TOKEN"`
	CrdbConn   string `yaml:"CRDB_CONN_STRING"`
	PrivateKey string `yaml:"PRIVATE_KEY"`
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
		vars := UploaderVars{}
		err = yaml.Unmarshal(data, &vars)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		os.Setenv("WEB3STORAGE_TOKEN", vars.W3SToken)
		os.Setenv("CRDB_CONN_STRING", vars.CrdbConn)
	}

	if targetFn == "StatusChecker" {
		data, err := os.ReadFile("checker.env.yml")
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		vars := StatusCheckerVars{}
		err = yaml.Unmarshal(data, &vars)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		os.Setenv("WEB3STORAGE_TOKEN", vars.W3SToken)
		os.Setenv("CRDB_CONN_STRING", vars.CrdbConn)
		os.Setenv("PRIVATE_KEY", vars.PrivateKey)
	}

	// Unmarshal the YAML data into the struct

	// read config from env files and set as env vars
	// if target functions is uploader load uploader env
	// if target functions is statuschecker load statuschecker env

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
