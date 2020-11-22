package main

import (
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"os"
	"strings"
)

var (
	bdb     *gorm.DB
	conf    Config
	debug   bool
	homeDir string
)

func parseOptions() {
	flag.BoolVar(&debug, "debug", false, "debug")

	// Backend MySQL
	flag.StringVar(&conf.Backend.Host, "bh", "127.0.0.1", "TBD")
	flag.StringVar(&conf.Backend.User, "bu", "root", "TBD")
	flag.IntVar(&conf.Backend.Port, "bP", 3306, "MySQL tmport")
	flag.StringVar(&conf.Backend.Password, "bp", "root", "TBD")
	flag.StringVar(&conf.Backend.DB, "bdb", "smt", "TBD")
	flag.StringVar(&conf.Backend.Socket, "bs", "", "TBD")

	// Target MySQL
	flag.StringVar(&conf.Target.Host, "th", "127.0.0.1", "TBD")
	flag.StringVar(&conf.Target.User, "tu", "bench", "TBD")
	flag.IntVar(&conf.Target.Port, "tP", 3306, "MySQL tmport")
	flag.StringVar(&conf.Target.Password, "tp", "bench", "TBD")
	flag.StringVar(&conf.Target.DB, "tdb", "smt_bench", "TBD")
	flag.StringVar(&conf.Target.Socket, "ts", "", "TBD")

	flag.Parse()
}

func main() {
	// get current directory to set homeDir
	pwd, _ := os.Getwd()
	l := strings.LastIndex(pwd, "/")
	homeDir = pwd[:l]

	parseOptions()
	err := readConf()
	if err != nil {
		fmt.Println(err)
		return
	}

	bdb, err = connectBackendMySQL()
	if err != nil {
		panic(err)
	}
	defer bdb.Close()

	// benchmark start
	err = start()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("[Notice] Completed!!")
}
