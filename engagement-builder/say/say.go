// Package say contains mechanisms to allow user-facing template rendering
// thoughout the whole application.
package say

// Identifiers.
// Identifiers are designed in such a way that if there is no template,
// we can degrade to simply return the identifier to user and the 
// meaning will be communicated.
const (
	ThankYou        ResourceKey = "Thank You"
	LearnMore       ResourceKey = "Learn more..."
	Edit            ResourceKey = "Edit"
	Delete          ResourceKey = "Delete"
	Yes             ResourceKey = "Yes"
	No              ResourceKey = "No"
	Cancel          ResourceKey = "Cancel"
)
