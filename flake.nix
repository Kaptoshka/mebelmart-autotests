{
  description = "GUI Autotest Framework";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        browsers = (builtins.fromJSON
          (builtins.readFile "${pkgs.playwright-driver}/browsers.json")).browsers;

        chromium-rev = (builtins.head
          (builtins.filter (x: x.name == "chromium") browsers)).revision;

        firefox-rev = (builtins.head
          (builtins.filter (x: x.name == "firefox") browsers)).revision;
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go_1_26
            gopls
            delve

            golangci-lint
            statix
            deadnix
            nixpkgs-fmt

            playwright-driver.browsers
          ];

          shellHook = ''
            echo "Go GUI Autotest Framework"
            echo "  Go version: $(go version)"
            echo "  Browser: chromium-${chromium-rev} (Nix store)"
            echo "  Browser: firefox-${firefox-rev} (Nix store)"

            export PLAYWRIGHT_BROWSER_PATH=${pkgs.playwright-driver.browsers}
            export PLAYWRIGHT_SKIP_VALIDATE_HOST_REQUIREMENTS=true
            export PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1
            export PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH="${pkgs.playwright-driver.browsers}/chromium-${chromium-rev}/chrome-linux/chrome"
            export PLAYWRIGHT_FIREFOX_EXECUTABLE_PATH="${pkgs.playwright-driver.browsers}/firefox-${firefox-rev}/firefox/firefox"
          '';
        };
      }
    );
}
