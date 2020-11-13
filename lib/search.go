// implements a parallel search over a set of edges with a given predicate to look for

package lib

import (
	"runtime"
	"sync"
)

type Search struct {
	H               *Graph
	Sp              []Special
	Edges           *Edges
	BalFactor       int
	Result          []int
	Generators      []*CombinationIterator
	ExhaustedSearch bool
}

type Predicate interface {
	Check(H *Graph, Sp []Special, sep *Edges, balancedFactor int) bool
}

func (s *Search) FindNext(pred Predicate) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	// log.Println("starting search")

	s.Result = []int{} // reset result

	var numProc = runtime.GOMAXPROCS(-1)

	var wg sync.WaitGroup
	wg.Add(numProc)
	finished := false
	// SEARCH:
	found := make(chan []int)
	wait := make(chan bool)
	//start workers
	for i := 0; i < numProc; i++ {
		go s.Worker(i, found, &wg, &finished, pred)
	}

	go func() {
		wg.Wait()
		wait <- true
	}()

	select {
	case s.Result = <-found:
		close(found) //to terminate other workers waiting on found
	case <-wait:
		s.ExhaustedSearch = true
	}

}

func (s Search) Worker(workernum int, found chan []int, wg *sync.WaitGroup, finished *bool, pred Predicate) {
	defer func() {
		if r := recover(); r != nil {
			// log.Printf("Worker %d 'forced' to quit, reason: %v", workernum, r)
			return
		}
	}()
	defer wg.Done()

	gen := s.Generators[workernum]

	for gen.HasNext() {
		if *finished {
			// log.Printf("Worker %d told to quit", workernum)
			return
		}
		j := gen.Combination

		sep := GetSubset(*s.Edges, j)
		if gen.BalSep || pred.Check(s.H, s.Sp, &sep, s.BalFactor) {
			gen.BalSep = true // cache result
			found <- j
			// log.Printf("Worker %d \" won \"", workernum)
			gen.Confirm()
			*finished = true
			return
		}
		gen.Confirm()
	}
}

type BalancedCheck struct {
}

func (b BalancedCheck) Check(H *Graph, Sp []Special, sep *Edges, balFactor int) bool {
	// log.Printf("Current considered sep %+v\n", sep)
	// log.Printf("Current present SP %+v\n", sp)

	//balancedness condition
	comps, compSps, _, _ := H.GetComponents(*sep, Sp)
	// log.Printf("Components of sep %+v\n", comps)

	balancednessLimit := (((H.Edges.Len() + len(Sp)) * (balFactor - 1)) / balFactor)

	for i := range comps {
		if comps[i].Edges.Len()+len(compSps[i]) > balancednessLimit {
			//log.Printf("Using %+v component %+v has weight %d instead of %d\n", sep,
			//        comps[i], comps[i].Edges.Len()+len(compSps[i]), ((g.Edges.Len() + len(sp)) / 2))
			return false
		}
	}

	// Make sure that "special seps can never be used as separators"
	for _, s := range Sp {
		if IntHash(s.Vertices) == IntHash(sep.Vertices()) {
			//log.Println("Special edge %+v\n used again", s)
			return false
		}
	}

	return true
}

type ParentCheck struct {
	Conn  []int
	Child []int
}

func (p ParentCheck) Check(H *Graph, Sp []Special, sep *Edges, balFactor int) bool {
	// log.Printf("Current considered sep %+v\n", sep)
	// log.Printf("Current present SP %+v\n", sp)

	//balancedness condition
	comps, compSps, _, _ := H.GetComponents(*sep, Sp)

	foundCompLow := false
	// var comp_low_index int
	var comp_low Graph
	var compSp_low []Special

	// log.Printf("Components of sep %+v\n", comps)

	balancednessLimit := (((H.Edges.Len() + len(Sp)) * (balFactor - 1)) / balFactor)

	for i := range comps {
		if comps[i].Edges.Len()+len(compSps[i]) > balancednessLimit {
			foundCompLow = true
			// comp_low_index = i //keep track of the index for composing comp_up later
			comp_low = comps[i]
			compSp_low = compSps[i]
		}
	}

	if !foundCompLow {
		return false
	}

	vertCompLow := append(comp_low.Vertices(), VerticesSpecial(compSp_low)...)
	childχ := Inter(p.Child, vertCompLow)

	if !Subset(Inter(vertCompLow, p.Conn), sep.Vertices()) {
		// log.Println("Conn not covered by parent")

		// log.Println("Conn: ", PrintVertices(Conn))
		// log.Println("V(parentλ) \\cap Conn", PrintVertices(Inter(parentλ.Vertices(), Conn)))
		// log.Println("V(Comp_low) \\cap Conn ", PrintVertices(Inter(vertCompLow, Conn)))
		return false
	}

	// Connectivity check
	if !Subset(Inter(vertCompLow, sep.Vertices()), childχ) {
		// log.Println("Child not connected to parent!")
		// log.Println("Parent lambda: ", PrintVertices(parentλ.Vertices()))
		// log.Println("Child lambda: ", PrintVertices(childλ.Vertices()))

		// log.Println("Child", childλ)

		return false
	}

	return true
}