package promotion

import uuid "github.com/satori/go.uuid"

type Issuer struct {
	PromotionID uuid.UUID `db:"promotion_id"`
	Cohort      string
	PublicKey   string `db:"public_key"`
}
