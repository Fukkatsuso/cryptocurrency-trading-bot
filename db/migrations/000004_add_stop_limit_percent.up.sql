USE trading_db;

ALTER TABLE trade_params ADD COLUMN
  stop_limit_percent FLOAT NOT NULL DEFAULT 0;
