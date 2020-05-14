CREATE TABLE IF NOT EXISTS platform_tasking_db.scheduler (
  id int PRIMARY KEY,
  last_update timestamp
);

INSERT INTO platform_tasking_db.scheduler (id , last_update) VALUES (0, toTimestamp(now() ) );