package main

// RUN InitGu!

import (
	"fmt"
	"hash/fnv"
	"os"
	"io/ioutil"
	"reflect"
	"bytes"
	"os/exec"
	"math/rand"
	"time"
	"strconv"
	"strings"
	"bufio"
	"runtime"
	"errors"
	"io/fs"
	//"github.com/eiannone/keyboard"
)

func ArgvAssing( argv []string ) (map[string][]string) {
	var argc int = len(argv)
	var dict map[string][]string = map[string][]string {}
	var now string = ""
	var i int
	for i = 1; i < argc; i++ {
		if argv[i][0] == '-' {
			now = argv[i]
			dict[now] = append(dict[now], "")
		} else {
			dict[now] = append(dict[now], argv[i])
		}
	}
	return dict
}

type _s_get struct {
	Exists bool
	First string
	Last string
	List []string
}

func Get( gts string ) (_s_get) {
	_list, _exists := args[gts]
	var _ll int = len(_list)
	var _first string = ""
	var _last string = ""
	if _ll > 0 {
		_first = _list[0]
		_last = _list[_ll-1]
	}
	return _s_get{
		_exists,
		_first,
		_last,
		_list,
	}
}

func Repr( v interface{} ) (string) {
	var vtype reflect.Type = typeof(v)
	var types map[reflect.Type]string = map[reflect.Type]string {}
	types[StringType] = "S["
	types[IntType] = "I["
	types[BoolType] = "B["
	types[FloatType] = "F["
	return fmt.Sprintf("%s%v]", types[vtype], v)
}

func typeof( v interface{} ) (reflect.Type) {
	return reflect.TypeOf(v)
}

func ReadFile( filename string ) (string) {
	file, err := os.Open(filename) // For read access.
	panic(err)
	_ = err
	FILE, err := ioutil.ReadAll(file)
	panic(err)
	var FILES bytes.Buffer
	FILES.Write(FILE)

	return string(FILES.Bytes())
}

func ReadFileBytes( filename string ) ([]byte) {
	file, err := os.Open(filename) // For read access.
	panic(err)
	_ = err
	FILE, err := ioutil.ReadAll(file)
	var FILES bytes.Buffer
	FILES.Write(FILE)

	return (FILES.Bytes())
}

func WriteFile( filename string, write string) {
	err := os.WriteFile(filename, []byte(write), 0644) // 1X 2W 4R
	panic(err)
}

var GetCh func()(string)

/*
func WindowsGetCh() (string) {
	char, _, err := keyboard.GetSingleKey()
	if (err != nil) {
		panic(err)
	}
	return string(char)
}
*/

func LinuxGetCh() (string) {
	var b []byte = make([]byte, 1)
	os.Stdin.Read(b)
	return string(b)
}

func InitGetCh() {
	//TODO(1) WindowsGetCh: burh
	if runtime.GOOS=="linux" {
		exec.Command("/usr/bin/stty", "-F", "/dev/tty","-echo", "cbreak", "min", "1").Run()
		GetCh=LinuxGetCh
	} else {
		//e := keyboard.Open()
		//panic(e)
		//GetCh=WindowsGetCh
	}
}

func GetChByte ( s *bufio.Reader ) ([]byte) {
	var b []byte = make([]byte, 8)
	s.Read(b)
	return b
}

func GetChNByte ( s *bufio.Reader, n int ) ([]byte) {
	var b []byte = make([]byte, n)
	s.Read(b)
	return b
}

func GetChBA(s *bufio.Reader, b *[]byte) (error) {
	_, a := s.Read(*b)
	return a
}

func spos(y int, x int) (string) {
	return fmt.Sprintf("\x1b[%d;%dH", y+1, x+1)
}

func pos(y int, x int) {
	fprintf(stdout, "\x1b[%d;%dH", y+1, x+1)
	stdout.Flush()
}

func printat(y int, x int, prt interface{}) {
	fmt.Printf("%s%v",spos(y,x),prt)
}

func SeedRand(seed int64) {
	rand.Seed(seed)
}

func InitRand() {
	SeedRand(time.Now().UTC().UnixNano())
}

func rint(min int , max int) (int) {
	return rand.Intn(max-(min-1))+min
}

func rbool() (bool) {
	return rand.Intn(2)==1
}

func rboolin(in int) (bool) {
	return rand.Intn(in+1)==1
}

func oldinput(prompt string) (string) {
	var b = ""
	print(prompt)
	var i = make([]byte, 1)
	for{
		os.Stdin.Read(i)
		stdout.Write(i)
		stdout.Flush()
		if i[0] == 10{break}
		if i[0] == 127{
			b = b[:len(b)-1]
			stdout.WriteString("\b \b\b")
		}
		b+=string(i[0])
		i = []byte{0}
	}
	return b
}

func CursorMode(mode string) () {
	var CursorModes map[string]int = map[string]int{
		"blinking block":1,
		"block":2,
		"blinking underline":3,
		"underline":4,
		"blinking I-beam":5,
		"I-beam":6,
	}
	fmt.Printf("\x1b[%d q", CursorModes[mode])
}

func GetTerminalSize() (int, int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, r := cmd.Output()
	panic(r)
	out = out[:len(out)-1]
	var ys string
	var xs string
	var spaced bool = false
	var x int
	var y int
	for  i := 0; i < len(out); i++ {
		if out[i] == byte(32) { // space
			spaced = true
		} else if spaced {
			xs+=string(out[i])
		} else {
			ys+=string(out[i])
		}
	}
	x, _ = (strconv.Atoi(xs))
	y, _ = (strconv.Atoi(ys))
	return y, x
}

func sleep(tm float64) {
	var slp = time.Duration(1000000000.0*tm)
	time.Sleep(slp)
}

func LockLink( link chan bool ) {
	for {if <-link{break} else {sleep(0.005)}}
}

func ReverseString( str string ) (string) {
	var now = make([]string, len(str))
	var ret = make([]string, len(str))
	for i:=0 ; i < len(str) ; i ++ {
		now[i] = string(str[i])
	}
	var j int = 0
	for i:=(len(now)-1) ; i >= 0; i -- {
		ret[i] = now[j]
		j++
	}
	return strings.Join(ret, "")
}

var _clear map[string]func() //create a map for storing clear funcs

func InitGu() {

	InitRand()
	//print("\x1b[38;2;255;255;255m\n\x1b[1;1H")
	_clear = make(map[string]func()) //Initialize funcs map
	_clear["linux"] = func() {
		cmd := exec.Command("/usr/bin/clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	_clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func clear() {
	value, ok := _clear[runtime.GOOS]
	if ok {
		value()
	} else {
		printf("Your platform [%s] is unsupported! I can't clear terminal screen :(", runtime.GOOS)
		Exit(1)
	}
}

func Exit( ecode int ) {
	stdout.Flush()
	stderr.Flush()
	os.Exit(ecode)
}

func Print( thing interface{} ) {
	printf("%v\n", thing)
}

func StoIA( str string ) ([]int) {
	values := make([]int, 0, len(str))
	for _, raw := range str {
		v, err := strconv.Atoi(string(raw))
		if err != nil {
			panic(err)
			continue
		}
		values = append(values, v)
	}
	return values
}

func RGB( r, g, b interface{} ) (string) {
	return fmt.Sprintf("\x1b[38;2;%v;%v;%vm", r, g, b )
}

var COLOR = map[string]string{
	"nc" : RGB(0xff, 0xff, 0xff),
	"red" : RGB(0xff, 0x0, 0x0),
	"green" : RGB(0x0, 0xff, 0x0),
	"blue" : RGB(0x0, 0x0, 0xff),
	"cyan" : RGB(0x80, 0x80, 0xff),
	"yellow" : RGB(0xff, 0xff, 0x0),
}

func bog(ifer bool, f1, f2 interface{}) (interface{}) {
	if (ifer) {
		return f1
	} else {
		return f2
	}
}

func IinA(a interface{}, arr []interface{}) (bool) {
	for _, b := range arr {
		if b == a {
			return true
		}
	}
	return false
}

func MakeArray(size int, value interface{}) ([]interface{}) {
	var array = make([]interface{}, size)
	for i := range array {
		array[i] = value
	}
	return array
}

func exists( filename string ) (bool) {
	_, err := os.Stat(filename)
	return !errors.Is(err, os.ErrNotExist)
}

func dprint(stream *bufio.Writer, Type string, Text string, Info ...interface{}) {
	fprintf(stream, "[%s]: %s\n", Type, spf(Text, Info...))
	stream.Flush()
}

func panic( err error ) {
	if ( err != nil ) {
		dprint(stderr, "ERROR", "%v\n", err)
		Exit(1)
	}
}

func fread( file *FILE, blen int ) ([]byte, int) {
	var err error
	var brd int
	var buff = make([]byte, blen)
	brd, err = file.Read(buff)
	panic(err)
	return buff, brd
}

func SetBit(flag, S byte	) (byte) {
	return S|flag
}

func UnsetBit(flag, S byte	) (byte) {
	return S&^flag
}

func IsBitSet(flag, S byte	) (bool) {
	return S&flag != 0
}

func ToggleBit(flag, S byte ) (byte) {
	return S^flag
}

func DecompressByte( b byte ) ([8]bool) {
	return [8]bool{
		b>>7& 1 == 1,
		b>>6& 1 == 1,
		b>>5& 1 == 1,
		b>>4& 1 == 1,
		b>>3& 1 == 1,
		b>>2& 1 == 1,
		b>>1& 1 == 1,
		b	& 1 == 1,
	}
}

func Pop(xs []interface{}, i int) (interface{}, []interface{}) {
	y := xs[i]
	ys := append(xs[:i], xs[i+1:]...)
	return y, ys
}

func Input(prompt string) (string) {
	var ipt string
	var err error
	print(prompt)
	ipt, err = stdin.ReadString('\n')
	panic(err)
	return ipt[:len(ipt)-1]
}

func HashStr(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func HashInt(i int) uint32 {
	h := fnv.New32a()
	h.Write([]byte(sprintf("%d", i)))
	return h.Sum32()
}

func PS( thing ...interface{} ) { // print simple
	printf("%v\n", thing)
}

func SP( thing ...interface{} ) (string) { // print simple
	return sprintf("%v\n", thing)
}

func HideCursor() () {
	fprintf(stdout, "\x1b[?25l")
	stdout.Flush()
}

func ShowCursor() () {
	fprintf(stdout, "\x1b[?25h")
	stdout.Flush()
}

func _ls(dr string) ([]fs.FileInfo) {
	dir, err := ioutil.ReadDir(dr)
	panic(err)
	return dir
}

func LsSize( dr string ) ( []int64 ) {
	var dir []fs.FileInfo = _ls(dr)
	var buff = make([]int64, len(dir))
	for i:=0;i<len(dir);i++ {
		buff[i] = dir[i].Size()
	}
	return buff
}

func LsSumSize( dr string ) (int64) {
	var dir []fs.FileInfo = _ls(dr)
	var buff int64 = 0
	for i:=0;i<len(dir);i++ {
		if (dir[i].IsDir()){
			buff+=LsSumSize(dr+(dir[i].Name())+"/")
		} else {
			buff += dir[i].Size()
		}
	}
	return buff
}

var _LsSumSizef_Postfix_to_char = []string{
	0:"b",
	1:"kb",
	2:"mb",
	3:"gb",
	4:"tb",
	5:"pt",
}

// this shit aint working
// dunno why
// good night
func _LsSumSizef( dr string , threshold float64 ) ( string ) {
	var size float64 = float64(LsSumSize(dr))
	var postfix int = 0
	var sz string
	//0-b 1-kb 2-mb 3-gb 4-tb
	for ;(size/1024)>threshold;{
		postfix++
		size = size/1024
	}
	sz = spf("%.2f%s", size, _LsSumSizef_Postfix_to_char[postfix])
	return sz
}

func LsSumSizef( dr string ) ( string ) {
	return _LsSumSizef(dr, 10)
}

func Ls(dr string) ([]string) {
	var dir []fs.FileInfo = _ls(dr)
	// make the array, append is cringe
	var buff = make([]string, len(dir))
	for i:=0;i<len(dir);i++ {
		buff[i] = dir[i].Name()
		if (dir[i].IsDir()){
			buff[i]+="/"
		}
	}
	return buff
}

func Color(fr,fg,fb, br,bg,bb interface{}) (string) {
	return spf("\x1b[38;2;%v;%v;%v;48;2;%v;%v;%vm", fr,fg,fb, br,bg,bb)
}

func Bkcolor(br,bg,bb interface{}) (string) {
	return spf("\x1b[48;2;%v;%v;%vm", br,bg,bb)
}

func Die(message string) () {
	panic(errors.New(message))
}

func CDie(message string) () {
	clear()
	stdout.Flush()
	stderr.Flush()
	pos(0,0)
	panic(errors.New(message))
}

func Assert(thing bool, message string) {
	if !thing { Die(message) }
}

func GetInt(prompt string) (int) {
	in := oldinput(prompt)
	i, err := strconv.Atoi(in)
	for ;err != nil; {
		in = oldinput(prompt)
		i, err = strconv.Atoi(in)
	}
	return i
}

type Log struct {
	filename string
	fd *FILE
	out Writer
	lenght int
	autosave bool
}

func MakeLog(filename string) (Log) {
	var f, err = fmake(filename)
	panic(err)
	return Log{filename, f, bufio.NewWriter(f), 0, true}
}

func (l Log) PS ( thing ...interface{} ) {
	fprintf(l.out, "%v\n", thing)
	if l.autosave {
		l.Save()
	}
}

func (l Log) write ( thing string ) {
	fprintf(l.out, thing)
	if l.autosave {
		l.Save()
	}
}

func (l Log) Save () {
	l.out.Flush()
}

func (l Log) End () {
	l.out.Flush()
	l.fd.Close()
}

//dodef
var (
	stdout *bufio.Writer = bufio.NewWriter(os.Stdout)
	stderr *bufio.Writer = bufio.NewWriter(os.Stderr)
	stdin  *bufio.Reader = bufio.NewReader(os.Stdin )
	args map[string][]string = ArgvAssing(os.Args)
	argv = os.Args[1:]
	argc = len(argv)
	format = fmt.Sprintf
	printf = fmt.Printf
	sprintf = fmt.Sprintf
	spf = fmt.Sprintf
	fprintf = fmt.Fprintf
	fmake = os.Create
	StringType reflect.Type = typeof("")
	IntType reflect.Type = typeof(2)
	BoolType reflect.Type = typeof(true)
	FloatType reflect.Type = typeof(0.1)
);

//const def
const (
	// max-mins
	// 8 bytes
	I8MAX =  int8(0x7f)
	I8MIN =  int8(-0x80)
	U8MAX =  uint8(0xFF)
	// 16 bytes
	I16MAX = int16(0x7FFF)
	I16MIN = int16(-0x8000)
	U16MAX = uint16(0xFFFF)
	// 32 bytes
	I32MAX = int32(0x7FFFFFFF)
	I32MIN = int32(-0x80000000)
	U32MAX = uint32(0xFFFFFFFF)
	// 64 bytes
	I64MAX = int64(0x7FFFFFFFFFFFFFFF)
	I64MIN = int64(-0x8000000000000000)
	U64MAX = uint64(0xFFFFFFFFFFFFFFFF)
	// file nums
	F_append = os.O_APPEND
	F_WR = os.O_WRONLY
)

//typedef
type FILE = os.File
type Reader = *bufio.Reader
type Writer = *bufio.Writer
