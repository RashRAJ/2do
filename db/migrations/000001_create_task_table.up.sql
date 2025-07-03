CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    completed BOOLEAN DEFAULT FALSE,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    priority VARCHAR(50),
    status VARCHAR(50) DEFAULT 'pending'
);

-- Add index for faster lookups
CREATE INDEX IF NOT EXISTS idx_tasks_title ON tasks(title);

-- Note: The 'status' field is defined in the Task model but not currently used in repository queries.
-- It's included here for future use.
