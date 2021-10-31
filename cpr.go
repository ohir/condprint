// (c) 2021 Ohir Ripe. MIT license.

/* Package condprint provides a fmt.Printf type function writing to
either os.Stdout or to any other io.StringWriter passed to the guard
function Verbose.

Function Verbose takes a boolean flag V and up to three optional parameters:
 - prefix string, to be prepended as-is to the output
 - prefix `func() string` returning string to be prepended to the output
 - target io.StringWriter that printers will write to (if V is true).
   Both prefix string and prefix function can be given, and their precedence
   on input parameters list will be preserved in the output.

If V is false at Verbose call site, the noop empty function is returned.
It is possible to simultanously have many printers writing to distinct
destinations (eg. *strings.Builder type) with distinct prefixes.

Returned printf function itself has two methods: If and IfNot, that make it
print according to condition in place (and as method names tell). Eg.

  p := condprint.Verbose(loglevel > 3, &logsink, tmstamp) // function tmstamp
  p("Writes always")
  p.If(details, "Only if condition is true")
  p.IfNot(details, "Only if condition is false")

*/

package condprint

import (
	"fmt"
	"io"
	"os"
)

type cprn func(string, ...interface{})

/*

function Verbose(V bool, ...opts) (printf) returns a printf type function
(for V being true), or an empty stub (for V being false).

Verbose may take three optional parameters:
 - prefix string, to be prepended as-is to the output
 - prefix function returning string that gets prepended to the output
 - target io.StringWriter printers will write to, if V is true.
 Both prefix string and prefix function can be given, and its precedence
 on input parameters list will be preserved in the output.

 Note: to early indicate options misuse, Verbose will panic if given
 optional parameter of a wrong type.
*/
func Verbose(V bool, opts ...interface{}) cprn { // Printf
	if !V {
		return func(string, ...interface{}) {}
	}
	var pfxf func() string  // prefix function, eg timenow()
	var pfx string          // prefix string
	var after bool          // pfx prints after a call to pfxf
	var out io.StringWriter // write there
	for _, opt := range opts {
		switch v := opt.(type) {
		case string:
			pfx = v
			if pfxf != nil {
				after = true
			}
		case io.StringWriter:
			out = v
		case func() string:
			pfxf = v
		default:
			panic(fmt.Sprintf("condprint.Verbose got option of a wrong type %T!\n", v))
		}
	}
	if out == nil {
		out = os.Stdout
	}
	if len(pfx) == 0 && pfxf == nil {
		return func(Fmt string, a ...interface{}) { // Printf
			out.WriteString(fmt.Sprintf(Fmt, a...))
		}
	} else {
		return func(Fmt string, a ...interface{}) {
			if !after && len(pfx) > 0 {
				out.WriteString(pfx)
			}
			if pfxf != nil {
				out.WriteString(pfxf())
			}
			if after && len(pfx) > 0 {
				out.WriteString(pfx)
			}
			out.WriteString(fmt.Sprintf(Fmt, a...))
		}
	}
}

// Method If writes on condition c being true.
func (p cprn) If(c bool, Fmt string, a ...interface{}) bool {
	if c {
		p(Fmt, a...)
	}
	return c
}

// Method IfNot writes on condition c being false.
func (p cprn) IfNot(c bool, Fmt string, a ...interface{}) bool {
	if !c {
		p(Fmt, a...)
	}
	return c
}
