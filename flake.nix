{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
  };
  outputs = inputs @ {flake-parts, ...}:
    flake-parts.lib.mkFlake {inherit inputs;} {
      systems = ["x86_64-linux" "aarch64-linux" "aarch64-darwin" "x86_64-darwin"];
      perSystem = {pkgs, ...}: {
        packages = rec {
          default = sbb-tui;
          sbb-tui = pkgs.buildGoModule {
            pname = "sbb-tui";
            version = "1.12.1";
            src = ./.;
            vendorHash = "sha256-K4DOu3rfSlKAa5JNKCzWWpnWZlXXxtN5Po7p1Spqe1w=";
          };
        };
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            gopls
            delve
          ];
        };
      };
    };
}
