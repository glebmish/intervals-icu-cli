// internal/cmd/wellness.go
package cmd

import (
	"github.com/glebmish/intervals-icu-cli/internal/validate"
	"github.com/spf13/cobra"
)

var wellnessCmd = &cobra.Command{
	Use:   "wellness",
	Short: "Daily wellness records",
}

func init() {
	wellnessCmd.AddCommand(
		newWellnessGetCmd(),
		newWellnessUpdateCmd(),
		newWellnessUpdateCurrentCmd(),
		newWellnessUpdateBulkCmd(),
		newWellnessUploadCmd(),
		newWellnessListCmd(),
	)
	rootCmd.AddCommand(wellnessCmd)
}

func newWellnessGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get wellness record for a date",
		RunE: func(cmd *cobra.Command, args []string) error {
			date, err := requireString(cmd, "date")
			if err != nil {
				return err
			}
			if err := validate.DateParam("date", date); err != nil {
				return err
			}
			params := map[string]string{"date": date}
			return doGet(cmd, "/api/v1/athlete/{id}/wellness/{date}", params)
		},
	}
	cmd.Flags().String("date", "", "Date (YYYY-MM-DD, required)")
	return cmd
}

func newWellnessUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update wellness record for a date",
		RunE: func(cmd *cobra.Command, args []string) error {
			date, err := requireString(cmd, "date")
			if err != nil {
				return err
			}
			if err := validate.DateParam("date", date); err != nil {
				return err
			}
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"date": date}
			return doMutate(cmd, "PUT", "/api/v1/athlete/{id}/wellness/{date}", params, jsonBody)
		},
	}
	cmd.Flags().String("date", "", "Date (YYYY-MM-DD, required)")
	return cmd
}

func newWellnessUpdateCurrentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-current",
		Short: "Update wellness record for current day",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			return doMutate(cmd, "PUT", "/api/v1/athlete/{id}/wellness", nil, jsonBody)
		},
	}
	return cmd
}

func newWellnessUpdateBulkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-bulk",
		Short: "Bulk update wellness records",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			return doMutate(cmd, "PUT", "/api/v1/athlete/{id}/wellness-bulk", nil, jsonBody)
		},
	}
	return cmd
}

func newWellnessUploadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload wellness data",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{}
			if v, _ := cmd.Flags().GetBool("ignore-missing-fields"); v {
				params["ignoreMissingFields"] = "true"
			}
			return doMutate(cmd, "POST", "/api/v1/athlete/{id}/wellness", params, jsonBody)
		},
	}
	cmd.Flags().Bool("ignore-missing-fields", false, "Ignore missing fields in payload")
	return cmd
}

func newWellnessListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List wellness records",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := map[string]string{}
			if v, _ := cmd.Flags().GetString("oldest"); v != "" {
				params["oldest"] = v
			}
			if v, _ := cmd.Flags().GetString("newest"); v != "" {
				params["newest"] = v
			}
			if v, _ := cmd.Flags().GetString("cols"); v != "" {
				params["cols"] = v
			}
			if v, _ := cmd.Flags().GetString("fields-param"); v != "" {
				params["fields"] = v
			}
			ext, _ := cmd.Flags().GetString("ext")
			path := "/api/v1/athlete/{id}/wellness" + ext
			return doGet(cmd, path, params)
		},
	}
	cmd.Flags().String("ext", "", "Format extension for path (e.g. .csv)")
	cmd.Flags().String("oldest", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().String("newest", "", "End date (YYYY-MM-DD)")
	cmd.Flags().String("cols", "", "Columns to include")
	cmd.Flags().String("fields-param", "", "Fields query parameter")
	return cmd
}
