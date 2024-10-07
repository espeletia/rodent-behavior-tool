package setup

import (
	"context"
	"echoes/internal/config"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const defaultUrl = "https://fra1.digitaloceanspaces.com"

func SetupS3Client(config config.S3Config, httpClient *http.Client) (*s3.Client, error) {
	options := []func(*awsConfig.LoadOptions) error{
		awsConfig.WithSharedConfigFiles([]string{config.Credentials}),
		awsConfig.WithHTTPClient(httpClient),
		awsConfig.WithRegion(config.Region),
	}

	if config.AccessKey != "" && config.SecretAccessKey != "" {
		creds := awsConfig.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     config.AccessKey,
				SecretAccessKey: config.SecretAccessKey,
				CanExpire:       false,
			}, nil
		}))
		options = append(options, creds)
	}

	if config.URL != "" && config.URL != defaultUrl {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if config.URL != "" {
				return aws.Endpoint{
					URL: config.URL,
				}, nil
			}

			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})
		options = append(options, awsConfig.WithEndpointResolverWithOptions(customResolver))
	}

	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(), options...)
	if err != nil {
		return nil, err
	}

	awsClient := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return awsClient, err
}

