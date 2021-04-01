USE trading_db;

CREATE TABLE IF NOT EXISTS eth_candles (
  time DATETIME PRIMARY KEY NOT NULL,
  open FLOAT,
  close FLOAT,
  high FLOAT,
  low FLOAT,
  volume FLOAT
);
