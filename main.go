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
	region = "us-east-1"
	bucket = "dynamo-price-research-hml"
)

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	svc := s3.New(sess)

	resp, err := svc.ListObjects(&s3.ListObjectsInput{Bucket: aws.String(bucket)})

	if err != nil {
		log.Fatal("Error in list", err)
	}

	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("key", *item.Key)
		fmt.Println("*******Content********")
		getObject(svc, item)
		fmt.Println("*******End Content********")
	}
}

func getObject(svc *s3.S3, item *s3.Object) {
	out, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    item.Key,
	})

	if err != nil {
		fmt.Println("Error in GetObject", err)
		return
	}

	content, _ := ioutil.ReadAll(out.Body)
	fmt.Println(string(content))
}
