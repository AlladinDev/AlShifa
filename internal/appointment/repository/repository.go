// Package repository provides repository functions for repository
package repository

import (
	"github.com/AlladinDev/AlShifa/internal/appointment/interfaces"
	"github.com/AlladinDev/AlShifa/internal/appointment/models"

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

//AddAppointment adds appointment and returns slot document id ,slot booked and error slot document id is required so that in case if appointment gets cancelled slot can be freed back using slot document id
func (r *Repository) AddAppointment(ctx context.Context, clinicMaxAppointments int, appointmentDetails models.Appointment) (slotDocumentID primitive.ObjectID, slot int, err error) {
	//here first update slot document using upsert to ensure if it is not present create it if present updates its slots booked by 1
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
			{Key: "token", Value: 1},
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

	//from here start the transaction
	session, sessionErr := r.db.Client().StartSession()
	if sessionErr != nil {
		return primitive.NilObjectID, 0, sessionErr
	}

	type slotResult struct {
		Slot           int                `json:"slot"`
		SlotDocumentID primitive.ObjectID `json:"slotDocumentID"`
	}

	transactionFn := func(mongoCtx mongo.SessionContext) (any, error) {
		slotObject := slotResult{}
		slotUpdationRes := r.db.Collection("Slot").FindOneAndUpdate(mongoCtx, slotFilter, updateQuery, updateOptions)
		var slotUpdated models.Slot
		if err := slotUpdationRes.Decode(&slotUpdated); err != nil {
			return slotObject, err
		}

		//here get the updated slotsBooked and set it in appointment details it will represent the slot for this appointment
		appointmentDetails.Slot = slotUpdated.Token

		//now save the appointment also
		_, err := r.db.Collection("Appointment").InsertOne(mongoCtx, appointmentDetails)
		if err != nil {
			return slotObject, err
		}

		//now as everything is correct update the slotResult with slotDoID and also slot
		slotObject.Slot = slotUpdated.SlotsBooked
		slotObject.SlotDocumentID = slotUpdated.ID
		return slotObject, nil
	}

	res, err := session.WithTransaction(ctx, transactionFn)
	if err != nil {
		return primitive.NilObjectID, 0, err
	}

	//now res is any try to decode it into slot doc result to get slot and also slot document id
	transactionSlotRes, ok := res.(slotResult)
	if !ok {
		return primitive.NilObjectID, 0, errors.New("failed to parse slot result into its desired form")
	}

	return transactionSlotRes.SlotDocumentID, transactionSlotRes.Slot, nil

}

func (r *Repository) FetchAppointments(ctx context.Context, filters bson.M) ([]models.Appointment, error) {

	cursor, err := r.db.Collection("Appointment").Find(ctx, filters)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var appointments []models.Appointment
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
