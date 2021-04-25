package awsdata

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"github.com/drk1wi/Modlishka/log"
)

type FullEvent struct {
	EventVersion string
	UserIdentity struct {
		Type        string
		PrincipalID string
		Arn         string
		AccountID   string
		UserName    string
		AccessKeyId string
	}
	EventTime         string
	EventSource       string
	EventName         string
	AwsRegion         string
	SourceIPAddress   string
	UserAgent         string
	RequestParameters struct {
		RoleArn string
	}
	ResponseElements    map[string]interface{}
	AdditionalEventData struct {
		LoginTo       string
		MobileVersion string
		MFAUsed       string
	}
	EventID            string
	EventType          string
	RecipientAccountID string
	SharedEventID      string
}

func LookupEvents(sess *session.Session, event string) []FullEvent {
	var input *cloudtrail.LookupEventsInput
	var events []FullEvent

	now := time.Now().UTC()
	count := 100
	now = now.Add(time.Duration(-count) * time.Hour)

	svc := cloudtrail.New(sess)
	eventCount := 0
	input = &cloudtrail.LookupEventsInput{
		LookupAttributes: []*cloudtrail.LookupAttribute{{
			AttributeKey:   aws.String("EventName"),
			AttributeValue: aws.String(event),
		}},
		StartTime: &now,
	}
	result, err := svc.LookupEvents(input)
	for _, event := range result.Events {
		eventCount++
		fullEvent := FullEvent{}
		json.Unmarshal([]byte(*event.CloudTrailEvent), &fullEvent)
		events = append(events, fullEvent)
	}
	for {
		if result.NextToken == nil {
			break
		}
		log.Infof("%d events found. Continuing...", eventCount)
		input.NextToken = result.NextToken

		result, err = svc.LookupEvents(input)
		if err != nil {
			log.Fatal(err)
		}
		for _, event := range result.Events {
			fullEvent := FullEvent{}
			json.Unmarshal([]byte(*event.CloudTrailEvent), &fullEvent)
			events = append(events, fullEvent)
		}
	}
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}

	return events
}

func DescribeTrails(sess *session.Session) []*cloudtrail.Trail {
	// Create CloudTrail client
	svc := cloudtrail.New(sess)
	resp, err := svc.DescribeTrails(&cloudtrail.DescribeTrailsInput{TrailNameList: nil})
	if err != nil {
		fmt.Println("Got error calling CreateTrail:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return resp.TrailList
}

func ListTrails() {

}

func GetEventSelectors() {

}

func StopLogging(sess *session.Session, trail string) bool {

	svc := cloudtrail.New(sess)

	_, err := svc.StopLogging(&cloudtrail.StopLoggingInput{
		Name: aws.String(trail),
	})
	if err != nil {
		fmt.Println("Got error calling StopLogging:")
		fmt.Println(err.Error())
		return false
	}

	return true
}

func DeleteTrail(sess *session.Session, trail string) bool {

	svc := cloudtrail.New(sess)

	_, err := svc.DeleteTrail(&cloudtrail.DeleteTrailInput{Name: aws.String(trail)})
	if err != nil {
		fmt.Println("Got error calling DeleteTrail:")
		fmt.Println(err.Error())
		return false
	}

	return true
}

func UpdateTrail() {

}
