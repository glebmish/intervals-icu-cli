// internal/cmd/chats.go
package cmd

import (
	"fmt"

	"github.com/glebmish/intervals-icu-cli/internal/validate"
	"github.com/spf13/cobra"
)

var chatsCmd = &cobra.Command{
	Use:   "chats",
	Short: "Coach/athlete messaging",
}

func init() {
	chatsCmd.AddCommand(
		newChatsGetCmd(),
		newChatsMessagesCmd(),
		newChatsSendCmd(),
		newChatsDeleteMessageCmd(),
		newChatsMarkSeenCmd(),
	)
	rootCmd.AddCommand(chatsCmd)
}

func newChatsGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a chat by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			chatID, err := requireString(cmd, "chat-id")
			if err != nil {
				return err
			}
			if err := validate.PathParam("chat-id", chatID); err != nil {
				return err
			}
			params := map[string]string{"chatId": chatID}
			return doGet(cmd, "/api/v1/chats/{chatId}", params)
		},
	}
	cmd.Flags().String("chat-id", "", "Chat ID (required)")
	return cmd
}

func newChatsMessagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "messages",
		Short: "Get messages in a chat",
		RunE: func(cmd *cobra.Command, args []string) error {
			chatID, err := requireString(cmd, "chat-id")
			if err != nil {
				return err
			}
			if err := validate.PathParam("chat-id", chatID); err != nil {
				return err
			}
			params := map[string]string{"chatId": chatID}
			if v, _ := cmd.Flags().GetString("before-id"); v != "" {
				params["beforeId"] = v
			}
			if v, _ := cmd.Flags().GetInt("limit"); v > 0 {
				params["limit"] = fmt.Sprintf("%d", v)
			}
			return doGet(cmd, "/api/v1/chats/{chatId}/messages", params)
		},
	}
	cmd.Flags().String("chat-id", "", "Chat ID (required)")
	cmd.Flags().String("before-id", "", "Get messages before this message ID")
	cmd.Flags().Int("limit", 0, "Max number of messages")
	return cmd
}

func newChatsSendCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "Send a chat message",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			return doMutate(cmd, "POST", "/api/v1/chats/send-message", nil, jsonBody)
		},
	}
	return cmd
}

func newChatsDeleteMessageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-message",
		Short: "Delete a chat message",
		RunE: func(cmd *cobra.Command, args []string) error {
			chatID, err := requireString(cmd, "chat-id")
			if err != nil {
				return err
			}
			if err := validate.PathParam("chat-id", chatID); err != nil {
				return err
			}
			msgID, err := requireString(cmd, "msg-id")
			if err != nil {
				return err
			}
			if err := validate.PathParam("msg-id", msgID); err != nil {
				return err
			}
			params := map[string]string{"chatId": chatID, "msgId": msgID}
			return doDelete(cmd, "/api/v1/chats/{chatId}/messages/{msgId}", params, "message", msgID)
		},
	}
	cmd.Flags().String("chat-id", "", "Chat ID (required)")
	cmd.Flags().String("msg-id", "", "Message ID (required)")
	return cmd
}

func newChatsMarkSeenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mark-seen",
		Short: "Mark a message as seen",
		RunE: func(cmd *cobra.Command, args []string) error {
			chatID, err := requireString(cmd, "chat-id")
			if err != nil {
				return err
			}
			if err := validate.PathParam("chat-id", chatID); err != nil {
				return err
			}
			msgID, err := requireString(cmd, "msg-id")
			if err != nil {
				return err
			}
			if err := validate.PathParam("msg-id", msgID); err != nil {
				return err
			}
			params := map[string]string{"chatId": chatID, "msgId": msgID}
			return doMutate(cmd, "PUT", "/api/v1/chats/{chatId}/messages/{msgId}/seen", params, "")
		},
	}
	cmd.Flags().String("chat-id", "", "Chat ID (required)")
	cmd.Flags().String("msg-id", "", "Message ID (required)")
	return cmd
}
