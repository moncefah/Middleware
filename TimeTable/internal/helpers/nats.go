package helpers

import "github.com/nats-io/nats.go"

var NatsConn *nats.Conn

func InitNats(url string) error {
	nc, err := nats.Connect(url)
	if err != nil {
		return err
	}
	NatsConn = nc
	return nil
}

func CloseNats() {
	if NatsConn != nil {
		NatsConn.Close()
	}
}
