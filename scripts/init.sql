DROP TABLE IF EXISTS bookings, rooms;

CREATE TABLE IF NOT EXISTS rooms
(
    id          SERIAL PRIMARY KEY,
    description text,
    price       int         NOT NULL,
    created     timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS bookings
(
    id         serial PRIMARY KEY,
    date_start date NOT NULL,
    date_end   date NOT NULL,
    room       int  NOT NULL,

    FOREIGN KEY (room) REFERENCES rooms (id) ON DELETE CASCADE
);