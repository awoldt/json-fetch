package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	Scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("\nEnter API endpoint to fetch json data from: ")
	var endpoint string
	fmt.Scanln(&endpoint)

	// make sure http or https endpoint
	if !strings.Contains(endpoint, "http") && !strings.Contains(endpoint, "https") {
		fmt.Printf("ERROR: must enter http or https endpoint")
	} else {
		req, err := http.Get(endpoint)
		if err != nil {
			fmt.Println("ERROR: could not fetch endpoint " + endpoint)
		}
		defer req.Body.Close()

		// content type header must be application/json
		if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
			fmt.Println("ERROR: Content-Type header must be 'application/json'. Got '" + req.Header.Get("Content-Type") + "' instead")
		} else {
			data, err2 := io.ReadAll(req.Body)
			if err2 != nil {
				fmt.Println("ERROR: could not read request body")
			}
			var filename string
			fmt.Print("\nWhat would you like to name this file: ")
			Scanner.Scan()
			filename = Scanner.Text()
			// while user enters invalid filename
			for {
				if filename == "" {
					fmt.Print("Must enter valid filename : ")
					Scanner.Scan()
					filename = Scanner.Text()

				} else {
					break
				}
			}

			// if filename has spaces, join with '-'
			if len(strings.Split(filename, " ")) > 1 {
				filename = strings.Join(strings.Split(filename, " "), "-")
			}

			_, err := os.OpenFile("data/"+filename+".json", os.O_RDONLY, 0644)
			// check if file exists
			// if does exist, ask for permission to overwrite
			if errors.Is(err, os.ErrNotExist) {
				wErr := os.WriteFile("data/"+filename+".json", []byte(data), 0644)
				if wErr != nil {
					fmt.Println("ERROR: could not write json data to file")
				}
				fmt.Println("\nDONE")

			} else {
				fmt.Println("\n" + filename + ".json already exists. Would you like to overwrite this file?\ny(n): ")
				var overwite string
				fmt.Scanln(&overwite)
				if overwite == "y" || overwite == "Y" {
					wErr := os.WriteFile("data/"+filename+".json", []byte(data), 0644)
					if wErr != nil {
						fmt.Println("ERROR: could not write json data to file")
					}
					fmt.Println("\nDONE")
				} else {
					fmt.Println("\nFile not saved.")
				}

			}

		}
	}
}
