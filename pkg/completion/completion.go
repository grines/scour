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

			//Grab Credentials
			readline.PcItem("token",
				readline.PcItem("profile",
					readline.PcItemDynamic(listProfiles(),
						readline.PcItem("us-east-1"),
						readline.PcItem("us-east-2"),
						readline.PcItem("us-west-1"),
						readline.PcItem("us-west-2"),
					),
				),
				readline.PcItem("AssumeRole",
					readline.PcItemDynamic(listRoles(sess),
						readline.PcItem("us-east-1"),
						readline.PcItem("us-east-2"),
						readline.PcItem("us-west-1"),
						readline.PcItem("us-west-2"),
					),
				),

				readline.PcItem("AssumeRaw",
					readline.PcItem("us-east-1"),
					readline.PcItem("us-east-2"),
					readline.PcItem("us-west-1"),
					readline.PcItem("us-west-2"),
				),
				readline.PcItem("GetSessionToken",
					readline.PcItem("us-east-1"),
					readline.PcItem("us-east-2"),
					readline.PcItem("us-west-1"),
					readline.PcItem("us-west-2"),
				),
			),

			//ATTACK Command completion
			readline.PcItem("attack",
				readline.PcItem("operations",
					readline.PcItem("s3-ransom",
						readline.PcItemDynamic(listBuckets(sess),
							readline.PcItem("arn:aws:kms:REGION:ACCOUNT-ID:key/KEY-ID"),
						),
					),
				),
				readline.PcItem("enum",
					readline.PcItem("IAM"),
					readline.PcItem("Roles"),
					readline.PcItem("EC2"),
					readline.PcItem("S3"),
					readline.PcItem("Groups"),
					readline.PcItem("CrossAccount"),
					readline.PcItem("Network"),
				),

				readline.PcItem("privesc",
					readline.PcItem("UserData",
						readline.PcItemDynamic(listEc2(sess),
							readline.PcItem("http://url.to.capture.post.data"),
						),
					),
					readline.PcItem("CreateLoginProfile",
						readline.PcItemDynamic(listUsers(sess)),
					),
					readline.PcItem("IAM"),
				),

				readline.PcItem("lateral",
					readline.PcItem("ConsoleLogin",
						readline.PcItemDynamic(listUsers(sess)),
					),
					readline.PcItem("CrossAccount"),
				),

				readline.PcItem("evasion",
					readline.PcItem("KillCloudTrail",
						readline.PcItemDynamic(listTrails(sess)),
					),
					readline.PcItem("StopCloudTrail",
						readline.PcItemDynamic(listTrails(sess)),
					),
					readline.PcItem("KillGuardDuty"),
					readline.PcItem("UpdateGuardDuty",
						readline.PcItemDynamic(listTrails(sess)),
					),
					readline.PcItem("TrustIPGuardDuty",
						readline.PcItem("https://scouriplist.s3-us-west-1.amazonaws.com/iplist.txt"),
					),
					readline.PcItem(""),
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
					readline.PcItem("Role",
						readline.PcItemDynamic(listRolesName(sess)),
					),
					readline.PcItem("CreateUser"),
					readline.PcItem("CodeBuild"),
					readline.PcItem("CrossAccount"),
					readline.PcItem("EC2",
						readline.PcItem("ami-013f17f36f8b1fefb",
							readline.PcItem("https://url.to.compiled.payload/payload"),
						),
					),
					readline.PcItem("SSM",
						readline.PcItemDynamic(listEc2(sess)),
					),
				),

				readline.PcItem("exfil",
					readline.PcItem("ec2Snapshot",
						readline.PcItemDynamic(listEc2(sess),
							readline.PcItem("<accountid>"),
						),
					),
				),
			),

			readline.PcItem("whoami"),

			//AWS cli implementation
			readline.PcItem("aws",

				readline.PcItem("iam",
					readline.PcItem("get-session-token"),
					readline.PcItem("list-users"),
					readline.PcItem("create-user"),
					readline.PcItem("list-groups"),
					readline.PcItem("create-group"),
					readline.PcItem("list-groups-for-user",
						readline.PcItemDynamic(listUsers(sess)),
					),
					readline.PcItem("add-user-to-group",
						readline.PcItemDynamic(listUsers(sess)),
					),
					readline.PcItem("list-roles"),
					readline.PcItem("create-role"),
					readline.PcItem("list-policies"),
					readline.PcItem("create-policy"),
					readline.PcItem("get-user",
						readline.PcItemDynamic(listUsers(sess)),
					),
					readline.PcItem("get-group"),
					readline.PcItem("get-role"),
					readline.PcItem("get-policy"),
					readline.PcItem("create-access-key",
						readline.PcItemDynamic(listUsers(sess)),
					),
					readline.PcItem("create-login-profile",
						readline.PcItemDynamic(listUsers(sess)),
					),
				),

				readline.PcItem("s3",
					readline.PcItem("list-buckets"),
				),

				readline.PcItem("ssm",
					readline.PcItem("send-command"),
				),

				readline.PcItem("ecs",
					readline.PcItem("send-command"),
				),

				readline.PcItem("ec2",
					readline.PcItem("describe-instances"),
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
