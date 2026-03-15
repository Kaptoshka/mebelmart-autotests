lint: lint-go lint-yaml lint-nix

lint-go:
	golangci-lint run --config .golangci.yml

lint-yaml:
	yamllint -c .yamllint.yml .

lint-nix:
	nixpkgs-fmt --check .
	statix check .
	deadnix --fail .
