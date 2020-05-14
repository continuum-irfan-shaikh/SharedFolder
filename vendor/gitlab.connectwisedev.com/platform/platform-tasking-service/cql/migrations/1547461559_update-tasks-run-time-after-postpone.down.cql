DROP MATERIALIZED VIEW IF EXISTS tasks_by_id;
DROP MATERIALIZED VIEW IF EXISTS tasks_by_id_managed_endpoint_id;
DROP MATERIALIZED VIEW IF EXISTS tasks_by_runtime;
DROP MATERIALIZED VIEW IF EXISTS tasks_order_by_last_task_instance_id;
DROP MATERIALIZED VIEW IF EXISTS task_by_runtime_unix;

ALTER TABLE tasks DROP original_next_run_time;

CREATE MATERIALIZED VIEW IF NOT EXISTS tasks_by_id AS
  SELECT *
  FROM tasks
  WHERE partner_id IS NOT NULL AND id IS NOT NULL AND managed_endpoint_id IS NOT NULL AND external_task IS NOT NULL
  PRIMARY KEY ((partner_id, id), managed_endpoint_id, external_task);

CREATE MATERIALIZED VIEW IF NOT EXISTS tasks_by_id_managed_endpoint_id AS
  SELECT *
  FROM tasks
  WHERE partner_id IS NOT NULL AND id IS NOT NULL AND managed_endpoint_id IS NOT NULL AND external_task IS NOT NULL
  PRIMARY KEY ((partner_id, id, managed_endpoint_id), external_task);

CREATE MATERIALIZED VIEW IF NOT EXISTS tasks_by_runtime AS
  SELECT *
  FROM tasks
  WHERE run_time_unix IS NOT NULL AND partner_id IS NOT NULL AND id IS NOT NULL AND managed_endpoint_id IS NOT NULL AND external_task IS NOT NULL
  PRIMARY KEY (run_time_unix, id, partner_id, managed_endpoint_id, external_task);

CREATE MATERIALIZED VIEW platform_tasking_db.tasks_order_by_last_task_instance_id AS
  SELECT *
  FROM platform_tasking_db.tasks
  WHERE partner_id IS NOT NULL AND id IS NOT NULL AND managed_endpoint_id IS NOT NULL AND last_task_instance_id IS NOT NULL AND external_task IS NOT NULL
  PRIMARY KEY ((partner_id, managed_endpoint_id, external_task), last_task_instance_id, id)
  WITH CLUSTERING ORDER BY (last_task_instance_id DESC, id DESC);

CREATE MATERIALIZED VIEW IF NOT EXISTS platform_tasking_db.task_by_runtime_unix AS
  SELECT *
  FROM platform_tasking_db.tasks
  WHERE run_time_unix IS NOT NULL AND partner_id IS NOT NULL AND external_task IS NOT NULL AND managed_endpoint_id IS NOT NULL AND id IS NOT NULL
  PRIMARY KEY (partner_id, external_task, run_time_unix, managed_endpoint_id, id);