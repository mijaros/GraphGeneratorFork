package algorithms

import "errors"

type treeNode struct {
	val     int
	width   int
	subTree int
	height  int
	balance int
	left    *treeNode
	right   *treeNode
	parent  *treeNode
}

var (
	ErrNoPoints     = errors.New("point set is empty")
	ErrPointNTooBig = errors.New("input is too large")
)

func (n *treeNode) traverse() {
	if n.parent == nil {
		return
	}

	n.parent.checkUp(n)
}

func (n *treeNode) updateHeight() {
	l, r := 0, 0
	subtree := 0
	if n.left != nil {
		l = n.left.height + 1
		subtree += n.left.subTree
	}
	if n.right != nil {
		r = n.right.height + 1
		subtree += n.right.subTree
	}
	max := l
	if r > l {
		max = r
	}
	n.height = max
	n.balance = l - r
	n.subTree = n.width + subtree
}

func (n *treeNode) updateWidth() {
	wl, wr := 0, 0
	if n.left != nil {
		wl = n.left.subTree
	}
	if n.right != nil {
		wr = n.right.subTree
	}
	n.subTree = wl + n.width + wr
	if n.parent != nil {
		n.parent.updateWidth()
	}
}

func (n *treeNode) leftChild(child *treeNode) bool {
	return n.left == child
}

func (n *treeNode) checkUp(child *treeNode) {
	if n.parent == nil {
		return
	}
	parBalance := n.parent.balance
	leftGrand := n.leftChild(child)
	leftCh := n.parent.leftChild(n)
	if parBalance > 1 || parBalance < -1 {
		switch {
		case leftGrand && leftCh:
			n.llRotation()
		case leftGrand && !leftCh:
			n.lrRotation()
		case !leftGrand && leftCh:
			n.rlRotation()
		case !leftGrand && !leftCh:
			n.rrRotation()
		}
	} else {
		n.parent.checkUp(n)
	}
}

func (t *Tree) GetPoint(n int) (int, error) {
	iter := t.root
	if iter.subTree <= n {
		return -1, ErrPointNTooBig
	}
	for iter != nil {
		leftBound := 0
		if iter.left != nil {
			leftBound = iter.left.subTree
		}
		midBound := leftBound + iter.width
		if n < leftBound {
			iter = iter.left
			continue
		}
		if n < midBound {
			return iter.val, nil
		}
		n -= midBound
		iter = iter.right
	}
	return -1, ErrNoPoints
}

func (t *Tree) RemovePoint(n int) {
	iter := t.root
	for iter != nil && iter.val != n {
		if n < iter.val {
			iter = iter.left
		} else {
			iter = iter.right
		}
	}
	if iter == nil {
		return
	}
	if iter.width > 0 {
		iter.width--
		iter.updateWidth()
		//if iter.width == 0 {
		//	t.removeNode(iter)
		//}
	}
}

func (n *treeNode) findSuccessor() *treeNode {
	if n.right == nil {
		return n.left
	}
	iter := n.right

	for iter.left != nil {
		iter = iter.left
	}
	return iter
}

func (n *treeNode) removeChild(ch *treeNode) {
	if n.left == ch {
		n.left = nil
	}
	if n.right == ch {
		n.right = nil
	}
}

func (n *treeNode) remove() {

}

func (n *treeNode) swap(o *treeNode) {
	n.subTree, o.subTree = o.subTree, n.subTree
	n.parent, o.parent = o.parent, n.parent
	if o.parent != nil {
		o.parent.updateChild(n, o)
	}
	o.left = n.left
	n.right, o.right = o.right, n.right

}

func (n *treeNode) decreaseHeight() {
	n.height--
	if n.left != nil {
		n.left.decreaseHeight()
	}
	if n.right != nil {
		n.right.decreaseHeight()
	}
}

func (t *Tree) removeNode(n *treeNode) {
	succ := n.findSuccessor()
	if succ == nil {
		if t.root == n {
			t.root = nil
		} else {
			n.parent.removeChild(n)
		}
		return
	}
	if succ == n.left {
		succ.decreaseHeight()
		succ.parent = n.parent
		if n.parent != nil {
			n.parent.updateChild(n, succ)
		}

	}
	if n == t.root {
		t.root = succ
	}
	succ.subTree = n.subTree
	parent := n.parent
	n.parent = succ.parent
	succ.parent = parent
	if succ.parent != nil {
		if succ.parent.leftChild(n) {
			succ.parent.left = succ
		} else {
			succ.parent.right = succ
		}
	}
	succ.left = n.left
	oldRight := n.right
	n.right = succ.right
	succ.right = oldRight
	if succ.left != nil {
		succ.left.parent = succ
	}
	if succ.right != nil {
		succ.right.parent = succ
	}
	if n.right != nil {
		n.right.parent = n
	}
	n.remove()
}

func (t *Tree) GetRank(n int) int {
	iter := t.root
	for iter != nil && iter.val != n {
		if n < iter.val {
			iter = iter.left
		} else {
			iter = iter.right
		}
	}
	if iter == nil {
		return 0
	}
	return iter.width
}

func (n *treeNode) updateChild(oldChild, newChild *treeNode) {
	if oldChild == n.left {
		n.left = newChild
	} else if oldChild == n.right {
		n.right = newChild
	}
}

func (n *treeNode) rrRotation() {
	left := n.left
	oldPar := n.parent
	par := oldPar.parent

	n.left = oldPar

	n.parent = par
	oldPar.right = left
	if left != nil {
		left.parent = oldPar
	}
	oldPar.parent = n
	oldPar.updateHeight()
	n.updateHeight()
	if par != nil {
		par.updateChild(oldPar, n)
		par.checkUp(n)
	}

}

func (n *treeNode) rlRotation() {
	left := n.left.right
	par := n.parent
	n.parent = n.left
	n.parent.parent = par
	n.parent.right = n
	n.left = left
	if left != nil {
		left.parent = n
	}
	n.updateHeight()
	n.parent.updateHeight()
	if par != nil {
		par.updateChild(n, n.parent)
		n.parent.rrRotation()
	}
}

func (n *treeNode) lrRotation() {
	right := n.right.left
	par := n.parent
	n.parent = n.right
	n.parent.left = n
	n.parent.parent = par
	n.right = right
	if right != nil {
		right.parent = n
	}
	n.updateHeight()
	n.parent.updateHeight()
	if par != nil {
		par.updateChild(n, n.parent)
		n.parent.llRotation()
	}

}

func (n *treeNode) llRotation() {
	right := n.right
	oldPar := n.parent
	par := oldPar.parent
	n.right = oldPar
	oldPar.left = right
	if right != nil {
		right.parent = oldPar
	}
	n.parent = par
	oldPar.parent = n
	oldPar.updateHeight()
	n.updateHeight()

	if par != nil {
		par.updateChild(oldPar, n)
		par.checkUp(n)
	}

}

func (n *treeNode) rebalance(child *treeNode) {
	n.updateHeight()

	if n.parent != nil {
		n.parent.rebalance(n)
	}
}

func (n *treeNode) insertNode(newNode *treeNode) {
	if n.val < newNode.val {
		if n.right == nil {
			n.right = newNode
			newNode.parent = n
			n.rebalance(newNode)
			n.traverse()
		} else {
			n.right.insertNode(newNode)
		}
	} else {
		if n.left == nil {
			n.left = newNode
			newNode.parent = n
			n.rebalance(newNode)
			n.traverse()
		} else {
			n.left.insertNode(newNode)
		}
	}
}

type Tree struct {
	root *treeNode
}

func (t *Tree) Length() int {
	if t.root == nil {
		return 0
	}
	return t.root.subTree
}

func New(number int, width []int) *Tree {
	nodes := make([]*treeNode, number)
	for k := range nodes {
		nodes[k] = &treeNode{
			val:     k,
			width:   width[k],
			subTree: width[k],
			balance: 0,
			left:    nil,
			right:   nil,
		}
	}
	root := nodes[0]
	root.subTree = 0
	nodes = nodes[1:]
	for len(nodes) > 0 {
		toAdd := nodes[0]
		nodes = nodes[1:]
		root.insertNode(toAdd)
		for root.parent != nil {
			root = root.parent
		}
	}

	tr := &Tree{root: root}
	return tr
}

func createPoints(nodes, deg int, currDeg []int) *Tree {
	exDeg := make([]int, nodes)
	for k := range exDeg {
		exDeg[k] = deg - currDeg[k]
	}

	return New(nodes, exDeg)
}
