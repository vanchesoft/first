package main

import (
	b64 "encoding/base64"
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
	maxGetRounded := 5
	// categoryOrigin := []string{"truck"}
	// categoryOrigin := []string{"heavy"}
	// categoryOrigin := []string{"bus"}
	// categoryOrigin := []string{"trailer"}
	// categoryOrigin := []string{"motorhome"}
	// categoryOrigin := []string{"random"}
	// categoryOrigin := []string{"truck", "heavy", "bus", "trailer", "motorhome", "random"}
	categoryOrigin := []string{"truck", "heavy", "bus", "trailer", "motorhome"}

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
		sEnc := b64.StdEncoding.EncodeToString(content)
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

				min := 5
				max := 40
				randomNumber := rand.Intn(max-min) + min

				dataLink, _ := b64.StdEncoding.DecodeString("aHR0cHM6Ly92YW5jaGUtYXBpLm5ldGxpZnkuYXBwL2FwaS9nZXQtdmVoaWNsZXM/JmNhdGVnb3J5PQ==")
				price := "&min_price=500.000&max_price=900.000"
				if category == "trailer" || category == "motorhome" {
					price = ""
				}
				resource := string(dataLink) + category + "&page=" + strconv.Itoa(randomNumber) + "&limit=1&sortby=price_desc" + price

				//Getting a client to make the https://api.mercadolibre.com/items/MLU439286635
				// var response *http.Response
				// ver err error
				response, err := http.Get(resource)
				check(err)

				body, err := ioutil.ReadAll(response.Body)
				check(err)
				// fmt.Fprintf(w, "%s", body)
				defer response.Body.Close()

				var result Response
				if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
					fmt.Println("Can not unmarshal JSON")
				}

				if result.Data != nil && len(result.Data) > 0 {
					// fmt.Println(result.Data)
					textWrite.Data = append(textWrite.Data, result.Data[0])

					count = count + 1
					if count >= maxGetRegister {
						break
					}
				}
			}

			categoryNumber = GetCategoryNumber(j)
			newFsConfigBytes, _ := json.Marshal(textWrite.Data)

			if withRaw {
				err := os.WriteFile("./d-h/"+categoryNumber+"/raw-"+strconv.Itoa(rounded)+".json", newFsConfigBytes, 0644)
				check(err)
			}

			salt := StringWithCharset(15, charset)

			sEnc := b64.StdEncoding.EncodeToString(newFsConfigBytes)
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
