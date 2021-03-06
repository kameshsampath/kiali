package route_rules

import (
	"github.com/kiali/kiali/kubernetes"
	"github.com/kiali/kiali/services/models"
)

type NoDestinationChecker struct {
	Namespace    string
	ServiceNames []string
	RouteRule    kubernetes.IstioObject
}

func (routeRule NoDestinationChecker) Check() ([]*models.IstioCheck, bool) {
	valid := false
	validations := make([]*models.IstioCheck, 0)

	for _, serviceName := range routeRule.ServiceNames {
		if valid = kubernetes.FilterByDestination(routeRule.RouteRule.GetSpec(), routeRule.Namespace, serviceName, ""); valid {
			break
		}
	}

	if !valid {
		validation := models.BuildCheck("Destination doesn't have a valid service", "error", "spec/destination")
		validations = append(validations, &validation)
	}

	return validations, valid
}
