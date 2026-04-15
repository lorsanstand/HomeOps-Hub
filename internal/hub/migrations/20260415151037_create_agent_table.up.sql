CREATE TABLE agents (
    id BIGINT UNIQUE PRIMARY KEY,
    agent_id VARCHAR(32) UNIQUE NOT NULL,
    agent_name VARCHAR(255),
    architecture VARCHAR(10) NOT NULL,
    system VARCHAR(10) NOT NULL,
    hostname VARCHAR(100) NOT NULL,
    version VARCHAR(10) NOT NULL,
    capabilities JSON,
    registered_at timestamp without time zone DEFAULT now() NOT NULL
);

CREATE UNIQUE INDEX idx_agent_id ON agent (id);
CREATE UNIQUE INDEX idx_agent_id_id On agent (agent_id)