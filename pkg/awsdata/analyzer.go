package awsdata

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

type Value []string

type PolicyDocument struct {
	Version   string
	Statement []StatementEntry
}

type StatementEntry struct {
	Effect    string
	Principal map[string]interface{}
	Action    interface{}
	Resource  interface{}
}

//AnalyzePolicy is used to detect over privileged IAM policies
func AnalyzePolicy(data string) bool {
	var doc PolicyDocument
	var isPrivileged bool
	var effect string

	//Privileged Actions
	privs := []string{"*", "ec2:*", "iam:*", "s3:*", "lambda:*", "*:*"}

	err1 := json.Unmarshal([]byte(data), &doc)

	if err1 != nil {
		fmt.Println(err1)
	}

	for _, v := range doc.Statement {
		actions := parseObject(v.Action)
		isPrivileged = contains(actions, privs)
		effect = v.Effect
		if effect == "Allow" {
			return isPrivileged
		}
	}

	return false
}

func GetTrustPolicy(data string) ([]string, []interface{}, map[string]interface{}) {
	var doc PolicyDocument
	var principalType []string
	//var effect string
	var identity []interface{}
	var principal map[string]interface{}

	//Privileged Actions
	//privs := []string{"*", "ec2:*", "iam:*", "s3:*", "lambda:*", "*:*"}

	err1 := json.Unmarshal([]byte(data), &doc)

	if err1 != nil {
		fmt.Println(err1)
	}

	for _, v := range doc.Statement {
		principal = v.Principal
		for k, p := range v.Principal {
			principalType = append(principalType, k)
			identity = append(identity, p)

		}
	}

	return principalType, identity, principal
}

func parseObject(t interface{}) []string {
	var actions []string
	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(t)

		for i := 0; i < s.Len(); i++ {
			y := s.Index(i).Interface().(string)
			actions = append(actions, y)
		}
	case reflect.String:
		s := reflect.ValueOf(t)
		y := s.Interface().(string)
		actions = append(actions, y)
	}
	return actions
}

func contains(s []string, str []string) bool {
	for _, v := range s {
		for _, p := range str {
			if v == p {
				//fmt.Printf("Match: %v\n", p)
				return true
			}
		}
	}
	return false
}

func contains1(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func AnalyzeInstanceProfile(sess *session.Session, profile string) bool {
	data := GetInstanceProfile(sess, profile, true)
	if data != nil {
		for _, v := range data.InstanceProfile.Roles {
			policies := ListAttachedRolePolicies(sess, *v.RoleName)
			for _, p := range policies {
				json := GetPolicyVersion(sess, *p.PolicyArn)
				status := AnalyzePolicy(json)
				if status == true {
					return true
				}
			}
		}
	}
	return false
}

func AnalyzeTrustPolicy(data string, userARN string, accountID string) bool {

	var doc PolicyDocument
	var isAllowed bool
	var effect string
	var checks []string
	var userCheck string
	var roleCheck string

	//Allowed ARNS
	root := "arn:aws:iam::accountid:root"
	user := "arn:aws:iam::accountid:user/userid"
	role := "arn:aws:iam::accountid:role/userid"

	regexRole := `assumed-role\/(.*)\/`
	regexUser := `user\/(.*)`

	r, _ := regexp.Compile(regexRole)
	if r.MatchString(userARN) {
		matches := r.FindStringSubmatch(string(userARN))
		roleCheck = strings.ReplaceAll(role, "accountid", accountID)
		roleCheck = strings.ReplaceAll(roleCheck, "userid", matches[1])
	}
	r, _ = regexp.Compile(regexUser)
	if r.MatchString(userARN) {
		matches := r.FindStringSubmatch(string(userARN))
		userCheck = strings.ReplaceAll(user, "accountid", accountID)
		userCheck = strings.ReplaceAll(userCheck, "userid", matches[1])
	}

	//Hacky way to check for access..
	rootCheck := strings.ReplaceAll(root, "accountid", accountID)

	checks = append(checks, userCheck)
	checks = append(checks, rootCheck)
	checks = append(checks, roleCheck)

	err1 := json.Unmarshal([]byte(data), &doc)

	if err1 != nil {
		fmt.Println(err1)
	}

	for _, v := range doc.Statement {
		arns := parser(v.Principal)
		//for _, arn := range arns {
		//	fmt.Println(arn)
		//}
		isAllowed = contains(arns, checks)
		//fmt.Println(isAllowed)
		effect = v.Effect
		if effect == "Allow" && v.Action == "sts:AssumeRole" {
			return isAllowed
		}
	}

	return false
}

func parser(t map[string]interface{}) []string {
	var arns []string
	for k, v := range t {
		switch c := v.(type) {
		case string:
			if k == "AWS" {
				arns = append(arns, c)
			}
		case []interface{}:
			for _, v2 := range c {
				if k == "AWS" {
					arn := fmt.Sprintf("%s", v2)
					arns = append(arns, arn)
				}
			}
		}

	}
	return arns
}

func AnalyzeRoleTrustRelationships(sess *session.Session) {
	svc := sts.New(sess)

	var params *sts.GetCallerIdentityInput
	resp, err := svc.GetCallerIdentity(params)

	accountID := *resp.Account
	userARN := *resp.Arn

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	data := [][]string{}

	roles := ListRoles(sess, true)
	for _, v := range roles.Roles {
		decodedValue, err := url.QueryUnescape(aws.StringValue(v.AssumeRolePolicyDocument))
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(decodedValue)
		status := AnalyzeTrustPolicy(decodedValue, userARN, accountID)
		if status == true {
			row := []string{*v.Arn, strconv.FormatBool(status)}
			data = append(data, row)
		}
	}
	header := []string{"Role", "Can Assume", "isPrivileged"}
	tableData(data, header)
	fmt.Println("Try: assume-role arn:aws:iam::accountid:role/rolename")
}

func TestInterface(t interface{}) []string {
	var idents []string
	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(t)

		for i := 0; i < s.Len(); i++ {
			id := fmt.Sprintf(s.Index(i).Elem().String())
			idents = append(idents, id)
		}
	case reflect.String:
		s := reflect.ValueOf(t)
		idents = append(idents, s.String())
	}
	return idents

}
