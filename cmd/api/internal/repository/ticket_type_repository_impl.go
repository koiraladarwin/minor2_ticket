package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/kafka"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/models"
)

var ErrTicketTypeNotFound = errors.New("ticket type not found")

type ticketTypeRepository struct {
	db    *sql.DB
	kafka *kafka.Producer
}

func NewTicketTypeRepository(db *sql.DB, kakfa *kafka.Producer) TicketTypeRepository {
	return &ticketTypeRepository{
		db:    db,
		kafka: kakfa,
	}
}

func (r *ticketTypeRepository) Create(ctx context.Context, ticketType *models.TicketType) error {
	query := `
	INSERT INTO ticket_types (
		event_id,
		name,
		description,
		price,
		quantity,
		remaining
	)
	VALUES (
		$1,$2,$3,$4,$5,$6
	)
	RETURNING
		id,
		created_at,
		updated_at;
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		ticketType.EventID,
		ticketType.Name,
		ticketType.Description,
		ticketType.Price,
		ticketType.Quantity,
		ticketType.Remaining,
	).Scan(
		&ticketType.ID,
		&ticketType.CreatedAt,
		&ticketType.UpdatedAt,
	)
}

func (r *ticketTypeRepository) GetByIDAndUserID(
	ctx context.Context,
	id uuid.UUID,
	userID uuid.UUID,
) (*models.TicketType, error) {

	query := `
	SELECT
		tt.id,
		tt.event_id,
		tt.name,
		tt.description,
		tt.price,
		tt.quantity,
		tt.remaining,
		tt.created_at,
		tt.updated_at
	FROM ticket_types tt
	INNER JOIN events e
		ON e.id = tt.event_id
	WHERE
		tt.id = $1
	AND
		e.created_by = $2;
	`

	var ticketType models.TicketType

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
		userID,
	).Scan(
		&ticketType.ID,
		&ticketType.EventID,
		&ticketType.Name,
		&ticketType.Description,
		&ticketType.Price,
		&ticketType.Quantity,
		&ticketType.Remaining,
		&ticketType.CreatedAt,
		&ticketType.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTicketTypeNotFound
	}

	if err != nil {
		return nil, err
	}

	return &ticketType, nil
}

func (r *ticketTypeRepository) GetByEventID(
	ctx context.Context,
	eventID uuid.UUID,
) ([]models.TicketType, error) {

	query := `
	SELECT
		id,
		event_id,
		name,
		description,
		price,
		quantity,
		remaining,
		created_at,
		updated_at
	FROM ticket_types
	WHERE event_id = $1
	ORDER BY price ASC;
	`

	rows, err := r.db.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ticketTypes []models.TicketType

	for rows.Next() {

		var ticketType models.TicketType

		err := rows.Scan(
			&ticketType.ID,
			&ticketType.EventID,
			&ticketType.Name,
			&ticketType.Description,
			&ticketType.Price,
			&ticketType.Quantity,
			&ticketType.Remaining,
			&ticketType.CreatedAt,
			&ticketType.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		ticketTypes = append(ticketTypes, ticketType)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ticketTypes, nil
}

func (r *ticketTypeRepository) Update(
	ctx context.Context,
	ticketType *models.TicketType,
) error {

	query := `
	UPDATE ticket_types
	SET
		name = $1,
		description = $2,
		price = $3,
		quantity = $4,
		remaining = $5,
		updated_at = NOW()
	WHERE id = $6;
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		ticketType.Name,
		ticketType.Description,
		ticketType.Price,
		ticketType.Quantity,
		ticketType.Remaining,
		ticketType.ID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrTicketTypeNotFound
	}

	return nil
}

func (r *ticketTypeRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {

	query := `
	DELETE FROM ticket_types
	WHERE id = $1;
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		id,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrTicketTypeNotFound
	}

	return nil
}
