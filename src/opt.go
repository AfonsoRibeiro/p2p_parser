package p2p_parser

import (
	"github.com/jnovack/flag"
)

type Otp struct {
	sourcepulsar                  string
	sourcetopic                   string
	sourcesubscription            string
	sourcename                    string
	sourcetrustcerts              string
	sourcecertfile                string
	sourcekeyfile                 string
	sourceallowinsecureconnection bool

	destpulsar                  string
	desttopic                   string
	destsubscription            string
	destname                    string
	desttrustcerts              string
	destcertfile                string
	destkeyfile                 string
	destallowinsecureconnection bool

	batchmaxpublishdelay uint
	batchmaxmessages     uint
	batchingmaxsize      uint

	parserthreads uint
}

func From_args() Otp {

	var opt Otp

	flag.StringVar(&opt.sourcepulsar, "source_pulsar", "pulsar://localhost:6650", "Source pulsar address")
	flag.StringVar(&opt.sourcetopic, "source_topic", "persistent://public/default/in", "Source topic name")
	flag.StringVar(&opt.sourcesubscription, "source_subscription", "jq_parser", "Source subscription name")
	flag.StringVar(&opt.sourcename, "source_name", "jq_parser_consumer", "Source consumer name")
	flag.StringVar(&opt.sourcetrustcerts, "source_trust_certs", "", "Path for source pem file, for ca.cert")
	flag.StringVar(&opt.sourcecertfile, "source_cert_file", "", "Path for source cert.pem file")
	flag.StringVar(&opt.sourcekeyfile, "source_key_file", "", "Path for source key-pk8.pem file")
	flag.BoolVar(&opt.sourceallowinsecureconnection, "source_allow_insecure_connection", false, "Source allow insecure connection")

	flag.StringVar(&opt.destpulsar, "dest_pulsar", "pulsar://localhost:6650", "Destination pulsar address")
	flag.StringVar(&opt.desttopic, "dest_topic", "persistent://public/default/out", "Destination topic name")
	flag.StringVar(&opt.destsubscription, "dest_subscription", "jq_parser", "Destination subscription name")
	flag.StringVar(&opt.destname, "dest_name", "jq_parser_consumer", "Destination producer name")
	flag.StringVar(&opt.desttrustcerts, "dest_trust_certs", "", "Path for destination pem file, for ca.cert")
	flag.StringVar(&opt.destcertfile, "dest_cert_file", "", "Path for destination cert.pem file")
	flag.StringVar(&opt.destkeyfile, "dest_key_file", "", "Path for destination key-pk8.pem file")
	flag.BoolVar(&opt.destallowinsecureconnection, "dest_allow_insecure_connection", false, "Dest allow insecure connection")

	flag.UintVar(&opt.batchmaxpublishdelay, "batch_max_publish_delay", 10, "How long to wait for batching in milliseconds")
	flag.UintVar(&opt.batchmaxmessages, "batch_max_messages", 300, "Max batch messages")
	flag.UintVar(&opt.batchingmaxsize, "batch_max_size", 131072, "Max batch size in bytes")

	flag.UintVar(&opt.parserthreads, "parser_threads", 4, "Number of parsing threads")

	flag.Parse()

	return opt

}
