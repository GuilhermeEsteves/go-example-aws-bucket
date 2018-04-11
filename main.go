package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	region      = "us-east-1"
	ourBucket   = "our-bucket"
	theirBucket = "their-bucket"
)

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	svc := s3.New(sess)

	resp, err := svc.ListObjects(&s3.ListObjectsInput{Bucket: aws.String(theirBucket)})

	if err != nil {
		log.Fatal("Error in list", err)
	}

	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		sendItemAnotherBucket(svc, item)
		deleteFileBucket(svc, item)
		fmt.Println("\n*******Content********")
		getObject(svc, item)
		fmt.Println("*******End Content********")
		break
	}
}

func getObject(svc *s3.S3, item *s3.Object) {
	out, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(ourBucket),
		Key:    item.Key,
	})

	if err != nil {
		fmt.Println("Error in GetObject", err)
		return
	}

	content, _ := ioutil.ReadAll(out.Body)
	fmt.Println(string(content))
}

func sendItemAnotherBucket(svc *s3.S3, item *s3.Object) {
	source := theirBucket + "/" + *item.Key
	_, err := svc.CopyObject(&s3.CopyObjectInput{Bucket: aws.String(ourBucket), CopySource: aws.String(source), Key: aws.String(*item.Key)})
	if err != nil {
		fmt.Println("Unable to copy item from bucket", err)
	}

	err = svc.WaitUntilObjectExists(&s3.HeadObjectInput{Bucket: aws.String(ourBucket), Key: aws.String(*item.Key)})
	if err != nil {
		fmt.Println("Error occurred while waiting for item  to be copied to bucket")
	}

	fmt.Println("File copied with success!")
}

func deleteFileBucket(svc *s3.S3, item *s3.Object) {
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(theirBucket), Key: aws.String(*item.Key)})
	if err != nil {
		fmt.Println("Unable to delete object from bucket")
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(theirBucket),
		Key:    aws.String(*item.Key),
	})

	if err != nil {
		fmt.Printf("Error occurred while waiting for object %q to be deleted\n", *item.Key)
	}

	fmt.Printf("File deleted with sucess!")
}
