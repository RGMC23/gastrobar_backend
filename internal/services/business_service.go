package services

import (
    "gastrobar-backend/internal/models"
    "gastrobar-backend/internal/repositories"

    "github.com/pkg/errors"
)

// BusinessService define las operaciones relacionadas con el negocio
type BusinessService interface {
    GetBusiness() (models.Business, error)
    UpdateBusiness(business models.Business) (models.Business, error)
}

type businessService struct {
    businessRepo repositories.BusinessRepository
}

// NewBusinessService crea una nueva instancia del servicio de negocio
func NewBusinessService(businessRepo repositories.BusinessRepository) BusinessService {
    return &businessService{
        businessRepo: businessRepo,
    }
}

// GetBusiness obtiene los datos del único negocio (Gastrobar)
func (s *businessService) GetBusiness() (models.Business, error) {
    business, err := s.businessRepo.Find()
    if err != nil {
        return models.Business{}, errors.Wrap(err, "failed to get business")
    }
    return business, nil
}

// UpdateBusiness actualiza los datos del negocio
func (s *businessService) UpdateBusiness(business models.Business) (models.Business, error) {
    updatedBusiness, err := s.businessRepo.Update(business)
    if err != nil {
        return models.Business{}, errors.Wrap(err, "failed to update business")
    }
    return updatedBusiness, nil
}