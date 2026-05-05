-- name: CreateAgent :exec
INSERT INTO agents (agent_id, agent_name, architecture, system, hostname, version, capabilities)
VALUES  (:agent_id, :agent_name, :architecture, :system, :hostname, :version, :capabilities);

-- name: GetAgentByID :one
SELECT * from agents
WHERE id = :id;

-- name: GetAgentByAgentID :one
SELECT * from agents
WHERE agent_id = :agent_id;

-- name: UpdateAgentByID :exec
UPDATE agents
SET agent_id = :agent_id,
    agent_name = :agent_name,
    architecture = :architecture,
    system = :system,
    hostname = :hostname,
    version = :version,
    capabilities = :capabilities
WHERE id = :id;