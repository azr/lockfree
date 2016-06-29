package skiplist

import (
	"sync"
	"sync/atomic"
)

//Header of a skip list
type Header struct {
	length                      uint32
	leftSentinel, rightSentinel *Node
}

//Node of a skip list
type Node struct {
	key         int
	nexts       []*Node
	marked      bool
	fullyLinked bool
	lock        sync.Mutex
}

func New() *Header {
	h := &Header{}
	h.Initialize()
	return h
}

// Initialize resets the list to a default empty state
func (h *Header) Initialize() {
	var left [maxlevel]*Node
	var right [1]*Node
	rightMost := &Node{
		nexts:       right[:],
		fullyLinked: true,
	}
	for i := range left {
		left[i] = rightMost
	}
	leftMost := &Node{
		key:         -1,
		nexts:       left[:],
		fullyLinked: true,
	}

	h.leftSentinel, h.rightSentinel = leftMost, rightMost
}

func (n *Node) lowerThan(v int) bool {
	if n.nexts[0] == nil {
		//special right sentinel case,
		//nothing can be higher than that
		return false
	}
	return n.key < v
}

//findNode searches for the closest node to v or the node containing it
//and stores the path to get there in preds/succs
//
//returns -1 if nothing was found
//returns the layer in preds/succs in wich the node could be found
func (h *Header) findNode(v int, preds, succs []*Node) (lFound int) {
	lFound = -1
	pred := h.leftSentinel
	for layer := maxlevel - 1; layer >= 0; layer-- {
		curr := pred.nexts[layer]
		for ; curr.lowerThan(v) && layer > 0; layer-- {
			pred = curr
			curr = pred.nexts[layer]
		}
		if lFound == -1 && v == curr.key {
			lFound = layer
		}
		preds[layer] = pred
		succs[layer] = curr
		if curr == nil {
			break
		}
	}
	return
}

//Add v into list into at a random level
//returns false if a node is already there
//returns true if it was added or if it was already in there
func (h *Header) Add(v int) bool {
	topLayer := generateLevel(maxlevel)
	var preds, succs [maxlevel]*Node
	for {
		lFound := h.findNode(v, preds[:], succs[:])
		if lFound != -1 { // node was found
			nodeFound := succs[lFound]
			if !nodeFound.marked {
				for !nodeFound.fullyLinked {
					//make sure everything is valid
				}
				//node already in there
				return false
			}
			//something is deleting that node
			//let's try again
			continue
		}
		highestLocked := -1

		var prevPred, pred, succ *Node
		valid := true
		for layer := 0; valid && layer < topLayer; layer++ {
			pred = preds[layer]
			succ = succs[layer]
			if pred != prevPred {
				if pred == nil {
					println(".")
				}
				pred.lock.Lock()
				highestLocked = layer
				prevPred = pred
			}
			valid = !pred.marked && !succ.marked && pred.nexts[layer] == succ
		}
		if !valid {
			continue
		}
		newNode := newNode(v, topLayer)
		for layer := 0; layer < topLayer; layer++ {
			newNode.nexts[layer] = succs[layer]
			preds[layer].nexts[layer] = newNode
		}
		newNode.fullyLinked = true
		unlock(preds, highestLocked)
		atomic.AddUint32(&h.length, 1)
		return true
	}
}

//Remove node containing v if any
//return false if a remove is already in progress on that node
func (h *Header) Remove(v int) bool {
	var nodeToDelete *Node
	isMarked := false
	topLayer := -1
	var preds, succs [maxlevel]*Node
	for {
		lFound := h.findNode(v, preds[:], succs[:])
		if !(isMarked || (lFound != -1 && succs[lFound].okToDelete(lFound))) {
			return false
		}
		if !isMarked {
			nodeToDelete = succs[lFound]
			topLayer = len(nodeToDelete.nexts) - 1
			nodeToDelete.lock.Lock()
			if nodeToDelete.marked {
				nodeToDelete.lock.Unlock()
				return false
			}
			nodeToDelete.marked = true
			isMarked = true
		}
		highestLocked := -1

		var prevPred, pred, succ *Node
		valid := true
		for layer := 0; valid && (layer <= topLayer); layer++ {
			pred = preds[layer]
			succ = succs[layer]
			if pred != prevPred {
				pred.lock.Lock()
				highestLocked = layer
				prevPred = pred
			}
			valid = !pred.marked && pred.nexts[layer] == succ
		}
		if !valid {
			continue
		}
		for layer := topLayer; layer >= 0; layer-- {
			preds[layer].nexts[layer] = nodeToDelete.nexts[layer]
		}
		nodeToDelete.lock.Unlock()
		unlock(preds, highestLocked)
		atomic.AddUint32(&h.length, ^uint32(0))
		return true
	}
}

func (n *Node) okToDelete(lFound int) bool {
	return n.fullyLinked && len(n.nexts) == lFound+1 && !n.marked
}

//Contains returns true if v was found in list
func (h *Header) Contains(v int) bool {
	var preds, succs [maxlevel]*Node
	lFound := h.findNode(v, preds[:], succs[:])
	return lFound != -1 && succs[lFound].fullyLinked && !succs[lFound].marked
}

func unlock(preds [maxlevel]*Node, highestLocked int) {
	for i := 0; i <= highestLocked; i++ {
		preds[i].lock.Unlock()
	}
}

//newNode instanciates a *Node with topLayer set right
// and a slice of `topLayer` sized nexts
func newNode(v, topLayer int) *Node {
	return &Node{
		key:   v,
		nexts: make([]*Node, topLayer),
	}
}

//Len returns the size of the list
func (h *Header) Len() int {
	return int(atomic.LoadUint32(&h.length))
}
