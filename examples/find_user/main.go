package main

import (
	"github.com/joho/godotenv"
	"gosmartschool/client"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	smartSchoolClient := client.NewSmartSchoolClient(os.Getenv("DOMAIN"))
	smartSchoolClient.PhpSessId = os.Getenv("PHPSESSID")

	err = smartSchoolClient.CheckIfAuthenticated()
	if err != nil {
		log.Fatalf("Failed to authenticate: %v", err)
	}
	log.Println("Authenticated!")

	user, err := smartSchoolClient.FindUsersByName("Alessio")
	if err != nil {
		log.Fatalf("Failed to find user: %v", err)
	}
	log.Println(user)
}
