package cmd

import (
	"fmt"
	"threshai/internal/core"
	"github.com/spf13/cobra"
)

var (
	topic  string
	length int
	tone   string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate content using ThreshAI",
	Long:  `Generate blogs, social posts, or other content with AI-powered precision.`,
	Run: func(cmd *cobra.Command, args []string) {
		request := core.ContentRequest{
			Topic:  topic,
			Length: length,
			Tone:   tone,
		}
		content := core.GenerateContent(request)
		fmt.Println("🚀 Generated Content:")
		fmt.Println(content)
	},
}

func init() {
	generateCmd.Flags().StringVarP(&topic, "topic", "t", "", "Topic for content generation")
	generateCmd.Flags().IntVarP(&length, "length", "l", 500, "Word count for generated content")
	generateCmd.Flags().StringVarP(&tone, "tone", "", "professional", "Tone of the content (e.g., professional, casual)")

	generateCmd.MarkFlagRequired("topic")
	RootCmd.AddCommand(generateCmd)
}