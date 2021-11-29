package main

import (
	"context"
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
	Id         string
	IssueDate  time.Time
	RandomWord string
}

func dynamoCloud() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = "ap-southeast-1"
		o.SharedConfigProfile = "mine"
		return nil
	})
	if err != nil {
		panic(err)
	}
	return dynamodb.NewFromConfig(cfg)
}

func dynamoLocal() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.SharedConfigProfile = "mine"
		return nil
	})
	if err != nil {
		panic(err)
	}
	return dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.EndpointResolver = dynamodb.EndpointResolverFromURL("http://localhost:8000")
	})
}

func main() {

	client := dynamoLocal()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// https://aws.github.io/aws-sdk-go-v2/docs/code-examples/dynamodb/scanitems/
		records, err := client.Scan(context.TODO(), &dynamodb.ScanInput{
			TableName: aws.String("Records"),
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
