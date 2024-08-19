package dhconn

import (
	"bufio"
	"fmt"
	"godehashed/parsedh"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const DURI string = "https://api.dehashed.com/search?query="

func makeRequest(url, username, password string) ([]byte, error) {
	client := http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request failed: %w", err)
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response failed: %w", err)
	}
	return body, nil
}

func DHConn(apikey, email, name, searchterm, uname, outfile, elist string, phone int) {
	username := strings.Split(strings.Trim(apikey, "\n"), ":")[0]
	password := strings.Split(strings.Trim(apikey, "\n"), ":")[1]
	pages := 1

	switch searchterm {
	case "email":
		fmt.Println("[*] Searching for emails.")
		url := fmt.Sprintf("%s%s:%s&size=10000&page=%d", DURI, searchterm, email, pages)
		body, err := makeRequest(url, username, password)
		if err != nil {
			log.Fatalf("Dehashed connection error: %v", err)
		}

		total, balance := parsedh.ParseDH(body, outfile)
		if balance == 0 {
			fmt.Println("[!] You're out of credits, please reload then run again.")
			os.Exit(0)
		}

		fmt.Println("[*] Total number of results:", total)

		if total >= 10000 {
			fmt.Println("[!] ATTENTION!! Total number of results greater than 10,000. Do you want to continue? ENTER: 'Y' or 'N' to exit.")
			var response string
			if _, err := fmt.Scanln(&response); err != nil {
				log.Fatalf("Input error: %v", err)
			}

			if response == "Y" || response == "y" {
				for pages = 2; pages <= 3; pages++ {
					fmt.Println("[!] Current Page is:", pages)
					time.Sleep(5 * time.Second)
					fmt.Println("[*] We are searching for emails.")
					url := fmt.Sprintf("%s%s:%s&size=10000&page=%d", DURI, searchterm, email, pages)
					body, err := makeRequest(url, username, password)
					if err != nil {
						log.Fatalf("Dehashed connection error: %v", err)
					}
					parsedh.ParseDH(body, outfile)
				}
			} else {
				fmt.Println("[*] Ending collection gathering")
			}
		}

	case "name":
		searchname := strings.ReplaceAll(name, " ", "+")
		fmt.Println("[*] We are searching for names.")
		url := fmt.Sprintf("%sname:%s", DURI, searchname)
		body, err := makeRequest(url, username, password)
		if err != nil {
			log.Fatalf("Dehashed connection error: %v", err)
		}

		total, balance := parsedh.ParseDH(body, outfile)
		if balance == 0 {
			fmt.Println("[!] You're out of credits, please reload then run again.")
			os.Exit(0)
		}

		fmt.Println("[*] Total number of results:", total)

	case "phone":
		fmt.Println("[*] We are searching for phone numbers.")
		url := fmt.Sprintf("%sphone:%d", DURI, phone)
		body, err := makeRequest(url, username, password)
		if err != nil {
			log.Fatalf("Dehashed connection error: %v", err)
		}

		total, balance := parsedh.ParseDH(body, outfile)
		if balance == 0 {
			fmt.Println("[!] You're out of credits, please reload then run again.")
			os.Exit(0)
		}

		fmt.Println("[*] Total number of results:", total)

	case "username":
		fmt.Println("[*] We are searching for usernames.")
		url := fmt.Sprintf("%susername:%s", DURI, uname)
		body, err := makeRequest(url, username, password)
		if err != nil {
			log.Fatalf("Dehashed connection error: %v", err)
		}

		total, balance := parsedh.ParseDH(body, outfile)
		if balance == 0 {
			fmt.Println("[!] You're out of credits, please reload then run again.")
			os.Exit(0)
		}

		fmt.Println("[*] Total number of results:", total)

	case "list":
		fmt.Println("[*] Going into List mode, will add a 2-second time delay to prevent blacklist.")
		fmt.Println("[*] ATTENTION: This can take a long time depending on the size of the list and will use A LOT OF CREDITS!")
		time.Sleep(2 * time.Second)

		file, err := os.Open(elist)
		if err != nil {
			log.Fatalf("Cannot read file: %v", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var dhlist []string

		for scanner.Scan() {
			dhlist = append(dhlist, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("Error reading file: %v", err)
		}

		for _, line := range dhlist {
			fmt.Println("[*] We are searching for emails.")
			url := fmt.Sprintf("%s%s", DURI, line)
			body, err := makeRequest(url, username, password)
			if err != nil {
				log.Fatalf("Dehashed connection error: %v", err)
			}

			total, balance := parsedh.ParseDH(body, outfile)
			if balance == 0 {
				fmt.Println("[!] You're out of credits, please reload then run again.")
				os.Exit(0)
			}

			fmt.Println("[*] Total number of results:", total)
			fmt.Println("[*] Delaying... Please wait.")
			time.Sleep(2 * time.Second)
		}

	default:
		fmt.Println("Please enter a valid search term: 'email', 'name', 'phone', 'username', or 'list'.")
		os.Exit(0)
	}
}
