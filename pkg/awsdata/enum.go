package awsdata

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
)

func IamEnumerate(sess *session.Session) {
	data := [][]string{}
	groupSlice := []string{}
	managedSlice := []string{}
	inlineSlice := []string{}

	users := GetAccountAuthorizationDetails(sess)
	for _, user := range users {
		for _, v := range user.AttachedManagedPolicies {
			managedSlice = append(managedSlice, *v.PolicyName)
		}
		for _, v := range user.GroupList {
			attached := ListAttachedGroupPolicies(sess, *v)
			for _, policy := range attached {
				if policy == nil {
					continue
				}
				groupSlice = append(groupSlice, *policy.PolicyName)
			}
		}
		for _, v := range user.UserPolicyList {
			inlineSlice = append(inlineSlice, *v.PolicyName)
			//decodedValue, err := url.QueryUnescape(*v.PolicyDocument)
			//if err != nil {
			//	log.Fatal(err)
			//	return
			//}
			//fmt.Println(decodedValue)
		}
		row := []string{*user.UserName, strings.Join(managedSlice, "\n"), strings.Join(inlineSlice, "\n"), strings.Join(groupSlice, "\n")}
		data = append(data, row)
	}

	//data = append(data, row)
	header := []string{"User", "Managed Policies", "Inline Policies", "Groups", "isAdmin", "isPrivileged"}
	tableData(data, header)
}
