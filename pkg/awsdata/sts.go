package awsdata

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

func GetCallerIdentity(sess *session.Session) string {
	svc := sts.New(sess)

	var params *sts.GetCallerIdentityInput
	resp, err := svc.GetCallerIdentity(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return ""
	}

	data := [][]string{
		{*resp.Account, *resp.Arn, *resp.UserId},
	}

	header := []string{"Account", "ARN", "UserID"}
	tableData(data, header)

	return string(*resp.Arn)
}

func GetSessionToken(sess *session.Session) *sts.GetSessionTokenOutput {
	svc := sts.New(sess)
	input := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(3600),
	}

	result, err := svc.GetSessionToken(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case sts.ErrCodeRegionDisabledException:
				fmt.Println(sts.ErrCodeRegionDisabledException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}

	return result
}

func GetFederationToken(sess *session.Session, username string) *sts.GetFederationTokenOutput {
	svc := sts.New(sess)
	input := &sts.GetFederationTokenInput{
		DurationSeconds: aws.Int64(3600),
		Name:            aws.String(username),
		Policy:          aws.String("{\"Version\":\"2012-10-17\",\"Statement\":[{\"Sid\":\"Stmt1\",\"Effect\":\"Allow\",\"Action\":\"*\",\"Resource\":\"*\"}]}"),
	}

	result, err := svc.GetFederationToken(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case sts.ErrCodeMalformedPolicyDocumentException:
				fmt.Println(sts.ErrCodeMalformedPolicyDocumentException, aerr.Error())
			case sts.ErrCodePackedPolicyTooLargeException:
				fmt.Println(sts.ErrCodePackedPolicyTooLargeException, aerr.Error())
			case sts.ErrCodeRegionDisabledException:
				fmt.Println(sts.ErrCodeRegionDisabledException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}

	return result
}
