build:
#	rm -f ./mutable_data/pathway_ledger.db
	mkdir -p ./mutable_data
	docker build -t jmixer .

sh:
	# run shell in detached mode
	docker container run \
		-v "${CURDIR}/mutable_data":/go/src/github.com/torsday/gemini_jobcoin_mixer/mutable_data \
		-it jmixer /bin/sh

