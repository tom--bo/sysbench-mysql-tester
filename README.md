# What is sysbench-mysql-tester (SMT)

sysbench-mysql-tester(SMT) executes sysbench senarios to MySQL.
This tool assume to use sysbench-mysql or sysbench-tpcc for MySQL benchmark.

## How to insatall

```
git clone ...
go build -o smt ...
./smt
```

## How to setup

- Prepare a MySQL as backend database
- Create backend database like `CREATE DATABASE smt`
- Execute setup scripts
  - `mysql -u mysql -p smt < scripts/definition.sql`
  - `mysql -u mysql -p smt < scripts/setup.sql`
- Insert records to senario table as you like
- sysbench might not support caching_sha2_password, then modify your my.cnf
  - ```
    [mysqld]
    default-authentication-plugin = mysql_native_password
    ```


## How to run

```
smt -h 192.168.1.3 -tu mysql -tp password -bh 192.168.1.2 -bu mysql -bp password
```

You can configure options by json-file.

```toml
# /etc/smt.cnf
[Base]
sysbench_path = "/usr/bin/"
sysbench_senario_dir = "/usr/share/sysbench/"

[Target]
host = "192.168.1.100"
user = "smt"
port = 3306
password = "smt"
database = "sysbench"

[Backend]
host = "127.0.0.1"
user = "smt"
port = 3306
password = "smt"
database = "smt"
```

Then, you need not to specify any option.



## Support

For now, SMT support only MySQL 8.0


