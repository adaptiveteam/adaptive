package aws_utils_go

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
)

// AWS Session
func sess(region, endpoint string) (*session.Session, *aws.Config) {
	var sess *session.Session
	var conf *aws.Config

	if endpoint != "" {
		sess = session.Must(session.NewSession(&aws.Config{
			Credentials:      credentials.NewStaticCredentials("foo", "var", ""),
			S3ForcePathStyle: aws.Bool(true),
			Region:           aws.String(endpoints.UsEast1RegionID),
			DisableSSL:       aws.Bool(true),
		}))
		conf = &aws.Config{
			Endpoint: aws.String(endpoint),
		}
	} else {
		sess = session.Must(session.NewSession(&aws.Config{
			Region: aws.String(region),
		}))
		conf = &aws.Config{}
	}
	return sess, conf
}
