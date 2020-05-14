
CREATE TABLE IF NOT EXISTS task_counters(
  partner_id text,
  endpoint_id uuid,
  internal_tasks counter,
  external_tasks counter,
  PRIMARY KEY (partner_id, endpoint_id)
);
