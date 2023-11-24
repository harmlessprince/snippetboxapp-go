.PHONY: serve

test:
	#go run cmd/web/!(*_test).go -addr=":9999"
	#go run cmd/web/* -addr=":9999" >>./logs/web.info.log 2>> ./logs/web.error.log
	go test -v ./cmd/web
serve:
	#cd cmd/web
	#go run $(find cmd/web -name '*.go' ! -name '*_test.go')
	go run cmd/web/^*_test.go


help:
	go run cmd/web/* -help

createSnippet:
	 curl -iL -X POST http://localhost:9999/create/snippet
getSnippet:
	 curl -iL -X GET http://localhost:9999/snippet?id=$(id)
fetchSnippets:
	 curl -iL -X GET http://localhost:9999