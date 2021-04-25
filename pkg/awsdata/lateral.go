package awsdata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
)

type Token struct {
	Signintoken string `json:"SigninToken"`
}

func LateralConsole(sess *session.Session, user string) {
	token := GetFederationToken(sess, user)

	//Construct signin url
	getSigninToken := fmt.Sprintf(`{"sessionId":"%s","sessionKey":"%s","sessionToken":"%s"}`, *token.Credentials.AccessKeyId, *token.Credentials.SecretAccessKey, *token.Credentials.SessionToken)

	params := url.Values{}
	params.Add("Action", "getSigninToken")
	params.Add("SessionDuration", "43200")
	params.Add("Session", getSigninToken)
	signInUrl := "https://signin.aws.amazon.com/federation?" + params.Encode()

	//Grab SignInToken
	consoleClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, signInUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := consoleClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	toke := Token{}
	jsonErr := json.Unmarshal(body, &toke)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	params = url.Values{}
	params.Add("Action", "login")
	params.Add("Issuer", "")
	params.Add("Destination", "https://console.aws.amazon.com/console/home")
	params.Add("SigninToken", toke.Signintoken)
	signInUrl = "https://signin.aws.amazon.com/federation?" + params.Encode()

	fmt.Println(signInUrl)
}
