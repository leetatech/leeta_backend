package infrastructure

import (
	"context"
	"github.com/leetatech/leeta_backend/services/auth/domain"
	"github.com/leetatech/leeta_backend/services/library/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
