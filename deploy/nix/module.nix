# NixOS module for Hyperuplink.
#
# Runs the binary from package.nix as a hardened system service. The layout
# mirrors the RPM and the Gentoo ebuild (dedicated `hyperuplink` user, state
# under /var/lib/hyperuplink, ImageMagick on PATH), except that the config is
# generated from `services.hyperuplink.settings` rather than shipped to /etc,
# and PostgreSQL and Redis can be provisioned along with it.
# See deploy/nix/README.md.
{
  config,
  lib,
  pkgs,
  ...
}:

let
  cfg = config.services.hyperuplink;

  format = pkgs.formats.toml { };
  configFile = format.generate "hyperuplink.toml" cfg.settings;

  stateDir = "/var/lib/hyperuplink";
  runtimeDir = "/run/hyperuplink";

  # Unix socket plus peer authentication, so the local case needs no password.
  localDatabaseConnection = "postgres:///${cfg.user}?host=/run/postgresql";

  # Everything in `settings` lands in the world-readable store, so with secrets
  # the config is rendered a second time at start-up, into the private runtime
  # directory, with the placeholders filled in from environmentFile.
  usesSecrets = cfg.environmentFile != null;
  activeConfig = if usesSecrets then "${runtimeDir}/hyperuplink.toml" else configFile;
in
{
  options.services.hyperuplink = {
    enable = lib.mkEnableOption "Hyperuplink, a super high speed internet bulletin board";

    package = lib.mkPackageOption pkgs "hyperuplink" { };

    user = lib.mkOption {
      type = lib.types.str;
      default = "hyperuplink";
      description = ''
        User the service runs as. It is created automatically when left at the
        default, and it doubles as the PostgreSQL role name when
        {option}`services.hyperuplink.database.createLocally` is set.
      '';
    };

    group = lib.mkOption {
      type = lib.types.str;
      default = "hyperuplink";
      description = "Group the service runs as. Created automatically when left at the default.";
    };

    environmentFile = lib.mkOption {
      type = lib.types.nullOr lib.types.path;
      default = null;
      example = "/var/lib/secrets/hyperuplink.env";
      description = ''
        Environment file (as read by {manpage}`systemd.exec(5)`'s
        `EnvironmentFile=`) holding secrets that must not be written to the Nix
        store.

        When this is set, `$VARIABLE` and `''${VARIABLE}` placeholders anywhere
        in {option}`services.hyperuplink.settings` are substituted at start-up,
        so a password can be kept out of the store:

        ```nix
        services.hyperuplink.settings.Database.Connection =
          "postgres://hyperuplink:$DB_PASSWORD@db.example.com:5432/hyperuplink";
        ```

        Note that this makes every literal `$` in `settings` significant; write
        `$$` for a dollar sign that should survive.
      '';
    };

    database.createLocally = lib.mkOption {
      type = lib.types.bool;
      default = true;
      description = ''
        Whether to set up a local PostgreSQL server with a database and role
        named after {option}`services.hyperuplink.user`, and point the service
        at it over a Unix socket with peer authentication, which needs no
        password and so keeps secrets out of the store entirely.

        Turn this off to use an external database and set
        {option}`services.hyperuplink.settings.Database.Connection` yourself.
      '';
    };

    redis.createLocally = lib.mkOption {
      type = lib.types.bool;
      default = true;
      description = ''
        Whether to set up a local Redis server for the background job queue and
        point the service at it. Turn this off to use an external Redis/Valkey
        and set {option}`services.hyperuplink.settings.Redis` yourself.
      '';
    };

    redis.port = lib.mkOption {
      type = lib.types.port;
      default = 6379;
      description = ''
        TCP port of the local Redis server. It has to be a TCP port, because the
        job queue (asynq) cannot speak to a Unix socket. Change it when the
        default collides with another Redis instance on this host.
      '';
    };

    openFirewall = lib.mkOption {
      type = lib.types.bool;
      default = false;
      description = ''
        Whether to open {option}`services.hyperuplink.settings.Server.Port` in
        the firewall. Leave this off when a reverse proxy on the same host
        terminates TLS, which is the recommended setup.
      '';
    };

    settings = lib.mkOption {
      type = lib.types.submodule { freeformType = format.type; };
      default = { };
      example = lib.literalExpression ''
        {
          Users.PromoteAdmin = [ "me@example.com" ];
          Server.ProxyHeader = "X-Forwarded-For";
          Logging.Level = "debug";
        }
      '';
      description = ''
        Configuration written to a TOML file and passed to the binary with `-c`.
        See deploy/hyperuplink.toml for the full set of keys.

        This ends up in the world-readable Nix store, so keep secrets out of it
        and use {option}`services.hyperuplink.environmentFile` instead.
      '';
    };
  };

  config = lib.mkIf cfg.enable {
    assertions = [
      {
        assertion =
          cfg.database.createLocally -> cfg.settings.Database.Connection == localDatabaseConnection;
        message = ''
          services.hyperuplink.database.createLocally provisions a database owned
          by the "${cfg.user}" role and connects to it over a Unix socket, so it
          cannot be combined with a custom
          services.hyperuplink.settings.Database.Connection. Set
          database.createLocally = false to point at your own database.
        '';
      }
      {
        assertion = !cfg.database.createLocally -> (cfg.settings.Database.Connection or "") != "";
        message = ''
          Hyperuplink needs a database. Either set
          services.hyperuplink.database.createLocally = true, or point
          services.hyperuplink.settings.Database.Connection at an existing
          PostgreSQL server.
        '';
      }
      {
        assertion = !cfg.redis.createLocally -> (cfg.settings.Redis.Addresses or [ ]) != [ ];
        message = ''
          Hyperuplink needs a Redis/Valkey server for its job queue. Either set
          services.hyperuplink.redis.createLocally = true, or list an existing
          server in services.hyperuplink.settings.Redis.Addresses.
        '';
      }
    ];

    warnings = lib.optional ((cfg.settings.Users.PromoteAdmin or [ ]) == [ ]) ''
      services.hyperuplink.settings.Users.PromoteAdmin is empty, so no account
      will be promoted to admin when it signs up, and the board will have nobody
      to administer it.
    '';

    services.hyperuplink.settings = {
      General.Mode = lib.mkDefault "production";
      Logging.Level = lib.mkDefault "info";

      Server = {
        BindIP = lib.mkDefault "0.0.0.0";
        Port = lib.mkDefault 3000;
        BodyLimit = lib.mkDefault 67108864;
        Concurrency = lib.mkDefault 262144;
        TrustLoopback = lib.mkDefault true;
      };

      # A local Postgres is reached over its socket as the service user, which
      # peer authentication accepts without a password.
      Database.Connection = lib.mkIf cfg.database.createLocally (lib.mkDefault localDatabaseConnection);

      Redis.Addresses = lib.mkIf cfg.redis.createLocally (
        lib.mkDefault [ "127.0.0.1:${toString cfg.redis.port}" ]
      );

      # Matches the media directory created by StateDirectory= below. Nothing
      # creates it at runtime, so the two have to agree.
      Storage = lib.mkDefault [
        {
          ID = "local-storage";
          Type = "Local";
          Local = {
            Path = "${stateDir}/media";
            PublicURI = "/media";
          };
        }
      ];
    };

    users.users = lib.mkIf (cfg.user == "hyperuplink") {
      hyperuplink = {
        description = "Hyperuplink service";
        group = cfg.group;
        home = stateDir;
        isSystemUser = true;
      };
    };

    users.groups = lib.mkIf (cfg.group == "hyperuplink") {
      hyperuplink = { };
    };

    services.postgresql = lib.mkIf cfg.database.createLocally {
      enable = true;
      ensureDatabases = [ cfg.user ];
      ensureUsers = [
        {
          name = cfg.user;
          # Needs the database to be named after the role, which is why the
          # database is not separately configurable.
          ensureDBOwnership = true;
        }
      ];
    };

    services.redis.servers.hyperuplink = lib.mkIf cfg.redis.createLocally {
      enable = true;
      bind = "127.0.0.1";
      port = cfg.redis.port;
    };

    networking.firewall.allowedTCPPorts = lib.mkIf cfg.openFirewall [ cfg.settings.Server.Port ];

    systemd.services.hyperuplink = {
      description = "Hyperuplink";
      wantedBy = [ "multi-user.target" ];

      wants = [ "network-online.target" ];
      # postgresql.target, not postgresql.service, because it also pulls in
      # postgresql-setup.service, which is what creates the database and the
      # role. Ordering after the bare server would race with it.
      after = [
        "network-online.target"
      ]
      ++ lib.optional cfg.database.createLocally "postgresql.target"
      ++ lib.optional cfg.redis.createLocally "redis-hyperuplink.service";
      requires =
        lib.optional cfg.database.createLocally "postgresql.target"
        ++ lib.optional cfg.redis.createLocally "redis-hyperuplink.service";

      preStart = lib.mkIf usesSecrets ''
        umask 077
        ${lib.getExe pkgs.envsubst} -i ${configFile} -o ${activeConfig}
      '';

      serviceConfig = {
        Type = "simple";
        User = cfg.user;
        Group = cfg.group;
        ExecStart = "${lib.getExe cfg.package} -c file://${activeConfig}";
        Restart = "on-failure";
        RestartSec = "5s";

        EnvironmentFile = lib.mkIf usesSecrets [ cfg.environmentFile ];
        RuntimeDirectory = lib.mkIf usesSecrets "hyperuplink";
        RuntimeDirectoryMode = lib.mkIf usesSecrets "0700";

        # Local media storage lives here, and nothing creates the media
        # subdirectory at runtime, so systemd has to.
        StateDirectory = [
          "hyperuplink"
          "hyperuplink/media"
        ];
        StateDirectoryMode = "0750";
        WorkingDirectory = stateDir;

        # Mirrors deploy/init/hyperuplink.service.
        NoNewPrivileges = true;
        PrivateTmp = true;
        PrivateDevices = true;
        ProtectSystem = "strict";
        ProtectHome = true;
        ProtectControlGroups = true;
        ProtectKernelModules = true;
        ProtectKernelTunables = true;
        ProtectKernelLogs = true;
        ProtectClock = true;
        ProtectHostname = true;
        ProtectProc = "invisible";
        RestrictNamespaces = true;
        RestrictRealtime = true;
        RestrictSUIDSGID = true;
        LockPersonality = true;
        # AF_UNIX is needed for the PostgreSQL socket.
        RestrictAddressFamilies = [
          "AF_INET"
          "AF_INET6"
          "AF_UNIX"
        ];
        SystemCallArchitectures = "native";
        SystemCallFilter = [
          "@system-service"
          "~@privileged @resources"
        ];
      };
    };
  };
}
