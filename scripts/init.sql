CREATE TABLE IF NOT EXISTS rooms
(
    id          SERIAL PRIMARY KEY,
    description text,
    price       int         NOT NULL,
    created     timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX cover_index ON rooms (id, description, price, created);
CREATE INDEX price_order_by_asc_rooms ON rooms (price ASC);
CREATE INDEX price_order_by_desc_rooms ON rooms (price DESC);
CREATE INDEX created_order_by_asc_rooms ON rooms (created ASC);
CREATE INDEX created_order_by_desc_rooms ON rooms (created DESC);

CREATE TABLE IF NOT EXISTS bookings
(
    id         serial PRIMARY KEY,
    date_start date NOT NULL,
    date_end   date NOT NULL,
    room       int  NOT NULL,

    FOREIGN KEY (room) REFERENCES rooms (id) ON DELETE CASCADE
);
CREATE INDEX cover_bookings ON bookings (id, date_start, date_end, room);
CREATE INDEX room_bookings ON bookings (room);
CREATE INDEX date_start_order_by_bookings ON bookings (date_start ASC);
