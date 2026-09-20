// package service

// import (
// 	"context"
// 	"net/http"
// 	"time"

// 	"github.com/AlladinDev/AlShifa/internal/appointment/interfaces"
// 	"github.com/AlladinDev/AlShifa/internal/appointment/models"
// 	"github.com/AlladinDev/AlShifa/internal/shared/constants"
// 	appInterfaces "github.com/AlladinDev/AlShifa/internal/shared/interfaces"
// 	"github.com/AlladinDev/AlShifa/internal/shared/response"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type AppointmentOrder struct {
// 	repo           interfaces.IRepository
// 	paymentGateway appInterfaces.IPaymentGateway
// 	planService    interfaces.IPlanService
// }

// func NewAppointmentOrderService(Repo interfaces.IRepository, PlanService interfaces.IPlanService, PaymentGateway appInterfaces.IPaymentGateway) interfaces.IAppointmentOrder {
// 	return &AppointmentOrder{
// 		repo:           Repo,
// 		planService:    PlanService,
// 		paymentGateway: PaymentGateway,
// 	}
// }
// func (r *AppointmentOrder) CreateOrder(ctx context.Context, appointmentID primitive.ObjectID, userID primitive.ObjectID, planId primitive.ObjectID) (appointmentOrder *models.AppointmentBookingPaymentOrder, err *response.IAppError) {
// 	//first using planid call the plan service to get the amount to deduct
// 	amountToDeduct, planErr := r.planService.FetchAmountToDeduct(ctx, planId)
// 	if planErr != nil {
// 		return nil, &response.IAppError{
// 			Message:    "Failed to create appointment order",
// 			Reason:     planErr.Error(),
// 			ErrorObj:   planErr,
// 			StatusCode: http.StatusInternalServerError,
// 		}
// 	}

// 	//now as we have amount to deduct call the payment gateway to create an order and get the provider id
// 	providerID, providerName := r.paymentGateway.CreateOrder(ctx, amountToDeduct, constants.IndianCurrency)
// 	order := models.AppointmentBookingPaymentOrder{
// 		ID:              primitive.NewObjectID(),
// 		Amount:          amountToDeduct,
// 		PaymentStatus:   constants.AppointmentBookingOrderCreated,
// 		CreatedAt:       time.Now(),
// 		ProviderOrderID: providerID,
// 		Provider:        providerName,
// 		Currency:        constants.IndianCurrencyCode,
// 		AppointmentID:   appointmentID,
// 	}

// }
package service