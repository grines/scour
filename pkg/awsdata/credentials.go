package awsdata

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
)

type Finding struct {
	Rule    string `json:"rule"`
	Finding string `json:"finding"`
}

func CredentialDiscoveryLambda(sess *session.Session, t string) {
	rando := SetTrackingAction(t, "lambda-creds")
	data := [][]string{}
	var envVariables []string

	functions := ListFunctions(sess)
	for _, f := range functions.Functions {
		function := GetFunction(sess, *f.FunctionName)
		for k, v := range function.Configuration.Environment.Variables {
			envVariables = append(envVariables, k+":"+*v)
		}
		row := []string{*function.Configuration.FunctionName, strings.Join(envVariables, ","), *function.Code.RepositoryType}
		data = append(data, row)
	}
	fmt.Println("UA Tracking: exec-env/" + rando)
	header := []string{"Function", "Env Variables", "Code Location"}
	tableData(data, header)
}

func CredentialDiscoveryUserData(sess *session.Session, t string) {
	rando := SetTrackingAction(t, "userdata-creds")
	data := [][]string{}

	instances := DescribeInstances(sess)
	for _, v := range instances {
		for _, i := range v.Instances {
			userData := DescribeInstanceAttribute(sess, *i.InstanceId, "userData")
			if userData.UserData.Value != nil {
				findings := secret(*userData.UserData.Value, *i.InstanceId)
				if len(findings) >= 1 {
					for _, f := range findings {
						row := []string{*i.InstanceId, f.Rule, f.Finding}
						data = append(data, row)
					}
				}
			}
		}

	}
	fmt.Println("UA Tracking: exec-env/" + rando)
	header := []string{"InstanceID", "Rule", "Finding"}
	tableData(data, header)
}

func CredentialDiscoverySSMParams(sess *session.Session, t string) {
	rando := SetTrackingAction(t, "ssm-params-creds")
	data := [][]string{}

	params := DescribeParameters(sess)
	for _, v := range params {
		params := GetParameter(sess, *v.Name)
		row := []string{*params.Parameter.Name, *params.Parameter.DataType, *params.Parameter.Value}
		data = append(data, row)

	}
	fmt.Println("UA Tracking: exec-env/" + rando)
	header := []string{"Param Name", "DataType", "Value"}
	tableData(data, header)
}

func CredentialDiscoveryECSEnv(sess *session.Session, t string) {
	rando := SetTrackingAction(t, "ecs-creds")
	data := [][]string{}

	defs := ListTaskDefinitions(sess)
	for _, v := range defs.TaskDefinitionArns {
		env := DescribeTaskDefinition(sess, *v)
		for _, t := range env.TaskDefinition.ContainerDefinitions {
			for _, e := range t.Environment {
				row := []string{*e.Name, *e.Value, *t.Name}
				data = append(data, row)
			}
		}

	}
	fmt.Println("UA Tracking: exec-env/" + rando)
	header := []string{"Envars name", "Value", "Definition"}
	tableData(data, header)
}

func secret(data64 string, instance string) []Finding {

	var findings []Finding
	r := BuildRules()

	decoded, err := base64.StdEncoding.DecodeString(data64)
	if err != nil {
		fmt.Println("decode error:", err)
		return nil
	}

	for _, rule := range r {
		r, _ := regexp.Compile(rule.Exp)
		if r.MatchString(string(decoded)) {
			row := Finding{
				Rule:    rule.Name,
				Finding: r.FindString(string(decoded)),
			}
			findings = append(findings, row)
		}

	}
	return findings
}
