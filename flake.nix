{
  description = "darktile - GPU rendered terminal emulator";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let
      # aarch64-darwin only builds on a Mac, not from here: nixpkgs has no
      # darwin cross toolchain, since cctools refuses to evaluate off a darwin
      # host. It is listed so `nix build` works when someone runs it there.
      # x86_64-darwin is absent because nixpkgs 26.11 dropped it.
      systems = [ "x86_64-linux" "aarch64-linux" "aarch64-darwin" ];
      forAllSystems = f: nixpkgs.lib.genAttrs systems (system: f nixpkgs.legacyPackages.${system});

      graphicsLibsFor = pkgs: with pkgs; [
        libGL
        libx11
        libxcursor
        libxi
        libxinerama
        libxrandr
        libxxf86vm
      ];

      # A flake carries no tag information: `self` exposes rev, shortRev and
      # revCount, and nothing else. So the released version is kept in a file
      # and the release workflow refuses to publish a tag that disagrees with
      # it, which is what stops this from drifting away from GoReleaser's
      # `{{.Version}}`. The suffix mirrors `git describe`, since a Nix build is
      # almost never exactly the tagged commit.
      releaseVersion = nixpkgs.lib.trim (builtins.readFile ./VERSION);
      version = releaseVersion + (if self ? rev then "-g${self.shortRev}" else "-dirty");
    in
    {
      packages = forAllSystems (pkgs: rec {
        darktile = pkgs.callPackage ./nix/package.nix {
          src = self;
          inherit version;
        };
        default = darktile;
      });

      overlays.default = final: prev: {
        darktile = final.callPackage ./nix/package.nix {
          src = self;
          inherit version;
        };
      };

      nixosModules.default = import ./nix/nixos-module.nix self;
      homeManagerModules.default = import ./nix/home-manager-module.nix self;

      devShells = forAllSystems (pkgs:
        let
          graphicsLibs = graphicsLibsFor pkgs;
        in
        {
          default = pkgs.mkShell {
            nativeBuildInputs = with pkgs; [
              go
              pkg-config
              git
              gopls
              delve
              goreleaser
              govulncheck

              # conformance and performance test suites
              python3
              goperf # provides benchstat
              graphviz # pprof renders call graphs through dot
              xvfb-run # render goldens need a display
              imagemagick # and a way to diff the screenshots
              util-linux # flock, so two render runs cannot stomp each other
            ];

            buildInputs = graphicsLibs;

            # ebiten from v2.9 on resolves libGL with dlopen instead of linking
            # it, so having these as buildInputs is not enough to run anything.
            #
            # MESA_PREFIX is for the render goldens: on NixOS the drivers behind
            # that libGL come from /run/opengl-driver, which exists nowhere else,
            # so on a CI runner there is nothing for libglvnd to load. The script
            # points at this mesa when it launches darktile. It cannot be
            # exported here as LIBGL_DRIVERS_PATH, because Xvfb would pick it up
            # too and crash mixing it with the system mesa.
            shellHook = ''
              export LD_LIBRARY_PATH="${pkgs.lib.makeLibraryPath graphicsLibs}''${LD_LIBRARY_PATH:+:$LD_LIBRARY_PATH}"
              export MESA_PREFIX="${pkgs.mesa}"
            '';
          };
        });
    };
}
