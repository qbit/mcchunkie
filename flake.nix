{
  description = "mcchunkie: a chat bot";

  inputs.nixpkgs.url = "nixpkgs/nixos-unstable";

  outputs =
    { self
    , nixpkgs
    ,
    }:
    let
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      overlay = _: prev: { inherit (self.packages.${prev.system}) mcchunkie; };
      nixosModule = import ./module.nix;
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          mcchunkie = pkgs.buildGoModule rec {
            pname = "mcchunkie";
            version = "v1.1.3";
            src = ./.;

            vendorHash = "sha256-J7ifgQzPwS6G4i+qYUpJmHQ9L3Ls3J6D9fvHD2J3a90=";

            # makes outbound http requests
            doCheck = false;

            ldflags = [ "-X suah.dev/mcchunkie/plugins.version=${version}" ];
          };
        });

      defaultPackage = forAllSystems (system: self.packages.${system}.mcchunkie);
      devShells = forAllSystems (system:
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
        });
    };
}
