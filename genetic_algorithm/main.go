package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

var no_routines = 10

var MutationRate = 0.05

// PopSize is the size of the population
var PopSize = 1000   // 5000
var generations = 10 //100
var aspiration = 100

// func RunGA(writer *csv.Writer) Organism {
func RunGA() Organism {
	// start := time.Now()
	rand.Seed(time.Now().UTC().UnixNano())

	population := createPopulation()

	generation := 0
	iterations_since_best_oragnism := 0
	var bestOragismFound Organism
	bestOragismFound.Fitness = math.MinInt64

	bestOrganism := Organism{}

	generation_best := make([]string, 0)

	for i := 0; i < generations; i++ {
		generation++
		bestOrganism = getBest(population)
		// fmt.Printf("best: %f\n", bestOrganism.Fitness)

		bestOrganism.Generation = i
		if bestOrganism.Fitness > bestOragismFound.Fitness {
			iterations_since_best_oragnism = 0
			bestOragismFound = bestOrganism
		} else {
			iterations_since_best_oragnism++
			if iterations_since_best_oragnism > aspiration {
				break
			}
		}
		generation_best = append(generation_best, fmt.Sprintf("%f", 1/bestOrganism.Fitness))

		maxFitness := bestOrganism.Fitness
		pool := createPool(population, maxFitness)
		population = naturalSelection(pool, population)
	}
	// fmt.Printf("%+4v\n", generation_best)

	// err := writer.Write(generation_best)
	// checkError("Cannot write to file", err)

	return bestOragismFound
}

func main() {
	// var err error

	fmt.Printf("Confirguration: Mutataion Rate[%0.3f]\tPopulation Size[%d]\tGenerations[%d]\tAspiration[%d]\tRoutines[%d]\n", MutationRate, PopSize, generations, aspiration, no_routines)
	fmt.Printf("%-10s\t%-10s\t%-10s\t%-10s\t%-16s\t%-10s\t%-10s\t%-10s\t%-10s\t%-10s\n", "One Seaters", "Two Seaters", "Four Seaters", "Staff", "Store to Seating", "Expand", "Profit", "Time Per Run", "Total Time", "Avg Generations")

	var best []Organism

	// file, err := os.Create(fmt.Sprintf("%s_%d_%f_results.csv", data_sets_cost[i], no_hubs, alpha))
	// checkError("Cannot create file", err)

	// writer := csv.NewWriter(file)
	primary_start_time := time.Now()
	for k := 0; k < no_routines; k++ {
		start := time.Now()
		// c := RunGA(writer)
		c := RunGA()
		elapsed := time.Since(start)
		c.ElapsedTime = elapsed
		best = append(best, c)
	}
	// writer.Flush()
	// file.Close()

	sort.Sort(OrganismVector(best))
	// average TNC
	// average_tnc := 0.0
	average_generations := 0.0
	for _, c := range best {
		// average_tnc += c.DNA.Cost
		average_generations += float64(c.Generation)
	}
	// average_tnc = average_tnc / float64(len(best))
	average_generations = average_generations / float64(len(best))
	fmt.Printf("%-10d\t%-10d\t%-10d\t%-10d\t%-16d\t%-10d\t%-10f\t%-10s\t%-10s\t%-10f\n",
		best[0].DNA.One_seaters,
		best[0].DNA.Two_seaters, best[0].DNA.Four_seaters, best[0].DNA.Staff,
		best[0].DNA.Store_to_seating, best[0].DNA.Expand, best[0].Fitness, best[0].ElapsedTime, time.Since(primary_start_time), average_generations)

}

//DNA
type SolutionDNA struct {
	One_seaters      int // 4 sq ft per seater
	Two_seaters      int // 8 sq ft per seater
	Four_seaters     int // 16 sq ft per seater
	Staff            int // 3 vs 4 employees
	Expand           int // 208 sq ft
	Store_to_seating int // 80 sq ft extra
}

// Organism for this genetic algorithm
type Organism struct {
	DNA         *SolutionDNA
	Fitness     float64 // normalized cost
	Generation  int
	ElapsedTime time.Duration
}

type OrganismVector []Organism

func (c OrganismVector) Len() int {
	return len(c)
}

func (c OrganismVector) Less(i, j int) bool {
	return c[i].Fitness > c[j].Fitness
}

func (c OrganismVector) Swap(i, j int) {
	c[j], c[i] = c[i], c[j]
}

// creates a Organism (DOES NOT REPAIR IT SO BE CAREFUL)
func createOrganism() (organism Organism) {

	rand.Seed(time.Now().UnixNano())

	organism = Organism{}
	organism.DNA = &SolutionDNA{}

	// generate random number of 1,2,and 4 seaters
	organism.DNA.One_seaters = rand.Intn(124) // max number of one seaters if all store + expand are only one seaters
	organism.DNA.Two_seaters = rand.Intn(62)  // max number of two seaters if all store + expand are only one seaters
	organism.DNA.Four_seaters = rand.Intn(31) // max number of four seaters if all store + expand are only one seaters
	organism.DNA.Staff = rand.Intn(5-3) + 3   // max number of four seaters if all store + expand are only one seaters

	/* actually let this be random also and do this in the repair function
	this will let u consider the case where keep the store but also expand*/
	organism.DNA.Store_to_seating = rand.Intn(2) // can only be binary [0,1]
	organism.DNA.Expand = rand.Intn(2)           // can only be binary [0,1]

	organism.calcFitness()

	return organism
}

// creates the initial population
func createPopulation() (population []Organism) {
	population = make([]Organism, PopSize)
	for i := 0; i < PopSize; i++ {
		org := createOrganism()
		org.Repair()
		org.calcFitness()
		population[i] = org
	}
	return population
}

// calculates the fitness of the Organism based on a linear regression model
func (d *Organism) calcFitness() {
	profit := 1293.93804311 +
		float64(d.DNA.Store_to_seating)*(-836.61388617) +
		math.Log1p(float64(d.DNA.One_seaters))*(1548.00232475) +
		math.Log1p(float64(d.DNA.Two_seaters))*(2102.60068773) +
		math.Log1p(float64(d.DNA.Four_seaters))*(1042.22007575) +
		float64(d.DNA.Staff)*(-306.67339854) +
		float64(d.DNA.Expand)*(-3205.69302007)

	// fmt.Printf("Profit %f\n", profit)
	d.Fitness = profit
	return
}

// create the breeding pool that creates the next generation
func createPool(population []Organism, maxFitness float64) (pool []Organism) {
	pool = make([]Organism, 0)
	// create a pool for next generation
	for i := 0; i < len(population); i++ {
		// population[i].calcFitness()
		num := int((population[i].Fitness / maxFitness) * 1000)
		// fmt.Printf("Fitness: %f |  max fitness: %f | num: %d\n", population[i].Fitness, maxFitness, num)
		for n := 0; n < num; n++ {
			pool = append(pool, population[i])
		}
	}
	return
}

func (d *Organism) Repair() {

	// all in sq ft
	solution_required_area := d.DNA.One_seaters*4 + d.DNA.Two_seaters*8 + d.DNA.Four_seaters*16
	current_shop_area := 208
	store_area := 80
	expansion_area := 208

	// if solution valid --> return with no changes
	if solution_required_area < (current_shop_area + d.DNA.Store_to_seating*store_area + d.DNA.Expand*expansion_area) {
		return
	}

	// fmt.Printf("solution area invalid: %d\n", solution_required_area)

	// if the solution area exceeds the total available space
	// keep randomly reducing seats until it is feasible
	for solution_required_area > (current_shop_area + store_area + expansion_area) {
		rndm_int := rand.Intn(3)
		if rndm_int == 0 && d.DNA.One_seaters > 1 {
			d.DNA.One_seaters--
			solution_required_area -= 4
		} else if rndm_int == 1 && d.DNA.Two_seaters > 1 {
			d.DNA.Two_seaters--
			solution_required_area -= 8
		} else if d.DNA.Four_seaters > 1 {
			d.DNA.Four_seaters--
			solution_required_area -= 16
		}
	}

	// now the required area is between current_shop_area and max area available

	if solution_required_area < current_shop_area {
		// this condition will never be entered since if the required area is less than current space avalale,
		//code will return before getting here
		return
	} else if solution_required_area < (current_shop_area + store_area) { // preference given to store area since there is no extra rent required
		d.DNA.Store_to_seating = 1
	} else if solution_required_area < (current_shop_area + expansion_area) { // the try to only expand keeping whatever store_to_seat was set
		d.DNA.Expand = 1
	} else { // if all fails set both to 1
		d.DNA.Store_to_seating = 1
		d.DNA.Expand = 1
	}

	// fmt.Printf("solution area after repair: %d\n", solution_required_area)

	return
}

// perform natural selection to create the next generation
func naturalSelection(pool []Organism, population []Organism) []Organism {
	next := make([]Organism, len(population))
	for i := 0; i < len(population); i++ {
		r1, r2 := rand.Intn(len(pool)), rand.Intn(len(pool))
		a := pool[r1]
		b := pool[r2]

		child := crossover(a, b)
		child.mutate()

		child.Repair()

		child.calcFitness()

		next[i] = child
	}
	return next
}

// crosses over 2 Organisms
func crossover(d1 Organism, d2 Organism) Organism {
	dna := SolutionDNA{}

	d1_array := []int{
		d1.DNA.One_seaters,
		d1.DNA.Two_seaters,
		d1.DNA.Four_seaters,
		d1.DNA.Staff,
		d1.DNA.Store_to_seating,
		d1.DNA.Expand,
	}

	d2_array := []int{
		d2.DNA.One_seaters,
		d2.DNA.Two_seaters,
		d2.DNA.Four_seaters,
		d2.DNA.Staff,
		d2.DNA.Store_to_seating,
		d2.DNA.Expand,
	}

	child_array := make([]int, 0)
	child := Organism{
		DNA:     &dna,
		Fitness: 0,
	}

	mid := rand.Intn(len(d1_array))
	for i := 0; i < len(d2_array); i++ {
		if i > mid {
			child_array = append(child_array, d1_array[i])
		} else {
			child_array = append(child_array, d2_array[i])
		}
	}

	child.DNA.One_seaters = child_array[0]
	child.DNA.Two_seaters = child_array[1]
	child.DNA.Four_seaters = child_array[2]
	child.DNA.Staff = child_array[3]
	child.DNA.Store_to_seating = child_array[4]
	child.DNA.Expand = child_array[5]

	return child
}

// mutate the Organism
func (d *Organism) mutate() {

	// for each solution factor randomly change it if random float less than the mutation rate

	if rand.Float64() < MutationRate {
		d.DNA.One_seaters = rand.Intn(124) // max number of one seaters if all store + expand are only one seaters
	}

	if rand.Float64() < MutationRate {
		d.DNA.Two_seaters = rand.Intn(62)
	}

	if rand.Float64() < MutationRate {
		d.DNA.Four_seaters = rand.Intn(31)
	}

	if rand.Float64() < MutationRate {
		d.DNA.Staff = rand.Intn(5-3) + 3
	}

	if rand.Float64() < MutationRate {
		d.DNA.Store_to_seating = rand.Intn(2)
	}

	if rand.Float64() < MutationRate {
		d.DNA.Expand = rand.Intn(2)
	}

	return
}

// Get the best organism
func getBest(population []Organism) Organism {
	best := 0.0
	index := 0
	for i := 0; i < len(population); i++ {
		if population[i].Fitness > best {
			index = i
			best = population[i].Fitness
		}
	}
	return population[index]
}
