package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	bdb   *gorm.DB
	conf  Config
	debug bool
)

func doFileExist(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func readConf() {
	confFile := "/etc/smt.cnf"
	if doFileExist(confFile) {
		f, err := os.Open(confFile)
		if err != nil {
			fmt.Print("[Warning] Can't open the config file: ")
			fmt.Println(err)
		}
		defer f.Close()

		_, err = toml.DecodeFile(confFile, &conf)
		if err != nil {
			fmt.Println(err)
		}

		if err != nil {
			fmt.Print("[Warning] Can't read config file: ")
			fmt.Println(err)
		}
	}
}

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

func gormConnect() {
	// connect to backend db
	var err error
	mysqlHost := conf.Backend.User + ":" + conf.Backend.Password + "@tcp(" + conf.Backend.Host + ":" + strconv.Itoa(conf.Backend.Port) + ")/" + conf.Backend.DB
	if conf.Backend.Socket != "" {
		mysqlHost = conf.Backend.User + ":" + conf.Backend.Password + "@unix(" + conf.Backend.Socket + ")/" + conf.Backend.DB + "?loc=Local&parseTime=true"
	}
	fmt.Println(mysqlHost)
	bdb, err = gorm.Open("mysql", mysqlHost)
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	parseOptions()
	readConf()

	// Connect to MySQL
	gormConnect()
	defer bdb.Close()

	// benchmark start
	err := start()
	if err != nil {
		fmt.Println(err)
	}
}
