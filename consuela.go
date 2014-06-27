package main

import (
	"bufio"
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

// func ApiCall(eventName, call, username, password string) []Unsubscribe {
// 	uri := "https://api.sendgrid.com/api/" + eventName + "." + call + ".json?api_user=" + username + "&api_key=" + password
// 	resp, err := http.Get(uri)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return resp
// }

func ApiGet(eventName, username, password string) map[string]string {
	uri := "https://api.sendgrid.com/api/" + eventName + ".get.json?api_user=" + username + "&api_key=" + password
	resp, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}

	robots, _ := ioutil.ReadAll(resp.Body)

	var unsubscribes []Unsubscribe
	json.Unmarshal(robots, &unsubscribes)

	// fmt.Println(unmarsh)

	returnedMap := make(map[string]string)

	// fmt.Println(len(unsubscribes))

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

func main() {
	username := "YOURUSERNAMEHERE"
	password := "YOURPASSWORDHERE"

	// fmt.Println("https://api.sendgrid.com/api/unsubscribes.get.json?api_user=" + username + "&api_key=" + password)

	fmt.Println("Opening your giant list")
	wholeList := CsvToMap("wholelist.csv")
	fmt.Println("Gathering your Unsubscribes, Bounces, Invalids, Blocks, and Spam Reports")
	unsubscribeAPIRespBody := ApiGet("unsubscribes", username, password)
	bounceAPIRespBody := ApiGet("bounces", username, password)
	invalidemailsAPIRespBody := ApiGet("invalidemails", username, password)
	blocksAPIRespBody := ApiGet("blocks", username, password)
	spamreportsAPIRespBody := ApiGet("spamreports", username, password)

	undesireables := []map[string]string{unsubscribeAPIRespBody, bounceAPIRespBody, invalidemailsAPIRespBody, blocksAPIRespBody, spamreportsAPIRespBody}
	undesireableNames := []string{"Unsubscribes", "Bounce", "Invalids", "Blocks", "Spam Reports"}

	responce := mapMerge(undesireables)

	fmt.Println("Outputing undesireables to 'DONOTSEND.csv'")
	outputFile, _ := os.Create("DONOTSEND.csv")
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

	// fmt.Println(responce)

	fmt.Println("Removing undesireables from the lists")
	outputtedList := compareLists(wholeList, responce)

	fmt.Println("Building your list")
	newlistOutputFile, _ := os.Create("newlist.csv")
	defer newlistOutputFile.Close()

	newlistCsvOutput := csv.NewWriter(newlistOutputFile)
	for k, _ := range outputtedList {
		newlistCsvOutput.Write([]string{k})
		newlistCsvOutput.Flush()
	}

	fmt.Println("Done :)")
}
