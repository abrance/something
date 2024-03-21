package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	env := "./.env"
	loadEnv(env)
}

func loadEnv(envPath string) {
	// Load the .env file
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
	fmt.Println("SN: ", os.Getenv("SN"))
}
