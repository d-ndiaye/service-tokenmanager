package userToken

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"service-token/pkg/config"
	"strconv"
)

const (
	userTokenCollection string = "userToken"
)

type Repository interface {
	FindByUsername(username string) (UserToken, error)
	Save(userToken UserToken) (UserToken, error)
	Ping() error
	GetByTenantAndCategory(tenantId string, category string) ([]UserToken, error)
	Delete(deviceToken string) error
}

type userTokenRepository struct {
	configMongo config.Mongodb
	client      *mongo.Client
}

func (tr userTokenRepository) Ping() error {
	return tr.client.Ping(context.TODO(), readpref.Primary())
}

func NewRepository(mongodb config.Mongodb) (Repository, error) {
	fmt.Println("mongodb://" + mongodb.Host + ":" + strconv.Itoa(mongodb.Port))
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://"+mongodb.Host+":"+strconv.Itoa(mongodb.Port)))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}
	r := userTokenRepository{
		client:      client,
		configMongo: mongodb,
	}
	return r, nil
}

// GetByTenantAndCategory Retrieves all the userTokens with the same tenant and category
// return []UserToken, nil
func (tr userTokenRepository) GetByTenantAndCategory(tenantId string, category string) ([]UserToken, error) {
	userTokenCollection := tr.client.Database(tr.configMongo.Dbname).Collection(userTokenCollection)
	var userTokens []UserToken
	filter := bson.M{
		"deviceToken.tenantNotificationConfig.tenantId": tenantId,
		"deviceToken.tenantNotificationConfig.category": bson.M{"$in": bson.A{category}},
	}
	result, err := userTokenCollection.Find(context.Background(), filter)
	if err != nil {
		return userTokens, err
	}
	defer result.Close(context.Background())
	for result.Next(context.Background()) {
		var ut UserToken
		err := result.Decode(&ut)
		if err != nil {
			return userTokens, err
		}
		userTokens = append(userTokens, ut)
	}
	if err := result.Err(); err != nil {
		return userTokens, err
	}
	return userTokens, nil
}

func (tr userTokenRepository) FindByUsername(username string) (UserToken, error) {
	tokenCollection := tr.client.Database(tr.configMongo.Dbname).Collection(userTokenCollection)
	userToken := UserToken{}
	result := tokenCollection.FindOne(context.TODO(), bson.M{"username": username})
	err := result.Decode(&userToken)
	if err != nil {
		return userToken, err
	}
	return userToken, nil
}

func (tr userTokenRepository) Delete(DeviceToken string) error {
	collection := tr.client.Database(tr.configMongo.Dbname).Collection(userTokenCollection)
	filter := bson.M{"deviceToken.id": DeviceToken}

	update := bson.M{
		"$pull": bson.M{"deviceToken": bson.M{"id": DeviceToken}},
	}
	_, err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (tr userTokenRepository) Save(userToken UserToken) (UserToken, error) {
	tokenCollection := tr.client.Database(tr.configMongo.Dbname).Collection(userTokenCollection)
	// Check if username already exists in database
	filter := bson.M{"username": userToken.Username}
	existingUserToken := UserToken{}
	err := tokenCollection.FindOne(context.Background(), filter).Decode(&existingUserToken)
	if err == nil {
		// username already exists, check for updates to DeviceToken or TenantNotificationConfig
		for _, deviceToken := range userToken.DeviceToken {
			existingDeviceTokenIndex := -1
			for i, existingDeviceToken := range existingUserToken.DeviceToken {
				if existingDeviceToken.Id == deviceToken.Id {
					existingDeviceTokenIndex = i
					break
				}
			}
			if existingDeviceTokenIndex == -1 {
				// DeviceToken doesn't exist yet, add it to the user's device tokens
				existingUserToken.DeviceToken = append(existingUserToken.DeviceToken, deviceToken)
			} else {
				// DeviceToken already exists, update the TenantNotificationConfig if necessary
				for _, tnc := range deviceToken.TenantNotificationConfig {
					existingTncIndex := -1
					for j, existingTNC := range existingUserToken.DeviceToken[existingDeviceTokenIndex].TenantNotificationConfig {
						if existingTNC.TenantId == tnc.TenantId {
							existingTncIndex = j
							break
						}
					}
					if existingTncIndex == -1 {
						// TenantNotificationConfig doesn't exist yet, add it to the device token's configs
						existingUserToken.DeviceToken[existingDeviceTokenIndex].TenantNotificationConfig = append(existingUserToken.DeviceToken[existingDeviceTokenIndex].TenantNotificationConfig, tnc)
					} else {
						// TenantNotificationConfig already exists, update its category
						existingUserToken.DeviceToken[existingDeviceTokenIndex].TenantNotificationConfig[existingTncIndex].Category = tnc.Category
					}
				}
			}
		}
		// Update existing userToken in the database
		update := bson.M{"$set": bson.M{"deviceToken": existingUserToken.DeviceToken}}
		_, err = tokenCollection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return UserToken{}, err
		}
		return existingUserToken, nil
	}
	// username does not exist, insert new userToken into the database
	_, err = tokenCollection.InsertOne(context.Background(), userToken)
	if err != nil {
		return UserToken{}, err
	}
	return userToken, nil
}
