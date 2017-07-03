package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"log"
	"github.com/aws/aws-sdk-go/service/sts"
	"fmt"
	"os"
	"flag"
	"time"
	//	"github.com/aws/aws-sdk-go/aws/credentials"
	//	"github.com/aws/aws-sdk-go/aws"
	//	"github.com/aws/aws-sdk-go/aws/credentials"
)

const (
	appName         = "assume"
	awsAccessKey    = "AWS_ACCESS_KEY_ID"
	awsAccessSecret = "AWS_SECRET_ACCESS_KEY"
	awsSessionToken = "AWS_SESSION_TOKEN"
)

var (
	verbosity       bool
	onEc2           bool
	duration        int64
	externalId      string
	roleArn         string
	roleSessionName string
	profile         string
)

func main() {

	flag.Int64Var(&duration, "d", 3600, "Credential duration")
	flag.StringVar(&externalId, "i", "", "External ID")
	flag.StringVar(&roleArn, "r", "", "Role ARN")
	flag.StringVar(&roleSessionName, "s", "", "Role session name")
	flag.StringVar(&profile, "p", "", "AWS profile to try if not on EC2")
	flag.BoolVar(&verbosity, "v", false, "Verbose output")

	flag.Parse()

	if roleArn == "" {
		errorExit("You need to specify a role ARN")
	}

	if externalId == "" {
		externalId = defaultExternalId()
	}

	if roleSessionName == "" {
		roleSessionName = defaultSessionName()
	}

	s := session.New()

	verbose("Checking for EC2...")
	svc := ec2metadata.New(s)
	onEc2 = svc.Available()
	verbose("EC2 found: %t", onEc2)

	if !onEc2 && profile == "" {
		verbose("Not on EC2, no profile specified, giving up")

	}

	//Try and assume role

	verbose("Assuming role...")
	verbose("Duration: %d", duration)
	verbose("ExternalId: %s", externalId)
	verbose("Role ARN: %s", roleArn)
	verbose("Role Session Name: %s", roleSessionName)

	input := &sts.AssumeRoleInput{
		DurationSeconds: &duration,
		ExternalId:      &externalId,
		RoleArn:         &roleArn,
		RoleSessionName: &roleSessionName,
	}

	var result *sts.AssumeRoleOutput
	var err error

	if onEc2 {

		svc := sts.New(s)

		result, err = svc.AssumeRole(input)
		if err != nil {
			errorExit(err)
		}

	} else {

		verbose("Not on EC2, trying profile: %s", profile)

		s, err := session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
			Profile: profile,
		})

		if err != nil {
			errorExit(err)
		}

		svc := sts.New(s)

		result, err = svc.AssumeRole(input)

		if err != nil {
			errorExit(err)
		}
	}

	verbose("Role assumed: %s", *result.AssumedRoleUser.Arn)
	verbose("Temporary credentials follow...\n")

	envOutput(awsAccessKey, *result.Credentials.AccessKeyId)
	envOutput(awsAccessSecret, *result.Credentials.SecretAccessKey)
	envOutput(awsSessionToken, *result.Credentials.SessionToken)

}



func envOutput(key string, value string) {
	fmt.Printf("%s=%s\n", key, value)
}

func defaultExternalId() string {
	return appName
}

func defaultSessionName() string {
	return fmt.Sprintf("%s-%d", appName, time.Now().UnixNano())
}

func errorExit(a ...interface{}) {

	fmt.Println(a...)
	os.Exit(1)
}

func verbose(format string, v ...interface{}) {

	if verbosity {
		log.Printf(format+"\n", v...)
	}
}
