package main

import (
	"errors"
)

var (
	// Custom Error
	RecordNotFound = errors.New("record not found")
)

type Config struct {
	Base BaseConfig
	Backend BackendMySQL
	Target  TargetMySQL
}

type BaseConfig struct {
	SysbenchPath string `toml:"sysbench_path"`
	SysbenchSenarioDir string `toml:"sysbench_senario_dir"`
}

type BackendMySQL struct {
	Host     string `toml:"host"`
	User     string `toml:"user"`
	Port     int    `toml:"port"`
	Password string `toml:"password"`
	DB       string `toml:"database"`
	Socket   string `toml:"socket"`
}

type TargetMySQL struct {
	Host     string `toml:"host"`
	User     string `toml:"user"`
	Port     int    `toml:"port"`
	Password string `toml:"password"`
	DB       string `toml:"database"`
	Socket   string `toml:"socket"`
}

type Count struct {
	Cnt int64 `gorm:"cnt"`
}

type Command struct {
	ID          int64  `gorm:"id"`
	CommandText string `gorm:"command_text"`
	lavel       string `gorm:"lavel"`
}

type Senario struct {
	ID                   int64  `gorm:"id"`
	SysbenchSenario      string `gorm:"sysbench_senario"`
	TableNum             int64  `gorm:"table_num"`
	TableSize            int64  `gorm:"table_size"`
	ThreadNum            int64  `gorm:"thread_num"`
	TimeSecond           int64  `gorm:"time_second"`
	Count                int64  `gorm:"count"`
	BeforeSenarioCommand int64  `gorm:"before_senario_command"`
	AfterPrepareCommand  int64  `gorm:"after_prepare_command"`
	AfterSenarioCommand  int64  `gorm:"after_senario_command"`
}

type Result struct {
	ID                    int64   `gorm:"id"`
	SenarioId             int64   `gorm:"senario_id"`
	SenarioCount          int64   `gorm:"senario_count"`
	SysbenchVersion       string  `gorm:"sysbench_version"`
	LuajitVersion         string  `gorm:"luajit_version"`
	Threads               int64   `gorm:"threads"`
	TotalRead             int64   `gorm:"total_read"`
	TotalWrite            int64   `gorm:"total_write"`
	TotalOther            int64   `gorm:"total_other"`
	TotalTx               int64   `gorm:"total_tx"`
	Tps                   float64 `gorm:"tps"`
	TotalQuery            int64   `gorm:"total_query"`
	Qps                   float64 `gorm:"qps"`
	IgnoredErrors         int64   `gorm:"ignored_errors"`
	Reconnects            int64   `gorm:"reconnects"`
	TotalTime             float64 `gorm:"total_time"`
	TotalEvents           int64   `gorm:"total_events"`
	MinLatency            float64 `gorm:"min_latency"`
	AvgLatency            float64 `gorm:"avg_latency"`
	MaxLatency            float64 `gorm:"max_latency"`
	P95thLatency          float64 `gorm:"p95th_latency"`
	SumLatency            float64 `gorm:"sum_latency"`
	ThreadsEventsAvg      float64 `gorm:"threads_events_avg"`
	ThreadsEventsStddev   float64 `gorm:"threads_events_stddev"`
	ThreadsExecTimeAvg    float64 `gorm:"threads_exec_time_avg"`
	ThreadsExecTimeStddev float64 `gorm:"threads_exec_time_stddev"`
}

type Variable struct {
	Name  string `gorm:"name"`
	Value string `gorm:"value"`
}

type One struct {
	One int64 `grm:"one"`
}
