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
