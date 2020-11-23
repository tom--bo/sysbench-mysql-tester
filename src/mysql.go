package main

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

var (
	// Custom Error
	RecordNotFound = errors.New("record not found")
)

type Senario struct {
	ID              int64  `gorm:"id"`
	SysbenchSenario string `gorm:"sysbench_senario"`
	TableNum        int64  `gorm:"table_num"`
	TableSize       int64  `gorm:"table_size"`
	ThreadNum       int64  `gorm:"thread_num"`
	TimeSecond      int64  `gorm:"time_second"`
	MycnfId         int64  `gorm:"mycnf_id"`
	ExpCount        int64  `gorm:"exp_count"`
}

type SenarioStatus string

const (
	QUEUED    SenarioStatus = "QUEUED"
	RUNNING   SenarioStatus = "RUNNING"
	SKIPPED   SenarioStatus = "SKIPPED"
	ERROR     SenarioStatus = "ERROR"
	COMPLETED SenarioStatus = "COMPLETED"
)

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

func connectMySQL(c string) (*gorm.DB, error) {
	tdb, err := gorm.Open("mysql", c)
	if err != nil {
		return nil, err
	}
	return tdb, err
}

func connectBackendMySQL() (*gorm.DB, error) {
	mysqlHost := conf.Backend.User + ":" + conf.Backend.Password + "@tcp(" + conf.Backend.Host + ":" + strconv.Itoa(conf.Backend.Port) + ")/" + conf.Backend.DB
	if conf.Backend.Socket != "" {
		mysqlHost = conf.Backend.User + ":" + conf.Backend.Password + "@unix(" + conf.Backend.Socket + ")/" + conf.Backend.DB + "?loc=Local&parseTime=true"
	}
	return connectMySQL(mysqlHost)
}

func connectTargetMySQL() (*gorm.DB, error) {
	mysqlHost := conf.Target.User + ":" + conf.Target.Password + "@tcp(" + conf.Target.Host + ":" + strconv.Itoa(conf.Target.Port) + ")/"
	if conf.Target.Socket != "" {
		mysqlHost = conf.Target.User + ":" + conf.Target.Password + "@unix(" + conf.Target.Socket + ")/"
	}
	return connectMySQL(mysqlHost)
}

func getQueuedSenarios() ([]Senario, error) {
	s := []Senario{}

	bdb.Raw("SELECT * FROM senario WHERE status like 'QUEUED' ORDER BY id").Scan(&s)
	if len(s) == 0 {
		return nil, RecordNotFound
	}

	return s, nil
}

func updateStatus(sid int64, status SenarioStatus, msg string) error {
	sql := fmt.Sprintf("UPDATE senario SET status = '%s', message = '%s' where id = %d", status, msg, sid)
	// fmt.Println(sql)
	bdb.Exec(sql)
	return nil
}

func registerResult(r Result) {
	bdb.Create(&r)
}

// Exec restart command and wait
func restartMySQL() error {
	fmt.Println("[Notice] Restart MySQL ...")
	tdb, err := connectTargetMySQL()
	if err != nil {
		return err
	}
	defer tdb.Close()

	tdb.Exec("RESTART;")

	timeoutCnt := 120
	for i := 0; i < timeoutCnt; i++ {
		time.Sleep(5 * time.Second)
		_, err := connectTargetMySQL()
		if err != nil {
			if i >= timeoutCnt-1 {
				return errors.New("Restart check timeout")
			}
			continue
		} else {
			break
		}
	}

	return nil
}

func innodb_redo_log(b bool) error {
	tdb, err := connectTargetMySQL()
	if err != nil {
		return err
	}
	defer tdb.Close()

	str := "disable"
	if b {
		str = "enable"
	}

	// drop db if exists
	tdb.Exec("ALTER INSTANCE " + str + " INNODB REDO_LOG;")
	return nil
}

// create DB
func recreateSchema() error {
	tdb, err := connectTargetMySQL()
	if err != nil {
		return err
	}
	defer tdb.Close()

	// drop db if exists
	tdb.Exec("DROP DATABASE IF EXISTS " + conf.Target.DB + ";")
	// create db
	tdb.Exec("CREATE DATABASE " + conf.Target.DB + ";")

	return nil
}
