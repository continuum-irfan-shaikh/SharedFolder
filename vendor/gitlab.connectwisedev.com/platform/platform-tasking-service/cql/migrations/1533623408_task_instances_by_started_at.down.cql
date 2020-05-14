DROP MATERIALIZED VIEW IF EXISTS task_instances_by_id;
DROP MATERIALIZED VIEW IF EXISTS task_instances_started_at;


ALTER TABLE task_instances DROP partner_id;


CREATE MATERIALIZED VIEW IF NOT EXISTS task_instances_by_id AS
SELECT *
FROM task_instances
WHERE id IS NOT NULL AND task_id IS NOT NULL AND started_at IS NOT NULL
PRIMARY KEY (id, task_id, started_at);