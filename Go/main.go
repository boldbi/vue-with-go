package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//Set EmbedSecret key from Bold BI Server. Please refer this link(https://help.syncfusion.com/bold-bi/on-premise/site-settings/embed-settings)
var embedSecret = ""

//Enter your BoldBI Server credentials
var userMail = ""

func main() {
	http.HandleFunc("/getDetails", getEmbedDetails)
	log.Fatal(http.ListenAndServe(":8086", nil))

}

func getEmbedDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	if len(body) > 0 {
		if queryString, err := unmarshal(string(body)); err != nil {
			log.Println("error converting", err)
		} else {
			serverAPIUrl := queryString.(map[string]interface{})["dashboardServerApiUrl"].(string)
			embedQueryString := queryString.(map[string]interface{})["embedQuerString"].(string)
			embedQueryString += "&embed_user_email=" + userMail
			timeStamp := time.Now().Unix()
			embedQueryString += "&embed_server_timestamp=" + strconv.FormatInt(timeStamp, 10)
			signatureString, err := getSignatureUrl(embedQueryString)
			embedDetails := "/embed/authorize?" + embedQueryString + "&embed_signature=" + signatureString
			query := serverAPIUrl + embedDetails
			log.Println(query)
			result, err := http.Get(query)
			if err != nil {
				log.Println(err)
			}
			log.Println(result)
			response, err := ioutil.ReadAll(result.Body)
			if err != nil {
				log.Fatalln(err)
			}
			w.Write(response)
		}
	}
}

func getSignatureUrl(queryData string) (string, error) {
	encoding := ([]byte(embedSecret))
	messageBytes := ([]byte(queryData))
	hmacsha1 := hmac.New(sha256.New, encoding)
	hmacsha1.Write(messageBytes)
	sha := base64.StdEncoding.EncodeToString(hmacsha1.Sum(nil))
	return sha, nil
}

func unmarshal(data string) (interface{}, error) {
	var iface interface{}
	decoder := json.NewDecoder(strings.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&iface); err != nil {
		return nil, err
	}
	return iface, nil
}
