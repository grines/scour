package awsdata

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/directconnect"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func DescribeDirectConnectGateways(sess *session.Session) *directconnect.DescribeDirectConnectGatewaysOutput {
	svc := directconnect.New(sess)
	input := &directconnect.DescribeDirectConnectGatewaysInput{}

	result, err := svc.DescribeDirectConnectGateways(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case lambda.ErrCodeServiceException:
				fmt.Println(lambda.ErrCodeServiceException, aerr.Error())
			case lambda.ErrCodeResourceNotFoundException:
				fmt.Println(lambda.ErrCodeResourceNotFoundException, aerr.Error())
			case lambda.ErrCodeTooManyRequestsException:
				fmt.Println(lambda.ErrCodeTooManyRequestsException, aerr.Error())
			case lambda.ErrCodeInvalidParameterValueException:
				fmt.Println(lambda.ErrCodeInvalidParameterValueException, aerr.Error())
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

func DescribeDirectConnectGatewayAssociations(sess *session.Session, gatewayID string) *directconnect.DescribeDirectConnectGatewayAssociationsOutput {
	svc := directconnect.New(sess)
	input := &directconnect.DescribeDirectConnectGatewayAssociationsInput{
		AssociationId: aws.String(gatewayID),
	}

	result, err := svc.DescribeDirectConnectGatewayAssociations(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case lambda.ErrCodeServiceException:
				fmt.Println(lambda.ErrCodeServiceException, aerr.Error())
			case lambda.ErrCodeResourceNotFoundException:
				fmt.Println(lambda.ErrCodeResourceNotFoundException, aerr.Error())
			case lambda.ErrCodeTooManyRequestsException:
				fmt.Println(lambda.ErrCodeTooManyRequestsException, aerr.Error())
			case lambda.ErrCodeInvalidParameterValueException:
				fmt.Println(lambda.ErrCodeInvalidParameterValueException, aerr.Error())
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
