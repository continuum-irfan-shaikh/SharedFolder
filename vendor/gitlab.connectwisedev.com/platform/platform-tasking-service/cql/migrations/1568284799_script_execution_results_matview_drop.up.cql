CREATE TABLE IF NOT EXISTS script_execution_results_by_task_instance_id_mv
(
    managed_endpoint_id uuid,
    task_instance_id    uuid,
    execution_status    int,
    std_err             text,
    std_out             text,
    updated_at          timestamp,
    PRIMARY KEY (task_instance_id, managed_endpoint_id)
);
