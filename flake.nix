{
  description = "E2E Autotest Framework";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        browsersJson = builtins.readFile "${pkgs.playwright-driver}/browsers.json";
        browsersData = builtins.fromJSON browsersJson;
        inherit (browsersData) browsers;

        chromium-rev = (builtins.head
          (builtins.filter (x: x.name == "chromium") browsers)).revision;

        firefox-rev = (builtins.head
          (builtins.filter (x: x.name == "firefox") browsers)).revision;

        webkit-rev = (builtins.head
          (builtins.filter (x: x.name == "webkit") browsers)).revision;
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            allure

            go_1_26
            gopls
            delve
            golangci-lint
            gotools

            statix
            deadnix
            nixpkgs-fmt

            yamllint
            yamlfmt

            nodejs
            playwright-driver.browsers
          ];

          shellHook = ''
            echo "Go E2E Autotest Framework"
            echo "  Go version: $(go version)"
            echo "  Node version: $(node --version)"
            echo "  Chromium rev: ${chromium-rev}"
            echo "  Firefox rev: ${firefox-rev}"
            echo "  WebKit rev: ${webkit-rev}"

            export ENV_FILE="$PWD/.env"

            export PLAYWRIGHT_BROWSER_PATH=${pkgs.playwright-driver.browsers}
            export PLAYWRIGHT_SKIP_VALIDATE_HOST_REQUIREMENTS=true
            export PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1

            export PLAYWRIGHT_NODEJS_PATH="${pkgs.nodejs}/bin/node"

            export PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH="${pkgs.playwright-driver.browsers}/chromium-${chromium-rev}/chrome-linux64/chrome"
            export PLAYWRIGHT_FIREFOX_EXECUTABLE_PATH="${pkgs.playwright-driver.browsers}/firefox-${firefox-rev}/firefox/firefox"
            export PLAYWRIGHT_WEBKIT_EXECUTABLE_PATH="${pkgs.playwright-driver.browsers}/webkit-${webkit-rev}/pw-run.sh"
          '';
        };
      }
    );
}
