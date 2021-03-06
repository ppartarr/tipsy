package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
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
	dataset := "phpbb"
	q := 1000
	ballSize := 3
	minPasswordLength := 6

	filename := buildFilename(q, ballSize, minPasswordLength, dataset)

	results := filepath.Join(checker, filename)

	// open results
	file, err := os.Open(results)
	if err != nil {
		log.Fatal(err.Error())
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
	if err := p.Save(4*vg.Inch, 4*vg.Inch, filepath.Join("plots", checker+"-"+dataset+".png")); err != nil {
		panic(err)
	}
}

func TestPlotChecker(t *testing.T) {
	// seed randomness
	rand.Seed(int64(0))

	checker := "always"
	datasets := []string{"rockyou", "phpbb", "muslim"}
	q := 1000
	ballSize := 3
	minPasswordLength := 6

	results := make(map[string]*Result, len(datasets))
	for _, dataset := range datasets {
		filename := buildFilename(q, ballSize, minPasswordLength, dataset)

		path := filepath.Join(checker, filename)

		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer file.Close()

		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal("unable to unmarshal json")
		}

		results[dataset] = &Result{}

		// get json
		err = json.Unmarshal(bytes, results[dataset])
		if err != nil {
			log.Fatal("unable to unmarshal json")
		}
	}

	// create plot
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	// customize plot
	p.Title.Text = "Security Loss of the Always Checker"
	p.X.Label.Text = "Number of guesses allowed"
	p.Y.Label.Text = "λᵍʳᵉᵉᵈʸᵩ − λᵩ"

	curves := make(map[string]plotter.XYs, len(datasets))
	for dataset, result := range results {
		// defender list
		defenderList := checkers.LoadFrequencyBlacklist(result.DefenderListFile, minPasswordLength)

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
		curves[dataset] = points
	}

	// add curves to theplot
	err = plotutil.AddLinePoints(p,
		"rockyou", curves["rockyou"],
		"phpbb", curves["phpbb"],
		"muslim match", curves["muslim"],
	)
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, filepath.Join("plots", checker+".png")); err != nil {
		panic(err)
	}
}

func TestPlotDataset(t *testing.T) {
	// seed randomness
	rand.Seed(int64(0))

	checkerz := []string{"always", "blacklist", "optimal"}
	dataset := "phpbb"
	q := 1000
	ballSize := 3
	minPasswordLength := 6

	results := make(map[string]*Result, len(checkerz))
	for _, checker := range checkerz {
		filename := buildFilename(q, ballSize, minPasswordLength, dataset)

		path := filepath.Join(checker, filename)

		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer file.Close()

		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal("unable to unmarshal json")
		}

		results[checker] = &Result{}

		// get json
		err = json.Unmarshal(bytes, results[checker])
		if err != nil {
			log.Fatal("unable to unmarshal json")
		}
	}

	// create plot
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	// customize plot
	p.Title.Text = "Security Loss using Muslim Match"
	p.X.Label.Text = "Number of guesses allowed"
	p.Y.Label.Text = "λᵍʳᵉᵉᵈʸᵩ − λᵩ"

	curves := make(map[string]plotter.XYs, len(checkerz))
	for checker, result := range results {
		// defender list
		defenderList := checkers.LoadFrequencyBlacklist(result.DefenderListFile, minPasswordLength)

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
		curves[checker] = points
	}

	// add curves to theplot
	err = plotutil.AddLinePoints(p,
		"always", curves["always"],
		"blacklist", curves["blacklist"],
		"optimal", curves["optimal"],
		// "muslim match", curves["muslim"],
	)
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, filepath.Join("plots", dataset+".png")); err != nil {
		panic(err)
	}
}
