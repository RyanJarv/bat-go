package promotion

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	// needed for magic migration
	"github.com/brave-intl/bat-go/wallet"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Datastore abstracts over the underlying datastore
type Datastore interface {
	ActivatePromotion(promotion *Promotion) error
	ClaimForWallet(promotion *Promotion, wallet *wallet.Info, blindedCreds JSONStringArray) (*Claim, error)
	CreateClaim(promotionID uuid.UUID, walletID string, value decimal.Decimal) (*Claim, error)
	CreatePromotion(promotionType string, numGrants int, value decimal.Decimal) (*Promotion, error)
	GetAvailablePromotionsForWallet(wallet *wallet.Info) ([]Promotion, error)
	GetClaimCreds(claimID uuid.UUID) (*ClaimCreds, error)
	SaveClaimCreds(claimCreds *ClaimCreds) error
	GetPromotion(promotionID uuid.UUID) (*Promotion, error)
	SaveIssuer(issuer *Issuer) error
	GetIssuer(promotionID uuid.UUID, cohort string) (*Issuer, error)
	SaveWallet(wallet *wallet.Info) error
	GetWallet(id uuid.UUID) (*wallet.Info, error)
}

// Postgres is a WIP Datastore
type Postgres struct {
	*sqlx.DB
}

func (pg *Postgres) NewMigrate() (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(pg.DB.DB, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///src/migrations",
		"postgres", driver)
	if err != nil {
		return nil, err
	}

	return m, err
}

func (pg *Postgres) Migrate() error {
	m, err := pg.NewMigrate()
	if err != nil {
		return err
	}

	err = m.Migrate(1)
	if err != migrate.ErrNoChange && err != nil {
		return err
	}
	return nil
}

// NewPostgres creates a new Postgres Datastore
func NewPostgres(databaseURL string, performMigration bool) (*Postgres, error) {
	if len(databaseURL) == 0 {
		databaseURL = os.Getenv("DATABASE_URL")
	}

	db, err := sqlx.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	pg := &Postgres{db}

	if performMigration {
		pg.Migrate()
	}

	return pg, nil
}

func (pg *Postgres) CreatePromotion(promotionType string, numGrants int, value decimal.Decimal) (*Promotion, error) {
	statement := `
	insert into promotions (promotion_type, remaining_grants, approximate_value, suggestions_per_grant)
	values ($1, $2, $3, $4)
	returning *`
	promotions := []Promotion{}
	suggestionsPerGrant := value.Div(decimal.NewFromFloat(0.25))
	err := pg.DB.Select(&promotions, statement, promotionType, numGrants, value, suggestionsPerGrant)
	if err != nil {
		return nil, err
	}

	return &promotions[0], nil
}

func (pg *Postgres) GetPromotion(promotionID uuid.UUID) (*Promotion, error) {
	statement := "select * from promotions where id = $1"
	promotions := []Promotion{}
	err := pg.DB.Select(&promotions, statement, promotionID)
	if err != nil {
		return nil, err
	}

	if len(promotions) > 0 {
		return &promotions[0], nil
	} else {
		return nil, nil
	}
}

func (pg *Postgres) ActivatePromotion(promotion *Promotion) error {
	_, err := pg.DB.Exec("update promotions set active = true where id = $1", promotion.ID)
	if err != nil {
		return err
	}

	return nil
}

func (pg *Postgres) SaveIssuer(issuer *Issuer) error {
	statement := `
	insert into issuers (promotion_id, cohort, public_key)
	values ($1, $2, $3)
	returning *`
	_, err := pg.DB.Exec(statement, issuer.PromotionID, issuer.Cohort, issuer.PublicKey)
	if err != nil {
		return err
	}

	return nil
}

func (pg *Postgres) GetIssuer(promotionID uuid.UUID, cohort string) (*Issuer, error) {
	statement := "select * from issuers where promotion_id = $1 and cohort = $2"
	issuers := []Issuer{}
	err := pg.DB.Select(&issuers, statement, promotionID.String(), cohort)
	if err != nil {
		return nil, err
	}

	if len(issuers) > 0 {
		return &issuers[0], nil
	} else {
		return nil, nil
	}
}

func (pg *Postgres) SaveWallet(wallet *wallet.Info) error {
	statement := `
	insert into wallets (id, provider, provider_id, public_key)
	values ($1, $2, $3, $4)
	returning *`
	_, err := pg.DB.Exec(statement, wallet.ID, wallet.Provider, wallet.ProviderID, wallet.PublicKey)
	if err != nil {
		return err
	}

	return nil
}

func (pg *Postgres) GetWallet(ID uuid.UUID) (*wallet.Info, error) {
	statement := "select * from wallets where id = $1"
	wallets := []wallet.Info{}
	err := pg.DB.Select(&wallets, statement, ID)
	if err != nil {
		return nil, err
	}

	if len(wallets) > 0 {
		return &wallets[0], nil
	} else {
		return nil, nil
	}
}

func (pg *Postgres) CreateClaim(promotionID uuid.UUID, walletID string, value decimal.Decimal) (*Claim, error) {
	statement := `
	insert into claims (promotion_id, wallet_id, approximate_value)
	values ($1, $2, $3)
	returning *`
	claims := []Claim{}
	err := pg.DB.Select(&claims, statement, promotionID, walletID, value)
	if err != nil {
		return nil, err
	}

	return &claims[0], nil
}

func (pg *Postgres) ClaimForWallet(promotion *Promotion, wallet *wallet.Info, blindedCreds JSONStringArray) (*Claim, error) {
	blindedCredsJSON, err := json.Marshal(blindedCreds)
	if err != nil {
		return nil, err
	}

	tx, err := pg.DB.Beginx()
	if err != nil {
		return nil, err
	}

	// This will error if remaining_grants is insufficent due to constraint or the promotion is inactive
	res, err := tx.Exec(`update promotions set remaining_grants = remaining_grants - 1 where id = $1 and active`, promotion.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	promotionCount, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return nil, err
	} else if promotionCount != 1 {
		tx.Rollback()
		return nil, errors.New("no matching active promotion")
	}

	claims := []Claim{}

	if promotion.Type == "ads" || promotion.Version < 5 {
		statement := `
    update claims
		set redeemed = true
		where promotion_id = $1 and wallet_id = $2
		returning *`
		err = tx.Select(&claims, statement, promotion.ID, wallet.ID)
	} else {
		statement := `
		insert into claims (promotion_id, wallet_id, approximate_value, redeemed)
		values ($1, $2, $3, true)
		returning *`
		err = tx.Select(&claims, statement, promotion.ID, wallet.ID, promotion.ApproximateValue)
	}
	if err != nil {
		tx.Rollback()
		return nil, err
	} else if len(claims) != 1 {
		tx.Rollback()
		return nil, fmt.Errorf("Incorrect number of claims updated / inserted: %d", len(claims))
	}
	claim := claims[0]

	// This will error if user has already claimed due to uniqueness constraint
	_, err = tx.Exec(`insert into claim_creds (claim_id, blinded_creds) values ($1, $2)`, claim.ID, blindedCredsJSON)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &claim, nil
}

func (pg *Postgres) GetAvailablePromotionsForWallet(wallet *wallet.Info) ([]Promotion, error) {
	statement := `
	select 
		promotions.*, 
		promotions.active and promotions.remaining_grants > 0 and (
			( promotions.promotion_type = 'ugp' and claims.id is null ) or 
			( ( promotion_type = 'ads' or promotions.version < 5 ) and claims.id is not null and not claims.redeemed )
		) as available
	from promotions left join claims on promotions.id = claims.promotion_id and claims.wallet_id = $1;`
	promotions := []Promotion{}

	err := pg.DB.Select(&promotions, statement, wallet.ID)
	if err != nil {
		return promotions, err
	}

	return promotions, nil
}

func (pg *Postgres) GetClaimCreds(claimID uuid.UUID) (*ClaimCreds, error) {
	claimCreds := []ClaimCreds{}
	err := pg.DB.Select(&claimCreds, "select * from claim_creds where claim_id = $1", claimID)
	if err != nil {
		return nil, err
	}

	if len(claimCreds) > 0 {
		return &claimCreds[0], nil
	} else {
		return nil, nil
	}
}

func (pg *Postgres) SaveClaimCreds(creds *ClaimCreds) error {
	_, err := pg.DB.Exec(`update claim_creds set signed_creds = $1, batch_proof = $2, public_key = $3 where claim_id = $4`, creds.SignedCreds, creds.BatchProof, creds.PublicKey, creds.ID)
	return err
}
