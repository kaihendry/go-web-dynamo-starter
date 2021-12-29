package main

import (
	"context"
	"embed"
	"net/http"
	"os"
	"text/template"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

//go:embed templates
var tmpl embed.FS

func (s *server) list() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("list")

		t, err := template.ParseFS(tmpl, "templates/*.html")
		if err != nil {
			log.WithError(err).Fatal("Failed to parse templates")
		}

		// https://aws.github.io/aws-sdk-go-v2/docs/code-examples/dynamodb/scanitems/
		records, err := s.client.Scan(context.TODO(), &dynamodb.ScanInput{
			TableName: aws.String(os.Getenv("TABLE_NAME")),
		})
		if err != nil {
			log.WithError(err).Fatal("couldn't get records")
		}

		log.WithField("table", os.Getenv("TABLE_NAME")).Info("got records")

		var selection []Record

		err = attributevalue.UnmarshalListOfMaps(records.Items, &selection)
		if err != nil {
			log.WithError(err).Fatal("couldn't parse records")
		}

		log.WithField("count", len(selection)).Info("parsed records")

		w.Header().Set("Content-Type", "text/html")
		err = t.ExecuteTemplate(w, "index.html", struct {
			Selection []Record
		}{
			selection,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.WithError(err).Fatal("Failed to execute templates")
		}
	}
}
