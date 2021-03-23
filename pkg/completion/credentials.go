package completion

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
)

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

//Assume role from roleArn
func assumeRole(arn string, region string) *session.Session {
	creds := stscreds.NewCredentials(sess, arn)
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: creds,
	})
	if err != nil {
		fmt.Println(err)
		return nil
	} else {
		connected = true
	}
	_, err = sess.Config.Credentials.Get()
	if err != nil {
		fmt.Println(err)
		return nil
	} else {
		connected = true
	}
	return sess
}

//Assume raw json token set
func assumeRaw(region string, data string) *session.Session {
	var token awsToken
	err := json.Unmarshal([]byte(data), &token)

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(token.AccessKeyID, token.SecretAccessKey, token.Token),
	})
	if err != nil {
		fmt.Println(err)
		return nil
	} else {
		connected = true
	}
	_, err = sess.Config.Credentials.Get()
	if err != nil {
		fmt.Println(err)
		return nil
	} else {
		connected = true
		target = token.AccessKeyID
	}
	return sess
}

//GetSessionToken for current profile
func stsSession(access string, secret string, sessiont string, region string) *session.Session {

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(access, secret, sessiont),
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
		target = access
	}
	return sess
}
