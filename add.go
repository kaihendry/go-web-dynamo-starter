package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func (s *server) add() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handle post only
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// parse body to a record
		var rec Record

		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// r.PostForm is a map of our POST form values
		err = decoder.Decode(&rec, r.PostForm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		rec.ID = r.RemoteAddr
		rec.Created = time.Now()
		rec.Expires = rec.Created.Add(time.Minute)

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

		http.Redirect(w, r, "/", http.StatusFound)

	}
}
