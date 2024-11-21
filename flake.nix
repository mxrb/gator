{
  description = "A simple Go package";

  inputs.nixpkgs.url = "nixpkgs/nixos-24.11";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages = rec {
          gator = pkgs.buildGoModule {
            pname = "gator";
            version = self.lastModifiedDate;
            src = ./.;
            vendorHash = "sha256-imzx9m29sUOSvXrBWSAH3oiRCkcfpsQJgH7tDKaH7VU=";
          };
          default = gator;
        };
        apps = rec {
          gator = flake-utils.lib.mkApp { drv = self.packages.${system}.gator; };
          default = gator;
        };
        devShells.default =
          let
            pgUser = "gator";
            pgPassword = "supersecretpassword";
            pgPort = "5432";
            connectionString = "postgres://${pgUser}:${pgPassword}@localhost:${pgPort}/${pgUser}?sslmode=disable";
          in
          pkgs.mkShell {
            buildInputs = with pkgs; [
              go
              gopls
              gotools
              go-tools
              goose
              sqlc
              (writeShellScriptBin "run-db" ''
                ${pkgs.podman}/bin/podman run --rm \
                --name gator-db \
                -e POSTGRES_USER=${pgUser} \
                -e POSTGRES_PASSWORD=${pgPassword} \
                -p 127.0.0.1:${pgPort}:5432 \
                docker.io/library/postgres:latest
              '')
            ];
            GOOSE_DRIVER = "postgres";
            GOOSE_DBSTRING = connectionString;
            shellHook = ''
              export GOOSE_MIGRATION_DIR=$(realpath ./sql/schema)
            '';
          };
      }
    );
}
