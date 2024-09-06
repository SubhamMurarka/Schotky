package Dynamo

import (
	"fmt"
	"log"

	"github.com/SubhamMurarka/Schotky/Config"
	"github.com/aws/aws-dax-go/dax"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	awsRegion   = Config.Cfg.AwsRegion
	TableName   = Config.Cfg.TableName
	DaxEndpoint = Config.Cfg.DaxEndpoint
	UseDax      = Config.Cfg.UseDax
	ExpiryTime  = Config.Cfg.ExpiryTime
)

// Define an interface for the DynamoDB client
type DynamoDaxAPI interface {
	InsertItem(shortURL, longURL string) (*dynamodb.PutItemOutput, error)
	SelectItem(shortURL string) (*dynamodb.GetItemOutput, error)
	CreateTable()
	EnableTTL()
}

// Struct for DynamoDB client
type DynamoDaxClient struct {
	DynamoClient *dynamodb.DynamoDB
	DaxClient    *dax.Dax
}

// Constructor for DynamoDBClient
func NewDynamoDaxClient() DynamoDaxAPI {
	var err error

	if awsRegion == "" {
		log.Fatalf("AWS_REGION environment variable is not set")
	}

	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(awsRegion),
		Endpoint: aws.String("http://localhost:8003"),
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	var daxSvc *dax.Dax

	if UseDax == "true" {
		conf := dax.DefaultConfig()
		conf.HostPorts = []string{DaxEndpoint}
		conf.Region = awsRegion
		daxSvc, err = dax.New(conf)
		if err != nil {
			log.Println("Failed to create DaxClient: ", err)
			// Optionally, handle the error (e.g., return a default DynamoDB client if DAX fails)
		}
	}

	return &DynamoDaxClient{
		DynamoClient: dynamodb.New(sess),
		DaxClient:    daxSvc,
	}
}

// Check if a table exists
func (d *DynamoDaxClient) checkTableExists() bool {
	input := &dynamodb.ListTablesInput{}
	result, err := d.DynamoClient.ListTables(input)
	if err != nil {
		log.Fatalf("Failed to list tables: %v", err)
	}

	for _, name := range result.TableNames {
		if *name == TableName {
			return true
		}
	}

	return false
}

// Create a DynamoDB table with the specified schema
func (d *DynamoDaxClient) CreateTable() {
	if d.checkTableExists() {
		fmt.Println("Table already exists")
		return
	}

	input := &dynamodb.CreateTableInput{
		TableName: aws.String(TableName),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ShortURL"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ShortURL"),
				KeyType:       aws.String("HASH"),
			},
		},
		BillingMode: aws.String("PAY_PER_REQUEST"), // Set to On-Demand mode
	}

	_, err := d.DynamoClient.CreateTable(input)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	fmt.Println("Table created successfully")
}

// Enable TTL for the ExpiryTime attribute in the DynamoDB table

func (d *DynamoDaxClient) EnableTTL() {
	// Get the current TTL status
	input := &dynamodb.DescribeTimeToLiveInput{
		TableName: aws.String(TableName),
	}

	result, err := d.DynamoClient.DescribeTimeToLive(input)
	if err != nil {
		fmt.Printf("failed to describe TTL: %v", err)
		return
	}

	if result.TimeToLiveDescription != nil && *result.TimeToLiveDescription.TimeToLiveStatus == dynamodb.TimeToLiveStatusEnabled {
		fmt.Println("TTL is already enabled")
		return
	}

	// Enable TTL if not enabled
	ttlInput := &dynamodb.UpdateTimeToLiveInput{
		TableName: aws.String(TableName),
		TimeToLiveSpecification: &dynamodb.TimeToLiveSpecification{
			AttributeName: aws.String("ExpiryTime"),
			Enabled:       aws.Bool(true),
		},
	}

	_, err = d.DynamoClient.UpdateTimeToLive(ttlInput)
	if err != nil {
		fmt.Printf("failed to enable TTL: %v", err)
		return
	}

	fmt.Println("TTL has been enabled successfully")
}

// Insert a new item into DynamoDB or DAX
func (d *DynamoDaxClient) InsertItem(shortURL, longURL string) (*dynamodb.PutItemOutput, error) {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item: map[string]*dynamodb.AttributeValue{
			"ShortURL": {
				S: aws.String(shortURL),
			},
			"LongURL": {
				S: aws.String(longURL),
			},
			"ExpiryTime": {
				N: aws.String(fmt.Sprintf("%d", ExpiryTime)),
			},
		},
	}

	// Attempt to put the item in DynamoDB
	result, err := d.DynamoClient.PutItem(input)
	if err != nil {
		log.Printf("Failed to put item: %v", err)
		return nil, err
	}

	fmt.Println("Item added successfully")
	return result, nil
}

// Retrieve an item from DynamoDB or DAX
func (d *DynamoDaxClient) SelectItem(shortURL string) (*dynamodb.GetItemOutput, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ShortURL": {
				S: aws.String(shortURL),
			},
		},
	}

	var result *dynamodb.GetItemOutput
	var err error

	if UseDax == "true" && d.DaxClient != nil {
		result, err = d.DaxClient.GetItem(input)
		if err != nil {
			log.Printf("DAX failed to get item: %v", err)
		}
	}

	// Fallback to DynamoDB if DAX is disabled or fails
	if result == nil || err != nil {
		result, err = d.DynamoClient.GetItem(input)
		if err != nil {
			log.Printf("Failed to get item from DynamoDB: %v", err)
			return nil, err
		}
	}

	if result.Item == nil {
		fmt.Println("Could not find the item.")
		return nil, nil
	}

	fmt.Printf("Retrieved item: %v\n", result.Item)
	return result, nil
}
