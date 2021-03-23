package completion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/grines/scour/pkg/awsdata"
)

func Commands(line string, t string) {
	switch {

	//Load aws profile from .aws/credentials
	case strings.HasPrefix(line, "profile "):
		help := HelpText("profile ec2user us-east-1", "Profile is used to load a profile from ~/.aws/credentials.", "enabled")
		parse := ParseCMD(line, 3, help)
		if parse != nil {
			target = parse[1]
			region = parse[2]
			sess = getProfile(target, region)
		}

	//GetSessionToken for current user
	case strings.HasPrefix(line, "get-session-token"):
		help := HelpText("get-session-token us-east-1", "GetSessionToken for user", "disabled")
		parse := ParseCMD(line, 2, help)
		if parse != nil {
			region = parse[1]
			token := awsdata.GetSessionToken(sess)
			if token != nil {
				sess = stsSession(*token.Credentials.AccessKeyId, *token.Credentials.SecretAccessKey, *token.Credentials.SessionToken, region)
			}
		}

	//AssumeRole from current user.
	case strings.HasPrefix(line, "assume-role") && connected == true:
		parts := strings.Split(line, " ")
		if len(parts) != 3 || strings.Contains(line, "help") {
			fmt.Println("Required: assume-role role-arn region")
			break
		}
		arn := parts[1]
		region = parts[2]
		sess = assumeRole(arn, region)
		target = arn

	//Assume raw json as credentials. (json format)
	case strings.HasPrefix(line, "assume-raw"):
		s := line
		r := regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)`)
		parts := r.FindAllStringSubmatch(s, -1)
		region = parts[1][0]
		json := parts[2][2]
		sess = assumeRaw(region, json)

	//The whoami/whoareyou
	case strings.HasPrefix(line, "whoami") && connected == true:
		awsdata.GetCallerIdentity(sess)

	case strings.HasPrefix(line, "get user ") && connected == true:
		parts := strings.Split(line, " ")
		if len(parts) != 3 || strings.Contains(line, "help") {
			fmt.Println("Required: get user username")
			break
		}
		username := parts[2]
		awsdata.GetUser(sess, username)

	case strings.HasPrefix(line, "get instance-profile ") && connected == true:
		parts := strings.Split(line, " ")
		profile := parts[2]
		//awsdata.GetInstanceProfile(sess, profile)
		status := awsdata.AnalyzeInstanceProfile(sess, profile)
		fmt.Printf("Privileged: %v\n", status)

	case strings.HasPrefix(line, "get policy ") && connected == true:
		parts := strings.Split(line, " ")
		arn := parts[2]
		json := awsdata.GetPolicyVersion(sess, arn)
		status := awsdata.AnalyzePolicy(json)
		fmt.Printf("Privileged: %v\n", status)

	case strings.HasPrefix(line, "get s3-policy ") && connected == true:
		parts := strings.Split(line, " ")
		bucket := parts[2]
		policy := awsdata.GetBucketPolicy(sess, bucket)
		if policy != nil {
			decodedValue, _ := url.QueryUnescape(aws.StringValue(policy.Policy))

			var prettyJSON bytes.Buffer
			error := json.Indent(&prettyJSON, []byte(decodedValue), "", "\t")
			if error != nil {
				log.Println("JSON parse error: ", error)
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

	//Execution of commands
	case strings.HasPrefix(line, "execute ssm-command") && connected == true:
		s := line
		r := regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)`)
		parts := r.FindAllString(s, -1)
		if len(parts) != 4 || strings.Contains(line, "help") {
			fmt.Println("Example: execute ssm-comand i-i00000000 \"cat /etc/passwd\"")
			fmt.Println("Execution module will run the supplied command on the given instance using SSM command execution.")
			break
		}
		instance := parts[2]
		cmd := parts[3]
		cmdid := awsdata.SendCommand(sess, instance, cmd)
		if cmdid != "" {
			time.Sleep(5 * time.Second)
			awsdata.GetCommandInvocation(sess, instance, cmdid)
		}

	//Persistence operations
	case strings.HasPrefix(line, "persist AccessKey") && connected == true:
		parts := strings.Split(line, " ")
		if len(parts) != 3 || strings.Contains(line, "help") {
			fmt.Println("Required: persist Accesskey username")
			break
		}
		username := parts[2]
		awsdata.PersistAccessKey(sess, username)

	case strings.HasPrefix(line, "persist EC2") && connected == true:
		parts := strings.Split(line, " ")
		payloadurl := parts[3]
		ami := parts[2]
		awsdata.PersistEC2(sess, ami, payloadurl)

	case strings.HasPrefix(line, "persist SSM") && connected == true:
		parts := strings.Split(line, " ")
		ami := parts[2]
		awsdata.PersistSSM(sess, ami)

	//Privesc Operations
	case strings.HasPrefix(line, "privesc UserData") && connected == true:
		parts := strings.Split(line, " ")
		payload := parts[2]
		instance := parts[3]
		awsdata.PrivescUserdata(sess, payload, instance)

	//List operations
	case strings.HasPrefix(line, "list ssm-parameters") && connected == true:
		awsdata.DescribeParameters(sess)
	case strings.HasPrefix(line, "list users") && connected == true:
		awsdata.ListUsers(sess, false)
	case strings.HasPrefix(line, "list buckets") && connected == true:
		awsdata.ListBuckets(sess, false)
	case strings.HasPrefix(line, "list groups") && connected == true:
		awsdata.ListGroups(sess)
	case strings.HasPrefix(line, "list ecs-task-definitions") && connected == true:
		awsdata.ListTaskDefinitions(sess)
	case strings.HasPrefix(line, "list groupsForUser ") && connected == true:
		parts := strings.Split(line, " ")
		username := parts[2]
		awsdata.ListGroupsForUser(sess, username)
	case strings.HasPrefix(line, "list roles") && connected == true:
		awsdata.ListRoles(sess)

	//Enumeration
	case strings.HasPrefix(line, "enum IAM") && connected == true:
		awsdata.IamEnumerate(sess, t)
		os.Setenv("AWS_EXECUTION_ENV", t)
	case strings.HasPrefix(line, "enum EC2") && connected == true:
		awsdata.EnumEC2(sess, t)
		os.Setenv("AWS_EXECUTION_ENV", t)
	case strings.HasPrefix(line, "enum S3") && connected == true:
		awsdata.EnumS3(sess)

	//Credential Discovery
	case strings.HasPrefix(line, "creds UserData") && connected == true:
		awsdata.CredentialDiscoveryUserData(sess, t)
		os.Setenv("AWS_EXECUTION_ENV", t)
	case strings.HasPrefix(line, "creds SSM") && connected == true:
		awsdata.CredentialDiscoverySSMParams(sess, t)
		os.Setenv("AWS_EXECUTION_ENV", t)
	case strings.HasPrefix(line, "creds ECS") && connected == true:
		awsdata.CredentialDiscoveryECSEnv(sess, t)
		os.Setenv("AWS_EXECUTION_ENV", t)

	//Show command history
	case line == "history":
		dat, err := ioutil.ReadFile("/tmp/readline.tmp")
		if err != nil {
			break
		}
		fmt.Print(string(dat))

	//exit
	case line == "bye":
		goto exit

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
exit:
}
