CREATE TABLE IF NOT EXISTS task_instances_started_at_mv
(
    task_id         uuid,
    started_at      timestamp,
    id              uuid,
    device_statuses frozen<map<uuid, int>>,
    failure_count   int,
    last_run_time   timestamp,
    name            text,
    origin_id       uuid,
    overall_status  int,
    partner_id      text,
    status          int,
    success_count   int,
    targets         set<text>,
    triggered_by    text,
    PRIMARY KEY (partner_id, started_at, id, task_id)
) WITH compaction = { 'class' :  'LeveledCompactionStrategy'  } AND CLUSTERING ORDER BY (started_at DESC, id DESC, task_id ASC);

CREATE TABLE IF NOT EXISTS task_instances_by_id_mv
(
    task_id         uuid,
    started_at      timestamp,
    id              uuid,
    device_statuses frozen<map<uuid, int>>,
    failure_count   int,
    last_run_time   timestamp,
    name            text,
    origin_id       uuid,
    overall_status  int,
    partner_id      text,
    status          int,
    success_count   int,
    targets         set<text>,
    triggered_by    text,
    PRIMARY KEY (id, task_id, started_at)
) WITH compaction = { 'class' :  'LeveledCompactionStrategy'  } AND CLUSTERING ORDER BY (task_id ASC, started_at DESC);