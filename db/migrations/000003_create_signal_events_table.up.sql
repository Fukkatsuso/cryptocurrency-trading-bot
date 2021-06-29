USE trading_db;

CREATE TABLE IF NOT EXISTS signal_events (
  time DATETIME PRIMARY KEY NOT NULL,
  product_code VARCHAR(50),
  side VARCHAR(50),
  price FLOAT,
  size FLOAT
);
