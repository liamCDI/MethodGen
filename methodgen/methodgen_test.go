package main

import (
	"reflect"
	"testing"
)

func TestProcTag(t *testing.T) {
	tst := "`mg:\"omit,len=10\"` `other:dontcare`"
	exp := make(map[string]string)
	exp["omit"] = ""
	exp["len"] = "10"

	got := ProcTag(tst)

	if reflect.DeepEqual(exp, got) == false {
		t.Errorf("expected %+v but got %+v", exp, got)
	}

}
