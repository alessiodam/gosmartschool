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
	smartSchoolClient.UniqueUsc = os.Getenv("UNIQUEUSC")
	smartSchoolClient.PhpSessId = os.Getenv("PHPSESSID")
	smartSchoolClient.Pid = os.Getenv("PID")

	err = smartSchoolClient.CheckIfAuthenticated()
	if err != nil {
		log.Fatalf("Failed to authenticate: %v", err)
	}
	log.Println("Authenticated!")

	courses, err := smartSchoolClient.GetCourses()
	if err != nil {
		log.Fatalf("Failed to get courses: %v", err)
	}

	for _, course := range courses {
		fmt.Printf("ID: %d, Platform ID: %d, Name: %s\n", course.ID, course.PlatformID, course.Name)
	}
}
