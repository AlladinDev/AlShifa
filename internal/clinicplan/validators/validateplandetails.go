package validators

import (
	"strings"

	"github.com/AlladinDev/AlShifa/internal/clinicplan/models"
	"github.com/AlladinDev/AlShifa/internal/shared/constants"
)

func ValidatePlanDetails(plan *models.Plan) map[string]string {
	if plan == nil {
		return map[string]string{
			"plan": "plan cannot be nil",
		}
	}

	errs := make(map[string]string)

	if plan.ID.IsZero() {
		errs["id"] = "id is required"
	}

	plan.Name = strings.TrimSpace(plan.Name)
	switch plan.Name {
	case constants.ClinicPlanBasic,
		constants.ClinicSilverPlan,
		constants.ClinicPlanPremium,
		constants.ClinicPlanGold:
		// valid
	default:
		errs["name"] = "invalid plan name"
	}

	if plan.AmountToDeduct <= 0 {
		errs["amountToDeduct"] = "must be greater than 0"
	}

	if plan.CreatedAt.IsZero() {
		errs["createdAt"] = "createdAt is required"
	}

	if plan.UpdatedAt.IsZero() {
		errs["updatedAt"] = "updatedAt is required"
	}

	if !plan.CreatedAt.IsZero() &&
		!plan.UpdatedAt.IsZero() &&
		plan.UpdatedAt.Before(plan.CreatedAt) {
		errs["updatedAt"] = "updatedAt cannot be before createdAt"
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}
