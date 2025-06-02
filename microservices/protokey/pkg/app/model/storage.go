package model

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	ValidKey = regexp.MustCompile(`^[a-zA-Z0-9_.\-@]{1,1000}$`)
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
	Type   CommandType      `json:"type"`
	Key    string           `json:"key,omitempty"`
	Value  string           `json:"value,omitempty"`
	Prefix string           `json:"prefix,omitempty"`
	RespCh chan interface{} `json:"-"`
	ErrCh  chan error       `json:"-"`
}

type Storage struct {
	storeCh    chan command
	file       *os.File
	mu         sync.Mutex
	pending    []command
	flushTimer *time.Ticker
}

func NewStorage() *Storage {
	s := &Storage{
		storeCh:    make(chan command),
		pending:    make([]command, 0, 100),
		flushTimer: time.NewTicker(1 * time.Second),
	}

	f, err := os.OpenFile("ProtoKey.data", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	s.file = f

	if err := s.loadFromFile(); err != nil {
		panic(err)
	}

	go s.flushWorker()
	return s
}

func (s *Storage) loadFromFile() error {
	f, err := os.Open("ProtoKey.data")
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	store := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		var cmd command
		if err := json.Unmarshal([]byte(line), &cmd); err != nil {
			continue
		}
		switch cmd.Type {
		case Set:
			store[cmd.Key] = cmd.Value
		}
	}

	go func() {
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
	}()

	return scanner.Err()
}

func (s *Storage) flushWorker() {
	for range s.flushTimer.C {
		s.mu.Lock()
		if len(s.pending) > 0 {
			var lines []string
			for _, cmd := range s.pending {
				if cmd.Type == Set {
					b, err := json.Marshal(cmd)
					if err == nil {
						lines = append(lines, string(b))
					}
				}
			}
			if len(lines) > 0 {
				data := strings.Join(lines, "\n") + "\n"
				if _, err := s.file.WriteString(data); err == nil {
					s.file.Sync()
					s.pending = s.pending[:0]
				}
			}
		}
		s.mu.Unlock()
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

	s.mu.Lock()
	s.pending = append(s.pending, cmd)
	s.mu.Unlock()

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
