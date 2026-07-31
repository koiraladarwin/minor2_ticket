CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    title VARCHAR(255) NOT NULL,
    description TEXT,
    venue VARCHAR(255) NOT NULL,
    banner_url TEXT,

    event_start_at TIMESTAMPTZ NOT NULL,
    event_end_at TIMESTAMPTZ NOT NULL,

    ticket_sale_start_at TIMESTAMPTZ NOT NULL,
    ticket_sale_end_at TIMESTAMPTZ NOT NULL,

    capacity INTEGER NOT NULL CHECK (capacity > 0),

    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT'
        CHECK (status IN ('DRAFT', 'PUBLISHED', 'CANCELLED', 'COMPLETED')),

    created_by UUID NOT NULL, -- User ID from Auth Service

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (event_start_at < event_end_at),
    CHECK (ticket_sale_start_at < ticket_sale_end_at),
    CHECK (ticket_sale_end_at <= event_start_at)
);


CREATE TABLE ticket_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    event_id UUID NOT NULL
        REFERENCES events(id) ON DELETE CASCADE,

    name VARCHAR(100) NOT NULL,
    description TEXT,

    price NUMERIC(10,2) NOT NULL CHECK (price >= 0),

    quantity INTEGER NOT NULL CHECK (quantity > 0),
    remaining INTEGER NOT NULL CHECK (remaining >= 0),
    CHECK (remaining <= quantity),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (event_id, name)
);

CREATE TABLE tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    event_id UUID NOT NULL
        REFERENCES events(id) ON DELETE CASCADE,

    ticket_type_id UUID NOT NULL
        REFERENCES ticket_types(id),

    user_id UUID NOT NULL, -- User ID from Auth Service

    ticket_number VARCHAR(50) NOT NULL UNIQUE,
    qr_code TEXT NOT NULL UNIQUE,

    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE'
        CHECK (status IN ('ACTIVE', 'USED', 'CANCELLED')),

    purchased_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    checked_in_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- One ticket per user per event
    UNIQUE (event_id, user_id)
);

CREATE INDEX idx_events_status
    ON events(status);

CREATE INDEX idx_events_start
    ON events(event_start_at);

CREATE INDEX idx_ticket_types_event
    ON ticket_types(event_id);

CREATE INDEX idx_tickets_user
    ON tickets(user_id);

CREATE INDEX idx_tickets_event
    ON tickets(event_id);

CREATE INDEX idx_tickets_type
    ON tickets(ticket_type_id);

CREATE INDEX idx_tickets_status
    ON tickets(status);