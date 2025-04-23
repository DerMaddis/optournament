package tournament

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTournamentDepth(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		Songs []Song
		Depth int
	}{
		{
			Songs: []Song{{"a"}, {"b"}},
			Depth: 0,
		},
		{
			Songs: []Song{{"a"}, {"b"}, {"c"}, {"d"}},
			Depth: 1,
		},
		{
			Songs: make([]Song, 8),
			Depth: 2,
		},
	}

	for _, test := range tests {
		t, err := New(test.Songs)
		assert.NoError(err)
		assert.Equal(test.Depth, t.Depth)
	}
}

func TestInvalidTournaments(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		Songs []Song
	}{
		{
			Songs: nil,
		},
		{
			Songs: []Song{},
		},
		{
			Songs: []Song{{"a"}},
		},
		{
			Songs: []Song{{"a"}, {"b"}, {"c"}},
		},
		{
			Songs: make([]Song, 17),
		},
	}

	for _, test := range tests {
		_, err := New(test.Songs)
		assert.Error(err)
	}
}

func TestTournament(t *testing.T) {
	assert := assert.New(t)

	songs := []Song{
		{"a"},
		{"b"},
		{"c"},
		{"d"},
		{"e"},
		{"f"},
		{"g"},
		{"h"},
		{"i"},
		{"j"},
		{"k"},
		{"l"},
		{"m"},
		{"n"},
		{"o"},
		{"p"},
	}

	tournament, err := New(songs)
	assert.NoError(err)

	// ---------------------------- Depth 0

	assert.Equal(8, len(tournament.Matchups[0]))
	for i, matchup := range tournament.Matchups[0] {
		assert.Equal(&songs[i*2], matchup.Song1)
		assert.Equal(&songs[i*2+1], matchup.Song2)
	}

	// ---------------------------- Depth 1

	// All Song1s win
	for range 8 {
		tournament.Submit(1, 0)
	}

	// {"a"}, W -
	// {"b"},   | New matchup
	// {"c"}, W -
	// {"d"},
	// {"e"}, W -
	// {"f"},  	| New matchup
	// {"g"}, W -
	// {"h"},
	// {"i"}, W -
	// {"j"},   | New matchup
	// {"k"}, W -
	// {"l"},
	// {"m"}, W -
	// {"n"},   | New matchup
	// {"o"}, W -
	// {"p"},

	assert.Equal(len(tournament.Matchups[1]), 4)
	for i, matchup := range tournament.Matchups[1] {
		assert.NotNil(matchup)
		switch i {
		case 0:
			assert.Equal("a", matchup.Song1.Url)
			assert.Equal("c", matchup.Song2.Url)
		case 1:
			assert.Equal("e", matchup.Song1.Url)
			assert.Equal("g", matchup.Song2.Url)
		case 2:
			assert.Equal("i", matchup.Song1.Url)
			assert.Equal("k", matchup.Song2.Url)
		case 3:
			assert.Equal("m", matchup.Song1.Url)
			assert.Equal("o", matchup.Song2.Url)
		}
	}

	// ------------------------ Depth 2
	// All Song2s win
	for range 4 {
		tournament.Submit(0, 1)
	}

	// {"a"}, W -
	// {"b"},   |
	// {"c"}, W - W -
	// {"d"},       |
	// {"e"}, W -   |
	// {"f"},  	|   |
	// {"g"}, W - W -
	// {"h"},
	// {"i"}, W -
	// {"j"},   |
	// {"k"}, W - W -
	// {"l"},       |
	// {"m"}, W -   |
	// {"n"},   |   |
	// {"o"}, W - W -
	// {"p"},
	assert.Equal(2, len(tournament.Matchups[2]))
	for i, matchup := range tournament.Matchups[2] {
		assert.NotNil(matchup)
		switch i {
		case 0:
			assert.Equal("c", matchup.Song1.Url)
			assert.Equal("g", matchup.Song2.Url)
		case 1:
			assert.Equal("k", matchup.Song1.Url)
			assert.Equal("o", matchup.Song2.Url)
		}
	}

	// ------------------- Depth 3
	tournament.Submit(1, 0)
	tournament.Submit(0, 1)

	// {"a"}, W -
	// {"b"},   |
	// {"c"}, W - W - W -
	// {"d"},       |   |
	// {"e"}, W -   |   |
	// {"f"},  	|   |   |
	// {"g"}, W - W -   |
	// {"h"},           |
	// {"i"}, W -       |
	// {"j"},   |       |
	// {"k"}, W - W -   |
	// {"l"},       |   |
	// {"m"}, W -   |   |
	// {"n"},   |   |   |
	// {"o"}, W - W - W -
	// {"p"},

	assert.Equal(1, len(tournament.Matchups[3]))
	assert.Equal("c", tournament.Matchups[3][0].Song1.Url)
	assert.Equal("o", tournament.Matchups[3][0].Song2.Url)
}
