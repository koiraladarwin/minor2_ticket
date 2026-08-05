package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/kafka"
	"github.com/koiraladarwin/minor2_ticket/cmd/api/internal/models"
)

var (
	TicketActive         = "ACTIVE"
	TicketUsed           = "USED"
	TicketCancelled      = "CANCELLED"
	ErrTicketNotFound    = errors.New("ticket not found")
	ErrTicketAlreadyUsed = errors.New("ticket already used")
	ErrAlreadyPurchased  = errors.New("user already purchased ticket for this event")
	ErrTicketSoldOut     = errors.New("ticket sold out")
	ErrEventUnavailable  = errors.New("event unavailable")
	ErrTicketSalesClosed = errors.New("ticket sales are closed")
)

type ticketRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewTicketRepository(db *sql.DB, redis *redis.Client) TicketRepository {
	return &ticketRepository{
		db:    db,
		redis: redis,
	}
}

func (r *ticketRepository) Purchase(
	ctx context.Context,
	ticketTypeID uuid.UUID,
	userID uuid.UUID,
) (*models.TicketDetail, error) {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	var eventID uuid.UUID
	var remaining int
	var eventStatus string
	var saleStart time.Time
	var saleEnd time.Time

	// Lock ticket type row to prevent overselling
	err = tx.QueryRowContext(
		ctx,
		`
		SELECT
			tt.event_id,
			tt.remaining,
			e.status,
			e.ticket_sale_start_at,
			e.ticket_sale_end_at
		FROM ticket_types tt
		JOIN events e
			ON e.id = tt.event_id
		WHERE tt.id = $1
		FOR UPDATE;
		`,
		ticketTypeID,
	).Scan(
		&eventID,
		&remaining,
		&eventStatus,
		&saleStart,
		&saleEnd,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrEventUnavailable
	}

	if err != nil {
		return nil, err
	}

	now := time.Now()

	if now.Before(saleStart) || now.After(saleEnd) {
		return nil, ErrTicketSalesClosed
	}

	if remaining <= 0 {
		return nil, ErrTicketSoldOut
	}

	// Check duplicate purchase
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
		eventID,
		userID,
	).Scan(&exists)

	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrAlreadyPurchased
	}

	// Decrease remaining ticket count
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

	ticketNumber := fmt.Sprintf(
		"TKT-%d",
		time.Now().UnixNano(),
	)

	qrCode := uuid.New().String()

	var ticketID uuid.UUID

	// Create ticket
	err = tx.QueryRowContext(
		ctx,
		`
		INSERT INTO tickets(
			event_id,
			ticket_type_id,
			user_id,
			ticket_number,
			qr_code,
			status
		)
		VALUES(
			$1,$2,$3,$4,$5,$6
		)
		RETURNING id;
		`,
		eventID,
		ticketTypeID,
		userID,
		ticketNumber,
		qrCode,
		models.TicketActive,
	).Scan(
		&ticketID,
	)

	if err != nil {
		return nil, err
	}

	// Commit purchase transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Fetch complete ticket details
	var ticket models.TicketDetail

	err = r.db.QueryRowContext(
		ctx,
		`
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
			e.created_by,

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

		WHERE t.id = $1;
		`,
		ticketID,
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
		&ticket.Event.CreatedBy,

		// ticket type
		&ticket.TicketType.ID,
		&ticket.TicketType.Name,
		&ticket.TicketType.Description,
		&ticket.TicketType.Price,
	)

	if err != nil {
		return nil, err
	}

	return &ticket, nil
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
		e.created_by,

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
			&ticket.Event.CreatedBy,

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
		e.created_by,

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
		&ticket.Event.CreatedBy,

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

func (r *ticketRepository) ScanTicket(
	ctx context.Context,
	ticketID string,
	scannedBy string,
) error {

	// Get ticket from Redis
	data, err := r.redis.Get(
		ctx,
		"ticket:"+ticketID,
	).Result()

	if err == redis.Nil {
		return ErrTicketNotFound
	}

	if err != nil {
		return err
	}

	var ticket kafka.TicketCache

	err = json.Unmarshal(
		[]byte(data),
		&ticket,
	)

	if err != nil {
		return err
	}

	// Check status
	if ticket.Status == TicketUsed {
		return ErrTicketAlreadyUsed
	}

	// Update cache status
	ticket.Status = TicketUsed

	updated, err := json.Marshal(ticket)
	if err != nil {
		return err
	}

	err = r.redis.Set(
		ctx,
		"ticket:"+ticketID,
		updated,
		0,
	).Err()

	if err != nil {
		return err
	}
	return nil
}
