DROP MATERIALIZED VIEW platform_tasking_db.task_instances_started_at;
DROP MATERIALIZED VIEW platform_tasking_db.task_instances_by_id;


ALTER TABLE platform_tasking_db.task_instances DROP device_statuses;
ALTER TABLE platform_tasking_db.task_instances DROP failure_count;
ALTER TABLE platform_tasking_db.task_instances DROP success_count;
ALTER TABLE platform_tasking_db.task_instances DROP overall_status;


CREATE MATERIALIZED VIEW IF NOT EXISTS task_instances_started_at AS
  SELECT * FROM platform_tasking_db.task_instances
  WHERE id IS NOT NULL AND task_id IS NOT NULL AND started_at IS NOT NULL AND partner_id IS NOT NULL
  PRIMARY KEY (partner_id, started_at, id, task_id)
  WITH CLUSTERING ORDER BY (started_at DESC, id DESC, task_id ASC);

CREATE MATERIALIZED VIEW IF NOT EXISTS task_instances_by_id AS
  SELECT * FROM platform_tasking_db.task_instances
  WHERE id IS NOT NULL AND task_id IS NOT NULL AND started_at IS NOT NULL
  PRIMARY KEY (id, task_id, started_at)
  WITH CLUSTERING ORDER BY (task_id ASC, started_at DESC);