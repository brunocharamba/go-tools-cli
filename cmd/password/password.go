package password

import (
	"fmt"
	"math/rand"

	"github.com/spf13/cobra"
)

var PasswordCmd = &cobra.Command{
	Use: "password",
	Short: "Handles passwords",
	Long: `Handles passwords with customizable options. For example:

	got password`,
}

var generateCmd = &cobra.Command{
	Use: "generate",
	Short: "Generate random passwords",
	Long: `Generate random passwords with customizable options. For example:

	got password generate -l 12 -d -s // generates 12-sized password with digits and special chars`,
	Run: generatePassword,
}

func init() {
	// generate
	generateCmd.Flags().IntP("length", "l", 16, "Length of the generated password")
	generateCmd.Flags().BoolP("digits", "d", false, "Password contains digits")
	generateCmd.Flags().BoolP("special-chars", "s", false, "Password contains special characters")
	PasswordCmd.AddCommand(generateCmd)
}

func generatePassword(cmd *cobra.Command, args []string) {
	var length, _ = cmd.Flags().GetInt("length")
	var isDigits, _ = cmd.Flags().GetBool("digits")
	var isSpecialChars, _ = cmd.Flags().GetBool("special-chars")

	var charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	if isDigits {
		charset += "0123456789"
	}

	if isSpecialChars {
		charset += "!@#$%^&*()_+{}[]|;:,.<>?-="
	}

	var password = make([]byte, length)

	for i := range password {
		password[i] = charset[rand.Intn(len(charset))]
	}

	fmt.Println("Password generated:", string(password))
}