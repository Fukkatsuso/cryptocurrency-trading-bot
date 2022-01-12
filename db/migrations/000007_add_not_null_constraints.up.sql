USE trading_db;

ALTER TABLE users MODIFY COLUMN
  session_id_hash VARCHAR(255) NOT NULL DEFAULT '';
