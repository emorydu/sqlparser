// Copyright 2023 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	p := Parser()

	// MySQL
	response, err := p.FingerAndParameter("SELECT * FROM t1 WHERE id IN (SELECT id FROM users WHERE username IN ('emorydu', 'hello') );")
	if err != nil {
		panic(err)
	}
	fmt.Println(response.Result())
	response, err = p.FingerAndParameter("SELECT * FROM users WHERE name = 'world'")
	if err != nil {
		panic(err)
	}
	fmt.Println(response.Result())

	// SQL Server
	response, err = p.FingerAndParameter("SELECT name FROM person WHERE countryid in (SELECT contryid FROM country WHERE contryname = '中国')")
	if err != nil {
		panic(err)
	}
	fmt.Println(response.Result())

	// SQLite
	response, err = p.FingerAndParameter("insert into users VALUEs('a', 'b', 'c')")
	if err != nil {
		panic(err)
	}
	fmt.Println(response.Result())

	// Postgres
	response, err = p.FingerAndParameter(`
UPDATE employees
SET salary = (
    SELECT salary * 1.1
    FROM employees
    WHERE department = 'Sales'
)
WHERE department = 'Sales';

`)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.Result())

	response, err = p.FingerAndParameter(`
SELECT COUNT(CGUID)FROM AOS_NOTICE WHERE ISTATUS=1 AND(CORGNID=418704614421011733 OR CORGNID='' OR CORGNID IS NULL)AND CTOUSERGUID=1;
`)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.Result())

}
