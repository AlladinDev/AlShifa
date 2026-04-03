// Package repository provides repository functions for repository
package repository

import (
	"github.com/AlladinDev/AlShifa/appointment/interfaces"
	"github.com/AlladinDev/AlShifa/appointment/models"
	sharedModels "github.com/AlladinDev/AlShifa/models"

	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		db: db,
	}
}

var _ interfaces.IRepository = (*Repository)(nil)

func (r *Repository) AddAppointment(ctx context.Context, clinicMaxAppointments int, appointmentDetails sharedModels.Appointment) (int, error) {

	//here first update slot document using upsert to ensure if it is not present create it if present updates its slots booked by 1
	session, err := r.db.Client().StartSession()
	if err != nil {
		return 0, err
	}
	defer session.EndSession(ctx)

	//this is the transaction function which will contain db operations logic which must be executed in single operation
	transactionFn := func(sessCtx mongo.SessionContext) (any, error) {
		slotFilter := bson.D{
			{Key: "bookingDate", Value: appointmentDetails.AppointmentDate},
			{Key: "doctorID", Value: appointmentDetails.DoctorID},
			{Key: "clinicID", Value: appointmentDetails.ClinicID},
			{Key: "slotsBooked", Value: bson.M{
				"$lt": clinicMaxAppointments,
			}},
		}

		updateQuery := bson.D{
			{Key: "$inc", Value: bson.D{
				{Key: "slotsBooked", Value: 1},
			}},
			{Key: "$setOnInsert", Value: bson.D{
				{Key: "bookingDate", Value: appointmentDetails.AppointmentDate},
				{Key: "doctorID", Value: appointmentDetails.DoctorID},
				{Key: "clinicID", Value: appointmentDetails.ClinicID},
			}},
		}

		//create the update options to update slot first then get its slotsBooked and assign it to current appointment as slot
		updateOptions := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
		slotUpdationRes := r.db.Collection("Slot").FindOneAndUpdate(sessCtx, slotFilter, updateQuery, updateOptions)
		var slotUpdated models.Slot
		if err := slotUpdationRes.Decode(&slotUpdated); err != nil {
			return nil, err
		}

		//here get the updated slotsBooked and set it in appointment details it will represent the slot for this appointment
		appointmentDetails.Slot = slotUpdated.SlotsBooked

		//now save the appointment also
		_, err := r.db.Collection("Appointment").InsertOne(sessCtx, appointmentDetails)
		if err != nil {
			return nil, err
		}

		return appointmentDetails, nil
	}

	res, transactionErr := session.WithTransaction(ctx, transactionFn)
	if transactionErr != nil {
		return 0, transactionErr
	}

	//now try to convert res into appointment model
	appointmentCreated, ok := res.(sharedModels.Appointment)
	if !ok {
		return 0, errors.New("failed to return appointment Created")
	}

	return appointmentCreated.Slot, nil
}

func (r *Repository) FetchAppointments(ctx context.Context, filters bson.M) ([]sharedModels.Appointment, error) {
	cursor, err := r.db.Collection("Appointment").Find(ctx, filters)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var appointments []sharedModels.Appointment
	if err := cursor.All(ctx, &appointments); err != nil {
		return nil, err
	}

	return appointments, nil
}

func (r *Repository) UpdateAppointmentStatus(ctx context.Context, appointmentID primitive.ObjectID, status bool) error {
	result, err := r.db.Collection("Appointment").UpdateOne(ctx, bson.M{"_id": appointmentID}, bson.M{"$set": bson.M{"status": status}})
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("this appointment doesnt exist")
	}
	if result.ModifiedCount == 0 {
		return errors.New("failed to Update Appointment")
	}

	return nil
}

func (r *Repository) FetchAppointmentDaysBooked(ctx context.Context, maxAppointments int, doctorID primitive.ObjectID, clinicID primitive.ObjectID) ([]models.Slot, error) {
	cursor, err := r.db.Collection("Slot").Find(ctx, bson.M{"clinicID": clinicID, "doctorID": doctorID, "slotsBooked": bson.M{"$eq": maxAppointments}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var appointmentDaysBooked []models.Slot
	if err := cursor.All(ctx, &appointmentDaysBooked); err != nil {
		return nil, err
	}

	return appointmentDaysBooked, nil
}
