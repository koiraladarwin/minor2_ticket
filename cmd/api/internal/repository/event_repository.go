package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/models"
)

var ErrEventNotFound = errors.New("event not found")

type EventRepository interface {
	Create(ctx context.Context, event *models.Event) error
	GetByIDAndUserID(ctx context.Context, id uuid.UUID, userId uuid.UUID) (*models.Event, error)
	GetAll(ctx context.Context) ([]models.Event, error)
	GetEventDetails(ctx context.Context, id uuid.UUID, userId uuid.UUID) (*models.EventDetails, error)
	Update(ctx context.Context, event *models.Event) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type eventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) EventRepository {
	return &eventRepository{
		db: db,
	}
}

func (r *eventRepository) Create(ctx context.Context, event *models.Event) error {
	query := `
	INSERT INTO events (
		title,
		description,
		venue,
		banner_url,
		event_start_at,
		event_end_at,
		ticket_sale_start_at,
		ticket_sale_end_at,
		capacity,
		status,
		created_by
	)
	VALUES (
		$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11
	)
	RETURNING
		id,
		created_at,
		updated_at;
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		event.Title,
		event.Description,
		event.Venue,
		event.BannerURL,
		event.EventStartAt,
		event.EventEndAt,
		event.TicketSaleStartAt,
		event.TicketSaleEndAt,
		event.Capacity,
		event.Status,
		event.CreatedBy,
	).Scan(
		&event.ID,
		&event.CreatedAt,
		&event.UpdatedAt,
	)
}

func (r *eventRepository) GetByIDAndUserID(ctx context.Context, id uuid.UUID, userId uuid.UUID) (*models.Event, error) {
	query := `
	SELECT
		id,
		title,
		description,
		venue,
		banner_url,
		event_start_at,
		event_end_at,
		ticket_sale_start_at,
		ticket_sale_end_at,
		capacity,
		status,
		created_by,
		created_at,
		updated_at
	FROM events
	WHERE id = $1
	AND created_by = $2;
	`

	var event models.Event

	err := r.db.QueryRowContext(ctx, query, id, userId).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.Venue,
		&event.BannerURL,
		&event.EventStartAt,
		&event.EventEndAt,
		&event.TicketSaleStartAt,
		&event.TicketSaleEndAt,
		&event.Capacity,
		&event.Status,
		&event.CreatedBy,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrEventNotFound
	}

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *eventRepository) GetEventDetails(
	ctx context.Context,
	id uuid.UUID,
	userId uuid.UUID,
) (*models.EventDetails, error) {

	query := `
	SELECT
		e.id,
		e.title,
		e.description,
		e.venue,
		e.banner_url,
		e.event_start_at,
		e.event_end_at,
		e.ticket_sale_start_at,
		e.ticket_sale_end_at,
		e.capacity,
		e.status,
		e.created_by,

		tt.id,
		tt.name,
		tt.description,
		tt.price,
		tt.quantity,
		tt.remaining,
		tt.created_at,
		tt.updated_at

	FROM events e
	LEFT JOIN ticket_types tt
	ON e.id = tt.event_id

	WHERE e.id = $1
	AND created_by = $2;
	`

	rows, err := r.db.QueryContext(ctx, query, id, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var event models.EventDetails
	var found bool

	for rows.Next() {

		var ticket models.TicketType

		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.Venue,
			&event.BannerURL,
			&event.EventStartAt,
			&event.EventEndAt,
			&event.TicketSaleStartAt,
			&event.TicketSaleEndAt,
			&event.Capacity,
			&event.Status,
			&event.CreatedBy,

			&ticket.ID,
			&ticket.Name,
			&ticket.Description,
			&ticket.Price,
			&ticket.Quantity,
			&ticket.Remaining,
			&ticket.CreatedAt,
			&ticket.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		found = true

		// because LEFT JOIN can return null ticket
		if ticket.ID != uuid.Nil {
			event.TicketTypes = append(
				event.TicketTypes,
				ticket,
			)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if !found {
		return nil, ErrEventNotFound
	}

	return &event, nil
}
func (r *eventRepository) GetAll(ctx context.Context) ([]models.Event, error) {
	query := `
	SELECT
		id,
		title,
		description,
		venue,
		banner_url,
		event_start_at,
		event_end_at,
		ticket_sale_start_at,
		ticket_sale_end_at,
		capacity,
		status,
		created_by,
		created_at,
		updated_at
	FROM events
	ORDER BY event_start_at ASC;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event

	for rows.Next() {
		var event models.Event

		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.Venue,
			&event.BannerURL,
			&event.EventStartAt,
			&event.EventEndAt,
			&event.TicketSaleStartAt,
			&event.TicketSaleEndAt,
			&event.Capacity,
			&event.Status,
			&event.CreatedBy,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *eventRepository) Update(ctx context.Context, event *models.Event) error {
	query := `
	UPDATE events
	SET
		title = $1,
		description = $2,
		venue = $3,
		banner_url = $4,
		event_start_at = $5,
		event_end_at = $6,
		ticket_sale_start_at = $7,
		ticket_sale_end_at = $8,
		capacity = $9,
		status = $10,
		updated_at = NOW()
	WHERE id = $11;
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		event.Title,
		event.Description,
		event.Venue,
		event.BannerURL,
		event.EventStartAt,
		event.EventEndAt,
		event.TicketSaleStartAt,
		event.TicketSaleEndAt,
		event.Capacity,
		event.Status,
		event.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrEventNotFound
	}

	return nil
}

func (r *eventRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM events WHERE id = $1;`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrEventNotFound
	}

	return nil
}
