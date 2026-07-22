-- name: GetChirps :exec
SELECT * FROM chirps
ORDER BY created_at ASC;