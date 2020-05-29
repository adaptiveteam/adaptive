import $file.Meta
import Meta._

import $file.Dsl
import Dsl._

import $file.GoTypes
import GoTypes._

import $file.Templates
import Templates._

// lazy val modelsImport = ImportClause(None, "github.com/adaptiveteam/adaptive/adaptive-utils-go/models")
val platformIdDef = TypeAlias("PlatformID".camel, string)
// val platformId = simpleType("PlatformId", "models.PlatformID", "S")
// lazy val platformId = modelsImport.simpleType("PlatformId", "PlatformID".camel, "\"\"", "S")


val AdaptiveCommunityIDDef = StringBasedEnum("AdaptiveCommunityID".camel, List(
	"Admin".camel ^^ "admin",
	spacedName("HR") ^^ "hr",
	"Coaching".camel ^^ "coaching",
	"User".camel ^^ "user",
	"Strategy".camel ^^ "strategy",
	"Capability".camel ^^ "capability",
	"Initiative".camel ^^ "initiative",
	"Competency".camel ^^ "competency"
))

val CommunityKindDef = StringBasedEnum("CommunityKind".camel, List(
	"AdminCommunity".camel ^^ "admin",          // one
	spacedName("HRCommunity") ^^ "hr",          // one
	"CoachingCommunity".camel ^^ "coaching",    // one
	"UserCommunity".camel ^^ "user",            // one
	"StrategyCommunity".camel ^^ "strategy",    // one
	"CompetencyCommunity".camel ^^ "competency",// one
	// "ObjectiveManagementCommunity".camel ^^ "objective-management",
	// "InitiativeManagementCommunity".camel ^^ "initiative-management",
	"ObjectiveCommunity".camel ^^ "objective",  // many
	"InitiativeCommunity".camel ^^ "initiative" // many
))

lazy val commonImport = ImportClause(Some("common"), "github.com/adaptiveteam/adaptive/daos/common")

lazy val coreImport = ImportClause(Some("core"), "github.com/adaptiveteam/adaptive/core-utils-go")

lazy val awsutilsImport = ImportClause(Some("awsutils"), awsUtilsUrl)

lazy val basicImportsForConnections = Imports(List(
    ImportClause(None, "github.com/aws/aws-sdk-go/aws"),
    coreImport,
    commonImport,
    ImportClause(None, "github.com/aws/aws-sdk-go/service/dynamodb"),
    ImportClause(None, "github.com/pkg/errors"),
    ImportClause(None, "fmt"),
    // ImportClause(None, "strconv")
))

val imports = List(
    awsutilsImport,
    ImportClause(None, "encoding/json"),
    ImportClause(None, "strings"),
) ::: basicImportsForConnections

val timeImport = ImportClause(None, "time")

def entitySpecificImports(entity: Entity): List[ImportClause] = {
    entity.traits.flatMap{
        case DeactivationTrait => List()
        case CreatedModifiedTimesTrait => List()
    }.distinct
}

def allEntitySpecificImports(entity: Entity): Imports = {
    entitySpecificImports(entity) ::: imports
}
// val AdaptiveCommunityID = TypeAlias("AdaptiveCommunityID".camel, string)
val AdaptiveCommunityID = commonImport.simpleType(goPublicName(AdaptiveCommunityIDDef.name), AdaptiveCommunityIDDef.name, "\"\"", "S")

val CommunityKind = commonImport.simpleType(goPublicName(CommunityKindDef.name), CommunityKindDef.name, "\"\"", "S")

val PriorityValueDef = StringBasedEnum("PriorityValue".camel, List(
	"UrgentPriority".camel ^^ "Urgent",
	"HighPriority".camel ^^ "High",
	"MediumPriority".camel ^^ "Medium",
	"LowPriority".camel ^^ "Low"
))

val ObjectiveStatusColorDef = StringBasedEnum("ObjectiveStatusColor".camel, List(
	"ObjectiveStatusRedKey".camel ^^ "Red",
	"ObjectiveStatusYellowKey".camel ^^ "Yellow",
	"ObjectiveStatusGreenKey".camel ^^ "Green"
))

val PlatformNameDef = StringBasedEnum("PlatformName".camel, List(
	"SlackPlatform".camel ^^ "slack",
	"MsTeamsPlatform".camel ^^ "ms-teams"
)) 




val ObjectiveStatusColor = commonImport.simpleType(goPublicName(ObjectiveStatusColorDef.name), ObjectiveStatusColorDef.name, "\"\"", "S")
//  ObjectiveStatusColorDef.typeAliasTypeInfo
val priorityValue = commonImport.simpleType(goPublicName(PriorityValueDef.name), PriorityValueDef.name, "\"\"", "S")
val platformId = commonImport.simpleType(goPublicName(platformIdDef.name), platformIdDef.name, "\"\"", "S")
val PlatformName = commonImport.simpleType(goPublicName(PlatformNameDef.name), PlatformNameDef.name, "\"\"", "S")
//val PlatformName = TypeAlias("PlatformName".camel, string)


val idField = underscoredName("id") :: string
val timezoneOffsetField = underscoredName("timezone_offset") :: int
val adaptiveScheduledTimeInUtcField = underscoredName("adaptive_scheduled_time_in_UTC") :: optionTimestamp
val platformIdField = spacedName("platform ID") :: platformId


val coachQuarterYearField = spacedName("coach quarter year") :: string
val coacheeQuarterYearField = spacedName("coachee quarter year") :: string
val coacheeField = "coachee".camel :: string
val quarterField = "quarter".camel :: int
val yearField = "year".camel :: int


val sourceField = "Source".camel :: string
val targetField = "Target".camel :: string
val quarterYearField = "QuarterYear".camel :: string
val channelIdField = "ChannelID".camel :: string
val channelIdOptionalField = "ChannelID".camel :: optionString
val channelIDField = channelIdField//"ChannelID".camel :: string
val channelIdFieldWithOldDbName = channelIdField.dbName("Channel".camel)  \\ "ChannelID is a channel identifier. TODO: rename db field `channel` to `channel_id`"

val userIdField = ("UserID".camel :: string) \\ 
    "UserID is the ID of the user to send an engagement to" \\
    "This usually corresponds to the platform user id"
val targetIdField = ("TargetID".camel :: string) \\ 
    "TargetID is the ID of the user for whom this is related to"
val answeredField = ("Answered".camel :: int) \\ 
    "Answered is a flag indicating that a user has responded to the engagement: 1 for answered, 0 for un-answered. "\\
    "This is required because, we need to keep the engagement even after a user has answered it. "\\
    "If the user wants to edit later, we will refer to the same engagement to post to user, like getting survey information"\\
    "So, we need a way to differentiate between answered and unanswered engagements"
val scriptField = ("Script".camel :: string) \\ 
    "Script that should be sent to a user to start engaging." \\
    "It's a serialized ebm.Message" \\
    "deprecated. Use `Message` directly."
val priorityField = ("Priority".camel :: priorityValue) \\ "Priority of the engagement" \\
    "Urgent priority engagements are immediately sent to a user" \\
    "Other priority engagements are queued up in the order of priority to be sent to user in next window"
val ignoredField = ("Ignored".camel :: int) \\ "Flag indicating if an engagement is ignored, 1 for yes, 0 for no"

val userIDField = userIdField // "UserID".camel :: string


val nameField = "Name".camel :: string
val descriptionField = "Description".camel :: string

val StrategyObjectiveType = TypeAlias("StrategyObjectiveType".camel, string)
// TODO: rename field in DB and then remove `dbName` 
val capabilityCommunityIDsField = (spacedName("capability community IDs") :: optionStringArray).
    dbName(spacedName("capability community IDs")) \\ "community id not require d for customer/financial objectives"
val createdByField = "CreatedBy".camel :: optionString
val modifiedByField = "ModifiedBy".camel :: optionString
// val requestedByField = "RequestedBy".camel :: optionString
val advocateField = "Advocate".camel :: optionString
val initiativeCommunityIDField = "InitiativeCommunityID".camel :: string
