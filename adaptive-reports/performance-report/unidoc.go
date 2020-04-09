package collaboration_report

import (
	"log"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/adaptiveteam/adaptive/aws-utils-go"
)

// UniDocLicenseKeySecretName - the name of unidoc license key
const UniDocLicenseKeySecretName = "dev/unidoc.license.key"

// SetUniDocGlobalLicenseIfAvailable reads license and sets it 
func SetUniDocGlobalLicenseIfAvailable() {
	sm := aws_utils_go.GetSecretsManagerFromEnv()
	uniDocLicense, err2 := sm.ReadSecretString(UniDocLicenseKeySecretName)
	if err2 == nil {
		err2 = license.SetLicenseKey(
			uniDocLicense,
			"adaptive.team",
		)
	}
	if err2 != nil {
		log.Printf("IGNORING ERROR in SetUniDocGlobalLicenseIfAvailable: %+v\n", err2)
	}
}

// NB: We cannot use `init` because many lambdas don't have access to SM
// func init() {
// 	SetUniDocGlobalLicenseIfAvailable()
// }
