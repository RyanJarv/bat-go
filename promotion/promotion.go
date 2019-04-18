package promotion

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

type Promotion struct {
	ID                  uuid.UUID `db:"id"`
	CreatedAt           time.Time `db:"created_at"`
	ExpiresAt           time.Time `db:"expires_at"`
	ClaimableUntil      time.Time
	Version             int             `db:"version"`
	SuggestionsPerGrant int             `db:"suggestions_per_grant"`
	ApproximateValue    decimal.Decimal `db:"approximate_value"`
	Type                string          `db:"promotion_type"`
	RemainingGrants     int             `db:"remaining_grants"`
	Active              bool            `db:"active"`
	Available           bool            `db:"available"`
	//PublicKeys          []string
}
