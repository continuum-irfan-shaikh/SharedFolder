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

  CREATE MATERIALIZED VIEW IF NOT EXISTS script_execution_results_by_task_instance_id AS
    SELECT *
    FROM script_execution_results
    WHERE managed_endpoint_id IS NOT NULL AND task_instance_id IS NOT NULL
    PRIMARY KEY (task_instance_id, managed_endpoint_id);
