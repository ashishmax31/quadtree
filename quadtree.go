package quadtree

import (
	"fmt"
	"math"
)

func CreateQuadtree(threshold int, boundary Rectangle) Quadtree {
	return Quadtree{Rectangle: boundary, Capacity: threshold}
}

func CreateBoundary(topRightX, topRightY, bottomLeftX, bottomLeftY, width, height float64) Rectangle {

	// TODO: Calculate width and height from TR and BL points.

	// width, height := helperMethods.GetWidthAndHeight(TR, BL)
	return Rectangle{Point{topRightX, topRightY}, Point{bottomLeftX, bottomLeftY}, width, height}
}

func CreateNewPoint(x, y float64) Point {
	return Point{x, y}
}

func (q *Quadtree) Insert(pnt Point) bool {
	if !q.pointInQuadrant(pnt) {
		return false
	}
	if len(q.Points) < q.Capacity {
		q.Points = append(q.Points, pnt)
		return true
	}
	if !q.Divided {
		q.subdivide()
	}
	if q.Ne.Insert(pnt) {
		return true
	} else if q.Se.Insert(pnt) {
		return true
	} else if q.Sw.Insert(pnt) {
		return true
	} else {
		return q.Nw.Insert(pnt)
	}
}

func (q *Quadtree) RemoveFromQuadTree(index int) {
	q.Points = append(q.Points[:index], q.Points[index+1:]...)
}

func (q *Quadtree) Search(pnt Point) (quad *Quadtree, index int) {
	if q.pointInQuadrant(pnt) {
		quad, index = q.getPoint(pnt)
		if index != -1 {
			return quad, index
		}
	}
	if q.Ne != nil && q.Ne.pointInQuadrant(pnt) {
		quad, index = q.Ne.Search(pnt)
	}
	if q.Se != nil && q.Se.pointInQuadrant(pnt) {
		quad, index = q.Se.Search(pnt)
	}
	if q.Sw != nil && q.Sw.pointInQuadrant(pnt) {
		quad, index = q.Sw.Search(pnt)

	}
	if q.Nw != nil && q.Nw.pointInQuadrant(pnt) {
		quad, index = q.Nw.Search(pnt)
	}
	return quad, index
}

// TODO: maybe improve this approach of passing back the result slice.
func (q *Quadtree) QueryForNearestPoints(p Point, searchRadius float64) []QueryResult {
	// searcharea = Circle{}
	searcharea = Circle{p, searchRadius}
	defer resetQueryResults()
	for {
		// fmt.Printf("Search area: %+v  \n", searcharea)
		q.query(searcharea, &results)
		if len(results) > 0 {
			break
		} else {
			// Proportionally increase the search radius if no points are found.
			searcharea = Circle{p, searcharea.Radius + searcharea.Radius*50/100}
		}
	}
	return results
}

func CalculateDistance(p1, p2 Point) float64 {
	return math.Sqrt((p2.x-p1.x)*(p2.x-p1.x) + (p2.y-p1.y)*(p2.y-p1.y))
}

func PrintQuad(quad *Quadtree, tpe string) {
	fmt.Printf("address: %p,  type: %s:  , data:%+v\n", quad, tpe, quad)
	if quad.Divided {
		PrintQuad(quad.Ne, "ne")
		PrintQuad(quad.Se, "se")
		PrintQuad(quad.Sw, "sw")
		PrintQuad(quad.Nw, "nw")
	}
}
