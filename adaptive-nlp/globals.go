package nlp

import ("os")

func ensureGlobalConnectionsAreOpen() {
	if globalConnections.Translate == nil {
		globalConnections = OpenConnections(os.Getenv("AWS_REGION"), globalMeaningCloudKey)
	}
}

// globalConnections is a global variable designed to be an Optimization.
// This enables us to create a single AWS client so we aren't creating it over and over.
// TODO: Get rid of this global variable. We should connect inside lambda and pass 
// connections to this library.
var globalConnections Connections
