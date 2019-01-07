// https://github.com/aws/aws-lambda-go#for-developers-on-windows
// https://mholt.github.io/json-to-go/
// https://dev.classmethod.jp/cloud/aws/aws-lambda-supports-go/
// https://docs.aws.amazon.com/ja_jp/lambda/latest/dg/go-programming-model-handler-types.html
//
// exapmple
// > set GOOS=linux
// > set GOARCH=amd64
// > set CGO_ENABLED=0
// > go build -o main main.go
// > %USERPROFILE%\Go\bin\build-lambda-zip.exe -o main.zip main

package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

// {
//   "context": {
//     "AudioPlayer": {
//       "offsetInMilliseconds": {{number}},
//       "playerActivity": {{string}},
//       "stream": {{AudioStreamInfoObject}},
//       "totalInMilliseconds": {{number}},
//     },
//     "System": {
//       "application": {
//         "applicationId": {{string}}
//       },
//       "device": {
//         "deviceId": {{string}},
//         "display": {
//           "contentLayer": {
//             "width": {{number}},
//             "height": {{number}}
//           },
//           "dpi": {{number}},
//           "orientation": {{string}},
//           "size": {{string}}
//         }
//       },
//       "user": {
//         "userId": {{string}},
//         "accessToken": {{string}}
//       }
//     }
//   },
//   "request": {{object}},
//   "session": {
//     "new": {{boolean}},
//     "sessionAttributes": {{object}},
//     "sessionId": {{string}},
//     "user": {
//       "userId": {{string}},
//       "accessToken": {{string}}
//     }
//   },
//   "version": {{string}}
// }
type MyEvent struct {
	Version string `json:"version"`
	Session struct {
		New               bool `json:"new"`
		SessionAttributes struct {
			SumX int `json:"sumx"`
			SumY int `json:"sumy"`
		} `json:"sessionAttributes"`
		SessionID string `json:"sessionId"`
		User      struct {
			UserID      string `json:"userId"`
			AccessToken string `json:"accessToken"`
		} `json:"user"`
	} `json:"session"`
	Context struct {
		System struct {
			Application struct {
				ApplicationID string `json:"applicationId"`
			} `json:"application"`
			User struct {
				UserID      string `json:"userId"`
				AccessToken string `json:"accessToken"`
			} `json:"user"`
			Device struct {
				DeviceID string `json:"deviceId"`
				Display  struct {
					Size         string `json:"size"`
					Orientation  string `json:"orientation"`
					Dpi          int    `json:"dpi"`
					ContentLayer struct {
						Width  int `json:"width"`
						Height int `json:"height"`
					} `json:"contentLayer"`
				} `json:"display"`
			} `json:"device"`
		} `json:"System"`
	} `json:"context"`
	Request struct {
		Type string `json:"type"`
	} `json:"request"`
}

type MyResponse struct {
	Message string `json:"Answer:"`
}

func MySum(event MyEvent) (MyResponse, error) {
	return MyResponse{Message: fmt.Sprintf("Version %s!!", event.Version)}, nil
	// return MyResponse{Message: "test"}, nil
}
func main() {
	lambda.Start(MySum)
}
