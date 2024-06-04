-- name: CreateTicket :one
INSERT INTO tickets (user_id, title, content, priority)
VALUES ($1, $2, $3, $4)
RETURNING ticket_id;

-- name: UpdateTicketStatus :exec
UPDATE tickets
SET status = $2
WHERE ticket_id = $1;

-- name: GetTicketById :one
SELECT ticket_id, user_id, title, content, priority, status
FROM tickets
WHERE ticket_id = $1;
