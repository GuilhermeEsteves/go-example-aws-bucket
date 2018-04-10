package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("sa-east-1")},
	)

	svc := s3.New(sess)

	resp, err := svc.ListObjects(&s3.ListObjectsInput{Bucket: aws.String("teste-leitura")})

	if err != nil {
		log.Fatal("Error ao listar")
	}
	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("key", *item.Key)
	}
}
