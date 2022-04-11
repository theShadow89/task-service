package main

import (
	"bytes"
	"context"
	"encoding/json"
	tasksservice "github.com/scastoldi/tasks-service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request Request) (Response, error) {
	var buf bytes.Buffer

	var jsonBody map[string]interface{}

	log.Debug("Received body: ", request.Body)

	if err := json.Unmarshal([]byte(request.Body), &jsonBody); err != nil {
		body, _ := json.Marshal(map[string]interface{}{
			"message": "Invalid Request",
		})
		json.HTMLEscape(&buf, body)

		return Response{
			StatusCode:      500,
			IsBase64Encoded: false,
			Body:            buf.String(),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, err
	}

	var taskTitle = jsonBody["title"].(string)

	_, err := tasksservice.Insert(taskTitle)
	if err != nil {

		body, _ := json.Marshal(map[string]interface{}{
			"message": "Error on insert Task " + taskTitle,
		})

		log.Error(err)

		json.HTMLEscape(&buf, body)

		return Response{
			StatusCode:      500,
			IsBase64Encoded: false,
			Body:            buf.String(),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	body, err := json.Marshal(map[string]interface{}{
		"message": "Task inserted successfully!",
	})

	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
