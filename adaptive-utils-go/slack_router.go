package adaptive_utils_go

import (
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/daos/common"
	"github.com/slack-go/slack"
)

// RequestHandler represents a function that handles requests from Slack and produces 
// some Slack notifications.
type RequestHandler func (slack.InteractionCallback, common.DynamoDBConnection) ([]models.PlatformSimpleNotification, error)

// DialogSubmissionHandler represents a function that handles dialog submission requests from Slack and produces 
// some Slack notifications.
type DialogSubmissionHandler func (slack.InteractionCallback, slack.DialogSubmissionCallback, common.DynamoDBConnection) ([]models.PlatformSimpleNotification, error)

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
	
	return func (slack.InteractionCallback, common.DynamoDBConnection) ([]models.PlatformSimpleNotification, error) {
		return []models.PlatformSimpleNotification{}, nil
	}
}

// Dispatch selects an appropriate handler from map of handlers or returns a no-op handler.
func (r DialogSubmissionHandlers)Dispatch(route string) DialogSubmissionHandler {
	handler, ok := r[route]
	if ok { return handler }
	return func (slack.InteractionCallback, slack.DialogSubmissionCallback, common.DynamoDBConnection) ([]models.PlatformSimpleNotification, error) {
		return []models.PlatformSimpleNotification{}, nil 
	}
}

// DispatchByRule selects an appropriate handler from map of handlers
// based on the value extracted from request itself
// or returns a no-op handler.
func (r RequestHandlers)DispatchByRule(rule RequestRoutingRule) RequestHandler {
	return func (request slack.InteractionCallback, conn common.DynamoDBConnection) ([]models.PlatformSimpleNotification, error) {
		route := rule(request)
		handler, ok := r[route]
		if ok { return handler(request, conn) }
		return []models.PlatformSimpleNotification{}, nil
	}
}

// DispatchByRule selects an appropriate handler from map of handlers
// based on the value extracted from request itself
// or returns a no-op handler.
func (r DialogSubmissionHandlers)DispatchByRule(rule RequestRoutingRule) DialogSubmissionHandler {
	return func (request slack.InteractionCallback, dialog slack.DialogSubmissionCallback, conn common.DynamoDBConnection) ([]models.PlatformSimpleNotification, error) {
		route := rule(request)
		handler, ok := r[route]
		if ok { return handler(request, dialog, conn) }
		return []models.PlatformSimpleNotification{}, nil 
	}
}

// RunAlso runs two handlers and combines their results.
func (rh RequestHandler)RunAlso(rh2 RequestHandler) RequestHandler {
	return func (request slack.InteractionCallback, conn common.DynamoDBConnection) (notes []models.PlatformSimpleNotification, err error) {
		notes, err = rh(request, conn)
		if err == nil {
			var notes2 []models.PlatformSimpleNotification
			notes2, err = rh2(request, conn)
			notes = append(notes, notes2...)
		}
		return
	}
}

// RunAlso runs two handlers and combines their results.
func (dsh DialogSubmissionHandler)RunAlso(dsh2 DialogSubmissionHandler) DialogSubmissionHandler {
	return func (request slack.InteractionCallback, dialog slack.DialogSubmissionCallback, conn common.DynamoDBConnection) (notes []models.PlatformSimpleNotification, err error) {
		notes, err = dsh(request, dialog, conn)
		if err == nil {
			var notes2 []models.PlatformSimpleNotification
			notes2, err = dsh2(request, dialog, conn)
			notes = append(notes, notes2...)
		}
		return
	}
}
