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
	"strings"
)

type SQLParser struct {
	sqlTemplate   string
	bitwiseParams []string
}

// Parser returns the default SQL parser.
func Parser() *SQLParser {
	return &SQLParser{
		sqlTemplate:   "",
		bitwiseParams: nil,
	}
}

type fingerprintVisitor struct {
	// parameters it is used to store the parameters after SQL parsing and store them by location
	parameters []string
}

func (f *fingerprintVisitor) Enter(n ast.Node) (node ast.Node, skipChildren bool) {
	if v, ok := n.(*driver.ValueExpr); ok {
		f.parameters = append(f.parameters, strings.Split(v.String(), " ")[1])
		v.Type.Charset = ""
		v.SetValue([]byte("?"))

	}
	return n, false
}

func (f *fingerprintVisitor) Leave(n ast.Node) (node ast.Node, ok bool) {
	return n, true
}

func (sp *SQLParser) FingerAndParameter(sql string) (*SQLParser, error) {
	p := parser.New()
	stmt, err := p.ParseOneStmt(sql, "", "")
	if err != nil {
		return nil, fmt.Errorf("parse sql: %w", err)
	}

	visitor := fingerprintVisitor{}
	stmt.Accept(&visitor)
	buf := new(bytes.Buffer)
	restoreCtx := format.NewRestoreCtx(format.RestoreKeyWordUppercase|format.RestoreNameBackQuotes, buf)
	err = stmt.Restore(restoreCtx)
	if err != nil {
		return nil, fmt.Errorf("restore sql: %w", err)
	}

	sp.sqlTemplate = buf.String()
	sp.bitwiseParams = visitor.parameters
	return sp, nil
}

// Result after calling Finger And Parameter, you can call this function to
// get the SQL template and parameter list.
func (sp *SQLParser) Result() (tmpl string, params []string) {
	return sp.sqlTemplate, sp.bitwiseParams
}
