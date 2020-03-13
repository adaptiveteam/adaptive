package lambda

import (
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	"strings"
)

var ( // const in fact, but golang is not smart enough to figure that out.
	SimulateDialogTitle ui.PlainText = "Simulate"
	SimulateUserLabel   ui.PlainText = "User to simulate"
	SimulateDateLabel   ui.PlainText = "Date to simulate"

	SimulateMenuTitle            ui.PlainText = "Simulate"
	SimulateCurrentQuarterAction              = "simulate_current_quarter"
	SimulateCurrentQuarterText   ui.PlainText = "Simulate a scenario for this quarter"
	SimulateNextQuarterAction                 = "simulate_next_quarter"
	SimulateNextQuarterText      ui.PlainText = "Simulate a scenario for next quarter"

	MenuPrompt         ui.PlainText = "What can I help you with today?"
	MenuPromptFallback ui.PlainText = "Adaptive at your service"

	CommunitiesMenuTitle       ui.PlainText = "Communities"
	CommunitySubscribeAction                = "subscribe"
	CommunitySubscribeText     ui.PlainText = "Be part of a community"
	CommunityUnsubscribeAction              = "unsubscribe"
	CommunityUnsubscribeText   ui.PlainText = "Leave a community"

	CommunityMenuActionName              = "menu_list"
	MainMenuEmbeddedPrompt  ui.PlainText = "Pick an option..."

	// CancelAction models.AttachActionName = "cancel"
	Confirm models.AttachActionName = "confirm"
	No      models.AttachActionName = "no"

	CommentsName                  = "Comments"
	CommentsLabel    ui.PlainText = "Comments"
	PercentDoneLabel ui.PlainText = "Percent Done"

	CommentsSurveyPlaceholder ui.PlainText = ebm.EmptyPlaceholder
	CommentsPlaceholder       ui.PlainText = ebm.EmptyPlaceholder

	CoachingDialogTitle       ui.PlainText = "Coaching"
	CoachingName                           = "coaching"
	CoachingLabel             ui.PlainText = "Coaching"
	CoachRejectionReasonLabel ui.PlainText = "Reason for not accepting the coach"

	AcknowledgeChannelToCommunitySubscription ui.RichText  = "Thank you, this channel is now subscribed to a community"
	ActionCancellationAction                               = "cancel"
	ActionCancellationAcknowledgementText     ui.RichText  = "Action has been cancelled"
	ActionCancellationText                    ui.PlainText = "Not now"

	SelectCommunityMenuText                                  ui.PlainText = "Select"
	SelectCommunityFallbackMenuText                          ui.PlainText = "Select a community"
	UserCommandUnknownText                                   ui.RichText  = "Sorry, could not process your message. Try saying 'hi' or 'hello'."
	UnsubscribedUserCommandRejectText                        ui.RichText  = "Sorry, cannot process your request. You are not subscribed to any community yet."
	UnsubscribedUserAndNoCommunityAvailableCommandRejectText ui.RichText  = ui.RichText("Sorry, I cannot help you. There are no available communities for which you are a coordinator for.").Italics()

	FetchingReportNotification   ui.RichText = "Fetching the report ..."
	GeneratingReportNotification ui.RichText = "Generating the report ..."
	ListOfCoachesWelcomeMessage  ui.RichText = "Hello! Below are the list of coaches available for this quarter"

	UserForReportSelectionPrompt   ui.RichText = "Whose report are you looking for?"
	ScheduleForCurrentQuarterTitle ui.RichText = ui.RichText("Here is the `schedule` for the current quarter :point_down:").Italics()
	ScheduleForNextQuarterTitle    ui.RichText = ui.RichText("Here is the `schedule` for the next quarter :point_down:").Italics()

	PostSubscribeEngagementTitle ui.PlainText = "Which community would you like to be part of?"
	SelectCommunityToLeavePrompt ui.PlainText = "Which community would you like to leave?"

	CoachingRejectionExplanationPrompt ui.PlainText = "Reason for not coaching"
	FarewellNotification               ui.RichText  = ui.RichText("Sorry to see you go. Come back any time.").Italics()
)

var (
	confirm = ebm.AttachmentActionConfirm{
		OkText:      models.YesLabel,
		DismissText: models.CancelLabel,
	}
)

func EngagementsRescheduledFrom(date string) ui.RichText {
	return ui.RichText(fmt.Sprintf("(Engagements rescheduled from %s)", date)).Italics()
}
func SubscribeToCommunityAlreadySubscribedErrorMessage(chanConcat string) ui.RichText {
	return ui.RichText(fmt.Sprintf("Already part of the following communities: `%s`", chanConcat)).Italics()
}
func GeneratingReportForUserNotification(targetUserID string) ui.RichText {
	return ui.RichText(fmt.Sprintf("Generating the report for <@%s>. You can see the status by clicking on the *reply* link below and then looking in the thread here :point_down:", targetUserID)).Italics()
}
func FetchingReportForUserNotification(targetUserID string) ui.RichText {
	return ui.RichText(fmt.Sprintf("Fetching <@%s>'s report for you. Thanks for asking me to help you. You can see the report by clicking on the *reply* link below and then looking in the thread here :point_down:", targetUserID)).Italics()
}
func CoachingRequestTitle(coacheeUserID string) ui.RichText {
	return ui.RichText(fmt.Sprintf(
		"<@%s> is requesting your coaching. Will you be able to coach them?", coacheeUserID)).Italics()
}
func CoachingRequestSentNotificationToCoacheeTitle(coachUserID string) ui.RichText {
	return ui.RichText(fmt.Sprintf(
		"Ok. I will ask <@%s> if they agree to be your coach for this quarter. I'll be back to you soon.", coachUserID)).Italics()
}
func CoachingRequestConfirmationToCoachee(coachUserID string, quarter, year int) ui.RichText {
	return ui.RichText(fmt.Sprintf("Your requested a  coach. <@%s>, has agreed to coach you for quarter `%d`, year `%d`.",
		coachUserID, quarter, year))
}
func CoachingRequestConfirmationToCoach(coacheeUserID string, quarter, year int) ui.RichText {
	return ui.RichText(fmt.Sprintf("Awesome! You will be coaching <@%s> for quarter `%d`, year `%d`.", coacheeUserID, quarter, year))
}

var (
	AdminRequestSentAcknowledgement = ui.RichText("Roger, I will notify Admin about my access with you.").Italics()
)

func NoUserCommunityErrorMessage(userTag ui.RichText) ui.RichText {
	return ui.RichText("There is no Adaptive User Community created yet. Please create it and invite " + userTag + " there or to any Adaptive Community.")
}
func NoUserInviteToUserCommunityErrorMessage(userTag ui.RichText) ui.RichText {
	return ui.RichText("Please invite " + userTag + " to the channel that is subscribed to user community or to any Adaptive Community.")
}
func UserRequestAdaptiveAccessNotifitationToAdminCommunity(userID string, commInfo string) ui.RichText {
	return ui.RichText(fmt.Sprintf("<@%s> is requesting access to Adaptive. %s", userID, commInfo)).Italics()
}

func CoachingRequestConfirmToCoachNotification(coacheeUserID string, quarter, year int) ui.RichText {
	return ui.RichText(fmt.Sprintf("Awesome! You will be coaching <@%s> for quarter `%d`, year `%d`", coacheeUserID, quarter, year)).Italics()
}

func CoachingRequestConfirmToCoacheeNotification(coachUserID string, quarter, year int) ui.RichText {
	return ui.RichText(fmt.Sprintf("Done. <@%s> is your coach for quarter `%d`, year `%d`", coachUserID, quarter, year)).Italics()
}

func CommunityIsUsedErrorMessage(community ui.PlainText) ui.RichText {
	return ui.RichText(fmt.Sprintf("`%s` community is already being used. Please use the command `@adaptive hi` to select another one.", community)).Italics()
}

func CommunitySuccessfullyAddedNotification(community ui.PlainText) ui.RichText {
	return ui.RichText(fmt.Sprintf("I am here to help with `%s` community requests and updates. You can tell me to leave at any time by using the command `/remove @adaptive`.", community)).Italics()
}

func CommunityWithTypeSuccessfullyAddedNotification(community ui.PlainText, typ ui.PlainText) ui.RichText {
	return ui.RichText(fmt.Sprintf("I am here to help with `%s %s Community` requests and updates. You can tell me to leave at any time by using the command `/remove @adaptive`.", community, strings.Title(string(typ)))).Italics()
}

func VisionNotification(vision ui.RichText) ui.RichText {
	return ui.RichText(fmt.Sprintf("*Here is the vision for our company: *\n _%s_", strings.TrimSpace(string(vision))))
}

// NB! This function is used in two places!
func NotifyAboutNewAbilitiesInCommunityNotification(commType ui.PlainText) ui.RichText {
	var contextEntity string
	switch string(commType) {
	case string(community.Capability):
		contextEntity = "objectives"
	case string(community.Initiative):
		contextEntity = "initiatives"
	}
	return ui.RichText(fmt.Sprintf("You can now create %s to associate with the below `%s Community`.",
		contextEntity, strings.Title(string(commType)))).Italics()
}

var (
	SubscriptionCancellationConfirmation = ui.RichText("Ok, your subscription selection has been canceled").Italics()
)

func LeavingCommunityNotification(comm ui.PlainText) ui.RichText {
	return ui.RichText(fmt.Sprintf("Goodbye! Leaving `%s` community.", comm)).Italics()
}

var (
	UnsubscriptionErrorMessage    = ui.RichText("There was some issue with unsubscription").Italics()
	EngagementRemovedNotification = ui.RichText("Removed the engagement")
	// BUG?: this message is used in `CoachConfirm` section
	CoachingRequestRejectionReasonTitleToCoachee = ui.RichText("You provided the following information for not accepting the coach")
	CoachingRequestRejectionReasonTitleToCoach   = ui.RichText("You provided the following information for not accepting the coachee")
)

func SimulationUserDateNotificatoin(emulDate, emulUser string) ui.RichText {
	return ui.RichText(fmt.Sprintf("Simulated `%s` date for <@%s>", emulDate, emulUser)).Italics()
}

var (
	QueryChannelFeaturePrompt          ui.PlainText = "What sorts of requests and updates should I allow in this channel?"
	adaptiveChannelJoinMessage                      = "Thank you for inviting me!"
	InvitationToChannelAcknowledgement              = ui.RichText(adaptiveChannelJoinMessage).Italics()
	InvitationToChannelRejection                    = ui.RichText(fmt.Sprintf("%s Sorry, but there are no open communities for me to join. Please create a new capability or initiative community for me to join, or ask somebody in the strategy community", adaptiveChannelJoinMessage)).Italics()
)

func WelcomeUserToCommunity(userId string) ui.RichText {
	return ui.RichText(fmt.Sprintf("Hello <@%s>! ", userId))
}

var (
	communityMessageSuffix            = "Your Adaptive menu has changed! Just say `hi` to me to see your new menu options."
	UserCommunityWelcomeMessage       = ui.RichText("Welcome to the Adaptive User Community. Am glad to be by your side. " + communityMessageSuffix)
	HRCommunityWelcomeMessage         = ui.RichText("Welcome to the Adaptive HR Community. " + communityMessageSuffix)
	CoachingCommunityWelcomeMessage   = ui.RichText("Welcome to the Adaptive Coaching Community. " + communityMessageSuffix)
	AdminCommunityWelcomeMessage      = ui.RichText("Welcome to the Adaptive Admin Community. " + communityMessageSuffix)
	StrategyCommunityWelcomeMessage   = ui.RichText("Welcome to the Adaptive Strategy Community. " + communityMessageSuffix)
	CompetencyCommunityWelcomeMessage = ui.RichText("Welcome to the Adaptive Competencies Community. " + communityMessageSuffix)
	InitiativeCommunityWelcomeMessage = ui.RichText("Welcome to the Adaptive Initiative Community. " + communityMessageSuffix)
	CapabilityCommunityWelcomeMessage = ui.RichText("Welcome to the Adaptive Objective Community. " + communityMessageSuffix)
)

func InvitedByUserToCommunityNotification(userID string, community string) ui.RichText {
	return ui.RichText(fmt.Sprintf(
		"<@%s> invited Adaptive to a channel for Adaptive %s Community requests and updates", userID, strings.Title(community))).Italics()
}

func RemovedByUserFromCommunityNotification(userID string, community string) ui.RichText {
	return ui.RichText(fmt.Sprintf(
		"<@%s> removed Adaptive from the Adaptive %s Community channel", userID, strings.Title(community))).Italics()
}
