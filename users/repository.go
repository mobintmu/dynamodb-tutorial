package repository

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DB struct {
	session        *session.Session
	dynamodb       *dynamodb.DynamoDB
	tableName      string
	tgIndex        string
	consistentRead bool
	referrerIndex  string
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
		session:        session,
		dynamodb:       svc,
		tableName:      "Users", //Table name
		tgIndex:        "TgIdIndex",
		referrerIndex:  "ReferrerIndex",
		consistentRead: true,
	}
}

func (db *DB) CreateTable() error {

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("user_id"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("tg_id"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("referrer"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("created_at"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("user_id"),
				KeyType:       aws.String("HASH"),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("TgIdIndex"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("tg_id"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
			},
			{
				IndexName: aws.String("ReferrerIndex"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("referrer"),
						KeyType:       aws.String("HASH"),
					},
					{
						AttributeName: aws.String("created_at"),
						KeyType:       aws.String("RANGE"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
			},
		},
		BillingMode: aws.String("PAY_PER_REQUEST"),
		TableName:   aws.String(db.tableName),
	}

	_, err := db.dynamodb.CreateTable(input)
	if err != nil {
		log.Fatalf("Got error calling CreateTable: %s", err)
		return err
	}
	return nil
}

func (db *DB) GetList() error {
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
				fmt.Println(err.Error())
			}
			return err
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

	return nil
}

func (db *DB) TableIsExist(tableName string) bool {
	input := &dynamodb.ListTablesInput{}
	result, err := db.dynamodb.ListTables(input)
	if err != nil {
		log.Fatalf("Got error calling ListTables: %s", err)
		return false
	}

	for _, n := range result.TableNames {
		if *n == tableName {
			return true
		}
	}
	return false
}

func (db *DB) AddItem(user User) error {

	attributeValue, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		log.Fatal("Got error marshalling new user item:", err)
		return err
	}
	input := &dynamodb.PutItemInput{
		Item:      attributeValue,
		TableName: aws.String(db.tableName),
	}
	// Put the item in the table
	_, err = db.dynamodb.PutItem(input)
	if err != nil {
		log.Fatal("Got error calling PutItem:", err)
		return err
	}

	return nil
}

func (db *DB) GetByUserID(userID string) (*User, error) {

	result, err := db.dynamodb.GetItem(&dynamodb.GetItemInput{
		ConsistentRead: aws.Bool(db.consistentRead),
		Key: map[string]*dynamodb.AttributeValue{
			"user_id": {
				S: aws.String(userID),
			},
		},
		TableName: aws.String(db.tableName),
	})

	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
		return nil, err
	}

	if result.Item == nil {
		msg := "Could not find '" + userID + "'"
		return nil, errors.New(msg)
	}

	user := User{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		log.Fatalf("Got error un marshalling: %v", err)
		return nil, err
	}
	return &user, nil

}

// GetByTgID get item by Global Secondary Index (GSI)
func (db *DB) GetByTgID(tgID string) ([]User, error) {

	input := &dynamodb.QueryInput{
		TableName:              aws.String(db.tableName),
		IndexName:              aws.String(db.tgIndex), // Replace with your GSI name
		KeyConditionExpression: aws.String("#tg = :v"),
		ExpressionAttributeNames: map[string]*string{
			"#tg": aws.String("tg_id"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v": {
				S: aws.String(tgID),
			},
		},
	}

	// Perform the query
	result, err := db.dynamodb.Query(input)
	if err != nil {
		log.Fatalf("failed to query items, %v", err)
		return nil, err
	}

	if result.Items == nil {
		msg := "Could not find '" + tgID + "'"
		return nil, errors.New(msg)
	}

	users := []User{}
	for _, item := range result.Items {
		user := User{}
		err := dynamodbattribute.UnmarshalMap(item, &user)
		if err != nil {
			log.Fatalf("Got error un marshalling: %v", err)
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (db *DB) GetByReferrer(referrer string) ([]User, error) {

	input := &dynamodb.QueryInput{
		TableName:              aws.String(db.tableName),
		IndexName:              aws.String(db.referrerIndex), // Replace with your GSI name
		KeyConditionExpression: aws.String("#ref = :v"),
		ExpressionAttributeNames: map[string]*string{
			"#ref": aws.String("referrer"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v": {
				S: aws.String(referrer),
			},
		},
	}

	// Perform the query
	result, err := db.dynamodb.Query(input)
	if err != nil {
		log.Fatalf("failed to query items, %v", err)
		return nil, err
	}

	if result.Items == nil {
		msg := "Could not find '" + referrer + "'"
		return nil, errors.New(msg)
	}

	users := []User{}
	for _, item := range result.Items {
		user := User{}
		err := dynamodbattribute.UnmarshalMap(item, &user)
		if err != nil {
			log.Fatalf("Got error un marshalling: %v", err)
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil

}

func (db *DB) IncrementUserReferrer(userID string) error {

	input := &dynamodb.UpdateItemInput{

		Key: map[string]*dynamodb.AttributeValue{
			"user_id": {
				S: aws.String(userID),
			},
		},
		TableName: aws.String(db.tableName),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":inc":  {N: aws.String("1")},
			":time": {N: aws.String(strconv.FormatInt(time.Now().UTC().UnixNano(), 10))},
		},
		UpdateExpression: aws.String("SET counter_referrer = counter_referrer + :inc, updated_at = :time"),
		ReturnValues:     aws.String("UPDATED_NEW"), // Optional: To get updated item
	}

	value, err := db.dynamodb.UpdateItem(input)
	if err != nil {
		log.Fatalf("Got error calling UpdateItem: %s", err)
		return err
	}

	fmt.Println(value)

	return nil
}

func (db *DB) CountAllData() (int64, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String("Users"),
	}

	result, err := db.dynamodb.Scan(input)
	if err != nil {
		log.Fatalf("failed to scan table, " + err.Error())
		return 0, err
	}

	return *result.Count, nil
}
