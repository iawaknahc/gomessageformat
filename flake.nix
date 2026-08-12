{
  inputs.nixpkgs.url = "https://nixos.org/channels/nixos-26.05/nixexprs.tar.xz";
  inputs.devshell = {
    url = "github:numtide/devshell";
    inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs =
    { nixpkgs, devshell, ... }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      forEachSystem = nixpkgs.lib.genAttrs systems;
      pkgsFor =
        system:
        import nixpkgs {
          inherit system;
          overlays = [ devshell.overlays.default ];
        };
    in
    {
      devShells = forEachSystem (
        system:
        let
          pkgs = pkgsFor system;
        in
        {
          default = pkgs.devshell.mkShell {
            packages = [
              pkgs.icu78
              pkgs.icu78.dev
              pkgs.pkg-config
            ];
            env = [
              {
                name = "PKG_CONFIG_PATH";
                prefix = "${pkgs.icu78.dev}/lib/pkgconfig";
              }
            ];
          };
        }
      );
    };
}
