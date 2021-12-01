package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	units "github.com/docker/go-units"
)

func (s *server) list() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// https://aws.github.io/aws-sdk-go-v2/docs/code-examples/dynamodb/scanitems/
		records, err := s.client.Scan(context.TODO(), &dynamodb.ScanInput{
			TableName: aws.String(os.Getenv("TABLE_NAME")),
		})
		if err != nil {
			log.WithError(err).Fatal("couldn't get records")
		}

		var selection []record
		err = attributevalue.UnmarshalListOfMaps(records.Items, &selection)
		if err != nil {
			log.WithError(err).Fatal("couldn't parse records")
		}

		for _, rec := range selection {
			fmt.Fprintf(w, units.HumanDuration(time.Now().UTC().Sub(rec.IssueDate))+" ago ")
			fmt.Fprintf(w, "%+v\n", rec)
		}
	}
}
