{
  description = "A devShell example";

  inputs = {
    nixpkgs.url      = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url  = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        overlays = [  ];
        pkgs = import nixpkgs {
          inherit system overlays;
        };
      in
      with pkgs;
      {
        devShell = mkShell {
          env = {
            RESTIC_PASSWORD = "password";
            RESTIC_REPOSITORY = "./restic-repo";
          };
          buildInputs = [
            go
            restic
            hyperfine
            gopls
          ];
        };
      }
    );
}
