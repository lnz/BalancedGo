package benchmark

import (
	"io/ioutil"
	"math/rand"
	"testing"

	lib "github.com/cem-okulmus/BalancedGo/lib"
	disj "github.com/cem-okulmus/BalancedGo/disj"
)


func setup(fname string, width int) (lib.Graph, lib.Edges) {
	dat, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}

	var G lib.Graph
	
	G, _ = lib.GetGraph(string(dat))

	var E_slice = G.Edges.Slice()

	var sep = make([]lib.Edge, width)
	sep = append(sep,  E_slice[rand.Intn(len(E_slice))])
	for i := 1; i < width; i++ {
		var re = E_slice[rand.Intn(len(E_slice))]
		for j := i-1; j>= 0; j -- {
			if re.Name == sep[j].Name {
				i = i-1
				break
			}
		}
		sep[i] = re
	}
	return G, lib.NewEdges(sep)
}

func BenchmarkGetCompKakuro2(b *testing.B) {
	G, sep := setup("Kakuro-hard-070-ext.xml.hg", 2)
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		G.GetComponents(sep)
	}
}

func BenchmarkGetCompFastKakuro2(b *testing.B) {
	G, sep := setup("Kakuro-hard-070-ext.xml.hg", 2)
	var vertices = make(map[int]*disj.Element, len(G.Vertices()))
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		G.GetComponents_fast(sep, vertices)
	}
}

func BenchmarkGetCompKakuro3(b *testing.B) {
	G, sep := setup("Kakuro-hard-070-ext.xml.hg", 3)
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		G.GetComponents(sep)
	}
}

func BenchmarkGetCompFastKakuro3(b *testing.B) {
	G, sep := setup("Kakuro-hard-070-ext.xml.hg", 2)
	var vertices = make(map[int]*disj.Element, len(G.Vertices()))
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		G.GetComponents_fast(sep, vertices)
	}
}


func BenchmarkGetCompNono2(b *testing.B) {
	G, sep := setup("Nonogram-180-table.xml.hg", 2)
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		G.GetComponents(sep)
	}
}

func BenchmarkGetCompFastNono2(b *testing.B) {
	G, sep := setup("Nonogram-180-table.xml.hg", 2)
	var vertices = make(map[int]*disj.Element, len(G.Vertices()))
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		G.GetComponents_fast(sep, vertices)
	}
}

func TestCompFast(t *testing.T) {
	dat, err := ioutil.ReadFile("Nonogram-180-table.xml.hg")
	if err != nil {
		panic(err)
	}

	var G lib.Graph
	
	G, _ = lib.GetGraph(string(dat))

	var E_slice = G.Edges.Slice()
	var e1 = E_slice[rand.Intn(len(E_slice))]

	var e2 = E_slice[rand.Intn(len(E_slice))]
	for ok := true; ok; ok = (e1.Name == e2.Name) {
		e2 = E_slice[rand.Intn(len(E_slice))]
	}

	var sep = lib.NewEdges([]lib.Edge{e1, e2})

	var vertices = make(map[int]*disj.Element, len(G.Vertices()))
	G.GetComponents_fast(sep, vertices);
}
