package completion

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/user"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
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
		if connected == true {
			users := awsdata.ListUsers(sess, true)
			return users
		}
		return nil
	}
}

func listEc2(sess *session.Session) func(string) []string {
	return func(line string) []string {
		if connected == true {
			var results []string
			instances := awsdata.DescribeInstances(sess)
			for _, v := range instances {
				for _, i := range v.Instances {
					results = append(results, *i.InstanceId)
				}
			}

			return results
		}
		return nil
	}
}

func listRoles(sess *session.Session) func(string) []string {
	return func(line string) []string {
		if connected == true {
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
		return nil
	}
}

func listRolesName(sess *session.Session) func(string) []string {
	return func(line string) []string {
		if connected == true {
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
				results = append(results, *role.RoleName)
			}
			return results
		}
		return nil
	}
}

func listTrails(sess *session.Session) func(string) []string {
	return func(line string) []string {
		if connected == true {
			var results []string
			svc := cloudtrail.New(sess)

			result, err := svc.DescribeTrails(&cloudtrail.DescribeTrailsInput{})

			if err != nil {
				fmt.Println("Error", err)
				return nil
			}

			for _, trail := range result.TrailList {
				if trail == nil {
					continue
				}
				results = append(results, *trail.Name)
			}
			return results
		}
		return nil
	}
}

func listBuckets(sess *session.Session) func(string) []string {
	return func(line string) []string {
		if connected == true {
			var results []string

			buckets := awsdata.ListBuckets(sess, true)

			for _, b := range buckets {
				results = append(results, b)
			}
			return results
		}
		return nil
	}
}

func GetRegions() []string {
	allRegions := []string{"us-east-1", "us-east-2", "us-west-1", "us-west-2", "ca-central-1",
		"eu-west-1", "eu-west-2", "eu-west-3", "eu-north-1", "ap-northeast-1", "ap-northeast-2",
		"ap-northeast-3", "ap-southeast-1", "ap-southeast-2", "ap-south-1", "sa-east-1"}

	return allRegions
}
