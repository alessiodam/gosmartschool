package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gosmartschool/client"
	"gosmartschool/client/microsoftlogin"
	"log"
	"os"
	"time"
)

func main() {
	_ = godotenv.Load()
	var smartschoolDomain string
	smartschoolDomain = os.Getenv("DOMAIN")
	if smartschoolDomain == "" {
		fmt.Print("Smartschool domain (school.smartschool.be): ")
		_, err := fmt.Scanln(&smartschoolDomain)
		if err != nil {
			log.Fatalf("Failed to read input: %v", err)
		}
	}
	smartSchoolClient := client.NewSmartSchoolClient(smartschoolDomain)

	var twoFactorSecurityQuestions microsoftlogin.TwoFactorSecurityQuestions

	if os.Getenv("MICROSOFT_EMAIL") != "" && os.Getenv("MICROSOFT_PASSWORD") != "" {
		twoFactorSecurityQuestions = microsoftlogin.TwoFactorSecurityQuestions{
			BirthdayAnswer: os.Getenv("BIRTHDAY"),
		}
	} else {
		var microsoftEmail, microsoftPassword string

		fmt.Print("Microsoft email: ")
		_, err := fmt.Scanln(&microsoftEmail)
		if err != nil {
			log.Fatalf("Failed to read input: %v", err)
		}

		fmt.Print("Microsoft password: ")
		_, err = fmt.Scanln(&microsoftPassword)
		if err != nil {
			log.Fatalf("Failed to read input: %v", err)
		}

		fmt.Print("Birthday answer (MM/DD/YYYY): ")
		_, err = fmt.Scanln(&twoFactorSecurityQuestions.BirthdayAnswer)
		if err != nil {
			log.Fatalf("Failed to read input: %v", err)
		}

		loginSuccess, err := smartSchoolClient.MicrosoftLogin(smartschoolDomain, microsoftEmail, microsoftPassword, twoFactorSecurityQuestions)
		if err != nil || !loginSuccess {
			log.Printf("Failed to login: %v", err)
		} else {
			log.Println("Authenticated!")
		}
	}
	time.Sleep(5 * time.Second)
	fmt.Println("Done! PHPSESSID: ", smartSchoolClient.PhpSessId)
}
