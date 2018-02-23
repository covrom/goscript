package parser

import (
	"testing"
)

func TestParseFile(t *testing.T) {
	astf, err := ParseFile("./testscript")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v\n", astf.Decls)
}

func TestRunAST(t *testing.T) {
	astf, err := ParseFile("./testscript")
	if err != nil {
		t.Error(err)
		return
	}
	if err := RunAST(astf, &struct {
		a int
	}{}); err != nil {
		t.Error(err)
		return
	}

}
