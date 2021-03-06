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

func RoleEnumerate(sess *session.Session, t string) {
	rando := SetTrackingAction(t, "role-enum")
	data := [][]string{}
	var ident string
	var isPrivileged bool

	roles := ListRoles(sess, true)
	fmt.Println("Building Table...\n")
	for _, r := range roles.Roles {
		decodedValue, err := url.QueryUnescape(aws.StringValue(r.AssumeRolePolicyDocument))
		if err != nil {
			log.Fatal(err)
		}

		principalType, identity, _ := GetTrustPolicy(decodedValue)
		ident = fmt.Sprintf("%v", identity)
		policies := ListAttachedRolePolicies(sess, *r.RoleName)
		for _, p := range policies {
			json := GetPolicyVersion(sess, *p.PolicyArn)
			status := AnalyzePolicy(json)
			if status == true {
				isPrivileged = true
			} else {
				isPrivileged = false
			}
		}
		row := []string{*r.RoleName, principalType[0], ident, strconv.FormatBool(isPrivileged)}
		data = append(data, row)
	}

	fmt.Println("UA Tracking: exec-env/" + rando)
	header := []string{"Role", "Principal Type", "Identity/Service", "isPrivileged"}
	tableData(data, header)
}

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

func UserGroupEnumerate(sess *session.Session, t string) {
	rando := SetTrackingAction(t, "usergroup-enum")
	data := [][]string{}

	GetCallerIdentity(sess)
	users := ListUsers(sess, true)
	for _, u := range users {
		managedSlice := []string{}
		privSlice := []string{}
		var isPrivileged bool

		inlineSlice := []string{}
		privINSlice := []string{}
		var isINPrivileged bool
		Groups := ListGroupsForUser(sess, u)
		for _, g := range Groups {
			for _, v := range ListAttachedGroupPolicies(sess, *g.GroupName) {
				managedSlice = append(managedSlice, *v.PolicyName)
				json := GetPolicyVersion(sess, *v.PolicyArn)
				isPrivileged = AnalyzePolicy(json)
				privSlice = append(privSlice, strconv.FormatBool(isPrivileged))
			}
			gpol := ListGroupPolicies(sess, *g.GroupName)
			for _, v := range gpol.PolicyNames {
				inlineSlice = append(inlineSlice, *v)
				json := GetGroupPolicy(sess, *g.GroupName, *v)
				isINPrivileged = AnalyzePolicy(json)
				privINSlice = append(privINSlice, strconv.FormatBool(isINPrivileged))
			}

		}
		row := []string{u, strings.Join(managedSlice, "\n"), strconv.FormatBool(contains1(privSlice, "true")), strings.Join(inlineSlice, "\n"), strconv.FormatBool(contains1(privINSlice, "true"))}
		data = append(data, row)
	}
	fmt.Println("UA Tracking: exec-env/" + rando)
	header := []string{"User", "Policies", "isPrivileged", "Inline Policies", "isPrivileged"}
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

func EnumS3(sess *session.Session, t string) {
	rando := SetTrackingAction(t, "s3-enum")
	data := [][]string{}

	buckets := ListBuckets(sess, true)

	for _, v := range buckets {
		var isWebsite bool
		var allowPublic bool
		var repStatus bool
		var allowAuthenticated bool
		var hasPolicy bool
		var permsPub []string
		var permsAuth []string

		policy := GetBucketPolicy(sess, v)
		if policy != nil {
			hasPolicy = true
		}
		region := GetBucketLocation(sess, v)
		reps := GetBucketReplication(sess, v)
		if reps == "exists" || reps == "AccessDenied" {
			repStatus = true
		}
		if reps == "ReplicationConfigurationNotFoundError" {
			repStatus = false
		}
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
		row := []string{v, strconv.FormatBool(hasPolicy), strconv.FormatBool(isWebsite), strconv.FormatBool(allowPublic), strings.Join(permsPub, ","), strconv.FormatBool(allowAuthenticated), strings.Join(permsAuth, ","), strconv.FormatBool(repStatus), region}
		data = append(data, row)
	}
	fmt.Println("UA Tracking: exec-env/" + rando)
	header := []string{"Bucket", "hasPolicy", "isWebsite", "Allow Public", "Permissions", "Allow Authenticated", "Permissions", "Replication", "Region"}
	tableData(data, header)
}

func GroupsEnumerate(sess *session.Session, t string) {
	rando := SetTrackingAction(t, "groups-enum")
	data := [][]string{}

	groups := ListGroups(sess, true)
	for _, group := range groups.Groups {
		managedSlice := []string{}
		privSlice := []string{}
		var isPrivileged bool

		inlineSlice := []string{}
		privINSlice := []string{}
		var isINPrivileged bool

		for _, v := range ListAttachedGroupPolicies(sess, *group.GroupName) {
			managedSlice = append(managedSlice, *v.PolicyName)
			json := GetPolicyVersion(sess, *v.PolicyArn)
			isPrivileged = AnalyzePolicy(json)
			privSlice = append(privSlice, strconv.FormatBool(isPrivileged))
		}
		gpol := ListGroupPolicies(sess, *group.GroupName)
		for _, v := range gpol.PolicyNames {
			inlineSlice = append(inlineSlice, *v)
			json := GetGroupPolicy(sess, *group.GroupName, *v)
			isINPrivileged = AnalyzePolicy(json)
			privINSlice = append(privINSlice, strconv.FormatBool(isINPrivileged))
		}
		row := []string{*group.GroupName, strings.Join(managedSlice, "\n"), strconv.FormatBool(contains1(privSlice, "true")), strings.Join(inlineSlice, "\n"), strconv.FormatBool(contains1(privINSlice, "true"))}
		data = append(data, row)
	}

	//data = append(data, row)
	fmt.Println("UA Tracking: exec-env/" + rando)
	header := []string{"Group", "Policies", "isPrivileged", "Inline Policies", "isPrivileged"}
	tableData(data, header)
}

func EnumCrossAccount(sess *session.Session, t string) {
	rando := SetTrackingAction(t, "crossaccount-lateral")
	data := [][]string{}

	events := LookupEvents(sess, "AssumeRole")
	fmt.Println(len(events))
	var RoleAccountID string
	for _, event := range events {
		roleArn := event.RequestParameters.RoleArn
		r := regexp.MustCompile(`::(.*):`)
		parts := r.FindAllStringSubmatch(roleArn, -1)
		for _, v := range parts {
			RoleAccountID = v[1]
		}
		if event.RecipientAccountID != RoleAccountID {
			if event.RequestParameters.RoleArn != "" {
				row := []string{event.UserIdentity.Type, event.RequestParameters.RoleArn, event.UserIdentity.Arn, event.EventTime}
				data = append(data, row)
			}
		}
	}
	fmt.Println("UA Tracking: exec-env/" + rando)
	header := []string{"Identity Type", "Lateral ARN", "Identity ARN", "Date"}
	tableData(data, header)
}

func EnumLateralNetwork(sess *session.Session, t string) {
	rando := SetTrackingAction(t, "lateral-network-enum")
	//data := [][]string{}

	gateways := DescribeDirectConnectGateways(sess)
	for _, g := range gateways.DirectConnectGateways {
		fmt.Println(g.DirectConnectGatewayId)
		DescribeDirectConnectGatewayAssociations(sess, *g.DirectConnectGatewayId)
	}
	vpns := DescribeVpnConnections(sess)
	fmt.Println(vpns.VpnConnections)
	peers := DescribeVpcPeeringConnections(sess)
	fmt.Println(peers.VpcPeeringConnections)
	fmt.Println("UA Tracking: exec-env/" + rando)
}

func EnumOrg(sess *session.Session, t string) {
	rando := SetTrackingAction(t, "account-enum")
	ListAccountAliases(sess)
	DescribeOrganization(sess)
	fmt.Println("UA Tracking: exec-env/" + rando)
}
