// Package service provides service functions for appointment module
package service

import (
	"net/http"
	"time"

	"github.com/AlladinDev/AlShifa/appointment/interfaces"
	"github.com/AlladinDev/AlShifa/appointment/models"
	"github.com/AlladinDev/AlShifa/constants"
	"github.com/AlladinDev/AlShifa/utils"

	"context"

	"github.com/AlladinDev/AlShifa/structs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo          interfaces.IRepository
	clinicService interfaces.IClinicModule
}

func NewService(repo interfaces.IRepository, clinicService interfaces.IClinicModule) *Service {
	return &Service{
		repo:          repo,
		clinicService: clinicService,
	}
}

var _ interfaces.IService = (*Service)(nil)

func (s *Service) AddAppointment(ctx context.Context, appointmentDetails models.Appointment) (int, *structs.IAppError) {
	//call clinicservice to get clinic and doctor details
	doctorName, clinicName, clinicAddress, maxAppointments, err := s.clinicService.ClinicDoctorDetails(ctx, appointmentDetails.ClinicID, appointmentDetails.DoctorID, appointmentDetails.AppointmentDate)
	if err != nil {
		return 0, err
	}

	//now update actual doctorName clinicName clinicAddress to appointment payload because we cant trust frontend it can send wrong names also so we just trust clinicid and doctorid and names we fetch ourselves
	appointmentDetails.DoctorName = doctorName
	appointmentDetails.ClinicName = clinicName
	appointmentDetails.ClinicAddress = clinicAddress

	//now add some default things like status createdAt
	appointmentDetails.Status = constants.StatusAppointmentPending
	appointmentDetails.CreatedAt = time.Now()
	appointmentDetails.ID = primitive.NewObjectID()

	//now call the repo to save data
	slotNumber, appointmentSavingErr := s.repo.AddAppointment(ctx, maxAppointments, appointmentDetails)
	if appointmentSavingErr != nil {
		return 0, &structs.IAppError{
			Message:    "Failed to book appointment",
			Reason:     appointmentSavingErr.Error(),
			ErrorObj:   appointmentSavingErr,
			StatusCode: http.StatusInternalServerError,
		}
	}

	return slotNumber, nil
}

// FetchAppointments retrieves appointments based on the provided filters.
// It performs role-based validation and enrichment of filters before querying the database.
//
// Behavior:
// This function enforces access control rules depending on the user role supplied in `filters`.
// It ensures that only authorized users can fetch appointments for a clinic.
//
// Supported Roles & Logic:
//
// 1. Clinic Owner:
//   - Expects `userID` (ownerID) and `clinicID` in filters.
//   - दोनों values are parsed into MongoDB ObjectIDs.
//   - Validates that the given clinic actually belongs to the provided owner.
//     This prevents unauthorized access where an owner might try to query another clinic’s data.
//   - If validation succeeds, the verified clinicID is reattached to filters.
//   - If validation fails, an error is returned.
//
// 2. Clinic Receptionist:
//   - Expects `userID` (receptionistID) in filters.
//   - Parses the receptionistID into MongoDB ObjectID.
//   - Fetches the clinicID associated with this receptionist from the database.
//     This ensures the receptionist can only access appointments of their assigned clinic.
//   - The derived clinicID is injected into filters.
//
// 3. Other Users:
//   - No additional validation is performed here.
//   - Assumes authentication/authorization is handled upstream (e.g., via request headers).
//
// Important Notes:
// - `userType` is only used for internal validation and is removed from filters before querying DB.
// - `userID` from request headers is considered trusted and is not re-validated here.
// - The function ensures that clinic-based access is strictly enforced for sensitive roles.
//
// Parameters:
// - ctx: Context for request lifecycle and cancellation.
// - filters: A BSON map containing query filters. Expected keys:
//     - "userType": Role of the user (e.g., clinicOwner, receptionist)
//     - "userID": ID of the requesting user
//     - "clinicID": (required for clinicOwner)
//
// Returns:
// - []sharedModels.Appointment: List of matching appointments.
// - *structs.IAppError: Structured error in case of failure.
//
// Error Handling:
// - Returns internal server errors for:
//     - Invalid ObjectID parsing
//     - Failed ownership validation
//     - Database query failures
//
// Design Rationale:
// - Centralizes role-based access validation inside the service layer.
// - Prevents unauthorized data access by validating ownership relationships.
// - Keeps repository layer clean by ensuring only valid, pre-processed filters reach it.
// - Improves maintainability by isolating role-specific logic in a single place.
func (s *Service) FetchAppointments(ctx context.Context, filters bson.M) ([]models.Appointment, *structs.IAppError) {
	//here we need to validate userid if usertype passed in filters is receptionist we need to grab this receptionist from db and get its clinicid
	//similary if usertype is clinicowner in filters there will be  ownerid and also clinicID we need to verify whether this owner owns this clinic or not
	//so these validations need to be done for user we can get its id from req.header we dont need to verify it , only for receptionist and clinic owner we need to further verify
	userType := filters["userType"]
	userID := filters["userID"]
	clinicIDPassed := filters["clinicID"]
	switch userType {
	case constants.RoleclinicOwner:
		//parse userid into mongodb format as it is passed as any
		idParsingErr, userMongoID := utils.ParseUserID(userID)
		if idParsingErr != nil {
			return nil, &structs.IAppError{
				Message:    "failed to fetch appointments",
				StatusCode: http.StatusInternalServerError,
				Reason:     "Failed to convert userid into mongodb format",
				ErrorObj:   idParsingErr,
			}
		}

		//parse clinicid as it will be passed as any so parse it into mongodb format
		clinicIDParsingErr, clinicMongodbID := utils.ParseUserID(clinicIDPassed)
		if clinicIDParsingErr != nil {
			return nil, &structs.IAppError{
				Message:    "failed to fetch appointments",
				StatusCode: http.StatusInternalServerError,
				Reason:     "Failed to convert clinic ID into mongodb format",
				ErrorObj:   idParsingErr,
			}
		}

		//now if user type is clinic owner search appointments by both ownerid and clinicid
		clinicID, idErr := s.clinicService.GetClinicIDIfExists(ctx, bson.M{"ownerID": userMongoID, "_id": clinicMongodbID})
		if idErr != nil {
			return nil, &structs.IAppError{
				Message:    "failed to fetch appointments",
				StatusCode: http.StatusInternalServerError,
				Reason:     "Failed to convert userid into mongodb format",
				ErrorObj:   idErr,
			}
		}

		filters["clinicID"] = clinicID

		//case to handle if user type is receptionist in this case using receptionist id get its clinicid and then get appointments using that clinicid
	case constants.RoleClinicReceptionist:
		receptionistIDErr, receptionistID := utils.ParseUserID(userID)
		if receptionistIDErr != nil {
			return nil, &structs.IAppError{
				Message:    "failed to fetch appointments",
				StatusCode: http.StatusInternalServerError,
				Reason:     "failed to convert receptionist id into mongodb format",
				ErrorObj:   receptionistIDErr,
			}
		}
		clinicID, idErr := s.clinicService.GetClinicIDByReceptionist(ctx, receptionistID)
		if idErr != nil {
			return nil, &structs.IAppError{
				Message:    "failed to fetch appointments",
				StatusCode: http.StatusInternalServerError,
				Reason:     idErr.Message,
				ErrorObj:   idErr,
			}
		}
		filters["clinicID"] = clinicID

	case constants.RoleDoctor:
		filters["doctorID"] = userID
	}

	//now delete the userType it is no longer needed now
	delete(filters, "userType")

	appointments, err := s.repo.FetchAppointments(ctx, filters)
	if err != nil {
		return nil, &structs.IAppError{
			Message:    "failed to fetch appointments",
			StatusCode: http.StatusInternalServerError,
			Reason:     err.Error(),
			ErrorObj:   err,
		}
	}

	return appointments, nil
}

func (s *Service) UpdateAppointmentStatus(ctx context.Context, appointmentID primitive.ObjectID, status bool) *structs.IAppError {
	if err := s.repo.UpdateAppointmentStatus(ctx, appointmentID, status); err != nil {
		return &structs.IAppError{
			Message:    "Failed to update appointment",
			Reason:     err.Error(),
			StatusCode: http.StatusInternalServerError,
			ErrorObj:   err,
		}
	}

	return nil
}

func (s *Service) FetchAppointmentDaysBooked(ctx context.Context, doctorID primitive.ObjectID, clinicID primitive.ObjectID) ([]models.Slot, *structs.IAppError) {
	//here first whether this clinic exists or not
	clinicSearchErr := s.clinicService.ClinicExists(ctx, clinicID)
	if clinicSearchErr != nil {
		return nil, clinicSearchErr
	}

	//now check whether this doctor exists or not
	doctorSearchErr := s.clinicService.DoctorExists(ctx, doctorID)
	if doctorSearchErr != nil {
		return nil, doctorSearchErr
	}

	//now check whether this doctor clinic mapping exists or not
	mappingErr := s.clinicService.DoctorClinicMappingExists(ctx, clinicID, doctorID)
	if mappingErr != nil {
		return nil, mappingErr
	}

	//now fetch maxAppointments for this clinic
	maxAppointments, err := s.clinicService.FetchMaxAppointments(ctx, clinicID)
	if err != nil {
		return nil, err
	}

	daysBooked, daysBookedErr := s.repo.FetchAppointmentDaysBooked(ctx, maxAppointments, doctorID, clinicID)
	if daysBookedErr != nil {
		return nil, &structs.IAppError{
			Message:    "Failed to fetch appointments days booked",
			StatusCode: http.StatusInternalServerError,
			Reason:     daysBookedErr.Error(),
			ErrorObj:   daysBookedErr,
		}
	}

	return daysBooked, nil
}
