# What is sysbench-mysql-tester (SMT)

sysbench-mysql-tester(SMT) executes sysbench senarios to MySQL.
This tool assume to use sysbench-mysql or sysbench-tpcc for MySQL benchmark.

## How to insatall

```
git clone https://github.com/tom--bo/sysbench-mysql-tester
cd sysbench_mysql_tester

make build
```


## How to setup

- Prepare a MySQL as backend database
- Create backend database like `CREATE DATABASE smt`
- Execute setup scripts
  - `mysql -u smt -p smt < init/definition.sql`
- Insert records to senario table as you like
- sysbench might not support caching_sha2_password, then modify your my.cnf
  - ```
    [mysqld]
    default-authentication-plugin = mysql_native_password
    ```
 TBD

## How to run

```sh
cd sysbench_mysql_tester
make run

# or

$ smt -h 192.168.1.3 -tu mysql -tp password -bh 192.168.1.2 -bu mysql -bp password
```

You can configure at `/etc/smt.cnf` by toml.  
Settings in the config file take precedence over command line options.

```toml
# please conf/smt_sample.cnf
# sample
[Base]
sysbench_path = "/usr/bin/"
sysbench_senario_dir = "/usr/share/sysbench/"

[Target]
host = "192.168.1.100"
user = "sysbench"
port = 3306
password = "sysbench"
database = "sysbench"

[Backend]
host = "127.0.0.1"
user = "smt"
port = 3306
password = "smt"
database = "smt"

[Scp]
user = "root"
password = "password"
path = "/etc/mysql/my.cnf"
```


