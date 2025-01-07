package programcomponents

import (
	"fmt"
	"math"
	"strconv"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Enums
type NodeType int

const (
	BASICNODE NodeType = iota
	CURRENTNODE
	PATHNODE
)

type ConnectionType int

const (
	BASICCONN ConnectionType = iota
	CANDIDATECONN
	PATHCONN
)

// Visual components
// Node
type Node struct {
	x        int32
	y        int32
	radius   float32
	index    int
	nodeType NodeType
}

func newNode(x_val int32, y_val int32, index_value int) Node {
	return Node{
		x:        x_val,
		y:        y_val,
		radius:   20.0,
		index:    index_value,
		nodeType: BASICNODE,
	}
}

func (node *Node) moveWhileDragged(x int32, y int32, graph *Graph) {
	node.x = x
	node.y = y

	for i := 0; i < len(graph.connections); i++ {
		conn := (*graph).connections[i]
		if conn.fromNode == node || conn.toNode == node {
			conn.updateCost()
		}
	}
}

func (node *Node) draw() {
	var color rl.Color

	switch node.nodeType {
	case BASICNODE:
		color = rl.Black
	case CURRENTNODE:
		color = rl.Blue
	case PATHNODE:
		color = rl.Green
	}
	rl.DrawCircle(node.x, node.y, node.radius, color)
	rl.DrawText(strconv.Itoa(node.index), node.x, node.y, 10, rl.White)
}

func (node *Node) updateType(newType NodeType) {
	node.nodeType = newType
}

// Connection
type Connection struct {
	fromNode *Node
	toNode   *Node
	connType ConnectionType
	cost     float32
}

func newConnection(fromNode **Node, toNode **Node) Connection {
	return Connection{
		fromNode: *fromNode,
		toNode:   *toNode,
		connType: BASICCONN,
		cost:     calculateDistance(*fromNode, *toNode),
	}
}

func (conn *Connection) draw() {
	var color rl.Color

	switch conn.connType {
	case BASICCONN:
		color = rl.Black
	case CANDIDATECONN:
		color = rl.Blue
	case PATHCONN:
		color = rl.Green
	}

	midX := (conn.fromNode.x + conn.toNode.x) / 2
	midY := (conn.fromNode.y + conn.toNode.y) / 2

	rl.DrawLine(conn.fromNode.x, conn.fromNode.y, conn.toNode.x, conn.toNode.y, color)
	rl.DrawText(strconv.FormatFloat(float64(conn.cost), 'f', 2, 64), midX, midY, 15, rl.Black)
}

func (conn *Connection) updateType(newType ConnectionType) {
	conn.connType = newType
}

func (conn *Connection) updateCost() {
	conn.cost = calculateDistance(conn.fromNode, conn.toNode)
}

// Thread Safe string
type SafeString struct {
	mu    sync.Mutex
	value string
}

func (ss *SafeString) setValue(newVal string) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	ss.value = newVal
}

func (ss *SafeString) getValue() string {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return ss.value
}

func (ss *SafeString) appendValue(newVal string) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	if newVal != fmt.Sprintf("%.2f\n", math.MaxFloat32) && newVal != fmt.Sprintf("%d\n", math.MaxInt) {
		ss.value += newVal
	} else {
		ss.value += "?\n"
	}
}

// Graph
type Graph struct {
	nodes           []*Node
	connections     []*Connection
	connectionMap   map[int][]*Connection
	startNode       *Node
	destNode        *Node
	nodeStr         *SafeString
	distanceDataStr *SafeString
	fromDataStr     *SafeString
	dDataStr        *SafeString
}

func newGraph() Graph {
	return Graph{
		nodes:           []*Node{},
		connections:     []*Connection{},
		connectionMap:   make(map[int][]*Connection),
		startNode:       nil,
		destNode:        nil,
		nodeStr:         &SafeString{value: "Nodes" + "\n"},
		distanceDataStr: &SafeString{value: "Distance" + "\n"},
		fromDataStr:     &SafeString{value: "From" + "\n"},
		dDataStr:        &SafeString{value: "D" + "\n"},
	}
}

func updateDataStrings(graph *Graph, numOfNodes int, distance []float32, from []int, d []float32) {
	graph.distanceDataStr.setValue("Distance" + "\n")
	graph.fromDataStr.setValue("From" + "\n")
	graph.dDataStr.setValue("D" + "\n")

	for i := 0; i < numOfNodes; i++ {
		graph.distanceDataStr.appendValue(fmt.Sprintf("%.2f\n", distance[i]))
		graph.fromDataStr.appendValue(fmt.Sprintf("%d\n", from[i]))
		graph.dDataStr.appendValue(fmt.Sprintf("%.2f\n", d[i]))
	}
}

func drawGraphData(graph *Graph) {
	rl.DrawText(graph.nodeStr.getValue(), 10, 10, 20, rl.Black)
	rl.DrawText(graph.fromDataStr.getValue(), 90, 10, 20, rl.Black)
	rl.DrawText(graph.distanceDataStr.getValue(), 150, 10, 20, rl.Black)
	rl.DrawText(graph.dDataStr.getValue(), 250, 10, 20, rl.Black)
}

// Algo data container
type AlgoDataContainer struct {
	from             *[]int
	distance         *[]float32
	d                *[]float32
	visited          *[]int
	numOfNodes       int
	currentNodeIndex int
}

func newAlgoDataContainer(graph *Graph) AlgoDataContainer {
	numOfNodesVal := len(graph.nodes)

	fromSlice := make([]int, numOfNodesVal)
	distanceSlice := make([]float32, numOfNodesVal)
	dSlice := make([]float32, numOfNodesVal)
	visitedSlice := []int{}

	return AlgoDataContainer{
		from:             &fromSlice,
		distance:         &distanceSlice,
		d:                &dSlice,
		visited:          &visitedSlice,
		numOfNodes:       numOfNodesVal,
		currentNodeIndex: graph.startNode.index,
	}
}
