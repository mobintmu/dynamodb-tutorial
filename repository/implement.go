package repository

import (
	// "github.com/aws/aws-sdk-go/aws"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	// "github.com/aws/aws-sdk-go/service/s3"
)

type DB struct {
	session  *session.Session
	dynamodb *dynamodb.DynamoDB
}

func New() *DB {

	session, err := session.NewSession(&aws.Config{
		Region:      aws.String("localhost"),
		Endpoint:    aws.String("http://localhost:8000"),
		Credentials: credentials.NewStaticCredentials("AKID", "DUMMYIDEXAMPLE", "DUMMYIDEXAMPLE"),
	})

	if err != nil {
		fmt.Println(err)
		return nil
	}

	// Create DynamoDB client
	svc := dynamodb.New(session)

	return &DB{
		session:  session,
		dynamodb: svc,
	}
}

func (db *DB) GetList() {
	// create the input configuration instance
	input := &dynamodb.ListTablesInput{}

	fmt.Printf("Tables:\n")

	for {
		// Get the list of tables
		result, err := db.dynamodb.ListTables(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeInternalServerError:
					fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return
		}

		for _, n := range result.TableNames {
			fmt.Println(*n)
		}

		// assign the last read table name as the start for our next call to the ListTables function
		// the maximum number of table names returned in a call is 100 (default), which requires us to make
		// multiple calls to the ListTables function to retrieve all table names
		input.ExclusiveStartTableName = result.LastEvaluatedTableName

		if result.LastEvaluatedTableName == nil {
			break
		}
	}
}
