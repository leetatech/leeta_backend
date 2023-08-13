package infrastructure

import (
	"context"
	"github.com/leetatech/leeta_backend/services/auth/domain"
	"github.com/leetatech/leeta_backend/services/library/leetError"
	"github.com/leetatech/leeta_backend/services/library/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

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

func (a authStoreHandler) CreateVendor(vendor models.Vendor) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := a.col(models.VendorCollectionName).InsertOne(ctx, vendor)
	if err != nil {
		return err
	}
	return nil
}

func (a authStoreHandler) CreateIdentity(identity models.Identity) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := a.col(models.IdentityCollectionName).InsertOne(ctx, identity)
	if err != nil {
		return err
	}
	return nil
}

func (a authStoreHandler) GetVendorByEmail(email string) (*models.Vendor, error) {
	vendor := &models.Vendor{}
	filter := bson.M{
		"email.address": email,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := a.col(models.VendorCollectionName).FindOne(ctx, filter).Decode(vendor)
	if err != nil {
		return nil, err
	}

	return vendor, nil
}

func (a authStoreHandler) CreateOTP(verification models.Verification) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := a.col(models.VerificationsCollectionName).InsertOne(ctx, verification)
	if err != nil {
		return err
	}
	return nil
}

func (a authStoreHandler) EarlyAccess(earlyAccess models.EarlyAccess) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := a.col(models.EarlyAccessCollectionName).InsertOne(ctx, earlyAccess)
	if err != nil {
		return err
	}
	return nil
}

func (a authStoreHandler) GetIdentityByCustomerID(id string) (*models.Identity, error) {
	identity := &models.Identity{}
	filter := bson.M{
		"customer_id": id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := a.col(models.IdentityCollectionName).FindOne(ctx, filter).Decode(identity)
	if err != nil {
		return nil, err
	}

	return identity, nil
}

func (a authStoreHandler) GetOTPForValidation(target string) (*models.Verification, error) {
	var verification models.Verification

	filter := bson.M{"target": target}
	option := options.FindOneOptions{
		Sort: bson.M{"_id": -1},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := a.col(models.VerificationsCollectionName).FindOne(ctx, filter, &option).Decode(&verification)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return nil, leetError.ErrorResponseBody(leetError.DatabaseNoRecordError, err)
		}
		return nil, leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return &verification, nil
}

func (a authStoreHandler) ValidateOTP(verificationId string) error {
	filter := bson.M{"id": verificationId}
	update := bson.M{"$set": bson.M{"validated": true}}
	_, err := a.col(models.VerificationsCollectionName).UpdateOne(context.Background(), filter, update)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}

	return nil
}

func (a authStoreHandler) UpdateCredential(customerID, password string) error {
	filter := bson.M{"customer_id": customerID, "credentials.type": string(models.CredentialsTypeLogin)}
	update := bson.M{"$set": bson.M{"credentials.$.password": password, "credentials.$.status": models.CredentialStatusActive, "credentials.$.update_ts": time.Now().Unix()}}
	_, err := a.col(models.IdentityCollectionName).UpdateOne(context.Background(), filter, update)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}
	return nil
}

func (a authStoreHandler) CreateAdmin(admin models.Admin) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := a.col(models.AdminCollectionName).InsertOne(ctx, admin)
	if err != nil {
		return err
	}
	return nil
}

func (a authStoreHandler) GetAdminByEmail(email string) (*models.Admin, error) {
	admin := &models.Admin{}
	filter := bson.M{
		"email": email,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := a.col(models.AdminCollectionName).FindOne(ctx, filter).Decode(admin)
	if err != nil {
		return nil, err
	}

	return admin, nil
}
