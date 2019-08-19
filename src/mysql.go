package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"os"
	"strconv"
	"time"
)

func connectTargetMySQL() (*gorm.DB, error) {
	mysqlHost := conf.Target.User + ":" + conf.Target.Password + "@tcp(" + conf.Target.Host + ":" + strconv.Itoa(conf.Target.Port) + ")/"
	if conf.Target.Socket != "" {
		mysqlHost = conf.Target.User + ":" + conf.Target.Password + "@unix(" + conf.Target.Socket + ")/"
	}

	tdb, err := gorm.Open("mysql", mysqlHost)
	if err != nil {
		return nil, err
	}
	return tdb, err
}

func resetPersist() {
	execSQL("RESET PERSIST;")
}

func setVariables(s Senario) error {
	tdb, err := connectTargetMySQL()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer tdb.Close()

	// set persistent parameters
	vars := getVariables(s)
	for _, v := range vars {
		sql := fmt.Sprintf("set persist %s = %s", v.Name, v.Value)
		tdb.Exec(sql)
	}

	return nil
}

// Exec restart command and wait
func restartMySQL() {
	fmt.Println("restart MySQL ...")
	tdb, err := connectTargetMySQL()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer tdb.Close()

	tdb.Exec("RESTART;")

	cnt := 0
	for !isMySQLUp() {
		cnt += 1
		if cnt > 120 {
			fmt.Println("MySQL will not start.")
			os.Exit(1)
		}
		fmt.Println("Target MySQL is not restarted. Wait 30 sec...")
		time.Sleep(30 * time.Second)
	}
}

// create DB
func createSchema() {
	tdb, err := connectTargetMySQL()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer tdb.Close()

	sql := "CREATE DATABASE " + conf.Target.DB + ";"
	tdb.Exec(sql)
}

// drop DB
func dropSchemaIfExists() {
	tdb, err := connectTargetMySQL()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer tdb.Close()

	sql := "DROP DATABASE IF EXISTS " + conf.Target.DB + ";"
	tdb.Exec(sql)
}

func isMySQLUp() bool {
	tdb, err := connectTargetMySQL()
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer tdb.Close()

	o := One{}
	bdb.Raw("SELECT 1 as one").Scan(&o)
	if o.One == 1 {
		return true
	}

	return false
}
