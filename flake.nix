{
  description = "mcchunkie: a chat bot";

  inputs.nixpkgs.url = "nixpkgs/nixos-unstable";

  outputs =
    {
      self,
      nixpkgs,
    }:
    let
      supportedSystems = [
        "x86_64-linux"
        "x86_64-darwin"
        "aarch64-linux"
        "aarch64-darwin"
      ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      overlays.default = _: prev: {
        inherit (self.packages.${prev.stdenv.hostPlatform.system}) mcchunkie;
      };
      nixosModules.default = import ./module.nix;
      packages = forAllSystems (
        system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          mcchunkie = pkgs.buildGoModule rec {
            pname = "mcchunkie";
            version = "v1.1.4";
            src = ./.;

            vendorHash = "sha256-JIx1t61rVFVd7iXrmBJbc5XFzmHWLtX5Ni0jwmuiFTw=";

            # makes outbound http requests
            doCheck = false;

            ldflags = [ "-X suah.dev/mcchunkie/plugins.version=${version}" ];
          };
        }
      );

      devShells = forAllSystems (
        system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.mkShell {
            shellHook = ''
              PS1='\u@\h:\@; '
              nix run github:qbit/xin#flake-warn
              echo "Go `${pkgs.go}/bin/go version`"
            '';
            nativeBuildInputs = with pkgs; [
              git
              go
              gopls
              go-tools
              nilaway
            ];
          };
        }
      );
    };
}
