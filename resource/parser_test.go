package resource

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestQuote(t *testing.T) {
	raw := Escape(`The portable canister is the Stationeer's basic unit of 
<color=#0080FFFF>{LINK:GasPage;gas}</color> delivery. Rated to a standard pressure o
f 8000kPa (80 atmospheres), empty gas canisters can be mounted to a 
{thing:DynamicGasCanisterAir} or {thing:StructureGasTankStorage} for refill. Careful n
ot to pressurize beyond 100kPA, or it may go 'bang'. Contains 64L of gas.`)
	if !strings.ContainsRune(raw, ' ') {
		t.Error("Should leave NBSP")
	}
}

func TestMerge(t *testing.T) {
	r := &parser{
		orig:   read(t, "Colors.po"),
		merged: read(t, "Colors.po"),
	}
	r.populate()
	if len(r.origID) == 0 {
		t.Error("not populated")
	}
	for k, v := range r.origID {
		fmt.Println(k)
		if len(v) == 0 {
			t.Error("nothing parsed")
		}
		for _, vi := range v {
			fmt.Println(vi)
		}
	}
	r.merge()
	if !bytes.Equal(r.result.Bytes(), r.orig) {
		t.Error("should not change")
	}
	f, c := create(t, "Colors.merged.po")
	io.Copy(f, r.result)
	c()
}

func TestParser(t *testing.T) {
	orig := read(t, "merge_original.po")
	r := &parser{
		orig:   orig,
		merged: read(t, "merge_result.po"),
	}
	r.populate()
	if len(r.origID) == 0 {
		t.Error("not populated")
	}
	for k, v := range r.origID {
		fmt.Println(k)
		if len(v) == 0 {
			t.Error("nothing parsed")
		}
		for _, vi := range v {
			fmt.Println(vi)
		}
	}
	r.merge()
	f, c := create(t, "merge_merged.po")
	io.Copy(f, r.result)
	c()
}
