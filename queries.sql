-- name: create-torrent
INSERT INTO torrents (name, status, info_hash)

SELECT $1, $2, $3
WHERE
    NOT EXISTS (
        SELECT name FROM torrents WHERE name = $1
    );
