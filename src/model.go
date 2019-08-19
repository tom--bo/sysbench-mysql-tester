package main

import "strconv"

func getQueuedSenarios() ([]Senario, error) {
	s := []Senario{}

	bdb.Raw("SELECT * FROM senario WHERE status like 'QUEUED'").Scan(&s)
	if len(s) == 0 {
		return nil, RecordNotFound
	}

	return s, nil
}

func getBenchCount(s Senario) int {
	var cnt Count
	bdb.Raw("SELECT count(*) as cnt FROM results where senario_id = ?", s.ID).Scan(&cnt)

	return int(cnt.Cnt)
}

func updateStatus(s Senario) error {
	sql := "UPDATE senario set status = 'COMPLETED' where id = " + strconv.Itoa(int(s.ID)) + ";"
	bdb.Exec(sql)
	return nil
}

func registerResult(r Result) {
	bdb.Create(&r)
}

func getVariables(s Senario) []Variable {
	v := []Variable{}
	bdb.Raw("SELECT v.name, sv.value from senario s inner join senario_variables sv on s.id = sv.senario_id inner join variables v on sv.variable_id = v.id where s.id = ?", s.ID).Scan(&v)

	return v
}

func getCommand(num int64) Command {
	c := Command{}
	bdb.Raw("SELECT * from command c where c.id = ?", num).Scan(&c)

	return c
}

func execSQL(s string) {
	bdb.Exec(s)
}
