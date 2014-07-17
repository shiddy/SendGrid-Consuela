package consuela

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//Used for removing only the emails from SendGrid's Web API
type JsonEmail struct {
	Email string `json:"email"`
}

//Used for making get calls from SendGrid's Web API with the ability to request what event you would like to reutrn
func ApiGet(eventName, username, password string) map[string]string {
	uri := "https://api.sendgrid.com/api/" + eventName + ".get.json?api_user=" + username + "&api_key=" + password
	resp, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}

	robots, _ := ioutil.ReadAll(resp.Body)

	var unsubscribes []JsonEmail
	json.Unmarshal(robots, &unsubscribes)

	returnedMap := make(map[string]string)

	for i := 0; i < len(unsubscribes); i++ {
		returnedMap[unsubscribes[i].Email] = ""
	}

	return returnedMap
}

//Merges an array of map[string]strings and returns a new map[string]string
func mapMerge(setArray []map[string]string) map[string]string {
	mergedList := make(map[string]string)
	for i := 0; i < len(setArray); i++ {
		for k, _ := range setArray[i] {
			mergedList[k] = ""
		}
	}

	return mergedList
}

// Opens a .csv when passed the file path and returns a map[string]string of the first column of the csv
func CsvToMap(file string) map[string]string {

	wl, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer wl.Close()

	reader := csv.NewReader(wl)

	reader.Comma = ';'
	reader.LazyQuotes = true

	lines, _ := reader.ReadAll()

	returnedMap := make(map[string]string)

	for i := 0; i < len(lines); i++ {
		returnedMap[lines[i][0]] = ""
	}

	return returnedMap
}

//Removes any instances of emails contained in unsubscribeList that are used in wholeList. Returns a list without duplicates.
func compareLists(wholeList, unsubscribeList map[string]string) map[string]string {
	for k, _ := range unsubscribeList {
		delete(wholeList, k)
	}

	return wholeList
}

//This takes the two map[string]string maps and creates a .zip file with the two maps turned into .csv files.
func zipOutput(username string, wholeList, donotsend map[string]string) {
	wholelistname := username + "_newlist.csv"
	donotsendname := username + "_DONOTSEND.csv"

	// Create a buffer to write our archive to.
	buffer := new(bytes.Buffer)

	// Create a new zip archive.
	writer := zip.NewWriter(buffer)

	// Add the file data to the archive.
	var files = []struct {
		Name string
		Body map[string]string
	}{
		{wholelistname, wholeList},
		{donotsendname, donotsend},
	}

	for _, file := range files {
		f, err := writer.Create(file.Name)
		if err != nil {
			log.Fatal(err)
		}
		for k, _ := range file.Body {
			_, err = f.Write([]byte(k + "\n"))
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// Checking for errors on close
	err := writer.Close()
	if err != nil {
		log.Fatal(err)
	}
	//write the zipped file to the disk
	ioutil.WriteFile(username+".zip", buffer.Bytes(), 0777)
}

// Pulls all of the user's data from the webAPI
//   If you intend to use this with another list of unwanted emails
//   you would run the CsvToMap fucntion with your wanted email list
//   and run a compareLists on the two maps.
func main() {
	username := "unsubscribe_check"
	password := "herpderp"

	fmt.Println("Opening your giant list")
	wholeList := CsvToMap("wholelist.csv")
	fmt.Println("Gathering your Unsubscribes, Bounces, Invalids, Blocks, and Spam Reports")
	unsubscribeAPIRespBody := ApiGet("unsubscribes", username, password)
	bounceAPIRespBody := ApiGet("bounces", username, password)
	invalidemailsAPIRespBody := ApiGet("invalidemails", username, password)
	blocksAPIRespBody := ApiGet("blocks", username, password)
	spamreportsAPIRespBody := ApiGet("spamreports", username, password)

	undesireables := []map[string]string{unsubscribeAPIRespBody, bounceAPIRespBody, invalidemailsAPIRespBody, blocksAPIRespBody, spamreportsAPIRespBody}

	responce := mapMerge(undesireables)

	fmt.Println("Removing undesireables from the lists")
	outputtedList := compareLists(wholeList, responce)

	fmt.Println("Building your list")
	zipOutput(username, outputtedList, responce)

	fmt.Println("Done :)")
}
