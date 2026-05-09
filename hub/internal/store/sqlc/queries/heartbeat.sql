-- name: InsertHeartbeat :exec
INSERT INTO heartbeats (agent_id, cpu_usage, memory_usage, disk_usage, heartbeat_timestamp)
VALUES (:agent_id, :cpu_usage, :memory_usage, :disk_usage, :heartbeat_timestamp);

-- name: SelectHeartbeatsAfter :many
SELECT * FROM heartbeats
WHERE agent_id = :agent_id AND heartbeat_timestamp > :timestamp;