package workflow

import (
	"github.com/pkg/errors"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"

)

// RequestHandler represents a function that will handle incoming request
type RequestHandler func (models.RelActionPath, models.NamespacePayload4) (furtherActions []TriggerImmediateEventForAnotherUser, err error)

// Routes is a mapping by head of action path to request handlers
type Routes map[string]RequestHandler


// NamedTemplate is a pair of name and workflow template.
// A list of NamedTemplate s can be converted to routing table.
type NamedTemplate struct {
	Name string
	Template
}

// Handler creates a RequestHandler from this set of routes.
func (routes Routes)Handler()RequestHandler {
	return func (actionPath models.RelActionPath, np models.NamespacePayload4) (furtherActions []TriggerImmediateEventForAnotherUser, err error) {
		head, tail := actionPath.HeadTail()
		r, ok := routes[head]
		if ok {
			// var furtherActions2 []TriggerImmediateEventForAnotherUser
			furtherActions, err = r(tail, np)
			// for _, fa := range furtherActions2 {
			// 	furtherActions = append(furtherActions, TriggerImmediateEventForAnotherUser{
			// 		UserID: fa.UserID, 
			// 	})
			// }
		} else {
			err = errors.New("No route found for " + actionPath.Encode() + ". Known keys: " + routes.KnownKeys())
		}
		return
	}
}
// KnownKeys returns comma-separated list of keys
func (routes Routes)KnownKeys() (res string) {
	for key := range routes {
		res = res + key + ", "
	}
	return
}

// ActionPathFromCallbackID - parses action path. Use NamespacePayload4.GetActionPath
func ActionPathFromCallbackID(np models.NamespacePayload4) models.ActionPath {
	return models.ParseActionPath(np.InteractionCallback.CallbackID)
}

// ToRoutingTable converts list of named templates to routing table
func ToRoutingTable(prefix models.Path, envWithoutPrefix Environment, namedTemplates []NamedTemplate) (routes Routes) {
	routes = make(Routes)
	for _, nt := range namedTemplates {
		env := envWithoutPrefix
		env.Prefix = append(prefix, nt.Name)
		routes[nt.Name] = nt.Template.GetRequestHandler(env)
	}
	return
}
