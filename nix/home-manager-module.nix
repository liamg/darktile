self:
{ config, lib, pkgs, ... }:

let
  cfg = config.programs.darktile;
  yaml = pkgs.formats.yaml { };
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

    # Go's yaml decoder lower-cases struct field names, so the keys here are
    # all lower case with no separator: `brightblack`, not `brightBlack`.
    settings = lib.mkOption {
      type = yaml.type;
      default = { };
      example = lib.literalExpression ''
        {
          opacity = 0.9;
          font = {
            family = "JetBrains Mono";
            size = 14.0;
            ligatures = true;
          };
        }
      '';
      description = ''
        Written to {file}`$XDG_CONFIG_HOME/darktile/config.yaml`.
        Left alone when empty, so darktile falls back to its own defaults.
      '';
    };

    theme = lib.mkOption {
      type = yaml.type;
      default = { };
      example = lib.literalExpression ''
        {
          background = "#1d1f21";
          foreground = "#c5c8c6";
          brightblack = "#666666";
        }
      '';
      description = ''
        Written to {file}`$XDG_CONFIG_HOME/darktile/theme.yaml`.
        Left alone when empty, so darktile falls back to its built-in theme.
      '';
    };
  };

  config = lib.mkIf cfg.enable {
    home.packages = [ cfg.package ];

    xdg.configFile."darktile/config.yaml" = lib.mkIf (cfg.settings != { }) {
      source = yaml.generate "darktile-config.yaml" cfg.settings;
    };

    xdg.configFile."darktile/theme.yaml" = lib.mkIf (cfg.theme != { }) {
      source = yaml.generate "darktile-theme.yaml" cfg.theme;
    };
  };
}
