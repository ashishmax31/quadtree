package quadtree

import (
	"github.com/ashishmax31/quadtree/helpermethods"
)

var searcharea Circle
var results []QueryResult

// Quadtree ...The core quadtree data structure
type Quadtree struct {
	Rectangle
	Points   []Point
	Ne       *Quadtree
	Se       *Quadtree
	Sw       *Quadtree
	Nw       *Quadtree
	Divided  bool
	Capacity int
	Parent   *Quadtree
}

// Point ... Datatype to hold the coordinates of a point
type Point struct {
	x float64
	y float64
}

// Rectangle ...  Datatype to hold the boundry of the quadtree.
type Rectangle struct {
	topRight   Point
	bottomLeft Point
	w          float64
	h          float64
}

type Circle struct {
	Center Point
	Radius float64
}

type QueryResult struct {
	Pnt  Point
	Ind  int
	Quad *Quadtree
}

func (q *Quadtree) subdivide() {
	// Each quad/rectangle is represented by its topright point, buttomleft point and its width and height.

	// NE quad identifiers
	neTR := Point{q.topRight.x, q.topRight.y}
	neBL := Point{(q.topRight.x + q.bottomLeft.x) / 2, (q.topRight.y + q.bottomLeft.y) / 2}

	// SE quad identifiers
	seTR := Point{q.topRight.x, (q.topRight.y + q.bottomLeft.y) / 2}
	seBL := Point{(q.topRight.x + q.bottomLeft.x) / 2, q.bottomLeft.y}

	// SW quad identifiers
	swTR := Point{(q.topRight.x + q.bottomLeft.x) / 2, (q.topRight.y + q.bottomLeft.y) / 2}
	swBL := Point{q.bottomLeft.x, q.bottomLeft.y}

	// NW quad identifiers
	nwTR := Point{(q.topRight.x + q.bottomLeft.x) / 2, q.topRight.y}
	nwBL := Point{q.bottomLeft.x, (q.topRight.y + q.bottomLeft.y) / 2}

	// Different quadrants represented by rectangles after subdivision
	ne := Rectangle{neTR, neBL, q.w / 2, q.h / 2}
	se := Rectangle{seTR, seBL, q.w / 2, q.h / 2}
	sw := Rectangle{swTR, swBL, q.w / 2, q.h / 2}
	nw := Rectangle{nwTR, nwBL, q.w / 2, q.h / 2}
	q.Ne = &Quadtree{Rectangle: ne, Capacity: q.Capacity, Parent: q}
	q.Se = &Quadtree{Rectangle: se, Capacity: q.Capacity, Parent: q}
	q.Sw = &Quadtree{Rectangle: sw, Capacity: q.Capacity, Parent: q}
	q.Nw = &Quadtree{Rectangle: nw, Capacity: q.Capacity, Parent: q}

	// Mark the newly subdivided quadrant as divided.
	q.Divided = true
}

func (q Quadtree) pointInQuadrant(pnt Point) bool {
	return ((pnt.x > q.bottomLeft.x) && (pnt.y > q.bottomLeft.y) && (pnt.x <= q.topRight.x) && (pnt.y <= q.topRight.y))
}

func (rect Rectangle) contains(pnt Point) bool {
	return ((pnt.x > rect.bottomLeft.x) && (pnt.y > rect.bottomLeft.y) && (pnt.x <= rect.topRight.x) && (pnt.y <= rect.topRight.y))
}

func (q *Quadtree) getPoint(pnt Point) (quad *Quadtree, index int) {
	for i, point := range q.Points {
		if point == pnt {
			quad = q
			return quad, i
		}
	}
	return &Quadtree{}, -1
}

func (q Quadtree) overlapsCircle(c Circle) bool {
	nearestx := helpermethods.Max(q.bottomLeft.x, helpermethods.Less(c.Center.x, q.bottomLeft.x+q.w))
	nearesty := helpermethods.Max(q.bottomLeft.y, helpermethods.Less(c.Center.y, q.bottomLeft.y+q.h))
	distance := (nearestx-c.Center.x)*(nearestx-c.Center.x) + (nearesty-c.Center.y)*(nearesty-c.Center.y)
	return (distance <= c.Radius*c.Radius)
}

func (c Circle) contains(pnt Point) bool {
	distance := (pnt.x-c.Center.x)*(pnt.x-c.Center.x) + (pnt.y-c.Center.y)*(pnt.y-c.Center.y)
	return (distance) < (c.Radius * c.Radius)
}

func (q *Quadtree) query(area Circle, found *[]QueryResult) {
	if !q.overlapsCircle(area) {
		return
	}
	if len(q.Points) > 0 {
		for i, point := range q.Points {
			if area.contains(point) {
				*found = append(*found, QueryResult{point, i, q})
			}
		}
	}
	if q.Divided {
		q.Ne.query(area, found)
		q.Se.query(area, found)
		q.Sw.query(area, found)
		q.Nw.query(area, found)
	}
}

func resetQueryResults() {
	results = []QueryResult{}
}

func (q Quadtree) getPointsInQuadrant(pnt Point, res *[]Point) {
	if q.contains(pnt) {
		for _, point := range q.Points {
			*res = append(*res, point)
		}
		if q.Divided {
			q.Ne.getPointsInQuadrant(pnt, res)
			q.Se.getPointsInQuadrant(pnt, res)
			q.Sw.getPointsInQuadrant(pnt, res)
			q.Nw.getPointsInQuadrant(pnt, res)
		}
	} else {
		return
	}
}
