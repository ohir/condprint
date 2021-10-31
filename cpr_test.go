// (c) 2021 Ohir Ripe. MIT license.

package condprint

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"testing"
)

const sumOK = "9f1b0a88203d6dd938f2b4854341f00bb74f0dfa1e5ca11cbecd63dc7b89d730"

var (
	logsink strings.Builder
	pffn    int
)

func TestBadArg(t *testing.T) {
	tre := func() (s string) {
		defer func() {
			e := recover()
			if e != nil {
				s = "OK"
			}
		}()
		_ = Verbose(true, 121)
		return "Bad"
	}
	if tre() != "OK" {
		t.Logf("Verbose should panic but it did not!")
		t.Fail()
	}
}

func TestPrint(t *testing.T) {
	pff := func() string { pffn++; return fmt.Sprintf(" pfxf called %d times! ", pffn) }
	pnff := func() string { pffn++; return fmt.Sprintf("\nPfxf called %d times! ", pffn) }
	p := Verbose(false)
	p("This should not print!")
	p = Verbose(true)
	p.If(false, "This neither!")
	p = Verbose(false, &logsink)
	p("This neither!")
	p = Verbose(true, &logsink)
	p("Test Begins with Printf, then ")
	p.If(true, "we use p.If, then ")
	p.If(false, " it should NOT print now (false) ")
	p.IfNot(true, " neither p.IfNot may print now (for true)")
	p.IfNot(false, "p.IfNot should have a say.\n")
	p("------------")
	p = Verbose(true, &logsink, "\n[PfxStr] ")
	p("Said with prefix PfxStr...\n")
	p = Verbose(true, &logsink, pnff)
	p("Said with prefix function...\n")
	p = Verbose(true, &logsink, pnff, " [PfxStr] ")
	p("Said with prefix function then PfxStr...\n")
	p = Verbose(true, &logsink, "\n[PfxStr] ", pff)
	p("Said with PfxStr first, then function...\n")
}

func TestRegression(t *testing.T) {
	var skipck bool
	if outfn := os.Getenv("MKGOLD"); outfn != "" {
		if outfn == "NOHASH" {
			skipck = true
		} else if outfn == "T" {
			os.Stdout.WriteString(logsink.String())
		} else if err := os.WriteFile(outfn, []byte(logsink.String()), 0660); err != nil {
			t.Fatalf("Can not dump to file %s [%v]", outfn, err)
		} else {
			fmt.Fprintf(os.Stderr, "\"Golden\" output has been written to %s\n", outfn)
		}
	}
	if !skipck {
		out := fmt.Sprintf("%x", sha256.Sum256([]byte(logsink.String())))
		if out != sumOK {
			t.Logf("Regression! Expected hash: %s", sumOK)
			t.Logf("        Hash now computed: %s", out)
			t.Logf("Diff MKGOLD output before setting a new hash to:")
			t.Logf("const sumOK = \"%s\"", out)
			t.Fail()
		}
	}
}
