package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/apex/gateway/v2"
	"github.com/apex/log"
	jsonhandler "github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	t, err := template.ParseFS(tmpl, "templates/*.html")
	if err != nil {
		log.WithError(err).Fatal("Failed to parse templates")
	}

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {

		customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:   "aws",
				SigningRegion: "us-west-2",
				URL:           "http://localhost:8000",
			}, nil
		})

		defaultConfig, err := config.LoadDefaultConfig(ctx.TODO(),
			config.WithRegion(aws),
			config.WithEndpointResolver(customResolver),
		)
		if err != nil {
			log.WithError(err).Fatal("couldn't get default config")
		}
		client := dynamodb.NewFromConfig(defaultConfig)

		param := &dynamodb.ListTablesInput{
			Limit: aws.Int32(10),
		}

		tables, err := client.ListTables(context.TODO(), param)
		if err != nil {
			log.WithError(err).Fatal("couldn't list tables")
		}
		for i, s := range tables.TableNames {
			println(i, s)
		}

		rw.Header().Set("Content-Type", "text/html")
		err = t.ExecuteTemplate(rw, "index.html", struct {
			Version string
		}{
			Version,
		})
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			log.WithError(err).Fatal("Failed to execute templates")
		}
	})

	if _, ok := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME"); ok {
		log.SetHandler(jsonhandler.Default)
		err = gateway.ListenAndServe("", nil)
	} else {
		log.SetHandler(text.Default)
		err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
	}
	log.WithError(err).Fatal("error listening")
}
