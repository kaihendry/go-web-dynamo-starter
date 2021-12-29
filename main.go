package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/apex/gateway/v2"
	"github.com/apex/log"
	jsonhandler "github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Record struct {
	ID      string    `dynamodbav:"id" json:"id"`
	Created time.Time `dynamodbav:"created,unixtime" json:"created"`
	Expires time.Time `dynamodbav:"expires,unixtime" json:"expires"`
	Color   string    `dynamodbav:"color" json:"color"`
}

type server struct {
	router *http.ServeMux
	client *dynamodb.Client
}

func (record *Record) TimeSinceCreation() string {
	return time.Since(record.Created).String()
}

func (record *Record) TimeUntilExpiry() string {
	return time.Until(record.Expires).String()
}

func newServer(local bool) *server {
	s := &server{router: &http.ServeMux{}}

	if local {
		log.SetHandler(text.Default)
		log.Info("local mode")
		s.client = dynamoLocal()
	} else {
		log.SetHandler(jsonhandler.Default)
		log.Info("cloud mode")
		s.client = dynamoCloud()
	}

	s.router.Handle("/", s.list())
	s.router.Handle("/latest", s.latest())
	s.router.Handle("/add", s.add())

	return s
}

func main() {
	_, awsDetected := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME")
	log.WithField("awsDetected", awsDetected).Info("starting up")
	s := newServer(!awsDetected)

	var err error

	if awsDetected {
		log.Info("starting cloud server")
		err = gateway.ListenAndServe("", s.router)
	} else {
		err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), s.router)
	}
	log.Info(".... starting")
	log.WithError(err).Fatal("error listening")
	log.Info(".... ending")
}
