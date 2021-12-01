package app

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

type record struct {
	ID         string    `dynamodbav:"id" json:"id"`
	IssueDate  time.Time `dynamodbav:"issueDate" json:"issueDate"`
	RandomWord string    `dynamodbav:"randomWord" json:"randomWord"`
}

type server struct {
	router *http.ServeMux
	client *dynamodb.Client
}

func newServer(local bool) *server {
	s := &server{router: &http.ServeMux{}}

	if local {
		log.Info("local mode")
		s.client = dynamoLocal()
	} else {
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
		log.SetHandler(jsonhandler.Default)
		err = gateway.ListenAndServe("", s.router)
	} else {
		log.SetHandler(text.Default)
		err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), s.router)
	}
	log.WithError(err).Fatal("error listening")
}
