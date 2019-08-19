package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	smp "github.com/tom--bo/sm-parser"
)

func start() error {
	// Get not completed senarios
	senarios, err := getQueuedSenarios()

	if err != nil {
		if err == RecordNotFound {
			fmt.Println("All senarios are completed!!")
			return nil
		}
		return err
	}

	// Start benchmark for senarios
	for i, s := range senarios {
		fmt.Printf("senario,%2d: \n", i)
		fmt.Println(s)

		err := beforeSenario(s)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		restartMySQL()

		prepare(s)
		afterPrepare(s)

		restartMySQL()

		bench(s)

		err = afterSenario(s)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	return nil
}

func beforeSenario(s Senario) error {
	dropSchemaIfExists()

	resetPersist()

	c := getCommand(s.BeforeSenarioCommand)
	cmds := strings.Split(c.CommandText, "\n")

	for _, s := range cmds {
		execSQL(s)
	}

	return nil
}

func prepare(s Senario) error {
	createSchema()
	prepareBench(s)

	return nil
}

func prepareBench(s Senario) {
	fmt.Println("prepareBench(s)")
	out, err := exec.Command(conf.Base.SysbenchPath+"sysbench",
		conf.Base.SysbenchSenarioDir + s.SysbenchSenario + ".lua",
		"--db-driver=mysql",
		"--table-size="+strconv.Itoa(int(s.TableSize)),
		"--tables="+strconv.Itoa(int(s.TableNum)),
		"--mysql-host="+conf.Target.Host,
		"--mysql-port="+strconv.Itoa(conf.Target.Port),
		"--mysql-user="+conf.Target.User,
		"--mysql-password="+conf.Target.Password,
		"--mysql-db="+conf.Target.DB,
		"prepare",
	).Output()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(out))
}

func afterPrepare(s Senario) error {
	c := getCommand(s.AfterPrepareCommand)
	cmds := strings.Split(c.CommandText, "\n")

	for _, s := range cmds {
		execSQL(s)
	}

	setVariables(s)
	return nil
}

func bench(s Senario) {
	fmt.Println("------- bench")

	cnt := getBenchCount(s)
	fmt.Printf("Senario count for this senario: %d\n", cnt)
	for i := cnt + 1; i <= int(s.Count); i++ {
		fmt.Println("benchmark num i: ", i)

		// start sysbench run
		err := run(i, s)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Something happen, Skip this senario!")
			break
		}
		// cool down
		time.Sleep(120 * time.Second)
	}
}

func run(i int, s Senario) error {
	fmt.Println(s.ID, ": ", i)

	// exec sysbench
	out, err := exec.Command(conf.Base.SysbenchPath+"sysbench",
		conf.Base.SysbenchSenarioDir + s.SysbenchSenario + ".lua",
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
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(out))

	var r smp.Result
	smp.ParseOutput(&r, string(out))

	ret := Result{SenarioId: s.ID, SenarioCount: int64(i)}
	ret = mapResult(ret, r)
	registerResult(ret) // if error happen, die

	return nil
}

func afterSenario(s Senario) error {
	err := updateStatus(s)
	if err != nil {
		return err
	}

	c := getCommand(s.AfterSenarioCommand)
	cmds := strings.Split(c.CommandText, "\n")

	for _, s := range cmds {
		execSQL(s)
	}
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
