
DROP MATERIALIZED VIEW IF EXISTS platform_tasking_db.tasks_by_runtime;
DROP MATERIALIZED VIEW IF EXISTS platform_tasking_db.tasks_by_target;
DROP MATERIALIZED VIEW IF EXISTS platform_tasking_db.tasks_by_id_target;
DROP MATERIALIZED VIEW IF EXISTS platform_tasking_db.tasks_by_id;
COPY platform_tasking_db.tasks (partner_id,id,target,created_at,created_by,description,dynamic_groups,end_run_time,external_task,last_task_instance_id,location,name,origin_id,parameters,regularity,result_webhook,run_time,run_time_unix,schedule,start_run_time,state,trigger,type) TO '/tmp/tasks-up.csv' WITH NULL = 'null' AND HEADER = TRUE;
DROP TABLE IF EXISTS platform_tasking_db.tasks;


CREATE TABLE IF NOT EXISTS platform_tasking_db.tasks (
    partner_id text,
    id uuid,
    managed_endpoint_id uuid,
    created_at timestamp,
    created_by text,
    description text,
    targets set<text>,
    end_run_time timestamp,
    external_task boolean,
    last_task_instance_id uuid,
    location text,
    name text,
    origin_id uuid,
    parameters text,
    regularity int,
    result_webhook text,
    run_time timestamp,
    run_time_unix timestamp,
    schedule text,
    start_run_time timestamp,
    state int,
    trigger text,
    type text,
    target_type int,
    PRIMARY KEY (partner_id, external_task, managed_endpoint_id, id)
   );


COPY platform_tasking_db.tasks (partner_id,id,managed_endpoint_id,created_at,created_by,description,targets,end_run_time,external_task,last_task_instance_id,location,name,origin_id,parameters,regularity,result_webhook,run_time,run_time_unix,schedule,start_run_time,state,trigger,type) FROM '/tmp/tasks-up.csv' WITH NULL = 'null' AND HEADER = TRUE AND TTL=7775984;

CREATE MATERIALIZED VIEW IF NOT EXISTS platform_tasking_db.tasks_by_runtime AS
    SELECT *
    FROM tasks
    WHERE run_time_unix IS NOT NULL AND partner_id IS NOT NULL AND id IS NOT NULL AND managed_endpoint_id IS NOT NULL AND external_task IS NOT NULL
    PRIMARY KEY (run_time_unix, id, partner_id, managed_endpoint_id, external_task);

CREATE MATERIALIZED VIEW IF NOT EXISTS platform_tasking_db.tasks_by_id AS
    SELECT *
    FROM tasks
    WHERE partner_id IS NOT NULL AND id IS NOT NULL AND managed_endpoint_id IS NOT NULL AND external_task IS NOT NULL
    PRIMARY KEY ((partner_id, id), managed_endpoint_id, external_task);

CREATE MATERIALIZED VIEW IF NOT EXISTS platform_tasking_db.tasks_by_id_managed_endpoint_id AS
    SELECT *
    FROM tasks
    WHERE partner_id IS NOT NULL AND id IS NOT NULL AND managed_endpoint_id IS NOT NULL AND external_task IS NOT NULL
    PRIMARY KEY ((partner_id, id, managed_endpoint_id), external_task);

CREATE MATERIALIZED VIEW IF NOT EXISTS platform_tasking_db.tasks_order_by_last_task_instance_id AS
    SELECT *
    FROM tasks
    WHERE partner_id IS NOT NULL AND id IS NOT NULL AND managed_endpoint_id IS NOT NULL AND last_task_instance_id IS NOT NULL AND external_task IS NOT NULL
    PRIMARY KEY ((partner_id, managed_endpoint_id, external_task), last_task_instance_id, id)
WITH CLUSTERING ORDER BY (last_task_instance_id DESC, id DESC);
