

# condprint
`import "github.com/ohir/condprint"`

Package condprint provides a `fmt.Printf` type function writing to either `os.Stdout` or to any other `io.StringWriter`. This function is instantiated via the `Verbose` guard function taking a boolean flag V and up to three optional parameters:
 - target `io.StringWriter`. Returned printf will write to it. Usually an address of local `strings.Builder` is given as a target, otherwise the default `os.Stdout` is used.
 - prefix string, to be prepended as-is to the output
 - prefix `func() string` returning string to be prepended to the output
   Both prefix string and prefix function can be given, and their precedence on input parameters list will be preserved in the output.

If V is false at Verbose call site, the noop empty function is returned instead.

The returned from Verbose printf function has two additional methods: `If` and `IfNot`, that print according to condition in place (and as method names tell). Eg.
``` go
  p := condprint.Verbose(loglevel > 3, &logsink, tmstamp) // function tmstamp
  p("Writes always")
  p.If(details, "Only if condition is true")
  p.IfNot(details, "Only if condition is false")
```
Both `If` and `IfNot` return condition intact.

### <a name="Verbose">func</a> [Verbose](/cpr.go?s=1862:1908#L53)
``` go
func Verbose(V bool, opts ...interface{}) func(fmt string, a ...interface{})
```
function Verbose(V bool, ...opts) (printf) returns a printf type function
(for V being true), or an empty stub (for V being false).

Verbose may take three optional parameters:


	- prefix string, to be prepended as-is to the output
	- prefix function returning string that gets prepended to the output
	- target io.StringWriter printers will write to, if V is true.
	Both prefix string and prefix function can be given, and its precedence
	on input parameters list will be preserved in the output.

__Note__: to indicate options misuse early, `Verbose` will panic if given optional parameter of a wrong type.
