package repositories

import (
    "database/sql"

    "gastrobar-backend/internal/models"

    "github.com/pkg/errors"
)

type BusinessRepository interface {
    Find() (models.Business, error)
    Update(business models.Business) (models.Business, error)
}

type businessRepository struct {
    db *sql.DB
}

func NewBusinessRepository(db *sql.DB) BusinessRepository {
    return &businessRepository{db: db}
}

func (r *businessRepository) Find() (models.Business, error) {
    var business models.Business
    err := r.db.QueryRow(`
        SELECT id, business_name, address, phone_number, email, corporate_reason, created_at, updated_at
        FROM business
        LIMIT 1`,
    ).Scan(&business.ID, &business.BusinessName, &business.Address, &business.PhoneNumber, &business.Email, &business.CorporateReason, &business.CreatedAt, &business.UpdatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return models.Business{}, errors.Wrap(err, "business not found")
        }
        return models.Business{}, errors.Wrap(err, "failed to query business")
    }
    return business, nil
}

func (r *businessRepository) Update(business models.Business) (models.Business, error) {
    var updatedBusiness models.Business
    err := r.db.QueryRow(`
        UPDATE business
        SET business_name = $1, address = $2, phone_number = $3, email = $4, corporate_reason = $5, updated_at = CURRENT_TIMESTAMP
        WHERE id = $6
        RETURNING id, business_name, address, phone_number, email, corporate_reason, created_at, updated_at`,
        business.BusinessName, business.Address, business.PhoneNumber, business.Email, business.CorporateReason, business.ID,
    ).Scan(&updatedBusiness.ID, &updatedBusiness.BusinessName, &updatedBusiness.Address, &updatedBusiness.PhoneNumber, &updatedBusiness.Email, &updatedBusiness.CorporateReason, &updatedBusiness.CreatedAt, &updatedBusiness.UpdatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return models.Business{}, errors.Wrap(err, "business not found")
        }
        return models.Business{}, errors.Wrap(err, "failed to update business")
    }
    return updatedBusiness, nil
}