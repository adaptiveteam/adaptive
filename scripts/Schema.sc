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

val PriorityValueDef = StringBasedEnum("PriorityValue".camel, List(
	"UrgentPriority".camel ^^ "Urgent",
	"HighPriority".camel ^^ "High",
	"MediumPriority".camel ^^ "Medium",
	"LowPriority".camel ^^ "Low"
))

// val priorityValueDef = TypeAlias("PriorityValue".camel, string)
// lazy val priorityValue = TypeAlias("PriorityValue", "PriorityValue".camel, "\"\"", "S")
// lazy val priorityValue = modelsImport.simpleType("PriorityValue", "PriorityValue".camel, "\"\"", "S")


val ObjectiveStatusColorDef = StringBasedEnum("ObjectiveStatusColor".camel, List(
	"ObjectiveStatusRedKey".camel ^^ "Red",
	"ObjectiveStatusYellowKey".camel ^^ "Yellow",
	"ObjectiveStatusGreenKey".camel ^^ "Green"
))

val PlatformNameDef = StringBasedEnum("PlatformName".camel, List(
	"SlackPlatform".camel ^^ "slack",
	"MsTeamsPlatform".camel ^^ "ms-teams"
)) 

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

// val commonDefs = 
// lazy val ObjectiveStatusColor = modelsImport.simpleType("ObjectiveStatusColor", "ObjectiveStatusColor".camel, "\"\"", "S")
val commonPackage = Package("common".camel, List(
    Module(
        Filename("types".camel, ".go"),
        List(GoModulePart(
            List(), // imports.importClauses,
            List(
                platformIdDef, 
                PriorityValueDef, ObjectiveStatusColorDef,
                PlatformNameDef,
                AdaptiveCommunityIDDef
            )
        ))
    )
))

def goField(decl: String): Field = goFieldParser(goTypes)(decl)

val coreImport = ImportClause(Some("core"), "github.com/adaptiveteam/adaptive/core-utils-go")
val commonImport = ImportClause(Some("common"), "github.com/adaptiveteam/adaptive/daos/common")

val ObjectiveStatusColor = commonImport.simpleType(goPublicName(ObjectiveStatusColorDef.name), ObjectiveStatusColorDef.name, "\"\"", "S")
//  ObjectiveStatusColorDef.typeAliasTypeInfo
val priorityValue = commonImport.simpleType(goPublicName(PriorityValueDef.name), PriorityValueDef.name, "\"\"", "S")
val platformId = commonImport.simpleType(goPublicName(platformIdDef.name), platformIdDef.name, "\"\"", "S")
val PlatformName = commonImport.simpleType(goPublicName(PlatformNameDef.name), PlatformNameDef.name, "\"\"", "S")
//val PlatformName = TypeAlias("PlatformName".camel, string)

// val AdaptiveCommunityID = TypeAlias("AdaptiveCommunityID".camel, string)
val AdaptiveCommunityID = commonImport.simpleType(goPublicName(AdaptiveCommunityIDDef.name), AdaptiveCommunityIDDef.name, "\"\"", "S")

val imports = Imports(List(
    ImportClause(None, "github.com/aws/aws-sdk-go/aws"),
    ImportClause(Some("awsutils"), "github.com/adaptiveteam/adaptive/aws-utils-go"),
    commonImport,
    coreImport,
    ImportClause(None, "github.com/aws/aws-sdk-go/service/dynamodb"),
    ImportClause(None, "github.com/pkg/errors"),
    ImportClause(None, "fmt"),
    ImportClause(None, "encoding/json"),
    ImportClause(None, "strings")
    // ImportClause(None, "strconv")
))

val timeImport = ImportClause(None, "time")

val idField = underscoredName("id") :: string
val timezoneOffsetField = underscoredName("timezone_offset") :: int
val adaptiveScheduledTimeInUtcField = underscoredName("adaptive_scheduled_time_in_UTC") :: optionTimestamp
val platformIdField = spacedName("platform ID") :: platformId

def entitySpecificImports(entity: Entity): List[ImportClause] = {
    entity.traits.flatMap{
        case DeactivationTrait => List()
        case CreatedModifiedTimesTrait => List()
    }.distinct
}

def allEntitySpecificImports(entity: Entity): Imports = {
    entitySpecificImports(entity) ::: imports
}
// val userProfile = Entity(
//     spacedName("user profile"),
//     List(
//         idField,
//         Field(underscoredName("display_name"), string),
//         Field(underscoredName("first_name"), optionString),
//         Field(underscoredName("last_name"), optionString),
//         Field(underscoredName("timezone"), string),
//         timezoneOffsetField,
//         Field(underscoredName("adaptive_scheduled_time"), optionTimestamp, Some("in 24 hr format, localtime")),
//         adaptiveScheduledTimeInUtcField
//     ))

val user = Entity(
        spacedName("user"),
        List(idField),
        List(
            underscoredName("display_name") :: string,
            underscoredName("first_name") :: optionString,
            underscoredName("last_name") :: optionString,
            underscoredName("timezone") :: string,
            "IsAdaptiveBot".camel :: optionBoolean,
            timezoneOffsetField,
            (underscoredName("adaptive_scheduled_time") :: optionTimestamp) \\ "in 24 hr format, localtime",
            adaptiveScheduledTimeInUtcField,
            platformIdField,
            "PlatformOrg".camel :: optionString,
            spacedName("is admin") :: boolean,
            // spacedName("deleted") :: boolean,
            spacedName("is shared") :: boolean
        ),
        Nil, List(DeactivationTrait, CreatedModifiedTimesTrait)
)

val userTableDefaultIndex = Index(idField, None)
val userTable = Table(user, userTableDefaultIndex, List(
    Index(platformIdField, None),
    Index(platformIdField, Some(timezoneOffsetField)),
    Index(platformIdField, Some(adaptiveScheduledTimeInUtcField)),
))

val userPackage = defaultPackage(userTable, imports)

val coachQuarterYearField = spacedName("coach quarter year") :: string
val coacheeQuarterYearField = spacedName("coachee quarter year") :: string
val coacheeField = "coachee".camel :: string
val quarterField = "quarter".camel :: int
val yearField = "year".camel :: int
// TODO: Remove CoachingRelationship as it is not being used
val coachingRelationship = Entity("CoachingRelationship".camel,
        List(
//            platformIdField,
            coachQuarterYearField
        ),
        List(
            Field(spacedName("coachee quarter year"), string),
            coacheeField,
            quarterField,
            yearField,
            Field(spacedName("coach requested"), boolean),
            Field(spacedName("coachee requested"), boolean)
        )
)

val coachingRelationshipTable = Table(coachingRelationship, 
    Index(coachQuarterYearField, Some(coacheeField)),
    List(
        Index(coachQuarterYearField, None),
        Index(quarterField, Some(yearField)),
        Index(coacheeQuarterYearField, None)
    )
)

val coachingRelationshipDao = Dao(coachingRelationshipTable)

val coachingRelationshipPackage = defaultPackage(coachingRelationshipTable, imports)

val sourceField = "Source".camel :: string
val targetField = "Target".camel :: string
val quarterYearField = "QuarterYear".camel :: string
val channelIdField = "ChannelID".camel :: string
val channelIDField = channelIdField//"ChannelID".camel :: string
val channelIdFieldWithOldDbName = channelIdField.dbName("Channel".camel)  \\ "ChannelID is a channel identifier. TODO: rename db field `channel` to `channel_id`"

val UserFeedback = Entity(
    "UserFeedback".camel, 
    List(idField),
    List(
        sourceField,
        targetField,
        "ValueID".camel :: string,
        "ConfidenceFactor".camel :: string,
        "Feedback".camel :: string,
        quarterYearField,
        channelIdFieldWithOldDbName \\ // TODO: rename field to channel_id
            "ChannelID, if any, to engage user in response to the feedback" \\ 
            "This is useful to reply to an event with no knowledge of the previous context",
        ("MsgTimestamp".camel :: string) \\ "A reference to the original timestamp that can be used to reply via threading",
        platformIdField
    )
)

val UserFeedbackTable = Table(UserFeedback, 
    Index(idField, None),
    List(
        Index(quarterYearField, Some(sourceField)),
        Index(quarterYearField, Some(targetField))
    )
)

val UserFeedbackPackage = defaultPackage(UserFeedbackTable, imports)

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


val ebmImport = ImportClause(Some("ebm"), "github.com/adaptiveteam/adaptive/engagement-builder/model")

// it is used for:
// feedback requests, coaching requests, schedules to ask for create objective,
// request for a partner to respond about progress updates,
// request for a partner about comments,
// closeout request,
// 
val UserEngagement = Entity(
    "UserEngagement".camel,
    List(
        userIdField, // GetItem requires all primary key attributes
        idField \\ "A unique id to identify the engagement"),
    List(
        platformIdField \\ 
            "PlatformID is the identifier of the platform." \\ 
            "It's used to get platform token required to send message to Slack/Teams.",
        targetIdField,
        ("Namespace".camel :: string) \\ "Namespace for the engagement",
        ("CheckIdentifier".camel :: optionString) \\ "Check identifier for the engagement",
        ("CheckValue".camel :: optionBoolean),
        scriptField,
        ("Message".camel :: ebmImport.struct("Message".camel)) \\ "Message is the message we want to send to user",
        priorityField,
        ("Optional".camel :: boolean) \\ "A boolean flag indicating if it's optional",
        answeredField,
        ignoredField,
        "EffectiveStartDate".camel :: optionTimestamp,
        "EffectiveEndDate".camel :: optionTimestamp,
        "PostedAt".camel :: optionTimestamp,
        // ("CreatedAt".camel :: timestamp) \\ "When a same engagement is written to dynamo, dynamo doesn't update it and it not treated as a new event" \\
        //     "This timestamp will help to identify newer same engagement",
        ("RescheduledFrom".camel :: string) \\ "Re-scheduled timestamp for the engagement, if any"
    ), 
    Nil, List(CreatedModifiedTimesTrait)
) \\ "UserEngagement encapsulates an engagement we want to provide to a user"
// 	// RFC3339 timestamp at which the engagement was posted to the user
// 	 string `json:"posted_at,omitempty"`
// }

// TODO:   stream_enabled   = true
//  stream_view_type = var.dynamo_stream_view_type
val UserEngagementTable = Table(
    UserEngagement,
    Index(userIdField, Some(idField)),
    List(Index(userIdField, Some(answeredField), 
        projectionType = ProjectionType.INCLUDE, nonKeyAttributes =
            List(
                scriptField, priorityField, targetIdField, ignoredField
            )))
)

val UserEngagementPackage = defaultPackage(UserEngagementTable, ebmImport :: allEntitySpecificImports(UserEngagement))

val PostponedEvent = Entity("PostponedEvent".camel,
    List(idField),
    List(
        userIdField,
        platformIdField,
        ("ActionPath".camel :: string) \\ "ActionPath is callback for triggering workflows", 
        ("ValidThrough".camel :: timestamp) \\ "ValidThrough is the last time moment when this event might still be valid"
    ),
    Nil, List(CreatedModifiedTimesTrait)
)

val PostponedEventTable = Table(PostponedEvent,
    Index(idField, None),
    List(
        Index(platformIdField, Some(userIdField)),
        Index(userIdField, None)
    )
)

val PostponedEventPackage = defaultPackage(PostponedEventTable, allEntitySpecificImports(PostponedEvent))

val userIDField = userIdField // "UserID".camel :: string
val communityIDField = "CommunityID".camel :: string
val AdaptiveCommunityUser = Entity(
    "AdaptiveCommunityUser".camel, 
    List(
        channelIDField,
        userIDField
    ),
    List(
        platformIdField,
        communityIDField,
    )
)
val AdaptiveCommunityUserTable = Table(AdaptiveCommunityUser, 
    Index(channelIDField, Some(userIDField)),
    List(
        Index(channelIDField, None),
        Index(userIDField, Some(communityIDField)),
        Index(userIDField, None),
        Index(platformIdField, Some(communityIDField))
    )
)
val AdaptiveCommunityUserPackage = defaultPackage(AdaptiveCommunityUserTable, imports)

// val channelField = "Channel".camel :: string
val AdaptiveCommunity = Entity(
    "AdaptiveCommunity".camel, 
    List(platformIdField, idField),
    List(
        channelIdFieldWithOldDbName, // TODO: rename db field to channel_id
        "Active".camel :: boolean,
        "RequestedBy".camel :: string
    ),
    Nil, 
    List(CreatedModifiedTimesTrait, DeactivationTrait)
)

// TODO: 
//  stream_enabled   = true
//  stream_view_type = var.dynamo_stream_view_type
val AdaptiveCommunityTable = Table(AdaptiveCommunity, 
    Index(idField, Some(platformIdField)),
    List(
        Index(channelIdFieldWithOldDbName, None),
        Index(platformIdField, None)
    )
)

val AdaptiveCommunityPackage = defaultPackage(AdaptiveCommunityTable, allEntitySpecificImports(AdaptiveCommunity))


val attrKeyField = ("AttrKey".camel::string) \\ "Key of the setting"
val UserAttribute = Entity(
    "UserAttribute".camel, 
    List(
        userIdField,
        attrKeyField
    ),
    List(
        ("AttrValue".camel::string) \\ "Value of the setting",
        ("Default".camel :: boolean) \\ 
            "A flag that tells whether setting is default or is explicitly set" \\ 
            "Every user will have default settings"
    )
) \\ "UserAttribute encapsulates key-value setting for a user"

val UserAttributeTable = Table(UserAttribute, 
    Index(userIdField, Some(attrKeyField)), // GetItem requires all primary key attributes
    List()
)

val UserAttributePackage = defaultPackage(UserAttributeTable, imports)

val dateField = "Date".camel::timestamp

val AdHocHoliday = Entity(
    "AdHocHoliday".camel,
    List(idField),
    List(
        platformIdField,
        dateField,
        "Name".camel::string,
        "Description".camel::string,
        "ScopeCommunities".camel::string
    ),
    Nil, List(DeactivationTrait)
) \\ "AdHocHoliday is a holiday on exact date."

val AdHocHolidayTable = Table(
    AdHocHoliday,
    Index(idField, None),
    List(
        Index(platformIdField, Some(dateField))
    )
)

val AdHocHolidayPackage = defaultPackage(AdHocHolidayTable, allEntitySpecificImports(AdHocHoliday))

/*

const (
)

*/
val completedField = ("Completed".camel :: int) \\ "1 for true, 0 for false"
val acceptedField = ("Accepted".camel :: int) \\ "1 for true, 0 for false"
val accountabilityPartnerField = "AccountabilityPartner".camel :: string
val DevelopmentObjectiveType = StringBasedEnum("DevelopmentObjectiveType".camel, List(
	"IndividualDevelopmentObjective".camel ^^ "individual",
	"StrategyDevelopmentObjective".camel ^^ "strategy"
))
val AlignedStrategyType = StringBasedEnum("AlignedStrategyType".camel, 
    List(
        "ObjectiveStrategyObjectiveAlignment".camel ^^ "strategy_objective", 
        "ObjectiveStrategyInitiativeAlignment".camel ^^ "strategy_initiative", 
        "ObjectiveCompetencyAlignment".camel ^^ "competency", 
        "ObjectiveNoStrategyAlignment".camel ^^ "none"))

val objectiveTypeField = ("ObjectiveType".camel :: DevelopmentObjectiveType.typeAliasTypeInfo)
    .dbName("type".camel)

val nameField = "Name".camel :: string
val descriptionField = "Description".camel :: string

val StrategyObjectiveType = TypeAlias("StrategyObjectiveType".camel, string)
// TODO: rename field in DB and then remove `dbName` 
val capabilityCommunityIDsField = (spacedName("capability community IDs") :: optionStringArray).
    dbName(spacedName("capability community IDs")) \\ "community id not require d for customer/financial objectives"
val createdByField = "CreatedBy".camel :: optionString
val modifiedByField = "ModifiedBy".camel :: optionString
val advocateField = "Advocate".camel :: string
val initiativeCommunityIDField = "InitiativeCommunityID".camel :: string

val UserObjective = Entity(
    "UserObjective".camel,
    List(idField),
    List(
        platformIdField,
        userIdField, // advocateField,
        nameField,
        descriptionField,
        accountabilityPartnerField,
        acceptedField,
        objectiveTypeField,
        "StrategyAlignmentEntityID".camel :: optionString,
        "StrategyAlignmentEntityType".camel :: AlignedStrategyType.typeAliasTypeInfo,
        "Quarter".camel :: int,
        "Year".camel :: int,
        ("CreatedDate".camel :: timestamp) \\ "Deprecated, use CreatedAt automated field", 
        "ExpectedEndDate".camel :: timestamp,
        completedField,
        "PartnerVerifiedCompletion".camel :: boolean,
        "CompletedDate".camel :: optionTimestamp,
        "PartnerVerifiedCompletionDate".camel :: optionTimestamp,
        "Comments".camel :: optionString, 
        ("Cancelled".camel :: int) \\ "1 for true, 0 for false",

        // TODO: add strategy objective and strategy initiative fields:
        // "AsMeasuredBy".camel :: optionString,
        // "Targets".camel :: optionString,
        // // objectiveTypeField,
        // // capabilityCommunityIDsField,
        // //
        // "DefinitionOfVictory".camel :: optionString,
        // initiativeCommunityIDField,
        // "Budget".camel :: optionString,
        // "CapabilityObjective".camel :: optionString,
        createdByField,
        modifiedByField
    ), 
    Nil, 
    List(CreatedModifiedTimesTrait)
)

val UserObjectiveTable = Table(
    UserObjective,
    Index(idField, None),//Some(userIdField)), We don't need sorting by user. 
    // Instead we want to be able to GetItem by just an ID. Sort key does not allow this
    // Index(userIdField, Some(idField)),
    List(
        Index(userIdField, Some(completedField)),
        // Index(idField, None),
        Index(acceptedField, None),
        Index(accountabilityPartnerField, None),
        Index(userIdField, Some(objectiveTypeField))
    )
)

val UserObjectivePackage = Package(UserObjective.name, 
    List(
        DevelopmentObjectiveType :: AlignedStrategyType :: 
            daoModule(UserObjectiveTable, allEntitySpecificImports(UserObjective)),
        daoConnectionModule(UserObjectiveTable, allEntitySpecificImports(UserObjective)),
        fieldNamesModule(UserObjectiveTable, allEntitySpecificImports(UserObjective))
    )
)

val createdOnField = "CreatedOn".camel :: string

// val ObjectiveStatusColor = StringBasedEnum("ObjectiveStatusColor".camel, List(
//   "ObjectiveStatusRed".camel ^^ "Red",
//   "ObjectiveStatusYellow".camel ^^ "Yellow",
//   "ObjectiveStatusGreen".camel ^^ "Green"
// ))

val UserObjectiveProgress = Entity(
    "UserObjectiveProgress".camel,
    List(
        idField,
        createdOnField
    ),
    List(
        platformIdField,
        "UserID".camel :: string,
        "PartnerID".camel :: string,
        "Comments".camel :: string,
        ("Closeout".camel :: int) \\ "1 for true, 0 for false",
        "PercentTimeLapsed".camel :: string,
        "StatusColor".camel :: ObjectiveStatusColor,//.typeAliasTypeInfo,
        "ReviewedByPartner".camel :: bool,
        "PartnerComments".camel :: optionString,
        "PartnerReportedProgress".camel :: optionString
    )
)

val UserObjectiveProgressTable = Table(
    UserObjectiveProgress,
    Index(idField, Some(createdOnField)),
    List(
        Index(idField, None),
        Index(createdOnField, None),
    )
)

val UserObjectiveProgressPackage = defaultPackage(UserObjectiveProgressTable, imports)

val AdaptiveValue = Entity(
    "AdaptiveValue".camel, 
    List(idField),
    List(
        platformIdField,
        nameField.dbName("ValueName".camel),
        "ValueType".camel :: string,
        descriptionField
    ), Nil,
    List(DeactivationTrait)
) \\ "AdaptiveValue is a value for a client (Reliability, Skill, Contribution, and Productivity)"

val AdaptiveValueTable = Table(AdaptiveValue, 
    Index(idField, None), 
    List(
        Index(platformIdField, None)
    )
)

val AdaptiveValuePackage = defaultPackage(AdaptiveValueTable, entitySpecificImports(AdHocHoliday) ::: imports)

val ClientPlatformToken = Entity(
    "ClientPlatformToken".camel, 
    List(
        platformIdField
    ),
    List(
        "Org".camel :: string,
        ("PlatformName".camel :: PlatformName) \\ "should be slack or ms-teams",
        "PlatformToken".camel :: string,
        "ContactFirstName".camel :: string,
        "ContactLastName".camel :: string,
        "ContactMail".camel :: string,
    )
)

val ClientPlatformTokenTable = Table(ClientPlatformToken, 
    Index(platformIdField, None), List()
)
val ClientPlatformTokenPackage = defaultPackage(ClientPlatformTokenTable, imports)

val StrategyObjective = Entity(
    "StrategyObjective".camel, 
    List(
        platformIdField \\ "range key",// GetItem requires all primary key attributes
        idField \\ "hash"),
    
    List(
        nameField,
        descriptionField,
        "AsMeasuredBy".camel :: string,
        "Targets".camel :: string,
        "ObjectiveType".camel :: TypeAliasTypeInfo(StrategyObjectiveType),
        advocateField,
        capabilityCommunityIDsField,
        "ExpectedEndDate".camel :: timestamp,
        createdByField
    ), Nil, List(CreatedModifiedTimesTrait))

val StrategyObjectiveTable = Table(StrategyObjective,
    Index(idField, Some(platformIdField)),
    List(
        Index(platformIdField, None),
        Index(capabilityCommunityIDsField, None)
    )
)

val StrategyObjectivePackage = Package(StrategyObjective.name, 
    List(StrategyObjectiveType :: 
        daoModule(StrategyObjectiveTable, allEntitySpecificImports(StrategyObjective)),
        daoConnectionModule(StrategyObjectiveTable, allEntitySpecificImports(StrategyObjective)),
        fieldNamesModule(StrategyObjectiveTable, allEntitySpecificImports(StrategyObjective))
    )
)

// ObjectiveTypeDictionary contains types of objectives.
val ObjectiveTypeDictionary = Entity(
    "ObjectiveTypeDictionary".camel, 
    List(idField \\ "hash"),
    List(
        platformIdField \\ "range key",
        nameField,
        descriptionField,
    ), Nil, List(CreatedModifiedTimesTrait, DeactivationTrait))

val ObjectiveTypeDictionaryTable = Table(ObjectiveTypeDictionary,
    Index(idField, Some(platformIdField)),
    List(
        Index(platformIdField, None)
    )
)

val ObjectiveTypeDictionaryPackage = defaultPackage(ObjectiveTypeDictionaryTable, allEntitySpecificImports(ObjectiveTypeDictionary))

val TypedObjective = Entity("TypedObjective".camel, 
    List(idField \\ "hash"),
    List(
        platformIdField \\ "range key",
        nameField,
        descriptionField,
        "AsMeasuredBy".camel :: string,
        "Targets".camel :: string,
        "TypeID".camel :: string, // TODO: Add Foreign Key constraint (PlatformID, TypeID)
        advocateField,
        capabilityCommunityIDsField,
        "ExpectedEndDate".camel :: timestamp,
        "CreatedBy".camel :: string
    ), Nil, List(CreatedModifiedTimesTrait, DeactivationTrait)
)

val TypedObjectiveTable = Table(TypedObjective,
    Index(idField, Some(platformIdField)),
    List(
        Index(platformIdField, None)
    )
)

val TypedObjectivePackage = defaultPackage(TypedObjectiveTable, allEntitySpecificImports(TypedObjective))

val StrategyInitiative = Entity(
    "StrategyInitiative".camel, 
    List(
        platformIdField,// GetItem requires all primary key attributes
        idField
        ),
    List(
        nameField,
        descriptionField,
        "DefinitionOfVictory".camel :: string,
        advocateField,
        initiativeCommunityIDField,
        "Budget".camel :: string,
        "ExpectedEndDate".camel :: timestamp,
        "CapabilityObjective".camel :: string,
        //createdAtField,
        createdByField,
        //"ModifiedAt".camel :: timestamp,
        modifiedByField
    ), Nil, List(CreatedModifiedTimesTrait))

val StrategyInitiativeTable = Table(StrategyInitiative, 
    Index(idField, Some(platformIdField)),
    List(
        Index(platformIdField, None),
        Index(initiativeCommunityIDField, None)
    )
)

val StrategyInitiativePackage = defaultPackage(StrategyInitiativeTable, allEntitySpecificImports(StrategyInitiative))

val VisionMission = Entity(
    "VisionMission".camel, 
    List(
        platformIdField // GetItem requires all primary key attributes
        ),
    List(
        idField, // TODO: remove id field
        "Mission".camel :: string,
        "Vision".camel :: string,
        advocateField,
        //createdAtField,
        createdByField
    ), Nil, List(CreatedModifiedTimesTrait))


val VisionMissionTable = Table(VisionMission, 
    Index(platformIdField, None),
    List()
)

val VisionMissionPackage = defaultPackage(VisionMissionTable, allEntitySpecificImports(VisionMission))

val channelCreatedField = ("ChannelCreated".camel :: int)\\ "0 for false"
val StrategyCommunity = Entity(
    "StrategyCommunity".camel, 
    List(idField),
    List(
        platformIdField,
        advocateField,
        //createdAtField,
        "Community".camel :: AdaptiveCommunityID,
        channelIdField,
        channelCreatedField,
        "AccountabilityPartner".camel :: string,
        "ParentCommunity".camel :: AdaptiveCommunityID,
        "ParentCommunityChannelID".camel :: string
    ), Nil, List(CreatedModifiedTimesTrait))

val StrategyCommunityTable = Table(StrategyCommunity, 
    Index(idField, None),
    List(
        Index(platformIdField, Some(channelCreatedField)),
        Index(platformIdField, None),
        Index(channelIdField, None)
    )
)

val StrategyCommunityPackage = defaultPackage(StrategyCommunityTable, allEntitySpecificImports(StrategyCommunity))

// val StrategyCommunityPackage = Package(StrategyCommunity.name, 
//     List(//AdaptiveCommunityID :: 
//         daoModule(StrategyCommunityTable, allEntitySpecificImports(StrategyCommunity)),
//         daoConnectionModule(table, imports),
//         fieldNamesModule(table, imports)
//     )
// )

val CapabilityCommunity = Entity(
    "CapabilityCommunity".camel, 
    List(
        platformIdField,// GetItem requires all primary key attributes
        idField
        ),
    List(
        nameField,
        descriptionField,
        advocateField,
    // createdAtField,
        createdByField
    ), Nil, List(CreatedModifiedTimesTrait))

val CapabilityCommunityTable = Table(CapabilityCommunity, 
    Index(idField, Some(platformIdField)),
    List(
        Index(platformIdField, None)
    )
)

val CapabilityCommunityPackage = defaultPackage(CapabilityCommunityTable, allEntitySpecificImports(CapabilityCommunity))

val StrategyInitiativeCommunity = Entity(
    "StrategyInitiativeCommunity".camel, 
    List(idField, platformIdField), // GetItem requires all primary key attributes
    List(
        nameField,
        descriptionField,
        advocateField,
        "CapabilityCommunityID".camel :: string,
        //createdAtField,
        createdByField
    ), Nil, List(CreatedModifiedTimesTrait))

val StrategyInitiativeCommunityTable = Table(StrategyInitiativeCommunity, 
    Index(idField, Some(platformIdField)),
    List(
        Index(platformIdField, None)
    )
)

val StrategyInitiativeCommunityPackage = defaultPackage(StrategyInitiativeCommunityTable, allEntitySpecificImports(StrategyInitiativeCommunity))

val dialogIDField = ("DialogID".camel :: string) \\ "This is an immutable UUID that developers can use"
val contextField = ("Context".camel :: string) \\ "This is the context path for the piece of dialog"
val subjectField = ("Subject".camel :: string) \\ "This is the dialog subject"
val DialogEntry = Entity(
    "DialogEntry".camel, 
    List(
        dialogIDField
    ),
    List(
        contextField,
        subjectField,
        ("Updated".camel :: string) \\ "This was when the dialog was last updated",
        ("Dialog".camel :: stringArray) \\ "These are the dialog options",
        ("Comments".camel :: stringArray) \\ "Comments to help cultivators understand the dialog intent",
        ("LearnMoreLink".camel :: string) \\ "This the link to the LearnMore page",
        ("LearnMoreContent".camel :: string) \\ "This is the actual content from the LearnMore page",
        "BuildBranch".camel :: string,
        "CultivationBranch".camel :: string,
        "MasterBranch".camel :: string,
        "BuildID".camel :: string
    )) \\ "DialogEntry stores all of the  relevant information for a piece of dialog including:"

val DialogEntryTable = Table(DialogEntry, 
    Index(dialogIDField, None),
    List(
        Index(contextField, Some(subjectField))
    ),
    encrypted = true
)

val DialogEntryPackage = defaultPackage(DialogEntryTable, imports)

val applicationAliasField = "ApplicationAlias".camel :: string
val ContextAliasEntry = Entity(
    "ContextAliasEntry".camel, 
    List(applicationAliasField),
    List(
        "Context".camel :: string,
        "BuildID".camel :: string
    )) \\ 
    "ContextAliasEntry contains all of the information needed for a context alias" \\
    "A context alias is a way to alias  a piece of context without spelling out" \\
    "the context path.  If the path changes you can still safely use the alias."

val ContextAliasEntryTable = Table(ContextAliasEntry, 
    Index(applicationAliasField, None), List())

val ContextAliasEntryPackage = defaultPackage(ContextAliasEntryTable, imports)

val teamIdField = "TeamID".camel :: platformId
// SlackTeam contains the information about Slack team that 
// we obtain during OAuth2 authorization
val SlackTeam = Entity(
    "SlackTeam".camel,//spacedName("OAuth state"),
    List(teamIdField),
    List(
        "AccessToken".camel :: optionString,
        "TeamName".camel :: optionString,
        userIdField,
        "EnterpriseID".camel :: optionString,
        "BotUserID".camel :: optionString, // we may add bot user
        ("BotAccessToken".camel :: optionString) \\ "bot_access_token",
        "Scopes".camel :: optionStringArray
    ),
    Nil, List(CreatedModifiedTimesTrait)
)
val SlackTeamTable = Table(SlackTeam,
    Index(teamIdField, None), 
    List()
)

val SlackTeamPackage = defaultPackage(SlackTeamTable, imports)

val packages = List(
    commonPackage,

    userPackage,
    coachingRelationshipPackage,
    UserFeedbackPackage,
    UserEngagementPackage,
    AdaptiveCommunityUserPackage,
    AdaptiveCommunityPackage,
    UserAttributePackage,
    AdHocHolidayPackage,
    UserObjectivePackage,
    UserObjectiveProgressPackage,
    AdaptiveValuePackage,
    ClientPlatformTokenPackage,
    StrategyObjectivePackage,
    StrategyInitiativePackage,
    VisionMissionPackage,
    StrategyCommunityPackage,
    CapabilityCommunityPackage,
    StrategyInitiativeCommunityPackage,
    DialogEntryPackage,
    ContextAliasEntryPackage,
    ObjectiveTypeDictionaryPackage,
    PostponedEventPackage,
    SlackTeamPackage
    )
val daosProject = GoProjectFolder("daos", packages)

val coreTerraformProject = TerraformProjectFolder("daos/terraform", List(
    userTable,
    coachingRelationshipTable,
    UserFeedbackTable,
    UserEngagementTable,
    AdaptiveCommunityUserTable,
    AdaptiveCommunityTable,
    AdHocHolidayTable,
    UserObjectiveTable,
    UserObjectiveProgressTable,
    AdaptiveValueTable,
    ClientPlatformTokenTable,
    StrategyObjectiveTable,
    StrategyInitiativeTable,
    VisionMissionTable,
    StrategyCommunityTable,
    CapabilityCommunityTable,
    StrategyInitiativeCommunityTable,
    DialogEntryTable,
    ContextAliasEntryTable,
    ObjectiveTypeDictionaryTable,
    PostponedEventTable,
    SlackTeamTable
    ))

val workspace: Workspace = List(daosProject, coreTerraformProject)

