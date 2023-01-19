package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/fatih/color"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func ReadJson(path string) map[string]interface{} {
readJsonAgain:

	jsonFile, err := os.Open(path)
	if err != nil {
		goto readJsonAgain
	}
	defer jsonFile.Close()

	b, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
	json.Unmarshal([]byte(b), &result)
	return result
}

func GetEmail(n int, domain string) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b) + domain
}
func RandInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}
func GetName() string {
	c := &http.Client{Timeout: time.Duration(30) * time.Second}
	r, err := c.Get("http://apis.kahoot.it/namerator")
	if err != nil {
		// log.Fatalln(err.Error())
	}
	defer r.Body.Close()
	var data map[string]interface{}
	json.NewDecoder(r.Body).Decode(&data)
	return data["name"].(string)

}
func Generate(email_domain string, password string, clients int) {
	var i int
restart:

	client := &http.Client{Timeout: time.Duration(30) * time.Second}
	gender := RandInt(1, 2)
	name := GetName()
	email := GetEmail(7, email_domain)
	// fmt.Println(gender)
	// fmt.Println(name)
	// fmt.Println(email)
	jsonData := map[string]interface{}{
		"account_details": map[string]interface{}{
			"birthdate": "2000-02-23",
			"consent_flags": map[string]interface{}{
				"eula_agreed":       true,
				"send_email":        true,
				"third_party_email": true,
			},
			"display_name": name,
			"email_and_password_identifier": map[string]interface{}{
				"email":    email,
				"password": password,
			},
			"gender": gender,
		},
		"callback_uri": "https://www.spotify.com/signup/challenge?forward_url=https%3A%2F%2Faccounts.spotify.com%2Fen%2Fstatus&locale=in-en",
		"client_info": map[string]interface{}{
			"api_key":         "a1e486e2729f46d6bb368d6b2bcda326",
			"app_version":     "v2",
			"capabilities":    []int{1},
			"installation_id": "432416f7-a0e9-4d0c-932a-28fab4720dce",
			"platform":        "www",
		},
		"tracking": map[string]interface{}{
			"creation_flow":  "",
			"creation_point": "https://www.spotify.com/in-en/signup/",
			"referrer":       "",
		},
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		fmt.Println(err)
		goto restart
	}
	request, err := http.NewRequest("POST", "https://spclient.wg.spotify.com/signup/public/v2/account/create", bytes.NewBuffer(jsonBytes))
	request.Header.Set("Host", "spclient.wg.spotify.com")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("sec-ch-ua", `"Not_A Brand";v="99", "Brave";v="109", "Chromium";v="109"`)
	request.Header.Set("sec-ch-ua-platform", `"Windows"`)
	request.Header.Set("sec-ch-ua-mobile", "?0")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "*/*")
	request.Header.Set("Sec-GPC", "1")
	request.Header.Set("Accept-Language", "en-US,en;q=0.7")
	request.Header.Set("Origin", "https://www.spotify.com")
	request.Header.Set("Sec-Fetch-Site", "same-site")
	request.Header.Set("Sec-Fetch-Mode", "cors")
	request.Header.Set("Sec-Fetch-Dest", "empty")
	request.Header.Set("Referer", "https://www.spotify.com/")
	// request.Header.Set("Accept-Encoding", "gzip, deflate, br")
	if err != nil {
		goto restart
	}
	response, err := client.Do(request)
	if err != nil {
		goto restart
	}
	// fmt.Println(response.StatusCode)
	defer response.Body.Close()
	tokenRequest, err := http.NewRequest("GET", "https://open.spotify.com/get_access_token?reason=transport&productType=web_player", nil)
	tokenRequest.Header.Set("Accept", "application/json")
	// tokenRequest.Header.Add("Accept-Encoding", "gzip, deflate, br")
	tokenRequest.Header.Add("Accept-Language", "en")
	tokenRequest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	tokenRequest.Header.Add("Spotify-App-Version", "1.1.52.204.ge43bc405")
	tokenRequest.Header.Add("App-Platform", "WebPlayer")
	tokenRequest.Header.Add("Host", "open.spotify.com")
	tokenRequest.Header.Add("Referer", "https://open.spotify.com/")
	if err != nil {
		goto restart
	}
	tokenResponse, err := client.Do(tokenRequest)
	if err != nil {
		goto restart
	}

	d, _ := ioutil.ReadAll(tokenResponse.Body)
	m := string(d)
	defer tokenResponse.Body.Close()
	var data map[string]interface{}
	err = json.Unmarshal([]byte(m), &data)
	if err != nil {
		goto restart
	}
	// fmt.Println(data["accessToken"])
	file, err := os.OpenFile("data/accounts.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		goto restart
	}
	defer file.Close()
	_, err = file.WriteString(email + ":" + password + ":" + data["accessToken"].(string) + "\n")
	if err != nil {
		goto restart
	}
	redText := color.RedString("INFO :")
	GreenText := color.GreenString("[CREATED]")
	magentaText := color.MagentaString(email + ":" + password + ":")
	BlueText := color.BlueString(data["accessToken"].(string))
	i = i + 1
	fmt.Println(redText, GreenText, magentaText, BlueText)
	cmd := exec.Command("cmd", "/c", "title", "[Spotify Account Creator] - Clients :"+strconv.Itoa(clients)+" | Dev: @icecastcve")
	cmd.Run()
	goto restart
}

func main() {
	cmd := exec.Command("cmd", "/c", "title", "[Spotify Account Creator] - [Reading Config] | Dev: @icecastcve")
	cmd.Run()

	jsonFile, err := os.Open("config.json")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer jsonFile.Close()

	b, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
	json.Unmarshal([]byte(b), &result)
	cmd = exec.Command("cmd", "/c", "title", "[Spotify Account Creator] - [Parsed Config] | Dev: @icecastcve")
	cmd.Run()
	k, _ := strconv.Atoi(result["settings"].(map[string]interface{})["threads"].(string))
	// fmt.Println(k)
	// color.Red("Spotify Account Creator")
	a := make(chan int)
	for i := 0; i < k; i++ {

		go func() {

			for {
				Generate(result["settings"].(map[string]interface{})["email_domain"].(string), result["accounts"].(map[string]interface{})["password"].(string), k)
			}
		}()
	}
	<-a
}
