ALTER TABLE tasks ADD COLUMN repeat_type VARCHAR(50);
ALTER TABLE tasks ADD COLUMN repeat_config JSONB;
ALTER TABLE tasks ADD COLUMN next_occurrence TIMESTAMPTZ;
ALTER TABLE tasks ADD COLUMN parent_id BIGINT;
ALTER TABLE tasks ADD COLUMN repeat_until TIMESTAMPTZ;
ALTER TABLE tasks ADD COLUMN repeat_time VARCHAR(5);

CREATE INDEX idx_tasks_next_occurrence ON tasks(next_occurrence) WHERE next_occurrence IS NOT NULL;

ALTER TABLE tasks ADD CONSTRAINT fk_tasks_parent FOREIGN KEY (parent_id) REFERENCES tasks(id) ON DELETE CASCADE;