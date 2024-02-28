package main

import (
	p2p_parser "example.com/p2p_parser/src"
	"github.com/sirupsen/logrus"
)

type does_nothing struct{}

func (dn *does_nothing) Parse(msg []byte) [][]byte {
	same := make([][]byte, 0)
	same = append(same, msg)
	return same
}

func logging(level string) {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
	})
	l, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Errorf("Failed parse log level. Reason: %+v", err)
	} else {
		logrus.SetLevel(l)
	}
}

func main() {

	opt := p2p_parser.From_args()

	logging("debug")

	parser := &does_nothing{}

	p2p_parser.P2P_parser(opt, parser)
}
