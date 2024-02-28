

run_example:
	go build -C example -o ../p2p_parser
	./p2p_parser

clean:
	rm p2p_parser