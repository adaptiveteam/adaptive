import $file.Meta
import Meta._

import $file.Dsl
import Dsl._

import $file.GoTypes
import GoTypes._

import $file.Templates
import Templates._

import $file.SchemaCommon
import SchemaCommon._

import $file.SchemaCommunities
import SchemaCommunities._

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
                AdaptiveCommunityIDDef,
                CommunityKindDef
            )
        ))
    )
))

def goField(decl: String): Field = goFieldParser(goTypes)(decl)

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
        "User".camel,
        List(platformIdField, idField),
        List(
            underscoredName("display_name") :: string,
            underscoredName("first_name") :: optionString,
            underscoredName("last_name") :: optionString,
            underscoredName("timezone") :: string,
            "IsAdaptiveBot".camel :: optionBoolean,
            timezoneOffsetField,
            (underscoredName("adaptive_scheduled_time") :: optionTimestamp) \\ "in 24 hr format, localtime",
            adaptiveScheduledTimeInUtcField,
            "PlatformOrg".camel :: optionString,
            spacedName("is admin") :: boolean,
            // spacedName("deleted") :: boolean,
            spacedName("is shared") :: boolean
        ),
        Nil, List(DeactivationTrait, CreatedModifiedTimesTrait)
)

val userTableDefaultIndex = 
    Index(idField, Some(platformIdField))
    
val userTable = Table(user, userTableDefaultIndex, List(
    Index(platformIdField, None),
    Index(platformIdField, Some(timezoneOffsetField)),
    Index(platformIdField, Some(adaptiveScheduledTimeInUtcField)),
))

val userPackage = defaultPackage(userTable, imports)

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
    Nil, List(CreatedModifiedTimesTrait, DeactivationTrait)
)

val PostponedEventTable = Table(PostponedEvent,
    Index(idField, None),
    List(
        Index(platformIdField, Some(userIdField)),
        Index(userIdField, None)
    )
)

val PostponedEventPackage = defaultPackage(PostponedEventTable, allEntitySpecificImports(PostponedEvent))


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
	"StrategyDevelopmentObjective".camel ^^ "strategy",
	"StrategyDevelopmentObjectiveIssue".camel ^^ "strategy_objective",
	"StrategyDevelopmentInitiative".camel ^^ "strategy_initiative"
))
val AlignedStrategyType = StringBasedEnum("AlignedStrategyType".camel, 
    List(
        "ObjectiveStrategyObjectiveAlignment".camel ^^ "strategy_objective", 
        "ObjectiveStrategyInitiativeAlignment".camel ^^ "strategy_initiative", 
        "ObjectiveCompetencyAlignment".camel ^^ "competency", 
        "ObjectiveNoStrategyAlignment".camel ^^ "none"))

val objectiveTypeField = ("ObjectiveType".camel :: DevelopmentObjectiveType.typeAliasTypeInfo)
    .dbName("type".camel)


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

val migrationIdField = ("MigrationID".camel :: string) \\ "Human-friendly identifier of migration. Should start with 3 digits for sorting purposes. Unique within platform"

val Migration = Entity(
    "Migration".camel,
    List(platformIdField, migrationIdField),
    List(
        "SuccessCount".camel :: int,
        "FailuresCount".camel :: int,
    ),
    Nil,
    List(CreatedModifiedTimesTrait),
)

val MigrationTable = Table(Migration,
    Index(platformIdField, Some(migrationIdField)), 
    List()
)

val MigrationPackage = defaultPackage(MigrationTable, imports)

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
    SlackTeamPackage,
    CommunityPackage,
    ChannelMemberPackage,
    MigrationPackage,
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
    SlackTeamTable,
    CommunityTable,
    ChannelMemberTable,
    MigrationTable,
    ))

val workspace: Workspace = List(daosProject, coreTerraformProject)

