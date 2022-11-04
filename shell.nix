{ pkgs ? import <nixpkgs> { } }:
pkgs.mkShell {
  shellHook = ''
    export NO_COLOR=true
    export PS1="\u@\h:\w; "
  '';

  nativeBuildInputs = with pkgs.buildPackages; [
    go
    go-tools
  ];
}
