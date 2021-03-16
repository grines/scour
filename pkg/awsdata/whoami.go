package awsdata

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

func Whoami(sess *session.Session) string {
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
