# adaptive-utils-go

[![Build Status](https://travis-ci.com/adaptiveteam/adaptive-utils-go.svg?token=BSM7265i3ndP9kG2qsqY&branch=develop)](https://travis-ci.com/adaptiveteam/adaptive-utils-go)

Common utilities and models that are used across Adaptive modules

## Models

**Engagement** - a structure that represents some future interaction with a particular user.

## Interaction with Slack

Slack has API that allows us to publish messages to chat and create interactive elements (buttons, dialogs, surveys, etc.). Interactive elements send requests to 
our application endpoint. The logic is very similar to that of web-server.

Routes in our platform are represented using `MessageCallback`. It's not a hierarchical path as in URL's. Currently it's a fixed data structure that has an exact number of path parts. (See `./models/chat.go` for definition.) It's field `Module` is often called `Namespace`. We might think of standardizing names.

NB: This callback path has a unique feature - collapsing user engagements based on `(year, month)` pair. It serves this additional role being a primary key for looking for engagements.

When we create an interactive UI element, we serialize `MessageCallback` into colon-separated path. And we parse it back when receiving message from Slack.

So far we don't have a hierarchical routing mechanism. We are allowing all lambdas to handle all `EventsAPIEvent` coming from Slack.

### Filled forms POST (AKA Dialog submissions)

Slack sends us completed forms using `InteractionTypeDialogSubmission` event type.
This event will have `_.CallbackID` which we have previously configured when creating the dialog. The field values from dialog will be available in `_.Submission map[string]string`.

```
if eventsAPIEvent.Type == string(slack.InteractionTypeDialogSubmission) {
```

