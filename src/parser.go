package p2p_parser

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/apache/pulsar-client-go/pulsar"
)

type Parser_code interface {
	Parse([]byte) [][]byte
}

func parser(consume_chan <-chan pulsar.ConsumerMessage, write_chan chan<- *write_struct, ack_chan chan<- pulsar.ConsumerMessage, parser Parser_code) {
	var n_read float64 = 0

	last_instant := time.Now()
	last_publish_time := time.Unix(0, 0)
	tick := time.NewTicker(time.Minute)
	defer tick.Stop()

	for {
		select {
		case msg := <-consume_chan:

			n_read++
			last_publish_time = msg.PublishTime()

			opsProcessed.Inc()

			parsed_vec := parser.Parse(msg.Payload())

			if len(parsed_vec) == 0 {
				ack_chan <- msg
			} else {
				write_chan <- &write_struct{
					msg:         msg,
					payload_vec: parsed_vec,
				}
			}

		case <-tick.C:
			since := time.Since(last_instant)
			last_instant = time.Now()
			logrus.Infof("Read rate: %.3f msg/s; (last pulsar time %v)", n_read/float64(since/time.Second), last_publish_time)
			n_read = 0
		}
	}
}

// func write_callback(agg []pulsar.ConsumerMessage, ack_chan chan<- pulsar.ConsumerMessage) func(msgID pulsar.MessageID, pm *pulsar.ProducerMessage, err error) {
// 	return func(msgID pulsar.MessageID, pm *pulsar.ProducerMessage, err error) {
// 		for i := 0; i < len(agg); i++ {
// 			ack_chan <- (agg)[i]
// 		}
// 	}
// }

type write_struct struct {
	msg         pulsar.ConsumerMessage
	payload_vec [][]byte
}

func write_pulsar(write_chan <-chan *write_struct, producer pulsar.Producer, ack_chan chan<- pulsar.ConsumerMessage) {
	var n_writen float64 = 0

	for write_info := range write_chan {

		for _, parsed := range write_info.payload_vec {
			n_writen++
			producer.SendAsync(
				context.Background(),
				&pulsar.ProducerMessage{
					Payload: parsed,
					Key:     write_info.msg.Key(),
				},
				func(msgID pulsar.MessageID, pm *pulsar.ProducerMessage, err error) {
					// TODO use atomic counters
					// if err == nil {
					// 	ack_chan <- msg
					// } / dont ack messages twice if more than 1 fit
				},
			)
		}
		ack_chan <- write_info.msg
	}
}

func acks(consumer pulsar.Consumer, ack_chan <-chan pulsar.ConsumerMessage) {
	last_instant := time.Now()
	tick := time.NewTicker(time.Minute)
	defer tick.Stop()

	var ack float64 = 0
	for {
		select {
		case msg := <-ack_chan:

			if err := consumer.Ack(msg); err != nil {
				logrus.Warnf("consumer.Acks err: %+v", err)
			}
			ack++

			//prom_metrics.Prom_metric.Inc_number_processed_msg()

		case <-tick.C:
			since := time.Since(last_instant)
			last_instant = time.Now()
			logrus.Infof("Ack rate: %.3f msg/s", ack/float64(since/time.Second))
			ack = 0
		}
	}
}
