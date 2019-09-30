package aws

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	ctxaws "multipart-upload-to-s3-using-presign-url/server/context/aws"

	"github.com/aws/aws-sdk-go/service/s3"
)

//AWSService is AWS service context
type AWSService struct{}

var region = "us-east-1"

// Get bucket name by bucket type
func getBucketNameByType(bucketType string) (bucketName string) {
	switch bucketType {
	case "image":
		bucketName = "videomine-images"
	case "video":
		bucketName = "videomine-videos"
	default:
		bucketName = "videomine-test"
	}
	return
}

// StartMultipartUpload initiates a multipart upload and returns an upload ID.
func (a *AWSService) StartMultipartUpload(r *http.Request, args *ctxaws.StartMultipartUploadArgs, reply *ctxaws.StartMultipartUploadResp) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	// Create S3 service client
	svc := s3.New(sess)
	mediaID := strconv.Itoa(rand.Intn(10))
	input := &s3.CreateMultipartUploadInput{
		Bucket:      aws.String(getBucketNameByType(args.BucketType)),
		Key:         aws.String(mediaID),
		ContentType: aws.String(args.ContentType),
	}

	resp, err := svc.CreateMultipartUpload(input)

	if err != nil {
		fmt.Println(err.Error())
		*reply = ctxaws.StartMultipartUploadResp{UploadID: "", Error: err.Error(), MediaID: ""}
	} else {
		*reply = ctxaws.StartMultipartUploadResp{UploadID: *resp.UploadId, Error: "", MediaID: mediaID}
	}
	return nil
}

// CompleteMultipartUpload Completes a multipart upload by assembling previously uploaded parts.
func (a *AWSService) CompleteMultipartUpload(r *http.Request, args *ctxaws.CompleteMultipartUploadArgs, reply *ctxaws.CompleteMultipartUploadResp) error {
	var parts []*s3.CompletedPart

	for _, part := range args.Parts {
		parts = append(parts, &s3.CompletedPart{
			ETag:       aws.String(part.ETag),
			PartNumber: aws.Int64(part.PartNumber),
		})
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	// Create S3 service client
	svc := s3.New(sess)

	completeInput := &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(getBucketNameByType(args.BucketType)),
		Key:      aws.String(args.MediaID),
		UploadId: aws.String(args.UploadID),
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: parts,
		},
	}
	_, err = svc.CompleteMultipartUpload(completeInput)

	if err != nil {
		fmt.Println(err.Error())
		*reply = ctxaws.CompleteMultipartUploadResp{Status: false, Error: err.Error()}
	} else {
		*reply = ctxaws.CompleteMultipartUploadResp{Status: true, Error: ""}
	}
	return nil
}

// GetMultipartUploadURL return aws presign url for given uploading part.
func (a *AWSService) GetMultipartUploadURL(r *http.Request, args *ctxaws.GetMultipartUploadURLArgs, reply *ctxaws.GetMultipartUploadURLResp) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	// Create S3 service client
	svc := s3.New(sess)
	req, _ := svc.UploadPartRequest(&s3.UploadPartInput{
		Bucket:     aws.String(getBucketNameByType(args.BucketType)),
		Key:        aws.String(args.MediaID),
		PartNumber: aws.Int64(args.PartNumber),
		UploadId:   aws.String(args.UploadID),
	})

	urlStr, err := req.Presign(30 * time.Minute)

	if err != nil {
		fmt.Println(err.Error())
		*reply = ctxaws.GetMultipartUploadURLResp{URL: "", Error: err.Error()}
	} else {
		*reply = ctxaws.GetMultipartUploadURLResp{URL: urlStr, Error: ""}
	}
	return nil
}

// GetUploadURL return aws presign url to direct upload on s3.
func (a *AWSService) GetUploadURL(r *http.Request, args *ctxaws.GetUploadURLArgs, reply *ctxaws.GetUploadURLResp) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	// Create S3 service client
	svc := s3.New(sess)
	// mediaID := strconv.Itoa(rand.Intn(10))

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(getBucketNameByType(args.BucketType)),
		Key:    aws.String("ijaj2.jpg"),
	})

	urlStr, err := req.Presign(30 * time.Minute)

	if err != nil {
		fmt.Println(err.Error())
		*reply = ctxaws.GetUploadURLResp{URL: "", Error: err.Error()}
	} else {
		*reply = ctxaws.GetUploadURLResp{URL: urlStr, Error: ""}
	}
	return nil
}
