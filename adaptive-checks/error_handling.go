package adaptive_checks

import (
	core "github.com/adaptiveteam/adaptive/core-utils-go"
)

// recoverToErrorVar recovers and places the recovered error into the given variable
var recoverToErrorVar = core.RecoverToErrorVar

// RecoverToLog recovers and places the recovered error into the given variable
var RecoverToLog = core.RecoverAsLogError
