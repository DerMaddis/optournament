package tournament

import (
	"fmt"
	"math"

	"github.com/dermaddis/op_tournament/internal/errs"
	"github.com/dermaddis/op_tournament/internal/sliceutils"
)

type Song struct {
	Url string
}

type Tournament struct {
	Songs []Song

	// Depth is the number of rounds.
	// For example: 2 songs => 1 round (finale),
	// 4 songs => 2 rounds (semi+finale),
	// 6 songs =>
	Depth int

	matchups [][]*Matchup

	currentDepth   int
	currentMatchup int
}

type Matchup struct {
	Song1 *Song
	Song2 *Song

	// 0: no winner yet
	// 1: song 1
	// 2: song 2
	Winner int
}

func New(songs []Song) (Tournament, error) {
	n := len(songs)
	if n < 2 {
		return Tournament{}, errs.ErrBadSongCount
	}

	// n is not power of 2
	if n&(n-1) != 0 {
		return Tournament{}, errs.ErrBadSongCount
	}

	// Depth is 0-based so a finale-only match has a depth of 0
	depth := int(math.Log2(float64(n / 2)))
	fmt.Printf("n: %v\n", n)
	fmt.Printf("depth: %v\n", depth)

	matchups := make([][]*Matchup, depth+1)

	matchups[0] = make([]*Matchup, n/2)
	for i, pair := range sliceutils.Pairs(songs) {
		matchups[0][i] = &Matchup{
			Song1:  &pair.One,
			Song2:  &pair.Two,
			Winner: 0,
		}
	}

	for d := 1; d <= depth; d++ {
		matchups[d] = make([]*Matchup, n/(2<<d))
	}

	return Tournament{
		Songs:          songs,
		Depth:          depth,
		matchups:       matchups,
		currentDepth:   0,
		currentMatchup: 0,
	}, nil
}

func (t *Tournament) CurrentMatchup() Matchup {
	return *t.matchups[t.currentDepth][t.currentMatchup]
}

func (t *Tournament) Submit(song1Score int, song2Score int) (done bool, err error) {
	if song1Score == song2Score {
		return false, errs.ErrBadScore
	}

	var winner *Song = t.CurrentMatchup().Song1
	if song2Score > song1Score {
		winner = t.CurrentMatchup().Song2
	}

	if t.currentDepth != t.Depth {
		t.insertToNextRound(winner)

		t.currentMatchup++
		if t.currentMatchup == len(t.matchups[t.currentDepth]) {
			// depth complete => move to next depth
			t.currentDepth++
			t.currentMatchup = 0
		}
	} else {
		return true, nil
	}

	return false, nil
}

func (t *Tournament) insertToNextRound(song *Song) {
	if t.currentDepth == t.Depth {
		return
	}

	nextRoundMatchups := t.matchups[t.currentDepth+1]
	fmt.Printf("nextRoundMatchups: %v\n", nextRoundMatchups)

	for i, matchup := range nextRoundMatchups {
		if matchup == nil {
			// non-existing matchup => create it
			nextRoundMatchups[i] = &Matchup{
				Song1:  song,
				Song2:  nil,
				Winner: 0,
			}
			return
		}

		if matchup.Song2 == nil {
			// not filled up matchup => fill it
			matchup.Song2 = song
			return
		}
	}
}
