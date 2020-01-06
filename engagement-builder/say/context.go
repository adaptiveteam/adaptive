// Package say contains mechanisms to allow user-facing template rendering
// thoughout the whole application.
package say

import (
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

// ResourceKey is the type of the key for accessing templates
type ResourceKey string

// ResourceContext provides mechanisms to obtain good user-facing
// templates.
type ResourceContext interface {
	// SayRich returns RichText template.
	SayRich(id ResourceKey) ui.RichText
	// SayPlain returns PlainText template.
	SayPlain(id ResourceKey) ui.PlainText
	// HasID checks if there is such key
	HasID(id ResourceKey) bool
	// HasPlainID checks if there is such key that can be used for plain texts
	HasPlainID(id ResourceKey) bool
}

// EmptyResourceContextImpl is a fake resource context that returns just the given id.
type EmptyResourceContextImpl struct {}

var (
	// EmptyResourceContext is an instance of empty resource context
	EmptyResourceContext  ResourceContext = EmptyResourceContextImpl{}
	// globalResourceContext is used by Say function to resolve
	globalResourceContext                 = EmptyResourceContext
)

// SayRich returns id ascribed to RichText
func (e EmptyResourceContextImpl) SayRich(id ResourceKey) ui.RichText {
	return ui.RichText(id)
}

// SayPlain returns id ascribed to PlainText
func (e EmptyResourceContextImpl) SayPlain(id ResourceKey) ui.PlainText {
	return ui.PlainText(id)
}

// HasID returns true
func (e EmptyResourceContextImpl) HasID(id ResourceKey) bool {
	return true
}

// HasPlainID returns true
func (e EmptyResourceContextImpl) HasPlainID(id ResourceKey) bool {
	return true
}


// UnsafeUpdateGlobalContext changes global context to also include the provided one.
func UnsafeUpdateGlobalContext(context ResourceContext) {
	globalResourceContext = combinedResourceContext{context, globalResourceContext}
}

type combinedResourceContext []ResourceContext

// SayRich returns RichText template.
func (contexts combinedResourceContext) SayRich(id ResourceKey) ui.RichText {
	for _, ctx := range contexts {
		if ctx.HasID(id) {
			return ctx.SayRich(id)
		}
	}
	return ui.RichText(id)
}

// SayPlain returns PlainText template.
func (contexts combinedResourceContext) SayPlain(id ResourceKey) ui.PlainText {
	for _, ctx := range contexts {
		if ctx.HasID(id) {
			return ctx.SayPlain(id)
		}
	}
	return ui.PlainText(id)
}

// HasID returns true if any of the contexts have the given id.
func (contexts combinedResourceContext) HasID(id ResourceKey) bool {
	for _, ctx := range contexts {
		if ctx.HasID(id) {
			return true
		}
	}
	return false
}

// HasPlainID returns true if any of the contexts have the given id suitable for plain text rendering.
func (contexts combinedResourceContext) HasPlainID(id ResourceKey) bool {
	for _, ctx := range contexts {
		if ctx.HasPlainID(id) {
			return true
		}
	}
	return false
}

// PlainTextResources is just a map of keys to plain text.
// the same plain text can be converted to RichText.
type PlainTextResources map[ResourceKey]ui.PlainText

// SayRich returns id ascribed to RichText
func (m PlainTextResources) SayRich(id ResourceKey) ui.RichText {
	return ui.RichText(m[id])
}

// SayPlain returns id ascribed to PlainText
func (m PlainTextResources) SayPlain(id ResourceKey) ui.PlainText {
	return ui.PlainText(m[id])
}

// HasID returns true
func (m PlainTextResources) HasID(id ResourceKey) bool {
	_, ok := m[id]
	return ok
}

// HasPlainID returns true
func (m PlainTextResources) HasPlainID(id ResourceKey) bool {
	_, ok := m[id]
	return ok
}

// RichTextResources is just a map of keys to rich text.
type RichTextResources map[ResourceKey]ui.RichText

// SayRich returns template for the given id
func (m RichTextResources) SayRich(id ResourceKey) ui.RichText {
	return ui.RichText(m[id])
}

// SayPlain returns empty string
func (m RichTextResources) SayPlain(id ResourceKey) ui.PlainText {
	return ""
}

// HasID returns true for available resources
func (m RichTextResources) HasID(id ResourceKey) bool {
	_, ok := m[id]
	return ok
}

// HasPlainID returns false because this collection doesn't have plain resources
func (m RichTextResources) HasPlainID(id ResourceKey) bool {
	return false
}



// Rich uses global context to return the RichText template
func Rich(id ResourceKey) ui.RichText {
	return globalResourceContext.SayRich(id)
}

// Plain uses global context to return the PlainText template 
func Plain(id ResourceKey) ui.PlainText {
	return globalResourceContext.SayPlain(id)
}

