package platform

import (
	"github.com/pkg/errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/ReneKroon/ttlcache"

	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	"github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/daos/clientPlatformToken"
	"github.com/adaptiveteam/adaptive/daos/slackTeam"
	
)

// DAO - wrapper around a Dynamo DB table to work with PlatformID -> PlatformToken mapping
type DAO interface {
	Read(teamID models.TeamID) (models.ClientPlatformToken, error)
	ReadUnsafe(teamID models.TeamID) models.ClientPlatformToken
	GetPlatformTokenUnsafe(teamID models.TeamID) string
}

// DAOImpl - a container for all information needed to access a DynamoDB table
type DAOImpl struct {
	Dynamo    *awsutils.DynamoRequest `json:"dynamo"`
	Namespace string                  `json:"namespace"`
	models.ClientPlatformTokenTableSchema
}

// NewDAO creates an instance of DAO that will provide access to ClientPlatformToken table
func NewDAO(dynamo *awsutils.DynamoRequest, namespace, table string) DAO {
	if table == "" { panic(errors.New("Cannot create ClientPlatformToken DAO without table")) }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		ClientPlatformTokenTableSchema: models.ClientPlatformTokenTableSchema{Name: table},
	}
}

// NewDAOFromSchema creates an instance of DAO that will provide access to adaptiveValues table
func NewDAOFromSchema(dynamo *awsutils.DynamoRequest, namespace string, schema models.Schema) DAO {	
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		ClientPlatformTokenTableSchema: schema.ClientPlatformTokens}
}

// Read reads ClientPlatformToken
func (d DAOImpl) Read(teamID models.TeamID) (models.ClientPlatformToken, error) {
	params := map[string]*dynamodb.AttributeValue{
		"platform_id": dynString(teamID.ToString()),
	}
	var out models.ClientPlatformToken
	err := d.Dynamo.GetItemFromTable(d.Name, params, &out)
	return out, err
}

// ReadUnsafe reads the ClientPlatformToken. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(teamID models.TeamID) models.ClientPlatformToken {
	out, err := d.Read(teamID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not find %s in %s", teamID, d.Name))
	return out
}

func dynString(str string) (attr *dynamodb.AttributeValue) {
	return &dynamodb.AttributeValue{S: aws.String(str)}
}
// GetPlatformTokenUnsafe reads platform token from database
func (d DAOImpl) GetPlatformTokenUnsafe(teamID models.TeamID) string {
	return d.ReadUnsafe(teamID).PlatformToken
}

var globalTokenCache *ttlcache.Cache

// GetToken retrieves the token from the cache or database.
func GetToken(teamID models.TeamID) func (common.DynamoDBConnection) (string, error) {
	return func (conn common.DynamoDBConnection) (token string, err error) {
		if globalTokenCache != nil {
			value, existsInCache := globalTokenCache.Get(teamID.ToString())
			if existsInCache {
				token = value.(string)
			}
		}
		if token == "" {
			if teamID.TeamID != "" {
				dao := slackTeam.NewDAO(conn.Dynamo, "GetToken", conn.ClientID)
				var teams [] slackTeam.SlackTeam
				teams, err = dao.ReadOrEmpty(teamID.TeamID)
				if err != nil {
					return
				}
				if len(teams) > 0 {
					token = teams[0].AccessToken
				}
			}
			if token == "" {
				if teamID.AppID != "" {
					dao := clientPlatformToken.NewDAO(conn.Dynamo, "GetToken2", conn.ClientID)
					var teams [] clientPlatformToken.ClientPlatformToken
					teams, err = dao.ReadOrEmpty(teamID.AppID)
					if err != nil {
						return
					}
					if len(teams) > 0 {
						token = teams[0].PlatformToken
					}
				}
			}
			if token != "" {
				globalTokenCache.Set(teamID.ToString(), token)
			}
		}
		return 
	}
}

// GetTokenForUser searches token for the given user
func GetTokenForUser(dynamo *awsutils.DynamoRequest, clientID string, userID string) (token string, err error) {
	var teamID models.TeamID
	teamID, err = GetTeamIDForUser(dynamo, clientID, userID)
	conn := common.DynamoDBConnection{
		Dynamo: dynamo,
		ClientID: clientID,
		PlatformID: teamID.ToPlatformID(),
	}
	if err == nil {
		token, err = GetToken(teamID)(conn)
	}
	return 
}
// GetTeamIDForUser -
func GetTeamIDForUser(dynamo *awsutils.DynamoRequest, clientID string, userID string) (teamID models.TeamID, err error) {
	dao := user.NewDAO(dynamo, "GetTeamIDForUser", clientID)
	var user models.User
	user, err = dao.Read(userID)
	if err == nil {
		teamID = models.ParseTeamID(user.PlatformID)
	}
	return 
}

var _ = func () int {
	globalTokenCache = InitLocalCache(globalTokenCache)
	return 0
}()
