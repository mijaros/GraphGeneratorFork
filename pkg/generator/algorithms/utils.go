package algorithms

import (
	"github.com/gammazero/deque"
	mrand "math/rand"
)

// getNthElem returns nth found element from map.
// it is used for getting random member.
func getNthElem[T any](val int, in map[int]T) int {
	it := 0
	for k, _ := range in {
		if it < val {
			it++
			continue
		}
		return k
	}
	return -1
}

// extractNodeDegFromGraph extracts slice of all degrees in current graph
func extractNodeDegFromGraph(graph []map[int]bool) []int {
	result := make([]int, len(graph))
	for i := range graph {
		result[i] = len(graph[i])
	}
	return result
}

// extractComponents returns slice of all components of connectivity in the graph
func extractComponents(graph []map[int]bool) []map[int]map[int]bool {
	components := make([]map[int]map[int]bool, 0)
	indexes := make([]int, len(graph))
	for k := range graph {
		indexes[k] = k
	}

	for len(indexes) != 0 {
		head := indexes[0]
		found := make(map[int]bool)
		iter := deque.New[int]()
		iter.PushBack(head)
		for iter.Len() != 0 {
			elem := iter.PopFront()
			if found[elem] {
				continue
			}
			found[elem] = true
			for k, t := range graph[elem] {
				if !t || found[k] {
					continue
				}
				iter.PushBack(k)
			}
		}
		newComp := make(map[int]map[int]bool)
		for k, v := range found {
			if !v {
				continue
			}
			newComp[k] = graph[k]
		}
		components = append(components, newComp)
		newIndexes := make([]int, 0, len(indexes))
		for _, k := range indexes {
			if found[k] {
				continue
			}
			newIndexes = append(newIndexes, k)
		}
		indexes = newIndexes
	}
	return components
}

// removeIndex removes element at passed index from the slice and shakes slice
// to be of correct size.
func removeIndex[T any](index int, slice []T) []T {
	if index == len(slice)-1 {
		slice = slice[:index]
	} else {
		slice = append(slice[:index], slice[index+1:]...)
	}
	return slice
}

// makeRegularGraphConnected is implementation of switching algorithm to make k-regular graph connected
// it finds all components of passed graph and randomly switches edges between those to connect
// the components.
func makeRegularGraphConnected(graph []map[int]bool, rand *mrand.Rand) []map[int]bool {
	components := extractComponents(graph)
	for len(components) > 1 {
		first := components[0]
		second := components[1]
		components = components[2:]
		firstNode := getNthElem(rand.Intn(len(first)), first)
		first2Node := getNthElem(rand.Intn(len(first[firstNode])), first[firstNode])
		secondNode := getNthElem(rand.Intn(len(second)), second)
		second2Node := getNthElem(rand.Intn(len(second[secondNode])), second[secondNode])
		delete(first[firstNode], first2Node)
		delete(first[first2Node], firstNode)
		delete(second[secondNode], second2Node)
		delete(second[second2Node], secondNode)
		first[firstNode][second2Node] = true
		first[first2Node][secondNode] = true
		second[secondNode][first2Node] = true
		second[second2Node][firstNode] = true
		newComp := make(map[int]map[int]bool)
		for k, v := range first {
			newComp[k] = v
		}
		for k, v := range second {
			newComp[k] = v
		}
		components = append(components, newComp)

	}
	mainComp := components[0]
	res := make([]map[int]bool, len(mainComp))
	for k, v := range mainComp {
		res[k] = v
	}
	components = extractComponents(res)
	if len(components) == 1 {
		return res
	}
	return makeRegularGraphConnected(res, rand) //Tail recursion in case of failure
} //can theoretically happen in case of a very
//bad sequence random numbers - haven't been observed yet.

// pointSet is
type pointSet struct {
	points    map[int]int
	sums      [][2]int
	len       int
	numPoints int
}

func (p *pointSet) Length() int {
	return p.len
}

func (p *pointSet) GetRank(elem int) int {
	return p.points[elem]
}

func (p *pointSet) GetPoint(elem int) int {
	left, right := 0, len(p.sums)
	mid := (right - left) / 2

	for mid > 0 && !(elem >= p.sums[mid-1][0] && elem <= p.sums[mid][0]) {
		if elem >= p.sums[mid][0] {
			left = mid
		} else {
			right = mid
		}
		mid = (right-left)/2 + left
	}
	return p.sums[mid][1]
}

func (p *pointSet) RemovePoint(elem int) {
	p.len -= 1
	p.points[elem] -= 1
	remIndex := -1
	for k := range p.sums {
		if p.sums[k][1] >= elem {
			p.sums[k][0]--
			if elem == p.sums[k][1] && p.points[elem] == 0 {
				remIndex = k
			}
		}
	}
	if remIndex != -1 {
		p.sums = removeIndex(remIndex, p.sums)
	}
	if p.points[elem] == 0 {
		delete(p.points, elem)
		p.numPoints -= 1
	}
}

func createPointsOld(nodes, deg int, currDeg []int) pointSet {
	points := pointSet{
		points: make(map[int]int),
		len:    0,
		sums:   make([][2]int, nodes),
	}
	counter := 0
	var filledElems []int
	for k := 0; k < nodes; k++ {
		if currDeg[k] >= deg {
			filledElems = append(filledElems, k)
			continue
		}
		points.points[k] = deg - currDeg[k]
		counter += points.points[k]
		points.sums[k] = [2]int{counter, k}
	}
	for i, v := range filledElems {
		points.sums = removeIndex(v-i, points.sums)
	}
	points.len = counter
	return points
}
