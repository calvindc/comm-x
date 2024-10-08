package dijkstra

import (
	"testing"
	"time"
)

func TestGraph_AllShortestPath(t *testing.T) {
	//假定权重都一样，路径带方向
	w := 1 //w=0 将会导致不能找到所有有效路径,必须保证weight不能为0,最小是1　权重都+1过
	v := []*Vertex{
		{
			ID: 0,
			Arcs: map[int]int{
				1: w,
				3: w,
			},
		},
		{
			ID: 1,
			Arcs: map[int]int{
				0: w,
				2: w,
			},
		},
		{
			ID: 2,
			Arcs: map[int]int{
				1: w,
				3: w,
			},
		},
		{
			ID: 3,
			Arcs: map[int]int{
				0: w,
				2: w,
			},
		},
		{
			ID: 4,
			Arcs: map[int]int{
				2: w,
				3: w,
			},
		},
	}
	g := NewGraph(v)
	result := g.AllShortestPath(0, 2, DefaultCostGetter)
	/*
		result:=[[0,1,2],[0,3,2]]
	*/
	if len(result) != 2 {
		t.Errorf("shoude be two shortest path,result=%v", result)
	}
	t.Logf("result=%v", result)
}

func Benchmark_AllShortestPathMassVertices(b *testing.B) {
	numNodes := 10000
	v := make([]*Vertex, numNodes)
	for i := 0; i < numNodes-1; i++ {
		v[i] = &Vertex{
			ID: i,
			Arcs: map[int]int{
				i + 1: i + 1,
			},
		}
	}
	//最后一个节点,什么都不指向
	v[numNodes-1] = &Vertex{
		ID: numNodes - 1,
	}
	g := NewGraph(v)
	b.N = 20
	for i := 0; i < b.N; i++ {
		result := g.AllShortestPath(0, numNodes-1, DefaultCostGetter)
		if len(result) != 1 && len(result[0]) != numNodes-1 {
			b.Error("shoude has only one shortest path")
			return
		}
	}

}

func TestGraph_AllShortestPath3(b *testing.T) {
	numNodes := 10000
	starttime := time.Now()

	v := make([]*Vertex, numNodes)
	for i := 0; i < numNodes-1; i++ {
		v[i] = &Vertex{
			ID: i,
			Arcs: map[int]int{
				i + 1: 1,
			},
		}
	}
	//最后一个节点,什么都不指向
	v[numNodes-1] = &Vertex{
		ID: numNodes - 1,
	}
	g := NewGraph(v)

	result := g.AllShortestPath(0, numNodes-1, DefaultCostGetter)
	if len(result) != 1 && len(result[0]) != numNodes-1 {
		b.Error("shoude has only one shortest path")
		return
	}

	b.Logf("test cost %s", time.Since(starttime))
	b.Logf("result=%v", result)
}
