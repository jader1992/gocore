package id

import "github.com/rs/xid"

type GocoreIDService struct {

}

func NewGocoreIDService(params ...interface{}) (interface{}, error) {
	return &GocoreIDService{}, nil
}

func (s *GocoreIDService) NewID() string {
	return xid.New().String()
}

