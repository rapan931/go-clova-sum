// clova extension for aws lambda
// build on windows
// > set GOOS=linux
// > set GOARCH=amd64
// > set CGO_ENABLED=0
// > go build -o sumQuiz go-clova-sum-quiz.go & build-lambda-zip.exe -o sumQuiz.zip sumQuiz

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type ClovaRequest struct {
	// Version string `json:"version"`
	Session struct {
		New               bool   `json:"new"`
		SessionID         string `json:"sessionId"`
		SessionAttributes struct {
			Question string `json:"question"`
			Correct  int    `json:"correct"`
		} `json:"sessionAttributes"`
		User struct {
			UserID      string `json:"userId"`
			AccessToken string `json:"accessToken"`
		} `json:"user"`
	} `json:"session"`
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
			Question string `json:"question"`
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

type SumQuizSlots struct {
	SumQuizLevel struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"sumQuizLevel"`
	SumQuiz struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"sumQuiz"`
}

type SumQuizAnswerSlots struct {
	Answer struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"answer"`
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

func SumQuiz(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	reqJsonBytes := ([]byte)(request.Body)
	clovaRequest := new(ClovaRequest)

	var err error
	if err = json.Unmarshal(reqJsonBytes, clovaRequest); err != nil {
		log.Println("[ERROR]", err)
		return events.APIGatewayProxyResponse{}, nil
	}

	response := NewClovaResponse()

	var text string
	sumQuizSlots := new(SumQuizSlots)
	sumQuizAnswerSlots := new(SumQuizAnswerSlots)

	// var result int
	switch clovaRequest.Request.Type {
	case "LaunchRequest":
		text = "たし算クイズ出して。と話してください。"

	case "IntentRequest":
		if clovaRequest.Request.Intent.Name == "SumQuizIntent" {

			if err = json.Unmarshal(clovaRequest.Request.Intent.Slots, sumQuizSlots); err != nil {
				log.Println("[ERROR]", err)
				return events.APIGatewayProxyResponse{}, nil
			}

			text = "1足す1は何？"

			response.Body.SessionAttributes.Question = text
			response.Body.SessionAttributes.Correct = 2
			response.Body.Response.ShouldEndSession = false

		} else if clovaRequest.Request.Intent.Name == "SumQuizAnswerIntent" {

			if err = json.Unmarshal(clovaRequest.Request.Intent.Slots, sumQuizAnswerSlots); err != nil {
				log.Println("[ERROR]", err)
				return events.APIGatewayProxyResponse{}, nil
			}

			var answer int
			if answer, err = strconv.Atoi(sumQuizAnswerSlots.Answer.Value); err != nil {
				log.Println("[ERROR]", answer)
				text = "すみません。理解できませんでした。"
				break
			}

			if answer == clovaRequest.Session.SessionAttributes.Correct {
				text = "正解！"
				response.Body.Response.ShouldEndSession = true
			} else {
				log.Println("[Info] Correct", clovaRequest.Session.SessionAttributes.Correct)
				log.Println("[Info] answer", answer)
				text = fmt.Sprintf("残念！もう一度！ %s", clovaRequest.Session.SessionAttributes.Correct)
			}

		} else {
			log.Println("[ERROR]", "Intent request parse error.")
			text = "すみません。理解できませんでした。"
		}

	default:
		log.Println("[ERROR]", "Request type parse error.")
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
