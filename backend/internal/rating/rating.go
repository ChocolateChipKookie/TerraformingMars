// Package rating computes player ratings from the chronological game history
// using the OpenSkill Plackett-Luce model (github.com/intinig/go-openskill).
//
// Ratings are derived, never stored: any game create/update/delete invalidates
// the cached snapshot, which is rebuilt on the next read.
package rating

import (
	"sync"

	openskill "github.com/intinig/go-openskill/rating"
	"github.com/intinig/go-openskill/types"
	"golang.org/x/sync/singleflight"
)

// Participant is one player's contribution to a single game.
type Participant struct {
	PlayerID    int
	TotalPoints int
}

// GameForRating is the minimal data needed to score one game.
type GameForRating struct {
	GameID       int
	Date         string
	Participants []Participant
}

// Sigma upper bound for a rating to count as "established". Tuned against real data:
// players with ≥5 games typically sit below this; sigma asymptotes near 3 for veterans.
const EstablishedSigmaThreshold = 6.0

// Entry is one snapshot in a player's rating history.
type Entry struct {
	GameID       int     `json:"game_id"`
	Date         string  `json:"date"`
	TotalPoints  int     `json:"total_points"`
	Mu           float64 `json:"mu"`
	Sigma        float64 `json:"sigma"`
	Ordinal      float64 `json:"ordinal"`
	DeltaOrdinal float64 `json:"delta_ordinal"`
	Established  bool    `json:"established"`
}

// Snapshot is the full computed timeline indexed by player ID.
type Snapshot struct {
	Timeline map[int][]Entry
}

// Current returns the latest entry for a player, or nil if they've never played.
func (s *Snapshot) Current(playerID int) *Entry {
	entries := s.Timeline[playerID]
	if len(entries) == 0 {
		return nil
	}
	return &entries[len(entries)-1]
}

// ForGame returns the entry recorded for each participant of a specific game.
func (s *Snapshot) ForGame(gameID int) map[int]Entry {
	out := map[int]Entry{}
	for playerID, entries := range s.Timeline {
		for _, e := range entries {
			if e.GameID == gameID {
				out[playerID] = e
				break
			}
		}
	}
	return out
}

// Compute replays games in the supplied order and returns the full timeline.
// Caller is responsible for sorting games chronologically (date asc, game_id asc).
// Games with fewer than two participants are skipped — OpenSkill needs at least two teams.
func Compute(games []GameForRating) *Snapshot {
	current := map[int]types.Rating{}
	timeline := map[int][]Entry{}

	for _, g := range games {
		if len(g.Participants) < 2 {
			continue
		}

		teams := make([]types.Team, len(g.Participants))
		scores := make([]int, len(g.Participants))
		prevOrdinals := make([]float64, len(g.Participants))

		for i, p := range g.Participants {
			r, ok := current[p.PlayerID]
			if !ok {
				r = openskill.New()
			}
			teams[i] = types.Team{r}
			scores[i] = p.TotalPoints
			prevOrdinals[i] = openskill.Ordinal(r)
		}

		updated := openskill.Rate(teams, &types.OpenSkillOptions{Score: scores})

		for i, t := range updated {
			pid := g.Participants[i].PlayerID
			r := t[0]
			ord := openskill.Ordinal(r)
			current[pid] = r
			timeline[pid] = append(timeline[pid], Entry{
				GameID:       g.GameID,
				Date:         g.Date,
				TotalPoints:  scores[i],
				Mu:           r.Mu,
				Sigma:        r.Sigma,
				Ordinal:      ord,
				DeltaOrdinal: ord - prevOrdinals[i],
				Established:  r.Sigma <= EstablishedSigmaThreshold,
			})
		}
	}

	return &Snapshot{Timeline: timeline}
}

// Service caches the latest snapshot and recomputes lazily on miss.
// Game mutations call Invalidate to bust the cache.
type Service struct {
	loadGames func() ([]GameForRating, error)

	mu       sync.RWMutex
	snapshot *Snapshot
	sf       singleflight.Group
}

func NewService(loadGames func() ([]GameForRating, error)) *Service {
	return &Service{loadGames: loadGames}
}

// Snapshot returns the cached timeline, computing it on first call after init/invalidation.
// Concurrent callers coalesce onto one computation via singleflight.
func (s *Service) Snapshot() (*Snapshot, error) {
	s.mu.RLock()
	snap := s.snapshot
	s.mu.RUnlock()
	if snap != nil {
		return snap, nil
	}

	v, err, _ := s.sf.Do("snapshot", func() (any, error) {
		s.mu.RLock()
		if s.snapshot != nil {
			cur := s.snapshot
			s.mu.RUnlock()
			return cur, nil
		}
		s.mu.RUnlock()

		games, err := s.loadGames()
		if err != nil {
			return nil, err
		}
		fresh := Compute(games)

		s.mu.Lock()
		s.snapshot = fresh
		s.mu.Unlock()
		return fresh, nil
	})
	if err != nil {
		return nil, err
	}
	return v.(*Snapshot), nil
}

// Invalidate clears the cached snapshot. Cheap; the next read recomputes.
func (s *Service) Invalidate() {
	s.mu.Lock()
	s.snapshot = nil
	s.mu.Unlock()
}
