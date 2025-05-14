package service

import (
	"errors"
	appevent "rankcalculator/pkg/app/event"
	"testing"

	"rankcalculator/pkg/app/model"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Моки
type mockText struct {
	mock.Mock
	text       string
	rank       float64
	similarity bool
	hash       string
}

func (m *mockText) SetSimilarity(similarity bool) { m.similarity = similarity }
func (m *mockText) GetText() string               { return m.text }
func (m *mockText) SetRank(r float64)             { m.rank = r }
func (m *mockText) GetRank() float64              { return m.rank }
func (m *mockText) GetSimilarity() bool           { return m.similarity }
func (m *mockText) GetHash() string               { return m.hash }

type mockRepo struct{ mock.Mock }

func (m *mockRepo) FindByHash(hash string) (model.Text, error) {
	args := m.Called(hash)
	return args.Get(0).(model.Text), args.Error(1)
}

func (m *mockRepo) Store(text model.Text) error {
	return m.Called(text).Error(0)
}

type mockDispatcher struct{ mock.Mock }

func (m *mockDispatcher) Dispatch(event appevent.Event) error {
	return m.Called(event).Error(0)
}

type mockCentrifugo struct{ mock.Mock }

func (m *mockCentrifugo) Publish(channel string, data interface{}) error {
	return m.Called(channel, data).Error(0)
}

// Тесты
func TestCalculateRank_Base(t *testing.T) {
	tests := []struct {
		name         string
		inputText    string
		expectedRank float64
	}{
		{
			name:         "Only letters (AБВГABCD)",
			inputText:    "AБВГABCD",
			expectedRank: 1.0, // 4/4
		},
		{
			name:         "Letters and punctuation (Hi!)",
			inputText:    "Hi!",
			expectedRank: 2.0 / 3.0,
		},
		{
			name:         "Only punctuation (?!@)",
			inputText:    "?!@",
			expectedRank: 0.0,
		},
		{
			name:         "Empty string",
			inputText:    "",
			expectedRank: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mockRepo)
			dispatcher := new(mockDispatcher)
			centrifugo := new(mockCentrifugo)

			text := &mockText{
				text: tt.inputText,
				hash: "testhash",
			}

			repo.On("FindByHash", "testhash").Return(text, nil)
			repo.On("Store", text).Return(nil)
			centrifugo.On("Publish", mock.Anything, mock.Anything).Return(nil)
			dispatcher.On("Dispatch", mock.Anything).Return(nil)

			svc := NewStatisticsService(repo, dispatcher, centrifugo)
			err := svc.CalculateRank("testhash")

			require.NoError(t, err)
			require.InDelta(t, tt.expectedRank, text.GetRank(), 0.0001, "unexpected rank value")
		})
	}
}

func TestCalculateRank_Unicode(t *testing.T) {
	tests := []struct {
		name         string
		inputText    string
		expectedRank float64
	}{
		{
			name:         "Only emoji (😊😊😊)",
			inputText:    "😊😊😊",
			expectedRank: 0.0,
		},
		{
			name:         "Letters and emoji (A😊Б)",
			inputText:    "A😊Б",
			expectedRank: 2.0 / 3.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mockRepo)
			dispatcher := new(mockDispatcher)
			centrifugo := new(mockCentrifugo)

			text := &mockText{
				text: tt.inputText,
				hash: "testhash",
			}

			repo.On("FindByHash", "testhash").Return(text, nil)
			repo.On("Store", text).Return(nil)
			centrifugo.On("Publish", mock.Anything, mock.Anything).Return(nil)
			dispatcher.On("Dispatch", mock.Anything).Return(nil)

			svc := NewStatisticsService(repo, dispatcher, centrifugo)
			err := svc.CalculateRank("testhash")

			require.NoError(t, err)
			require.InDelta(t, tt.expectedRank, text.GetRank(), 0.0001, "unexpected rank value")
		})
	}
}

func TestCalculateRank_StoreError(t *testing.T) {
	repo := new(mockRepo)
	dispatcher := new(mockDispatcher)
	centrifugo := new(mockCentrifugo)

	text := &mockText{
		text: "Hello",
		hash: "abc123",
	}

	repo.On("FindByHash", "abc123").Return(text, nil)
	repo.On("Store", text).Return(errors.New("store error"))
	centrifugo.On("Publish", mock.Anything, mock.Anything).Return(nil)

	svc := NewStatisticsService(repo, dispatcher, centrifugo)
	err := svc.CalculateRank("abc123")

	require.Error(t, err)
	require.Contains(t, err.Error(), "store error")
}

func TestCalculateRank_CentrifugoError(t *testing.T) {
	repo := new(mockRepo)
	dispatcher := new(mockDispatcher)
	centrifugo := new(mockCentrifugo)

	text := &mockText{
		text: "Hi😊",
		hash: "xyz789",
	}

	repo.On("FindByHash", "xyz789").Return(text, nil)
	repo.On("Store", text).Return(nil)
	centrifugo.On("Publish", mock.Anything, mock.Anything).Return(errors.New("centrifugo failed"))

	svc := NewStatisticsService(repo, dispatcher, centrifugo)
	err := svc.CalculateRank("xyz789")

	require.Error(t, err)
	require.Contains(t, err.Error(), "centrifugo failed")
}

func TestCalculateRank_DispatchError(t *testing.T) {
	repo := new(mockRepo)
	dispatcher := new(mockDispatcher)
	centrifugo := new(mockCentrifugo)

	text := &mockText{
		text: "ABC",
		hash: "dispatch123",
	}

	repo.On("FindByHash", "dispatch123").Return(text, nil)
	repo.On("Store", text).Return(nil)
	centrifugo.On("Publish", mock.Anything, mock.Anything).Return(nil)
	dispatcher.On("Dispatch", mock.Anything).Return(errors.New("dispatch fail"))

	svc := NewStatisticsService(repo, dispatcher, centrifugo)
	err := svc.CalculateRank("dispatch123")

	require.Error(t, err)
	require.Contains(t, err.Error(), "dispatch fail")
}
