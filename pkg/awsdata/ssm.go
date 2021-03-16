package awsdata

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func DescribeParameters(sess *session.Session) []*ssm.ParameterMetadata {
	svc := ssm.New(sess)
	input := &ssm.DescribeParametersInput{}

	result, err := svc.DescribeParameters(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
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

	return result.Parameters
}

func GetParameter(sess *session.Session, paramname string) *ssm.GetParameterOutput {
	svc := ssm.New(sess)
	input := &ssm.GetParameterInput{
		Name: aws.String(paramname),
	}

	result, err := svc.GetParameter(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
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

//SendCommand to ssm enabled ec2
func SendCommand(sess *session.Session, instanceID string, cmd string) string {

	var instanceList, commandList []string

	cmd = strings.Replace(cmd, "\"", "", 2)
	//commandList = strings.Split(cmd, " ")
	//fmt.Println(commandList)
	commandList = append(commandList, cmd)

	svc := ssm.New(sess)
	instanceList = append(instanceList, instanceID)

	input := &ssm.SendCommandInput{
		InstanceIds:  aws.StringSlice(instanceList),
		DocumentName: aws.String("AWS-RunShellScript"),
		Parameters: map[string][]*string{
			/*
				ssm.SendCommandInput objects require parameters for the DocumentName chosen
				For AWS-RunShellScript, the only required parameter is "commands",
				which is the shell command to be executed on the target. To emulate
				the original script, we also set "executionTimeout" to 10 minutes.
			*/
			"commands":         aws.StringSlice(commandList),
			"executionTimeout": aws.StringSlice([]string{"600"}),
		},
	}

	// Example sending a request using the SendCommandRequest method.
	req, resp := svc.SendCommandRequest(input)

	err := req.Send()
	if err == nil { // resp is now filled
		return *resp.Command.CommandId
	}
	fmt.Println(err)
	return ""
}

func GetCommandInvocation(sess *session.Session, instanceID string, cmdID string) {

	svc := ssm.New(sess)

	input := &ssm.GetCommandInvocationInput{
		CommandId:  aws.String(cmdID),
		InstanceId: aws.String(instanceID),
	}

	req, resp := svc.GetCommandInvocationRequest(input)

	err := req.Send()
	if err != nil { // resp is now filled
		fmt.Println(err)
	}
	fmt.Println("Command Output: ")
	fmt.Println(*resp.StandardOutputContent)

}
