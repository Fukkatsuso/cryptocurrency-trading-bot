USE trading_db;

CREATE TABLE IF NOT EXISTS trade_params (
  product_code VARCHAR(50) NOT NULL,
  size FLOAT NOT NULL,
  sma_enable BOOLEAN NOT NULL DEFAULT 0,
  sma_period1 INT NOT NULL,
  sma_period2 INT NOT NULL,
  sma_period3 INT NOT NULL,
  ema_enable BOOLEAN NOT NULL DEFAULT 0,
  ema_period1 INT NOT NULL,
  ema_period2 INT NOT NULL,
  ema_period3 INT NOT NULL,
  bbands_enable BOOLEAN NOT NULL DEFAULT 0,
  bbands_n INT NOT NULL,
  bbands_k INT NOT NULL,
  ichimoku_enable BOOLEAN NOT NULL DEFAULT 0,
  rsi_enable BOOLEAN NOT NULL DEFAULT 0,
  rsi_period INT NOT NULL,
  rsi_buy_thread FLOAT NOT NULL,
  rsi_sell_thread FLOAT NOT NULL,
  macd_enable BOOLEAN NOT NULL DEFAULT 0,
  macd_fast_period INT NOT NULL,
  macd_slow_period INT NOT NULL,
  macd_signal_period INT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
