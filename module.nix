{
  lib,
  config,
  pkgs,
  ...
}:
let
  cfg = config.services.mcchunkie;
in
{
  options = with lib; {
    services.mcchunkie = {
      enable = lib.mkEnableOption "Enable mcchunkie";

      user = mkOption {
        type =
          with types;
          oneOf [
            str
            int
          ];
        default = "mcchunkie";
        description = ''
          The user the service will use.
        '';
      };

      group = mkOption {
        type =
          with types;
          oneOf [
            str
            int
          ];
        default = "mcchunkie";
        description = ''
          The group the service will use.
        '';
      };

      dataDir = mkOption {
        type = types.path;
        default = "/var/lib/mcchunkie";
        description = "Path mcchunkie will use to store data";
      };

      disabledPlugins = mkOption {
        type = with types; listOf str;
        default = [ ];
        description = "Plugins to disable.";
      };

      disabledChats = mkOption {
        type = with types; listOf str;
        default = [ ];
        description = "Chat services to disable.";
      };

      package = mkOption {
        type = types.package;
        default = pkgs.mcchunkie;
        defaultText = literalExpression "pkgs.mcchunkie";
        description = "The package to use for mcchunkie";
      };
    };
  };

  config = lib.mkIf (cfg.enable) {
    users.groups.${cfg.group} = { };
    users.users.${cfg.user} = {
      description = "mcchunkie service user";
      isSystemUser = true;
      home = "${cfg.dataDir}";
      createHome = true;
      group = "${cfg.group}";
    };

    systemd.services.mcchunkie = {
      enable = true;
      description = "mcchunkie server";
      after = [ "network-online.target" ];
      wants = [ "network-online.target" ];
      wantedBy = [ "multi-user.target" ];

      environment = {
        HOME = "${cfg.dataDir}";
      };

      serviceConfig = {
        User = cfg.user;
        Group = cfg.group;

        ExecStart =
          let
            inherit (builtins) concatStringsSep;
            mkMcOpt =
              flag: list:
              let
                realFlag = if list != [ ] then flag else "";
              in
              (concatStringsSep " " [
                realFlag
                (concatStringsSep "," list)
              ]);
            disabledPluginStr = if cfg.disabledPlugins != [ ] then (mkMcOpt "-dp" cfg.disabledPlugins) else "";
            disabledChatStr = if cfg.disabledChats != [ ] then (mkMcOpt "-dc" cfg.disabledChats) else "";
          in
          ''
            ${cfg.package}/bin/mcchunkie -db ${cfg.dataDir}/db ${disabledPluginStr} ${disabledChatStr}
          '';
      };
    };
  };
}
