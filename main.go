package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	log "github.com/rs/zerolog/log"
)

type Response struct {
	Data   []interface{} `json:"data"`
	Paging interface{}   `json:"paging"`
}

type DataFormated struct {
	Data []interface{} `json:"data"`
}

const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func main() {

	var categoryNumber string
	var category string

	withRaw := false
	encodedOnly := false
	maxGetRegister := 8
	maxGetRounded := 20
	// categoryOrigin := []string{"truck"}
	// categoryOrigin := []string{"heavy"}
	// categoryOrigin := []string{"bus"}
	// categoryOrigin := []string{"trailer"}
	// categoryOrigin := []string{"motorhome"}
	// categoryOrigin := []string{"random"}
	categoryOrigin := []string{"truck", "heavy", "bus", "trailer", "motorhome"}
	// categoryOrigin := []string{"truck"}

	if encodedOnly {

		var payload []interface{}
		registerNumber := "1"
		categoryNumber = GetCategoryNumber(categoryOrigin[0])

		content, err := ioutil.ReadFile("./d-h/" + categoryNumber + "/raw-" + registerNumber)

		err = json.Unmarshal(content, &payload)
		if err != nil {
			log.Error().Any("error", err).Msg("Error during Unmarshal()")
		}
		check(err)

		salt := StringWithCharset(15, charset)
		sEnc := base64.StdEncoding.EncodeToString(content)
		salt = append(salt, sEnc...)
		err = os.WriteFile("./d-h/"+categoryNumber+"/"+registerNumber, salt, 0644)
		check(err)

		os.Exit(1)
	}
	for _, j := range categoryOrigin {

		fmt.Println("Start Category", j)
		rounded := 0
		for {
			var textWrite DataFormated
			// not
			count := 0
			var numbersJustResolv []int
			for {

				category = j
				if j == "random" {

					minR := 0
					maxR := 6
					randomNumberR := rand.Intn(maxR-minR) + minR
					if randomNumberR == 0 {
						category = "truck"
					}
					if randomNumberR == 1 {
						category = "heavy"
					}
					if randomNumberR == 2 {
						category = "bus"
					}
					if randomNumberR == 3 {
						category = "trailer"
					}
					if randomNumberR == 4 {
						category = "motorhome"
					}
					if randomNumberR == 5 {
						category = "truck"
					}
				}
				minP := 1
				maxP := 5
				randomNumberP := rand.Intn(maxP-minP) + minP

				dataLink, _ := base64.StdEncoding.DecodeString("aHR0cDovL2xvY2FsaG9zdDozMDAwL2FwaS9nZXQtYWRzP3NlYXJjaD0mY2F0ZWdvcnk9")
				// dataLink, _ := base64.StdEncoding.DecodeString("aHR0cHM6Ly92YW5jaGUuY29tLmJyL2FwaS9nZXQtYWRzP3NlYXJjaD0mY2F0ZWdvcnk9")
				price := ""
				sort := "&sortby=price_desc"
				if category != "truck" {
					sort = ""
					price = "&min_price=10.000&max_price=*"
				}
				if category == "trailer" {
					sort = "&sortby=price_desc"
					price = "&min_price=15.000&max_price=*"
				}
				if category == "heavy" {
					sort = "&sortby=price_desc"
					price = "&min_price=10.000&max_price=*"
				}
				resource := string(dataLink) + category + "&page=" + strconv.Itoa(randomNumberP) + "&limit=50" + sort + price

				//Getting a client to make the https://api.mercadolibre.com/items/MLU439286635
				// var response *http.Response
				// ver err error
				response, err := http.Get(resource)
				check(err)

				body, err := ioutil.ReadAll(response.Body)
				check(err)

				// fmt.Println("resource body", string(body))
				fmt.Println("AAAAAAAAAA", resource)

				// fmt.Fprintf(w, "%s", body)
				defer response.Body.Close()

				var result Response
				if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
					fmt.Println("Can not unmarshal JSON", err)
				}

				if result.Data != nil && len(result.Data) > 0 {
					// fmt.Println(result.Data)
					min := 5
					max := 50
					var randomNumber int
					for {
						randomNumber = rand.Intn(max-min) + min
						if !contains(numbersJustResolv, randomNumber) {
							fmt.Println("Find number", randomNumber)

							numbersJustResolv = append(numbersJustResolv, randomNumber)

							if result.Data != nil && len(result.Data) > randomNumber {
								textWrite.Data = append(textWrite.Data, result.Data[randomNumber])

								count = count + 1
								if count >= maxGetRegister {
									break
								}
							}
						}
					}
				}

				if count >= maxGetRegister {
					break
				}
			}

			categoryNumber = GetCategoryNumber(j)
			newFsConfigBytes, _ := json.Marshal(textWrite.Data)

			if withRaw {
				err := os.WriteFile("./d-h/"+categoryNumber+"/raw-"+strconv.Itoa(rounded)+".json", newFsConfigBytes, 0644)
				check(err)
			}

			salt := StringWithCharset(15, charset)

			sEnc := base64.RawStdEncoding.EncodeToString(newFsConfigBytes)
			salt = append(salt, sEnc...)
			err := os.WriteFile("./d-h/"+categoryNumber+"/"+strconv.Itoa(rounded), salt, 0644)
			check(err)

			rounded = rounded + 1
			if rounded >= maxGetRounded {
				break
			}
		}

		fmt.Println("---- End Category", j)
	}

	fmt.Println("end")
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func check(e error) {
	if e != nil {
		log.Error().Msg("Error: " + e.Error())
		log.Debug().Msg("Returning Authentication URL: ")
		panic(e)
	}
}

func GetCategoryNumber(category string) string {
	if category == "truck" {
		return "0"
	}
	if category == "heavy" {
		return "1"
	}
	if category == "bus" {
		return "2"
	}
	if category == "trailer" {
		return "3"
	}
	if category == "motorhome" {
		return "4"
	}
	return "5"
}

func StringWithCharset(length int, charset string) []byte {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return b
}
