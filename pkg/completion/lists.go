package completion

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/user"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/grines/scour/pkg/awsdata"
)

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
