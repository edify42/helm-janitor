package main

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/edify42/helm-janitor/cmd/scan"
	log "github.com/sirupsen/logrus"
)

// runs the generic handler to execute helm delete...
// when the ttl expires.

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	logLevel := "info"
	if os.Getenv("LOG_LEVEL") != "" {
		logLevel = os.Getenv("LOG_LEVEL")
	}
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Errorf("Dodgy log level set: %s", logLevel)
		log.SetLevel(log.WarnLevel)
	} else {
		log.SetLevel(level)
	}
}

func HandleRequest() error {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("ap-southeast-2"),
		//config.WithClientLogMode(aws.LogSigning),
	)
	if err != nil {
		log.Fatal(err)
	}
	roleARN := os.Getenv("ROLE_ARN")
	sessionName := "sessionName"
	stsClient := sts.NewFromConfig(cfg)
	// provider := stscreds.NewAssumeRoleProvider(stsClient, roleARN)
	// cfg.Credentials = aws.NewCredentialsCache(provider)
	// // without the following, I'm getting an error message: api error SignatureDoesNotMatch: The request signature we calculated does not match the signature you provided.
	// creds, err := cfg.Credentials.Retrieve(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Infof(creds.AccessKeyID)
	// log.Infof(creds.SecretAccessKey)
	// log.Infof(creds.SessionToken)

	input := &sts.AssumeRoleInput{
		RoleArn:         &roleARN,
		RoleSessionName: &sessionName,
	}

	result, err := stsClient.AssumeRole(ctx, input)
	if err != nil {
		log.Fatalf("Got an error assuming the role: %v", err)
	}

	log.Info(*result.AssumedRoleUser.Arn)

	os.Setenv("AWS_ACCESS_KEY_ID", *result.Credentials.AccessKeyId)
	os.Setenv("AWS_SECRET_ACCESS_KEY", *result.Credentials.SecretAccessKey)
	os.Setenv("AWS_SESSION_TOKEN", *result.Credentials.SessionToken)

	scanner := scan.NewScanClient()
	scanner.Dryrun = true
	scanner.AllNamespaces = true
	scanner.Context = ctx
	scanner.Init()
	log.Info("starting...")
	scan.RunV2(scanner)
	return nil
}

func main() {
	log.Infof("starting")
	if os.Getenv("DEBUG") == "true" {
		HandleRequest()
	}
	// lambda.Start(HandleRequest)
	log.Infof("finished")
}
