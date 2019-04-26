package promotion

import (
	"context"
	"time"

	"github.com/brave-intl/bat-go/wallet"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

type Promotion struct {
	ID                  uuid.UUID       `json:"id" db:"id"`
	CreatedAt           time.Time       `json:"createdAt" db:"created_at"`
	ExpiresAt           time.Time       `json:"expiresAt" db:"expires_at"`
	Version             int             `json:"version" db:"version"`
	SuggestionsPerGrant int             `json:"suggestionsPerGrant" db:"suggestions_per_grant"`
	ApproximateValue    decimal.Decimal `json:"approximateValue" db:"approximate_value"`
	Type                string          `json:"type" db:"promotion_type"`
	RemainingGrants     int             `json:"-" db:"remaining_grants"`
	Active              bool            `json:"-" db:"active"`
	Available           bool            `json:"available" db:"available"`
	//ClaimableUntil      time.Time
	//PublicKeys          []string
}

func (service *Service) GetOrCreateWallet(ctx context.Context, walletID uuid.UUID) (*wallet.Info, error) {
	wallet, err := service.datastore.GetWallet(walletID)
	if err != nil {
		return nil, errors.Wrap(err, "Error looking up wallet")
	}

	if wallet == nil {
		wallet, err = service.ledgerClient.GetWallet(ctx, walletID)
		if err != nil {
			return nil, errors.Wrap(err, "Error looking up wallet")
		}
		if wallet != nil {
			err = service.datastore.SaveWallet(wallet)
			if err != nil {
				return nil, errors.Wrap(err, "Error saving wallet")
			}
		}
	}
	return wallet, nil
}

func (service *Service) GetAvailablePromotions(ctx context.Context, walletID uuid.UUID) ([]Promotion, error) {
	wallet, err := service.GetOrCreateWallet(ctx, walletID)
	if err != nil {
		return []Promotion{}, err
	}
	return service.datastore.GetAvailablePromotionsForWallet(wallet)
}
