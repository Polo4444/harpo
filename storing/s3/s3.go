package storing_s3

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/Polo44444/harpo/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	s3config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type s3Provider struct {
	c               *s3.Client
	accessKeyID     string
	secretAccessKey string
	bucket          string
	region          string
	endpoint        string

	// When true, force a path-style endpoint to be used where the bucket name is part
	// of the path.
	//
	// Defaults to false if no value is
	// provided.
	//
	// AWS::S3::ForcePathStyle
	forcePath bool
}

func BuildS3Config(accessKeyID, secretAccessKey, bucket, region, endpoint string, forcePath bool) models.ProviderConfig {

	return models.ProviderConfig{
		"access_key_id":     accessKeyID,
		"secret_access_key": secretAccessKey,
		"bucket":            bucket,
		"region":            region,
		"endpoint":          endpoint,
		"force_path":        forcePath,
	}
}

func NewS3Provider(config models.ProviderConfig) (*s3Provider, error) {

	prvd := &s3Provider{
		accessKeyID:     config["access_key_id"].(string),
		secretAccessKey: config["secret_access_key"].(string),
		bucket:          config["bucket"].(string),
		region:          config["region"].(string),
		endpoint:        config["endpoint"].(string),
		forcePath:       config["force_path"].(bool),
	}

	s3Cfg, err := s3config.LoadDefaultConfig(
		context.TODO(),
		s3config.WithCredentialsProvider(
			credentials.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID: prvd.accessKeyID, SecretAccessKey: prvd.secretAccessKey,
				},
			},
		),
		s3config.WithRegion(prvd.region))

	if err != nil {
		return nil, fmt.Errorf("failed to load s3 config, %w", err)
	}

	prvd.c = s3.NewFromConfig(
		s3Cfg,
		s3.WithEndpointResolverV2(NewResolverV2(prvd.bucket, prvd.region, prvd.endpoint, prvd.forcePath)))

	return prvd, nil
}

func (s *s3Provider) UploadWithReader(ctx context.Context, filePath string, data io.Reader, contentType string) error {

	uploader := manager.NewUploader(s.c)
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(filePath),
		Body:        data,
		ContentType: aws.String(contentType),
		ACL:         types.ObjectCannedACLPrivate,
	})

	return err
}

func (s *s3Provider) DownloadWithWriter(ctx context.Context, filePath string, writer io.Writer) error {

	downloader := manager.NewDownloader(s.c)
	_, err := downloader.Download(ctx, NewWriterAtFromWriter(writer), &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filePath),
	})

	return err
}

func (s *s3Provider) Info(ctx context.Context, filePath string) (*models.FileInfo, error) {

	headObjectOutput, err := s.c.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filePath),
	})
	if err != nil {
		return nil, err
	}

	return &models.FileInfo{
		Size: *headObjectOutput.ContentLength,
	}, nil
}

func (s *s3Provider) Delete(ctx context.Context, filePath string) error {

	_, err := s.c.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filePath),
	})

	return err
}

func (s *s3Provider) DeleteMany(ctx context.Context, filePaths []string) error {

	// build objectIdentifiers
	objects := make([]types.ObjectIdentifier, len(filePaths))
	for i, fp := range filePaths {
		objects[i] = types.ObjectIdentifier{
			Key: aws.String(fp),
		}
	}

	_, err := s.c.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(s.bucket),
		Delete: &types.Delete{
			Objects: objects,
			Quiet:   aws.Bool(true),
		},
	})

	return err
}

// Test tests the connection to the provider.
func (s *s3Provider) Test(ctx context.Context) error {

	testFilePrefix := "harpo.test"
	testFileContent := bytes.NewBufferString("Hi dear! It's Harpo!🤗")

	// Upload
	err := s.UploadWithReader(ctx, testFilePrefix, testFileContent, "text/plain")
	if err != nil {
		return fmt.Errorf("1. Failed to upload test file, %w", err)
	}

	// Download
	buf := new(bytes.Buffer)
	err = s.DownloadWithWriter(ctx, testFilePrefix, buf)
	if err != nil {
		return fmt.Errorf(`1. Upload successfull!
		2. failed to download test file, %w`, err)
	}

	// Check content
	if testFileContent.String() != buf.String() {
		return fmt.Errorf(`1. Upload successfull!
		2. Download successfull!
		3. Test file content is not the same`)
	}

	// Info
	_, err = s.Info(ctx, testFilePrefix)
	if err != nil {
		return fmt.Errorf(`1. Upload successfull!
		2. Download successfull!
		3. Test file content is the same!
		4. Failed to get info of test file, %w`, err)
	}

	// Delete
	err = s.Delete(ctx, testFilePrefix)
	if err != nil {
		return fmt.Errorf(`1. Upload successfull!
		2. Download successfull!
		3. Test file content is the same!
		4. Info successfull!
		5. Failed to delete test file, %w`, err)
	}

	return nil
}

// Close closes the provider.
// S3 provider does not need to be closed.
func (s *s3Provider) Close(ctx context.Context) error {
	return nil
}
