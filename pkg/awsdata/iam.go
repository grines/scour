package awsdata

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

func ListUsers(sess *session.Session) {

	// Create a IAM service client.
	svc := iam.New(sess)

	result, err := svc.ListUsers(&iam.ListUsersInput{
		MaxItems: aws.Int64(10),
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	for i, user := range result.Users {
		if user == nil {
			continue
		}
		fmt.Printf("%d user %s created %v\n", i, *user.UserName, user.CreateDate)
	}
}

func ListGroups(sess *session.Session) {

	// Create a IAM service client.
	svc := iam.New(sess)

	result, err := svc.ListGroups(&iam.ListGroupsInput{
		MaxItems: aws.Int64(10),
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	for i, user := range result.Groups {
		if user == nil {
			continue
		}
		fmt.Printf("%d user %s created %v\n", i, *user.GroupName, user.CreateDate)
	}
}

func ListRoles(sess *session.Session) {

	// Create a IAM service client.
	svc := iam.New(sess)

	result, err := svc.ListRoles(&iam.ListRolesInput{
		MaxItems: aws.Int64(10),
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	for i, role := range result.Roles {
		if role == nil {
			continue
		}
		fmt.Printf("%d role %s created %v\n", i, *role.Arn, role.CreateDate)
	}
}

func ListGroupsForUser(sess *session.Session, username string) []*iam.Group {

	// Create a IAM service client.
	svc := iam.New(sess)

	result, err := svc.ListGroupsForUser(&iam.ListGroupsForUserInput{
		UserName: aws.String(username),
	})

	if err != nil {
		fmt.Println("Error", err)
		return nil
	}

	for i, group := range result.Groups {
		if group == nil {
			continue
		}
		fmt.Printf("%d role %s created %v\n", i, *group.Arn, group.CreateDate)
	}
	return result.Groups
}

func ListAttachedGroupPolicies(sess *session.Session, groupname string) []*iam.AttachedPolicy {

	// Create a IAM service client.
	svc := iam.New(sess)

	result, err := svc.ListAttachedGroupPolicies(&iam.ListAttachedGroupPoliciesInput{
		GroupName: aws.String(groupname),
	})

	if err != nil {
		fmt.Println("Error", err)
		return nil
	}

	return result.AttachedPolicies
}

func ListAttachedRolePolicies(sess *session.Session, rolename string) {

	// Create a IAM service client.
	svc := iam.New(sess)

	result, err := svc.ListAttachedRolePolicies(&iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(rolename),
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	for i, role := range result.AttachedPolicies {
		if role == nil {
			continue
		}
		fmt.Printf("%d policy %s name %v\n", i, *role.PolicyArn, role.PolicyName)
	}
}

func ListUserPolicies(sess *session.Session, username string) {

	// Create a IAM service client.
	svc := iam.New(sess)

	result, err := svc.ListUserPolicies(&iam.ListUserPoliciesInput{
		UserName: aws.String(username),
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	for i, user := range result.PolicyNames {
		if user == nil {
			continue
		}
		fmt.Printf("%d policy %v\n", i, *user)
	}
}

func ListGroupPolicies(sess *session.Session, group string) {

	// Create a IAM service client.
	svc := iam.New(sess)

	result, err := svc.ListGroupPolicies(&iam.ListGroupPoliciesInput{
		GroupName: aws.String(group),
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	for i, policy := range result.PolicyNames {
		if policy == nil {
			continue
		}
		fmt.Printf("%d policy %v\n", i, *policy)
	}
}

func GetUser(sess *session.Session, username string) {
	svc := iam.New(sess)
	input := &iam.GetUserInput{
		UserName: aws.String(username),
	}

	result, err := svc.GetUser(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				fmt.Println(iam.ErrCodeNoSuchEntityException, aerr.Error())
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

	fmt.Println(result)
}

func GetAccountAuthorizationDetails(sess *session.Session) []*iam.UserDetail {
	svc := iam.New(sess)

	user := "User"
	input := &iam.GetAccountAuthorizationDetailsInput{Filter: []*string{&user}}
	resp, err := svc.GetAccountAuthorizationDetails(input)
	if err != nil {
		fmt.Println("Got error getting account details")
		fmt.Println(err.Error())
	}

	return resp.UserDetailList
}

func GetPolicy(sess *session.Session, arn string) {
	// Create a IAM service client.
	svc := iam.New(sess)

	result, err := svc.GetPolicy(&iam.GetPolicyInput{
		PolicyArn: &arn,
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	fmt.Printf("%s - %s - %s\n", arn, *result.Policy.Description, *result.Policy)
}

func CreateAccessKey(sess *session.Session, username string) {
	svc := iam.New(sess)
	input := &iam.CreateAccessKeyInput{
		UserName: aws.String(username),
	}

	result, err := svc.CreateAccessKey(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				fmt.Println(iam.ErrCodeNoSuchEntityException, aerr.Error())
			case iam.ErrCodeLimitExceededException:
				fmt.Println(iam.ErrCodeLimitExceededException, aerr.Error())
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

	fmt.Println(result)
}

func CreateInstanceProfile(sess *session.Session, name string) {
	svc := iam.New(sess)
	input := &iam.CreateInstanceProfileInput{
		InstanceProfileName: aws.String(name),
	}

	result, err := svc.CreateInstanceProfile(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeEntityAlreadyExistsException:
				fmt.Println(iam.ErrCodeEntityAlreadyExistsException, aerr.Error())
			case iam.ErrCodeInvalidInputException:
				fmt.Println(iam.ErrCodeInvalidInputException, aerr.Error())
			case iam.ErrCodeLimitExceededException:
				fmt.Println(iam.ErrCodeLimitExceededException, aerr.Error())
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

	fmt.Println(result)
}

func CreateRole(sess *session.Session, name string) {

	policy := `{
		"Version": "2012-10-17",
		"Statement": [
		  {
			"Sid": "",
			"Effect": "Allow",
			"Principal": {
			  "Service": "ec2.amazonaws.com"
			},
			"Action": "sts:AssumeRole"
		  }
		]
	}`

	svc := iam.New(sess)
	input := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(policy),
		Path:                     aws.String("/"),
		RoleName:                 aws.String(name),
	}

	result, err := svc.CreateRole(input)
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

	fmt.Println(result)
}

func AttachRolePolicy(sess *session.Session, name string) {
	svc := iam.New(sess)
	input := &iam.AttachRolePolicyInput{
		PolicyArn: aws.String("arn:aws:iam::aws:policy/AdministratorAccess"),
		RoleName:  aws.String(name),
	}

	result, err := svc.AttachRolePolicy(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				fmt.Println(iam.ErrCodeNoSuchEntityException, aerr.Error())
			case iam.ErrCodeLimitExceededException:
				fmt.Println(iam.ErrCodeLimitExceededException, aerr.Error())
			case iam.ErrCodeInvalidInputException:
				fmt.Println(iam.ErrCodeInvalidInputException, aerr.Error())
			case iam.ErrCodeUnmodifiableEntityException:
				fmt.Println(iam.ErrCodeUnmodifiableEntityException, aerr.Error())
			case iam.ErrCodePolicyNotAttachableException:
				fmt.Println(iam.ErrCodePolicyNotAttachableException, aerr.Error())
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

	if len(result.String()) > 4 {
		fmt.Println(result.String())
	}
}

func AddRoleToInstanceProfile(sess *session.Session, name string) {
	svc := iam.New(sess)
	input := &iam.AddRoleToInstanceProfileInput{
		InstanceProfileName: aws.String(name),
		RoleName:            aws.String(name),
	}

	result, err := svc.AddRoleToInstanceProfile(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				fmt.Println(iam.ErrCodeNoSuchEntityException, aerr.Error())
			case iam.ErrCodeEntityAlreadyExistsException:
				fmt.Println(iam.ErrCodeEntityAlreadyExistsException, aerr.Error())
			case iam.ErrCodeLimitExceededException:
				fmt.Println(iam.ErrCodeLimitExceededException, aerr.Error())
			case iam.ErrCodeUnmodifiableEntityException:
				fmt.Println(iam.ErrCodeUnmodifiableEntityException, aerr.Error())
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

	if len(result.String()) > 4 {
		fmt.Println(result.String())
	}
}
