package main

import (
	"os"
	"fmt"
	"math"
	"math/rand"
	"bufio"
	"strconv"
)

const EarthRadius = 6372.8

func Square(x float64) float64 {
	return x * x
}

func RadianFromDegrees(Degrees float64) float64 {
	return Degrees * 0.01745329251994329577
}

func ReferenceHarvensine(X0, Y0, X1, Y1, EarthRadius float64) float64 {

	lat1 := Y0
	lat2 := Y1
	lon1 := X0
	lon2 := X1

	dLat := RadianFromDegrees(lat2 - lat1)
	dLon := RadianFromDegrees(lon2 - lon1)
	lat1 = RadianFromDegrees(lat1)
	lat2 = RadianFromDegrees(lat2)
	
	a := Square(math.Sin(dLat/2)) + math.Cos(lat1)*math.Cos(lat2)*Square(math.Sin(dLon/2))
	c := 2.0 * math.Asin(math.Sqrt(a))

	result := EarthRadius * c
	return result
}

func RandomRange(r *rand.Rand, min, max float64) float64 {
	return min + r.Float64()*(max-min)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func uniformGenerator(seed int64, numberOfPoints int64, w1 *bufio.Writer, w2 *bufio.Writer) float64 {

	r := rand.New(rand.NewSource(seed))

	maxx := 180.0
	maxy := 90.0

	sum := 0.0
	mul := 1 / (float64)(numberOfPoints)

	w1.WriteString("{\n\"pairs\": [\n")
	sep := ","
	for i := 0; i < int(numberOfPoints); i++ {
		x0 := RandomRange(r, -maxx, maxx)
		x1 := RandomRange(r, -maxx, maxx)
		y0 := RandomRange(r, -maxy, maxy)
		y1 := RandomRange(r, -maxy, maxy)

		harversineDistance := ReferenceHarvensine(x0, y0, x1, y1, EarthRadius)
		sum += (mul * harversineDistance)
		w2.WriteString(fmt.Sprintf("%f\n", harversineDistance))

		if i == (int)(numberOfPoints - 1) {
			sep = ""
		}
		w1.WriteString(fmt.Sprintf("{\"x0\": %.16f,\"y0\": %.16f,\"x1\": %.16f,\"y1\": %.16f}%s\n",
                x0,
                y0,
                x1,
                y1,
                sep))
	}

	w1.WriteString("]\n}\n")	
	return sum
}

func clusterGenerator(seed int64, numberOfPoints int64, w1 *bufio.Writer, w2 *bufio.Writer) float64 {
	r := rand.New(rand.NewSource(seed))

	numberofClusters:= 1+(numberOfPoints / 64)
	currentCluster := numberofClusters
	centerx := 0.0
	centery := 0.0
	radiusx := 180.0
	radiusy := 90.0
	maxx := 180.0
	maxy := 90.0
	sum := 0.0
	mul := 1 / (float64)(numberOfPoints)

	sep := ","
	for i := 0; i < int(numberOfPoints); i++ {
		if currentCluster == 0 {
			centerx = RandomRange(r, -maxx, maxx)
			centery = RandomRange(r, -maxy, maxy)
			radiusx = RandomRange(r, 0, maxx)
			radiusy = RandomRange(r, 0, maxy)
			currentCluster = numberofClusters
		}
		x0 := RandomRange(r, math.Max(centerx-radiusx, -maxx), math.Min(centerx+radiusx, maxx))
		y0 := RandomRange(r, math.Max(centery-radiusy, -maxy), math.Min(centery+radiusy, maxy))
		x1 := RandomRange(r, math.Max(centerx-radiusx, -maxx), math.Min(centerx+radiusx, maxx))
		y1 := RandomRange(r, math.Max(centery-radiusy, -maxy), math.Min(centery+radiusy, maxy))

		harversineDistance := ReferenceHarvensine(x0, y0, x1, y1, EarthRadius)
		sum += (mul * harversineDistance)
		
		w2.WriteString(fmt.Sprintf("%f\n", harversineDistance))

		if i == (int)(numberOfPoints - 1) {
			sep = ""
		}
		w1.WriteString(fmt.Sprintf("{\"x0\": %.16f,\"y0\": %.16f,\"x1\": %.16f,\"y1\": %.16f}%s\n",
                x0,
                y0,
                x1,
                y1,
                sep))
	}

	w1.WriteString("]\n}\n")	
	return sum
}

func main() {
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) <3 {
		fmt.Printf("Usage: %s <uniform/cluster> <random seed> <number of pairs to generate>\n", os.Args[0])
		os.Exit(1)
	}

	method := argsWithoutProg[0]
	seed, _ := strconv.ParseInt(argsWithoutProg[1], 10, 64)
	numberOfPoints, _ := strconv.ParseInt(argsWithoutProg[2], 10, 64)
	sum := 0.0

	file1, err := os.Create("harversinePoints.json")	
	checkError(err)

	file2, err := os.Create("harversineValues.txt")
	checkError(err)

	w1 := bufio.NewWriter(file1)
	w2 := bufio.NewWriter(file2)
	defer w1.Flush()
	defer w2.Flush()

	if method == "uniform" {
		sum = uniformGenerator(seed, numberOfPoints, w1, w2)
	} else if method == "cluster" {
		sum = clusterGenerator(seed, numberOfPoints, w1, w2)
	} else {
		fmt.Printf("Usage: %s <uniform/cluster> <random seed> <number of pairs to generate>\n", os.Args[0])
		os.Exit(1)
	}

	fmt.Printf("Method: %s, Seed: %d, Number Of Points: %d, Sum: %f\n", method, seed, numberOfPoints, sum)
}