CREATE TABLE IF NOT EXISTS senario (
	id INT AUTO_INCREMENT NOT NULL,
	sysbench_senario VARCHAR(30),
	table_num INT,
	table_size INT,
	thread_num INT,
	time_second INT,
	count INT,
	before_senario_command INT,
    after_prepare_command INT,
    after_senario_command INT,
	status ENUM('QUEUED', 'SKIPPED', 'COMPLETED') default 'QUEUED',
	PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS variables (
  id INT AUTO_INCREMENT NOT NULL,
  name VARCHAR(100),
  PRIMARY KEY(id),
  KEY(name)
);

CREATE TABLE IF NOT EXISTS commands (
  id INT AUTO_INCREMENT NOT NULL,
  command_text TEXT,
  lavel ENUM('BEFORE_SENARIO', 'AFTER_PREPARE', 'AFRTER_SENARIO'),
  PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS senario_variables (
  id INT AUTO_INCREMENT NOT NULL,
  senario_id INT,
  variable_id INT,
  value VARCHAR(1000),
  PRIMARY KEY(id),
  UNIQUE KEY(senario_id, variable_id)
);

CREATE TABLE IF NOT EXISTS results (
	id INT AUTO_INCREMENT NOT NULL,
	senario_id INT,
	senario_count INT,
	sysbench_version VARCHAR(20),
	luajit_version VARCHAR(20),
	threads INT,
	total_read BIGINT,
	total_write BIGINT,
	total_other BIGINT,
	total_tx BIGINT,
	tps DOUBLE,
	total_query BIGINT,
	qps DOUBLE,
	ignored_errors BIGINT,
	reconnects BIGINT,
	total_time DOUBLE,
	total_events BIGINT,
	min_latency DOUBLE,
	avg_latency DOUBLE,
	max_latency DOUBLE,
	p95th_latency DOUBLE,
	sum_latency DOUBLE,
	threads_events_avg DOUBLE,
	threads_events_stddev DOUBLE,
	threads_exec_time_avg DOUBLE,
	threads_exec_time_stddev DOUBLE,
	created_at DATETIME not null default current_timestamp,
	PRIMARY KEY(id),
	KEY(senario_id)
);

