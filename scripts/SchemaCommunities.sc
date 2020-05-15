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

val AdaptiveCommunityTable = Table(AdaptiveCommunity, 
    Index(idField, Some(platformIdField)),
    List(
        Index(channelIdFieldWithOldDbName, None),
        Index(platformIdField, None)
    )
)

val AdaptiveCommunityPackage = defaultPackage(AdaptiveCommunityTable, allEntitySpecificImports(AdaptiveCommunity))




val channelCreatedField = ("ChannelCreated".camel :: int)\\ "0 for false"
val StrategyCommunity = Entity(
    "StrategyCommunity".camel, 
    List(idField),
    List(
        platformIdField,
        advocateField,
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
        createdByField
    ), Nil, List(CreatedModifiedTimesTrait))

val StrategyInitiativeCommunityTable = Table(StrategyInitiativeCommunity, 
    Index(idField, Some(platformIdField)),
    List(
        Index(platformIdField, None)
    )
)

val StrategyInitiativeCommunityPackage = defaultPackage(StrategyInitiativeCommunityTable, allEntitySpecificImports(StrategyInitiativeCommunity))



val communityKindField = "CommunityKind".camel :: CommunityKind

val Community = Entity(
    "Community".camel, 
    List(platformIdField, idField),
    List(
        channelIdOptionalField,
        communityKindField,

        "ParentCommunityID".camel :: optionString,

        nameField,
        descriptionField,

        advocateField \\ "Owner, responsible person",
        ("AccountabilityPartner".camel :: optionString) \\ "Nudging person",

        createdByField,
        modifiedByField
        // requestedByField
    ),
    Nil, 
    List(CreatedModifiedTimesTrait, DeactivationTrait)
)

val CommunityTable = Table(Community, 
    Index(idField, Some(platformIdField)),
    List(
        Index(channelIdOptionalField, Some(platformIdField)),
        Index(platformIdField, Some(communityKindField))
    )
)

val CommunityPackage = defaultPackage(CommunityTable, allEntitySpecificImports(Community))


val ChannelMember = Entity(
    "ChannelMember".camel, 
    List(
        channelIDField,
        userIDField
    ),
    List(
        platformIdField,
    )
)
val ChannelMemberTable = Table(ChannelMember, 
    Index(channelIDField, Some(userIDField)),
    List(
        Index(channelIDField, None),
        Index(userIDField, None),
        Index(platformIdField, Some(userIDField))
    )
)
val ChannelMemberPackage = defaultPackage(ChannelMemberTable, imports)
