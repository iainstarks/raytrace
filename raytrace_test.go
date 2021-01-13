package main

import (
	"flag"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"math"
	"os"
	r "reflect"
	"testing"
)

// go test support for godog
var opt = godog.Options{Output: colors.Uncolored(os.Stdout), Format: "progress"}

var a bool

func equalFloats(a float64, b float64) bool {
	return (math.Abs(a-b) < 0.0001)
}

var tupleMap = make(map[string]*tuple)

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

func TestMain(m *testing.M) {
	flag.Parse()
	opt.Paths = flag.Args()

	//	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
	//		FeatureContext(s)
	//	}, opt)

	status := godog.TestSuite{
		Name:                 "raytrace",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options:              &opt,
	}.Run()
	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {})
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(func(*godog.Scenario) {
	})

	ctx.Step(`^(\w+) ← tuple\(([\d.-]+), ([\d.-]+), ([\d.-]+), ([\d.-]+)\)$`, addTuple)
	ctx.Step(`^(\w+) = tuple\(([\d.-]+), ([\d.-]+), ([\d.-]+), ([\d.-]+)\)$`, matchesTuple)
	ctx.Step(`^-(\w+) = tuple\(([\d.-]+), ([\d.-]+), ([\d.-]+), ([\d.-]+)\)$`, isNegatedTuple)
	ctx.Step(`^(\w+) ← point\(([\d.-]+), ([\d.-]+), ([\d.-]+)\)$`, addPoint)
	ctx.Step(`^(\w+) ← vector\(([\d.-]+), ([\d.-]+), ([\d.-]+)\)$`, addVector)
	ctx.Step(`^(\w+) = vector\(([\d.-]+), ([\d.-]+), ([\d.-]+)\)$`, matchesVector)
	ctx.Step(`^(\w+) is a point$`, isAPoint)
	ctx.Step(`^(\w+) is a vector$`, isAVector)
	ctx.Step(`^(\w+) is not a point$`, isNotAPoint)
	ctx.Step(`^(\w+) is not a vector$`, isNotAVector)
	ctx.Step(`^(\w+)\.w = ([\d.-]+)$`, hasW)
	ctx.Step(`^(\w+)\.x = ([\d.-]+)$`, hasX)
	ctx.Step(`^(\w+)\.y = ([\d.-]+)$`, hasY)
	ctx.Step(`^(\w+)\.z = ([\d.-]+)$`, hasZ)
	ctx.Step(`^(\w+) \+ (\w+) = tuple\(([\d.-]+), ([\d.-]+), ([\d.-]+), ([\d.-]+)\)$`, addTuples)
	ctx.Step(`^(\w+) - (\w+) = vector\(([\d.-]+), ([\d.-]+), ([\d.-]+)\)$`, subtractVector)
	ctx.Step(`^(\w+) - (\w+) = point\(([\d.-]+), ([\d.-]+), ([\d.-]+)\)$`, subtractPoint)
	ctx.Step(`^(\w+) \* ([\d.-]+) = tuple\(([\d.-]+), ([\d.-]+), ([\d.-]+), ([\d.-]+)\)$`, multiplyTuples)
	ctx.Step(`^(\w+) \/ ([\d.-]+) = tuple\(([\d.-]+), ([\d.-]+), ([\d.-]+), ([\d.-]+)\)$`, divideTuples)
	ctx.Step(`^magnitude\((\w+)\) = (\d+)$`, compare_magnitude)
	ctx.Step(`^magnitude\((\w+)\) = √(\d+)$`, compare_magnitude_squared)
	ctx.Step(`^(\w+) ← normalize\((\w+)\)$`, normalize_and_add)
	ctx.Step(`^normalize\((\w+)\) = approximately vector\(([\d.-]+), ([\d.-]+), ([\d.-]+)\)$`, normalize_and_compare_vector)
	ctx.Step(`^normalize\((\w+)\) = vector\(([\d.-]+), ([\d.-]+), ([\d.-]+)\)$`, normalize_and_compare_vector)
	ctx.Step(`^cross\((\w+), (\w+)\) = vector\(([\d.-]+), ([\d.-]+), ([\d.-]+)\)$`, cross_compare)
	ctx.Step(`^dot\((\w+), (\w+)\) = ([\d.-]+)$`, dot_compare)
	// chapter 2
	ctx.Step(`^(\w+) ← color\(([\d.-]+), ([\d.-]+), ([\d.-]+)\)$`, addColor)
	ctx.Step(`^(\w+) - (\w+) = color\(([\d.-]+), ([\d.-]+), ([\d.-]+)\)$`, subtractColorsAndCompare)
	ctx.Step(`^(\w+) \* ([\d.-]+) = color\(([\d.-]+), ([\d.-]+), ([\d.-]+)\)$`, scalarMultiplyColorAndCompare)
	ctx.Step(`^(\w+) \* (\w+) = color\(([\d.-]+), ([\d.-]+), ([\d.-]+)\)$`, hadamardProductColorsAndCompare)
	ctx.Step(`^(\w+) \+ (\w+) = color\(([\d.-]+), ([\d.-]+), ([\d.-]+)\)$`, addColorsAndCompare)
	ctx.Step(`^(\w+)\.red = ([\d.-]+)$`, hasX)
	ctx.Step(`^(\w+)\.green = ([\d.-]+)$`, hasY)
	ctx.Step(`^(\w+)\.blue = ([\d.-]+)$`, hasZ)
	ctx.Step(`^(\w+) ← reflect\((\w+), (\w+)\)$`, reflectAndAdd)
	/*
		ctx.Step(`^r = vector\((\d+), (\d+), (\d+)\)$`, rVector)
	*/
}

func addColor(name string, x, y, z float64) error {
	return addVector(name, x, y, z)
}

func subtractColorsAndCompare(name1, name2 string, x, y, z float64) error {
	return subtractVector(name1, name2, x, y, z)
}

func reflectAndAdd(resultName, name1, name2 string) error {
	if a, ok := tupleMap[name1]; ok {
		if b, ok := tupleMap[name2]; ok {
			result := reflect(a, b)
			addTuple(resultName, result.x, result.y, result.z, result.w)
			return nil
		}
	}
	return fmt.Errorf("not found %s or %s", name1, name2)
}

func hadamardProductColorsAndCompare(name1, name2 string, x, y, z float64) error {
	if a, ok := tupleMap[name1]; ok {
		if b, ok := tupleMap[name2]; ok {
			t := hadamardProduct(a, b)
			if equalFloats(t.x, x) && equalFloats(t.y, y) && equalFloats(t.z, z) && equalFloats(t.w, 0.0) {
				return nil
			}
			return fmt.Errorf("cross product does not match values %f %f %f %f != %f %f %f %f", t.x, t.y, t.z, t.w, x, y, z, 0.0)
		}
	}
	return fmt.Errorf("not found %s or %s", name1, name2)
}

func addColorsAndCompare(name1, name2 string, x, y, z float64) error {
	return addTuples(name1, name2, x, y, z, 0.0)
}

func scalarMultiplyColorAndCompare(name string, n, x, y, z float64) error {
	return multiplyTuples(name, n, x, y, z, 0.0)
}

func cross_compare(name1, name2 string, x, y, z float64) error {
	if a, ok := tupleMap[name1]; ok {
		if b, ok := tupleMap[name2]; ok {
			t := cross(a, b)
			if equalFloats(t.x, x) && equalFloats(t.y, y) && equalFloats(t.z, z) && equalFloats(t.w, 0.0) {
				return nil
			}
			return fmt.Errorf("cross product does not match values %f %f %f %f != %f %f %f %f", t.x, t.y, t.z, t.w, x, y, z, 0.0)
		}
	}
	return fmt.Errorf("not found %s or %s", name1, name2)
}

func dot_compare(name1, name2 string, v float64) error {
	if a, ok := tupleMap[name1]; ok {
		if b, ok := tupleMap[name2]; ok {
			t := dot(a, b)
			if equalFloats(t, v) {
				return nil
			}
			return fmt.Errorf("Dot product does not match expected")
		}
	}
	return fmt.Errorf("not found %s or %s", name1, name2)
}

func normalize_and_add(normal_name string, name string) error {
	if t, ok := tupleMap[name]; ok {
		n := normalize(t)
		tupleMap[normal_name] = n
		return nil
	}
	return fmt.Errorf("not found %s", name)
}

func normalize_and_compare_vector(name string, x float64, y float64, z float64) error {
	if t, ok := tupleMap[name]; ok {
		n := normalize(t)
		if equalFloats(n.x, x) && equalFloats(n.y, y) && equalFloats(n.z, z) && equalFloats(n.w, 0.0) {
			return nil
		}
		return fmt.Errorf("Normalized vector does not match values %f %f %f %f != %f %f %f %f", n.x, n.y, n.z, n.w, x, y, z, 0.0)
	}
	return fmt.Errorf("not found %s", name)
}

// BDD mapping functions
func fTrue(name string, v string) error {
	if name != "f" {
		return fmt.Errorf("Not a %s", name)
	}
	//if v != "1.2" {
	//	return fmt.Errorf("No match %s", v)
	//}
	a = true
	return nil
}

func matchesTuple(name string, x, y, z, w float64) error {
	if t, ok := tupleMap[name]; ok {
		if equalFloats(t.x, x) && equalFloats(t.y, y) && equalFloats(t.z, z) && equalFloats(t.w, w) {
			return nil
		}
		return fmt.Errorf("tuple does not match values %f %f %f %f != %f %f %f %f", t.x, t.y, t.z, t.w, x, y, z, w)
	}
	return fmt.Errorf("not found %s", name)
}

func matchesVector(name string, x, y, z float64) error {
	return matchesTuple(name, x, y, z, 0.0)
}

func isNegatedTuple(name string, x, y, z, w float64) error {
	if t, ok := tupleMap[name]; ok {
		if t.x == -x && t.y == -y && t.z == -z && t.w == -w {
			return nil
		}
		return fmt.Errorf("Not a point and should be")
	}
	return fmt.Errorf("not found %s", name)
}

// TODO these will need to use safe floating point comparison not ==
func addTuples(name1 string, name2 string, x, y, z, w float64) error {
	if a, ok := tupleMap[name1]; ok {
		if b, ok := tupleMap[name2]; ok {
			t := add(a, b)
			if t.x == x && t.y == y && t.z == z && t.w == w {
				return nil
			}
			return fmt.Errorf("no match ( %f %f %f %f ) != ( %f %f %f %f )", t.x, t.y, t.z, t.w, x, y, z, w)
		}
	}
	return fmt.Errorf("not found %s or %s", name1, name2)
}

func multiplyTuples(name string, m, x, y, z, w float64) error {
	if t, ok := tupleMap[name]; ok {
		r := multiply(t, m)
		if r.x == x && r.y == y && r.z == z && r.w == w {
			return nil
		}
		return fmt.Errorf("no match ( %f %f %f %f ) != ( %f %f %f %f )", r.x, r.y, r.z, r.w, x, y, z, w)
		return nil
	}
	return fmt.Errorf("not found %s ", name)
}

func divideTuples(name string, m, x, y, z, w float64) error {
	if t, ok := tupleMap[name]; ok {
		if t.x/m == x && t.y/m == y && t.z/m == z && t.w/m == w {
			return nil
		}
		return fmt.Errorf("Not a point and should be")
		return nil
	}
	return fmt.Errorf("not found %s ", name)
}

func subtractVector(name1, name2 string, x, y, z float64) error {
	return subtractTuples(name1, name2, x, y, z, 0.0)
}

func subtractPoint(name1, name2 string, x, y, z float64) error {
	return subtractTuples(name1, name2, x, y, z, 1.0)
}

func subtractTuples(name1, name2 string, x, y, z, w float64) error {
	if a, ok := tupleMap[name1]; ok {
		if b, ok := tupleMap[name2]; ok {
			t := sub(a, b)
			if equalFloats(t.x, x) && equalFloats(t.y, y) && equalFloats(t.z, z) && equalFloats(t.w, w) {
				return nil
			}
			return fmt.Errorf("Not a point and should be %f %f", t.x, x)
			return nil
		}
	}
	return nil
}

func addTuple(name string, x, y, z, w float64) error {
	tupleMap[name] = NewTuple(x, y, z, w)
	return nil
}

func addPoint(name string, x, y, z float64) error {
	tupleMap[name] = NewTuple(x, y, z, 1.0)
	return nil
}

func addVector(name string, x, y, z float64) error {
	tupleMap[name] = NewTuple(x, y, z, 0.0)
	return nil
}

func compare_magnitude(name string, v float64) error {
	return compare_magnitude_squared(name, v*v)
}

func compare_magnitude_squared(name string, v float64) error {
	if t, ok := tupleMap[name]; ok {
		m := magnitude_squared(t)
		if equalFloats(v, m) {
			return nil
		}
		return fmt.Errorf("Magnitude squared does not match %f != %f", v, m)

	}
	return fmt.Errorf("not found %s ", name)
}

func isPoint(t *tuple) bool {
	return equalFloats(t.w, 1.0)
}

func isVector(t *tuple) bool {
	return !isPoint(t)
}

func isX(t *tuple, x float64) bool {
	return equalFloats(t.x, x)
}

func isY(t *tuple, y float64) bool {
	return equalFloats(t.y, y)
}

func isZ(t *tuple, z float64) bool {
	return equalFloats(t.z, z)
}

func isW(t *tuple, w float64) bool {
	return equalFloats(t.w, w)
}

func isAPoint(name string) error {
	if t, ok := tupleMap[name]; ok {
		if isPoint(t) {
			return nil
		}
		return fmt.Errorf("Not a point and should be")
	}
	return fmt.Errorf("not found %s", name)
}

func isAVector(name string) error {
	if t, ok := tupleMap[name]; ok {
		if !isPoint(t) {
			return nil
		}
		return fmt.Errorf("Not a point and should be")
	}
	return fmt.Errorf("not found %s", name)
}

func isNotAPoint(name string) error {
	return isAVector(name)
}

func isNotAVector(name string) error {
	return isAPoint(name)
}

func hasW(name string, w float64) error {
	if t, ok := tupleMap[name]; ok {
		if isW(t, w) {
			return nil
		}
		return fmt.Errorf("w %f does not match %f", t.w, w)
	}
	return fmt.Errorf("not found %s", name)
}

func hasX(name string, x float64) error {
	if t, ok := tupleMap[name]; ok {
		if isX(t, x) {
			return nil
		}
		return fmt.Errorf("x %f does not match %f", t.x, x)
	}
	return fmt.Errorf("not found %s", name)
}

func hasY(name string, y float64) error {
	if t, ok := tupleMap[name]; ok {
		if isY(t, y) {
			return nil
		}
		return fmt.Errorf("w %f does not match %f", t.y, y)
	}
	return fmt.Errorf("not found %s", name)
}

func hasZ(name string, z float64) error {
	if t, ok := tupleMap[name]; ok {
		if isZ(t, z) {
			return nil
		}
		return fmt.Errorf("z %f does not match %f", t.z, z)
	}
	return fmt.Errorf("not found %s", name)
}

// PENDING

//Scenario: point() creates tuples with w=1
//Given p ← point(4, -4, 3)
//Then p = tuple(4, -4, 3, 1)

func TestNewPoint(t *testing.T) {
	type args struct {
		x float64
		y float64
		z float64
	}
	tests := []struct {
		name string
		args args
		want *tuple
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPoint(tt.args.x, tt.args.y, tt.args.z); !r.DeepEqual(got, tt.want) {
				t.Errorf("NewPoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTuple(t *testing.T) {
	type args struct {
		x float64
		y float64
		z float64
		w float64
	}
	tests := []struct {
		name string
		args args
		want *tuple
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTuple(tt.args.x, tt.args.y, tt.args.z, tt.args.w); !r.DeepEqual(got, tt.want) {
				t.Errorf("NewTuple() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewVector(t *testing.T) {
	type args struct {
		x float64
		y float64
		z float64
	}
	tests := []struct {
		name string
		args args
		want *tuple
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewVector(tt.args.x, tt.args.y, tt.args.z); !r.DeepEqual(got, tt.want) {
				t.Errorf("NewVector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_add(t *testing.T) {
	type args struct {
		a *tuple
		b *tuple
	}
	tests := []struct {
		name string
		args args
		want *tuple
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := add(tt.args.a, tt.args.b); !r.DeepEqual(got, tt.want) {
				t.Errorf("add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cross(t *testing.T) {
	type args struct {
		a *tuple
		b *tuple
	}
	tests := []struct {
		name string
		args args
		want *tuple
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cross(tt.args.a, tt.args.b); !r.DeepEqual(got, tt.want) {
				t.Errorf("cross() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dot(t *testing.T) {
	type args struct {
		a *tuple
		b *tuple
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dot(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("dot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_magnitude(t *testing.T) {
	type args struct {
		t *tuple
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := magnitude(tt.args.t); got != tt.want {
				t.Errorf("magnitude() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_magnitude_squared(t *testing.T) {
	type args struct {
		t *tuple
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := magnitude_squared(tt.args.t); got != tt.want {
				t.Errorf("magnitude_squared() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_normalize(t *testing.T) {
	type args struct {
		t *tuple
	}
	tests := []struct {
		name string
		args args
		want *tuple
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalize(tt.args.t); !r.DeepEqual(got, tt.want) {
				t.Errorf("normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tick(t *testing.T) {
	type args struct {
		e *env
		p *projectile
	}
	tests := []struct {
		name string
		args args
		want projectile
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tick(tt.args.e, tt.args.p); !r.DeepEqual(got, tt.want) {
				t.Errorf("tick() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_v(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{"1", args{"1"}, 1.0, false},
		{"1.1", args{"1.1"}, 1.1, false},
		{"-1.1", args{"-1.1"}, -1.1, false},
		// would be nice but it's not having this...		{ "root2/2", args{ "√2/2" }, 1.1, false },

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := v(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("v() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("v() got = %v, want %v", got, tt.want)
			}
		})
	}
}
