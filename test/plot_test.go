package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/ppartarr/tipsy/checkers"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// TestPlot generates the plots
func TestPlot(t *testing.T) {
	checker := "always"
	q := 1000
	ballSize := 3
	minPasswordLength := 6

	filename := strconv.Itoa(q) + "-" + strconv.Itoa(ballSize) + "-" + strconv.Itoa(minPasswordLength) + ".json"

	results := filepath.Join(checker, filename)

	// open results
	file, err := os.Open(results)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("unable to unmarshal json")
	}

	result := &Result{}

	// get json
	err = json.Unmarshal(bytes, result)
	if err != nil {
		log.Fatal("unable to unmarshal json")
	}

	// open attacker & defender list
	// attackerList := checkers.LoadFrequencyBlacklist(result.AttackerListFile, minPasswordLength)
	defenderList := checkers.LoadFrequencyBlacklist(result.DefenderListFile, minPasswordLength)

	// seed randomness
	rand.Seed(int64(0))

	// create plot
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	// customize plot
	p.Title.Text = "Security Loss"
	p.X.Label.Text = "Number of guesses allowed"
	p.Y.Label.Text = "λᵍʳᵉᵉᵈʸᵩ − λᵩ"

	// add points
	step := 100
	points := make(plotter.XYs, (q / step))
	for i := step; i < q; i = i + step {
		guesses := result.GuessList[:i+step]
		naiveGuesses := result.NaiveGuessList[:i+step]
		guessListBall := guessListBall(guesses, result.Correctors)
		lambdaQGreedy := ballProbability(guessListBall, defenderList)
		lambdaQ := ballProbability(naiveGuesses, defenderList)
		secloss := (lambdaQGreedy - lambdaQ)
		fmt.Println("lambda q greedy: ", lambdaQGreedy)
		fmt.Println("lambda q: ", lambdaQ)
		fmt.Println("secloss: ", secloss)
		points[i/step].X = float64(i)
		points[i/step].Y = secloss
	}

	err = plotutil.AddLinePoints(p,
		"Always", points,
		// "Second", randomPoints(15),
		// "Third", randomPoints(15),
	)
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, filepath.Join("plots", checker+".png")); err != nil {
		panic(err)
	}
}
