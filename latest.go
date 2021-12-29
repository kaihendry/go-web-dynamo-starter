package main

import (
	"context"
	"html/template"
	"net/http"
	"os"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func (s *server) latest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// parse id from get request
		id := r.URL.Query().Get("id")

		log.WithField("id", id).Info("getting latest")

		expr, err := expression.NewBuilder().WithKeyCondition(expression.Key("id").Equal(expression.Value(id))).Build()
		if err != nil {
			log.WithError(err).Fatal("couldn't build query expression")
		}

		records, err := s.client.Query(context.TODO(), &dynamodb.QueryInput{
			KeyConditionExpression:    expr.KeyCondition(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			TableName:                 aws.String(os.Getenv("TABLE_NAME")),
			ScanIndexForward:          aws.Bool(false),
			Limit:                     aws.Int32(1),
		})
		if err != nil {
			log.WithError(err).Fatal("couldn't get records")
		}

		var selection []Record
		err = attributevalue.UnmarshalListOfMaps(records.Items, &selection)
		if err != nil {
			log.WithError(err).Fatal("couldn't parse records")
		}

		t, err := template.ParseFS(tmpl, "templates/*.html")
		if err != nil {
			log.WithError(err).Fatal("Failed to parse templates")
		}

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
