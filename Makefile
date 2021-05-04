root.key:
	openssl genrsa -out root.key 4096

root.crt: root.key
	openssl req -new -x509 -sha256 -key root.key -out root.crt \
		-subj "/C=PE/ST=Foo/L=Bar/O=Acme Inc. /OU=IT Department/CN=localhost"

pkg/credit/credit.pb.go: credit.proto
	mkdir -p pkg/credit
		 #/usr/local/bin/protoc ./credit.proto --plugin="/$$HOME/go/bin/protoc-gen-go" --go_out=plugins=grpc:.
	PATH=$$PATH:/usr/local/go/bin && \
		 protoc ./credit.proto --plugin="/$$HOME/go/bin/protoc-gen-go" --go_out=plugins=grpc:.

go.sum: go.mod
	export PATH=$$PATH:/usr/local/go/bin && \
		go get github.com/golang/protobuf/protoc-gen-go && \
		go mod download

.PHONY: server-reqs
server-reqs: go.sum root.key root.crt pkg/credit/credit.pb.go ## Run server.

.PHONY: server
server: server-reqs
	export PATH=$$PATH:/usr/local/go/bin && \
		go run ./server.go

.PHONY: broken-server
broken-server: server-reqs
	export PATH=$$PATH:/usr/local/go/bin && \
		go run ./server.go -with-h2=false

venv:
	python -m venv venv
	. venv/bin/activate && pip install -U setuptools grpcio-tools

.PHONY: client
client: root.crt venv
	. venv/bin/activate && GRPC_VERBOSITY=INFO GRPC_TRACE= python client.py

.PHONY: go-client
go-client: root.crt
	GODEBUG=x509ignoreCN=0 /tmp/grpcurl -proto ./credit.proto -cacert ./root.crt localhost:50051 credit.CreditService.Credit


.PHONY: deps
deps:
	echo "Install golang"
	curl -L https://golang.org/dl/go1.16.3.linux-amd64.tar.gz --output /tmp/go.tgz
	tar -C /usr/local -xzf /tmp/go.tgz
	
	echo "Installing protoc"
	apt update -q
	apt install -qy protobuf-compiler
	curl -L https://github.com/protocolbuffers/protobuf/releases/download/v3.15.8/protoc-3.15.8-linux-x86_64.zip --output /tmp/protoc.zip
	unzip /tmp/protoc.zip bin/protoc -d /usr/local/go
	
	echo "Installing grpcurl"
	curl -s -L https://github.com/fullstorydev/grpcurl/releases/download/v1.8.1/grpcurl_1.8.1_linux_x86_64.tar.gz --output - | tar xzf -C /tmp -

.PHONY: clean
clean:
	rm -f ./root.key ./root.crt
	rm -f go.sum
	rm -rf ./pkg/
	rm -rf ./venv
