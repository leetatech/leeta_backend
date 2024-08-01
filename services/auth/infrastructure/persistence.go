package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/leetatech/leeta_backend/pkg/errs"
	"github.com/leetatech/leeta_backend/services/auth/domain"
	"github.com/leetatech/leeta_backend/services/dtos"
	"github.com/leetatech/leeta_backend/services/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var ErrItemNotFound = errors.New("item not found")

type authStoreHandler struct {
	client       *mongo.Client
	databaseName string
}

func (a authStoreHandler) col(collectionName string) *mongo.Collection {
	return a.client.Database(a.databaseName).Collection(collectionName)
}

func New(client *mongo.Client, databaseName string) domain.AuthRepository {
	return &authStoreHandler{client: client, databaseName: databaseName}
}

func (a authStoreHandler) CreateIdentity(ctx context.Context, identity models.Identity) error {
	_, err := a.col(models.IdentityCollectionName).InsertOne(ctx, identity)
	if err != nil {
		return err
	}
	return nil
}

func (a authStoreHandler) CreateGuestRecord(ctx context.Context, guest models.Guest) error {
	_, err := a.col(models.GuestsCollectionName).InsertOne(ctx, guest)
	if err != nil {
		return err
	}
	return nil
}

func (a authStoreHandler) VendorByEmail(ctx context.Context, email string) (*models.Vendor, error) {
	vendor := &models.Vendor{}
	filter := bson.M{
		dtos.EmailAddress: email,
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

func (a authStoreHandler) SaveEarlyAccess(ctx context.Context, earlyAccess models.EarlyAccess) error {
	_, err := a.col(models.EarlyAccessCollectionName).InsertOne(ctx, earlyAccess)
	if err != nil {
		return err
	}
	return nil
}

func (a authStoreHandler) IdentityByUserID(ctx context.Context, id string) (*models.Identity, error) {
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

func (a authStoreHandler) FindUnvalidatedVerificationByTarget(ctx context.Context, target string) (*models.Verification, error) {
	var verification models.Verification

	filter := bson.M{"target": target, "validated": false}
	filterOptions := options.FindOne().SetSort(bson.M{"_id": -1})

	err := a.col(models.VerificationsCollectionName).FindOne(ctx, filter, filterOptions).Decode(&verification)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrItemNotFound
		}
		return nil, fmt.Errorf("error finding unvalidated verification for target %s: %w", target, err)
	}

	return &verification, nil
}

func (a authStoreHandler) ValidateOTP(ctx context.Context, verificationId string) error {
	filter := bson.M{"id": verificationId}
	update := bson.M{"$set": bson.M{"validated": true}}
	_, err := a.col(models.VerificationsCollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return errs.Body(errs.DatabaseError, err)
	}

	return nil
}

func (a authStoreHandler) UpdateCredential(ctx context.Context, userID, password string) error {
	filter := bson.M{dtos.UserID: userID, dtos.CredentialsType: string(models.CredentialsTypeLogin)}
	update := bson.M{"$set": bson.M{"credentials.$.password": password, "credentials.$.status": models.CredentialStatusActive, "credentials.$.update_ts": time.Now().Unix()}}
	_, err := a.col(models.IdentityCollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return errs.Body(errs.DatabaseError, err)
	}
	return nil
}

func (a authStoreHandler) AdminByEmail(ctx context.Context, email string) (*models.Admin, error) {
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

func (a authStoreHandler) UserByEmail(ctx context.Context, email string) (*models.Customer, error) {
	customer := &models.Customer{}
	filter := bson.M{
		dtos.EmailAddress: email,
	}

	err := a.col(models.UsersCollectionName).FindOne(ctx, filter).Decode(customer)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (a authStoreHandler) SetEmailVerificationStatus(ctx context.Context, email string, status bool) error {
	filter := bson.M{dtos.EmailAddress: email}
	update := bson.M{"$set": bson.M{dtos.EmailVerifiedStatus: status}}
	_, err := a.col(models.UsersCollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return errs.Body(errs.DatabaseError, err)
	}
	return nil
}

func (a authStoreHandler) GuestRecord(ctx context.Context, deviceId string) (guest models.Guest, err error) {
	filter := bson.M{
		dtos.DeviceId: deviceId,
	}

	err = a.col(models.GuestsCollectionName).FindOne(ctx, filter).Decode(&guest)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			err = ErrItemNotFound
		}
		return
	}

	return
}

func (a authStoreHandler) UpdateGuestRecord(ctx context.Context, guest models.Guest) error {
	filter := bson.M{
		dtos.DeviceId: guest.DeviceID,
	}

	update := bson.M{
		"$set": guest,
	}

	_, err := a.col(models.GuestsCollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (a authStoreHandler) GetUserByEmailOrPhone(ctx context.Context, target string) (*models.Customer, error) {
	customer := &models.Customer{}

	filter := bson.M{"$or": []bson.M{
		{dtos.EmailAddress: target},
		{dtos.PhoneNumber: target},
	}}

	err := a.col(models.UsersCollectionName).FindOne(ctx, filter).Decode(customer)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (a authStoreHandler) UpdatePhoneVerify(ctx context.Context, phone string, status bool) error {
	filter := bson.M{dtos.PhoneNumber: phone}
	update := bson.M{"$set": bson.M{dtos.PhoneVerificationStatus: status}}
	_, err := a.col(models.UsersCollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return leetError.ErrorResponseBody(leetError.DatabaseError, err)
	}
	return nil
}
