build:
	mkdir -p bin
	CGO_ENABLED=0 go build -trimpath -o bin/ddns-go .
clean:
	rm -r bin/
