.PHONY: test diagrams audit-demo lint
test:
	cd src && go test -cover ./...
lint:
	cd src && go vet ./...
diagrams:
	for f in diagrams/src/*.dot; do b=$$(basename $$f .dot); \
	  dot -Tsvg $$f -o diagrams/export/$$b.svg; dot -Tpng -Gdpi=130 $$f -o diagrams/export/$$b.png; done
audit-demo:
	cd src && go run ./cmd/pam-demo-audit > ../tests/data/sample-audit.jsonl && \
	go run ./cmd/pam-audit verify -f ../tests/data/sample-audit.jsonl
