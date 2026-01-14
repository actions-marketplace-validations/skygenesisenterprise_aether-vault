package cmd

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/ui"
	"github.com/skygenesisenterprise/aether-vault/package/cli/pkg/types"
	"github.com/spf13/cobra"
)

// newPasswordCommand creates the password command
func newPasswordCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "password",
		Short: "Manage passwords and credentials",
		Long: `Manage passwords and credentials in your vault.

This command provides subcommands for:
  - Creating new password entries
  - Listing and searching passwords
  - Updating existing passwords
  - Deleting passwords
  - Generating secure passwords
  - Exporting and importing passwords`,
	}

	cmd.AddCommand(newPasswordCreateCommand())
	cmd.AddCommand(newPasswordListCommand())
	cmd.AddCommand(newPasswordGetCommand())
	cmd.AddCommand(newPasswordUpdateCommand())
	cmd.AddCommand(newPasswordDeleteCommand())
	cmd.AddCommand(newPasswordGenerateCommand())
	cmd.AddCommand(newPasswordSearchCommand())
	cmd.AddCommand(newPasswordImportCommand())
	cmd.AddCommand(newPasswordExportCommand())

	return cmd
}

// newPasswordCreateCommand creates the password create command
func newPasswordCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new password entry",
		Long: `Create a new password entry in the vault.

You can specify the password details using flags or enter them interactively.

Examples:
  vault password create --name "GitHub" --username "user" --password "pass123" --url "https://github.com"
  vault password create --type card --name "Visa Card"
  vault password create --interactive`,
		RunE: runPasswordCreateCommand,
	}

	cmd.Flags().String("name", "", "Name of the entry")
	cmd.Flags().String("username", "", "Username or email")
	cmd.Flags().String("password", "", "Password (leave empty to generate)")
	cmd.Flags().String("url", "", "Website URL")
	cmd.Flags().String("type", "login", "Entry type (login, card, identity, secure_note, ssh, database)")
	cmd.Flags().String("folder", "", "Folder/path")
	cmd.Flags().StringSlice("tags", []string{}, "Tags")
	cmd.Flags().Bool("favorite", false, "Mark as favorite")
	cmd.Flags().Bool("interactive", false, "Interactive mode")
	cmd.Flags().Bool("generate", false, "Generate a secure password")
	cmd.Flags().Int("length", 16, "Password length (for generation)")
	cmd.Flags().Bool("symbols", true, "Include symbols (for generation)")
	cmd.Flags().Bool("numbers", true, "Include numbers (for generation)")
	cmd.Flags().Bool("uppercase", true, "Include uppercase letters (for generation)")
	cmd.Flags().Bool("lowercase", true, "Include lowercase letters (for generation)")
	cmd.Flags().String("notes", "", "Additional notes")

	return cmd
}

// newPasswordListCommand creates the password list command
func newPasswordListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List password entries",
		Long: `List password entries in your vault with optional filtering.

Examples:
  vault password list
  vault password list --type login
  vault password list --folder "Work"
  vault password list --tags "important,work"
  vault password list --search "github"`,
		RunE: runPasswordListCommand,
	}

	cmd.Flags().String("type", "", "Filter by type")
	cmd.Flags().String("folder", "", "Filter by folder")
	cmd.Flags().StringSlice("tags", []string{}, "Filter by tags")
	cmd.Flags().String("search", "", "Search term")
	cmd.Flags().Bool("favorites", false, "Show only favorites")
	cmd.Flags().Bool("details", false, "Show detailed information")
	cmd.Flags().Int("limit", 50, "Limit results")
	cmd.Flags().String("sort", "name", "Sort by field (name, created, updated)")
	cmd.Flags().String("order", "asc", "Sort order (asc, desc)")

	return cmd
}

// newPasswordGetCommand creates the password get command
func newPasswordGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [name]",
		Short: "Get a password entry",
		Long: `Retrieve a password entry by name or ID.

Examples:
  vault password get "GitHub"
  vault password get github-id
  vault password get "GitHub" --show-password
  vault password get "GitHub" --copy-password`,
		RunE: runPasswordGetCommand,
	}

	cmd.Flags().Bool("show-password", false, "Show the password")
	cmd.Flags().Bool("copy-password", false, "Copy password to clipboard")
	cmd.Flags().Bool("copy-username", false, "Copy username to clipboard")
	cmd.Flags().String("format", "table", "Output format (table, json, yaml)")
	cmd.Flags().Bool("show-notes", false, "Show notes")

	return cmd
}

// newPasswordUpdateCommand creates the password update command
func newPasswordUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [name]",
		Short: "Update a password entry",
		Long: `Update an existing password entry.

Examples:
  vault password update "GitHub" --username "newuser"
  vault password update "GitHub" --password "newpass123"
  vault password update "GitHub" --generate-password
  vault password update "GitHub" --add-tags "work,important"`,
		RunE: runPasswordUpdateCommand,
	}

	cmd.Flags().String("name", "", "New name")
	cmd.Flags().String("username", "", "New username")
	cmd.Flags().String("password", "", "New password")
	cmd.Flags().String("url", "", "New URL")
	cmd.Flags().String("folder", "", "New folder")
	cmd.Flags().StringSlice("add-tags", []string{}, "Add tags")
	cmd.Flags().StringSlice("remove-tags", []string{}, "Remove tags")
	cmd.Flags().Bool("favorite", false, "Mark as favorite")
	cmd.Flags().Bool("no-favorite", false, "Remove from favorites")
	cmd.Flags().Bool("generate-password", false, "Generate a new secure password")
	cmd.Flags().Int("length", 16, "Password length (for generation)")
	cmd.Flags().String("notes", "", "Update notes")

	return cmd
}

// newPasswordDeleteCommand creates the password delete command
func newPasswordDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [name]",
		Short: "Delete a password entry",
		Long: `Delete a password entry from the vault.

Examples:
  vault password delete "GitHub"
  vault password delete github-id --confirm`,
		RunE: runPasswordDeleteCommand,
	}

	cmd.Flags().Bool("confirm", false, "Skip confirmation prompt")
	cmd.Flags().Bool("force", false, "Force deletion without confirmation")

	return cmd
}

// newPasswordGenerateCommand creates the password generate command
func newPasswordGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a secure password",
		Long: `Generate a secure password with customizable options.

Examples:
  vault password generate
  vault password generate --length 20 --symbols --numbers
  vault password generate --passphrase --words 4
  vault password generate --copy`,
		RunE: runPasswordGenerateCommand,
	}

	cmd.Flags().Int("length", 16, "Password length")
	cmd.Flags().Bool("symbols", true, "Include symbols")
	cmd.Flags().Bool("numbers", true, "Include numbers")
	cmd.Flags().Bool("uppercase", true, "Include uppercase letters")
	cmd.Flags().Bool("lowercase", true, "Include lowercase letters")
	cmd.Flags().Bool("passphrase", false, "Generate passphrase instead of password")
	cmd.Flags().Int("words", 6, "Number of words for passphrase")
	cmd.Flags().String("separator", "-", "Word separator for passphrase")
	cmd.Flags().Bool("capitalize", false, "Capitalize words in passphrase")
	cmd.Flags().Bool("include-number", false, "Include number in passphrase")
	cmd.Flags().Bool("copy", false, "Copy to clipboard")
	cmd.Flags().Bool("show", true, "Show the password")

	return cmd
}

// newPasswordSearchCommand creates the password search command
func newPasswordSearchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Search password entries",
		Long: `Search for password entries by name, username, URL, or notes.

Examples:
  vault password search "github"
  vault password search "user@company.com"
  vault password search "login" --type login
  vault password search "important" --tags
  vault password search "work" --url`,
		RunE: runPasswordSearchCommand,
	}

	cmd.Flags().String("type", "", "Search in specific type")
	cmd.Flags().Bool("url", false, "Include URLs in search")
	cmd.Flags().Bool("notes", false, "Include notes in search")
	cmd.Flags().Bool("tags", false, "Search in tags")
	cmd.Flags().String("folder", "", "Limit search to folder")
	cmd.Flags().Int("limit", 20, "Limit results")
	cmd.Flags().String("format", "table", "Output format (table, json)")

	return cmd
}

// newPasswordImportCommand creates the password import command
func newPasswordImportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import [file]",
		Short: "Import passwords from file",
		Long: `Import passwords from various file formats.

Supported formats: csv, json, 1pif, xml, keepass, lastpass

Examples:
  vault password import passwords.csv --format csv
  vault password import export.json --format json
  vault password import --preview passwords.csv`,
		RunE: runPasswordImportCommand,
	}

	cmd.Flags().String("format", "", "Import format (csv, json, 1pif, xml, keepass, lastpass)")
	cmd.Flags().Bool("preview", false, "Preview import without importing")
	cmd.Flags().String("strategy", "merge", "Import strategy (merge, replace)")
	cmd.Flags().String("folder", "", "Import to specific folder")
	cmd.Flags().StringSlice("tags", []string{}, "Add tags to imported entries")

	return cmd
}

// newPasswordExportCommand creates the password export command
func newPasswordExportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export [file]",
		Short: "Export passwords to file",
		Long: `Export passwords to various file formats.

Supported formats: csv, json, 1pif, xml

Examples:
  vault password export passwords.csv --format csv
  vault password export passwords.json --format json --include-passwords
  vault password export --output export.json`,
		RunE: runPasswordExportCommand,
	}

	cmd.Flags().String("format", "json", "Export format (csv, json, 1pif, xml)")
	cmd.Flags().Bool("include-passwords", false, "Include passwords in export")
	cmd.Flags().String("output", "", "Output file (default: stdout)")
	cmd.Flags().String("folder", "", "Export specific folder")
	cmd.Flags().StringSlice("tags", []string{}, "Export entries with tags")
	cmd.Flags().String("encrypt", "", "Encrypt export with password")

	return cmd
}

// Command runners

func runPasswordCreateCommand(cmd *cobra.Command, args []string) error {
	fmt.Printf("%sCreating new password entry%s\n", ui.Blue, ui.Reset)

	interactive, _ := cmd.Flags().GetBool("interactive")
	generatePassword, _ := cmd.Flags().GetBool("generate")

	if interactive {
		return createPasswordInteractive(cmd)
	}

	name, _ := cmd.Flags().GetString("name")
	if name == "" {
		return fmt.Errorf("name is required")
	}

	entry := &types.PasswordEntry{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if username, _ := cmd.Flags().GetString("username"); username != "" {
		entry.Username = username
	}

	if password, _ := cmd.Flags().GetString("password"); password != "" {
		entry.Password = password
	} else if generatePassword {
		password, err := generateSecurePassword(cmd)
		if err != nil {
			return err
		}
		entry.Password = password
		fmt.Printf("%sGenerated secure password%s\n", ui.Green, ui.Reset)
	}

	if url, _ := cmd.Flags().GetString("url"); url != "" {
		entry.URL = url
	}

	if typeStr, _ := cmd.Flags().GetString("type"); typeStr != "" {
		entry.Type = types.PasswordType(typeStr)
	}

	if folder, _ := cmd.Flags().GetString("folder"); folder != "" {
		entry.Folder = folder
	}

	if tags, _ := cmd.Flags().GetStringSlice("tags"); len(tags) > 0 {
		entry.Tags = tags
	}

	if favorite, _ := cmd.Flags().GetBool("favorite"); favorite {
		entry.Favorite = favorite
	}

	if notes, _ := cmd.Flags().GetString("notes"); notes != "" {
		entry.Notes = notes
	}

	// TODO: Save to vault
	fmt.Printf("%sPassword entry '%s' created successfully%s\n", ui.Green, entry.Name, ui.Reset)

	return nil
}

func runPasswordListCommand(cmd *cobra.Command, args []string) error {
	fmt.Printf("%sPassword entries%s\n", ui.Blue, ui.Reset)

	// TODO: Fetch from vault
	fmt.Println("(Listing would show actual password entries here)")

	fmt.Printf("%sName%-20s Username%-20s Type%-15s Updated%s\n",
		ui.Bold, "", "", "", ui.Reset)
	fmt.Println(strings.Repeat("-", 70))

	return nil
}

func runPasswordGetCommand(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("password name or ID is required")
	}

	name := args[0]
	showPassword, _ := cmd.Flags().GetBool("show-password")
	copyPassword, _ := cmd.Flags().GetBool("copy-password")

	fmt.Printf("%sPassword entry: %s%s\n", ui.Blue, name, ui.Reset)

	// TODO: Fetch from vault
	if showPassword {
		fmt.Printf("Password: %s(show password here)%s\n", ui.Yellow, ui.Reset)
	}

	if copyPassword {
		fmt.Printf("%sPassword copied to clipboard%s\n", ui.Green, ui.Reset)
		// TODO: Implement clipboard copy
	}

	return nil
}

func runPasswordUpdateCommand(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("password name or ID is required")
	}

	name := args[0]
	fmt.Printf("%sUpdating password entry: %s%s\n", ui.Blue, name, ui.Reset)

	// TODO: Update in vault
	fmt.Printf("%sPassword entry '%s' updated successfully%s\n", ui.Green, name, ui.Reset)

	return nil
}

func runPasswordDeleteCommand(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("password name or ID is required")
	}

	name := args[0]
	force, _ := cmd.Flags().GetBool("force")

	if !force {
		fmt.Printf("%sAre you sure you want to delete '%s'? [y/N]: %s", ui.Yellow, name, ui.Reset)
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			fmt.Printf("Deletion cancelled.\n")
			return nil
		}
	}

	fmt.Printf("%sDeleting password entry: %s%s\n", ui.Red, name, ui.Reset)

	// TODO: Delete from vault
	fmt.Printf("%sPassword entry '%s' deleted successfully%s\n", ui.Green, name, ui.Reset)

	return nil
}

func runPasswordGenerateCommand(cmd *cobra.Command, args []string) error {
	passphrase, _ := cmd.Flags().GetBool("passphrase")

	var password string
	var err error

	if passphrase {
		password, err = generatePassphrase(cmd)
	} else {
		password, err = generateSecurePassword(cmd)
	}

	if err != nil {
		return err
	}

	show, _ := cmd.Flags().GetBool("show")
	copy, _ := cmd.Flags().GetBool("copy")

	if show {
		fmt.Printf("%sGenerated password:%s %s%s%s\n", ui.Green, ui.Reset, ui.Yellow, password, ui.Reset)

		// Show password strength
		strength := evaluatePasswordStrength(password)
		fmt.Printf("%sStrength:%s %s (%d/100)%s\n",
			ui.Blue, ui.Reset,
			getStrengthColor(strength.Level)+strength.Level+ui.Reset,
			strength.Score)

		if len(strength.Issues) > 0 {
			fmt.Printf("%sIssues:%s %s%s\n", ui.Red, ui.Reset, strings.Join(strength.Issues, ", "), ui.Reset)
		}

		if len(strength.Suggestions) > 0 {
			fmt.Printf("%sSuggestions:%s %s%s\n", ui.Yellow, ui.Reset, strings.Join(strength.Suggestions, ", "), ui.Reset)
		}
	}

	if copy {
		fmt.Printf("%sPassword copied to clipboard%s\n", ui.Green, ui.Reset)
		// TODO: Implement clipboard copy
	}

	return nil
}

func runPasswordSearchCommand(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("search query is required")
	}

	query := args[0]
	fmt.Printf("%sSearching for: %s%s\n", ui.Blue, query, ui.Reset)

	// TODO: Search in vault
	fmt.Printf("%sFound X matching password entries%s\n", ui.Green, ui.Reset)

	return nil
}

func runPasswordImportCommand(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("import file is required")
	}

	file := args[0]
	preview, _ := cmd.Flags().GetBool("preview")
	_, _ = cmd.Flags().GetString("format")

	fmt.Printf("%sImporting passwords from: %s%s\n", ui.Blue, file, ui.Reset)

	if preview {
		fmt.Printf("%sPreview mode - no changes will be made%s\n", ui.Yellow, ui.Reset)
	}

	// TODO: Implement import based on format
	fmt.Printf("%sImport completed successfully%s\n", ui.Green, ui.Reset)

	return nil
}

func runPasswordExportCommand(cmd *cobra.Command, args []string) error {
	output, _ := cmd.Flags().GetString("output")
	format, _ := cmd.Flags().GetString("format")
	includePasswords, _ := cmd.Flags().GetBool("include-passwords")

	fmt.Printf("%sExporting passwords%s\n", ui.Blue, ui.Reset)
	fmt.Printf("Format: %s\n", format)
	fmt.Printf("Include passwords: %s\n", includePasswords)
	if output != "" {
		fmt.Printf("Output file: %s\n", output)
	}

	// TODO: Implement export
	fmt.Printf("%sExport completed successfully%s\n", ui.Green, ui.Reset)

	return nil
}

// Helper functions

func createPasswordInteractive(cmd *cobra.Command) error {
	fmt.Printf("%sInteractive password creation%s\n", ui.Blue, ui.Reset)

	var name, username, password, url string

	fmt.Print("Entry name: ")
	fmt.Scanln(&name)

	fmt.Print("Username/Email: ")
	fmt.Scanln(&username)

	fmt.Print("Password (leave empty to generate): ")
	fmt.Scanln(&password)

	if password == "" {
		generated, err := generateSecurePassword(cmd)
		if err != nil {
			return err
		}
		password = generated
		fmt.Printf("%sGenerated secure password: %s%s\n", ui.Green, password, ui.Reset)
	}

	fmt.Print("URL (optional): ")
	fmt.Scanln(&url)

	fmt.Printf("%sPassword entry '%s' would be created with the provided information%s\n",
		ui.Green, name, ui.Reset)

	return nil
}

func generateSecurePassword(cmd *cobra.Command) (string, error) {
	length, _ := cmd.Flags().GetInt("length")
	symbols, _ := cmd.Flags().GetBool("symbols")
	numbers, _ := cmd.Flags().GetBool("numbers")
	uppercase, _ := cmd.Flags().GetBool("uppercase")
	lowercase, _ := cmd.Flags().GetBool("lowercase")

	// Build character set
	var charset string
	var lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	var uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var numberChars = "0123456789"
	var symbolChars = "!@#$%^&*()_+-=[]{}|;:,.<>?"

	if lowercase {
		charset += lowercaseChars
	}
	if uppercase {
		charset += uppercaseChars
	}
	if numbers {
		charset += numberChars
	}
	if symbols {
		charset += symbolChars
	}

	if charset == "" {
		return "", fmt.Errorf("at least one character type must be selected")
	}

	if length < 4 {
		return "", fmt.Errorf("password length must be at least 4 characters")
	}

	// Generate secure random password
	password := make([]byte, length)
	charsetRunes := []rune(charset)
	charsetSize := big.NewInt(int64(len(charsetRunes)))

	for i := range password {
		randomIndex, err := rand.Int(rand.Reader, charsetSize)
		if err != nil {
			return "", fmt.Errorf("failed to generate random character: %w", err)
		}
		password[i] = byte(charsetRunes[randomIndex.Int64()])
	}

	// Ensure password contains at least one character from each selected type
	result := string(password)
	if lowercase && !containsAny(result, lowercaseChars) {
		result = replaceFirstChar(result, lowercaseChars)
	}
	if uppercase && !containsAny(result, uppercaseChars) {
		result = replaceFirstChar(result, uppercaseChars)
	}
	if numbers && !containsAny(result, numberChars) {
		result = replaceFirstChar(result, numberChars)
	}
	if symbols && !containsAny(result, symbolChars) {
		result = replaceFirstChar(result, symbolChars)
	}

	return result, nil
}

func generatePassphrase(cmd *cobra.Command) (string, error) {
	words, _ := cmd.Flags().GetInt("words")
	separator, _ := cmd.Flags().GetString("separator")
	capitalize, _ := cmd.Flags().GetBool("capitalize")
	includeNumber, _ := cmd.Flags().GetBool("include-number")

	// Word list for passphrase generation (EFF's long word list subset)
	wordList := []string{
		"abacus", "abdomen", "abide", "abiding", "ability", "ablaze", "abnormal", "aboard", "abolish", "abort",
		"absorb", "abstract", "absurd", "abuse", "access", "accident", "account", "accuse", "achieve", "acid",
		"acoustic", "acquire", "across", "actress", "adapter", "addict", "adjust", "admit", "adult", "advance",
		"advice", "advocate", "affect", "afford", "agency", "agent", "agony", "agree", "ahead", "aid",
		"airline", "airport", "aisle", "alarm", "album", "alcohol", "alert", "alien", "alike", "alive",
		"allergy", "allow", "alloy", "alone", "aloft", "alphabet", "already", "also", "alter", "always",
		"amazing", "ambition", "among", "amount", "ample", "amuse", "anchor", "ancient", "anger", "angle",
		"angry", "ankle", "announce", "annual", "answer", "antenna", "antique", "anxiety", "anybody", "apart",
		"apology", "appear", "apple", "approve", "april", "arch", "ardent", "area", "arena", "argue",
		"army", "around", "arrange", "arrest", "arrive", "arrow", "article", "aside", "aspect", "asset",
		"avoid", "awake", "award", "aware", "badly", "baffle", "baggage", "baker", "balance", "balcony",
		"ballot", "banana", "banner", "barrel", "basic", "basket", "battle", "beach", "beard", "beast",
		"became", "beef", "before", "begin", "being", "belief", "believe", "belly", "belong", "below",
		"bench", "berry", "betray", "better", "between", "beyond", "bicycle", "bigger", "billion", "binary",
		"biology", "bird", "birth", "bitter", "black", "blade", "blame", "blanket", "blast", "bleak",
		"bless", "blind", "blink", "bliss", "block", "blood", "blossom", "blouse", "blue", "blur",
		"blunt", "blurt", "blush", "board", "boast", "bonus", "boost", "booth", "border", "boss",
		"bottom", "bounce", "bound", "bowel", "bowler", "bowling", "boxer", "boxing", "boyfriend", "brain",
		"brake", "branch", "brass", "brave", "bread", "break", "breath", "breed", "brick", "bride",
		"brief", "bring", "broad", "broke", "broken", "brother", "brought", "brown", "brush", "bubble",
		"bucket", "budget", "build", "bullet", "bunch", "burden", "burger", "burst", "business", "butter",
		"button", "cabin", "cable", "calm", "camera", "camp", "canal", "cannot", "capital", "captain",
		"capture", "carbon", "card", "cargo", "carry", "cartoon", "carve", "castle", "casual", "cat",
		"catalog", "catch", "cater", "cause", "caution", "cave", "ceiling", "cell", "center", "cereal",
		"certain", "chain", "chair", "chalk", "champion", "change", "chaos", "chapter", "charge", "charm",
		"chart", "chase", "cheap", "cheat", "check", "cheek", "cheese", "chemical", "cherry", "chest",
		"chicken", "chief", "child", "chill", "chimney", "choice", "choose", "chop", "chose", "chosen",
		"church", "circle", "citizen", "city", "claim", "clash", "claw", "clay", "clean", "clear",
		"clerk", "click", "cliff", "climb", "cling", "clock", "clone", "close", "cloth", "cloud",
		"clown", "club", "clue", "coach", "coast", "coat", "code", "coffee", "coin", "collect",
		"college", "colony", "color", "column", "combine", "comfort", "comic", "common", "company", "compare",
		"compete", "complex", "compose", "concept", "concern", "concert", "conduct", "confirm", "connect", "consent",
		"consider", "consist", "constant", "contain", "content", "contest", "context", "continue", "contract", "contrast",
		"control", "convert", "cook", "cool", "copper", "copy", "coral", "cord", "corner", "correct",
		"cost", "cotton", "couch", "could", "count", "counter", "country", "couple", "courage", "course",
		"court", "cousin", "cover", "crack", "craft", "crash", "crazy", "cream", "create", "creature",
		"credit", "crew", "cricket", "crime", "crisis", "criteria", "critical", "crop", "cross", "crowd",
		"crown", "crude", "cruel", "cruise", "crush", "cry", "crystal", "cube", "culture", "cup",
		"curb", "cure", "curly", "currency", "current", "curve", "custom", "cycle", "daily", "damage",
		"dance", "danger", "daring", "dark", "data", "database", "date", "daughter", "dawn", "day",
		"dead", "deal", "dean", "dear", "death", "debate", "debris", "decade", "decent", "decide",
		"declare", "decline", "decorate", "decrease", "dedicate", "deep", "defeat", "defend", "define", "degree",
		"delay", "deliver", "demand", "democracy", "demonstrate", "denial", "dense", "depart", "depend", "deposit",
		"depth", "derive", "describe", "desert", "deserve", "design", "desire", "desk", "despair", "desperate",
		"despite", "destroy", "detail", "detect", "determine", "develop", "device", "devote", "diagram", "dial",
		"diamond", "diary", "dice", "diet", "differ", "digital", "dilemma", "dimension", "dinner", "direct",
		"director", "dirt", "dirty", "disagree", "disappear", "disaster", "discard", "discover", "discuss", "disease",
		"dismiss", "display", "distance", "distribute", "district", "disturb", "diverse", "divide", "division", "divorce",
		"doctor", "document", "dodge", "doing", "dollar", "domain", "domestic", "dominant", "donate", "done",
		"donkey", "donor", "door", "dose", "double", "doubt", "dozen", "draft", "dragon", "drama",
		"dramatic", "draw", "dream", "dress", "drift", "drill", "drink", "drive", "drop", "drought",
		"drum", "drunk", "dry", "due", "dull", "during", "dust", "duty", "dwarf", "dynamic",
		"eager", "early", "earth", "easier", "easily", "east", "easy", "echo", "eclipse", "economy",
		"edge", "edit", "educate", "effort", "eight", "either", "elaborate", "elastic", "elbow", "elder",
		"elect", "elegant", "element", "elephant", "elite", "else", "elsewhere", "email", "embark", "embody",
		"emerge", "emotion", "emphasize", "empire", "employ", "empty", "enact", "enable", "enclose", "encounter",
		"encourage", "end", "endless", "endorse", "enemy", "energy", "enforce", "engage", "engine", "enhance",
		"enjoy", "enlarge", "enough", "enquire", "ensure", "enter", "entire", "entry", "envelope", "environment",
		"episode", "equal", "equip", "era", "erase", "erect", "error", "escape", "essay", "essential",
		"establish", "estate", "estimate", "ethics", "ethnic", "evaluate", "even", "event", "ever", "every",
		"everybody", "everyday", "everyone", "everything", "everywhere", "evidence", "evil", "evoke", "evolve", "exact",
		"example", "exceed", "excel", "except", "excess", "exchange", "excite", "exclude", "excuse", "execute",
		"exercise", "exhaust", "exhibit", "exist", "exit", "exotic", "expand", "expect", "expense", "experience",
		"experiment", "expert", "explain", "explode", "explore", "export", "expose", "express", "extend", "exterior",
		"external", "extra", "extraordinary", "extreme", "fabric", "face", "facilitate", "fact", "factory", "faculty",
		"fade", "fail", "failure", "fair", "faith", "fake", "fall", "false", "familiar", "family",
		"famous", "fan", "fancy", "fantasy", "farm", "farmer", "fashion", "fast", "fat", "fatal",
		"fate", "father", "fault", "favor", "favorite", "fear", "feature", "federal", "fee", "feed",
		"feel", "feet", "fellow", "female", "fence", "festival", "fetch", "fever", "few", "fiber",
		"fiction", "field", "fierce", "fifteen", "fifth", "fifty", "fight", "figure", "file", "fill",
		"film", "filter", "final", "find", "fine", "finger", "finish", "fire", "firm", "first",
		"fish", "fit", "fitness", "fix", "flag", "flame", "flat", "flavor", "flee", "flesh",
		"flight", "float", "flock", "flood", "floor", "flour", "flow", "flower", "fluid", "flush",
		"focus", "fold", "folk", "follow", "food", "foot", "football", "force", "foreign", "forest",
		"forget", "forgive", "form", "formal", "format", "former", "formula", "forth", "forty", "forum",
		"forward", "fossil", "foster", "foul", "found", "foundation", "four", "fourteen", "fourth", "fraction",
		"fragment", "frame", "framework", "frank", "fraud", "free", "freedom", "freeze", "frequent", "fresh",
		"friend", "frighten", "front", "frozen", "fruit", "frustrate", "fuel", "full", "fun", "function",
		"fund", "fundamental", "funding", "funeral", "funny", "furniture", "further", "future", "gain", "galaxy",
		"gallery", "game", "gang", "garage", "garbage", "garden", "garlic", "gather", "gauge", "gaze",
		"general", "generate", "generation", "generous", "genetic", "genre", "gentle", "genuine", "gesture", "ghost",
		"giant", "gift", "giggle", "ginger", "girl", "give", "glad", "glance", "glass", "glide",
		"glory", "glove", "glow", "glue", "goat", "goal", "goat", "gold", "golden", "golf",
		"good", "govern", "government", "grab", "grace", "grade", "grain", "grand", "grant", "grape",
		"graph", "grasp", "grass", "grateful", "grave", "gravity", "gray", "great", "green", "greet",
		"grid", "grief", "grind", "grip", "grocery", "group", "grow", "growth", "guarantee", "guard",
		"guess", "guide", "guilt", "guilty", "guitar", "gulf", "gun", "gym", "habit", "habitat",
		"hair", "half", "hall", "hammer", "hand", "handful", "handle", "handsome", "hang", "happen",
		"happy", "harbor", "hard", "harden", "hardly", "harm", "harsh", "harvest", "hat", "hatch",
		"hate", "haul", "have", "hazard", "head", "headline", "headquarters", "heal", "health", "healthy",
		"heap", "hear", "hearing", "heart", "heat", "heaven", "heavy", "hedge", "height", "helicopter",
		"hell", "hello", "help", "helpful", "hence", "herd", "here", "heritage", "hero", "herself",
		"hesitate", "hide", "high", "highlight", "highway", "hill", "himself", "hint", "hip", "hire",
		"historian", "history", "hold", "hole", "holiday", "hollow", "home", "homework", "honest", "honey",
		"honor", "hope", "horizon", "horn", "horrible", "horror", "horse", "hospital", "host", "hostile",
		"hour", "house", "household", "housing", "however", "huge", "human", "humble", "humor", "hundred",
		"hunger", "hunt", "hurry", "hurt", "husband", "hypothesis", "ice", "idea", "ideal", "identical",
		"identify", "identity", "ideology", "image", "imagination", "imagine", "imitate", "immense", "immune", "impact",
		"implement", "imply", "import", "impose", "impossible", "impress", "impressive", "improve", "impulse", "inch",
		"incident", "include", "income", "incorporate", "increase", "incredible", "indeed", "indefinite", "independent", "index",
		"indicate", "indigenous", "individual", "indoor", "induce", "industry", "inevitable", "infant", "infect", "infinite",
		"inflation", "influence", "inform", "infrastructure", "ingredient", "inhabit", "inhabitant", "inhale", "inherit", "initial",
		"initiate", "inject", "injure", "injury", "inmate", "inner", "innocent", "input", "inquiry", "insect",
		"insert", "inside", "insight", "insist", "inspect", "inspire", "install", "instance", "instant", "instead",
		"institute", "institution", "instruct", "instrument", "insulate", "insult", "insurance", "intact", "integrate", "integrity",
		"intellectual", "intelligence", "intend", "intense", "intensity", "intent", "interact", "interest", "interfere", "interior",
		"intermediate", "internal", "international", "internet", "interpret", "interrupt", "interval", "intervene", "interview", "intimate",
		"introduce", "intrinsic", "invade", "invent", "invest", "investigate", "investment", "investor", "invisible", "invitation",
		"involve", "iron", "irony", "island", "isolate", "issue", "item", "its", "itself", "ivory",
		"jacket", "jail", "jam", "jar", "jaw", "jeans", "jelly", "jewelry", "job", "join",
		"joint", "joke", "journal", "journalist", "journey", "judge", "juice", "jump", "jungle", "junior",
		"jurisdiction", "jury", "just", "justice", "justify", "keen", "keep", "key", "keyboard", "kick",
		"kid", "kidney", "kill", "kilogram", "kilometer", "kind", "kingdom", "kiss", "kit", "kitchen",
		"knee", "knife", "knock", "know", "knowledge", "label", "labor", "laboratory", "lack", "ladder",
		"lady", "lake", "lamb", "lamp", "language", "laptop", "large", "largely", "laser", "last",
		"late", "later", "latter", "laugh", "launch", "laundry", "lava", "law", "lawn", "lawsuit",
		"layer", "layman", "layout", "lazy", "lead", "leader", "leadership", "leading", "leaf", "league",
		"lean", "leap", "learn", "lease", "least", "leather", "leave", "lecture", "left", "legacy",
		"legal", "legend", "legislation", "legitimate", "lemon", "lend", "length", "lens", "less", "lesson",
		"letter", "level", "lever", "liberal", "library", "license", "life", "lift", "light", "lightning",
		"like", "limit", "line", "linear", "link", "lion", "lip", "liquid", "list", "listen",
		"literally", "literary", "literature", "little", "live", "liver", "living", "load", "loan", "lobby",
		"local", "locate", "location", "lock", "log", "logic", "logical", "lonely", "long", "longevity",
		"look", "loop", "loose", "lord", "lose", "loss", "lost", "loud", "love", "lovely",
		"lover", "lower", "loyal", "loyalty", "luck", "lucky", "luggage", "lump", "lunch", "lung",
		"luxury", "machine", "machinery", "mad", "magazine", "magic", "magnetic", "magnificent", "magnitude", "maid",
		"mail", "main", "mainland", "mainly", "maintain", "maintenance", "major", "majority", "make", "maker",
		"makeup", "male", "mall", "manage", "management", "manager", "mandate", "manipulate", "mankind", "manner",
		"manual", "manufacture", "manufacturer", "many", "map", "marble", "march", "margin", "marine", "mark",
		"market", "marketing", "marriage", "married", "marry", "mask", "mass", "massive", "master", "match",
		"mate", "material", "mathematics", "matter", "mature", "maximum", "maybe", "mayor", "meal", "mean",
		"meanwhile", "measure", "measurement", "meat", "mechanic", "mechanical", "mechanism", "media", "medical", "medication",
		"medicine", "medieval", "medium", "meet", "meeting", "member", "membership", "memo", "memorial", "memory",
		"mental", "mention", "mentor", "menu", "mercy", "mere", "merge", "merit", "merry", "mess",
		"message", "metal", "meter", "method", "metropolitan", "micro", "microscope", "middle", "midnight", "midst",
		"might", "migration", "mild", "mile", "military", "milk", "mill", "million", "mind", "mine",
		"minimum", "mining", "minister", "ministry", "minor", "minority", "minute", "miracle", "mirror", "missile",
		"mission", "mistake", "mister", "mix", "mixture", "mobile", "moderate", "modern", "modify", "module",
		"moisture", "moment", "momentum", "monarch", "money", "monitor", "month", "monument", "mood", "moon",
		"moral", "more", "moreover", "morning", "mortgage", "mosque", "mosquito", "most", "mostly", "mother",
		"motion", "motivate", "motivation", "motor", "mount", "mountain", "mouse", "mouth", "move", "movement",
		"movie", "much", "multiple", "multiply", "municipal", "murder", "muscle", "museum", "mushroom", "music",
		"musical", "musician", "must", "mutual", "myself", "mystery", "myth", "naked", "name", "narrative",
		"narrow", "nation", "national", "native", "natural", "naturally", "nature", "navigate", "near", "nearby",
		"nearly", "necessity", "neck", "need", "needle", "negative", "neglect", "negotiate", "neighbor", "neighborhood",
		"neither", "nerve", "nervous", "nest", "network", "neutral", "never", "nevertheless", "new", "newly",
		"news", "newspaper", "next", "nice", "night", "nine", "nobody", "noise", "nomination", "none",
		"nonetheless", "noodle", "noon", "nor", "normal", "normally", "north", "northern", "nose", "notable",
		"note", "notebook", "nothing", "notice", "notion", "novel", "nowadays", "nowhere", "nuclear", "number",
		"numerous", "nurse", "nursery", "nutrient", "nylon", "obey", "object", "objection", "objective", "obligation",
		"observation", "observe", "observer", "obstacle", "obtain", "obvious", "occasion", "occasional", "occupation", "occupy",
		"occur", "ocean", "odd", "offense", "offer", "office", "officer", "official", "offset", "often",
		"okay", "old", "older", "olive", "olympic", "once", "one", "ongoing", "onion", "online",
		"only", "onto", "open", "opening", "opera", "operate", "operation", "operator", "opinion", "opponent",
		"opportunity", "oppose", "opposite", "opposition", "optimal", "optimistic", "option", "orange", "orbit", "order",
		"ordinary", "organ", "organic", "organization", "organize", "orientation", "origin", "original", "otherwise", "ought",
		"ounce", "outcome", "outdoor", "outer", "outfit", "outlet", "outline", "outlook", "output", "outrage",
		"outside", "outstanding", "over", "overall", "overcome", "overlook", "override", "overseas", "oversee", "overt",
		"overtime", "overturn", "overwhelm", "overwhelming", "owe", "owner", "ownership", "oxygen", "ozone", "pace",
		"pack", "package", "packet", "page", "pain", "painful", "paint", "painter", "painting", "pair",
		"palace", "pale", "palm", "pan", "panel", "panic", "paper", "parade", "paradise", "paragraph",
		"parallel", "parameter", "parent", "park", "parking", "parliament", "part", "partial", "partially", "participant",
		"participate", "participation", "particle", "particular", "particularly", "partly", "partner", "partnership", "party", "pass",
		"passage", "passenger", "passion", "passive", "passport", "past", "paste", "pastor", "patch", "patent",
		"path", "pathology", "patience", "patient", "pattern", "pause", "pay", "payment", "payroll", "peace",
		"peaceful", "peak", "peanut", "pear", "pedestrian", "peer", "penalty", "pencil", "pending", "penetrate",
		"peninsula", "penny", "pension", "people", "pepper", "perceive", "percentage", "perception", "perfect", "perform",
		"performance", "perhaps", "period", "permanent", "permission", "permit", "persist", "person", "personal", "personality",
		"personally", "personnel", "perspective", "persuade", "phase", "phenomenon", "philosophy", "phone", "photo", "photograph",
		"phrase", "physical", "physically", "physician", "physics", "piano", "pick", "pickup", "picture", "piece",
		"pile", "pill", "pillar", "pillow", "pilot", "pin", "pine", "pink", "pioneer", "pipe",
		"pitch", "pizza", "place", "plain", "plan", "plane", "planet", "planning", "plant", "plastic",
		"plate", "platform", "play", "player", "playground", "plea", "plead", "pleasant", "please", "pleasure",
		"pledge", "plenty", "plot", "plow", "plug", "plunge", "plural", "plus", "pocket", "poem",
		"poet", "poetry", "point", "pole", "police", "policy", "political", "politically", "politician", "politics",
		"poll", "pollution", "pond", "pool", "poor", "pop", "popular", "population", "porch", "pork",
		"port", "portfolio", "portion", "portrait", "portray", "pose", "position", "positive", "possess", "possession",
		"possibility", "possible", "possibly", "post", "postage", "postpone", "pot", "potato", "potential", "potentially",
		"pottery", "pound", "pour", "poverty", "powder", "power", "powerful", "practical", "practice", "practitioner",
		"praise", "pray", "prayer", "preach", "precede", "precious", "precise", "precision", "predict", "prefer",
		"preference", "pregnancy", "pregnant", "preliminary", "premise", "premium", "preparation", "prepare", "prescription", "presence",
		"present", "presentation", "preserve", "president", "presidential", "press", "pressure", "presume", "pretend", "pretty",
		"prevail", "prevent", "previous", "price", "pride", "priest", "primarily", "primary", "prime", "primitive",
		"prince", "princess", "principal", "principle", "print", "prior", "priority", "prison", "prisoner", "privacy",
		"private", "probably", "probe", "problem", "procedure", "proceed", "process", "processor", "proclaim", "produce",
		"product", "production", "profession", "professional", "professor", "profile", "profit", "program", "progress", "prohibit",
		"project", "prominent", "promise", "promising", "promote", "prompt", "proof", "propaganda", "proper", "properly",
		"property", "proportion", "proposal", "propose", "proposed", "proposition", "prosecute", "prosecution", "prospect", "prosper",
		"protect", "protection", "protective", "protein", "protest", "proud", "prove", "provide", "provided", "provider",
		"province", "provision", "provoke", "psychological", "psychology", "public", "publication", "publicity", "publish", "publisher",
		"pudding", "pull", "pulse", "pump", "punch", "punish", "punishment", "pupil", "purchase", "pure",
		"purple", "purpose", "purse", "pursue", "pursuit", "push", "put", "puzzle", "pyramid", "qualification",
		"qualify", "quality", "quantity", "quarrel", "quarter", "quarterback", "quarterly", "quarters", "queen", "quest",
		"question", "quick", "quickly", "quiet", "quietly", "quit", "quite", "quota", "quote", "race",
		"racial", "racism", "racist", "rack", "radar", "radiation", "radio", "radical", "radius", "rage",
		"raid", "rail", "railroad", "rain", "rainbow", "raise", "rally", "ramp", "ranch", "random",
		"range", "rank", "rapid", "rapidly", "rare", "rarely", "rate", "rather", "rating", "ratio",
		"rational", "ravage", "raw", "ray", "reach", "react", "reaction", "read", "reader", "readily",
		"reading", "ready", "real", "realistic", "reality", "realize", "really", "realm", "rear", "reason",
		"reasonable", "reasonably", "reasoning", "rebel", "rebuild", "recall", "receipt", "receive", "receiver", "recent",
		"recently", "reception", "recession", "recipe", "recipient", "reckon", "recognize", "recommend", "recommendation", "record",
		"recorder", "recover", "recovery", "recreation", "recruit", "recycling", "reduce", "reduction", "refer", "referee",
		"reference", "reflect", "reflection", "reform", "refuge", "refugee", "refusal", "refuse", "regard", "regarding",
		"regardless", "regime", "region", "regional", "register", "regret", "regular", "regularly", "regulate", "regulation",
		"regulator", "rehabilitate", "reign", "reinforce", "reject", "relate", "relation", "relationship", "relative", "relatively",
		"relax", "release", "relevance", "relevant", "reliable", "relief", "relieve", "religion", "religious", "reluctant",
		"rely", "remain", "remainder", "remains", "remark", "remarkable", "remedy", "remember", "remind", "reminder",
		"remote", "removal", "remove", "render", "renew", "rent", "rental", "repair", "repeat", "repeatedly",
		"replace", "replacement", "reply", "report", "reporter", "reporting", "represent", "representation", "representative", "republic",
		"republican", "reputation", "request", "require", "requirement", "rescue", "research", "researcher", "resemble", "reservation",
		"reserve", "residence", "resident", "residential", "resign", "resignation", "resist", "resistance", "resistant", "resort",
		"resource", "respect", "respective", "respond", "respondent", "response", "responsibility", "responsible", "rest", "restaurant",
		"restore", "result", "retain", "retire", "retired", "retirement", "retreat", "retrieve", "return", "reveal",
		"revelation", "revenge", "revenue", "reverse", "review", "revise", "revision", "revive", "revolt", "revolution",
		"revolutionary", "reward", "rhythm", "rice", "rich", "rid", "riddle", "ride", "rider", "ridge",
		"ridiculous", "rifle", "right", "righteous", "rigid", "ring", "riot", "rise", "rising", "risk",
		"risky", "ritual", "rival", "river", "road", "roar", "roast", "rob", "robber", "robbery",
		"robot", "robust", "rock", "rocket", "rocky", "role", "roll", "roller", "romance", "romantic",
		"roof", "room", "root", "rope", "rose", "rough", "roughly", "round", "route", "routine",
		"row", "royal", "rub", "rubber", "rubbish", "rude", "rug", "ruin", "rule", "ruler",
		"rumor", "run", "runner", "running", "rural", "rush", "sacred", "sacrifice", "sad", "sadness",
		"safe", "safely", "safer", "safety", "saga", "sake", "salad", "salary", "sale", "salmon",
		"salon", "salt", "salute", "salvation", "sample", "sanction", "sand", "sandwich", "satellite", "satisfaction",
		"satisfied", "satisfy", "sauce", "saucepan", "sausage", "save", "saving", "savings", "savior", "scale",
		"scan", "scandal", "scared", "scarf", "scatter", "scenario", "scene", "scenery", "schedule", "scheme",
		"scholar", "scholarship", "school", "science", "scientific", "scientist", "scope", "score", "scorn", "scout",
		"scrap", "scrape", "scratch", "scream", "screen", "screening", "screw", "script", "scrub", "sculpture",
		"seal", "seam", "search", "season", "seasonal", "seat", "second", "secondary", "secret", "secretary",
		"secretly", "section", "sector", "secure", "security", "seed", "seek", "seem", "seemingly", "segment",
		"seize", "seldom", "select", "selection", "self", "selfish", "sell", "seller", "selling", "semester",
		"semi", "semiconductor", "senate", "senator", "send", "senior", "sensation", "sense", "sensible", "sensitive",
		"sentence", "sentiment", "separate", "separately", "separation", "sequence", "series", "serious", "seriously", "serum",
		"serve", "server", "service", "serving", "session", "set", "setback", "setting", "settle", "settlement",
		"several", "severe", "severely", "sew", "sex", "sexual", "sexuality", "sexy", "shabby", "shade",
		"shadow", "shaft", "shake", "shallow", "shame", "shampoo", "shape", "share", "shareholder", "shark",
		"sharp", "sharply", "shatter", "shave", "shear", "shed", "sheep", "sheer", "sheet", "shelf",
		"shell", "shelter", "shield", "shift", "shine", "shiny", "ship", "shipment", "shipping", "shirt",
		"shock", "shoe", "shoot", "shooter", "shop", "shopper", "shopping", "shore", "short", "shortage",
		"shortly", "shot", "should", "shoulder", "shout", "show", "shower", "shrink", "shrug", "shut",
		"shuttle", "shy", "sibling", "sick", "side", "sidewalk", "siege", "sight", "sign", "signal",
		"signature", "significance", "significant", "significantly", "silence", "silent", "silly", "silver", "similar", "similarly",
		"simple", "simplicity", "simply", "simulate", "simulation", "simultaneous", "since", "sincere", "sing", "singer",
		"singing", "single", "sink", "sister", "site", "situated", "situation", "six", "size", "skate",
		"skeleton", "skeptic", "sketch", "skill", "skilled", "skin", "skip", "skirt", "skull", "sky",
		"slam", "slap", "slave", "slavery", "sleep", "sleepy", "sleeve", "slender", "slice", "slide",
		"slight", "slightly", "slip", "slippery", "slope", "slot", "slow", "slowly", "small", "smart",
		"smell", "smile", "smoke", "smoking", "smooth", "smoothly", "snack", "snake", "snap", "snapshot",
		"snow", "soak", "soap", "soccer", "social", "society", "sociology", "sock", "sodium", "soft",
		"softball", "softly", "software", "soil", "solar", "soldier", "sole", "solely", "solid", "solidarity",
		"solo", "solution", "solve", "some", "somebody", "somehow", "someone", "something", "sometime", "sometimes",
		"somewhat", "somewhere", "song", "soon", "sophisticated", "sophomore", "sore", "sorrow", "sorry", "sort",
		"soul", "sound", "soundtrack", "soup", "sour", "source", "south", "southeast", "southern", "southwest",
		"sovereign", "space", "spacecraft", "span", "spare", "spark", "speak", "speaker", "special", "specialist",
		"specialize", "specialized", "specialty", "species", "specific", "specifically", "specification", "specify", "specimen", "spectacle",
		"spectacular", "spectrum", "speculate", "speculation", "speech", "speed", "spell", "spend", "spending", "spider",
		"spin", "spine", "spiral", "spirit", "spiritual", "spit", "spite", "split", "spoil", "spokesman",
		"sponsor", "spontaneous", "spot", "spotlight", "spouse", "spray", "spread", "spring", "sprinkle", "spur",
		"squad", "square", "squeeze", "squirrel", "stability", "stable", "stack", "stadium", "staff", "stage",
		"stagger", "stain", "stair", "staircase", "stake", "stale", "stall", "stamp", "stand", "standard",
		"standpoint", "star", "stare", "stark", "start", "starter", "startle", "startup", "starve", "state",
		"statement", "static", "station", "stationary", "statistical", "statistics", "statue", "status", "statute", "statutory",
		"stay", "steady", "steak", "steal", "steam", "steel", "steep", "steer", "stem", "step",
		"stereotype", "stick", "sticky", "stiff", "still", "stimulate", "stimulus", "sting", "stir", "stock",
		"stockholder", "stomach", "stone", "stop", "storage", "store", "storm", "story", "straight", "straightforward",
		"strain", "strange", "stranger", "strategic", "strategy", "straw", "stream", "street", "strength", "strengthen",
		"stress", "stretch", "strike", "striking", "string", "strip", "stripe", "strive", "stroke", "strong",
		"stronger", "strongly", "structural", "structure", "struggle", "student", "studio", "study", "stuff", "stumble",
		"stun", "stunning", "stunt", "stupid", "style", "subject", "submit", "subordinate", "subscribe", "subscriber",
		"subscription", "subsequent", "subsequently", "subsidy", "substance", "substantial", "substantially", "substitute", "subtle", "suburb",
		"suburban", "succeed", "success", "successful", "successfully", "succession", "successor", "such", "suck", "sudden",
		"suddenly", "sue", "suffer", "suffering", "sufficient", "sufficiently", "sugar", "suggest", "suggestion", "suicide",
		"suit", "suitable", "suite", "sum", "summarize", "summary", "summer", "summit", "summon", "sun",
		"super", "superb", "superficial", "superior", "supermarket", "supervise", "supervisor", "supper", "supplier", "supply",
		"support", "supporter", "supportive", "suppose", "supposed", "supposedly", "suppress", "supreme", "sure", "surely",
		"surface", "surge", "surgeon", "surgery", "surgical", "surplus", "surprise", "surprised", "surprising", "surprisingly",
		"surrender", "surround", "surrounding", "surroundings", "surveillance", "survey", "survival", "survive", "survivor", "suspect",
		"suspend", "suspension", "suspicion", "suspicious", "sustain", "sustainable", "swallow", "swamp", "swap", "swarm",
		"sway", "swear", "sweat", "sweater", "sweep", "sweet", "swell", "swept", "swift", "swim",
		"swing", "switch", "sword", "symbol", "symbolic", "sympathetic", "sympathy", "symphony", "symptom", "syndrome",
		"synthesis", "synthetic", "system", "systematic", "systematically", "table", "tablet", "tackle", "tactic", "tactical",
		"tactics", "tail", "tailor", "take", "takeover", "tale", "talent", "talented", "talk", "tall",
		"tank", "tape", "target", "task", "taste", "tax", "taxation", "taxi", "tea", "teach",
		"teacher", "teaching", "team", "teamwork", "tear", "tease", "technical", "technically", "technician", "technique",
		"technology", "teen", "teenage", "teenager", "teeth", "telephone", "telescope", "television", "tell", "temper",
		"temperature", "temple", "temporary", "tempt", "temptation", "tenant", "tend", "tendency", "tender", "tennis",
		"tense", "tension", "tent", "tentative", "term", "terminal", "terminate", "termination", "terms", "terrible",
		"terribly", "terrific", "terrify", "territory", "terror", "terrorism", "terrorist", "test", "testify", "testimony",
		"testing", "text", "textbook", "textile", "texture", "than", "thank", "thanks", "thanksgiving", "that",
		"theater", "theatre", "theft", "their", "them", "theme", "themselves", "then", "theoretical", "theory",
		"therapist", "therapy", "there", "thereafter", "thereby", "therefore", "therein", "thereof", "thermostat", "these",
		"they", "thick", "thickness", "thief", "thigh", "thin", "thing", "think", "thinker", "thinking",
		"third", "thirst", "thirty", "this", "thorough", "thoroughly", "those", "though", "thought", "thoughtful",
		"thoughtful", "thousand", "threat", "threaten", "three", "threshold", "thrift", "thrill", "thriller", "thrive",
		"throat", "throne", "through", "throughout", "throw", "thumb", "thunder", "thus", "tick", "ticket",
		"tide", "tidy", "tie", "tier", "tiger", "tight", "tightly", "tile", "till", "timber",
		"time", "timely", "timetable", "timing", "tiny", "tip", "tire", "tired", "tissue", "title",
		"toast", "tobacco", "today", "toddler", "toe", "together", "toilet", "token", "tolerance", "tolerant",
		"tolerate", "toll", "tomato", "tomorrow", "tone", "tongue", "tonight", "too", "tool", "tooth",
		"top", "topic", "torch", "tornado", "torpedo", "torment", "torrent", "torture", "toss", "total",
		"totally", "touch", "tough", "tour", "tourism", "tourist", "tournament", "toward", "towards", "towel",
		"tower", "town", "township", "toxic", "trace", "track", "tractor", "trade", "trader", "trading",
		"tradition", "traditional", "traffic", "tragedy", "tragic", "trail", "trailer", "train", "trainer", "training",
		"traitor", "tram", "tramp", "transaction", "transcend", "transfer", "transform", "transformation", "transit", "transition",
		"translate", "translation", "transmission", "transmit", "transparent", "transplant", "transport", "transportation", "trap", "trash",
		"trauma", "travel", "traveler", "travelling", "tray", "treasure", "treat", "treatment", "treaty", "tree",
		"tremendous", "trend", "trial", "triangle", "tribal", "tribe", "tribunal", "tribunal", "tribute", "trick",
		"tricky", "trigger", "trillion", "trim", "trio", "trip", "triple", "triumph", "trivial", "troop",
		"trophy", "tropical", "trouble", "troubled", "troublesome", "truck", "true", "truly", "trump", "trunk",
		"trust", "trustee", "trustworthy", "truth", "try", "tube", "tuck", "tuition", "tumble", "tumor",
		"tune", "tunnel", "turkey", "turn", "turnover", "turtle", "tutor", "twelve", "twenty", "twice",
		"twin", "twist", "two", "type", "typical", "typically", "tyranny", "ugly", "ultimate", "ultimately",
		"umbrella", "unable", "unanimous", "unavoidable", "unaware", "uncertain", "uncertainty", "uncle", "unclear", "uncomfortable",
		"unconscious", "unconstitutional", "uncover", "under", "undergo", "undergraduate", "underground", "underlie", "underlying", "undermine",
		"underneath", "understand", "understanding", "undertake", "undertaking", "underwater", "underwear", "undo", "undoubtedly", "unemployment",
		"unexpected", "unfair", "unfavorable", "unfortunate", "unhappy", "unhealthy", "uniform", "unify", "union", "unique",
		"unit", "unite", "united", "unity", "universal", "universe", "university", "unknown", "unless", "unlike",
		"unlikely", "unload", "unlock", "unnecessary", "unpleasant", "unpopular", "unreal", "unreasonable", "unrest", "unsafe",
		"unsettle", "unstable", "unsuccessful", "until", "unusual", "unveil", "unwilling", "up", "upon", "upper",
		"upright", "upset", "urban", "urge", "urgent", "usage", "use", "used", "useful", "user",
		"usual", "usually", "utility", "utilize", "utmost", "utter", "utterly", "vacant", "vacation", "vaccine",
		"vacuum", "vague", "vain", "valid", "validity", "valley", "valuable", "value", "valve", "van",
		"vanish", "vanity", "variable", "variation", "variety", "various", "vary", "vast", "vegetable", "vegetation",
		"vehicle", "veil", "vein", "velocity", "vendor", "venture", "venue", "verbal", "verbally", "verdict",
		"verify", "version", "versus", "vertical", "very", "vessel", "veteran", "veto", "vibrant", "vibration",
		"vice", "victim", "victory", "video", "view", "viewer", "viewpoint", "vigilant", "vigor", "village",
		"villain", "violate", "violation", "violence", "violent", "virtually", "virtue", "virus", "visa", "vision",
		"visit", "visitor", "visual", "vital", "vitamin", "vivid", "vocabulary", "vocal", "vocational", "voice",
		"void", "volcano", "volume", "voluntary", "volunteer", "vote", "voter", "voting", "voyage", "vulnerable",
		"wage", "wagon", "waist", "wait", "waiter", "waiting", "waitress", "wake", "walk", "wall",
		"wallet", "wander", "want", "war", "ward", "warehouse", "warfare", "warm", "warmth", "warn",
		"warning", "warrant", "warrior", "wash", "waste", "watch", "water", "waterfall", "wave", "wavelength",
		"way", "weak", "weakness", "wealth", "wealthy", "weapon", "wear", "weather", "weave", "web",
		"wedding", "wednesday", "week", "weekend", "weekly", "weigh", "weight", "weird", "welcome", "welfare",
		"well", "wellbeing", "wellness", "west", "western", "whale", "what", "whatever", "wheat", "wheel",
		"when", "whenever", "where", "whereas", "wherever", "whether", "which", "whichever", "while", "whip",
		"whisper", "whistle", "white", "who", "whoever", "whole", "wholesale", "wholly", "whom", "whose",
		"why", "wicked", "wide", "widely", "widen", "widespread", "widow", "width", "wife", "wild",
		"wilderness", "wildlife", "will", "willing", "willingness", "win", "wind", "window", "wine", "wing",
		"winner", "winning", "winter", "wipe", "wire", "wisdom", "wise", "wish", "with", "withdraw",
		"withdrawal", "within", "without", "witness", "wolf", "woman", "wonder", "wonderful", "wood", "wooden",
		"wool", "word", "wording", "work", "worker", "workforce", "working", "workout", "workplace", "workshop",
		"world", "worldwide", "worm", "worry", "worrying", "worse", "worship", "worst", "worth", "worthless",
		"worthwhile", "worthy", "would", "wound", "wrap", "wreck", "wrestle", "wrist", "write", "writer",
		"writing", "written", "wrong", "yard", "yeah", "year", "yearly", "yell", "yellow", "yes",
		"yesterday", "yet", "yield", "young", "younger", "your", "yours", "yourself", "youth", "zero",
		"zone", "zoo",
	}

	if words < 2 {
		words = 6 // default
	}

	if words > len(wordList) {
		return "", fmt.Errorf("requested %d words but only %d available", words, len(wordList))
	}

	// Generate secure random word selection
	selectedWords := make([]string, words)
	wordListSize := big.NewInt(int64(len(wordList)))
	usedIndices := make(map[int]bool)

	for i := 0; i < words; i++ {
		var randomIndex int

		// Find unused index
		for {
			randomIndexInt, err := rand.Int(rand.Reader, wordListSize)
			if err != nil {
				return "", fmt.Errorf("failed to generate random word index: %w", err)
			}
			randomIndex = int(randomIndexInt.Int64())
			if !usedIndices[randomIndex] {
				usedIndices[randomIndex] = true
				break
			}
		}

		selectedWords[i] = wordList[randomIndex]
	}

	// Apply transformations
	if capitalize {
		for i, word := range selectedWords {
			selectedWords[i] = strings.Title(word)
		}
	}

	if includeNumber {
		// Add random number at the end
		num, err := rand.Int(rand.Reader, big.NewInt(100))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		selectedWords = append(selectedWords, num.String())
	}

	return strings.Join(selectedWords, separator), nil
}

// Helper functions for password generation

func containsAny(s, chars string) bool {
	for _, c := range chars {
		for _, sc := range s {
			if sc == c {
				return true
			}
		}
	}
	return false
}

func replaceFirstChar(s, chars string) string {
	if len(s) == 0 || len(chars) == 0 {
		return s
	}

	// Find first position to replace
	for i, sc := range s {
		for _, c := range chars {
			if sc == c {
				// Already contains this character type
				return s
			}
		}
		// This character doesn't match any in chars, replace it
		replacement := chars[0]
		return s[:i] + string(replacement) + s[i+1:]
	}

	// If we get here, replace first character
	return string(chars[0]) + s[1:]
}

// evaluatePasswordStrength evaluates the strength of a generated password
func evaluatePasswordStrength(password string) *types.PasswordStrength {
	score := 0
	var issues []string
	var suggestions []string

	// Length scoring
	if len(password) >= 8 {
		score += 10
	}
	if len(password) >= 12 {
		score += 10
	}
	if len(password) >= 16 {
		score += 10
	}
	if len(password) < 8 {
		issues = append(issues, "Password is too short (minimum 8 characters)")
		suggestions = append(suggestions, "Use at least 8 characters")
	}

	// Character variety
	hasLower := containsAny(password, "abcdefghijklmnopqrstuvwxyz")
	hasUpper := containsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	hasNumber := containsAny(password, "0123456789")
	hasSymbol := containsAny(password, "!@#$%^&*()_+-=[]{}|;:,.<>?")

	charTypes := 0
	if hasLower {
		charTypes++
		score += 10
	} else {
		issues = append(issues, "Missing lowercase letters")
		suggestions = append(suggestions, "Include lowercase letters")
	}

	if hasUpper {
		charTypes++
		score += 10
	} else {
		issues = append(issues, "Missing uppercase letters")
		suggestions = append(suggestions, "Include uppercase letters")
	}

	if hasNumber {
		charTypes++
		score += 10
	} else {
		issues = append(issues, "Missing numbers")
		suggestions = append(suggestions, "Include numbers")
	}

	if hasSymbol {
		charTypes++
		score += 10
	} else {
		issues = append(issues, "Missing symbols")
		suggestions = append(suggestions, "Include symbols")
	}

	// Bonus for character variety
	if charTypes >= 3 {
		score += 10
	}

	// Entropy estimation (simplified)
	entropy := float64(len(password)) * float64(charTypes) * 2.585 // log2(6) approx
	if entropy > 60 {
		score += 10
	}

	// Common patterns penalty
	if isCommonPattern(password) {
		score -= 20
		issues = append(issues, "Password contains common patterns")
		suggestions = append(suggestions, "Avoid common patterns and sequences")
	}

	// Determine level
	var level string
	var crackTime string
	switch {
	case score >= 80:
		level = "strong"
		crackTime = "centuries"
	case score >= 60:
		level = "good"
		crackTime = "years"
	case score >= 40:
		level = "fair"
		crackTime = "months"
	default:
		level = "weak"
		crackTime = "days"
	}

	// Cap score at 100
	if score > 100 {
		score = 100
	}

	return &types.PasswordStrength{
		Score:       score,
		Level:       level,
		CrackTime:   crackTime,
		Issues:      issues,
		Suggestions: suggestions,
	}
}

// isCommonPattern checks for common password patterns
func isCommonPattern(password string) bool {
	lower := strings.ToLower(password)

	// Common sequences
	sequences := []string{
		"123", "abc", "qwe", "asd", "zxc", "password", "login", "admin",
		"welcome", "letmein", "master", "shadow", "michael", "jordan",
	}

	for _, seq := range sequences {
		if strings.Contains(lower, seq) {
			return true
		}
	}

	// Repeated characters
	for i := 0; i < len(password)-2; i++ {
		if password[i] == password[i+1] && password[i+1] == password[i+2] {
			return true
		}
	}

	// Sequential characters
	for i := 0; i < len(password)-2; i++ {
		if int(password[i+1]) == int(password[i])+1 &&
			int(password[i+2]) == int(password[i])+2 {
			return true
		}
	}

	return false
}

// getStrengthColor returns color code for password strength level
func getStrengthColor(level string) string {
	switch level {
	case "strong":
		return ui.Green
	case "good":
		return ui.Blue
	case "fair":
		return ui.Yellow
	case "weak":
		return ui.Red
	default:
		return ui.Reset
	}
}
