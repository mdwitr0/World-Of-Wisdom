package hashcash

import (
	"fmt"
	"main/internal/shared/helpers"
	"math/rand"
	"strconv"
	"strings"
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

func (s *Stamp) Hash2Stamps(hash string) error {
	splitHash := strings.Split(hash, ":")
	if len(splitHash) < 6 {
		return fmt.Errorf("invalid hash")
	}
	version, err := strconv.Atoi(splitHash[0])
	if err != nil {
		return fmt.Errorf("invalid version")
	}
	zerosCount, err := strconv.Atoi(splitHash[1])
	if err != nil {
		return fmt.Errorf("invalid zeros count")
	}
	date, err := strconv.ParseInt(splitHash[2], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid date")
	}
	counter, err := strconv.Atoi(splitHash[5])
	if err != nil {
		return fmt.Errorf("invalid counter")
	}

	s.Version = version
	s.ZerosCount = zerosCount
	s.Date = date
	s.Resource = splitHash[3]
	s.Rand = splitHash[4]
	s.Counter = counter
	return nil
}

func (s *Stamp) Stamp2Hash() string {
	return s.ToString()
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
