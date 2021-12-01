package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func (s *server) add() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Handle post only
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// parse body to a record
		var rec record
		err := json.NewDecoder(r.Body).Decode(&rec)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.WithField("record", rec).Info("adding")

		// marshall record into dynamo
		av, err := attributevalue.MarshalMap(rec)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// what if issueDate is wrong?
		_, err = s.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(os.Getenv("TABLE_NAME")),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// output how many records we have, using an empty ProjectionExpression
		records, err := s.client.Scan(context.TODO(), &dynamodb.ScanInput{
			TableName:            aws.String(os.Getenv("TABLE_NAME")),
			ProjectionExpression: aws.String("id"),
		})

		fmt.Fprintf(w, "added %d records\n", records.Count)

	}
}
