{
  description = "terraform-provider-jellyfin dev flake.";

  inputs = {
    nixpkgs = {
      url = "github:nixos/nixpkgs/nixos-unstable";
    };
  };

  outputs = {
    nixpkgs,
    ...
  }:
  let
    system = "x86_64-linux";
    pkgs = import nixpkgs { inherit system; };
  in
  {
    devShells.${pkgs.stdenv.hostPlatform.system}.default = pkgs.mkShell {
      packages = with pkgs; [
        go
        gopls
        gnumake
        opentofu
      ];
    };
  };
}
