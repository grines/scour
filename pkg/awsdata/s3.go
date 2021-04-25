package awsdata

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func ListBuckets(sess *session.Session, hide bool) []string {
	data := [][]string{}
	var buckets []string

	svc := s3.New(sess)
	input := &s3.ListBucketsInput{}

	result, err := svc.ListBuckets(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}

	for _, bucket := range result.Buckets {
		if bucket == nil {
			continue
		}
		//fmt.Printf("%d user %s created %v\n", i, *user.UserName, user.CreateDate)
		buckets = append(buckets, *bucket.Name)
		row := []string{*bucket.Name, bucket.CreationDate.String()}
		data = append(data, row)
	}
	header := []string{"BucketName", "CreateDate"}
	if hide == false {
		tableData(data, header)
	}
	return buckets
}

func GetBucketPolicy(sess *session.Session, bucket string) *s3.GetBucketPolicyOutput {
	svc := s3.New(sess)
	input := &s3.GetBucketPolicyInput{
		Bucket: aws.String(bucket),
	}

	result, err := svc.GetBucketPolicy(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				//fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			//fmt.Println(err.Error())
		}
		return nil
	}

	return result

}

func GetBucketLocation(sess *session.Session, bucket string) string {
	svc := s3.New(sess)
	input := &s3.GetBucketLocationInput{
		Bucket: aws.String(bucket),
	}

	result, err := svc.GetBucketLocation(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				//fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			//fmt.Println(err.Error())
		}
		return ""
	}

	if result.LocationConstraint != nil {
		return *result.LocationConstraint
	}
	return ""
}

func GetBucketReplication(sess *session.Session, bucket string) string {
	svc := s3.New(sess)
	input := &s3.GetBucketReplicationInput{
		Bucket: aws.String(bucket),
	}

	_, err := svc.GetBucketReplication(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return aerr.Code()
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return aerr.Code()
		}
	}
	status := "exists"
	return status
}

func GetBucketACL(sess *session.Session, bucket string) *s3.GetBucketAclOutput {
	svc := s3.New(sess)
	input := &s3.GetBucketAclInput{
		Bucket: aws.String(bucket),
	}

	result, err := svc.GetBucketAcl(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				//fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			//fmt.Println(err.Error())
		}
		return nil
	}

	return result
}

func GetBucketWebsite(sess *session.Session, bucket string) *s3.GetBucketWebsiteOutput {
	svc := s3.New(sess)
	input := &s3.GetBucketWebsiteInput{
		Bucket: aws.String(bucket),
	}

	result, err := svc.GetBucketWebsite(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				//fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			//fmt.Println(err.Error())
		}
		return nil
	}

	return result
}

func GetBucketVersioning(sess *session.Session, bucket string) {
	svc := s3.New(sess)
	input := &s3.GetBucketVersioningInput{
		Bucket: aws.String(bucket),
	}

	result, err := svc.GetBucketVersioning(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

func CopyObject(sess *session.Session, bucket string, key string, kmsArn string) *s3.CopyObjectOutput {
	svc := s3.New(sess)
	input := &s3.CopyObjectInput{
		Bucket:               aws.String(bucket),
		CopySource:           aws.String("/" + bucket + "/" + key),
		Key:                  aws.String(key),
		SSEKMSKeyId:          aws.String(kmsArn),
		ServerSideEncryption: aws.String("aws:kms"),
	}

	result, err := svc.CopyObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeObjectNotInActiveTierError:
				fmt.Println(s3.ErrCodeObjectNotInActiveTierError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}

	return result
}

func ListObjects(sess *session.Session, bucket string) *s3.ListObjectsV2Output {
	svc := s3.New(sess)
	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		MaxKeys: aws.Int64(1000),
	}

	result, err := svc.ListObjectsV2(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				fmt.Println(s3.ErrCodeNoSuchBucket, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}

	return result
}
