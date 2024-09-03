package main

import (
	"fmt"
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

	messages, err := smartSchoolClient.ListMessages("inbox", 0)
	if err != nil {
		log.Fatalf("Failed to list messages: %v", err)
	}
	for _, message := range messages {
		fmt.Printf("ID: %s, From: %s, Subject: %s, Date: %s, Unread: %s\n",
			message.ID, message.From, message.Subject, message.Date, message.Unread)
	}
}
