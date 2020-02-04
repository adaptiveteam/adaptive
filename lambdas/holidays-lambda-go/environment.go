package lambda

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	eholidays "github.com/adaptiveteam/adaptive/adaptive-engagements/holidays"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
)

var (
	platformNotificationTopic      = utils.NonEmptyEnv("PLATFORM_NOTIFICATION_TOPIC")
	namespace                      = utils.NonEmptyEnv("LOG_NAMESPACE")
	region                         = utils.NonEmptyEnv("AWS_REGION")
	engagementTable                = utils.NonEmptyEnv("USER_ENGAGEMENTS_TABLE_NAME")
	adHocHolidaysTable             = utils.NonEmptyEnv("HOLIDAYS_AD_HOC_TABLE")
	adHocHolidaysPlatformDateIndex = utils.NonEmptyEnv("HOLIDAYS_PLATFORM_DATE_INDEX")

	sns                   = awsutils.NewSNS(region, "", namespace)
	d                     = awsutils.NewDynamo(region, "", namespace)
	dns                   = common.DynamoNamespace{Dynamo: d, Namespace: namespace}
	adHocHolidaysTableDao = eholidays.NewDAO(&dns, adHocHolidaysTable, adHocHolidaysPlatformDateIndex)

	platform = utils.Platform{
		Sns:                       *sns,
		PlatformNotificationTopic: platformNotificationTopic,
		Namespace:                 namespace,
		IsInteractiveDebugEnabled: false,
	}
	userProfileLambda = utils.UserProfileLambda{
		Namespace:         namespace,
		ProfileLambdaName: utils.NonEmptyEnv("USER_PROFILE_LAMBDA_NAME"),
		Region:            region,
	}
	// community
	communityUsersTable              = utils.NonEmptyEnv("COMMUNITY_USERS_TABLE_NAME")
	communityUsersUserCommunityIndex = utils.NonEmptyEnv("COMMUNITY_USERS_USER_COMMUNITY_INDEX")
	communityUsersUserIndex          = utils.NonEmptyEnv("COMMUNITY_USERS_USER_INDEX")
	// TODO: use DAO for the query
	//communityUserDAO = communityUser.NewDAOFromSchema(d, namespace, schema)
)
