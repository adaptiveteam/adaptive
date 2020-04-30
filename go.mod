module github.com/adaptiveteam/adaptive

go 1.14

require (
	github.com/360EntSecGroup-Skylar/excelize/v2 v2.0.2
	github.com/Merovius/go-misc v0.0.0-20190816102157-30fd8989cac1
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/ReneKroon/ttlcache v1.6.0
	github.com/acroca/go-symbols v0.1.1 // indirect
	github.com/antlr/antlr4 v0.0.0-20200103163232-691acdc23f1f // indirect
	github.com/avast/retry-go v2.4.2+incompatible
	github.com/aws/aws-dax-go v1.1.2
	github.com/aws/aws-lambda-go v1.13.2
	github.com/aws/aws-sdk-go v1.27.1
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/containerd/continuity v0.0.0-20200107194136-26c1120b8d41 // indirect
	github.com/cweill/gotests v1.5.3 // indirect
	github.com/go-delve/delve v1.4.0 // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-test/deep v1.0.4
	github.com/gofrs/uuid v3.2.0+incompatible // indirect
	github.com/gonum/blas v0.0.0-20181208220705-f22b278b28ac // indirect
	github.com/gonum/floats v0.0.0-20181209220543-c233463c7e82 // indirect
	github.com/gonum/integrate v0.0.0-20181209220457-a422b5c0fdf2 // indirect
	github.com/gonum/internal v0.0.0-20181124074243-f884aa714029 // indirect
	github.com/gonum/lapack v0.0.0-20181123203213-e4cdc5a0bff9 // indirect
	github.com/gonum/matrix v0.0.0-20181209220409-c518dec07be9 // indirect
	github.com/gonum/stat v0.0.0-20181125101827-41a0da705a5b
	github.com/google/uuid v1.1.1
	github.com/gotestyourself/gotestyourself v2.2.0+incompatible // indirect
	github.com/gruntwork-io/terratest v0.23.4
	github.com/jinzhu/gorm v1.9.12
	github.com/karrick/godirwalk v1.15.6 // indirect
	github.com/lib/pq v1.3.0 // indirect
	github.com/mattn/go-colorable v0.1.4
	github.com/mdempsky/gocode v0.0.0-20200405233807-4acdcbdea79d // indirect
	github.com/nlopes/slack v0.6.0
	github.com/onsi/ginkgo v1.10.2
	github.com/onsi/gomega v1.7.0
	github.com/opencontainers/runc v0.1.1
	github.com/ory/dockertest v3.3.5+incompatible
	github.com/pkg/errors v0.9.0
	github.com/rogpeppe/godef v1.1.2 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/sqs/goreturns v0.0.0-20181028201513-538ac6014518 // indirect
	github.com/stamblerre/gocode v1.0.0 // indirect
	github.com/stretchr/testify v1.4.0
	github.com/thoas/go-funk v0.5.0
	github.com/unidoc/unidoc v0.0.0-20190420234413-f1b8d1021126
	github.com/unidoc/unipdf/v3 v3.1.0
	go.uber.org/zap v1.10.0
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	gotest.tools v2.2.0+incompatible
	rsc.io/quote/v3 v3.1.0 // indirect
	rsc.io/sampler v1.99.99 // indirect
)

replace github.com/nlopes/slack => github.com/adaptiveteam/slack v0.13.0

replace gopkg.in/urfave/cli.v1 => gopkg.in/urfave/cli.v1 v1.20.0

replace github.com/aws/aws-sdk-go => github.com/aws/aws-sdk-go v1.25.16
