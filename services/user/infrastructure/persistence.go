package infrastructure

import (
	"context"
	"github.com/leetatech/leeta_backend/pkg/errs"
	"github.com/leetatech/leeta_backend/services/models"
	"github.com/leetatech/leeta_backend/services/user/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type userStoreHandler struct {
	client       *mongo.Client
	databaseName string
}

func (u userStoreHandler) col(collectionName string) *mongo.Collection {
	return u.client.Database(u.databaseName).Collection(collectionName)
}

func New(client *mongo.Client, databaseName string) domain.UserRepository {
	return &userStoreHandler{client: client, databaseName: databaseName}
}

func (u userStoreHandler) VendorDetailsUpdate(request domain.VendorDetailsUpdateRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": request.ID}
	update := bson.M{"$set": bson.M{"first_name": request.FirstName, "last_name": request.LastName, "status": request.Status, "status_ts": time.Now().Unix()}}
	result, err := u.col(models.UsersCollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return errs.Body(errs.DatabaseError, err)
	}
	if result.MatchedCount == 0 {
		return errs.Body(errs.DatabaseNoRecordError, err)
	}
	return nil
}

func (u userStoreHandler) RegisterVendorBusiness(request models.Business) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := u.col(models.BusinessCollectionName).InsertOne(ctx, request)
	if err != nil {
		return err
	}
	return nil
}

func (u userStoreHandler) GetVendorByID(id string) (*models.Vendor, error) {
	vendor := &models.Vendor{}
	filter := bson.M{
		"user.id": id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := u.col(models.UsersCollectionName).FindOne(ctx, filter).Decode(vendor)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return nil, errs.Body(errs.DatabaseNoRecordError, err)

		default:
			return nil, errs.Body(errs.DatabaseError, err)
		}
	}

	return vendor, nil
}

func (u userStoreHandler) GetCustomerByID(id string) (*models.Customer, error) {
	customer := &models.Customer{}
	filter := bson.M{
		"user.id": id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := u.col(models.UsersCollectionName).FindOne(ctx, filter).Decode(customer)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return nil, errs.Body(errs.DatabaseNoRecordError, err)

		default:
			return nil, errs.Body(errs.DatabaseError, err)
		}
	}

	return customer, nil
}

func (u userStoreHandler) UpdateUserRecord(request *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"user.id": request.ID,
	}

	update := bson.M{
		"$set": bson.M{
			"user": *request,
		},
	}

	_, err := u.col(models.UsersCollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
