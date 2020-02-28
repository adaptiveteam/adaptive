module github.com/adaptiveteam/adaptive

go 1.13

require (
	github.com/360EntSecGroup-Skylar/excelize/v2 v2.0.2
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/ReneKroon/ttlcache v1.6.0
	github.com/antlr/antlr4 v0.0.0-20200103163232-691acdc23f1f // indirect
	github.com/avast/retry-go v2.4.2+incompatible
	github.com/aws/aws-dax-go v1.1.2
	github.com/aws/aws-lambda-go v1.13.2
	github.com/aws/aws-sdk-go v1.27.1
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/containerd/continuity v0.0.0-20200107194136-26c1120b8d41 // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-test/deep v1.0.4
	github.com/gofrs/uuid v3.2.0+incompatible // indirect
	github.com/google/uuid v1.1.1
	github.com/gruntwork-io/terratest v0.23.4
	github.com/mattn/go-colorable v0.1.4
	github.com/nlopes/slack v0.6.0
	github.com/onsi/ginkgo v1.10.2
	github.com/onsi/gomega v1.7.0
	github.com/opencontainers/runc v0.1.1
	github.com/ory/dockertest v3.3.5+incompatible
	github.com/pkg/errors v0.9.0
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.4.0
	github.com/thoas/go-funk v0.5.0
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	gotest.tools v2.2.0+incompatible
)

replace github.com/nlopes/slack => github.com/adaptiveteam/slack v0.13.0

replace gopkg.in/urfave/cli.v1 => gopkg.in/urfave/cli.v1 v1.20.0

replace github.com/aws/aws-sdk-go => github.com/aws/aws-sdk-go v1.25.16
