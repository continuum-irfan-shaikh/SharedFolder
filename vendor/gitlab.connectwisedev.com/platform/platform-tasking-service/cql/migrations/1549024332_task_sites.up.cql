CREATE TABLE IF NOT EXISTS task_sites (
  task_id uuid PRIMARY KEY,
  sites   frozen<set<text>>
)