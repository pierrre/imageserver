// Package s3 provides a imageserver.Server implementation that gets the Image from an S3 URL.
package s3

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/pierrre/imageserver"
)

// Server is a imageserver.Server implementation that gets the Image from an S3 URL.
//
// It parses the "source" param as a bucket-relative URL, then does a HEAD followed by
// a GET request for the item.
// It returns an error if there are any problems.
type Server struct {
	Session    *aws.Config
	BucketName string
}

// NewS3Server takes AWS Region, AccessKey, SecretKey, and S3 Bucket name, and returns a Server or an error.
//
// If the values for Region, AccessKey, or SecretKey are not provided they are read from the well-known
// environment variable names they are typically stored.
func NewS3Server(awsRegion, awsAccessKey, awsSecretKey, awsS3Bucket string) (imageserver.Server, error) {
	if awsS3Bucket == "" {
		return nil, fmt.Errorf("the s3 bucket must be specified")
	}

	awsSession, err := newAWSSession(awsRegion, awsAccessKey, awsSecretKey)
	if err != nil {
		return nil, err
	}

	return &Server{
		Session:    awsSession,
		BucketName: awsS3Bucket,
	}, nil
}

func newAWSSession(awsRegion, awsAccessKey, awsSecretKey string) (*aws.Config, error) {
	config := aws.NewConfig()

	// Region
	if awsRegion != "" {
		// CLI trumps
		config.Region = awsRegion
	} else if os.Getenv("AWS_DEFAULT_REGION") != "" {
		// Env is good, too
		config.Region = os.Getenv("AWS_DEFAULT_REGION")
	} else {
		return nil, fmt.Errorf("cannot find AWS region")
	}

	// Creds
	if awsAccessKey != "" && awsSecretKey != "" {
		// CLI trumps
		config.Credentials = credentials.NewStaticCredentialsProvider(
			awsAccessKey,
			awsSecretKey,
			"")
	} else if os.Getenv("AWS_ACCESS_KEY_ID") != "" {
		// Env is good, too
		config.Credentials = credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"")
	}
	return config, nil
}

// Get implements imageserver.Server.
func (srv *Server) Get(params imageserver.Params) (*imageserver.Image, error) {
	bucketPath, err := params.GetString("source")
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(*srv.Session)

	// HEAD request
	hoo, err := client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: awsString(srv.BucketName),
		Key:    awsString(bucketPath),
	})
	if err != nil {
		return nil, err
	}

	// Determine the file type
	format := identifyFormat(*hoo.ContentType) // TODO this right
	if format == "" {
		return nil, fmt.Errorf("content-type of '%s' not a valid image type", *hoo.ContentType)
	}

	// pre-allocate in memory buffer, where headObject type is *s3.HeadObjectOutput
	buf := make([]byte, int(hoo.ContentLength))
	// wrap with aws.WriteAtBuffer
	w := manager.NewWriteAtBuffer(buf)

	// GET it
	downloader := manager.NewDownloader(client)
	if _, err := downloader.Download(context.TODO(), w,
		&s3.GetObjectInput{
			Bucket: awsString(srv.BucketName),
			Key:    awsString(bucketPath),
		}); err != nil {
		return nil, err
	}

	// Return the Image
	return &imageserver.Image{
		Format: format,
		Data:   w.Bytes(),
	}, nil
}

// awsString mimics aws.String() as many AWS SDK functions
// demand *string and not string :shrug:
func awsString(v string) *string {
	return &v
}

// identifyFormat returns the right side of an "image/" content-type string,
// or empty
func identifyFormat(contentType string) string {
	if contentType == "" {
		return ""
	} else if !strings.HasPrefix(contentType, "image/") {
		return ""
	}

	return strings.TrimPrefix(contentType, "image/")
}
