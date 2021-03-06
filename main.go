package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

// Read this to understand how to use the Google Sheets API https://developers.google.com/sheets/api/guides/concepts
// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, dir string) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := fmt.Sprintf("%v/token.json", dir)
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	// TODO: Note that sheet id is hardcoded
	spreadsheetID := "1j5zrdUcqJWCTp9_05DmUTZk3PyBKVsTv36BO26Wv3SE"
	hardcodedPerson := ""
	if len(os.Args) == 2 && os.Args[1] == "sheet" {
		open(fmt.Sprintf("https://docs.google.com/spreadsheets/d/%v", spreadsheetID))
		os.Exit(0)
	}
	if len(os.Args) == 3 && os.Args[1] == "chat" {
		hardcodedPerson = os.Args[2]
	}
	dir := fmt.Sprintf("%v/src/github.com/marianogappa/fsm", os.Getenv("GOPATH"))
	b, err := ioutil.ReadFile(fmt.Sprintf("%v/credentials.json", dir))
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config, dir)

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// TODO: Note that sheet name (i.e. 2018) is hardcoded
	readRange := "2018"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		log.Fatalf("No data found.")
	}

	// Special case for head row: need the field names
	// TODO: probably want to die if any of these are missing: ["name", "last_comm", "chat_link", "keep"]
	fields := []string{}
	fieldNameIndex := map[string]int{}
	nameFieldIndex := -1
	for i, field := range resp.Values[0] {
		strField := fmt.Sprintf("%v", field)
		if strField == "name" {
			nameFieldIndex = i
		}
		fieldNameIndex[strField] = i
		fields = append(fields, strField)
	}
	if nameFieldIndex == -1 {
		log.Fatal("Couldn't find name field")
	}

	// Construct a table structure from the get API result
	table := []map[string]string{}
	nameIndex := map[string]int{}
	for i, row := range resp.Values[1:] {
		tableRow := map[string]string{}
		for j, field := range row {
			tableRow[fields[j]] = fmt.Sprintf("%v", field)
		}
		table = append(table, tableRow)
		nameIndex[fmt.Sprintf("%v", row[nameFieldIndex])] = i
	}

	if len(os.Args) >= 2 && os.Args[1] == "log" {
		headN := math.MaxInt32
		if len(os.Args) == 3 {
			n, err := strconv.Atoi(os.Args[2])
			if err == nil && n > 0 {
				headN = n
			}
		}
		chatLogs := sortNamesByLastComm(table)
		for i := 0; i < headN && i < len(chatLogs); i++ {
			oldenTimes, _ := time.Parse("2006-01-02", "1971-01-01")
			if chatLogs[i].lastComm.Before(oldenTimes) {
				break
			}
			fmt.Printf("%v %v\n", chatLogs[i].lastComm.Format("2006-01-02"), chatLogs[i].name)
		}
		os.Exit(0)
	}

	// Decide which people have a last comm documented older than desired comm frequency.
	personsToTalkTo := []map[string]string{}
	if hardcodedPerson != "" {
		if _, ok := nameIndex[hardcodedPerson]; !ok {
			log.Fatalf("Couldn't find person %v\n", hardcodedPerson)
		}
		personsToTalkTo = append(personsToTalkTo, table[nameIndex[hardcodedPerson]])
	} else {
		personsToTalkTo = calculatePersonsToTalkTo(table)
	}

	// There might be more than one person to talk to. Choose one randomly.
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for _, i := range r.Perm(len(personsToTalkTo)) {
		if personsToTalkTo[i]["chat_link"] != "" {
			// Open their chat link
			open(personsToTalkTo[i]["chat_link"])

			// Document this comm by updating "last_comm" with today's date
			// Ideally we'd want to abstract out this table functionality with a robust library. One day.
			letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
			// TODO: investigate why +2. +1 makes sense to me (title row)
			// TODO: Note that sheet name (i.e. 2018) is hardcoded
			cell := fmt.Sprintf("2018!%v%v", string(letters[fieldNameIndex["last_comm"]]), nameIndex[personsToTalkTo[i]["name"]]+2)
			srv.Spreadsheets.Values.Update(
				spreadsheetID,
				cell,
				&sheets.ValueRange{Values: [][]interface{}{{time.Now().Format("2006-01-02")}}},
			).ValueInputOption("RAW").Do()

			// Only open one chat per run
			break
		}
	}
}

func calculatePersonsToTalkTo(table []map[string]string) []map[string]string {
	personsToTalkTo := []map[string]string{}
	for _, person := range table {
		lastComm, err := time.Parse("2006-01-02", person["last_comm"])
		if err != nil {
			lastComm, _ = time.Parse("2006-01-02", "1970-01-01")
		}
		dayFrequency, err := strconv.Atoi(person["frequency"])
		if err != nil || dayFrequency == 0 {
			continue
		}
		if lastComm.Add(time.Duration(24*dayFrequency) * time.Hour).Before(time.Now()) {
			personsToTalkTo = append(personsToTalkTo, person)
		}
	}
	return personsToTalkTo
}

type chatLog struct {
	name     string
	lastComm time.Time
}

func sortNamesByLastComm(table []map[string]string) []chatLog {
	chatLogs := []chatLog{}
	for _, person := range table {
		lastComm, err := time.Parse("2006-01-02", person["last_comm"])
		if err != nil {
			lastComm, _ = time.Parse("2006-01-02", "1970-01-01")
		}
		chatLogs = append(chatLogs, chatLog{name: person["name"], lastComm: lastComm})
	}
	sort.Slice(chatLogs, func(i, j int) bool { return chatLogs[i].lastComm.After(chatLogs[j].lastComm) })
	return chatLogs
}

func open(url string) {
	exec.Command("open", url).Run()
}
