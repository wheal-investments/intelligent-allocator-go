package main

import (
	"fmt"
	"sort"
	"wheal-investments-algorithm/funds"
	"wheal-investments-algorithm/ga"
)

func main() {
	//The size of the population
	sizeOfPopulation := 1000

	//The number of generations
	numGenerations := 1000

	//The number of eliete members of the population
	numElietes := 50

	//The proability of mutation or crossover
	probMutation := 0.5

	//The probabilitty of completely new member of population
	probNewChromosome := 0.05

	//Create a new population of the desired size
	population := ga.NewPopulation(sizeOfPopulation)

	//Store the inital fittest ever chromosome
	fittestEverChromosome := population.Fittest()

	//Main generation loop
	for generation := 0; generation < numGenerations; generation++ {

		//Create a new populatio
		var newPopulation ga.Population

		//Sort the population by fitness
		sort.Slice(population.Chromosomes,
			func(i, j int) bool {
				return population.Chromosomes[i].Fitness > population.Chromosomes[j].Fitness
			})

		//Get the elite population and always include the fittest ever chromosome
		elitePopulation := ga.Population{
			Chromosomes: append(population.Chromosomes[0:numElietes], fittestEverChromosome),
		}

		//Loop to populate new population
		for len(newPopulation.Chromosomes) <= sizeOfPopulation {

			//Select a random chromosome from the elite population
			chromosome := elitePopulation.SelectRoulette()

			//Generate a random number
			randomNumber := ga.Random().Float64()

			//If should mutate
			if randomNumber < probMutation {
				//Equal proability of each mutation/crossover type
				mutationRandom := ga.Random().Intn(4)

				switch mutationRandom {
				case 1:
					chromosome.MutateIncrement()
				case 2:
					chromosome.MutateSwap()
				case 3:
					chromosome = ga.SingleCrossover(chromosome, elitePopulation.SelectRoulette())
				case 4:
					chromosome = ga.MultipleCrossover(chromosome, elitePopulation.SelectRoulette())
				}
			}

			//If should genenerate entirely new chromosome
			if randomNumber > (1.0 - probNewChromosome) {
				chromosome = ga.GenerateRandomChromosome()
			}

			//Add the new chromosome to the new population
			newPopulation.Chromosomes = append(newPopulation.Chromosomes, chromosome)
		}

		//Calculate the total fitness of the new population
		newPopulation.CalculateFitness()

		//Get the fittest chromosome of the new population
		fittest := newPopulation.Fittest()

		//If the fittest ever chromosome
		if fittest.Fitness > fittestEverChromosome.Fitness {
			fittestEverChromosome = fittest
		}

		//Print the fittest chromosome to the screen
		generationText := fmt.Sprintf("Gen %04d:", generation)
		fmt.Println(generationText, parametersText(fittest.GetActualFundParameters()), fittest.Fitness)

		//Set the new population as the population for the next generation
		population = newPopulation
	}

	//Get the fittest chromosome of the population
	fittest := population.Fittest()

	fmt.Println("Answer:", allocationText(fittest.GetFundAllocationPercentage()))
}

//Returns the actual fund parameters in human readable form
func parametersText(fundParameters funds.FundParameters) string {
	var parametersText string
	for _, value := range fundParameters {
		parametersText += fmt.Sprintf("%.2f ", value)
	}
	return parametersText
}

//Returns the allocation in human readable form
func allocationText(fundAllocation ga.FundAllocation) string {
	var allocationsText string
	for _, value := range fundAllocation {
		allocationsText += fmt.Sprintf("%.2f ", value*100)
	}
	return allocationsText
}
