package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"math/rand"
	"runtime"
	"strings"
	"time"

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

func Sqrt(x float64) float64 {
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
	return z
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
	fmt.Printf("Sqrt(%v), Me (%v) vs math (%v)\n", x, Sqrt(x), math.Sqrt(x))

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
}
