package tasks_service

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"regexp"
)

// Task has fields for the DynamoDB keys (Title)
type Task struct {
	Title string `json:"title"`
}

// List wraps up the DynamoDB calls to list all tasks
func List() ([]Task, error) {
	// Build the Dynamo client object
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	tasks := []Task{}

	// Get back the task title
	proj := expression.NamesList(expression.Name("title"))

	expr, err := expression.NewBuilder().WithProjection(proj).Build()

	if err != nil {
		fmt.Println("Got error building expression:")
		fmt.Println(err.Error())
		return tasks, err
	}

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("tasks"),
	}

	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)
	fmt.Println("Result", result)

	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println((err.Error()))
		return tasks, err
	}

	numItems := 0
	for _, i := range result.Items {
		item := Task{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
			return tasks, err
		}

		fmt.Println("Title: ", item.Title)
		tasks = append(tasks, item)
		numItems++
	}

	fmt.Println("Found", numItems, "tasks")
	if err != nil {
		fmt.Println(err.Error())
		return tasks, err
	}

	return tasks, nil
}

// Insert write a task to DynamoDB
func Insert(title string) (Task, error) {
	// Create the dynamo client object
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	// create Task object
	var thisItem = Task{
		Title: title,
	}

	// Take out non-alphanumberic except space characters from the title for easier slug building on reads
	reg, err := regexp.Compile("[^a-zA-Z0-9\\s]+")
	thisItem.Title = reg.ReplaceAllString(thisItem.Title, "")
	fmt.Println("Task to add:", thisItem)

	// Marshall the Item into a Map DynamoDB can deal with
	av, err := dynamodbattribute.MarshalMap(thisItem)
	if err != nil {
		fmt.Println("Got error marshalling map:")
		fmt.Println(err.Error())
		return thisItem, err
	}

	// Create Item in table and return
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("tasks"),
	}
	_, err = svc.PutItem(input)
	return thisItem, err

}
