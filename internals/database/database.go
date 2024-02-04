package database

import (
	"fmt"
	"log"
	"os"

	"trip/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func connect() *dynamodb.DynamoDB {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		// TODO
		panic(err)
	}
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region: aws.String(os.Getenv("AWS_REGION")),
		},
	}))
	return dynamodb.New(sess)
}

func CreateTrip(dto types.Trip) (types.Trip, error) {
	db := connect()
	trip := types.Trip{
		Uuid:        string(uuid.New().String()),
		Origin:      dto.Origin,
		Destination: dto.Destination,
		Status:      dto.Status,
		Price:       dto.Price,
	}
	av, marshalErr := dynamodbattribute.MarshalMap(trip)
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(`trip`),
	}
	if marshalErr != nil {
		log.Fatalf("Got error marshalling new trip item: %s", marshalErr)
		return trip, marshalErr
	}
	_, putItemErr := db.PutItem(input)
	if putItemErr != nil {
		log.Fatalf("Got error calling PutItem: %s", putItemErr)
		return trip, putItemErr
	}
	return trip, nil
}

func ListTrips(dto types.ListTripsParamsDto) ([]types.Trip, error) {
	db := connect()

	filters := expression.Name("Status").Equal(expression.Value(dto.Status))

	proj := expression.NamesList(
		expression.Name("Origin"),
		expression.Name("Destination"),
		expression.Name("Status"),
		expression.Name("Price"),
	)

	expr, builderErr := expression.NewBuilder().WithFilter(filters).WithProjection(proj).Build()
	if builderErr != nil {
		log.Fatalf("Got error building expression: %s", builderErr)
		return nil, builderErr
	}

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(`trip`),
	}

	// Make the DynamoDB Query API call
	result, queryErr := db.Scan(params)
	if queryErr != nil {
		log.Fatalf("Query API call failed: %s", queryErr)
		return nil, queryErr
	}

	trips := []types.Trip{}
	for _, i := range result.Items {
		trip := types.Trip{}
		err := dynamodbattribute.UnmarshalMap(i, &trip)

		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
			return trips, err
		}
		trips = append(trips, trip)
	}

	return trips, nil
}
