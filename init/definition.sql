CREATE TABLE IF NOT EXISTS senario (
	id INT AUTO_INCREMENT NOT NULL,
	sysbench_senario VARCHAR(30),
	table_num INT UNSIGNED NOT NULL,
	table_size INT UNSIGNED NOT NULL,
	thread_num INT UNSIGNED NOT NULL,
	time_second INT UNSIGNED NOT NULL,
	mycnf_id INT UNSIGNED NOT NULL,
	exp_count INT UNSIGNED NOT NULL,
	status ENUM('QUEUED', 'RUNNING', 'SKIPPED', 'ERROR','COMPLETED') default 'QUEUED',
	message VARCHAR(100) NOT NULL default '',
	created_at DATETIME NOT NULL DEFAULT current_timestamp,
	updated_at DATETIME NOT NULL DEFAULT current_timestamp ON UPDATE current_timestamp,
	PRIMARY KEY(id),
	key(status)
);

CREATE TABLE IF NOT EXISTS results (
	id INT AUTO_INCREMENT NOT NULL,
	senario_id INT UNSIGNED,
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
