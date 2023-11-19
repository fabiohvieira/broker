package sqs

import (
	"context"
	"github.com/fabiohvieira/broker/broker"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SQSTestSuite struct {
	suite.Suite
}

func TestSQSSuite(t *testing.T) {
	suite.Run(t, new(SQSTestSuite))
}

func (s *SQSTestSuite) SetupTest() {
}

func (s *SQSTestSuite) TestNewSQSClient_WhenAllOk_ThenReturnClient() {
	// act
	cli, err := NewSQSClient(context.Background())

	// assert
	s.NoError(err)
	s.NotNil(cli)
}

func (s *SQSTestSuite) TestSendMessage_WhenAllOk_ThenReturnNil() {
	// arrange
	cli, _ := NewSQSClient(context.Background())
	b := New(cli)
	message := broker.Message{
		Body: []byte("test"),
	}

	// act
	err := b.SendMessage(context.Background(), message, "topic")

	// assert
	s.NoError(err)
}
