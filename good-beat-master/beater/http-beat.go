package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/freebirdrides/good-beat/config"
)

type HttpBeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
	events chan *common.MapStr
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &HttpBeat{
		done:   make(chan struct{}),
		config: config,
		events: make(chan *common.MapStr, config.MaxEvents),
	}
	return bt, nil
}

func (bt *HttpBeat) Run(b *beat.Beat) (err error) {
	logp.Info("good-beat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()
	ticker := time.NewTicker(bt.config.Period)
	counter := 1

	go func() {
		err = Run(TokenAuth{bt.config.AuthToken}, bt.events)
	}()

	for err == nil {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
			events := []common.MapStr{}
		loop:
			for {
				select {
				case evt := <-bt.events:
					events = append(events, *evt)
				default:
					if len(events) > 0 {
						bt.client.PublishEvents(events)
						logp.Info("Sent %d events", len(events))
					}
					break loop
				}
			}
		}
		counter++
	}

	return err
}

func (bt *HttpBeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
