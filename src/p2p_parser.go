package p2p_parser

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
)

func new_client(url string, trust_cert_file string, cert_file string, key_file string, allow_insecure_connection bool) pulsar.Client {
	var client pulsar.Client
	var err error
	var auth pulsar.Authentication

	if len(cert_file) > 0 || len(key_file) > 0 {
		auth = pulsar.NewAuthenticationTLS(cert_file, key_file)
	}

	client, err = pulsar.NewClient(pulsar.ClientOptions{
		URL:                        url,
		TLSAllowInsecureConnection: allow_insecure_connection,
		Authentication:             auth,
		TLSTrustCertsFilePath:      trust_cert_file,
		Logger:                     log.NewLoggerWithLogrus(logrus.New()),
	})

	if err != nil {
		logrus.Errorf("Failed connect to pulsar. Reason: %+v", err)
	}
	return client
}

func P2P_parser(opt Otp, parser_code Parser_code) {

	// Clients
	source_client := new_client(opt.sourcepulsar, opt.sourcetrustcerts, opt.sourcecertfile, opt.sourcekeyfile, opt.sourceallowinsecureconnection)
	dest_client := new_client(opt.destpulsar, opt.desttrustcerts, opt.destcertfile, opt.destkeyfile, opt.destallowinsecureconnection)

	defer source_client.Close()
	defer dest_client.Close()

	consume_chan := make(chan pulsar.ConsumerMessage, 2000)
	write_chan := make(chan *write_struct, 2000)
	ack_chan := make(chan pulsar.ConsumerMessage, 2000)

	consumer, err := source_client.Subscribe(pulsar.ConsumerOptions{
		Topic:                       opt.sourcetopic,
		SubscriptionName:            opt.sourcesubscription,
		Name:                        opt.sourcename,
		Type:                        pulsar.Shared,
		SubscriptionInitialPosition: pulsar.SubscriptionPositionLatest,
		MessageChannel:              consume_chan,
		ReceiverQueueSize:           2000,
	})
	if err != nil {
		logrus.Fatalln("Failed create consumer. Reason: ", err)
	}

	producer, err := dest_client.CreateProducer(pulsar.ProducerOptions{
		Topic:                   opt.desttopic,
		Name:                    opt.destname,
		BatchingMaxPublishDelay: time.Millisecond * time.Duration(opt.batchmaxpublishdelay),
		BatchingMaxMessages:     opt.batchmaxmessages,
		BatchingMaxSize:         opt.batchingmaxsize,
	})
	if err != nil {
		logrus.Fatalln("Failed create producer. Reason: ", err)
	}

	defer consumer.Close()
	defer producer.Close()

	for i := 0; i < int(opt.parserthreads); i++ {
		go parser(consume_chan, write_chan, ack_chan, parser_code)
	}

	go write_pulsar(write_chan, producer, ack_chan)

	acks(consumer, ack_chan)
}
