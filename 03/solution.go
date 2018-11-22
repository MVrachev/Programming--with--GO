package main

import (
	"math"

	"github.com/fmi/go-homework/geom"
)

const EPSILON float64 = 0.0000001

type Triangle struct {
	A, B, C geom.Vector
}

func NewTriangle(a, b, c geom.Vector) Triangle {
	return Triangle{
		A: a,
		B: b,
		C: c,
	}
}

func createVectorFromTwoPoints(a, b geom.Vector) (ab geom.Vector) {
	ab.X = b.X - a.X
	ab.Y = b.Y - a.Y
	ab.Z = b.Z - a.Z
	return ab
}

func findT(n, p, D geom.Vector, d float64) float64 {
	return (d - geom.Dot(n, p)) / geom.Dot(n, D)
}

func findIntersectDot(ray geom.Ray, t float64) geom.Vector {
	return geom.Add(ray.Origin, geom.Mul(ray.Direction, t))
}

func floatEquals(a, b float64) bool {
	if (a-b) < EPSILON && (b-a) < EPSILON {
		return true
	}
	return false
}

func isDotInTriangle(tr Triangle, normalVector, q geom.Vector) bool {

	ba := createVectorFromTwoPoints(tr.B, tr.A)
	qa := createVectorFromTwoPoints(q, tr.A)

	equationRes := geom.Dot((geom.Cross(ba, qa)), normalVector)
	if !floatEquals(0, equationRes) && !(equationRes > 0) {
		return false
	}

	cb := createVectorFromTwoPoints(tr.C, tr.B)
	qb := createVectorFromTwoPoints(q, tr.B)
	equationRes = geom.Dot((geom.Cross(cb, qb)), normalVector)

	if !floatEquals(0, equationRes) && !(equationRes > 0) {
		return false
	}

	ac := createVectorFromTwoPoints(tr.A, tr.C)
	qc := createVectorFromTwoPoints(q, tr.C)

	equationRes = geom.Dot((geom.Cross(ac, qc)), normalVector)
	if !floatEquals(0, equationRes) && !(equationRes > 0) {
		return false
	}

	return true
}

func (tr Triangle) Intersect(ray geom.Ray) bool {

	ab := createVectorFromTwoPoints(tr.A, tr.B)
	ac := createVectorFromTwoPoints(tr.A, tr.C)
	normalVector := geom.Cross(ab, ac)

	d := geom.Dot(normalVector, tr.A)

	if floatEquals(0, geom.Dot(normalVector, ray.Direction)) {
		return false
	}

	t := findT(normalVector, ray.Origin, ray.Direction, d)

	if t < 0 {
		return false
	}
	intersectionPoint := findIntersectDot(ray, t)

	return isDotInTriangle(tr, normalVector, intersectionPoint)
}

// ------------------------------------ QUAD ------------------------------------

type Quad struct {
	A, B, C, D geom.Vector
}

func NewQuad(a, b, c, d geom.Vector) Quad {
	return Quad{
		A: a,
		B: b,
		C: c,
		D: d,
	}
}

func (q Quad) Intersect(ray geom.Ray) bool {

	trFirst := NewTriangle(q.A, q.B, q.C)
	trSecond := NewTriangle(q.A, q.D, q.C)

	if trFirst.Intersect(ray) || trSecond.Intersect(ray) {
		return true
	}

	return false
}

// ------------------------------------ SPHERE ------------------------------------

type Sphere struct {
	origin geom.Vector
	radius float64
}

func NewSphere(origin geom.Vector, r float64) Sphere {
	return Sphere{
		origin: origin,
		radius: r,
	}
}

/*
	Check if ray intersects with Sphere
*/
func (s Sphere) Intersect(ray geom.Ray) bool {

	Q := geom.Sub(ray.Origin, s.origin)
	a := geom.Dot(ray.Direction, ray.Direction) // should be = 1
	b := 2 * geom.Dot(ray.Direction, Q)         // 2 * U * Q
	c := geom.Dot(Q, Q) - s.radius*s.radius
	d := b*b - 4*a*c

	// then solutions are complex, so no intersections
	if d < 0 {
		return false
	} else {
		t1 := (-b - math.Sqrt(d)) / 2 * a * c
		t2 := (-b + math.Sqrt(d)) / 2 * a * c

		if t1 < 0 || t2 < 0 {
			return false
		}
	}
	return true
}

//
func main() {

}
