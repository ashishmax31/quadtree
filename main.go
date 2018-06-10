package main

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"time"
)

var insertcount int
var searchcount int
var inputs []point
var linearcount int
var min float64
var startingPointpoint = point{0.5, 0.5}

type quadtree struct {
	rectangle
	points   []point
	ne       *quadtree
	se       *quadtree
	sw       *quadtree
	nw       *quadtree
	divided  bool
	capacity int
	outer    bool
}

type point struct {
	x float64
	y float64
}

type rectangle struct {
	topRight   point
	bottomLeft point
	w          float64
	h          float64
}
type circle struct {
	center point
	radius float64
}

func createQuadTree(capacity int, boundary rectangle) quadtree {
	return quadtree{capacity: capacity, rectangle: boundary, outer: true}
}

var first = true

func (q *quadtree) insert(pnt point) bool {

	if !q.pointInQuadrant(pnt) {
		return false
	}
	if len(q.points) < q.capacity {
		q.points = append(q.points, pnt)
		return true
	} else {
		if !q.divided {
			q.subdivide()
		}
		if q.ne.insert(pnt) {
			return true
		} else if q.se.insert(pnt) {
			return true
		} else if q.sw.insert(pnt) {
			return true
		} else {
			return q.nw.insert(pnt)
		}
	}

}

func (q *quadtree) subdivide() {
	// Each rectangle is represented by its topright point, buttomleft point and its width and height.
	neTR := point{q.topRight.x, q.topRight.y}
	neBL := point{(q.topRight.x + q.bottomLeft.x) / 2.0, (q.topRight.y + q.bottomLeft.y) / 2.0}
	seTR := point{q.topRight.x, (q.topRight.y + q.bottomLeft.y) / 2.0}
	seBL := point{(q.topRight.x + q.bottomLeft.x) / 2.0, q.bottomLeft.y}
	swTR := point{(q.topRight.x + q.bottomLeft.x) / 2.0, (q.topRight.y + q.bottomLeft.y) / 2.0}
	swBL := point{q.bottomLeft.x, q.bottomLeft.y}
	nwTR := point{(q.topRight.x + q.bottomLeft.x) / 2.0, q.topRight.y}
	nwBL := point{q.bottomLeft.x, (q.topRight.y + q.bottomLeft.y) / 2.0}

	// Different quadrants represented by rectangles after subdivision
	ne := rectangle{neTR, neBL, q.w / 2, q.h / 2}
	se := rectangle{seTR, seBL, q.w / 2, q.h / 2}
	sw := rectangle{swTR, swBL, q.w / 2, q.h / 2}
	nw := rectangle{nwTR, nwBL, q.w / 2, q.h / 2}
	q.ne = &quadtree{rectangle: ne, capacity: q.capacity}
	q.se = &quadtree{rectangle: se, capacity: q.capacity}
	q.sw = &quadtree{rectangle: sw, capacity: q.capacity}
	q.nw = &quadtree{rectangle: nw, capacity: q.capacity}
	q.divided = true
}

func (q quadtree) pointInQuadrant(pnt point) bool {
	return ((pnt.x > q.bottomLeft.x) && (pnt.y > q.bottomLeft.y) && (pnt.x <= q.topRight.x) && (pnt.y <= q.topRight.y))
}

func (rect rectangle) contains(pnt point) bool {
	return ((pnt.x > rect.bottomLeft.x) && (pnt.y > rect.bottomLeft.y) && (pnt.x <= rect.topRight.x) && (pnt.y <= rect.topRight.y))
}
func findMinDistance(currentpos point) (float64, point) {
	var minIndex int
	min = calculateDistance(currentpos, inputs[0])
	for i, item := range inputs {
		distance := calculateDistance(currentpos, item)
		if distance < min {
			min = distance
			minIndex = i
		}
	}
	currentpos = inputs[minIndex]
	removeItem(minIndex)
	return min, currentpos
}

func removeItem(index int) {
	inputs = append(inputs[:index], inputs[index+1:]...)
}
func calculateDistance(p1, p2 point) float64 {
	return math.Sqrt((p2.x-p1.x)*(p2.x-p1.x) + (p2.y-p1.y)*(p2.y-p1.y))
}

func linearsearch(pnt point) {
	for i, point := range inputs {
		if point == pnt {
			fmt.Printf("inputs searching: %d\n", i)
			fmt.Printf("Found point through inputs search after %d lookups\n", i)
		}
	}
}

func printQuad(quad quadtree, tpe string) {
	fmt.Printf("%s: %+v\n", tpe, quad)
	if quad.divided {
		printQuad(*quad.ne, "ne")
		printQuad(*quad.se, "se")
		printQuad(*quad.sw, "sw")
		printQuad(*quad.nw, "nw")
	}
}

func (q *quadtree) search(pnt point) (quad *quadtree, index int) {
	if q.pointInQuadrant(pnt) {
		quad, index = q.getPoint(pnt)
		if index != -1 {
			return quad, index
		}
	}
	if q.ne != nil && q.ne.pointInQuadrant(pnt) {
		quad, index = q.ne.search(pnt)
	}
	if q.se != nil && q.se.pointInQuadrant(pnt) {
		quad, index = q.se.search(pnt)
	}
	if q.sw != nil && q.sw.pointInQuadrant(pnt) {
		quad, index = q.sw.search(pnt)

	}
	if q.nw != nil && q.nw.pointInQuadrant(pnt) {
		quad, index = q.nw.search(pnt)
	}
	return quad, index
}

func (q *quadtree) getPoint(pnt point) (quad *quadtree, index int) {
	for i, point := range q.points {
		if point == pnt {
			quad = q
			return quad, i
		}
	}
	return &quadtree{}, -1
}

func randFloats(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func (q quadtree) findMinDistanceInNode(pnt point) (float64, *quadtree, int) {
	for i := range q.points {
		distance := calculateDistance(pnt, q.points[i])
		if distance < minimum {
			minimum = distance
			minIndex = i
		}
	}
	return min, &q, minIndex

}

func (q quadtree) overlapsCircle(c circle) bool {
	nearestx := max(q.bottomLeft.x, less(c.center.x, q.bottomLeft.x+q.w))
	nearesty := max(q.bottomLeft.y, less(c.center.y, q.bottomLeft.y+q.h))
	distance := math.Abs(nearestx-c.center.x)*math.Abs(nearestx-c.center.x) + math.Abs(nearesty-c.center.y)*math.Abs(nearesty-c.center.y)
	return (distance <= c.radius*c.radius)
}

func (q quadtree) overlap(rect rectangle) bool {
	// tl1 := point{q.bottomLeft.x, q.topRight.y}
	// br1 := point{q.topRight.x, q.bottomLeft.y}
	// tl2 := point{rect.bottomLeft.x, rect.topRight.y}
	// br2 := point{rect.topRight.x, rect.bottomLeft.y}
	// if tl1.x > br2.x || tl2.x > br1.x {
	// 	return false
	// }
	// if tl1.y < br2.y || tl2.y < br1.y {
	// 	return false
	// }
	// return true
	// Check if one rect is above the other
	fmt.Printf("search area: %+v \n quad area: %+v \n", rect, q.rectangle)

	if (rect.topRight.y < q.bottomLeft.y) || (rect.bottomLeft.y > q.topRight.y) {
		return false
	}
	// check if one rect is on the right of the right most edge
	if (q.topRight.x < rect.bottomLeft.x) || (rect.topRight.x < q.bottomLeft.x) {
		return false
	}
	// // return q.bottomLeft.x < rect.topRight.x && q.topRight.x > rect.bottomLeft.x && q.topRight.y > rect.bottomLeft.y && q.bottomLeft.y < rect.topRight.y
	fmt.Println("overlap")
	return true
}

var minimum = 100.00
var minquad *quadtree
var minIndex int

func (q quadtree) traverseTree(pnt point) (float64, *quadtree, int) {
	minimum, minquad, minIndex = q.findMinDistanceInNode(pnt)
	if q.divided {
		if q.ne.pointInQuadrant(pnt) {
			q.ne.traverseTree(pnt)

		}
		if q.se.pointInQuadrant(pnt) {
			q.se.traverseTree(pnt)

		}
		if q.sw.pointInQuadrant(pnt) {
			q.sw.traverseTree(pnt)

		}
		if q.nw.pointInQuadrant(pnt) {
			q.nw.traverseTree(pnt)
		}

	}

	return minimum, minquad, minIndex
}

func readInputs(quad *quadtree) {
	fh, err := os.Open("inputs")
	defer fh.Close()
	if err != nil {
		panic(err)
	}
	var x, y float64
	for {
		_, err := fmt.Fscanf(fh, "%f %f\n", &x, &y)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic("Something went wrong")
		}
		// inputs = append(inputs, point{x, y})
		quad.insert(point{x, y})

	}
}

func readInputsfromfile() {
	fh, err := os.Open("inputs")
	defer fh.Close()
	if err != nil {
		panic(err)
	}
	var x, y float64
	for {
		_, err := fmt.Fscanf(fh, "%f %f\n", &x, &y)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic("Something went wrong")
		}
		inputs = append(inputs, point{x, y})
		// quad.insert(point{x, y})

	}

}

func insertInputsToQuadTree(quad quadtree) {
	for i := range inputs {
		quad.insert(inputs[i])
	}
}

func (c circle) contains(pnt point) bool {
	distance := math.Abs(pnt.x-c.center.x)*math.Abs(pnt.x-c.center.x) + math.Abs(pnt.y-c.center.y)*math.Abs(pnt.y-c.center.y)
	fmt.Printf("Point checking: %v  distance, radius square=", pnt)
	fmt.Println(distance, c.radius*c.radius)
	return (distance) < (c.radius * c.radius)
}

func (q *quadtree) removeFromQuadTree(index int) {
	q.points = append(q.points[:index], q.points[index+1:]...)
}

func (q *quadtree) query(area circle, found *[]point) {
	// fmt.Printf("%v\n", q)
	if !q.overlapsCircle(area) {
		return
	} else {
		if len(q.points) > 0 {
			fmt.Printf("points in overlapping quad %v \n", q.points)
			for _, point := range q.points {
				fmt.Printf("area %v\n", area)
				// fmt.Printf("area contains point: %v \n", area.contains(point))
				if area.contains(point) {
					fmt.Println("inserting point to array")
					*found = append(*found, point)
					// fmt.Printf("found: %v \n", found)
				}
			}
		}
		if q.divided {

			q.ne.query(area, found)
			q.se.query(area, found)
			q.sw.query(area, found)
			q.nw.query(area, found)
		}
	}
}
func (q quadtree) getPointsInQuadrant(pnt point, res *[]point) {
	if q.contains(pnt) {
		for _, point := range q.points {
			*res = append(*res, point)
		}
		if q.divided {
			q.ne.getPointsInQuadrant(pnt, res)
			q.se.getPointsInQuadrant(pnt, res)
			q.sw.getPointsInQuadrant(pnt, res)
			q.nw.getPointsInQuadrant(pnt, res)
		}
	} else {
		return
	}
}

func bruteForce() {
	var dist float64
	inputLen := len(inputs)
	fmt.Printf("%v \n", inputs)
	currentpos.x = 0.5
	currentpos.y = 0.5
	var eaten int
	for {
		var t float64
		t, currentpos = findMinDistance(currentpos)
		fmt.Printf("%v \n", currentpos)
		dist = dist + t
		// fmt.Printf("%v \n", currentpos)

		eaten++
		fmt.Printf("remaining: %v \n", len(inputs))
		if eaten == inputLen {
			break
		}

	}
	// fmt.Printf("test:%f\n", calculateDistance(point{0.5, 0.5}, point{0.6, 0.6}))
	fmt.Printf("traveled : %f \n", dist)
}
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func less(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

var rect = rectangle{point{1.0, 1.0}, point{0.0, 0.0}, 1, 1}
var quad = createQuadTree(10, rect)

// Read and store the inputs from the file to inputs[]
var currentpos point

func main() {
	rand.Seed(time.Now().UnixNano())
	// Base rectangle
	// readInputsfromfile()
	readInputs(&quad)
	// quad.insert(point{0.5, 0.5})
	// t, _ := quad.search(point{0.48950098, 0.56534211})
	// printQuad(*t, "main")
	// bruteForce()
	treeQuad()
	// quad.Printf("%v \n", quad)
	// q, index := quad.search(point{0.8, 0.8})
	// q.removeFromQuadTree(index)
	// printQuad(quad)

	// se := point{0.71690689, 0.30471351}
	// se1 := point{0.20123, 0.305}
	// // insertInputsToQuadTree(quad)

	// currentpos.x = 0.5
	// currentpos.y = 0.5
	// res := []point{}
	// var radius = 0.07
	// var searcharea circle
	// var eaten int
	// var totalDistance float64

	// // searchArea := rectangle{topRight: point{currentpos.x + area, currentpos.y + area}, bottomLeft: point{currentpos.x - area, currentpos.y - area}}
	// for {
	// 	// quad.getPointsInQuadrant(currentpos, &res)
	// 	// fmt.Printf("result : %v \n", res)
	// 	if eaten == 100 {
	// 		break
	// 	}
	// 	// 	// fmt.Printf("minimum iteration: %v \n", minimum)
	// 	// 	// fmt.Printf("area iteration: %v \n", area)
	// 	// 	// fmt.Printf("currentpos: %v \n", currentpos)
	// 	// 	// https: //stackoverflow.com/questions/1373035/how-do-i-scale-one-rectangle-to-the-maximum-size-possible-within-another-rectang
	// 	// 	// sAreaWidth := math.Abs(searchArea.topRight.x - searchArea.bottomLeft.x)
	// 	// 	// sAreaHeight := math.Abs(searchArea.topRight.y - searchArea.bottomLeft.y)
	// 	// 	// maxScale := less((1 / sAreaWidth), (1 / sAreaHeight))
	// 	// 	// maxtrX := 1 - searchArea.topRight.x
	// 	// 	// maxtrY := 1 - searchArea.topRight.y
	// 	// 	// maxblX := searchArea.bottomLeft.x
	// 	// 	// maxblY := searchArea.bottomLeft.y

	// 	// 	// searchArea.topRight.x = (searchArea.topRight.x + maxtrX/10)
	// 	// 	// searchArea.topRight.y = (searchArea.topRight.y + maxtrY/10)
	// 	// 	// searchArea.bottomLeft.x = (searchArea.bottomLeft.x - maxblX/10)
	// 	// 	// searchArea.bottomLeft.y = (searchArea.bottomLeft.x - maxblY/10)

	// 	// 	// searchArea = rectangle{topRight: point{(currentpos.x + maxtrX/10) % maxtrX, (currentpos.y + maxtrY/10)}, bottomLeft: point{currentpos.x - area, currentpos.y - area}}
	// 	searcharea = circle{center: currentpos, radius: radius}
	// 	quad.query(searcharea, &res)
	// 	// 	// quad.getPointsInQuadrant(currentpos, &res)
	// 	// 	// fmt.Printf("result : %v \n", res)
	// 	// 	if len(res) == 0 {
	// 	// 		area += 0.05
	// 	// 	} else {
	// 	if len(res) > 0 {
	// 		for i, point := range res {
	// 			dist := calculateDistance(currentpos, point)
	// 			if dist < minimum {
	// 				minimum = dist
	// 				minIndex = i
	// 			}
	// 		}
	// 		// 		// fmt.Println(minimum)
	// 		totalDistance += minimum
	// 		currentpos = res[minIndex]
	// 		q, ind := quad.search(res[minIndex])
	// 		// 		// printQuad(quad)
	// 		q.removeFromQuadTree(ind)
	// 		// 		// printQuad(quad)
	// 		eaten++
	// 		// 		area = 0.05
	// 		res = []point{}
	// 		minimum = 100
	// 		radius = 0.07
	// 		// 	}
	// 	} else {
	// 		radius = radius + 0.07
	// 	}

	// }

	// fmt.Printf("total distance %v \n", totalDistance)

	// currentpos := point{0.5, 0.5}
	// for i := 1; i < len(inputs); i++ {
	// 	min = calculateDistance(currentpos, quad.points[0])
	// 	nearestpoint := quad.traverseTree(currentPos)
	// 	currentpos = nearestpoint

	// }
	// var dist float64
	// /*  */var eaten int
	// inputLen := 6
	// for {
	// 	var t float64
	// 	fmt.Printf("%v \n", currentpos)
	// 	t, q, i := quad.traverseTree(currentpos)
	// 	currentpos = q.points[i]
	// 	dist = dist + t
	// 	q.removeFromQuadTree(i)
	// 	// fmt.Printf("%v \n", currentpos)
	// 	eaten++
	// 	// fmt.Printf("remaining: %v \n", len(inputs))
	// 	if eaten == inputLen {
	// 		break
	// 	}
	// }
	// quad.insert(se1)
	// fmt.Println("___________________________________________________________________________________")

	// fmt.Println("Searching point in:")
	// fmt.Println(se.x, se.y)
	// quad.search(se)
	// fmt.Println("___________________________________________________________________________________")
	// fmt.Println("___________________________________________________________________________________")
	// linearsearch(se)

}

func treeQuad() {
	currentpos.x = 0.5
	currentpos.y = 0.5
	res := []point{}
	var radius = 0.1
	var searcharea circle
	var eaten int
	var totalDistance float64
	// printQuad(quad, "main")

	// searchArea := rectangle{topRight: point{currentpos.x + area, currentpos.y + area}, bottomLeft: point{currentpos.x - area, currentpos.y - area}}
	for {
		// quad.getPointsInQuadrant(currentpos, &res)
		// fmt.Printf("result : %v \n", res)
		if eaten == 100000 {
			break
		}
		// 	// fmt.Printf("minimum iteration: %v \n", minimum)
		// 	// fmt.Printf("area iteration: %v \n", area)
		// 	// fmt.Printf("currentpos: %v \n", currentpos)
		// 	// https: //stackoverflow.com/questions/1373035/how-do-i-scale-one-rectangle-to-the-maximum-size-possible-within-another-rectang
		// 	// sAreaWidth := math.Abs(searchArea.topRight.x - searchArea.bottomLeft.x)
		// 	// sAreaHeight := math.Abs(searchArea.topRight.y - searchArea.bottomLeft.y)
		// 	// maxScale := less((1 / sAreaWidth), (1 / sAreaHeight))
		// 	// maxtrX := 1 - searchArea.topRight.x
		// 	// maxtrY := 1 - searchArea.topRight.y
		// 	// maxblX := searchArea.bottomLeft.x
		// 	// maxblY := searchArea.bottomLeft.y

		// 	// searchArea.topRight.x = (searchArea.topRight.x + maxtrX/10)
		// 	// searchArea.topRight.y = (searchArea.topRight.y + maxtrY/10)
		// 	// searchArea.bottomLeft.x = (searchArea.bottomLeft.x - maxblX/10)
		// 	// searchArea.bottomLeft.y = (searchArea.bottomLeft.x - maxblY/10)

		// 	// searchArea = rectangle{topRight: point{(currentpos.x + maxtrX/10) % maxtrX, (currentpos.y + maxtrY/10)}, bottomLeft: point{currentpos.x - area, currentpos.y - area}}
		searcharea = circle{center: currentpos, radius: radius}
		quad.query(searcharea, &res)
		// 	// quad.getPointsInQuadrant(currentpos, &res)
		// 	// fmt.Printf("result : %v \n", res)
		// 	if len(res) == 0 {
		// 		area += 0.05
		// 	} else {
		if len(res) > 0 {
			for i, point := range res {
				dist := calculateDistance(currentpos, point)
				if dist < minimum {
					minimum = dist
					minIndex = i
				}
			}
			// 		// fmt.Println(minimum)
			fmt.Fprintf(os.Stderr, "%v\n", currentpos)

			totalDistance += minimum
			currentpos = res[minIndex]
			q, ind := quad.search(res[minIndex])
			// 		// printQuad(quad)
			q.removeFromQuadTree(ind)
			// 		// printQuad(quad)
			eaten++
			// 		area = 0.05
			res = []point{}
			minimum = 100
			radius = 0.1
			// 	}
		} else {
			radius = radius + 0.05
		}

	}

	fmt.Printf("total distance %v \n", totalDistance)
}
