package chat

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/solderneer/axiom-backend/graph/model"
	"github.com/influxdata/influxdb-client-go/v2"
)

type MessageRange struct {
	To string
	Start time.Time
	End time.Time
}

func FromModelMessageRange(m model.MessageRange) (*MessageRange, error) {
	start, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", m.Start)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", m.End)
	if err != nil {
		return nil, err
	}

	mr := &MessageRange {
		To: m.To,
		Start: start,
		End: end,
	}

	return mr, nil
}

type Chat struct {
	dbClient influxdb2.Client
	org string
	bucket string
	channels map[string] chan *model.Message
	mux sync.Mutex
}

func NewChat(influxURL string, authToken string, org string, bucket string) *Chat {
	c := &Chat {
		dbClient: influxdb2.NewClient(influxURL, authToken),
		org: org,
		bucket: bucket,
	}

	return c
}

func (c *Chat) Close() {
	c.dbClient.Close()
}

func (c *Chat) GetMessages(ctx context.Context, from string, r model.MessageRange) ([]*model.Message, error) {
	mr, err := FromModelMessageRange(r)
	if err != nil {
		return nil, err
	}

	api := c.dbClient.QueryAPI(c.org)
	query := fmt.Sprintf(`from(bucket:"%s")|> range(start: %d, end: %d)|> filter(fn: (r) => r._measurement == "msg" and r.to == "%s" and r.from == "%s")`, c.bucket, mr.Start.Unix(), mr.End.Unix(), mr.To, from)
	res, err := api.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	messages := make([]*model.Message, 0)
	for res.Next() {
		record := res.Record()
		vals := record.Values()
		m := &model.Message {
			To: vals["to"].(string),
			From: vals["form"].(string),
			Timestamp: record.Time().String(),
			Message: vals["msg"].(string),
		}

		messages = append(messages, m)
	}

	return messages, nil
}

func (c *Chat) SendMessage(ctx context.Context, sender string, message model.SendMessage) error {
	timestamp := time.Now()

	msg := influxdb2.NewPointWithMeasurement("msg").AddTag("to", message.To).AddTag("from", sender).AddField("msg", message.Message).SetTime(timestamp)

	api := c.dbClient.WriteAPIBlocking(c.org, c.bucket)
	err := api.WritePoint(ctx, msg)

	if err != nil {
		return err
	}

	c.mux.Lock()
	c.channels[message.To] <- &model.Message {
		To: message.To,
		From: sender,
		Timestamp: timestamp.String(),
		Message: message.Message,
	}
	c.mux.Unlock()

	return nil
}

func (c *Chat) SubscribeMessages(uid string, done <-chan struct{}) <-chan *model.Message {
	channel := make(chan *model.Message)
	c.mux.Lock()
	c.channels[uid] = channel
	c.mux.Unlock()

	go func() {
		<-done
		c.mux.Lock()
		delete(c.channels, uid)
		c.mux.Unlock()
	}()

	return channel
}
