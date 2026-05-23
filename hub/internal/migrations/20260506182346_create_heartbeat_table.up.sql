CREATE TABLE heartbeats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    agent_id VARCHAR(32) NOT NULL,
    cpu_usage FLOAT NOT NULL ,
    memory_usage FLOAT NOT NULL ,
    disk_usage FLOAT NOT NULL ,
    heartbeat_timestamp TIMESTAMP NOT NULL,
    FOREIGN KEY (agent_id) REFERENCES agents(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_heartbeat_agent_id ON heartbeats (agent_id);