package main

import (
	"log"
	"github.com/spf13/viper"

	"go.uber.org/fx"

	"github.com/elmntri/zeitgeber-aws-modules/bucket_connector"
	
	"github.com/elmntri/zeitgeber-common-modules/http_server"
	"github.com/elmntri/zeitgeber-common-modules/logger"

	"github.com/IMingElmntri/aws-module-api-example/aws_apis"
)

func main() {

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	app := fx.New(
		logger.Module(),
		bucket_connector.Module("bucket_connector"),
		http_server.Module("web_service"),
		aws_apis.Module("aws_apis"),
	)

	app.Run()
}