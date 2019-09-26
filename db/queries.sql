-- name: create-torrents-table
CREATE TABLE IF NOT EXISTS torrents (
	id serial PRIMARY KEY,
  info_hash BYTEA NOT NULL UNIQUE,
	name text NOT NULL,
  status text NOT NULL,
	created_at timestamp with time zone DEFAULT current_timestamp
)

-- name: create-pieces-table
CREATE TABLE IF NOT EXISTS pieces (
  info_hash BYTEA,
	index integer,
  piece BYTEA,
	created_at timestamp with time zone DEFAULT current_timestamp,
  FOREIGN KEY (info_hash) REFERENCES torrents (info_hash)
)
