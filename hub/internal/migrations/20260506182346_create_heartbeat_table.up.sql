CREATE TABLE heartbeats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    agent_id VARCHAR(32) UNIQUE NOT NULL,
    cpu_usage FLOAT,
    memory_usage FLOAT,
    disk_usage FLOAT,
    heartbeat_timestamp TIMESTAMP
    FOREIGN KEY (agent_id) REFERENCES agents(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_heartbeat_agent_id ON heartbeats (agent_id);