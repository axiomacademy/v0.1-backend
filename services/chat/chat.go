package chat

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/solderneer/axiom-backend/graph/model"
)

type MessageRange struct {
	To    string
	Start time.Time
	End   time.Time
}

type Chat struct {
	dbClient influxdb2.Client
	org      string
	bucket   string
	channels map[string]chan *model.Message
	mux      sync.Mutex
}

const defaultInfluxURL = "http://localhost:8086"
const defaultAuthToken = "user:pass"
const defaultOrg = "axiom"
const defaultBucket = "messages"

// Initialise a Chat struct with InfluxDB connection information taken from environment variables.
func InitChat() *Chat {
	influxURL := os.Getenv("INFLUX_URL")
	if influxURL == "" {
		influxURL = defaultInfluxURL
	}

	authToken := os.Getenv("INFLUX_AUTH_TOKEN")
	if authToken == "" {
		authToken = defaultAuthToken
	}

	org := os.Getenv("INFLUX_ORG")
	if org == "" {
		org = defaultOrg
	}

	bucket := os.Getenv("INFLUX_BUCKET")
	if bucket == "" {
		bucket = defaultBucket
	}

	return NewChat(influxURL, authToken, org, bucket)
}

// Initialise a Chat struct with InfluxDB connection information supplied as arguments.
func NewChat(influxURL string, authToken string, org string, bucket string) *Chat {
	c := &Chat{
		dbClient: influxdb2.NewClient(influxURL, authToken),
		org:      org,
		bucket:   bucket,
	}

	return c
}

func (c *Chat) Close() {
	c.dbClient.Close()
}

// Retrieve the messages to and from two certain users in a certain time range.
func (c *Chat) GetMessages(ctx context.Context, from string, r model.MessageRange) ([]*model.Message, error) {
	api := c.dbClient.QueryAPI(c.org)
	query := fmt.Sprintf(`from(bucket:"%s")|> range(start: %d, end: %d)|> filter(fn: (r) => r._measurement == "msg" and r.to == "%s" and r.from == "%s")`, c.bucket, r.Start.Unix(), r.End.Unix(), r.To, from)
	res, err := api.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	messages := make([]*model.Message, 0)
	for res.Next() {
		record := res.Record()
		vals := record.Values()
		m := &model.Message{
			To:        vals["to"].(string),
			From:      vals["form"].(string),
			Timestamp: record.Time(),
			Message:   vals["msg"].(string),
		}

		messages = append(messages, m)
	}

	return messages, nil
}

// Send a message.
func (c *Chat) SendMessage(ctx context.Context, sender string, message model.SendMessage) error {
	timestamp := time.Now()

	msg := influxdb2.NewPointWithMeasurement("msg").AddTag("to", message.To).AddTag("from", sender).AddField("msg", message.Message).SetTime(timestamp)

	api := c.dbClient.WriteAPIBlocking(c.org, c.bucket)
	err := api.WritePoint(ctx, msg)

	if err != nil {
		return err
	}

	c.mux.Lock()
	c.channels[message.To] <- &model.Message{
		To:        message.To,
		From:      sender,
		Timestamp: timestamp,
		Message:   message.Message,
	}
	c.mux.Unlock()

	return nil
}

// Returns a channel that messages to a certain user ID are sent to for the purposes of SSEs.
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
