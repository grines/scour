package completion

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
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

type help struct {
	helpText     string
	infoText     string
	autocomplete string
}
