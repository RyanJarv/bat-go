package promotion

import (
	"context"

	uuid "github.com/satori/go.uuid"
)

type Issuer struct {
	PromotionID uuid.UUID `db:"promotion_id"`
	Cohort      string
	PublicKey   string `db:"public_key"`
}

func (service *Service) CreateIssuer(ctx context.Context, promotionID uuid.UUID, cohort string) (*Issuer, error) {
	issuer := &Issuer{PromotionID: promotionID, Cohort: cohort, PublicKey: ""}

	err := service.cbClient.CreateIssuer(ctx, issuer.Name(), 100)
	if err != nil {
		return nil, err
	}

	resp, err := service.cbClient.GetIssuer(ctx, issuer.Name())
	if err != nil {
		return nil, err
	}

	issuer.PublicKey = resp.PublicKey

	service.datastore.SaveIssuer(issuer)
	if err != nil {
		return nil, err
	}

	return issuer, nil
}

func (issuer *Issuer) Name() string {
	return issuer.PromotionID.String() + ":" + issuer.Cohort
}

func (service *Service) GetOrCreateIssuer(ctx context.Context, promotionID uuid.UUID, cohort string) (*Issuer, error) {
	issuer, err := service.datastore.GetIssuer(promotionID, cohort)
	if err != nil {
		return nil, err
	}

	if issuer == nil {
		issuer, err = service.CreateIssuer(ctx, promotionID, cohort)
	}

	return issuer, err
}
