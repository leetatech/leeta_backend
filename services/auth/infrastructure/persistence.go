package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/leetatech/leeta_backend/services/auth/domain"
	"github.com/leetatech/leeta_backend/services/library/leetError"
	"github.com/leetatech/leeta_backend/services/library/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

var ErrItemNotFound = errors.New("item not found")

type authStoreHandler struct {
	client       *mongo.Client
	databaseName string
	logger       *zap.Logger
}

func (a authStoreHandler) col(collectionName string) *mongo.Collection {
	return a.client.Database(a.databaseName).Collection(collectionName)
}

func NewAuthPersistence(client *mongo.Client, databaseName string, logger *zap.Logger) domain.AuthRepository {
	return &authStoreHandler{client: client, databaseName: databaseName, logger: logger}
}

func (a authStoreHandler) CreateIdentity(ctx context.Context, identity models.Identity) error {
	_, err := a.col(models.IdentityCollectionName).InsertOne(ctx, identity)
	if err != nil {
		return err
	}
	return nil
}

func (a authStoreHandler) GetVendorByEmail(ctx context.Context, email string) (*models.Vendor, error) {
	vendor := &models.Vendor{}
	filter := bson.M{
		EmailAddress: email,
	}

	err := a.col(models.UsersCollectionName).FindOne(ctx, filter).Decode(vendor)
	if err != nil {
		return nil, err
	}

	return vendor, nil
}

func (a authStoreHandler) CreateOTP(ctx context.Context, verification models.Verification) error {
	_, err := a.col(models.VerificationsCollectionName).InsertOne(ctx, verification)
	if err != nil {
		return err
	}
	return nil
}

func (a authStoreHandler) EarlyAccess(ctx context.Context, earlyAccess models.EarlyAccess) error {
	_, err := a.col(models.EarlyAccessCollectionName).InsertOne(ctx, earlyAccess)
	if err != nil {
		return err
	}
	return nil
}

func (a authStoreHandler) GetIdentityByCustomerID(ctx context.Context, id string) (*models.Identity, error) {
	identity := &models.Identity{}
	filter := bson.M{
		"user_id": id,
	}

	err := a.col(models.IdentityCollectionName).FindOne(ctx, filter).Decode(identity)
	if err != nil {
		return nil, err
	}

	return identity, nil
}

func (a authStoreHandler) GetOTPForValidation(ctx context.Context, target string) (*models.Verification, error) {
	var verification models.Verification

	filter := bson.M{"target": target, "validated": false}
	option := options.FindOneOptions{
		Sort: bson.M{"_id": -1},
	}

	err := a.col(models.VerificationsCollectionName).FindOne(ctx, filter, &option).Decode(&verification)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return nil, ErrItemNotFound
		}
		return nil, fmt.Errorf("error finding otp for validation %w", err)
	}

	return &verification, nil
}

func (a authStoreHandler) ValidateOTP(ctx context.Context, verificationId string) error {
	filter := bson.M{"id": verificationId}
	update := bson.M{"$set": bson.M{"validated": true}}
	_, err := a.col(models.VerificationsCollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return nil
}

func (a authStoreHandler) UpdateCredential(ctx context.Context, customerID, password string) error {
	filter := bson.M{"customer_id": customerID, "credentials.type": string(models.CredentialsTypeLogin)}
	update := bson.M{"$set": bson.M{"credentials.$.password": password, "credentials.$.status": models.CredentialStatusActive, "credentials.$.update_ts": time.Now().Unix()}}
	_, err := a.col(models.IdentityCollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}
	return nil
}

func (a authStoreHandler) GetAdminByEmail(ctx context.Context, email string) (*models.Admin, error) {
	admin := &models.Admin{}
	filter := bson.M{
		"email": email,
	}

	err := a.col(models.UsersCollectionName).FindOne(ctx, filter).Decode(admin)
	if err != nil {
		return nil, err
	}

	return admin, nil
}

func (a authStoreHandler) CreateUser(ctx context.Context, user any) error {
	_, err := a.col(models.UsersCollectionName).InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (a authStoreHandler) GetUserByEmail(ctx context.Context, email string) (*models.Customer, error) {
	customer := &models.Customer{}
	filter := bson.M{
		EmailAddress: email,
	}

	err := a.col(models.UsersCollectionName).FindOne(ctx, filter).Decode(customer)
	if err != nil {
		return nil, err
	}

	return customer, nil
}
