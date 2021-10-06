package main

import (
	"time"

	"github.com/google/uuid"
)

type gsMsg struct {
	Id        uuid.UUID `json:"id"`
	Msg       string    `json:"msg"`
	Timestamp string    `json:"timestamp"`
	Type      string    `json:"type"` // MSG or CLOSE
	Topic     string    `json:"topic"`
}

func (gsMsg) Create(msg string, topic string) gsMsg {
	id := uuid.New()
	timestamp := time.Now().UTC().Format(time.RFC3339)
	return gsMsg{Id: id, Timestamp: timestamp, Msg: msg, Type: "MSG", Topic: topic}
}

func (gsMsg) CreateCloseMsg() gsMsg {
	id := uuid.New()
	timestamp := time.Now().UTC().Format(time.RFC3339)
	return gsMsg{Id: id, Timestamp: timestamp, Msg: "", Type: "CLOSE"}
}

func (gsMsg gsMsg) GetMsg() string {
	return gsMsg.Msg
}

func (gsMsg gsMsg) Stringify() string {
	return gsMsg.Timestamp + " " + gsMsg.Id.String() + " " + gsMsg.Topic + " " + gsMsg.Msg
}
