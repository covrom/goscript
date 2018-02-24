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
	st := struct {
		A   int
		B   string
		D   string
		Flt float64 `goscript:"f"`
	}{}
	if err := RunAST(astf, &st); err != nil {
		t.Error(err)
		return
	}
	t.Logf("Parsed struct %#v", st)

}
