package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/models"
)

var (
	ErrTicketNotFound    = errors.New("ticket not found")
	ErrAlreadyPurchased  = errors.New("user already purchased ticket for this event")
	ErrTicketSoldOut     = errors.New("ticket sold out")
	ErrEventUnavailable  = errors.New("event unavailable")
	ErrTicketSalesClosed = errors.New("ticket sales are closed")
)

type ticketRepository struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) TicketRepository {
	return &ticketRepository{
		db: db,
	}
}

func (r *ticketRepository) Purchase(
	ctx context.Context,
	ticketTypeID uuid.UUID,
	userID uuid.UUID,
) (*models.Ticket, error) {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	type purchaseInfo struct {
		EventID   uuid.UUID
		Remaining int
		Status    string
		SaleStart time.Time
		SaleEnd   time.Time
	}

	var info purchaseInfo

	query := `
	SELECT
		tt.event_id,
		tt.remaining,
		e.status,
		e.ticket_sale_start_at,
		e.ticket_sale_end_at
	FROM ticket_types tt
	INNER JOIN events e
		ON e.id = tt.event_id
	WHERE tt.id = $1
	FOR UPDATE;
	`

	err = tx.QueryRowContext(
		ctx,
		query,
		ticketTypeID,
	).Scan(
		&info.EventID,
		&info.Remaining,
		&info.Status,
		&info.SaleStart,
		&info.SaleEnd,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrEventUnavailable
	}

	if err != nil {
		return nil, err
	}

	now := time.Now()

	if now.Before(info.SaleStart) || now.After(info.SaleEnd) {
		return nil, ErrTicketSalesClosed
	}

	if info.Remaining <= 0 {
		return nil, ErrTicketSoldOut
	}

	var exists bool

	err = tx.QueryRowContext(
		ctx,
		`
		SELECT EXISTS(
			SELECT 1
			FROM tickets
			WHERE event_id = $1
			AND user_id = $2
		);
		`,
		info.EventID,
		userID,
	).Scan(&exists)

	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrAlreadyPurchased
	}

	_, err = tx.ExecContext(
		ctx,
		`
		UPDATE ticket_types
		SET
			remaining = remaining - 1,
			updated_at = NOW()
		WHERE id = $1;
		`,
		ticketTypeID,
	)

	if err != nil {
		return nil, err
	}
	ticket := &models.Ticket{
		EventID:      info.EventID,
		TicketTypeID: ticketTypeID,
		UserID:       userID,

		// Simple unique values for now.
		// These can be replaced later with a nicer format.
		TicketNumber: fmt.Sprintf("TKT-%d", time.Now().UnixNano()),
		QRCode:       uuid.New().String(),

		Status: models.TicketActive,
	}

	insertQuery := `
	INSERT INTO tickets (
		event_id,
		ticket_type_id,
		user_id,
		ticket_number,
		qr_code,
		status
	)
	VALUES (
		$1,$2,$3,$4,$5,$6
	)
	RETURNING
		id,
		purchased_at,
		created_at,
		updated_at;
	`

	err = tx.QueryRowContext(
		ctx,
		insertQuery,
		ticket.EventID,
		ticket.TicketTypeID,
		ticket.UserID,
		ticket.TicketNumber,
		ticket.QRCode,
		ticket.Status,
	).Scan(
		&ticket.ID,
		&ticket.PurchasedAt,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return ticket, nil
}

func (r *ticketRepository) GetByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]models.TicketDetail, error) {

	query := `
	SELECT
		-- ticket
		t.id,
		t.ticket_number,
		t.qr_code,
		t.status,
		t.purchased_at,
		t.checked_in_at,
		t.user_id,

		-- event
		e.id,
		e.title,
		e.description,
		e.venue,
		e.banner_url,
		e.event_start_at,
		e.event_end_at,
		e.status,

		-- ticket type
		tt.id,
		tt.name,
		tt.description,
		tt.price

	FROM tickets t

	JOIN events e
		ON e.id = t.event_id

	JOIN ticket_types tt
		ON tt.id = t.ticket_type_id

	WHERE t.user_id = $1

	ORDER BY t.purchased_at DESC;
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tickets := make([]models.TicketDetail, 0)

	for rows.Next() {

		var ticket models.TicketDetail

		err := rows.Scan(
			// ticket
			&ticket.ID,
			&ticket.TicketNumber,
			&ticket.QRCode,
			&ticket.Status,
			&ticket.PurchasedAt,
			&ticket.CheckedInAt,
			&ticket.UserID,

			// event
			&ticket.Event.ID,
			&ticket.Event.Title,
			&ticket.Event.Description,
			&ticket.Event.Venue,
			&ticket.Event.BannerURL,
			&ticket.Event.StartAt,
			&ticket.Event.EndAt,
			&ticket.Event.Status,

			// ticket type
			&ticket.TicketType.ID,
			&ticket.TicketType.Name,
			&ticket.TicketType.Description,
			&ticket.TicketType.Price,
		)

		if err != nil {
			return nil, err
		}

		tickets = append(tickets, ticket)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tickets, nil
}
func (r *ticketRepository) GetByIDAndUserID(
	ctx context.Context,
	id uuid.UUID,
	userID uuid.UUID,
) (*models.TicketDetail, error) {

	query := `
	SELECT
		-- ticket
		t.id,
		t.ticket_number,
		t.qr_code,
		t.status,
		t.purchased_at,
		t.checked_in_at,
		t.user_id,

		-- event
		e.id,
		e.title,
		e.description,
		e.venue,
		e.banner_url,
		e.event_start_at,
		e.event_end_at,
		e.status,

		-- ticket type
		tt.id,
		tt.name,
		tt.description,
		tt.price

	FROM tickets t

	JOIN events e
		ON e.id = t.event_id

	JOIN ticket_types tt
		ON tt.id = t.ticket_type_id

	WHERE t.id = $1
	AND t.user_id = $2;
	`

	var ticket models.TicketDetail

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
		userID,
	).Scan(
		// ticket
		&ticket.ID,
		&ticket.TicketNumber,
		&ticket.QRCode,
		&ticket.Status,
		&ticket.PurchasedAt,
		&ticket.CheckedInAt,
		&ticket.UserID,

		// event
		&ticket.Event.ID,
		&ticket.Event.Title,
		&ticket.Event.Description,
		&ticket.Event.Venue,
		&ticket.Event.BannerURL,
		&ticket.Event.StartAt,
		&ticket.Event.EndAt,
		&ticket.Event.Status,

		// ticket type
		&ticket.TicketType.ID,
		&ticket.TicketType.Name,
		&ticket.TicketType.Description,
		&ticket.TicketType.Price,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTicketNotFound
	}

	if err != nil {
		return nil, err
	}

	return &ticket, nil
}
