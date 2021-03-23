package completion

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/chzyer/readline"
	"github.com/gookit/color"
)

func Start(t string) {

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

		//Autocompletion configuration
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
				readline.PcItem("policies"),
			),
			readline.PcItem("get",
				readline.PcItem("user",
					readline.PcItemDynamic(listUsers(sess)),
				),
				readline.PcItem("group"),
				readline.PcItem("role"),
				readline.PcItem("policy"),
				readline.PcItem("s3-policy",
					readline.PcItemDynamic(listBuckets(sess)),
				),
				readline.PcItem("s3-acl",
					readline.PcItemDynamic(listBuckets(sess)),
				),
			),
		)

		//readlines configuration
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
			l.SetPrompt(green("Connected") + " <" + blue(target+"/"+region) + "> ")
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

		//Trimwhitespace and send to commands switch
		line = strings.TrimSpace(line)
		Commands(line, t)
	}
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
