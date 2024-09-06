package services

import (
	"errors"
	"log"
	"sync/atomic"

	"github.com/SubhamMurarka/Schotky/Dynamo"
	helper "github.com/SubhamMurarka/Schotky/Helpers"
	zookeepercounter "github.com/SubhamMurarka/Schotky/ZookeeperCounter"
)

type ShortenService struct {
	Dyno  Dynamo.DynamoDaxAPI
	Zkr   zookeepercounter.ZooKeeperClient
	start int64
	end   int64
}

type ShortenServices interface {
	ShortUrl(LongUrl string) (string, error)
	IncrementCounter()
}

func NewShortenServiceObj(d Dynamo.DynamoDaxAPI, z zookeepercounter.ZooKeeperClient) ShortenServices {
	return &ShortenService{
		Dyno:  d,
		Zkr:   z,
		start: 0,
		end:   -1,
	}
}

func (s *ShortenService) ShortUrl(LongUrl string) (string, error) {
	//Get new range for the server
	if s.start > s.end {
		s.start, s.end = s.Zkr.GetNewRange()
	}

	if !helper.RemoveDomainError(LongUrl) {
		return "", errors.New("DOMAIN HAS ERROR")
	}

	LongUrl = helper.EnforceHTTP(LongUrl)

	id := helper.ConvertToBase62(s.start)

	s.IncrementCounter()

	_, err := s.Dyno.InsertItem(id, LongUrl)
	if err != nil {
		log.Fatal("Not saved to Dynamo", err)
	}

	return id, nil
}

func (s *ShortenService) IncrementCounter() {
	atomic.AddInt64(&s.start, 1)
}
