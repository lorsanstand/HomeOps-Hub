-- name: CreateAgent :exec
INSERT INTO agents (agent_id, agent_name, architecture, system, hostname, version, capabilities)
VALUES  ($1, $2, $3, $4, $5, $6, $7);

-- name: GetAgentByID :one
SELECT * from agents
WHERE id=$1;

-- name: GetAgentByAgentID :one
SELECT * from agents
WHERE agent_id=$1;

-- name: UpdateAgentByID :exec
UPDATE agents
SET agent_id=$1, agent_name=$2, architecture=$3, system=$4, hostname=$5, version=$6, capabilities=$7
WHERE id=$8;