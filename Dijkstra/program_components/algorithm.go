package programcomponents

import (
	"fmt"
	"math"
	"slices"
	"time"
)

// Calculation
func isPointInsideCircle(pX int32, pY int32, cX int32, cY int32, cR float32) bool {
	distance := float64(((cX - pX) * (cX - pX)) + (cY-pY)*(cY-pY))
	distance = math.Sqrt(distance)
	return float32(distance) <= cR
}

func calculateDistance(node1 *Node, node2 *Node) float32 {
	distance := float64(((node1.x - node2.x) * (node1.x - node2.x)) + ((node1.y - node2.y) * (node1.y - node2.y)))
	return float32(math.Sqrt(distance))
}

func dijkstraAlgo(graph *Graph) {
	setupSlices := func(data *AlgoDataContainer) {
		for i := 0; i < data.numOfNodes; i++ {
			var distanceVal float32
			var dValue float32
			var fromValue int

			if i == data.currentNodeIndex {
				distanceVal = 0.0
				dValue = 0.0
				fromValue = 0
			} else {
				distanceVal = math.MaxFloat32
				dValue = math.MaxFloat32
				fromValue = math.MaxInt
			}

			(*data.distance)[i] = distanceVal
			(*data.d)[i] = dValue
			(*data.from)[i] = fromValue
		}

		graph.updateDataStrings(data.numOfNodes, *data.distance, *data.from, *data.d)
	}

	calculateConnectedNodeDistances := func(data *AlgoDataContainer) {
		connections := graph.connectionMap[data.currentNodeIndex]
		for i := 0; i < len(connections); i++ {
			toNodeIndex := connections[i].toNode.index
			(*data.d)[toNodeIndex] = connections[i].cost + (*data.distance)[data.currentNodeIndex]
		}

		graph.updateDataStrings(data.numOfNodes, *data.distance, *data.from, *data.d)
	}

	setNewDistancesIfSmaller := func(data *AlgoDataContainer) {
		for i := 0; i < data.numOfNodes; i++ {
			if (*data.d)[i] < (*data.distance)[i] && i != data.currentNodeIndex {
				(*data.distance)[i] = (*data.d)[i]
				(*data.from)[i] = data.currentNodeIndex
			}
		}
		graph.updateDataStrings(data.numOfNodes, *data.distance, *data.from, *data.d)
	}

	pickClosestNodeIndex := func(data *AlgoDataContainer) int {
		var minValue float32 = math.MaxFloat32
		var minIndex int = math.MaxInt

		for i := 0; i < len(*data.d); i++ {
			if !slices.Contains(*data.visited, i) && (*data.d)[i] < minValue {
				minValue = (*data.d)[i]
				minIndex = i
			}
		}
		return minIndex
	}

	colorResultInGraph := func(data *AlgoDataContainer) {
		cni := graph.destNode.index
		nodes := &(graph.nodes)
		conns := &(graph.connections)

		for graph.startNode.nodeType != PATHNODE {
			(*nodes)[cni].updateType(PATHNODE)
			fni := (*data.from)[cni]

			for i := 0; i < len(graph.connections); i++ {
				if (*conns)[i].fromNode.index == fni && (*conns)[i].toNode.index == cni {
					(*conns)[i].updateType(PATHCONN)
					break
				}
			}
			cni = fni
		}
	}

	isFinished := func(data *AlgoDataContainer) bool {
		numOfVisitableNodes := 0
		for i := 0; i < data.numOfNodes; i++ {
			if (*data.d)[i] < math.MaxFloat32 && !slices.Contains(*data.visited, i) {
				numOfVisitableNodes++
			}
		}

		return numOfVisitableNodes == 0
	}

	if graph.startNode != nil && graph.destNode != nil {
		data := newAlgoDataContainer(graph)
		setupSlices(&data)

		for {
			*data.visited = append(*data.visited, data.currentNodeIndex)
			graph.nodes[data.currentNodeIndex].updateType(CURRENTNODE)

			calculateConnectedNodeDistances(&data)
			setNewDistancesIfSmaller(&data)

			data.currentNodeIndex = pickClosestNodeIndex(&data)

			if isFinished(&data) {
				break
			}

			time.Sleep(1 * time.Second)
		}
		colorResultInGraph(&data)
	} else {
		fmt.Println("The graphs start or dest node is nil, select them!")
	}
	fmt.Println("Algo finished")
}
