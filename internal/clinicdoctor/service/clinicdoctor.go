package service

import (
	"context"
	"time"

	"github.com/AlladinDev/AlShifa/internal/clinicdoctor/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TService struct {
	repo interfaces.IRepo
}

func NewService(Repo interfaces.IRepo) interfaces.IService {
	return &TService{
		repo: Repo,
	}
}

func (s TService) IsAppointmentDatePossible(ctx context.Context, appointmentDate time.Time, doctorID primitive.ObjectID, clinicID primitive.ObjectID) (bool, error) {
	return s.repo.CheckAppointmentDatePossible(ctx, clinicID, doctorID, appointmentDate)
}
