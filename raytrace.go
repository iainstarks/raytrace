package main

import (
	"fmt"
	"math"
	"strconv"
)

type tuple struct {
	x float64
	y float64
	z float64
	w float64
}

type env struct {
	gravity *tuple
	wind    *tuple
}

type projectile struct {
	position *tuple
	velocity *tuple
}

func main() {
	// we'll implement the little end of chapter program in here - so definitions of Tuples, Points, Vectors - and functions to work with them
	// should go in here too
	p := projectile{NewPoint(0, 1, 0), normalize(NewVector(1, 1, 0))}
	e := env{NewVector(0.0, -0.1, 0), NewVector(-0.01, 0, 0)}
	i := 0
	for p.position.y > 0.0 && i < 100 {
		p = tick(&e, &p)
		fmt.Printf("%d Pos %f, %f, %f \n", i, p.position.x, p.position.y, p.position.z)
		i++
	}
	fmt.Println("Done.")
}

func v(v string) (float64, error) {
	// √2/2
	return strconv.ParseFloat(v, 64)
}

func tick(e *env, p *projectile) projectile {
	newpos := add(p.position, p.velocity)
	newvector := add(e.gravity, e.wind)
	newvelocity := add(p.velocity, newvector)
	newprojectile := projectile{newpos, newvelocity}
	return newprojectile
}

func add(a *tuple, b *tuple) *tuple {
	t := tuple{a.x + b.x, a.y + b.y, a.z + b.z, a.w + b.w}
	return &t
}

func sub(a *tuple, b *tuple) *tuple {
	t := tuple{a.x - b.x, a.y - b.y, a.z - b.z, a.w - b.w}
	return &t
}

func multiply(a *tuple, m float64) *tuple {
	t := tuple{a.x * m, a.y * m, a.z * m, a.w * m}
	return &t
}

func hadamardProduct(a *tuple, b *tuple) *tuple {
	t := tuple{a.x * b.x, a.y * b.y, a.z * b.z, a.w * b.w}
	return &t
}

func NewTuple(x float64, y float64, z float64, w float64) *tuple {
	t := tuple{x, y, z, w}
	return &t
}

func NewPoint(x float64, y float64, z float64) *tuple {
	t := tuple{x, y, z, 1.0}
	return &t
}

func NewVector(x float64, y float64, z float64) *tuple {
	t := tuple{x, y, z, 0.0}
	return &t
}

func magnitude(t *tuple) float64 {
	return math.Sqrt(magnitude_squared(t))
}

func magnitude_squared(t *tuple) float64 {
	m := t.x*t.x + t.y*t.y + t.z*t.z + t.w*t.w
	return m
}

func normalize(t *tuple) *tuple {
	m := magnitude(t)
	n := NewTuple(t.x/m, t.y/m, t.z/m, t.w/m)
	return n
}

func cross(a *tuple, b *tuple) *tuple {
	return NewTuple(a.y*b.z-a.z*b.y,
		a.z*b.x-a.x*b.z,
		a.x*b.y-a.y*b.x,
		0.0) // cross applies only to vectors
}

func dot(a *tuple, b *tuple) float64 {
	return a.x*b.x +
		a.y*b.y +
		a.z*b.z +
		a.w*b.w
}

// relect((1,-1,0), ( 0,1,0)) should apparently be (1,1,0)
// sub would be (1,-2,0) so sub*2 would be 2,-4,0 and dot would be 0-1+0+0=-1 so we ought to end up with -2,4,0 in my book which i do
// i expect there is some precedence that I am mis implementing
// the book says:
// a - b * 2 * dot(a, b)
// and math.stackedchange.com says
// r = d - 2(d.n)n   but n has to be normalised - which sounds the same to me as the book, except the book doesn't mention normalisation
// r = 1,-1,0  - 2 * (-1) * (0,1,0)  = [1,-1,0]  - [0, -2, 0]  = [1, 1, 0] which is the correct answer!
func reflect(a *tuple, b *tuple) *tuple {
	t := sub(a, multiply(multiply(b, 2), dot(a, b)))
	return t
}
