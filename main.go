package main

import (
	"context"
		"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/apex/gateway/v2"
	"github.com/apex/log"
	jsonhandler "github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	units "github.com/docker/go-units"
)

type Record struct {
	Id         string    `dynamodbav:"id" json:"id"`
	IssueDate  time.Time `dynamodbav:"issueDate" json:"issueDate"`
	RandomWord string    `dynamodbav:"randomWord" json:"randomWord"`
}

func dynamoCloud() *dynamodb.Client {
	log.Info("cloud config")
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		return nil
	})
	if err != nil {
		log.WithError(err).Fatal("failed to load config")
	}
	return dynamodb.NewFromConfig(cfg)
}

func dynamoLocal() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.SharedConfigProfile = "mine"
		return nil
	})
	if err != nil {
		log.WithError(err).Fatal("failed to load config")
	}
	return dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.EndpointResolver = dynamodb.EndpointResolverFromURL("http://localhost:8000")
	})
}

func main() {

	client := dynamoCloud()

	log.Info("starting")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// https://aws.github.io/aws-sdk-go-v2/docs/code-examples/dynamodb/scanitems/
		records, err := client.Scan(context.TODO(), &dynamodb.ScanInput{
			TableName: aws.String(os.Getenv("TABLE_NAME")),
		})
		if err != nil {
			log.WithError(err).Fatal("couldn't get records")
		}

		var selection []Record
		err = attributevalue.UnmarshalListOfMaps(records.Items, &selection)
		if err != nil {
			log.WithError(err).Fatal("couldn't parse records")
		}

		for _, record := range selection {
			fmt.Fprintf(w, units.HumanDuration(time.Now().UTC().Sub(record.IssueDate))+" ago ")
			fmt.Fprintf(w, "%+v\n", record)
		}

	})

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {

		// Handle post only
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// parse body to a record
		var record Record
		err := json.NewDecoder(r.Body).Decode(&record)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.WithField("record", record).Info("adding")

		// marshall record into dynamo
		av, err := attributevalue.MarshalMap(record)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// what if issueDate is wrong?
		_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(os.Getenv("TABLE_NAME")),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// output how many records we have, using an empty ProjectionExpression
		records, err := client.Scan(context.TODO(), &dynamodb.ScanInput{
			TableName:            aws.String(os.Getenv("TABLE_NAME")),
			ProjectionExpression: aws.String("id"),
		})

		fmt.Fprintf(w, "added %d records\n", records.Count)

	})

	// get latest record
	// https://stackoverflow.com/a/12809659/4534
	http.HandleFunc("/latest", func(w http.ResponseWriter, r *http.Request) {
		log.Info("getting latest")

		// sort by issueDate
		records, err := client.Query(context.TODO(), &dynamodb.QueryInput{
			TableName:        aws.String(os.Getenv("TABLE_NAME")),
			ScanIndexForward: aws.Bool(true),
			Limit:            aws.Int32(1),
		})
		if err != nil {
			log.WithError(err).Fatal("couldn't get records")
		}

		var selection []Record
		err = attributevalue.UnmarshalListOfMaps(records.Items, &selection)
		if err != nil {
			log.WithError(err).Fatal("couldn't parse records")
		}

		for _, record := range selection {
			fmt.Fprintf(w, units.HumanDuration(time.Now().UTC().Sub(record.IssueDate))+" ago ")
			fmt.Fprintf(w, "%+v\n", record)
		}
	})

	var err error

	if _, ok := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME"); ok {
		log.SetHandler(jsonhandler.Default)
		err = gateway.ListenAndServe("", nil)
	} else {
		log.SetHandler(text.Default)
		err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
	}
	log.WithError(err).Fatal("error listening")
}
