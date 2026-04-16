-- name: CreateAgent :exec
INSERT INTO agents (agent_id, agent_name, architecture, system, hostname, version, capabilities)
VALUES  ($1, $2, $3, $4, $5, $6, $7);