package promotion

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

type JSONStringArray []string

func (arr *JSONStringArray) Scan(src interface{}) error {
	var jt types.JSONText

	if err := jt.Scan(src); err != nil {
		return err
	}

	if err := jt.Unmarshal(arr); err != nil {
		return err
	}

	return nil
}

func (arr *JSONStringArray) Value() (driver.Value, error) {
	var jt types.JSONText

	data, err := json.Marshal((*[]string)(arr))
	if err != nil {
		return nil, err
	}

	if err := jt.UnmarshalJSON(data); err != nil {
		return nil, err
	}

	return jt.Value()
}

func (arr *JSONStringArray) MarshalJSON() ([]byte, error) {
	return json.Marshal((*[]string)(arr))
}

func (arr *JSONStringArray) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, (*[]string)(arr)); err != nil {
		return err
	}

	return nil
}

type Claim struct {
	ID               uuid.UUID       `db:"id"`
	CreatedAt        time.Time       `db:"created_at"`
	PromotionID      uuid.UUID       `db:"promotion_id"`
	WalletID         uuid.UUID       `db:"wallet_id"`
	ApproximateValue decimal.Decimal `db:"approximate_value"`
	Redeemed         bool            `db:"redeemed"`
}

type ClaimCreds struct {
	ID           uuid.UUID        `db:"claim_id"`
	BlindedCreds JSONStringArray  `db:"blinded_creds"`
	SignedCreds  *JSONStringArray `db:"signed_creds"`
	BatchProof   *string          `db:"batch_proof"`
	PublicKey    *string          `db:"public_key"`
}

func (service *Service) ClaimPromotionForWallet(ctx context.Context, promotionID uuid.UUID, walletID uuid.UUID, blindedCreds []string) (*uuid.UUID, error) {
	promotion, err := service.datastore.GetPromotion(promotionID)
	if err != nil {
		return nil, err
	}
	if promotion == nil {
		return nil, errors.New("promotion did not exist")
	}

	wallet, err := service.datastore.GetWallet(walletID)
	if err != nil || wallet == nil {
		return nil, errors.Wrap(err, "Error getting wallet")
	}

	// TODO lookup and return existing claim if exists?

	// TODO lookup reputation server

	cohort := "control"
	issuer, err := service.GetOrCreateIssuer(ctx, promotionID, cohort)
	if err != nil {
		return nil, err
	}

	if len(blindedCreds) != promotion.SuggestionsPerGrant {
		return nil, errors.New("wrong number of blinded tokens included")
	}

	claim, err := service.datastore.ClaimForWallet(promotion, wallet, JSONStringArray(blindedCreds))
	if err != nil {
		return nil, err
	}

	go service.SignClaimCreds(ctx, claim.ID, *issuer, blindedCreds)

	return &claim.ID, nil
}

func (service *Service) SignClaimCreds(ctx context.Context, claimID uuid.UUID, issuer Issuer, blindedCreds []string) {
	resp, err := service.cbClient.SignCredentials(ctx, issuer.Name(), blindedCreds)
	if err != nil {
		// FIXME
		fmt.Println(err)
	}

	signedTokens := JSONStringArray(resp.SignedTokens)

	creds := &ClaimCreds{
		ID:           claimID,
		BlindedCreds: blindedCreds,
		SignedCreds:  &signedTokens,
		BatchProof:   &resp.BatchProof,
		PublicKey:    &issuer.PublicKey,
	}

	err = service.datastore.SaveClaimCreds(creds)
	if err != nil {
		// FIXME
		fmt.Println(err)
	}
}
