USE trading_db;

ALTER TABLE trade_params MODIFY
  bbands_k INT NOT NULL;
