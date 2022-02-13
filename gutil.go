package gutil

// RUN InitGu!

import (
	"fmt"
	"hash/fnv"
	"os"
	"log"
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
)

func argvAssing( argv[] string ) (map[string][]string) {
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

func get( gts string ) (_s_get) {
	_list, _exists := args[gts]
	var _ll int = len(_list)
	var _first string = ""
	var _last string = ""
	_ = _last
	_ = _first
	if _ll > 1 {
		var _first string = _list[0]
		var _last string = _list[_ll-1]
		_ = _last
		_ = _first
	}
	return _s_get{
		Exists: _exists,
		First: _first,
		Last: _last,
		List: _list,
	}
}

var stringType reflect.Type = typeof("")
var intType reflect.Type = typeof(2)
var boolType reflect.Type = typeof(true)
var floatType reflect.Type = typeof(0.1)
func Repr( v interface{} ) (string) {
	var vtype reflect.Type = typeof(v)
	var types map[reflect.Type]string = map[reflect.Type]string {}
	types[stringType] = "S["
	types[intType] = "I["
	types[boolType] = "B["
	types[floatType] = "F["
	return fmt.Sprintf("%s%v%s", types[vtype], v, "]")
}

func typeof( v interface{} ) (reflect.Type) {
	return reflect.TypeOf(v)
}

func ReadFile( filename string ) (string) {
	file, err := os.Open(filename) // For read access.
	if err != nil {
		log.Fatal(err)
	}
	_ = err
	FILE, err := ioutil.ReadAll(file)
	var FILES bytes.Buffer
	FILES.Write(FILE)

	return string(FILES.Bytes())
}

func ReadFileBytes( filename string ) ([]byte) {
	file, err := os.Open(filename) // For read access.
	if err != nil {
		log.Fatal(err)
	}
	_ = err
	FILE, err := ioutil.ReadAll(file)
	var FILES bytes.Buffer
	FILES.Write(FILE)

	return (FILES.Bytes())
}

func WriteFile( filename string, write string) {
	err := os.WriteFile(filename, []byte(write), 0644) // 1X 2W 4R
	if err != nil{
		log.Fatal(err)
	}
}

func InitGetCh() {
	exec.Command("stty", "-F", "/dev/tty","-echo", "cbreak", "min", "1").Run()
}

func GetCh() (string) {
	var b []byte = make([]byte, 1)
	os.Stdin.Read(b)
	return string(b)
}

func GetChByte() ([]byte) { // 6 bytes
	var b []byte = make([]byte, 6)
	os.Stdin.Read(b)
	return b
}

func GetChBA(b *[]byte) ([]byte) {
	os.Stdin.Read(*b)
	return *b
}

func spos(y int, x int) (string) {
	return fmt.Sprintf("\x1b[%d;%dH", y+1, x+1)
}

func pos(y int, x int) {
	fmt.Printf("\x1b[%d;%dH", y+1, x+1)
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

func oldinput() (string) {
	var b = ""
	var i = make([]byte, 1)
	for{
		os.Stdin.Read(i)
		i[0]++
		if i[0] == 11{break}
		b+=string(i[0]-1)
		i = []byte{0}
	}
	return b
}

func hideCursor() {
	fmt.Print("\x1b[?25l")
}

func showCursor() {
	fmt.Print("\x1b[?25h")
}

func cursorMode(mode string) {
	var CursorModes map[string]int = map[string]int{
		"blinking block":1,
		"block":2,
		"blinking underline":3,
		"underline":4,
		"blinking I-beam":5,
		"I-beam":6,
	}
	fmt.Print("\033[%dq", CursorModes[mode])
}

func getTerminalSize() (int, int) {
	cmd := exec.Command("stty", "size")
	//cmd.Stdin = os.Stdin
	out, _ := cmd.Output()
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
	return join(ret, "")
}

var _clear map[string]func() //create a map for storing clear funcs

func InitGu() {
	InitRand()
	print("\x1b[38;2;255;255;255m\n\x1b[1;1H")
	_clear = make(map[string]func()) //Initialize funcs map
	_clear["linux"] = func() {
		cmd := exec.Command("clear")
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
		//exit(1)
	}
}

func exit( ecode int ) {
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
			log.Print(err)
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

func panic( err error ) {
	if ( err != nil ) {
		dprint(stderr, "ERROR", "%v\n", err)
		exit(1)
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

func pop(xs []interface{}, i int) (interface{}, []interface{}) {
	y := xs[i]
	ys := append(xs[:i], xs[i+1:]...)
	return y, ys
}

func input(text string) (string) {
	var ipt string
	var err error
	print(text)
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

type HashMap struct {
	length int
	items []interface{}
	hashes []int // before %len
}

func Hash( thing interface{} ) (int) {
	var st string = sprintf("%v", thing)
	for ;len(st) < 3;{
		st+="0"
	}
	st+="\n"
	return int(HashStr(st))
}

func MakeHashMap(keys []interface{}, results[]interface{}) (*HashMap) {
	var kl = len(keys)
	if ((kl) != len(results)) {
		dprint(stderr, "ERROR", "[MakeHashMap] keys' len %d != results' len %d\n", kl, len(results))
		exit(1)
	}
	var length = kl
	_ = kl

	// make Hash Array
	var hashes = make([]int, length)
	for i:=0;i<length;i++ {
		hashes[i] = Hash(keys[i])
	}

	// make Ordered Items
	var oi = make([]interface{}, length)
	for i:=0;i<length;i++ {
		oi[hashes[i]%length] = results[i]
	}

	//var hs HashMap = HashMap{length, oi, hashes}
	return &HashMap{length, oi, hashes}
}

func dprint(stream *bufio.Writer, Type string, Text string, Info ...interface{}) {
	if Type == "ERROR" {
		Type = COLOR["red"]+"ERROR"+COLOR["nc"]
	} // else if ... for other colors
	fprintf(stream, fs("[%s]: %s", Type, fs(Text, Info...)))
	stream.Flush()
}

// HS Any -> Any
func HSGet( h *HashMap, key interface{} ) (interface{}) {
	return h.items[Hash(key)%h.length]
}

// hash map (unsafe) add
func HSUAdd( h *HashMap, key interface{}, result interface{}) {
	h.length++
	h.hashes = append(h.hashes, (Hash(key)))
	h.items = append(h.items, result)
	//return h
}

// hash map add
func HSAdd( h *HashMap, key interface{}, result interface{}) {
	var nh = Hash(key)
	var nhl = nh%(h.length+1)
	for i:=0;i<h.length;i++ {
		if (h.hashes[i]%(h.length+1) == nhl) {
			// TODO: stop this shit
			dprint(stderr, "ERROR", "hash cruching: [%v->]%d/%d(=%d) == [%v->]%d/%d(=%d)!\n",
			h.items[i],
			// old hash
			h.hashes[i], h.length+1,	h.hashes[i]%(h.length+1),
			// added hash
			key,
			nh,			 h.length+1,	nhl)
		}
	}
	h.length++
	h.hashes = append(h.hashes, nh)
	h.items = append(h.items, result)
}

func PS( thing interface{} ) { // print single
	printf("%v\n", thing)
}

func HideCursor() () {
	fprintf(stdout, "\x1b[?25l")
	stdout.Flush()
}

func ShowCursor() () {
	fprintf(stdout, "\x1b[?25h")
	stdout.Flush()
}

//dodef
var (
	stdout *bufio.Writer = bufio.NewWriter(os.Stdout)
	stderr *bufio.Writer = bufio.NewWriter(os.Stderr)
	stdin  *bufio.Reader = bufio.NewReader(os.Stdin )
	args map[string][]string = argvAssing(os.Args)
	argv = os.Args[1:]
	argc = len(os.Args)-1
	format = fmt.Sprintf
	printf = fmt.Printf
	sprintf = fmt.Sprintf
	fs = fmt.Sprintf
	fprintf = fmt.Fprintf
	join = strings.Join
	split = strings.Split
	fopen = os.Open
	fmake = os.Create
	fwriter = bufio.NewWriter
	freader = bufio.NewReader
	NULL = interface{}(nil)
)

//const def
const (
	// nums
	I8MAX =  int8(0x7f)
	I8MIN =  int8(-0x80)
	//U8MAX =  uint8(0xFF)

	I16MAX = int16(0x7FFF)
	I16MIN = int16(-0x8000)
	//U16MAX = uint16(0xFFFF)

	I32MAX = int32(0x7FFFFFFF)
	I32MIN = int32(-0x80000000)
	//U32MAX = uint32(0xFFFFFFFF)

	I64MAX = int64(0x7FFFFFFFFFFFFFFF)
	I64MIN = int64(-0x8000000000000000)
	//U64MAX = uint64(0xFFFFFFFFFFFFFFFF)

	// file nums
	F_append = os.O_APPEND
	F_WR = os.O_WRONLY

)



//typedef
type FILE = os.File
type reader = *bufio.Reader
type writer = *bufio.Writer
