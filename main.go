// Copyright 2023 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/format"
	driver "github.com/pingcap/tidb/types/parser_driver"
	"log"
	"strings"
)

type FingerprintVisitor struct {
}

var sqlExprValues []string

func (f *FingerprintVisitor) Enter(n ast.Node) (node ast.Node, skipChildren bool) {
	if v, ok := n.(*driver.ValueExpr); ok {
		sqlExprValues = append(sqlExprValues, strings.Split(v.String(), " ")[1])
		v.Type.Charset = ""
		v.SetValue([]byte("?"))
	}
	return n, false
}

func (f *FingerprintVisitor) Leave(n ast.Node) (node ast.Node, ok bool) {
	return n, true
}

type SQLParserService struct {
	Template   string
	Parameters []string
}

func NewSQLParserService() *SQLParserService {
	return &SQLParserService{
		Template:   "",
		Parameters: nil,
	}
}

func main() {
	//buf, done := FingerAndParameter()
	//if done {
	//	return
	//}
	//fmt.Println(buf.String())
	//fmt.Println(sqlExprValues)

	parserService := NewSQLParserService()
	_, _ = parserService.FingerAndParameter("SELECT * FROM t1 WHERE id IN (SELECT id FROM users WHERE username IN ('emorydu', 'hello') );")
	fmt.Println(parserService.Template)
	fmt.Println(parserService.Parameters)
	parserService.FingerAndParameter("SELECT * FROM users WHERE name = 'world'")
	fmt.Println(parserService.Template)
	fmt.Println(parserService.Parameters)

}

func (r *SQLParserService) FingerAndParameter(sql string) (*bytes.Buffer, bool) {
	r.Parameters = []string{}
	//sql := "SELECT * FROM t1 WHERE id IN (SELECT id FROM users WHERE username IN ('emorydu', 'hello') );"
	p := parser.New()
	stmt, err := p.ParseOneStmt(sql, "", "")
	if err != nil {
		log.Println(err)
		return nil, true
	}
	stmt.Accept(&FingerprintVisitor{})

	buf := new(bytes.Buffer)
	restoreCtx := format.NewRestoreCtx(format.RestoreKeyWordUppercase|format.RestoreNameBackQuotes, buf)
	err = stmt.Restore(restoreCtx)
	if err != nil {
		log.Println(err)
		return nil, true
	}
	// TODO
	r.Template = buf.String()
	r.Parameters = sqlExprValues
	sqlExprValues = []string{}
	return buf, false
}
