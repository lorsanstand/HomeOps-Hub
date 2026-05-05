CREATE TABLE agents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    agent_id VARCHAR(32) UNIQUE NOT NULL,
    agent_name VARCHAR(255),
    architecture VARCHAR(10) NOT NULL,
    system VARCHAR(10) NOT NULL,
    hostname VARCHAR(100) NOT NULL,
    version VARCHAR(10) NOT NULL,
    capabilities TEXT,
    registered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX idx_agent_id ON agents (id);
CREATE UNIQUE INDEX idx_agent_id_id ON agents (agent_id);