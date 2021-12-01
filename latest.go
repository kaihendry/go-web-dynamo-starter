package app

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

func (s *server) latest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("getting latest")

		// sort by issueDate
		records, err := s.client.Query(context.TODO(), &dynamodb.QueryInput{
			TableName:        aws.String(os.Getenv("TABLE_NAME")),
			ScanIndexForward: aws.Bool(true),
			Limit:            aws.Int32(1),
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
