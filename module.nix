{ lib, config, pkgs, ... }:
let cfg = config.services.mcchunkie;
in {
  options = with lib; {
    services.mcchunkie = {
      enable = lib.mkEnableOption "Enable mcchunkie";

      user = mkOption {
        type = with types; oneOf [ str int ];
        default = "mcchunkie";
        description = ''
          The user the service will use.
        '';
      };

      group = mkOption {
        type = with types; oneOf [ str int ];
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
      wantedBy = [ "network-online.target" ];
      after = [ "network-online.target" ];

      environment = { HOME = "${cfg.dataDir}"; };

      serviceConfig = {
        User = cfg.user;
        Group = cfg.group;

        ExecStart = ''
          ${cfg.package}/bin/mcchunkie -db ${cfg.dataDir}/db
        '';
      };
    };
  };
}
