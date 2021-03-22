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

func GetBucketLocation(sess *session.Session, bucket string) {
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
		return
	}

	fmt.Println(*result)
}

func GetBucketReplication(sess *session.Session, bucket string) *s3.GetBucketReplicationOutput {
	svc := s3.New(sess)
	input := &s3.GetBucketReplicationInput{
		Bucket: aws.String(bucket),
	}

	result, err := svc.GetBucketReplication(input)
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
