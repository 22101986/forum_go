package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

var ( // data from local.env, global and "constants"
	AppEnv     string
	ServerPort string
	ServerHost string
	CertPath   string
	KeyPath    string
	DBHost     string
	DBPath     string
)

func loadEnvFile(filepath string) error {
	file, err := os.Open(filepath) // opening file
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file) // scan line by line ignoring comments and empty lines
	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		// splitting lines around "=" to map them in environment
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Printf("Ligne ignorée (pas au format clé=valeur) : %s", line)
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		err = os.Setenv(key, value)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	// giving values to the global variables
	AppEnv = os.Getenv("APP_ENV")
	ServerPort = os.Getenv("SERVER_PORT")
	ServerHost = os.Getenv("SERVER_HOST")
	CertPath = os.Getenv("CERT_PATH")
	KeyPath = os.Getenv("KEY_PATH")
	DBHost = os.Getenv("DB_HOST")
	DBPath = os.Getenv("DB_PATH")

	return nil
}
