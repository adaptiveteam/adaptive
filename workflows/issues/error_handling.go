package issues

import (
	core "github.com/adaptiveteam/adaptive/core-utils-go"
)

// recoverToErrorVar recovers and places the recovered error into the given variable
var recoverToErrorVar = core.RecoverToErrorVar

func (w workflowImpl) recoverToErrorVar(name string, err *error) {
	if err != nil {
		w.AdaptiveLogger.WithError(*err).Errorln("Before recoverToErrorVar " + name)
	}
	core.RecoverToErrorVar(name, err)
	if err != nil {
		w.AdaptiveLogger.WithError(*err).Errorln("After recoverToErrorVar " + name)
	}

}

// RecoverToLog recovers and places the recovered error into the given variable
var RecoverToLog = core.RecoverAsLogError
