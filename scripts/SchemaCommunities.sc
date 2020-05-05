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

