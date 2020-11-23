package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	smp "github.com/tom--bo/sm-parser"
)

func start() error {
	senarios, err := getQueuedSenarios()
	if err != nil {
		if err == RecordNotFound {
			fmt.Println("All senarios are completed!!")
			return nil
		} else {
			return err
		}
	}

	for _, s := range senarios {
		fmt.Printf("[Notice] senario: %04d  ------\n", s.ID)
		updateStatus(s.ID, RUNNING, "")

		// send cnf with logs off
		err = sendTmpMycnf(s.MycnfId)
		if err != nil {
			updateStatus(s.ID, ERROR, "Send my.cnf: "+err.Error())
			continue
		}
		// disable redo-log
		innodb_redo_log(false)
		recreateSchema()
		restartMySQL()

		err = prepareBenchmark(s)
		if err != nil {
			fmt.Println(err.Error())
			updateStatus(s.ID, ERROR, "Prepare: "+err.Error())
			continue
		}

		err = sendMycnf(s.MycnfId)
		if err != nil {
			updateStatus(s.ID, ERROR, "Send my.cnf: "+err.Error())
			continue
		}
		// enable redo-log
		innodb_redo_log(true)
		restartMySQL()
		err = benchmark(s)
		if err != nil {
			fmt.Println(err.Error())
			updateStatus(s.ID, ERROR, "Benchmark: "+err.Error())
			continue
		} else {
			updateStatus(s.ID, COMPLETED, "")
		}
	}

	return nil
}

func sendTmpMycnf(mycnfId int64) error {
	origName := fmt.Sprintf("my_%04d.cnf", mycnfId)
	fileName := "my_tmp.cnf"
	base := homeDir + "/mycnfs/"

	logoffs := []string{
		"## temporal variables to logs off\n",
		"disable_log_bin\n",
		"slow_query_log = off\n",
		"general_log = off\n",
		"innodb_doublewrite = off\n",
		"innodb_flush_log_at_trx_commit = 2\n",
	}

	// make tmp cnf from original my_xxxx.cnf
	err := exec.Command("cp", base+origName, base+fileName).Run()
	if err != nil {
		fmt.Println(err)
		return err
	}
	// add log off variables

	f, err := os.OpenFile(base+fileName, os.O_APPEND|os.O_WRONLY, 755)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()
	for _, val := range logoffs {
		if _, err := f.Write([]byte(val)); err != nil {
			fmt.Println(err.Error())
			return err
		}
	}

	// send
	srcPath := homeDir + "/mycnfs/" + fileName
	dstPath := conf.Scp.User + "@" + conf.Target.Host + ":" + conf.Scp.Path
	err = exec.Command("sshpass", "-p", conf.Scp.Password, "scp", srcPath, dstPath).Run()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func sendMycnf(mycnfId int64) error {
	fileName := fmt.Sprintf("my_%04d.cnf", mycnfId)

	srcPath := homeDir + "/mycnfs/" + fileName
	dstPath := conf.Scp.User + "@" + conf.Target.Host + ":" + conf.Scp.Path
	err := exec.Command("sshpass", "-p", conf.Scp.Password, "scp", srcPath, dstPath).Run()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func prepareBenchmark(s Senario) error {
	fmt.Println("prepareBenchmark")

	time.Sleep(5 * time.Second)

	out, err := exec.Command(conf.Base.SysbenchPath+"sysbench",
		"--db-driver=mysql",
		"--table-size="+strconv.Itoa(int(s.TableSize)),
		"--tables="+strconv.Itoa(int(s.TableNum)),
		"--mysql-host="+conf.Target.Host,
		"--mysql-port="+strconv.Itoa(conf.Target.Port),
		"--mysql-user="+conf.Target.User,
		"--mysql-password="+conf.Target.Password,
		"--mysql-db="+conf.Target.DB,
		s.SysbenchSenario,
		"prepare").Output()
	if err != nil {
		fmt.Println(string(out))
		return err
	}
	return nil
}

func benchmark(s Senario) error {
	fmt.Println("[debug] benchmark() ---")
	for i := 1; i <= int(s.ExpCount); i++ {
		err := run(i, s)
		if err != nil {
			return err
		}
		// cool down
		time.Sleep(120 * time.Second)
	}
	return nil
}

func run(i int, s Senario) error {
	fmt.Println("[debug] run() ---")
	out, err := exec.Command(conf.Base.SysbenchPath+"sysbench",
		conf.Base.SysbenchSenarioDir+s.SysbenchSenario+".lua",
		"--db-driver=mysql",
		"--table-size="+strconv.Itoa(int(s.TableSize)),
		"--tables="+strconv.Itoa(int(s.TableNum)),
		"--threads="+strconv.Itoa(int(s.ThreadNum)),
		"--mysql-host="+conf.Target.Host,
		"--mysql-port="+strconv.Itoa(conf.Target.Port),
		"--mysql-user="+conf.Target.User,
		"--mysql-password="+conf.Target.Password,
		"--mysql-db="+conf.Target.DB,
		"--time="+strconv.Itoa(int(s.TimeSecond)),
		"run",
	).Output()
	if err != nil {
		return err
	}
	// fmt.Println(string(out))

	var r smp.Result
	smp.ParseOutput(&r, string(out))

	ret := Result{SenarioId: s.ID, SenarioCount: int64(i)}
	ret = mapResult(ret, r)
	registerResult(ret) // if error happen, die

	return nil
}

func mapResult(ret Result, r smp.Result) Result {
	// TODO: fix (same filedname)
	ret.SysbenchVersion = r.SysbenchVersion
	ret.LuajitVersion = r.LuajitVersion
	ret.Threads = int64(r.Threads)
	ret.TotalRead = int64(r.TotalRead)
	ret.TotalWrite = int64(r.TotalWrite)
	ret.TotalOther = int64(r.TotalOther)
	ret.TotalTx = int64(r.TotalTx)
	ret.Tps = r.Tps
	ret.TotalQuery = int64(r.TotalQuery)
	ret.Qps = r.Qps
	ret.IgnoredErrors = int64(r.IgnoredErrors)
	ret.Reconnects = int64(r.Reconnects)
	ret.TotalTime = r.TotalTime
	ret.TotalEvents = int64(r.TotalEvents)
	ret.MinLatency = r.MinLatency
	ret.AvgLatency = r.AvgLatency
	ret.MaxLatency = r.MaxLatency
	ret.P95thLatency = r.P95thLatency
	ret.SumLatency = r.SumLatency
	ret.ThreadsEventsAvg = r.ThreadsEventsAvg
	ret.ThreadsEventsStddev = r.ThreadsEventsStddev
	ret.ThreadsExecTimeAvg = r.ThreadsExecTimeAvg
	ret.ThreadsExecTimeStddev = r.ThreadsExecTimeStddev

	return ret
}
