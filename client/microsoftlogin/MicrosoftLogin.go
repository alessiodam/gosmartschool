package microsoftlogin

import (
	"fmt"
	"github.com/tebeka/selenium"
	"gosmartschool/assets"
	"log"
	"runtime"
	"strings"
	"time"
)

const (
	seleniumPort = 65534
)

type TwoFactorSecurityQuestions struct {
	BirthdayAnswer string `json:"birthdayAnswer"`
}

func MicrosoftLogin(domain string, microsoftEmail string, microsoftPassword string, twoFactorSecurityQuestions TwoFactorSecurityQuestions) (string, error) {
	var chromeDriverPath string
	var err error

	log.Println("Extracting temporary ChromeDriver binary...")
	switch runtime.GOOS {
	case "windows":
		chromeDriverPath, err = assets.ExtractFile("chromedriver.exe")
		if err != nil {
			return "", err
		}
	case "linux":
		chromeDriverPath, err = assets.ExtractFile("chromedriver")
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	log.Printf("ChromeDriver path: %s", chromeDriverPath)

	var opts []selenium.ServiceOption
	service, err := selenium.NewChromeDriverService(chromeDriverPath, seleniumPort, opts...)
	if err != nil {
		return "", fmt.Errorf("error starting the ChromeDriver server: %v", err)
	}
	defer func() {
		if err := service.Stop(); err != nil {
			log.Printf("Error stopping ChromeDriver service: %v", err)
		}
	}()

	caps := selenium.Capabilities{
		"browserName": "chrome",
		"chromeOptions": map[string]interface{}{
			"args": []string{
				"--disable-search-engine-choice-screen",
				"--no-sandbox",
			},
		},
	}

	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", seleniumPort))
	if err != nil {
		return "", fmt.Errorf("error connecting to WebDriver: %v", err)
	}
	defer func() {
		if err := wd.Quit(); err != nil {
			log.Printf("Error quitting WebDriver: %v", err)
		}
	}()

	if err := wd.Get("https://" + domain + "/login/sso/init/office365"); err != nil {
		return "", err
	}

	time.Sleep(2 * time.Second)

	emailField, err := waitForElement(wd, selenium.ByID, "i0116", 10*time.Second)
	if err != nil {
		return "", fmt.Errorf("error finding the email field: %v", err)
	}
	if err := emailField.SendKeys(microsoftEmail); err != nil {
		return "", err
	}
	log.Println("Entered email")

	time.Sleep(2 * time.Second)

	nextButton, err := waitForElement(wd, selenium.ByID, "idSIButton9", 10*time.Second)
	if err != nil {
		return "", fmt.Errorf("error finding the next button: %v", err)
	}
	if err := nextButton.Click(); err != nil {
		return "", err
	}
	log.Println("Clicked next button")

	time.Sleep(2 * time.Second)

	passwordField, err := waitForElement(wd, selenium.ByID, "i0118", 10*time.Second)
	if err != nil {
		return "", fmt.Errorf("error finding the password field: %v", err)
	}
	if err := passwordField.SendKeys(microsoftPassword); err != nil {
		return "", err
	}
	log.Println("Entered password")

	time.Sleep(2 * time.Second)

	signInButton, err := waitForElement(wd, selenium.ByID, "idSIButton9", 10*time.Second)
	if err != nil {
		return "", fmt.Errorf("error finding the sign-in button: %v", err)
	}
	if err := signInButton.Click(); err != nil {
		return "", err
	}
	log.Println("Clicked sign-in button")

	staySignedInButton, err := waitForElement(wd, selenium.ByID, "idSIButton9", 10*time.Second)
	if err == nil {
		if err := staySignedInButton.Click(); err != nil {
			return "", err
		}
		log.Println("Clicked 'Stay signed in' button")
	} else {
		log.Println("No 'Stay signed in' prompt found")
	}

	time.Sleep(5 * time.Second)

	pageSource, err := wd.PageSource()
	if err != nil {
		return "", fmt.Errorf("error getting page source: %v", err)
	}

	if strings.Contains(pageSource, "geboortedatum") {
		log.Println("Successfully logged in! But we have to fill in the date of birth.")
		geboorteDatumVerificationField, err := waitForElement(wd, selenium.ByID, "account_verification_form__security_question_answer", 10*time.Second)
		if err != nil {
			return "", fmt.Errorf("error finding the geboortedatum verification field: %v", err)
		}
		if err := geboorteDatumVerificationField.SendKeys(twoFactorSecurityQuestions.BirthdayAnswer); err != nil {
			return "", err
		}
		log.Println("Entered geboortedatum")
		if err := geboorteDatumVerificationField.SendKeys(selenium.EnterKey); err != nil {
			return "", err
		}
		log.Println("Clicked enter key")
	} else {
		log.Println("Doesn't seem to ask for date of birth. Let's see if we're logged in")
	}

	time.Sleep(3 * time.Second)

	pageSource, err = wd.PageSource()
	if err != nil {
		return "", fmt.Errorf("error getting page source: %v", err)
	}
	if !strings.Contains(pageSource, "Start") {
		return "", fmt.Errorf("login failed")
	}

	cookies, err := wd.GetCookies()
	if err != nil {
		return "", fmt.Errorf("error getting cookies: %v", err)
	}
	var phpSessId string
	for _, cookie := range cookies {
		if cookie.Name == "PHPSESSID" {
			phpSessId = cookie.Value
		}
	}
	if phpSessId == "" {
		return "", fmt.Errorf("could not find PHPSESSID cookie")
	}

	return phpSessId, nil
}

func waitForElement(wd selenium.WebDriver, by, value string, timeout time.Duration) (selenium.WebElement, error) {
	var elem selenium.WebElement
	var err error
	for start := time.Now(); time.Since(start) < timeout; time.Sleep(500 * time.Millisecond) {
		elem, err = wd.FindElement(by, value)
		if err == nil {
			return elem, nil
		}
	}
	return nil, fmt.Errorf("element not found: %v", err)
}
