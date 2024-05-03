-- Questdb sql goes here
CREATE TABLE IF NOT EXISTS 'sm_powermeter' (
  home_id SYMBOL,
  ts TIMESTAMP,
  power DOUBLE,
  min_power DOUBLE,
  average_power DOUBLE,
  max_power DOUBLE,
  last_meter_consumption DOUBLE,
  last_meter_production DOUBLE,
  accumulated_consumption DOUBLE,
  accumulated_production DOUBLE,
  accumulated_cost DOUBLE,
  accumulated_production_last_hour DOUBLE,
  accumulated_consumption_last_hour DOUBLE,
  currency STRING
) timestamp (ts) PARTITION BY MONTH WAL;
