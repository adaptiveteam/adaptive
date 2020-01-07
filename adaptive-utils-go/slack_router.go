package adaptive_utils_go

import (
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/nlopes/slack"
)

// RequestHandler represents a function that handles requests from Slack and produces 
// some Slack notifications.
type RequestHandler func (slack.InteractionCallback) []models.PlatformSimpleNotification 

// DialogSubmissionHandler represents a function that handles dialog submission requests from Slack and produces 
// some Slack notifications.
type DialogSubmissionHandler func (slack.InteractionCallback, slack.DialogSubmissionCallback) []models.PlatformSimpleNotification 

// RequestRoutingRule extracts route from request
type RequestRoutingRule func (slack.InteractionCallback) string

// RequestHandlers is a routing map for requests
type RequestHandlers map[string]RequestHandler

// DialogSubmissionHandlers is a routing map for dialog submissions
type DialogSubmissionHandlers map[string]DialogSubmissionHandler

// Dispatch selects an appropriate handler from map of handlers or returns a no-op handler.
func (r RequestHandlers)Dispatch(route string) RequestHandler {
	handler, ok := r[route]
	if ok { return handler }
	
	return func (slack.InteractionCallback) []models.PlatformSimpleNotification {
		return []models.PlatformSimpleNotification{} 
	}
}

// Dispatch selects an appropriate handler from map of handlers or returns a no-op handler.
func (r DialogSubmissionHandlers)Dispatch(route string) DialogSubmissionHandler {
	handler, ok := r[route]
	if ok { return handler }
	return func (slack.InteractionCallback, slack.DialogSubmissionCallback) []models.PlatformSimpleNotification {
		return []models.PlatformSimpleNotification{} 
	}
}

// DispatchByRule selects an appropriate handler from map of handlers
// based on the value extracted from request itself
// or returns a no-op handler.
func (r RequestHandlers)DispatchByRule(rule RequestRoutingRule) RequestHandler {
	return func (request slack.InteractionCallback) []models.PlatformSimpleNotification {
		route := rule(request)
		handler, ok := r[route]
		if ok { return handler(request) }
		return []models.PlatformSimpleNotification{} 
	}
}

// DispatchByRule selects an appropriate handler from map of handlers
// based on the value extracted from request itself
// or returns a no-op handler.
func (r DialogSubmissionHandlers)DispatchByRule(rule RequestRoutingRule) DialogSubmissionHandler {
	return func (request slack.InteractionCallback, dialog slack.DialogSubmissionCallback) []models.PlatformSimpleNotification {
		route := rule(request)
		handler, ok := r[route]
		if ok { return handler(request, dialog) }
		return []models.PlatformSimpleNotification{} 
	}
}

// RunAlso runs two handlers and combines their results.
func (rh RequestHandler)RunAlso(rh2 RequestHandler) RequestHandler {
	return func (request slack.InteractionCallback) [] models.PlatformSimpleNotification {
		return append(rh(request), rh2(request)...)
	}
}

// RunAlso runs two handlers and combines their results.
func (dsh DialogSubmissionHandler)RunAlso(dsh2 DialogSubmissionHandler) DialogSubmissionHandler {
	return func (request slack.InteractionCallback, dialog slack.DialogSubmissionCallback) [] models.PlatformSimpleNotification {
		return append(dsh(request, dialog), dsh2(request, dialog)...)
	}
}
