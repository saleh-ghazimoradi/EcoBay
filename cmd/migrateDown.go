package cmd

import (
	"fmt"
	"github.com/saleh-ghazimoradi/EcoBay/slg"
	"github.com/saleh-ghazimoradi/EcoBay/utils"

	"github.com/spf13/cobra"
)

// migrateDownCmd represents the migrateDown command
var migrateDownCmd = &cobra.Command{
	Use:   "migrateDown",
	Short: "A brief description of your command",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("migrateDown called")
		db, err := utils.DBConnection(utils.DBMigrateDrop)
		if err != nil {
			slg.Logger.Error("Failed to connect to database", "error", err)
			return
		}
		if err := utils.DBMigrateDrop(db); err != nil {
			slg.Logger.Error("Failed to migrate drop", "error", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateDownCmd)
}
