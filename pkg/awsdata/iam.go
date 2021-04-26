package awsdata

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/organizations"
)

func ListUsers(sess *session.Session, hide bool) []string {
	data := [][]string{}
	var users []string

	// Create a IAM service client.
	svc := iam.New(sess)

	result, err := svc.ListUsers(&iam.ListUsersInput{
		MaxItems: aws.Int64(10),
	})

	if err != nil {
		fmt.Println("Error", err)
		return nil
	}

	for _, user := range result.Users {
		if user == nil {
			continue
		}
		//fmt.Printf("%d user %s created %v\n", i, *user.UserName, user.CreateDate)
		users = append(users, *user.UserName)
		row := []string{*user.UserName, user.CreateDate.String()}
		data = append(data, row)
	}
	header := []string{"UserName", "CreateDate"}
	if hide == false {
		tableData(data, header)
	}
	return users
}

func ListGroups(sess *session.Session, hide bool) *iam.ListGroupsOutput {
	data := [][]string{}
	var groups []string

	// Create a IAM service client.
	svc := iam.New(sess)

	result, err := svc.ListGroups(&iam.ListGroupsInput{
		MaxItems: aws.Int64(100),
	})

	if err != nil {
		fmt.Println("Error", err)
		return nil
	}

	for _, group := range result.Groups {
		if group == nil {
			continue
		}
		groups = append(groups, *group.GroupName)
		row := []string{*group.GroupName, group.CreateDate.String()}
		data = append(data, row)
	}
	header := []string{"GroupName", "CreateDate"}
	if hide == false {
		tableData(data, header)
	}
	return result
}

func ListRoles(sess *session.Session, hide bool) *iam.ListRolesOutput {
	data := [][]string{}
	var roles []string
	// Create a IAM service client.
	svc := iam.New(sess)

	result, err := svc.ListRoles(&iam.ListRolesInput{
		MaxItems: aws.Int64(100),
	})

	if err != nil {
		fmt.Println("Error", err)
		return nil
	}

	if hide == false {
		for _, role := range result.Roles {
			if role == nil {
				continue
			}
			roles = append(roles, *role.RoleName)
			row := []string{*role.RoleName, *role.Arn, role.CreateDate.String()}
			data = append(data, row)
		}
	}
	header := []string{"RoleName", "RoleARN", "CreateDate"}
	if hide == false {
		tableData(data, header)
	}
	return result
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

	for _, group := range result.Groups {
		if group == nil {
			continue
		}
		//fmt.Printf("%d role %s created %v\n", i, *group.Arn, group.CreateDate)
	}
	return result.Groups
}

func ListPolicies(sess *session.Session) {
	data := [][]string{}
	var policies []string

	// Create a IAM service client.
	svc := iam.New(sess)

	result, err := svc.ListPolicies(&iam.ListPoliciesInput{
		MaxItems:     aws.Int64(1000),
		OnlyAttached: aws.Bool(true),
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	for _, policy := range result.Policies {
		if policy == nil {
			continue
		}
		policies = append(policies, *policy.Arn)
		n := *policy.AttachmentCount
		count := strconv.FormatInt(n, 10)
		row := []string{*policy.PolicyName, *policy.Arn, count}
		data = append(data, row)
	}
	header := []string{"Policy Name", "Arn", "Attachment Count"}
	tableData(data, header)
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

func ListAttachedRolePolicies(sess *session.Session, rolename string) []*iam.AttachedPolicy {

	// Create a IAM service client.
	svc := iam.New(sess)

	result, err := svc.ListAttachedRolePolicies(&iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(rolename),
	})

	if err != nil {
		fmt.Println("Error", err)
		return nil
	}

	//for i, role := range result.AttachedPolicies {
	//	if role == nil {
	//		continue
	//	}
	//	fmt.Printf("%d policy %s name %v\n", i, *role.PolicyArn, role.PolicyName)
	//}
	return result.AttachedPolicies
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

func ListGroupPolicies(sess *session.Session, group string) *iam.ListGroupPoliciesOutput {

	// Create a IAM service client.
	svc := iam.New(sess)

	result, err := svc.ListGroupPolicies(&iam.ListGroupPoliciesInput{
		GroupName: aws.String(group),
	})

	if err != nil {
		fmt.Println("Error", err)
		return nil
	}

	//for i, policy := range result.PolicyNames {
	//	if policy == nil {
	//		continue
	//	}
	//	fmt.Printf("%d policy %v\n", i, *policy)
	//}
	return result
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

func GetPolicyVersion(sess *session.Session, arn string) string {
	// Create a IAM service client.
	svc := iam.New(sess)

	result, err := svc.GetPolicyVersion(&iam.GetPolicyVersionInput{
		PolicyArn: &arn,
		VersionId: aws.String("v1"),
	})

	if err != nil {
		fmt.Println("Error", err)
		return ""
	}

	//fmt.Printf("%s\n", *result.PolicyVersion.Document)
	decodedValue, err := url.QueryUnescape(aws.StringValue(result.PolicyVersion.Document))
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return decodedValue
}

func GetUserPolicy(sess *session.Session, username string, policy string) string {
	// Create a IAM service client.
	svc := iam.New(sess)

	result, err := svc.GetUserPolicy(&iam.GetUserPolicyInput{
		UserName:   aws.String(username),
		PolicyName: aws.String(policy),
	})

	if err != nil {
		fmt.Println("Error", err)
		return ""
	}

	//fmt.Printf("%s\n", *result.PolicyVersion.Document)
	decodedValue, err := url.QueryUnescape(aws.StringValue(result.PolicyDocument))
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return decodedValue

	//data := []byte(decodedValue)
}

func GetGroupPolicy(sess *session.Session, group string, policy string) string {
	// Create a IAM service client.
	svc := iam.New(sess)

	result, err := svc.GetGroupPolicy(&iam.GetGroupPolicyInput{
		GroupName:  aws.String(group),
		PolicyName: aws.String(policy),
	})

	if err != nil {
		fmt.Println("Error", err)
		return ""
	}

	//fmt.Printf("%s\n", *result.PolicyVersion.Document)
	decodedValue, err := url.QueryUnescape(aws.StringValue(result.PolicyDocument))
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return decodedValue

	//data := []byte(decodedValue)
}

func CreateAccessKey(sess *session.Session, username string) *iam.CreateAccessKeyOutput {
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
		return nil
	}

	return result
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

func GetRole(sess *session.Session, rolename string) *iam.GetRoleOutput {
	svc := iam.New(sess)
	input := &iam.GetRoleInput{
		RoleName: aws.String(rolename),
	}

	result, err := svc.GetRole(input)
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
		return nil
	}

	return result
}

func CreateRole(sess *session.Session, name string, policy string) string {

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
		return ""
	}

	return *result.Role.Arn
}

func AttachRolePolicy(sess *session.Session, name string, policyArn string) {
	svc := iam.New(sess)
	input := &iam.AttachRolePolicyInput{
		PolicyArn: aws.String(policyArn),
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

func GetInstanceProfile(sess *session.Session, profile string, hide bool) *iam.GetInstanceProfileOutput {
	svc := iam.New(sess)
	input := &iam.GetInstanceProfileInput{
		InstanceProfileName: aws.String(profile),
	}

	result, err := svc.GetInstanceProfile(input)
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
		return nil
	}
	if hide == false {
		fmt.Println(*result.InstanceProfile)
	}
	return result
}

func CreateUser(sess *session.Session, username string) bool {
	svc := iam.New(sess)
	input := &iam.CreateUserInput{
		UserName: aws.String(username),
	}

	_, err := svc.CreateUser(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeLimitExceededException:
				fmt.Println(iam.ErrCodeLimitExceededException, aerr.Error())
			case iam.ErrCodeEntityAlreadyExistsException:
				fmt.Println(iam.ErrCodeEntityAlreadyExistsException, aerr.Error())
			case iam.ErrCodeNoSuchEntityException:
				fmt.Println(iam.ErrCodeNoSuchEntityException, aerr.Error())
			case iam.ErrCodeInvalidInputException:
				fmt.Println(iam.ErrCodeInvalidInputException, aerr.Error())
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
		return false
	}

	return true
}

func AttachUserPolicy(sess *session.Session, username string, arn string) bool {
	svc := iam.New(sess)
	input := &iam.AttachUserPolicyInput{
		PolicyArn: aws.String(arn),
		UserName:  aws.String(username),
	}

	_, err := svc.AttachUserPolicy(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				fmt.Println(iam.ErrCodeNoSuchEntityException, aerr.Error())
			case iam.ErrCodeLimitExceededException:
				fmt.Println(iam.ErrCodeLimitExceededException, aerr.Error())
			case iam.ErrCodeInvalidInputException:
				fmt.Println(iam.ErrCodeInvalidInputException, aerr.Error())
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
		return false
	}

	return true
}

func TrustPolicyAWS(accountID string) string {
	var policy string

	policy = fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
			  {
				"Effect": "Allow",
				"Principal": {
				  "AWS": "arn:aws:iam::%s:root"
				},
				"Action": "sts:AssumeRole",
				"Condition": {}
			  }
			]
		  }`, accountID)
	return policy
}

func TrustPolicyService(service string) string {
	var policy string

	policy = fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
		  {
			"Sid": "",
			"Effect": "Allow",
			"Principal": {
			  "Service": "%s"
			},
			"Action": "sts:AssumeRole"
		  }
		]
	}`, service)

	return policy
}

func CreatePolicy(sess *session.Session, name string, policy string) string {
	svc := iam.New(sess)
	input := &iam.CreatePolicyInput{
		PolicyName:     aws.String(name),
		PolicyDocument: aws.String(policy),
	}

	result, err := svc.CreatePolicy(input)
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
		return ""
	}
	return *result.Policy.Arn
}

func UpdateAssumeRolePolicy(sess *session.Session, role string, policy string) {

	svc := iam.New(sess)
	input := &iam.UpdateAssumeRolePolicyInput{
		PolicyDocument: aws.String(policy),
		RoleName:       aws.String(role),
	}

	result, err := svc.UpdateAssumeRolePolicy(input)
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
	if result.String() != "" {
		fmt.Println("Role updated")
	}
}

func AddUserToGroup(sess *session.Session, group string, user string) bool {
	svc := iam.New(session.New())
	input := &iam.AddUserToGroupInput{
		GroupName: aws.String(group),
		UserName:  aws.String(user),
	}

	_, err := svc.AddUserToGroup(input)
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
		return false
	}

	return true
}

func CreateLoginProfile(sess *session.Session, user string) (*iam.CreateLoginProfileOutput, string) {
	svc := iam.New(sess)
	password := "h]6EszR}vJ*m"
	input := &iam.CreateLoginProfileInput{
		Password:              aws.String(password),
		PasswordResetRequired: aws.Bool(false),
		UserName:              aws.String(user),
	}

	result, err := svc.CreateLoginProfile(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeEntityAlreadyExistsException:
				fmt.Println(iam.ErrCodeEntityAlreadyExistsException, aerr.Error())
			case iam.ErrCodeNoSuchEntityException:
				fmt.Println(iam.ErrCodeNoSuchEntityException, aerr.Error())
			case iam.ErrCodePasswordPolicyViolationException:
				fmt.Println(iam.ErrCodePasswordPolicyViolationException, aerr.Error())
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
		return nil, ""
	}

	return result, password
}

func UpdateLoginProfile(sess *session.Session, user string) (*iam.UpdateLoginProfileOutput, string) {
	svc := iam.New(sess)
	password := "h]6EszR}vJ*m"
	input := &iam.UpdateLoginProfileInput{
		Password:              aws.String(password),
		PasswordResetRequired: aws.Bool(false),
		UserName:              aws.String(user),
	}

	result, err := svc.UpdateLoginProfile(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeEntityAlreadyExistsException:
				fmt.Println(iam.ErrCodeEntityAlreadyExistsException, aerr.Error())
			case iam.ErrCodeNoSuchEntityException:
				fmt.Println(iam.ErrCodeNoSuchEntityException, aerr.Error())
			case iam.ErrCodePasswordPolicyViolationException:
				fmt.Println(iam.ErrCodePasswordPolicyViolationException, aerr.Error())
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
		return nil, ""
	}

	return result, password
}

func ListAccountAliases(sess *session.Session) {
	svc := iam.New(sess)
	input := &iam.ListAccountAliasesInput{}

	result, err := svc.ListAccountAliases(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
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

func DescribeOrganization(sess *session.Session) {
	svc := organizations.New(sess)
	input := &organizations.DescribeOrganizationInput{}

	result, err := svc.DescribeOrganization(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case organizations.ErrCodeAccessDeniedException:
				fmt.Println(organizations.ErrCodeAccessDeniedException, aerr.Error())
			case organizations.ErrCodeAWSOrganizationsNotInUseException:
				fmt.Println(organizations.ErrCodeAWSOrganizationsNotInUseException, aerr.Error())
			case organizations.ErrCodeConcurrentModificationException:
				fmt.Println(organizations.ErrCodeConcurrentModificationException, aerr.Error())
			case organizations.ErrCodeServiceException:
				fmt.Println(organizations.ErrCodeServiceException, aerr.Error())
			case organizations.ErrCodeTooManyRequestsException:
				fmt.Println(organizations.ErrCodeTooManyRequestsException, aerr.Error())
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

func AssociateIamInstanceProfile(sess *session.Session, instance string, role string) {
	svc := ec2.New(sess)
	input := &ec2.AssociateIamInstanceProfileInput{
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Name: aws.String(role),
		},
		InstanceId: aws.String(instance),
	}

	result, err := svc.AssociateIamInstanceProfile(input)
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

	fmt.Println(result)
}

func ListInstanceProfiles(sess *session.Session) {
	svc := iam.New(sess)
	input := &iam.ListInstanceProfilesInput{}

	result, err := svc.ListInstanceProfiles(input)
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

	fmt.Println(result)
}
