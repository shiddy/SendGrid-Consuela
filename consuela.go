package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Unsubscribe struct {
	Email string `json:"email"`
}

func ApiGet(eventName, username, password string) map[string]string {
	uri := "https://api.sendgrid.com/api/" + eventName + ".get.json?api_user=" + username + "&api_key=" + password
	resp, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}

	robots, _ := ioutil.ReadAll(resp.Body)

	var unsubscribes []Unsubscribe
	json.Unmarshal(robots, &unsubscribes)

	returnedMap := make(map[string]string)

	for i := 0; i < len(unsubscribes); i++ {
		returnedMap[unsubscribes[i].Email] = ""
	}

	return returnedMap
}

func mapMerge(setArray []map[string]string) map[string]string {
	mergedList := make(map[string]string)
	for i := 1; i < len(setArray); i++ {
		for k, _ := range setArray[i] {
			mergedList[k] = ""
		}
	}

	return mergedList
}

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

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func compareLists(wholeList, unsubscribeList map[string]string) map[string]string {
	for k, _ := range unsubscribeList {
		delete(wholeList, k)
	}

	return wholeList
}

func zipOutput1(username string, wholeList, responce map[string]string, undesireables []map[string]string) {
	undesireableNames := []string{"Unsubscribes", "Bounce", "Invalids", "Blocks", "Spam Reports"}
	wholelistname := username + "DONOTSEND.csv"
	donotsendname := username + "newlist.csv"

	outputFile, _ := os.Create(wholelistname)
	defer outputFile.Close()
	csvOutput := csv.NewWriter(outputFile)

	for i := 0; i < len(undesireables); i++ {
		csvOutput.Write([]string{undesireableNames[i]})
		for k, _ := range undesireables[i] {
			csvOutput.Write([]string{k})
			csvOutput.Flush()
		}
		csvOutput.Write([]string{})
	}

	newlistOutputFile, _ := os.Create(donotsendname)
	defer newlistOutputFile.Close()

	newlistCsvOutput := csv.NewWriter(newlistOutputFile)
	for k, _ := range wholeList {
		newlistCsvOutput.Write([]string{k})
		newlistCsvOutput.Flush()
	}

	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	w := zip.NewWriter(buf)

	// Add some files to the archive.
	var files = []struct {
		Name, Body string
	}{
		{wholelistname, "This archive contains your new list."},
		{donotsendname, "This archive contains the emails that should not be sent to again."},
	}
	for _, file := range files {
		f, err := w.Create(file.Name)
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.Write([]byte(file.Body))
		if err != nil {
			log.Fatal(err)
		}
	}

	// Make sure to check the error on Close.
	err := w.Close()
	if err != nil {
		log.Fatal(err)
	}
	//write the zipped file to the disk
	ioutil.WriteFile(username+"_.zip", buf.Bytes(), 0777)
}

func zipOutput(username string, wholeList, responce map[string]string) {
	wholelistname := username + "_newlist.csv"
	donotsendname := username + "_DONOTSEND.csv"

	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	w := zip.NewWriter(buf)

	// Add some files to the archive.
	var files = []struct {
		Name string
		Body map[string]string
	}{
		{wholelistname, wholeList},
		{donotsendname, responce},
	}

	for _, file := range files {
		f, err := w.Create(file.Name)
		if err != nil {
			log.Fatal(err)
		}
		for k, _ := range file.Body {
			_, err = f.Write([]byte(k + "\n"))
			// f.Write([]byte("\n"))
			if err != nil {
				log.Fatal(err)
			}
			// newlistCsvOutput.Write([]string{k})
		}
	}

	// Make sure to check the error on Close.
	err := w.Close()
	if err != nil {
		log.Fatal(err)
	}
	//write the zipped file to the disk
	ioutil.WriteFile(username+".zip", buf.Bytes(), 0777)
}

// Pulls all of the user's data from the webAPI
//   If you intend to use this with another list of unwanted emails
//   you would run the CsvToMap fucntion with your wanted email list
//   and run a compareLists on the two maps.
func main() {
	username := "USERNAME"
	password := "PASSWORD"

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
