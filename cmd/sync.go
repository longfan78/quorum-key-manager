package cmd

import (
	"context"

	auth "github.com/longfan78/quorum-key-manager/src/auth/entities"
	"github.com/longfan78/quorum-key-manager/src/auth/service/roles"
	"github.com/longfan78/quorum-key-manager/src/entities"
	storesservice "github.com/longfan78/quorum-key-manager/src/stores"
	manifeststores "github.com/longfan78/quorum-key-manager/src/stores/api/manifest"
	manifestvaults "github.com/longfan78/quorum-key-manager/src/vaults/api/manifest"
	"github.com/longfan78/quorum-key-manager/src/vaults/service/vaults"

	"github.com/longfan78/quorum-key-manager/cmd/flags"
	"github.com/longfan78/quorum-key-manager/src/infra/log/zap"
	manifestreader "github.com/longfan78/quorum-key-manager/src/infra/manifests/yaml"
	"github.com/longfan78/quorum-key-manager/src/infra/postgres/client"
	"github.com/longfan78/quorum-key-manager/src/stores/connectors/stores"
	"github.com/longfan78/quorum-key-manager/src/stores/database/postgres"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newSyncCommand() *cobra.Command {
	var logger *zap.Logger
	var storesService storesservice.Stores
	var mnfs map[string][]entities.Manifest
	var storeName string

	userInfo := auth.NewWildcardUser()

	syncCmd := &cobra.Command{
		Use:   "sync",
		Short: "Resource synchronization management tool",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			storeName = flags.GetStoreName(viper.GetViper())

			var err error

			// Infra dependencies
			if logger, err = getLogger(); err != nil {
				return err
			}

			postgresClient, err := client.New(flags.NewPostgresConfig(viper.GetViper()))
			if err != nil {
				return err
			}

			if mnfs, err = getManifests(ctx); err != nil {
				return err
			}

			// Instantiate register vaults
			roles := roles.New(logger)
			vaultService := vaults.New(roles, logger)
			if err := manifestvaults.NewVaultsHandler(vaultService).Register(ctx, mnfs[entities.VaultKind]); err != nil {
				return err
			}

			// Instantiate register stores
			storesService = stores.NewConnector(roles, postgres.New(logger, postgresClient), vaultService, logger)
			if err := manifeststores.NewStoresHandler(storesService).Register(ctx, mnfs[entities.StoreKind]); err != nil {
				return err
			}

			return nil
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			syncZapLogger(logger)
		},
	}

	flags.PGFlags(syncCmd.Flags())
	flags.SyncFlags(syncCmd.Flags())
	flags.ManifestFlags(syncCmd.Flags())

	syncSecretsCmd := &cobra.Command{
		Use:   "secrets",
		Short: "indexing secrets from remote vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := storesService.ImportSecrets(cmd.Context(), storeName, userInfo)
			if err != nil {
				cmd.SilenceUsage = true
				return err
			}
			return nil
		},
	}
	syncCmd.AddCommand(syncSecretsCmd)

	syncKeysCmd := &cobra.Command{
		Use:   "keys",
		Short: "indexing keys from remote vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := storesService.ImportKeys(cmd.Context(), storeName, userInfo)
			if err != nil {
				cmd.SilenceUsage = true
				return err
			}
			return nil
		},
	}
	syncCmd.AddCommand(syncKeysCmd)

	syncEthereumCmd := &cobra.Command{
		Use:   "ethereum",
		Short: "indexing ethereum accounts remote vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := storesService.ImportEthereum(cmd.Context(), storeName, userInfo)
			if err != nil {
				cmd.SilenceUsage = true
				return err
			}
			return nil
		},
	}
	syncCmd.AddCommand(syncEthereumCmd)

	return syncCmd
}

func getLogger() (*zap.Logger, error) {
	return zap.NewLogger(flags.NewLoggerConfig(viper.GetViper()))
}

func getManifests(ctx context.Context) (map[string][]entities.Manifest, error) {
	manifestReader, err := manifestreader.New(flags.NewManifestConfig(viper.GetViper()))
	if err != nil {
		return nil, err
	}

	return manifestReader.Load(ctx)
}
