package programcomponents

import (
	"fmt"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Constants
const (
	windowWidth       int32   = 1100
	windowHeight      int32   = 900
	buttonWidth       float32 = 80.0
	buttonHeight      float32 = 50.0
	buttonDistance    float32 = 5.0
	nodeRadius        float32 = 20.0
	uiWidthOffset     float32 = 330.0
	buttonY           float32 = 0
	connValueFontSize int32   = 15
	nodeIndexFontSize int32   = 10
	dataTableFontSize int32   = 20
)

// Create and add components
func addNodeToGraph(nodes *[]*Node, connectionMap *map[int][]*Connection, index *int) {
	node := newNode(windowWidth/2, windowHeight/2, *index)
	*nodes = append(*nodes, &node)
	(*connectionMap)[node.index] = []*Connection{}
}

func addConnectionToGraph(conns *[]*Connection, connectionMap *map[int][]*Connection, fromNode *Node, toNode *Node) {
	conn := newConnection(&fromNode, &toNode)
	*conns = append(*conns, &conn)
	(*connectionMap)[fromNode.index] = append((*connectionMap)[fromNode.index], &conn)
}

// Handle user import
func handleInput(graph *Graph, draggedNode **Node, fromNodeToConn **Node, toNodeToConn **Node) {
	// Handle mouse input
	// Left click
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		for i := 0; i < len(graph.nodes); i++ {
			n := (graph.nodes)[i]
			if isPointInsideCircle(rl.GetMouseX(), rl.GetMouseY(), n.x, n.y, n.radius) {
				*draggedNode = n
				break
			}
		}
	} else if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
		*draggedNode = nil
	}

	// Right click
	if rl.IsMouseButtonPressed(rl.MouseButtonRight) {
		for i := 0; i < len(graph.nodes); i++ {
			n := (graph.nodes)[i]
			if isPointInsideCircle(rl.GetMouseX(), rl.GetMouseY(), n.x, n.y, n.radius) {
				if *fromNodeToConn == nil {
					*fromNodeToConn = n
					break
				} else {
					*toNodeToConn = n
					addConnectionToGraph(&graph.connections, &graph.connectionMap, *fromNodeToConn, *toNodeToConn)
					*fromNodeToConn = nil
					*toNodeToConn = nil
					break
				}
			}
		}
	}

	// Handle Keyboard input
	if rl.IsKeyPressed(rl.KeyA) {
		for i := 0; i < len(graph.nodes); i++ {
			n := (graph.nodes)[i]
			if isPointInsideCircle(rl.GetMouseX(), rl.GetMouseY(), n.x, n.y, n.radius) {
				n.nodeType = PATHNODE
				if graph.startNode == nil {
					graph.startNode = n
				} else if graph.startNode != nil && graph.destNode == nil {
					graph.destNode = n
				}
				break
			}
		}
	}
}

// Update
func update(draggedNode *Node, graph *Graph) {
	if draggedNode != nil {
		draggedNode.moveWhileDragged(rl.GetMouseX(), rl.GetMouseY(), graph)
	}
}

// Drawing methods
func drawState(graph *Graph) {
	for i := 0; i < len(graph.connections); i++ {
		graph.connections[i].draw()
	}

	for i := 0; i < len(graph.nodes); i++ {
		graph.nodes[i].draw()
	}
}

func drawAndHandleGui(graph *Graph, index *int) {
	drawIndex := 0

	if gui.Button(rl.NewRectangle((buttonWidth+10)*float32(drawIndex)+uiWidthOffset, buttonY, buttonWidth, buttonHeight), "Node") {
		addNodeToGraph(&graph.nodes, &graph.connectionMap, index)
		graph.nodeStr.appendValue(fmt.Sprintf("%d\n", *index))
		*index = *index + 1
	}
	drawIndex++

	if gui.Button(rl.NewRectangle((buttonWidth+10)*float32(drawIndex)+uiWidthOffset, buttonY, buttonWidth, buttonHeight), "Start") {
		go dijkstraAlgo(graph)
	}
	drawIndex++

	if gui.Button(rl.NewRectangle((buttonWidth+10)*float32(drawIndex)+uiWidthOffset, buttonY, buttonWidth, buttonHeight), "New Graph") {
		*graph = newGraph()
		*index = 0
	}
	drawIndex++
}

// Main program loop
func MainLoop() {
	rl.InitWindow(windowWidth, windowHeight, "Dijkstra")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	graph := newGraph()

	var draggedNode *Node = nil
	var fromNodeToConnect *Node = nil
	var toNodeToConnect *Node = nil

	val := 0
	var index *int = &val

	for !rl.WindowShouldClose() {
		// Handle input
		handleInput(&graph, &draggedNode, &fromNodeToConnect, &toNodeToConnect)

		// Update
		update(draggedNode, &graph)

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		drawAndHandleGui(&graph, index)
		drawState(&graph)
		graph.drawGraphData()
		rl.EndDrawing()
	}
}
