package matrix

import (
	"math"

	"github.com/Avash027/Intra-Page-Ranker/crawler"
)

type MarkovMatrix struct {
	UrlToIntMapper   map[string]int
	IntToUrlMapper   map[int]string
	Alpha            float64
	AdjacencyMatrix  [][]float64
	TransitionMatrix [][]float64
	Epsilon          float64
}

type MarkovMatrixOpts struct {
	Size    int
	Alpha   float64
	Epsilon float64
}

type RankData struct {
	Url    string
	Weight float64
	Rank   int
}

func CreateMarkovMatrix(opts MarkovMatrixOpts) *MarkovMatrix {
	alpha := opts.Alpha

	urlToIntMapper := make(map[string]int)
	intToUrlMapper := make(map[int]string)

	return &MarkovMatrix{
		UrlToIntMapper: urlToIntMapper,
		IntToUrlMapper: intToUrlMapper,
		Alpha:          alpha,
		Epsilon:        opts.Epsilon,
	}
}

func (m *MarkovMatrix) InitalizeLinks(links []crawler.Link) {

	urlIndex := make(map[string]int)
	for _, link := range links {

		if _, ok := urlIndex[link.From]; !ok {
			urlIndex[link.From] = len(urlIndex)
		}

		if _, ok := urlIndex[link.To]; !ok {
			urlIndex[link.To] = len(urlIndex)
		}

		m.UrlToIntMapper[link.From] = urlIndex[link.From]
		m.UrlToIntMapper[link.To] = urlIndex[link.To]

		m.IntToUrlMapper[urlIndex[link.From]] = link.From
		m.IntToUrlMapper[urlIndex[link.To]] = link.To

	}

	numOfLinks := len(urlIndex)

	m.AdjacencyMatrix = make([][]float64, numOfLinks+1)
	m.TransitionMatrix = make([][]float64, numOfLinks+1)

	for i := 0; i <= numOfLinks; i++ {
		m.AdjacencyMatrix[i] = make([]float64, numOfLinks+1)
		m.TransitionMatrix[i] = make([]float64, numOfLinks+1)

	}

	for _, link := range links {
		m.AdjacencyMatrix[urlIndex[link.From]][urlIndex[link.To]] = 1
	}

}

func (m *MarkovMatrix) CreateTransitionMatrix() {
	N := len(m.AdjacencyMatrix)
	alpha := m.Alpha

	for i, row := range m.AdjacencyMatrix {
		sum := 0.0
		count := 0
		for _, val := range row {
			if val == 1 {
				count++
			}
			sum += val
		}

		if count == 0 {
			for j := range row {
				m.TransitionMatrix[i][j] = 1 / float64(N)
			}
		} else {
			for j, val := range row {
				if val == 1 {
					m.TransitionMatrix[i][j] = 1 / (sum * 1.0)
				}
			}
		}
	}

	for i := range m.TransitionMatrix {
		for j := range m.TransitionMatrix[i] {
			m.TransitionMatrix[i][j] = (1-alpha)*m.TransitionMatrix[i][j] + alpha/float64(N)
		}
	}
}

func (m *MarkovMatrix) ComputePageRank() []RankData {
	startIndex := 0

	startState := m.AdjacencyMatrix[startIndex]

	for {

		newState := multiplyMatrix(startState, m.TransitionMatrix)
		steadyStateReached := true

		for i, val := range startState {
			if math.Abs(val-newState[0][i]) > m.Epsilon {
				steadyStateReached = false
			}
		}

		if steadyStateReached {
			startState = newState[0]
			break
		}
		startState = newState[0]
	}

	finalState := make([]RankData, len(startState))

	for i, val := range startState {
		finalState[i] = RankData{
			Weight: val,
			Url:    m.IntToUrlMapper[i],
		}
	}

	return finalState
}

func multiplyMatrix(mat1 []float64, mat2 [][]float64) [][]float64 {
	// Get the dimensions of the input matrices

	rows1 := 1
	cols1 := len(mat1)
	rows2 := len(mat2)
	cols2 := len(mat2[0])

	if cols1 != rows2 {

		return nil
	}

	result := make([][]float64, rows1)
	for i := 0; i < rows1; i++ {
		result[i] = make([]float64, cols2)
	}

	for i := 0; i < rows1; i++ {
		for j := 0; j < cols2; j++ {
			for k := 0; k < cols1; k++ {
				result[i][j] += mat1[k] * mat2[k][j]
			}
		}
	}

	return result
}
