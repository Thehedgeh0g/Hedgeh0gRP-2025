package model

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ValidKey = regexp.MustCompile(`^[a-zA-Z0-9_.\-]{1,1000}$`)
)

type CommandType int

const (
	Set CommandType = iota
	Get
	Keys
)

var (
	ErrKeyNotFound    = errors.New("key not found")
	ErrUnknownCommand = errors.New("unknown command")
)

type command struct {
	Type   CommandType
	Key    string
	Value  string
	Prefix string
	RespCh chan interface{}
	ErrCh  chan error
}

type Storage struct {
	storeCh chan command
}

func NewStorage() *Storage {
	s := &Storage{
		storeCh: make(chan command),
	}
	go s.worker()
	return s
}

func (s *Storage) worker() {
	store := make(map[string]string)
	for cmd := range s.storeCh {
		switch cmd.Type {
		case Set:
			store[cmd.Key] = cmd.Value
			cmd.ErrCh <- nil
			cmd.RespCh <- struct{}{}

		case Get:
			val, ok := store[cmd.Key]
			if !ok {
				cmd.ErrCh <- ErrKeyNotFound
				cmd.RespCh <- ""
			} else {
				cmd.ErrCh <- nil
				cmd.RespCh <- val
			}

		case Keys:
			var result []string
			for k := range store {
				if strings.HasPrefix(k, cmd.Prefix) {
					result = append(result, k)
				}
			}
			cmd.ErrCh <- nil
			cmd.RespCh <- result

		default:
			cmd.ErrCh <- ErrUnknownCommand
			closeChannels(cmd)
		}
	}
}

func closeChannels(cmd command) {
	close(cmd.RespCh)
	close(cmd.ErrCh)
}

func (s *Storage) Set(key string, value string) error {
	respCh := make(chan interface{}, 1)
	errCh := make(chan error, 1)
	cmd := command{
		Type:   Set,
		Key:    key,
		Value:  value,
		RespCh: respCh,
		ErrCh:  errCh,
	}
	s.storeCh <- cmd
	<-respCh
	return <-errCh
}

func (s *Storage) Get(key string) (string, error) {
	respCh := make(chan interface{}, 1)
	errCh := make(chan error, 1)
	cmd := command{
		Type:   Get,
		Key:    key,
		RespCh: respCh,
		ErrCh:  errCh,
	}
	s.storeCh <- cmd
	val := <-respCh
	err := <-errCh
	return val.(string), err
}

func (s *Storage) Keys(prefix string) ([]string, error) {
	respCh := make(chan interface{}, 1)
	errCh := make(chan error, 1)
	cmd := command{
		Type:   Keys,
		Prefix: prefix,
		RespCh: respCh,
		ErrCh:  errCh,
	}
	s.storeCh <- cmd
	val := <-respCh
	err := <-errCh
	return val.([]string), err
}
