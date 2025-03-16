CREATE TABLE sessions (
    cookie CHAR(44) PRIMARY KEY,
    user_id INT NOT NULL,
    last_accessed TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT fk_session_user FOREIGN KEY (user_id) REFERENCES individuals (id) ON DELETE CASCADE
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_last_accessed ON sessions(last_accessed);
