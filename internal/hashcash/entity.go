package hashcash

import (
	"fmt"
	"main/internal/shared/helpers"
	"math/rand"
	"strconv"
	"time"
)

type Stamp struct {
	Version    int    `json:"version"`
	ZerosCount int    `json:"zerosCount"`
	Date       int64  `json:"date"`
	Resource   string `json:"resource"`
	Rand       string `json:"rand"`
	Counter    int    `json:"counter"`
}

func (s *Stamp) ToString() string {
	return fmt.Sprintf("%d:%d:%d:%s::%s:%d", s.Version, s.ZerosCount, s.Date, s.Resource, s.Rand, s.Counter)
}

func (s *Stamp) IsSolved() bool {
	hashString := helpers.Data2Sha1Hash(s.ToString())
	if s.ZerosCount > len(hashString) {
		return false
	}

	for _, ch := range hashString[:s.ZerosCount] {
		if ch != zeroByte {
			return false
		}
	}

	return true
}

func (s *Stamp) ComputeHash(maxIterations int) (Stamp, error) {
	for s.Counter <= maxIterations || maxIterations <= 0 {
		if s.IsSolved() {
			return *s, nil
		}
		s.Counter++
	}
	return *s, ErrorMaxIterationsExceeded
}

func (s *Stamp) IssueStamp(resource string, zerosCount int) *Stamp {
	randNum := rand.Intn(200000)
	randStr := strconv.Itoa(randNum)

	s.Version = 1
	s.ZerosCount = zerosCount
	s.Date = time.Now().Unix()
	s.Resource = resource
	s.Rand = randStr
	s.Counter = 0

	return s
}
