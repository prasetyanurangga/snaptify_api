// Package spotify:
// go-spotify provides an easy-to-use API
// to access Spotify's WEB API
package imageLabel

import (
	"encoding/json"
	"errors"
	"fmt"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/parnurzeal/gorequest"
)

const (
	BASE_URL = "https://image-labeling1.p.rapidapi.com/img/label"
)


type ImageLabel struct {
	apiKey     string
}

func New(apiKey string) ImageLabel {
	imageLabel := ImageLabel{apiKey: apiKey}
	return imageLabel
}

func (imageLabel *ImageLabel)  Get( data map[string]interface{}) ([]byte, []error) {

	targetURL := createTargetURL()

	request := gorequest.New()

	request.Post(targetURL)
	request.Set("content-type", "application/json")
	request.Set("x-rapidapi-host", "image-labeling1.p.rapidapi.com")
    request.Set("x-rapidapi-key", imageLabel.apiKey)

	// Add the data to the request if it
	// isn't null
	if data != nil {
		jsonData, _ := getJsonBytesFromMap(data)
		if jsonData != nil {
			request.Send(string(jsonData))
		}
	}

	_, body, errs := request.End()

	fmt.Println(errs)


	result := []byte(body)
	if unauthorizedResponse(result) {
		result = nil
		errs = []error{
			errors.New("Authorization Error. Make sure you called Spotify.Authorize() method!"),
			errors.New(body)}
	}

	return result, errs
}

// Checks for the response content to see if we
// received a not authorized error.
func unauthorizedResponse(body []byte) bool {

	// Parse response to simplejson object
	js, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Println("[unauthorizedResponse] Error parsing Json!")
		return true
	}

	// check whether we got an error or not.
	_, exists := js.CheckGet("error")
	if exists {
		return true
	}

	return false
}

func createTargetURL() string {
	result := fmt.Sprintf("%s", BASE_URL)
	return result
}


// Extracts Json Bytes from map[string]interface
func getJsonBytesFromMap(data map[string]interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Invalid data object, can't parse to json:")
		fmt.Println("Error:", err)
		fmt.Println("Data:", data)
		return nil, err
	}
	return jsonData, nil
}