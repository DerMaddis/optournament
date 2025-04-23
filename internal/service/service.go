package service

import (
	"fmt"
	"log/slog"
	"slices"
	"sync"

	"github.com/dermaddis/op_tournament/internal/errs"
	"github.com/dermaddis/op_tournament/internal/model/tournament"
	"github.com/rs/xid"
)

type tournamentWrapper struct {
	tournament *tournament.Tournament

	// All the users are in here and their vote flag is false by default.
	// Once they voted for a certain matchup, they are set to true and
	// reset once the matchup is done.
	usersAndVotes map[string]bool

	creatorId string
	started   bool

	inviteIds []string

	score1 int
	score2 int
}

type Service struct {
	log *slog.Logger

	tournaments   map[string]*tournamentWrapper
	tournamentsMu *sync.RWMutex
}

func New(log *slog.Logger) *Service {
	return &Service{
		log:           log,
		tournaments:   map[string]*tournamentWrapper{},
		tournamentsMu: &sync.RWMutex{},
	}
}

func (s *Service) NewTournament(songs []tournament.Song, userId string) (string, error) {
	tournament, err := tournament.New(songs)
	if err != nil {
		return "", fmt.Errorf("failed to create tournament: %w", err)
	}

	wrapper := tournamentWrapper{
		tournament:    &tournament,
		usersAndVotes: map[string]bool{userId: false},
		creatorId:     userId,
		started:       false,
		inviteIds:     nil,
	}

	for range 9 {
		wrapper.tournament.Submit(1, 0)
	}

	s.tournamentsMu.Lock()
	s.tournaments["d04apidug4q34mvd0urg"] = &wrapper
	s.tournamentsMu.Unlock()

	return tournament.Id, nil
}

func (s *Service) NewInviteId(tournamentId, userId string) (string, error) {
	s.tournamentsMu.Lock()
	defer s.tournamentsMu.Unlock()

	tWrapper, ok := s.tournaments[tournamentId]
	if !ok {
		return "", errs.ErrNotFound.Tournament
	}

	inviteId := xid.New()
	tWrapper.inviteIds = append(tWrapper.inviteIds, inviteId.String())

	return inviteId.String(), nil
}

func (s *Service) Tournament(tournamentId, userId string) (*tournament.Tournament, error) {
	s.tournamentsMu.RLock()
	defer s.tournamentsMu.RUnlock()

	tWrapper, ok := s.tournaments[tournamentId]
	if !ok {
		return nil, errs.ErrNotFound.Tournament
	}

	_, ok = tWrapper.usersAndVotes[userId]
	if !ok {
		return nil, errs.ErrNotFound.Tournament
	}

	return tWrapper.tournament, nil
}

func (s *Service) Join(tournamentId, inviteId, userId string) error {
	s.tournamentsMu.Lock()
	defer s.tournamentsMu.Unlock()

	tWrapper, ok := s.tournaments[tournamentId]
	if !ok {
		return errs.ErrNotFound.Tournament
	}

	if tWrapper.started {
		return errs.ErrAlreadyStarted
	}

	_, ok = tWrapper.usersAndVotes[userId]
	if ok {
		return errs.ErrAlreadyInTournament
	}

	if !slices.Contains(tWrapper.inviteIds, inviteId) {
		return errs.ErrNotFound.Invite
	}

	tWrapper.usersAndVotes[userId] = false

	return nil
}

func (s *Service) Start(tournamentId, userId string) error {
	s.tournamentsMu.Lock()
	defer s.tournamentsMu.Unlock()

	tWrapper, ok := s.tournaments[tournamentId]
	if !ok {
		return errs.ErrNotFound.Tournament
	}

	_, ok = tWrapper.usersAndVotes[userId]
	if !ok {
		return errs.ErrNotFound.Tournament
	}

	if tWrapper.creatorId != userId {
		return errs.ErrNoPermission.Start
	}

	tWrapper.started = true
	return nil
}

func (s *Service) Vote(tournamentId string, vote int, userId string) error {
	if vote != 1 && vote != 2 {
		return errs.ErrBadVote
	}

	s.tournamentsMu.Lock()
	defer s.tournamentsMu.Unlock()
	tWrapper, ok := s.tournaments[tournamentId]
	if !ok {
		return errs.ErrNotFound.Tournament

	}

	voted, ok := tWrapper.usersAndVotes[userId]
	if !ok {
		return errs.ErrNotFound.Tournament
	}
	if voted {
		return errs.ErrAlreadyVoted
	}

	tWrapper.usersAndVotes[userId] = true
	switch vote {
	case 1:
		tWrapper.score1++
	case 2:
		tWrapper.score2++
	}

	allVotesIn := tWrapper.score1+tWrapper.score2 == len(tWrapper.usersAndVotes)
	if allVotesIn {
		tWrapper.tournament.Submit(tWrapper.score1, tWrapper.score2)

		tWrapper.score1 = 0
		tWrapper.score2 = 0

		for userId := range tWrapper.usersAndVotes {
			tWrapper.usersAndVotes[userId] = false
		}
	}

	return nil
}

func (s *Service) ForceSubmit(tournamentId, userId string) error {
	s.tournamentsMu.Lock()
	defer s.tournamentsMu.Unlock()

	tWrapper, ok := s.tournaments[tournamentId]
	if !ok {
		return errs.ErrNotFound.Tournament
	}

	_, ok = tWrapper.usersAndVotes[userId]
	if !ok {
		return errs.ErrNotFound.Tournament
	}

	if tWrapper.creatorId != userId {
		return errs.ErrNoPermission.ForceSumbit
	}

	tWrapper.tournament.Submit(tWrapper.score1, tWrapper.score2)

	return nil
}
