package main

import (
	"bytes"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"
	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
)

func openSerialPort(portName string, baudRate int) (serial.Port, error) {
	mode := &serial.Mode{
		BaudRate: baudRate,
	}
	port, err := serial.Open(portName, mode)
	if err != nil {
		return nil, fmt.Errorf("failed to open port: %v", err)
	}
	return port, nil
}

func readOneCharFromSerial(port serial.Port) (byte, error) {
	// Create a buffer to hold one character
	buffer := make([]byte, 1)

	for {
		// Read one character from the serial port
		n, err := port.Read(buffer)
		if err != nil {
			return 0, fmt.Errorf("failed to read from port: %v", err)
		}

		if n == 0 {
			return 0, fmt.Errorf("no data read from port")
		}

		// Check if the character is '0' or '1'
		if buffer[0] == '0' || buffer[0] == '1' {
			return buffer[0], nil
		}
	}
}

func PlaceHigh(apiKey, betValue, currency string, mode bool) bool {
	var result bool
	url := "https://duckdice.io/api/play?api_key=" + apiKey

	// Create a bet
	sampleBet := BetPayload{
		Symbol: currency,
		Chance: "50",
		IsHigh: true,
		Amount: betValue,
		Faucet: mode,
	}

	// Create a CookieJar to store cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println("Error creating cookie jar:", err)
		return false
	}

	// Marshal the payload into JSON
	jsonPayload, err := json.Marshal(sampleBet)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return false
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return false
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client that uses the CookieJar
	client := &http.Client{
		Jar: jar,
	}
	client.Transport = cloudflarebp.AddCloudFlareByPass(client.Transport)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return false
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return false
	}

	// Retrying when captcha triggered
	for string(body)[2:9] == "DOCTYPE" {
		waiter := rand.Uint32N(27) + 3
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return false
		}

		// Set the Content-Type header
		req.Header.Set("Content-Type", "application/json")

		client.Transport = cloudflarebp.AddCloudFlareByPass(client.Transport)
		fmt.Println("CAPTCHA!😠😠😠")
		fmt.Printf("Waiting %d seconds", waiter)

		// Implented 10 second wait
		for range waiter {
			time.Sleep(time.Second)
			fmt.Print(".")
		}
		fmt.Println()
		fmt.Println("Retrying...")

		//Send request again
		captchaResp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return false
		}
		defer captchaResp.Body.Close()

		// Read the new response
		body, err = io.ReadAll(captchaResp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return false
		}
	}

	var betResp BetResponse
	err = json.Unmarshal([]byte(body), &betResp)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return false
	}

	result = betResp.Bet.Result

	fmt.Println("Roll is", betResp.Bet.Number)

	return result
}

func PlaceLow(apiKey, betValue, currency string, mode bool) bool {
	var result bool
	url := "https://duckdice.io/api/play?api_key=" + apiKey

	// Create a bet
	sampleBet := BetPayload{
		Symbol: currency,
		Chance: "50",
		IsHigh: false,
		Amount: betValue,
		Faucet: mode,
	}

	// Create a CookieJar to store cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println("Error creating cookie jar:", err)
		return false
	}

	// Marshal the payload into JSON
	jsonPayload, err := json.Marshal(sampleBet)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return false
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return false
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client that uses the CookieJar
	client := &http.Client{
		Jar: jar,
	}
	client.Transport = cloudflarebp.AddCloudFlareByPass(client.Transport)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return false
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return false
	}

	// Retrying when captcha triggered
	for string(body)[2:9] == "DOCTYPE" {
		waiter := rand.Uint32N(27) + 3
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return false
		}

		// Set the Content-Type header
		req.Header.Set("Content-Type", "application/json")

		client.Transport = cloudflarebp.AddCloudFlareByPass(client.Transport)
		fmt.Println("CAPTCHA!😠😠😠")
		fmt.Printf("Waiting %d seconds", waiter)

		// Implented 10 second wait
		for range waiter {
			time.Sleep(time.Second)
			fmt.Print(".")
		}
		fmt.Println()
		fmt.Println("Retrying...")

		//Send request again
		captchaResp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return false
		}
		defer captchaResp.Body.Close()

		// Read the new response
		body, err = io.ReadAll(captchaResp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return false
		}
	}

	var betResp BetResponse
	err = json.Unmarshal([]byte(body), &betResp)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return false
	}

	result = betResp.Bet.Result

	fmt.Println("Roll is", betResp.Bet.Number)

	return result
}


func main() {
	apiFile, err := os.Open("API.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(apiFile)
	scanner.Scan()
	apiKey := scanner.Text() // Replace with your actual API key
	url := "https://duckdice.io/api/bot/user-info?api_key=" + apiKey
	err = apiFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Create a CookieJar to store cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println("Error creating cookie jar:", err)
		return
	}

	// Create an HTTP client that uses the CookieJar
	client := &http.Client{
		Jar: jar,
	}
	client.Transport = cloudflarebp.AddCloudFlareByPass(client.Transport)

	// Create a new HTTP GET request
	response, err := client.Get(url)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer response.Body.Close() // Ensure the response body is closed later

	// Read the response data
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Retrying when captcha triggered
	for string(body)[2:9] == "DOCTYPE" {
		waiter := rand.Uint32N(27) + 3

		client.Transport = cloudflarebp.AddCloudFlareByPass(client.Transport)
		fmt.Println("CAPTCHA!😠😠😠")
		fmt.Printf("Waiting %d seconds", waiter)

		// Implemented 10 second wait
		for range waiter {
			time.Sleep(time.Second)
			fmt.Print(".")
		}
		fmt.Println()
		fmt.Println("Retrying...")

		// Send request again
		captchaResp, err := client.Get(url)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer captchaResp.Body.Close()

		// Read the new response
		body, err = io.ReadAll(captchaResp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}
	}

	var userInfo UserInfoResponse
	err = json.Unmarshal([]byte(body), &userInfo)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	fmt.Println("Username:", userInfo.Username)
	fmt.Println("Balances:")
	for _, balans := range userInfo.Balances {
		if balans.Main == "" {
			fmt.Println(balans.Faucet, " ", balans.Currency, "(Faucet)")
		} else {
			fmt.Println(balans.Main, " ", balans.Currency)
		}
	}

	var bet float64
	var betOrig float64
	var currency string
	fMode := true
	//fmt.Print("Insert bet value: ")
	//fmt.Scan(&bet)
	bet = 0.006
	amount := fmt.Sprintf("%.2f", bet) 
	betOrig = bet
	//fmt.Print("Choose currency: ")
	//fmt.Scan(&currency)
	currency = "doge"
	// List available serial ports (optional)
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatal(err)
	}
	for _, port := range ports {
		fmt.Printf("Found port: %s\n", port.Name)
	}
	var selection string
	//fmt.Print("Serial port(from available): ")
	//fmt.Scan(&selection)
	selection = "/dev/ttyACM0"
	// Open the serial port once
	port, err := openSerialPort(selection, 115200)
	if err != nil {
		log.Fatalf("Error opening serial port: %v", err)
	}
	defer port.Close()

	for i := 0; i < 99999999; i++ {
		char, err := readOneCharFromSerial(port)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		fmt.Printf("Read character: %c\n", char)

		if char == '0' {
			rez := PlaceLow(apiKey, amount, currency, fMode)
			if rez {
				fmt.Println("Bet successful.✅")
				amount = fmt.Sprintf("%.2f", betOrig) // Convert result to string
				bet = betOrig
			} else {
				fmt.Println("Bet unsuccessful.☯")
				bet = bet * 2
				amount = fmt.Sprintf("%.2f", bet) // Convert result to string
			}
		} else if char == '1' {
			rez := PlaceHigh(apiKey, amount, currency, fMode)
			if rez {
				fmt.Println("Bet successful.✅")
				amount = fmt.Sprintf("%.2f", betOrig) // Convert result to string
				bet = betOrig

			
			} else {
				fmt.Println("Bet unsuccessful.☯")
				bet = bet * 2
				amount = fmt.Sprintf("%.2f", bet) // Convert result to string

			}
		}
	}
}
