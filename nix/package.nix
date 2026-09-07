{ lib
, stdenv
, buildGoModule
, pkg-config
, makeWrapper
, makeDesktopItem
, copyDesktopItems
, libGL
, libx11
, libxcursor
, libxi
, libxinerama
, libxrandr
, libxxf86vm
, src
, version
}:

let
  # On darwin ebiten and glfw go through Cocoa, which the stdenv already
  # provides, so none of the X11 stack applies and there is nothing to add in
  # its place.
  isLinux = stdenv.hostPlatform.isLinux;

  graphicsLibs = lib.optionals isLinux [
    libGL
    libx11
    libxcursor
    libxi
    libxinerama
    libxrandr
    libxxf86vm
  ];
in
buildGoModule {
  pname = "darktile";
  inherit src version;

  # The dependencies are committed under vendor/.
  vendorHash = null;

  subPackages = [ "cmd/darktile" ];

  nativeBuildInputs = [ pkg-config ]
    ++ lib.optionals isLinux [ makeWrapper copyDesktopItems ];
  buildInputs = graphicsLibs;

  ldflags = [
    "-s"
    "-w"
    "-X github.com/liamg/darktile/internal/app/darktile/version.Version=${version}"
  ];

  # ebiten resolves libGL with dlopen rather than linking it, so nothing here
  # ends up in the binary's DT_NEEDED and the loader has no way to find it.
  postInstall = lib.optionalString isLinux ''
    wrapProgram $out/bin/darktile \
      --prefix LD_LIBRARY_PATH : ${lib.makeLibraryPath graphicsLibs}
  '';

  desktopItems = lib.optionals isLinux [
    (makeDesktopItem {
      name = "darktile";
      exec = "darktile";
      icon = "utilities-terminal";
      desktopName = "Darktile";
      comment = "GPU rendered terminal emulator";
      categories = [ "System" "TerminalEmulator" ];
      terminal = false;
    })
  ];

  meta = {
    description = "GPU rendered terminal emulator";
    homepage = "https://github.com/liamg/darktile";
    license = lib.licenses.mit;
    mainProgram = "darktile";
    # Darwin builds only work when run on a Mac. nixpkgs cannot cross-compile
    # to darwin from anywhere else: cctools, the linker, refuses to evaluate
    # off a darwin host, so there is no toolchain to cross with.
    platforms = lib.platforms.linux ++ lib.platforms.darwin;
  };
}
