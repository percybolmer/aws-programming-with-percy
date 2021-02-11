package main

import (
	"fmt"
	"os"

	// aws-lambda package has to be imported
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// a variable that is a dynamoDB client
var dynaDB *dynamodb.DynamoDB

var latestID int

// User is a simple struct to accept as input
type User struct {
	Name     string `json:"name"`
	LastName string `json:"lastname"`
	ID       int    `json:"ID"`
}

// Greetings is the response to a user entering
type Greetings struct {
	Greet string `json:"greet"`
}

// String implements the stringer interface for user
func (u User) String() string {
	return fmt.Sprintf("Welcome %s %s", u.Name, u.LastName)
}

// GreetVisitor is a lambda function that will greet a visitor
func GreetVisitor(event User) (Greetings, error) {
	// Add an ID to the user
	event.ID = latestID + 1
	latestID++
	// Marhsal the user into a Map
	item, err := dynamodbattribute.MarshalMap(event)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// Create a dynamoDB input item, set the marshalled data and the users table
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String("Users"),
	}
	// Send the Item into the DB
	_, err = dynaDB.PutItem(input)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return Greetings{Greet: event.String()}, nil
}

func main() {
	// Set up a Session to dynamoDB. If you want to setup access from CLI then you need to generate User keys.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// Create a dynamoDB client
	dynaDB = dynamodb.New(sess)

	// Lambda start is used to handle the function
	lambda.Start(GreetVisitor)
}
