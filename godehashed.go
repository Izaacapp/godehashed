package main

import (
	"flag"
	"fmt"
	"godehashed/dhconn"
	"godehashed/parsedh"
	"os" // Replace ioutil with os
	"log"
)

func Importkey(data string) string {
	// Use os.ReadFile instead of ioutil.ReadFile
	keydata, err := os.ReadFile(data)
	if err != nil {
		log.Fatal("Error reading key: ", err)
	}

	apikey := string(keydata)
	return apikey
}

func main() {

	keyname := ""
	email := ""
	outfile := ""
	name := ""
	phone := 0
	searchterm := ""
	uname := ""
	elist := ""

	flag.StringVar(&keyname, "i", "", "Name of apikey to import.")
	flag.StringVar(&name, "n", "", "Name we are searching for.")
	flag.StringVar(&email, "e", "", "Email we are searching for")
	flag.StringVar(&uname, "u", "", "Username we are searching for")
	flag.IntVar(&phone, "p", 0, "Phone number we are searching for")
	flag.StringVar(&searchterm, "s", "", "Specify what we are searching for: name, email or username. Then add corresponding switch.")
	flag.StringVar(&outfile, "o", "", "Outfile file name, will output in CSV Format.")
	flag.StringVar(&elist, "l", "", "Search a list of emails.")
	flag.Parse()

	if keyname == "" {
		fmt.Println("You must include an API key to use '-i'")
	} else {
		apikey := Importkey(keyname)
		parsedh.SetHeader(outfile)
		dhconn.DHConn(apikey, email, name, searchterm, uname, outfile, elist, phone)
	}

}
