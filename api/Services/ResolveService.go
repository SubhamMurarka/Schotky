package services

import (
	"fmt"
	"log"

	"github.com/SubhamMurarka/Schotky/Dynamo"
)

type ResolveService struct {
	Dyno Dynamo.DynamoDaxAPI
}

type ResolveServices interface {
	ResolveURL(ShortUrl string) (string, error)
}

func NewResolveServiceObj(d Dynamo.DynamoDaxAPI) ResolveServices {
	return &ResolveService{
		Dyno: d,
	}
}

func (s *ResolveService) ResolveURL(ShortUrl string) (string, error) {
	// Call SelectItem to retrieve the item associated with the ShortUrl
	result, err := s.Dyno.SelectItem(ShortUrl)
	if err != nil {
		log.Printf("Failed to retrieve item: %v", err)
		return "", err
	}

	// Check if the result item is nil (i.e., the item was not found)
	if result == nil || result.Item == nil {
		log.Printf("No item found for ShortUrl: %s", ShortUrl)
		return "", fmt.Errorf("ShortUrl %s not found", ShortUrl)
	}

	// Extract the LongURL from the result item
	longURLAttr, ok := result.Item["LongURL"]
	if !ok || longURLAttr.S == nil {
		log.Printf("LongURL attribute not found for ShortUrl: %s", ShortUrl)
		return "", fmt.Errorf("LongURL not found for ShortUrl %s", ShortUrl)
	}

	longURL := *longURLAttr.S
	return longURL, nil
}
