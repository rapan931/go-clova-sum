// clova extension for aws lambda
// build on windows
// > set GOOS=linux
// > set GOARCH=amd64
// > set CGO_ENABLED=0
// > go build -o sum-quiz go-clova-sum-quiz.go & build-lambda-zip.exe -o sum-quiz.zip sum-quiz

package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type ClovaRequest struct {
	// Version string `json:"version"`
	// Session struct {
	// 	New               bool `json:"new"`
	// 	SessionID string `json:"sessionId"`
	// 	User      struct {
	// 		UserID      string `json:"userId"`
	// 		AccessToken string `json:"accessToken"`
	// 	} `json:"user"`
	// } `json:"session"`
	// Context struct {
	// 	System struct {
	// 		Application struct {
	// 			ApplicationID string `json:"applicationId"`
	// 		} `json:"application"`
	// 		User struct {
	// 			UserID      string `json:"userId"`
	// 			AccessToken string `json:"accessToken"`
	// 		} `json:"user"`
	// 		Device struct {
	// 			DeviceID string `json:"deviceId"`
	// 			Display  struct {
	// 				Size         string `json:"size"`
	// 				Orientation  string `json:"orientation"`
	// 				Dpi          int    `json:"dpi"`
	// 				ContentLayer struct {
	// 					Width  int `json:"width"`
	// 					Height int `json:"height"`
	// 				} `json:"contentLayer"`
	// 			} `json:"display"`
	// 		} `json:"device"`
	// 	} `json:"System"`
	// } `json:"context"`
	SessionAttributes struct {
		Question string `json:"string"`
		Correct  int    `json:"correct"`
	} `json:"sessionAttributes"`
	Request struct {
		Type   string `json:"type"`
		Intent struct {
			Name  string          `json:"name"`
			Slots json.RawMessage `json:"slots"`
		} `json:"intent"`
	} `json:"request"`
}

type ClovaResponse struct {
	StatusCode int `json:"statusCode"`
	Headers    struct {
		ContentType string `json:"Content-Type"`
	} `json:"headers"`
	Body struct {
		Version           string `json:"version"`
		SessionAttributes struct {
			Question string `json:"string"`
			Correct  int    `json:"correct"`
		} `json:"sessionAttributes"`
		Response struct {
			OutputSpeech struct {
				Type   string `json:"type"`
				Values struct {
					Type  string `json:"type"`
					Lang  string `json:"lang"`
					Value string `json:"value"`
				} `json:"values"`
			} `json:"outputSpeech"`
			Card             struct{}   `json:"card"`
			Directives       []struct{} `json:"directives"`
			Reprompt         struct{}   `json:"reprompt"`
			ShouldEndSession bool       `json:"shouldEndSession"`
		} `json:"response"`
	} `json:"body"`
}

func NewClovaResponse() *ClovaResponse {
	response := ClovaResponse{}
	response.Body.Version = "1.0"
	response.Body.Response.ShouldEndSession = false
	response.Body.Response.OutputSpeech.Type = "SimpleSpeech"
	response.Body.Response.OutputSpeech.Values.Type = "PlainText"
	response.Body.Response.OutputSpeech.Values.Lang = "ja"
	return &response
}

// func SumQuizResponseText() string {
// 	return "hoge"
// }
//
// func SumQuizAnswerResponseText() string {
// 	return "hoge"
// }

func SumQuiz(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	reqJsonBytes := ([]byte)(request.Body)
	clovaRequest := new(ClovaRequest)

	var err error
	if err = json.Unmarshal(reqJsonBytes, clovaRequest); err != nil {
		log.Println("[ERROR]", err)
	}

	response := NewClovaResponse()

	var text string
	// var result int
	switch clovaRequest.Request.Type {
	case "LaunchRequest":
		text = "たし算クイズ出して。と話してください。"

	case "IntentRequest":
		if clovaRequest.Request.Intent.Name == "SumQuizIntent" {
			text = ""
		} else if clovaRequest.Request.Intent.Name == "SumQuizAnswerIntent" {
			text = ""
		} else {
			text = "すみません。理解できませんでした。"
		}

	default:
		log.Println("[ERROR]", "Intent request parse error.")
		text = "すみません。理解できませんでした。"
		break
	}

	response.Body.Response.OutputSpeech.Values.Value = text

	resJsonBytes, _ := json.Marshal(response.Body)
	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		StatusCode: 200,
		Body:       string(resJsonBytes),
	}, nil
}

func main() {
	log.SetFlags(log.Lshortfile)
	lambda.Start(SumQuiz)
}
