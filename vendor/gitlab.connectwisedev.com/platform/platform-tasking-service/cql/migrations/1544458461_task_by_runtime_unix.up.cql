//
// ATTENTION!!! DO NOT EDIT THIS FILE. IT IS CREATED FOLLOWING SCHEMA MIGRATION STRATEGY
//

CREATE MATERIALIZED VIEW IF NOT EXISTS platform_tasking_db.task_by_runtime_unix AS
    SELECT *
    FROM platform_tasking_db.tasks
    WHERE run_time_unix IS NOT NULL AND partner_id IS NOT NULL AND external_task IS NOT NULL AND managed_endpoint_id IS NOT NULL AND id IS NOT NULL
    PRIMARY KEY (partner_id, external_task, run_time_unix, managed_endpoint_id, id)