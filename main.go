package main

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"math"
	"math/cmplx"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"golang.org/x/tour/tree"
	"golang.org/x/tour/wc"
)

const (
	// Create a huge number by shifting a 1 bit left 100 places.
	// In other words, the binary number that is 1 followed by 100 zeroes.
	Big = 1 << 100
	// Shift it right again 99 places, so we end up with 1<<1, or 2.
	Small = Big >> 99
)

var (
	ToBe   bool       = false
	MaxInt uint64     = 1<<64 - 1
	z      complex128 = cmplx.Sqrt(-5 + 12i)
)

type MyFloat float64

func (f MyFloat) Abs() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

type Abser interface {
	Abs() float64
}

type Vertex struct {
	X float64
	Y float64
}

func (v *Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vertex) Scale(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func needInt(x int) int { return x*10 + 1 }
func needFloat(x float64) float64 {
	return x * 0.1
}

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %v", float64(e))
}

func Sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, ErrNegativeSqrt(x)
	}
	z := float64(1)
	i := 0
	for {
		p := z
		z -= ((z*z - x) / (2 * z))

		if d := math.Abs(p - z); d < 1e-6 {
			break
		}
		if i < 10 {
			i++
		} else {
			break
		}
	}
	return z, nil
}

func describe(i interface{}) {
	fmt.Printf("(%v, %T)\n", i, i)
}

func returnMultiple() (int, int) {
	return 1, 2
}

type MyFunctionOptions struct {
	Name           string
	Age            int
	withNameAndAge func(nameAndAge string)
}

func MyFunction(options MyFunctionOptions) {
	if options.Name == "" {
		options.Name = "Default Name"
	}
	if options.Age == 0 {
		options.Age = 25
	}
	if options.withNameAndAge != nil {
		options.withNameAndAge(fmt.Sprint(options.Name, " ", options.Age))
	}
}

type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return fmt.Sprintf("%v (%v years)", p.Name, p.Age)
}

type IPAddr [4]byte

func (ipaddr IPAddr) String() string {
	result := ""
	for i, v := range ipaddr {
		result += fmt.Sprintf("%d", v)
		if i != len(ipaddr)-1 {
			result += "."
		}
	}
	return result
}

type rot13Reader struct {
	r io.Reader
}

func (r rot13Reader) Read(b []byte) (int, error) {
	n, err := r.r.Read(b)
	if err == nil {
		for i := range b[:n] {
			if b[i] >= 'A' && b[i] <= 'Z' {
				b[i] = 'A' + (b[i]-'A'+13)%26
			} else if b[i] >= 'a' && b[i] <= 'z' {
				b[i] = 'a' + (b[i]-'a'+13)%26
			}
		}
	}
	return n, err
}

type Image struct {
	w int
	h int
}

func (m Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, m.w, m.h)
}

func (m Image) ColorModel() color.Model {
	return color.RGBAModel
}

func (m Image) At(x, y int) color.Color {
	r := uint8((x + y) / 2)
	g := uint8(x * y)
	return color.RGBA{r, g, 255, 255}
}

func Index[T comparable](s []T, x T) int {
	for i, v := range s {
		if v == x {
			return i
		}
	}
	return -1
}

type List[T any] struct {
	next *List[T]
	val  T
}

func (l List[any]) String() string {
	if l.next == nil {
		return fmt.Sprintf("%v", l.val)
	} else {
		return fmt.Sprintf("%v -> %v", l.val, l.next)
	}
}

type Stack[T any] struct {
	Value    T
	Previous *Stack[T]
}

func Push[T any](s *Stack[T], t T) *Stack[T] {
	return &Stack[T]{Value: t, Previous: s}
}

func Pop[T any](s *Stack[T]) (*Stack[T], T) {
	return s.Previous, s.Value
}

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	var s *Stack[*tree.Tree]
	n := t
	for n != nil {
		s = Push(s, n)
		n = n.Left
	}
	for s != nil {
		s, n = Pop(s)
		ch <- n.Value
		if n.Right != nil {
			n = n.Right
			for n != nil {
				s = Push(s, n)
				n = n.Left
			}
		}
	}
	close(ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	c1 := make(chan int)
	c2 := make(chan int)

	go Walk(t1, c1)
	go Walk(t2, c2)

	for {
		i1, ok1 := <-c1
		i2, ok2 := <-c2
		if !ok1 && !ok2 {
			return true
		} else if ok1 != ok2 || i1 != i2 {
			return false
		}
	}
}

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type Seen struct {
	mut  sync.Mutex
	urls map[string]bool
}

func (s *Seen) Has(url string) bool {
	s.mut.Lock()
	defer s.mut.Unlock()
	return s.urls[url]
}

func (s *Seen) Set(url string) {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.urls[url] = true
}

var seen = &Seen{urls: make(map[string]bool)}
var crawlWg = sync.WaitGroup{}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	crawlWg.Add(1)
	defer crawlWg.Done()

	if depth <= 0 {
		return
	}
	// Don't fetch the same URL twice.
	if seen.Has(url) {
		return
	}
	seen.Set(url)

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		// Fetch URLs in parallel.
		go Crawl(u, depth-1, fetcher)
	}
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}

func main() {
	var i = 42
	var f = float64(i)
	var u = uint(f)

	fmt.Printf("Type: %T Value: %v\n", ToBe, ToBe)
	fmt.Printf("Type: %T Value: %v\n", MaxInt, MaxInt)
	fmt.Printf("Type: %T Value: %v\n", z, z)
	fmt.Printf("Type: %T Value: %v\n", i, i)
	fmt.Printf("Type: %T Value: %v\n", f, f)
	fmt.Printf("Type: %T Value: %v\n", u, u)
	fmt.Println(needInt(Small))
	fmt.Println(needFloat(Small))
	fmt.Println(needFloat(Big))

	sum := 0
	for i := 0; i < 22; i++ {
		sum += i
	}
	fmt.Println(sum)

	for sum = 1; sum < 1000; {
		sum += sum
	}
	fmt.Println(sum)

	sum = 1
	for {
		sum += sum
		if sum > 1000 {
			break
		}
	}
	fmt.Println(sum)

	rand.Seed(time.Now().UnixNano())

	if v := rand.Float64(); v < 0.5 {

		fmt.Printf("Yes: %f\n", v)
	} else {
		fmt.Printf("No: %f\n", v)

	}

	fmt.Printf("-0 == 0: %v\n", -0 == 0)

	x := rand.Float64() * 500
	sqrt, _ := Sqrt(x)
	fmt.Printf("Sqrt(%v), Me (%v) vs math (%v)\n", x, sqrt, math.Sqrt(x))

	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("OS X.")
	case "linux":
		fmt.Println("Linux.")
	default:
		// freebsd, openbsd,
		// plan9, windows...
		fmt.Printf("%s.\n", os)
	}

	switch today := time.Now().Weekday(); time.Saturday {
	case today + 0:
		fmt.Println("Today.")
	case today + 1:
		fmt.Println("Tomorrow.")
	case today + 2:
		fmt.Println("In two days.")
	default:
		fmt.Println("Too far away.")
	}

	a := float64(2)
	sl := float64(2)
	defer fmt.Printf("%v ^ %v = %v\n", a, sl, math.Pow(a, sl))
	sl = 4

	defer func() {
		fmt.Printf("%v ^ %v = %v\n", a, sl, math.Pow(a, sl))
	}()

	i, j := 42, 2701

	p := &i         // point to i
	fmt.Println(*p) // read i through the pointer
	*p = 21         // set i through the pointer
	fmt.Println(i)  // see the new value of i

	p = &j         // point to j
	*p = *p / 37   // divide j through the pointer
	fmt.Println(j) // see the new value of j

	s := "Hello"
	q := &s
	fmt.Println(*q + " Aravindan!")

	fmt.Println(Vertex{1, 2})

	v := Vertex{1, 2}
	v.X = 4
	fmt.Println(v.X)

	r := &v
	r.X = 1e9

	t := v
	t.X = 1e9
	fmt.Printf("v: %v, r: %v, t: %v\n", v, r, t)

	// int
	modifyInt := func(n int) int {
		return n + 5
	}

	age := 30
	fmt.Println("Before function call: ", age)
	fmt.Println("Function call:", modifyInt(age))
	fmt.Println("After function call: ", age)

	// float
	modifyFloat := func(n float64) float64 {
		return n + 5.0
	}
	cash := 10.50
	fmt.Println("Before function call: ", cash)
	fmt.Println("Function call:", modifyFloat(cash))
	fmt.Println("After function call: ", cash)

	// bool
	modifyBool := func(n bool) bool {
		return !n
	}
	old := false
	fmt.Println("Before function call: ", old)
	fmt.Println("Function call:", modifyBool(old))
	fmt.Println("After function call: ", old)

	// string
	modifyString := func(n string) string {
		return n + " Golang"
	}
	message := "Go"
	fmt.Println("Before function call: ", message)
	fmt.Println("Function call:", modifyString(message))
	fmt.Println("After function call: ", message)

	// array

	modifyArray := func(coffee [3]string) [3]string {
		coffee[2] = "germany"
		return coffee
	}
	country := [3]string{"nigeria", "egypt", "sweden"}
	fmt.Println("Before function call: ", country)
	// [nigeria egypt sweden]
	fmt.Println("Function call:", modifyArray(country))
	// [nigeria egypt germany]
	fmt.Println("After function call: ", country)
	// [nigeria egypt sweden]

	// Profile contains user data
	type Profile struct {
		Age          int
		Name         string
		Salary       float64
		TechInterest bool
	}
	myProfile := Profile{
		Age:          15,
		Name:         "Adeshina",
		Salary:       300,
		TechInterest: false,
	}
	modifyStruct := func(p Profile) Profile {
		p.Age = 85
		p.Name = "Balqees"
		p.Salary = 500.45
		p.TechInterest = true

		return p
	}
	fmt.Println("Before function call: ", myProfile)
	fmt.Println("Function call:", modifyStruct(myProfile))
	fmt.Println("After function call: ", myProfile)
	fmt.Printf("modifyStruct profiles equal %v\n", modifyStruct(myProfile) == myProfile)

	modifyStructByPointer := func(p *Profile) Profile {
		p.Age = 85
		p.Name = "Balqees"
		p.Salary = 500.45
		p.TechInterest = true

		return *p
	}
	fmt.Println("Before function call: ", myProfile)
	fmt.Println("Function call:", modifyStructByPointer(&myProfile))
	fmt.Println("After function call: ", myProfile)
	fmt.Printf("modifyStructByPointer profiles equal %v\n", modifyStructByPointer(&myProfile) == myProfile)

	arr := []string{"hello", "world"}
	arr = append(arr, []string{"hello", "world"}...)
	fmt.Printf("Array: %v\n", arr)

	primes := [6]int{2, 3, 5, 7, 11, 13}
	slice := primes[1:4]

	fmt.Println(slice)

	arr1 := []int{1, 2, 3}
	slice1 := arr1[1:2]
	slice1 = append(slice1, 4)
	slice1 = append(slice1, []int{5, 6}...)
	fmt.Println(slice1)
	fmt.Println(arr1)

	slice1 = arr1[:0]
	fmt.Println(cap(slice1))
	fmt.Println(len(slice1))
	slice1 = slice1[:3]
	fmt.Println(slice1)

	var s1 []int
	fmt.Println(s1, len(s1), cap(s1))
	if s1 == nil {
		fmt.Println("nil!")
	}
	s1 = append(s1, []int{6, 7, 8}...)
	fmt.Println(s1, len(s1), cap(s1))

	if 1 == 1 {
		fmt.Println("1 == 1")
	}
	if true == true {
		fmt.Println("true == true")
	}
	if "hello" == "hello" {
		fmt.Println("\"hello\" == \"hello\"")
	}
	if [10]int{} == [10]int{} {
		fmt.Println("[10]int{} == [10]int{}")
	}
	if [10]int{} == [10]int{1} {
		fmt.Println("[10]int{} == [10]int{1}")
	}
	if [10]int{1} == [10]int{1} {
		fmt.Println("[10]int{1} == [10]int{1}")
	}

	if (Profile{}) == (Profile{}) {
		fmt.Println("(Profile{}) == (Profile{})")
	}

	if (&Profile{}) == (&Profile{}) {
		fmt.Println("(&Profile{}) == (&Profile{})")
	}

	if (Profile{Name: "Ara"}) == (Profile{Name: "Ara"}) {
		fmt.Println("(Profile{Name: \"Ara\"}) == (Profile{Name: \"Ara\"})")
	}

	if (Profile{Name: "Arav"}) == (Profile{Name: "Ara"}) {
		fmt.Println("(Profile{Name: \"Arav\"}) == (Profile{Name: \"Ara\"})")
	}

	p1 := &Profile{Name: "Ara"}
	p2 := &Profile{Name: "Ara"}
	p3 := &p1
	p4 := &p3
	fmt.Printf("p3 == %v\n", p3)
	fmt.Printf("p4 == %v\n", p4)
	if p1 == p2 {
		fmt.Printf("%v == %v\n", p1, p2)
	} else {
		fmt.Printf("%v != %v\n", &p1, &p2)
	}

	sl1 := make([]int, 0, 5) // len(b)=0, cap(b)=5

	sl1 = sl1[:cap(sl1)] // len(b)=5, cap(b)=5
	sl1 = sl1[1:]        // len(b)=4, cap(b)=4

	type Number int

	var number Number = 1
	fmt.Printf("number: %v\n", number)

	pow := []int{1, 2, 4, 8, 16, 32, 64, 128}
	for i, v := range pow {
		fmt.Printf("2**%d = %d\n", i, v)
	}

	// pic.Show(func(dx, dy int) [][]uint8 {
	// 	s := make([][]uint8, dy)
	// 	for i := range s {
	// 		s[i] = make([]uint8, dx)
	// 		for j := range s[i] {
	// 			// s[i][j] = uint8((i + j) / 2)
	// 			s[i][j] = uint8(i * j)
	// 		}
	// 	}
	// 	return s
	// })

	type Location struct {
		Lat, Long float64
	}

	m := make(map[string]Location)
	m["Bell Labs"] = Location{
		40.68433, -74.39967,
	}
	fmt.Println(m)

	m2 := make(map[Location]Location)
	m2[Location{}] = Location{
		40.68433, -74.39967,
	}
	fmt.Println(m2[Location{}])

	m3 := make(map[string]map[string]Location)
	m3["Locations"] = map[string]Location{
		"Bell Labs": {
			40.68433, -74.39967,
		},
	}
	fmt.Println(m3)

	delete(m3, "Locations")

	m3Locations, ok := m3["Locations"]
	fmt.Printf("m3Locations: %v, ok: %v, m3Locations == nil: %v\n", m3Locations, ok, m3Locations == nil)

	wc.Test(func(s string) map[string]int {
		m := map[string]int{}
		for _, v := range strings.Fields(s) {
			m[v]++
		}
		return m
	})

	adder := func() func(int) int {
		sum := 0
		return func(x int) int {
			sum += x
			return sum
		}
	}

	pos, neg := adder(), adder()
	for i := 0; i < 10; i++ {
		fmt.Println(
			pos(i),
			neg(-2*i),
		)
	}

	fibonacci := func() func() int {
		n2, n1 := -1, 1
		return func() int {
			n0 := n2 + n1
			n2 = n1
			n1 = n0
			return n0
		}
	}

	f1 := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f1())
	}

	v2 := Vertex{3, 4}
	fmt.Println("v2.Abs()", v2.Abs())

	f2 := MyFloat(-math.Sqrt2)
	fmt.Println("f2.Abs()", f2.Abs())

	v3 := Vertex{3, 4}
	v3.Scale(10)
	fmt.Println("v3.Scale(10)", v3.Abs())

	var a1 Abser
	a1 = &Vertex{3, 4}
	fmt.Printf("a1: %v\n", a1.Abs())

	v4 := Vertex{3, 4}
	(&v4).X = 1
	p5 := &v4
	v4.Abs()
	(*p5).Abs()

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Printf("Recovered from %v\n", r)
	// 	}
	// }()

	// var v5 *Vertex = nil
	// fmt.Printf("v5.Abs(): %v\n", v5.Abs())

	var v6 Vertex
	v6.Abs()

	// var a2 Abser
	// a2.Abs()

	describe(v6)

	x1, y1 := returnMultiple()
	fmt.Printf("x1: %v, y1: %v\n", x1, y1)

	var i7 interface{} = "hello"

	s3 := i7.(string)
	fmt.Println(s3)

	s3, ok = i7.(string)
	fmt.Println(s3, ok)

	f5, ok := i7.(float64)
	fmt.Println(f5, ok)

	// f5 = i7.(float64) // panic
	// fmt.Println(f5)

	switch v := i7.(type) {
	case int:
		fmt.Printf("Twice %v is %v\n", v, v*2)
	case string:
		fmt.Printf("%q is %v bytes long\n", v, len(v))
	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}

	MyFunction(MyFunctionOptions{})
	MyFunction(MyFunctionOptions{Name: "Aravindan"})
	MyFunction(MyFunctionOptions{Age: 30})
	MyFunction(MyFunctionOptions{withNameAndAge: func(nameAndAge string) {
		fmt.Println("Name and Age:", nameAndAge)
	}})
	MyFunction(MyFunctionOptions{Name: "Aravindan", withNameAndAge: func(nameAndAge string) {
		fmt.Println("Name and Age:", nameAndAge)
	}})
	MyFunction(MyFunctionOptions{Age: 30, withNameAndAge: func(nameAndAge string) {
		fmt.Println("Name and Age:", nameAndAge)
	}})

	var fn func()
	describe(fn)
	fmt.Println("fn == nil", fn == nil)

	var i3 interface{}
	describe(i3)
	fmt.Println("i3 == nil", i3 == nil)

	hosts := map[string]IPAddr{
		"loopback":  {127, 0, 0, 1},
		"googleDNS": {8, 8, 8, 8},
	}
	for name, ip := range hosts {
		fmt.Printf("%v: %v\n", name, ip)
	}

	fmt.Printf("io.EOF: %v\n", io.EOF.Error())

	r1 := strings.NewReader("Hello, Reader!")

	b := make([]byte, 8)
	for {
		n, err := r1.Read(b)
		fmt.Printf("n = %v err = %v b = %v\n", n, err, b)
		fmt.Printf("b[:n] = %q\n", b[:n])
		if err == io.EOF {
			break
		}
	}

	fmt.Printf("'A': %v\n", 'A')

	s6 := strings.NewReader("Lbh penpxrq gur pbqr!")
	r6 := rot13Reader{s6}
	io.Copy(os.Stdout, &r6)

	im1 := image.NewRGBA(image.Rect(0, 0, 100, 100))
	fmt.Println(im1.Bounds())
	fmt.Println(im1.At(0, 0).RGBA())

	// im2 := Image{100, 100}
	// pic.ShowImage(im2)

	// Index works on a slice of ints
	si := []int{10, 20, 15, -10}
	fmt.Println(Index(si, 15))

	// Index also works on a slice of strings
	ss := []string{"foo", "bar", "baz"}
	fmt.Println(Index(ss, "hello"))

	list1 := List[int]{val: 1}
	describe(list1)

	list2 := List[int]{val: 0}
	for l, i := &list2, 1; i < 10; i++ {
		l.next = &List[int]{val: i}
		l = l.next
	}
	fmt.Println(list2)

	go func() { fmt.Println(" world!") }()
	fmt.Print("hello")
	time.Sleep(100 * time.Millisecond)

	c1 := make(chan string, 2)

	go func() {
		time.Sleep(1000 * time.Millisecond)
		c1 <- "hello after sometime!"
	}()

	go func() {
		time.Sleep(2000 * time.Millisecond)
		c1 <- "hello after more time!"
	}()

	x2, x3 := <-c1, <-c1
	fmt.Printf("x2: %v, x3: %v\n", x2, x3)

	fibonacci2 := func(n int, c chan int) {
		x, y := 0, 1
		for i := 0; i < n; i++ {
			c <- x
			x, y = y, x+y
		}
		close(c)
	}

	c2 := make(chan int, 10)
	go fibonacci2(cap(c2), c2)
	for i := range c2 {
		fmt.Println(i)
	}

	x4, ok := <-c2
	fmt.Printf("x4: %v, ok: %v\n", x4, ok)

	fibonacci3 := func(c, quit chan int) {
		x, y := 0, 1
		for {
			select {
			case c <- x:
				x, y = y, x+y
			case <-quit:
				fmt.Println("quit")
				return
			}
		}
	}

	c3 := make(chan int)
	quit := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c3)
		}
		quit <- 0
	}()
	fibonacci3(c3, quit)

	tick := time.Tick(1000 * time.Millisecond)
	boom := time.After(5000 * time.Millisecond)

Awaiter:
	for {
		select {
		case <-tick:
			fmt.Println("tick.")
		case <-boom:
			fmt.Println("BOOM!")
			break Awaiter
		default:
			fmt.Println("    .")
			time.Sleep(500 * time.Millisecond)
		}
	}

	ch := make(chan int)
	go Walk(tree.New(1), ch)
	fmt.Print("Tree: ")
	for i := range ch {
		fmt.Print(i, " ")
	}
	fmt.Println()

	fmt.Println("Same(tree.New(1), tree.New(1))", Same(tree.New(1), tree.New(1)))
	fmt.Println("Same(tree.New(1), tree.New(2))", Same(tree.New(1), tree.New(2)))
	t1trim := tree.New(1)
	t1trim.Right = nil
	ch2 := make(chan int)
	go Walk(t1trim, ch2)
	fmt.Print("Tree t1trim: ")
	for i := range ch2 {
		fmt.Print(i, " ")
	}
	fmt.Println()
	fmt.Println("Same(tree.New(1), t1trim)", Same(tree.New(1), t1trim))

	slice2 := []int{1, 2, 3}
	fmt.Printf("slice2: %v\n", slice2)

	array2 := [...]int{1, 2, 3}
	fmt.Printf("array2: %v\n", array2)

	Crawl("https://golang.org/", 4, fetcher)
	crawlWg.Wait()
	time.Sleep(1000 * time.Millisecond)
}
