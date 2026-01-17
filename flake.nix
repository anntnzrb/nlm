{
  description = "nlm - NotebookLM CLI";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
    go-overlay.url = "github:purpleclay/go-overlay";
  };

  outputs = inputs@{ flake-parts, nixpkgs, go-overlay, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = nixpkgs.lib.systems.flakeExposed;

      perSystem =
        { config, system, ... }:
        let
          pkgs = import nixpkgs {
            inherit system;
            overlays = [ go-overlay.overlays.default ];
          };
          go = if pkgs ? go-bin then pkgs.go-bin.latest else pkgs.go;
        in
        {
          packages.nlm = pkgs.buildGoModule {
            pname = "nlm";
            version = "0-unstable-2026-01-17";

            src = ./.;
            vendorHash = "sha256-HGDejtwcHfOTUGwXjqCpwbe1tsOwULBAvLwm1VramRM=";
            subPackages = [ "cmd/nlm" ];
            inherit go;
          };

          packages.default = config.packages.nlm;

          apps.default = {
            type = "app";
            program = "${config.packages.nlm}/bin/nlm";
          };

          devShells.default = pkgs.mkShell {
            packages = [
              go
              pkgs.golangci-lint
              pkgs.gopls
              pkgs.gofumpt
              pkgs.gci
            ];
          };
        };
    };
}
