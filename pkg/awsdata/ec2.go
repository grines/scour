package awsdata

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/drk1wi/Modlishka/log"
)

func DescribeInstances(sess *session.Session) []*ec2.Reservation {
	svc := ec2.New(sess)
	input := &ec2.DescribeInstancesInput{}

	result, err := svc.DescribeInstances(input)
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
	return result.Reservations
}

func DescribeInstanceAttribute(sess *session.Session, InstanceId string, attribute string) *ec2.DescribeInstanceAttributeOutput {
	svc := ec2.New(sess)
	input := &ec2.DescribeInstanceAttributeInput{
		Attribute:  aws.String(attribute),
		InstanceId: aws.String(InstanceId),
	}

	result, err := svc.DescribeInstanceAttribute(input)
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

func RunInstance(sess *session.Session, data *ec2.RunInstancesInput) string {

	// Create EC2 service client
	svc := ec2.New(sess)

	// Specify the details of the instance that you want to create.
	runResult, err := svc.RunInstances(data)

	if err != nil {
		log.Errorf("Could not create instance %v", err)
		return ""
	}

	log.Infof("Created instance %v", *runResult.Instances[0].InstanceId)
	return *runResult.Instances[0].InstanceId
}

func ec2Status(sess *session.Session, instanceID string) []*ec2.InstanceStatus {
	svc := ec2.New(sess)
	input := &ec2.DescribeInstanceStatusInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	}

	result, err := svc.DescribeInstanceStatus(input)
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

	fmt.Println("Checking if EC2 is Running...")
	return result.InstanceStatuses
}

func ModifyInstanceAttribute(sess *session.Session, input *ec2.ModifyInstanceAttributeInput) bool {

	svc := ec2.New(sess)

	_, err := svc.ModifyInstanceAttribute(input)
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
		return false
	}
	return true

}

func StopInstance(sess *session.Session, instanceID string) {
	svc := ec2.New(sess)
	input := &ec2.StopInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	}

	result, err := svc.StopInstances(input)
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
		return
	}

	log.Infof("Stopping Instance %v - State: %v", instanceID, *result.StoppingInstances[0].CurrentState.Name)
}

func StartInstance(sess *session.Session, instanceID string) {
	svc := ec2.New(sess)
	input := &ec2.StartInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	}

	result, err := svc.StartInstances(input)
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
		return
	}

	log.Infof("Starting Instance %v - State: %v", instanceID, *result.StartingInstances[0].CurrentState.Name)
}

func DescribeSecurityGroup(sess *session.Session, sg string) *ec2.DescribeSecurityGroupsOutput {
	svc := ec2.New(sess)
	input := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{
			aws.String(sg),
		},
	}

	result, err := svc.DescribeSecurityGroups(input)
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

func DescribeVpnConnections(sess *session.Session) *ec2.DescribeVpnConnectionsOutput {
	svc := ec2.New(sess)
	input := &ec2.DescribeVpnConnectionsInput{}

	result, err := svc.DescribeVpnConnections(input)
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

func DescribeVpcPeeringConnections(sess *session.Session) *ec2.DescribeVpcPeeringConnectionsOutput {
	svc := ec2.New(sess)
	input := &ec2.DescribeVpcPeeringConnectionsInput{}

	result, err := svc.DescribeVpcPeeringConnections(input)
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
