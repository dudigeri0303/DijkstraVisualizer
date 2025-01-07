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
		updateDataStrings(graph, data.numOfNodes, *data.distance, *data.from, *data.d)
	}

	calculateConnectedNodeDistances := func(data *AlgoDataContainer) {
		connections := graph.connectionMap[data.currentNodeIndex]
		for i := 0; i < len(connections); i++ {
			toNodeIndex := connections[i].toNode.index

			if !slices.Contains(*data.visited, toNodeIndex) {
				(*data.d)[toNodeIndex] = connections[i].cost + (*data.distance)[data.currentNodeIndex]
			}
		}
		updateDataStrings(graph, data.numOfNodes, *data.distance, *data.from, *data.d)
	}

	setNewDistancesIFSmaller := func(data *AlgoDataContainer) {
		for i := 0; i < len(*data.d); i++ {
			if (*data.d)[i] < (*data.distance)[i] && i != data.currentNodeIndex {
				(*data.distance)[i] = (*data.d)[i]
				(*data.from)[i] = data.currentNodeIndex
			}
		}
		updateDataStrings(graph, data.numOfNodes, *data.distance, *data.from, *data.d)
	}

	pickClosestNodeIndex := func(data *AlgoDataContainer) int {
		var minValue float32 = math.MaxFloat32
		var minIndex int = math.MaxInt

		for i := 0; i < len(*data.d); i++ {
			if (*data.d)[i] > 0 && (*data.d)[i] < minValue {
				minValue = (*data.d)[i]
				minIndex = i
			}
		}
		return minIndex
	}

	updateAndRestDataWithResults := func(data *AlgoDataContainer, minIndex int) {
		*data.visited = append(*data.visited, minIndex)
		data.currentNodeIndex = minIndex

		for i := 0; i < len(*data.d); i++ {
			if i != data.currentNodeIndex {
				(*data.d)[i] = math.MaxFloat32
			} else {
				(*data.d)[i] = 0
			}
		}
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

	if graph.startNode != nil && graph.destNode != nil {
		data := newAlgoDataContainer(graph)
		setupSlices(&data)

		for len(*data.visited) != data.numOfNodes {
			calculateConnectedNodeDistances(&data)
			setNewDistancesIFSmaller(&data)

			minIndex := pickClosestNodeIndex(&data)
			if minIndex == math.MaxInt {
				break
			}

			updateAndRestDataWithResults(&data, minIndex)
			time.Sleep(1 * time.Second)
		}
		colorResultInGraph(&data)
	} else {
		fmt.Println("The graphs start or dest node is nil, select them!")
	}
	fmt.Println("Algo finished")
}
