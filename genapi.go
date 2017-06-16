// +build ignore

package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"
)

var (
	gopackage = os.Getenv("GOPACKAGE")
	gofile    = os.Getenv("GOFILE")

	module  string
	hmodule syscall.Handle

	checkfunc = false
	buildsrc  = false
)

var (
	syscallfuncs = [...]string{
		"syscall.Syscall",
		"syscall.Syscall6",
		"syscall.Syscall9",
		"syscall.Syscall12",
		"syscall.Syscall15",
	}
)

type argdecl struct {
	name string
	typ  string
}

func (a *argdecl) callarg() string {
	if strings.HasPrefix(a.typ, "*") || strings.HasSuffix(a.typ, "Func") {
		return "uintptr(unsafe.Pointer(" + a.name + "))"
	} else if a.typ == "bool" {
		return "boolcast(" + a.name + ")"
	} else if a.typ != "uintptr" {
		return "uintptr(" + a.name + ")"
	}
	return a.name
}

type funcdecl struct {
	decl  string
	name  string
	entry string
	pp    string
	args  []*argdecl
	rets  []string

	ptrname string
	reterr  bool
}

func (fn *funcdecl) ptrvar() string {
	if fn.ptrname == "" {
		fn.ptrname = "pfn" + fn.name
	}
	return fn.ptrname
}

func (fn *funcdecl) returnReceiver() string {
	receivers := ""
	switch len(fn.rets) {
	case 0:
	case 1:
		if fn.rets[0] == "error" {
			receivers = "_, _, en"
			fn.reterr = true
		} else {
			receivers = "r1, _, _"
		}
	case 2:
		receivers = "r1, _, en"
		fn.reterr = true
	case 3:
		receivers = "r1, r2, en"
		fn.reterr = true
	default:
		fatal("too many receivers")
	}
	return receivers
}

func (fn *funcdecl) returnValues() string {
	values := ""
	i := 1
	for _, ret := range fn.rets {
		if values != "" {
			values += ", "
		}
		if ret == "error" {
			values += "err"
		} else {
			if ret == "uintptr" {
				values += fmt.Sprintf("r%d", i)
			} else {
				values += fmt.Sprintf("%s(r%d)", ret, i)
			}
			i++
		}
	}
	return values
}

func (fn *funcdecl) callFunc() string {
	if len(fn.args) == 0 {
		return syscallfuncs[0]
	} else if len(fn.args) <= 15 {
		return syscallfuncs[(len(fn.args)-1)/3]
	} else {
		fatal("too many arguments")
	}
	return ""
}

func (fn *funcdecl) callRemaining() string {
	remaining := ""
	require := 0
	if len(fn.args) == 0 {
		require = 3
	} else if len(fn.args) <= 15 {
		require = ((len(fn.args)-1)/3 + 1) * 3
	}
	for i := len(fn.args); i < require; i++ {
		if remaining != "" {
			remaining += ", "
		}
		remaining += "0"
	}
	return remaining
}

var (
	hasfunc  = map[string]bool{}
	funclist = []*funcdecl{}
)

func fatal(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
}

func parserets(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	if strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")") {
		var rets []string
		retlist := strings.Split(s[1:len(s)-1], ",")
		for _, ret := range retlist {
			rets = append(rets, strings.TrimSpace(ret))
		}
		return rets
	}
	return []string{s}
}

func parseargs(s string) (string, []*argdecl) {
	var (
		preprocess string
		args       []*argdecl
	)
	argslist := strings.Split(s, ",")
	for i, arg := range argslist {
		arg = strings.TrimSpace(arg)
		if arg == "" {
			continue
		}
		sppos := strings.IndexByte(arg, ' ')
		name := arg[:sppos]
		typ := strings.TrimSpace(arg[sppos+1:])
		if strings.HasPrefix(typ, "[]") {
			pname := fmt.Sprintf("p%s%d", name, i)
			ptype := "*" + typ[2:]
			nname := fmt.Sprintf("n%s%d", name, i)
			preprocess += fmt.Sprintf("var (\n%s %s\n%s = len(%s)\n)\n", pname, ptype, nname, name)
			preprocess += fmt.Sprintf("if %s > 0 {\n%s = &%s[0]\n}\n", nname, pname, name)
			args = append(args, &argdecl{
				name: pname,
				typ:  ptype,
			})
			args = append(args, &argdecl{
				name: nname,
				typ:  "int",
			})
		} else {
			args = append(args, &argdecl{
				name: name,
				typ:  typ,
			})
		}
	}
	return preprocess, args
}

func parsefunc(line string) *funcdecl {
	entry := ""
	if strings.HasPrefix(line, "entry=") {
		sppos := strings.IndexByte(line, ' ')
		entry = line[6:sppos]
		line = strings.TrimSpace(line[sppos:])
	}
	lbpos := strings.IndexByte(line, '(')
	rbpos := strings.IndexByte(line[lbpos:], ')') + lbpos
	pp, args := parseargs(line[lbpos+1 : rbpos])
	name := strings.TrimSpace(line[5:lbpos])
	if entry == "" {
		entry = name
	}
	return &funcdecl{
		decl:  line,
		name:  name,
		entry: entry,
		pp:    pp,
		args:  args,
		rets:  parserets(line[rbpos+1:]),
	}
}

var (
	pfngetprocaddress uintptr
)

func getprocaddr(hmodule syscall.Handle, procname string) (uintptr, syscall.Errno) {
	ptr := uintptr(0)
	if procname[0] == '#' {
		for i := 1; i < len(procname); i++ {
			c := procname[i]
			if c < '0' || c > '9' {
				break
			}
			ptr = ptr*10 + uintptr(c-'0')
		}
	} else {
		ptr = *(*uintptr)(unsafe.Pointer(&procname))
	}
	proc, _, err := syscall.Syscall(pfngetprocaddress, 2,
		uintptr(hmodule),
		ptr,
		0)
	if err != 0 {
		return 0, err
	}
	return proc, 0
}

func parsecomment(text string) {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		fn := parsefunc(line)
		if fn == nil {
			fatal("parse func failed")
		}
		if hasfunc[fn.name] {
			continue
		}
		hasfunc[fn.name] = true

		if checkfunc {
			_, err := getprocaddr(hmodule, fn.entry+"\000")
			if err != 0 {
				fmt.Println("genapi:", "[checkfunc]", fn.name, err)
				continue
			}
			if fn.name[0] < 'A' || fn.name[0] > 'Z' {
				fmt.Println("genapi:", "[checkfunc]", fn.name, "unexported name")
				continue
			}
		}

		funclist = append(funclist, fn)
	}
}

func render(w io.Writer) error {
	bw := bufio.NewWriter(w)

	// file header
	// Code generated by protoc-gen-go.
	// source: LinkerProtocol.proto
	bw.WriteString("// generated by genapi.go\n")
	bw.WriteString("// GOFILE=" + gofile + " GOPACKAGE=" + gopackage + "\n")
	bw.WriteString("// DO NOT EDIT!\n")
	bw.WriteString("package " + gopackage + "\n")
	bw.WriteString("\n")
	bw.WriteString("import (\n")
	fmt.Fprintf(bw, "\"%s\"\n", "syscall")
	fmt.Fprintf(bw, "\"%s\"\n", "unsafe")
	bw.WriteString(")\n")
	bw.WriteString("\n")
	bw.WriteString("var _ unsafe.Pointer // keep unsafe\n")
	bw.WriteString("\n")

	// func pointers
	bw.WriteString("var (\n")
	for _, fn := range funclist {
		fmt.Fprintf(bw, "%s uintptr\n", fn.ptrvar())
	}
	bw.WriteString(")\n")
	bw.WriteString("\n")

	// utils
	bw.WriteString(`func mustload(libname string) syscall.Handle {
	hlib, err := syscall.LoadLibrary(libname)
	if err != nil {
		panic(err)
	}
	return hlib
}

var (
	pfngetprocaddress uintptr
)

func mustfind(hmodule syscall.Handle, procname string) uintptr {
	ptr := uintptr(0)
	if procname[0] == '#' {
		for i := 1; i < len(procname); i++ {
			c := procname[i]
			if c < '0' || c > '9' {
				break
			}
			ptr = ptr*10 + uintptr(c-'0')
		}
	} else {
		ptr = *(*uintptr)(unsafe.Pointer(&procname))
	}
	proc, _, err := syscall.Syscall(pfngetprocaddress, 2,
		uintptr(hmodule),
		ptr,
		0)
	if proc == 0 {
		panic(err)
	}
	return proc
}

func boolcast(b bool) uintptr {
	if b {
		return 1
	}
	return 0
}
`)
	bw.WriteString("\n")

	// exported functions
	for _, fn := range funclist {
		fmt.Fprintf(bw, "%s {\n", fn.decl)
		if fn.pp != "" {
			bw.WriteString(fn.pp)
		}
		recv := fn.returnReceiver()
		if recv != "" {
			recv += ":="
		}
		fmt.Fprintf(bw, "%s%s(%s, %d,\n", recv, fn.callFunc(), fn.ptrvar(), len(fn.args))
		for _, arg := range fn.args {
			fmt.Fprintf(bw, "%s,\n", arg.callarg())
		}
		fmt.Fprintf(bw, "%s)\n", fn.callRemaining())
		if len(fn.rets) > 0 {
			if fn.reterr {
				bw.WriteString(`var err error
if en != 0 {
	err = en
}
`)
			}
			fmt.Fprintf(bw, "return %s", fn.returnValues())
		}
		fmt.Fprintf(bw, "}\n")
		bw.WriteString("\n")
	}

	// init function
	bw.WriteString("func init() {\n")
	bw.WriteString("hkernel32 := mustload(\"kernel32.dll\")\n")
	bw.WriteString(`var err error
pfngetprocaddress, err = syscall.GetProcAddress(hkernel32, "GetProcAddress")
if err != nil {
	panic(err)
}
`)
	if gopackage != "kernel32" {
		fmt.Fprintf(bw, "h%s := mustload(\"%s\")\n", gopackage, module)
		fmt.Fprintf(bw, "_ = h%s\n", gopackage)
	}
	for i := range funclist {
		fn := funclist[i]
		fmt.Fprintf(bw, "%s = mustfind(h%s, \"%s\\000\")\n", fn.ptrvar(), gopackage, fn.entry)
	}
	bw.WriteString("}\n")
	return bw.Flush()
}

func mustload(libname string) syscall.Handle {
	hlib, err := syscall.LoadLibrary(libname)
	if err != nil {
		panic(err)
	}
	return hlib
}

func main() {
	if gopackage == "" || gofile == "" {
		fatal("missing required arguments")
	}

	flag.BoolVar(&buildsrc, "buildsrc", false, "build after gen")
	flag.BoolVar(&checkfunc, "checkfunc", false, "check on parsing")
	flag.Parse()

	module = gopackage + ".dll"
	if checkfunc {
		hkernel32 := mustload("kernel32.dll")
		pfngetprocaddress, _ = syscall.GetProcAddress(hkernel32, "GetProcAddress")
		hmodule = mustload(gopackage + ".dll")
	}

	fmt.Println("genapi:", "parsing comments from", gofile)
	fset := token.NewFileSet()
	astf, err := parser.ParseFile(fset, gofile, nil, parser.ParseComments)
	if err != nil {
		fatal(err)
	}
	for _, g := range astf.Comments {
		for _, c := range g.List {
			s := c.Text
			if strings.HasPrefix(s, "/*") && strings.HasSuffix(s, "*/") {
				parsecomment(s[2 : len(s)-2])
			}
		}
	}

	outputfile := gofile
	if strings.HasSuffix(gofile, ".go") {
		outputfile = gofile[:len(gofile)-3]
	}
	outputfile += ".api.go"

	fmt.Println("genapi:", "generating", outputfile)
	f, err := os.Create(outputfile)
	if err != nil {
		fatal("fatal:", "os.Create", err)
	}
	err = render(f)
	if err != nil {
		fatal("fatal:", "render", err)
	}
	f.Close()

	fmt.Println("genapi:", "formatting", outputfile)
	err = exec.Command("go", "fmt", outputfile).Run()
	if err != nil {
		fatal("fatal:", "go fmt", err)
	}

	if buildsrc {
		wd, _ := os.Getwd()
		fmt.Println("genapi:", "building", wd)
		err = exec.Command("go", "build").Run()
		if err != nil {
			fatal("fatal:", "go build", err)
		}
	}

	fmt.Println("genapi:", outputfile, "done")
}
