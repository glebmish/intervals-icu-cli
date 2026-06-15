// internal/cmd/shared_events.go
package cmd

import (
	"fmt"
	"os"

	"github.com/glebmish/intervals-icu-cli/internal/api"
	"github.com/glebmish/intervals-icu-cli/internal/format"
	"github.com/glebmish/intervals-icu-cli/internal/validate"
	"github.com/spf13/cobra"
	"io"
)

var sharedEventsCmd = &cobra.Command{
	Use:   "shared-events",
	Short: "Public shared events",
}

func init() {
	sharedEventsCmd.AddCommand(
		newSharedEventsGetCmd(),
		newSharedEventsGetBySlugCmd(),
		newSharedEventsCreateCmd(),
		newSharedEventsUpdateCmd(),
		newSharedEventsDeleteCmd(),
		newSharedEventsUploadImageCmd(),
	)
	rootCmd.AddCommand(sharedEventsCmd)
}

func newSharedEventsGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a shared event by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			sharedEventID, err := requireString(cmd, "shared-event-id")
			if err != nil {
				return err
			}
			if err := validate.PathParam("shared-event-id", sharedEventID); err != nil {
				return err
			}
			params := map[string]string{"sharedEventId": sharedEventID}
			return doGet(cmd, "/api/v1/shared-event/{sharedEventId}", params)
		},
	}
	cmd.Flags().String("shared-event-id", "", "Shared event ID (required)")
	return cmd
}

func newSharedEventsGetBySlugCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-by-slug",
		Short: "Get a shared event by slug",
		RunE: func(cmd *cobra.Command, args []string) error {
			slug, err := requireString(cmd, "slug")
			if err != nil {
				return err
			}
			if err := validate.PathParam("slug", slug); err != nil {
				return err
			}
			params := map[string]string{"slug": slug}
			return doGet(cmd, "/api/v1/shared-events-by-slug/{slug}", params)
		},
	}
	cmd.Flags().String("slug", "", "Shared event slug (required)")
	return cmd
}

func newSharedEventsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a shared event",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{}
			if v, _ := cmd.Flags().GetString("link-to-event-id"); v != "" {
				params["linkToEventId"] = v
			}
			return doMutate(cmd, "POST", "/api/v1/shared-event", params, jsonBody)
		},
	}
	cmd.Flags().String("link-to-event-id", "", "Link to existing event ID")
	return cmd
}

func newSharedEventsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a shared event",
		RunE: func(cmd *cobra.Command, args []string) error {
			sharedEventID, err := requireString(cmd, "shared-event-id")
			if err != nil {
				return err
			}
			if err := validate.PathParam("shared-event-id", sharedEventID); err != nil {
				return err
			}
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"sharedEventId": sharedEventID}
			return doMutate(cmd, "PUT", "/api/v1/shared-event/{sharedEventId}", params, jsonBody)
		},
	}
	cmd.Flags().String("shared-event-id", "", "Shared event ID (required)")
	return cmd
}

func newSharedEventsDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a shared event",
		RunE: func(cmd *cobra.Command, args []string) error {
			sharedEventID, err := requireString(cmd, "shared-event-id")
			if err != nil {
				return err
			}
			if err := validate.PathParam("shared-event-id", sharedEventID); err != nil {
				return err
			}
			params := map[string]string{"sharedEventId": sharedEventID}
			return doDelete(cmd, "/api/v1/shared-event/{sharedEventId}", params, "shared-event", sharedEventID)
		},
	}
	cmd.Flags().String("shared-event-id", "", "Shared event ID (required)")
	return cmd
}

func newSharedEventsUploadImageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload-image",
		Short: "Upload an image for a shared event",
		RunE: func(cmd *cobra.Command, args []string) error {
			sharedEventID, err := requireString(cmd, "shared-event-id")
			if err != nil {
				return err
			}
			if err := validate.PathParam("shared-event-id", sharedEventID); err != nil {
				return err
			}
			file, data, err := requireUploadFile(cmd)
			if err != nil {
				return err
			}

			params := map[string]string{"sharedEventId": sharedEventID}
			c := api.FromContext(cmd.Context())

			dryRun, _ := cmd.Flags().GetBool("dry-run")
			if dryRun {
				out, err := c.DryRun("POST", "/api/v1/shared-event/{sharedEventId}/image", params, []byte(fmt.Sprintf("<binary file: %s, %d bytes>", file, len(data))))
				if err != nil {
					return err
				}
				return format.DryRunOutput(os.Stdout, out)
			}

			resp, err := c.Do("POST", "/api/v1/shared-event/{sharedEventId}/image", params, data)
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
	cmd.Flags().String("shared-event-id", "", "Shared event ID (required)")
	addUploadFlag(cmd, "Path to image file (required)")
	return cmd
}
