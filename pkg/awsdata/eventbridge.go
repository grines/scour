package awsdata

import (
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/drk1wi/Modlishka/log"
)

func PutRule(sess *session.Session, rulename string) string {
	rando := String(10)
	ruleName := rulename + rando
	svc := eventbridge.New(sess)
	input := &eventbridge.PutRuleInput{
		Name:               aws.String(ruleName),
		ScheduleExpression: aws.String("cron(*/5 * * * ? *)"),
	}

	result, err := svc.PutRule(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeLimitExceededException:
				fmt.Println(iam.ErrCodeLimitExceededException, aerr.Error())
			case iam.ErrCodeInvalidInputException:
				fmt.Println(iam.ErrCodeInvalidInputException, aerr.Error())
			case iam.ErrCodeEntityAlreadyExistsException:
				fmt.Println(iam.ErrCodeEntityAlreadyExistsException, aerr.Error())
			case iam.ErrCodeMalformedPolicyDocumentException:
				fmt.Println(iam.ErrCodeMalformedPolicyDocumentException, aerr.Error())
			case iam.ErrCodeConcurrentModificationException:
				fmt.Println(iam.ErrCodeConcurrentModificationException, aerr.Error())
			case iam.ErrCodeServiceFailureException:
				fmt.Println(iam.ErrCodeServiceFailureException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return ""
	}
	log.Infof("* Created EventBridge rule %s successfully", *result.RuleArn)
	return *result.RuleArn
}

func PutTarget(sess *session.Session, targetArn string, roleArn string, ruleArn string) {

	regexRule := `rule\/(.*)`
	var ruleName string

	r, _ := regexp.Compile(regexRule)
	if r.MatchString(ruleArn) {
		matches := r.FindStringSubmatch(string(ruleArn))
		ruleName = matches[1]
	}

	var targets []*eventbridge.Target

	targets = append(targets, &eventbridge.Target{
		Arn:     aws.String(targetArn),
		Id:      aws.String("1"),
		RoleArn: aws.String(roleArn),
	})

	svc := eventbridge.New(sess)
	input := &eventbridge.PutTargetsInput{
		Rule:    aws.String(ruleName),
		Targets: targets,
	}

	_, err := svc.PutTargets(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeLimitExceededException:
				fmt.Println(iam.ErrCodeLimitExceededException, aerr.Error())
			case iam.ErrCodeInvalidInputException:
				fmt.Println(iam.ErrCodeInvalidInputException, aerr.Error())
			case iam.ErrCodeEntityAlreadyExistsException:
				fmt.Println(iam.ErrCodeEntityAlreadyExistsException, aerr.Error())
			case iam.ErrCodeMalformedPolicyDocumentException:
				fmt.Println(iam.ErrCodeMalformedPolicyDocumentException, aerr.Error())
			case iam.ErrCodeConcurrentModificationException:
				fmt.Println(iam.ErrCodeConcurrentModificationException, aerr.Error())
			case iam.ErrCodeServiceFailureException:
				fmt.Println(iam.ErrCodeServiceFailureException, aerr.Error())
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
	log.Infof("* Added EventBridge CodeBuild target %s successfully", targetArn)
}
