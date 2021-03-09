package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/chzyer/readline"
	"github.com/gookit/color"
)

var sess *session.Session

func main() {
	start()
}

func getProfile(pname string) *session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
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
			readline.PcItem("list",
				readline.PcItem("users"),
				readline.PcItem("groups"),
				readline.PcItem("roles"),
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
			readline.PcItem("creds",
				readline.PcItem("UserData"),
				readline.PcItem("Lambda"),
				readline.PcItem("ECR"),
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
		case strings.HasPrefix(line, "target "):
			parts := strings.Split(line, " ")
			target = parts[1]
		case strings.HasPrefix(line, "profile "):
			parts := strings.Split(line, " ")
			target = parts[1]
			sess = getProfile(target)
		case strings.HasPrefix(line, "whoami"):
			whoami(sess)
		case strings.HasPrefix(line, "list users"):
			listUsers(sess)
		case strings.HasPrefix(line, "exploit metadata"):

		case strings.HasPrefix(line, "exploit shell"):
			parts := strings.Split(line, " ")
			fmt.Println(parts)
		case line == "history":
			dat, err := ioutil.ReadFile("/tmp/readline.tmp")
			if err != nil {
				break
			}
			fmt.Print(string(dat))
		case line == "bye":
			goto exit
		case line == "sleep":
			log.Println("sleep 4 second")
			time.Sleep(4 * time.Second)
		default:
			cmdString := line
			if cmdString == "exit" {
				os.Exit(1)
			}
		}
	}
exit:
}

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func whoami(sess *session.Session) {
	svc := sts.New(sess)

	var params *sts.GetCallerIdentityInput
	resp, err := svc.GetCallerIdentity(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func listUsers(sess *session.Session) {

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
