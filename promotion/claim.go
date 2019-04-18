package promotion

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx/types"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

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
}

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
