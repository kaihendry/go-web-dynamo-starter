package main

import (
	"context"

	"github.com/apex/log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func dynamoCloud() *dynamodb.Client {
	log.Info("cloud config")
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		return nil
	})
	if err != nil {
		log.WithError(err).Fatal("failed to load config")
	}
	return dynamodb.NewFromConfig(cfg)
}

func dynamoLocal() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.SharedConfigProfile = "mine"
		return nil
	})
	if err != nil {
		log.WithError(err).Fatal("failed to load config")
	}
	return dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.EndpointResolver = dynamodb.EndpointResolverFromURL("http://localhost:8000")
	})
}
