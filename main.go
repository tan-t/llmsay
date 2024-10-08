package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

const (
	GPT4    = "gpt-4"
	GPT4O   = "gpt-4o"
	GPT35   = "gpt-3.5-turbo"
	CLAUDE2 = "claude-2"
)

var modelToProvider = map[string]string{
	GPT4:    "openai",
	GPT4O:   "openai",
	GPT35:   "openai",
	CLAUDE2: "anthropic",
}

var rootCmd = &cobra.Command{
	Use:  "llmsay",
	Long: "llmsay is a command line interface which allows you to access major llm models like GPT, Claude, Gemini",
	Args: cobra.RangeArgs(0, 1),
	// メイン処理
	RunE: func(cmd *cobra.Command, args []string) error {
		prompt := ""

		if len(args) == 1 {
			// 引数が1つ渡された場合
			prompt = args[0]
		}

		if stat, _ := os.Stdin.Stat(); (stat.Mode() & os.ModeCharDevice) == 0 {
			// 標準入力からデータがある場合
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				// だいたいの場合、プロンプトの「後ろ」に標準入力を入れることになるので
				prompt = prompt + scanner.Text()
			}
		}

		if prompt == "" {
			// 引数も標準入力もない場合、プロンプトを表示
			fmt.Print("Enter prompt: ")
			reader := bufio.NewReader(os.Stdin)
			prompt, _ = reader.ReadString('\n')
			prompt = strings.TrimSpace(prompt)
		}

		// フラグを取得
		model, err := cmd.Flags().GetString("model")
		if err != nil {
			return err
		}

		path, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}

		m := map[string]Config{}
		_, err = toml.DecodeFile(path, &m)
		if err != nil {
			return err
		}

		provider, ok := modelToProvider[model]
		if !ok {
			return fmt.Errorf("Unknown model: %s", model)
		}

		key := m[provider].Key
		client := getClient(provider, key)
		err = client.StreamCompletion(model, prompt)
		return err
	},
}

func getClient(provider, apiKey string) LLMClient {
	switch provider {
	case "openai":
		return NewOpenAIClient(apiKey)
	case "anthropic":
		return NewAnthropicClient(apiKey)
	default:
		log.Fatalf("Unknown provider: %s", provider)
		return nil
	}
}

type Config struct {
	Key string `toml:"key"`
}

var config = &cobra.Command{
	Use:  "configure",
	Long: "Configure api key for each provider",
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := cmd.Flags().GetString("provider")
		if err != nil {
			return err
		}

		k, err := cmd.Flags().GetString("key")
		if err != nil {
			return err
		}

		path, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}

		os.MkdirAll(filepath.Dir(path), os.ModePerm)

		m := map[string]Config{}

		if _, exErr := os.Stat(path); exErr == nil {
			_, err = toml.DecodeFile(path, &m)
			if err != nil {
				return err
			}
		}

		m[p] = Config{Key: k}

		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
		err = toml.NewEncoder(file).Encode(m)
		return err
	},
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	home, err := os.UserConfigDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	rootCmd.Flags().StringP("model", "m", "gpt-4o", "model name")
	rootCmd.Flags().StringP("file", "f", filepath.Join(home, "llmsay/config.toml"), "config file path")

	rootCmd.AddCommand(config)
	config.Flags().StringP("provider", "p", "", "[REQUIRED] provider name(openai,anthropic,gemini)")
	config.Flags().StringP("key", "k", "", "[REQUIRED] provider api key")
	config.Flags().StringP("file", "f", filepath.Join(home, "llmsay/config.toml"), "config file path")
	config.MarkFlagRequired("provider")
}
