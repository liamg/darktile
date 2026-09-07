self:
{ config, lib, pkgs, ... }:

let
  cfg = config.programs.darktile;
in
{
  options.programs.darktile = {
    enable = lib.mkEnableOption "the darktile terminal emulator";

    package = lib.mkOption {
      type = lib.types.package;
      default = self.packages.${pkgs.stdenv.hostPlatform.system}.darktile;
      defaultText = lib.literalExpression "darktile.packages.\${system}.darktile";
      description = "The darktile package to install.";
    };
  };

  # darktile only ever reads config.yaml and theme.yaml from the user's config
  # directory, so there is nothing for a system-wide module to configure. Use
  # the home-manager module for settings and themes.
  config = lib.mkIf cfg.enable {
    environment.systemPackages = [ cfg.package ];
  };
}
