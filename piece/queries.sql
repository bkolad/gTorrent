-- name: save-piece
INSERT INTO pieces (info_hash, index, piece) VALUES($1, $2, $3)

-- name: get-piece
SELECT piece FROM pieces where  info_hash = $1 AND index = $2
