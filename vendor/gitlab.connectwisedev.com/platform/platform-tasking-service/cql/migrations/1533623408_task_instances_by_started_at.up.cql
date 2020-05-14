DROP MATERIALIZED VIEW IF EXISTS task_instances_by_id;

ALTER TABLE task_instances ADD partner_id text;

CREATE MATERIALIZED VIEW IF NOT EXISTS task_instances_started_at AS
SELECT *
FROM task_instances
WHERE id IS NOT NULL AND task_id IS NOT NULL AND started_at IS NOT NULL AND partner_id IS NOT NULL
PRIMARY KEY (partner_id, started_at, id, task_id);

CREATE MATERIALIZED VIEW IF NOT EXISTS task_instances_by_id AS
SELECT *
FROM task_instances
WHERE id IS NOT NULL AND task_id IS NOT NULL AND started_at IS NOT NULL
PRIMARY KEY (id, task_id, started_at);