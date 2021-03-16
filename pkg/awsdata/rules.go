package awsdata

type rules []struct {
	Name string `json:"name"`
	Exp  string `json:"exp"`
}

func BuildRules() rules {
	r := rules{
		{
			Name: "Slack Webhook",
			Exp:  "https://hooks.slack.com/services/T[a-zA-Z0-9_]{8}/B[a-zA-Z0-9_]{8}/[a-zA-Z0-9_]{24}",
		},
		{
			Name: "Generic Password",
			Exp:  "pass=(.*)",
		},
		{
			Name: "Generic Password",
			Exp:  "password=(.*)",
		},
		{
			Name: "AWS API Key",
			Exp:  "AKIA[0-9A-Z]{16}",
		},
	}

	return r
}
