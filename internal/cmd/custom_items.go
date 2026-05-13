// internal/cmd/custom_items.go
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/glebmish/intervals-icu-cli/internal/api"
	"github.com/glebmish/intervals-icu-cli/internal/format"
	"github.com/glebmish/intervals-icu-cli/internal/validate"
	"github.com/spf13/cobra"
)

var customItemsCmd = &cobra.Command{
	Use:   "custom-items",
	Short: "Custom dashboard items",
}

func init() {
	customItemsCmd.AddCommand(
		newCustomItemsListCmd(),
		newCustomItemsGetCmd(),
		newCustomItemsCreateCmd(),
		newCustomItemsUpdateCmd(),
		newCustomItemsDeleteCmd(),
		newCustomItemsUploadImageCmd(),
		newCustomItemsUpdateIndexesCmd(),
	)
	rootCmd.AddCommand(customItemsCmd)
}

func newCustomItemsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List custom dashboard items",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doGet(cmd, "/api/v1/athlete/{id}/custom-item", nil)
		},
	}
	return cmd
}

func newCustomItemsGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a custom item by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			itemID, err := requireString(cmd, "item-id")
			if err != nil {
				return err
			}
			if err := validate.PathParam("item-id", itemID); err != nil {
				return err
			}
			params := map[string]string{"itemId": itemID}
			return doGet(cmd, "/api/v1/athlete/{id}/custom-item/{itemId}", params)
		},
	}
	cmd.Flags().String("item-id", "", "Custom item ID (required)")
	return cmd
}

func newCustomItemsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a custom dashboard item",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			return doMutate(cmd, "POST", "/api/v1/athlete/{id}/custom-item", nil, jsonBody)
		},
	}
	return cmd
}

func newCustomItemsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a custom dashboard item",
		RunE: func(cmd *cobra.Command, args []string) error {
			itemID, err := requireString(cmd, "item-id")
			if err != nil {
				return err
			}
			if err := validate.PathParam("item-id", itemID); err != nil {
				return err
			}
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"itemId": itemID}
			return doMutate(cmd, "PUT", "/api/v1/athlete/{id}/custom-item/{itemId}", params, jsonBody)
		},
	}
	cmd.Flags().String("item-id", "", "Custom item ID (required)")
	return cmd
}

func newCustomItemsDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a custom dashboard item",
		RunE: func(cmd *cobra.Command, args []string) error {
			itemID, err := requireString(cmd, "item-id")
			if err != nil {
				return err
			}
			if err := validate.PathParam("item-id", itemID); err != nil {
				return err
			}
			params := map[string]string{"itemId": itemID}
			return doDelete(cmd, "/api/v1/athlete/{id}/custom-item/{itemId}", params, "custom-item", itemID)
		},
	}
	cmd.Flags().String("item-id", "", "Custom item ID (required)")
	return cmd
}

func newCustomItemsUploadImageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload-image",
		Short: "Upload an image for a custom item",
		RunE: func(cmd *cobra.Command, args []string) error {
			itemID, err := requireString(cmd, "item-id")
			if err != nil {
				return err
			}
			if err := validate.PathParam("item-id", itemID); err != nil {
				return err
			}
			file, err := requireString(cmd, "file")
			if err != nil {
				return err
			}

			data, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("reading file %s: %w", file, err)
			}

			params := map[string]string{"itemId": itemID}
			c := api.FromContext(cmd.Context())

			dryRun, _ := cmd.Flags().GetBool("dry-run")
			if dryRun {
				fmt.Print(c.DryRun("POST", "/api/v1/athlete/{id}/custom-item/{itemId}/image", params, []byte(fmt.Sprintf("<binary file: %s, %d bytes>", file, len(data)))))
				return nil
			}

			resp, err := c.Do("POST", "/api/v1/athlete/{id}/custom-item/{itemId}/image", params, data)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("reading response: %w", err)
			}
			return format.Write(os.Stdout, body, fmtOpts(cmd))
		},
	}
	cmd.Flags().String("item-id", "", "Custom item ID (required)")
	cmd.Flags().String("file", "", "Path to image file (required)")
	return cmd
}

func newCustomItemsUpdateIndexesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-indexes",
		Short: "Update custom item indexes",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			return doMutate(cmd, "PUT", "/api/v1/athlete/{id}/custom-item-indexes", nil, jsonBody)
		},
	}
	return cmd
}
