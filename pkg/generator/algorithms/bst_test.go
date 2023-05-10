package algorithms

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func verifyPoints(t *testing.T, tree *Tree, points []int) {
	for k, v := range points {
		vert, err := tree.GetPoint(k)
		assert.Nil(t, err)
		assert.Equal(t, v, vert)
	}
}

func buildTree(nodes, degree int) (*Tree, []int) {
	degs := make([]int, nodes)
	points := make([]int, 0)
	for k := range degs {
		degs[k] = degree
		for j := 0; j < degree; j++ {
			points = append(points, k)
		}
	}
	tree := New(nodes, degs)
	return tree, points
}

//func buildTreeFromSpanning(degrees []int) (*Tree, []int) {
//
//}

func TestBuildBaic(t *testing.T) {
	tree, points := buildTree(10, 3)
	verifyPoints(t, tree, points)
}

func TestBasicRemoval(t *testing.T) {
	tree, points := buildTree(20, 4)

	order := []int{
		0, 10, 12, 16, 5, 5, 3, 2, 2}

	verifyPoints(t, tree, points)
	for _, v := range order {
		vert := points[v]
		points = removeIndex(v, points)
		tree.RemovePoint(vert)
		verifyPoints(t, tree, points)
	}
}

func TestBuildMultiple(t *testing.T) {
	testcases := [][2]int{
		{12, 5},
		{58, 31},
		{100, 98},
		{2, 1},
		{5, 2},
	}

	for k := range testcases {
		tc := testcases[k]
		t.Run(fmt.Sprintf("%d:%d", tc[0], tc[1]), func(t *testing.T) {
			tree, points := buildTree(tc[0], tc[1])
			verifyPoints(t, tree, points)
		})
	}
}

func TestFullRemoval(t *testing.T) {
	randomOrder := []int{
		293, 109, 19, 82, 83, 193, 121, 282, 162, 184,
		119, 38, 212, 165, 148, 21, 83, 173, 202, 97, 235,
		23, 102, 243, 221, 65, 176, 244, 203, 82, 41, 133, 20,
		154, 17, 206, 155, 191, 259, 137, 94, 173, 80, 54, 102,
		22, 22, 87, 125, 243, 134, 113, 156, 124, 214, 171, 121,
		200, 108, 167, 56, 21, 62, 8, 123, 14, 183, 196, 154, 89,
		33, 45, 110, 167, 9, 56, 94, 78, 145, 157, 96, 53, 23, 210,
		117, 154, 34, 35, 70, 194, 193, 103, 110, 86, 124, 165, 25,
		67, 78, 24, 127, 112, 6, 85, 154, 52, 71, 103, 2, 75, 2, 104,
		47, 57, 19, 148, 27, 109, 72, 154, 81, 7, 41, 52, 90, 147, 120,
		101, 148, 136, 29, 13, 88, 73, 70, 33, 69, 6, 156, 140, 116, 17,
		33, 109, 6, 16, 17, 42, 2, 45, 58, 113, 58, 13, 7, 121, 63, 8, 94,
		9, 131, 56, 75, 84, 3, 12, 120, 62, 36, 82, 126, 88, 118, 68, 41, 91,
		48, 20, 72, 75, 116, 40, 27, 106, 111, 3, 56, 64, 56, 51, 40,
		78, 63, 26, 17, 58, 90, 21, 2, 15, 53, 29, 10, 46, 61, 66, 65,
		25, 56, 67, 79, 79, 29, 1, 63, 53, 38, 12, 68, 9, 58, 75, 36, 10,
		74, 15, 10, 58, 49, 39, 67, 4, 25, 27, 33, 13, 61, 17, 23, 13, 19,
		47, 29, 33, 55, 25, 7, 27, 46, 34, 25, 34, 45, 23, 22, 31, 23, 25,
		2, 26, 11, 18, 37, 1, 10, 26, 21, 7, 4, 22, 21, 15, 24, 26, 0, 5,
		16, 3, 3, 5, 15, 11, 15, 15, 4, 12, 8, 6, 6, 8, 6, 5, 1, 6, 0, 2,
		3, 2, 1, 0}

	tree, degs := buildTree(100, 3)

	for _, v := range randomOrder {
		node := degs[v]
		degs = removeIndex(v, degs)
		tree.RemovePoint(node)
		verifyPoints(t, tree, degs)
	}
}
