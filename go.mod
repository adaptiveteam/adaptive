module github.com/adaptiveteam/adaptive

require (
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/ReneKroon/ttlcache v1.5.0
	github.com/adaptiveteam/adaptive-utils-go v0.16.0
	github.com/adaptiveteam/dialog-fetcher v0.4.0
	github.com/antlr/antlr4 v0.0.0-20190910151933-cd81586d3d6a // indirect
	github.com/avast/retry-go v2.4.2+incompatible
	github.com/aws/aws-dax-go v1.1.1
	github.com/aws/aws-lambda-go v1.13.1
	github.com/aws/aws-sdk-go v1.23.19
	github.com/containerd/continuity v0.0.0-20190827140505-75bee3e2ccb6 // indirect

	github.com/go-test/deep v1.0.3
	github.com/google/go-cmp v0.3.1 // indirect

	github.com/google/uuid v1.1.1
	github.com/mattn/go-colorable v0.1.2
	github.com/mattn/go-isatty v0.0.9 // indirect
	github.com/nlopes/slack v0.6.0
	github.com/onsi/ginkgo v1.10.1
	github.com/onsi/gomega v1.7.0
	github.com/ory/dockertest v3.3.5+incompatible
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20190909003024-a7b16738d86b // indirect
	golang.org/x/sys v0.0.0-20190910064555-bbd175535a8b // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gotest.tools v2.2.0+incompatible
)

replace github.com/nlopes/slack => github.com/adaptiveteam/slack v0.10.0

replace gopkg.in/urfave/cli.v1 => gopkg.in/urfave/cli.v1 v1.20.0
