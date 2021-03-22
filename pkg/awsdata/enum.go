package awsdata

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func IamEnumerate(sess *session.Session, t string) {
	rando := SetTrackingAction(t, "iam-enum")
	data := [][]string{}

	users := GetAccountAuthorizationDetails(sess)
	for _, user := range users {
		groupSlice := []string{}
		managedSlice := []string{}
		inlineSlice := []string{}
		privSlice := []string{}
		var isPrivileged bool
		for _, v := range user.AttachedManagedPolicies {
			managedSlice = append(managedSlice, *v.PolicyName)
			json := GetPolicyVersion(sess, *v.PolicyArn)
			isPrivileged = AnalyzePolicy(json)
			privSlice = append(privSlice, strconv.FormatBool(isPrivileged))
		}
		for _, v := range user.GroupList {
			attached := ListAttachedGroupPolicies(sess, *v)
			for _, policy := range attached {
				if policy == nil {
					continue
				}
				groupSlice = append(groupSlice, *policy.PolicyName)
				json := GetPolicyVersion(sess, *policy.PolicyArn)
				isPrivileged = AnalyzePolicy(json)
				privSlice = append(privSlice, strconv.FormatBool(isPrivileged))
			}
		}
		for _, v := range user.UserPolicyList {
			inlineSlice = append(inlineSlice, *v.PolicyName)
			decodedValue, err := url.QueryUnescape(*v.PolicyDocument)
			if err != nil {
				log.Fatal(err)
				return
			}
			isPrivileged = AnalyzePolicy(decodedValue)
			privSlice = append(privSlice, strconv.FormatBool(isPrivileged))
		}
		row := []string{*user.UserName, strings.Join(managedSlice, "\n"), strings.Join(inlineSlice, "\n"), strings.Join(groupSlice, "\n"), strconv.FormatBool(contains1(privSlice, "true"))}
		data = append(data, row)
	}

	//data = append(data, row)
	fmt.Println("UA Tracking: exec-env/" + rando)
	header := []string{"User", "Managed Policies", "Inline Policies", "Groups", "isPrivileged"}
	tableData(data, header)
}

func EnumEC2(sess *session.Session, t string) {
	rando := SetTrackingAction(t, "ec2-enum")
	data := [][]string{}

	ec2s := DescribeInstances(sess)
	for _, e := range ec2s {
		for _, i := range e.Instances {
			var isPub bool
			var isPrivileged bool
			var profileName string

			roleARN := ""
			var sgroups []string
			var ports []string

			if i.IamInstanceProfile != nil {
				roleARN = *i.IamInstanceProfile.Arn
				s := roleARN
				r := regexp.MustCompile(`instance-profile\/(.*)`)
				parts := r.FindAllStringSubmatch(s, -1)
				for _, v := range parts {
					profileName = v[1]
				}
				isPrivileged = AnalyzeInstanceProfile(sess, profileName)
			} else {
				roleARN = "None"
			}
			for _, v := range i.SecurityGroups {
				sgroups = append(sgroups, *v.GroupId)
				data := DescribeSecurityGroup(sess, *v.GroupId)
				for _, sg := range data.SecurityGroups {
					for _, perms := range sg.IpPermissions {

						portInt := checkPointerInt(perms.FromPort)
						port := strconv.Itoa(int(portInt))
						if port == "0" {
							port = "All"
						}
						for _, i := range perms.IpRanges {
							var pubIP *ec2.IpRange
							pubIP = &ec2.IpRange{
								CidrIp: aws.String("0.0.0.0/0"),
							}
							if *i.CidrIp == *pubIP.CidrIp {
								isPub = true
								ports = append(ports, port+"*")
							} else {
								ports = append(ports, port+"-")
							}
						}
					}
				}
			}
			row := []string{checkPointer(i.InstanceId), roleARN, checkPointer(i.VpcId), checkPointer(i.PublicIpAddress), checkPointer(i.PrivateIpAddress), strings.Join(sgroups, "\n"), strings.Join(ports, "\n"), checkPointer(i.State.Name), strconv.FormatBool(isPrivileged), strconv.FormatBool(isPub)}
			data = append(data, row)
		}
	}
	fmt.Println("UA Tracking: exec-env/" + rando)
	header := []string{"InstanceID", "Instance Profile", "VPC", "PublicIP", "PrivateIP", "Security Groups", "Ports", "State", "isPrivileged", "isPublic"}
	tableData(data, header)
}

func checkPointer(pointer *string) string {
	if pointer != nil {
		return *pointer
	}
	return "None"

}

func checkPointerInt(pointer *int64) int64 {
	if pointer != nil {
		return *pointer
	}
	return 0

}

func EnumS3(sess *session.Session) {
	data := [][]string{}

	buckets := ListBuckets(sess, true)

	for _, v := range buckets {
		var isWebsite bool
		var allowPublic bool
		var allowAuthenticated bool
		var hasPolicy bool
		var permsPub []string
		var permsAuth []string

		policy := GetBucketPolicy(sess, v)
		if policy != nil {
			hasPolicy = true
		}
		//GetBucketReplication(sess, v) *g.Grantee.URI == "http://acs.amazonaws.com/groups/global/AuthenticatedUsers"
		acl := GetBucketACL(sess, v)
		if acl != nil {
			for _, g := range acl.Grants {
				if g.Grantee.URI != nil {
					if *g.Grantee.URI == "http://acs.amazonaws.com/groups/global/AllUsers" {
						allowPublic = true
						permsPub = append(permsPub, *g.Permission)

					}
					if *g.Grantee.URI == "http://acs.amazonaws.com/groups/global/AuthenticatedUsers" {
						allowAuthenticated = true
						permsAuth = append(permsAuth, *g.Permission)

					}
				}
			}
		}
		website := GetBucketWebsite(sess, v)
		if website != nil {
			isWebsite = true
		}
		row := []string{v, strconv.FormatBool(hasPolicy), strconv.FormatBool(isWebsite), strconv.FormatBool(allowPublic), strings.Join(permsPub, ","), strconv.FormatBool(allowAuthenticated), strings.Join(permsAuth, ",")}
		data = append(data, row)
	}
	header := []string{"Bucket", "hasPolicy", "isWebsite", "Allow Public", "Permissions", "Allow Authenticated", "Permissions"}
	tableData(data, header)
	fmt.Println("---Follow on actions---")
	fmt.Println("get s3-policy policy")
	fmt.Println("get s3-acl acl")
}
