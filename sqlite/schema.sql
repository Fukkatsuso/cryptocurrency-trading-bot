CREATE TABLE `eth_candles` (
  `time` TEXT NOT NULL,
  `open` REAL DEFAULT NULL,
  `close` REAL DEFAULT NULL,
  `high` REAL DEFAULT NULL,
  `low` REAL DEFAULT NULL,
  `volume` REAL DEFAULT NULL,
  PRIMARY KEY (`time`)
);

CREATE TABLE `signal_events` (
  `time` TEXT NOT NULL,
  `product_code` TEXT DEFAULT NULL,
  `side` TEXT DEFAULT NULL,
  `price` REAL DEFAULT NULL,
  `size` REAL DEFAULT NULL,
  PRIMARY KEY (`time`)
);

CREATE TABLE `trade_params` (
  `trade_enable` INTEGER NOT NULL DEFAULT '1',
  `product_code` TEXT NOT NULL,
  `size` REAL NOT NULL,
  `sma_enable` INTEGER NOT NULL DEFAULT '0',
  `sma_period1` INTEGER NOT NULL,
  `sma_period2` INTEGER NOT NULL,
  `sma_period3` INTEGER NOT NULL,
  `ema_enable` INTEGER NOT NULL DEFAULT '0',
  `ema_period1` INTEGER NOT NULL,
  `ema_period2` INTEGER NOT NULL,
  `ema_period3` INTEGER NOT NULL,
  `bbands_enable` INTEGER NOT NULL DEFAULT '0',
  `bbands_n` INTEGER NOT NULL,
  `bbands_k` REAL NOT NULL,
  `ichimoku_enable` INTEGER NOT NULL DEFAULT '0',
  `rsi_enable` INTEGER NOT NULL DEFAULT '0',
  `rsi_period` INTEGER NOT NULL,
  `rsi_buy_thread` REAL NOT NULL,
  `rsi_sell_thread` REAL NOT NULL,
  `macd_enable` INTEGER NOT NULL DEFAULT '0',
  `macd_fast_period` INTEGER NOT NULL,
  `macd_slow_period` INTEGER NOT NULL,
  `macd_signal_period` INTEGER NOT NULL,
  `created_at` TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `stop_limit_percent` REAL NOT NULL DEFAULT 0
);
