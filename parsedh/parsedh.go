package parsedh

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
)

type DHAPI struct {
	Balance int `json:"balance"`
	Entries []struct {
		ID             string `json:"id"`
		Email          string `json:"email"`
		IPAddress      string `json:"ip_address"`
		Username       string `json:"username"`
		Password       string `json:"password"`
		HashedPassword string `json:"hashed_password"`
		Name           string `json:"name"`
		Vin            string `json:"vin"`
		Address        string `json:"address"`
		Phone          string `json:"phone"`
		DatabaseName   string `json:"database_name"`
	} `json:"entries"`
	Success bool   `json:"success"`
	Took    string `json:"took"`
	Total   int    `json:"total"`
}

func ParseDH(body []byte, outfile string) (int, int) {
	var jsonAPI DHAPI

	err := json.Unmarshal(body, &jsonAPI)
	if err != nil {
		log.Fatalf("There was an error unmarshaling body: %v", err)
	}

	total := jsonAPI.Total
	balance := jsonAPI.Balance
	dhdata := jsonAPI.Entries

	var user []string
	for _, value := range dhdata {
		user = append(user, fmt.Sprintf("Database: %s Username: %s Email: %s Password: %s Hash: %s Phone: %s Name: %s Address: %s\n",
			value.DatabaseName, value.Username, value.Email, value.Password, value.HashedPassword, value.Phone, value.Name, value.Address))
	}

	// If outfile is not empty, will export data to outfile.
	if outfile != "" {
		csvFile, err := os.OpenFile(outfile, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}
		defer csvFile.Close()

		writer := csv.NewWriter(csvFile)
		defer writer.Flush()

		for _, entry := range jsonAPI.Entries {
			row := []string{
				entry.DatabaseName,
				entry.Username,
				entry.Email,
				entry.Password,
				entry.HashedPassword,
				entry.Name,
				entry.Address,
				entry.Phone,
			}
			if err := writer.Write(row); err != nil {
				log.Fatalf("Error writing to CSV: %v", err)
			}
		}
	}

	// By default, output will be displayed to console
	for _, u := range user {
		if u != "" {
			fmt.Println(u)
		}
	}

	fmt.Println("[*] Your API balance remaining:", strconv.Itoa(balance))
	return total, balance
}

func SetHeader(outfile string) {
	if outfile != "" {
		csvFile, err := os.OpenFile(outfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}
		defer csvFile.Close()

		writer := csv.NewWriter(csvFile)
		defer writer.Flush()

		header := []string{"Database Name", "Username", "Email", "Password", "Hashed Password", "Name", "Address", "Phone"}
		if err := writer.Write(header); err != nil {
			log.Fatalf("Error writing header to CSV: %v", err)
		}
	}
}
