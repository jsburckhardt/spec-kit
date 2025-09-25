# GoSpecify - Complete Go Implementation Specification

## Project Overview

**GoSpecify** is a complete Go reimplementation of the Python Specify CLI that bootstraps projects for Spec-Driven Development (SDD). This implementation will be a single, self-contained binary that embeds all templates and scripts internally, eliminating external dependencies while maintaining full compatibility with the existing ecosystem.

## Architecture Design

### Core Principles

1. **Single Binary**: All functionality embedded in one executable
2. **Zero External Dependencies**: Templates, scripts, and assets bundled internally
3. **Cross-Platform Native**: Leverage Go's compilation strengths
4. **API Compatibility**: Maintain identical CLI interface and behavior
5. **Performance**: Faster startup and execution than Python version

### Technology Stack

```go
// Core Dependencies
github.com/spf13/cobra          // CLI framework (industry standard)
github.com/spf13/viper          // Configuration management
github.com/charmbracelet/bubbletea  // Terminal UI framework
github.com/charmbracelet/bubbles    // UI components
github.com/charmbracelet/lipgloss   // Styling and layout
github.com/schollz/progressbar/v3   // Progress visualization
github.com/AlecAivazis/survey/v2    // Interactive prompts (alternative)

// Standard Library
embed                           // Embed static assets
net/http                        // GitHub API client
archive/zip                     // Template extraction
os/exec                         // Git operations
runtime                         // OS detection
path/filepath                   // Cross-platform paths
```

## Project Structure

```
cmd/
├── gospecify/
│   └── main.go                 // Entry point
├── root.go                     // Root command setup
├── init.go                     // Project initialization
├── check.go                    // Tool validation
└── version.go                  // Version information

internal/
├── config/
│   ├── agents.go               // AI assistant definitions
│   ├── constants.go            // Global constants
│   └── types.go                // Core data structures
├── templates/
│   ├── embedded.go             // Embedded template assets
│   ├── processor.go            // Template processing logic
│   └── extractor.go            // Template extraction
├── scripts/
│   ├── embedded.go             // Embedded script assets  
│   ├── executor.go             // Cross-platform script execution
│   └── generator.go            // Dynamic script generation
├── ui/
│   ├── banner.go               // ASCII art and branding
│   ├── progress.go             // Progress tracking trees
│   ├── selector.go             // Interactive selection
│   └── styles.go               // Consistent styling
├── github/
│   ├── client.go               // GitHub API integration
│   ├── auth.go                 // Token authentication
│   └── releases.go             // Release asset handling
├── git/
│   ├── operations.go           // Git repository operations
│   └── detection.go            // Repository detection
└── tools/
    ├── checker.go              // Tool availability validation
    └── platforms.go            // Cross-platform tool paths

pkg/
├── specify/
│   ├── project.go              // Project configuration
│   ├── assistant.go            // AI assistant logic
│   └── workflow.go             // SDD workflow management
└── errors/
    └── types.go                // Custom error types

assets/                         // Source files for embedding
├── templates/
│   ├── agent-file-template.md
│   ├── plan-template.md
│   ├── spec-template.md
│   ├── tasks-template.md
│   └── commands/
│       ├── analyze.md
│       ├── clarify.md
│       ├── constitution.md
│       ├── implement.md
│       ├── plan.md
│       ├── specify.md
│       └── tasks.md
└── scripts/
    ├── bash/
    │   ├── check-prerequisites.sh
    │   ├── create-new-feature.sh
    │   ├── setup-plan.sh
    │   └── update-agent-context.sh
    └── powershell/
        ├── check-prerequisites.ps1
        ├── create-new-feature.ps1  
        ├── setup-plan.ps1
        └── update-agent-context.ps1

scripts/
├── build.sh                   // Build automation
├── embed.sh                   // Asset embedding automation
└── release.sh                 // Release packaging

docs/
├── README.md                  // Go-specific documentation
├── ARCHITECTURE.md            // Technical architecture
└── MIGRATION.md              // Migration from Python version
```

## Core Data Structures

### AI Assistant Configuration

```go
type AIAssistant struct {
    Key         string            `json:"key"`
    Name        string            `json:"name"`
    Directory   string            `json:"directory"`
    Format      FileFormat        `json:"format"`
    CLITool     string            `json:"cli_tool,omitempty"`
    ArgFormat   string            `json:"arg_format"`
    IsIDEBased  bool             `json:"is_ide_based"`
    Website     string            `json:"website"`
}

type FileFormat string

const (
    FormatMarkdown FileFormat = "md"
    FormatTOML     FileFormat = "toml"
    FormatPrompt   FileFormat = "prompt.md"
)

// Pre-configured AI assistants (matches Python AI_CHOICES)
var AIAssistants = map[string]AIAssistant{
    "copilot": {
        Key:        "copilot",
        Name:       "GitHub Copilot",
        Directory:  ".github/prompts/",
        Format:     FormatPrompt,
        ArgFormat:  "$ARGUMENTS",
        IsIDEBased: true,
        Website:    "https://github.com/features/copilot",
    },
    "claude": {
        Key:       "claude",
        Name:      "Claude Code",
        Directory: ".claude/commands/",
        Format:    FormatMarkdown,
        CLITool:   "claude",
        ArgFormat: "$ARGUMENTS",
        Website:   "https://docs.anthropic.com/en/docs/claude-code/setup",
    },
    "gemini": {
        Key:       "gemini",
        Name:      "Gemini CLI", 
        Directory: ".gemini/commands/",
        Format:    FormatTOML,
        CLITool:   "gemini",
        ArgFormat: "{{args}}",
        Website:   "https://github.com/google-gemini/gemini-cli",
    },
    // ... all other assistants from Python version
}
```

### Project Configuration

```go
type ProjectConfig struct {
    Name             string    `json:"name"`
    Path             string    `json:"path"`
    AIAssistant      string    `json:"ai_assistant"`
    ScriptType       string    `json:"script_type"`
    NoGit           bool      `json:"no_git"`
    Force           bool      `json:"force"`
    IgnoreTools     bool      `json:"ignore_tools"`
    SkipTLS         bool      `json:"skip_tls"`
    Debug           bool      `json:"debug"`
    GitHubToken     string    `json:"github_token,omitempty"`
    Here            bool      `json:"here"`
    CreatedAt       time.Time `json:"created_at"`
}

type ScriptType struct {
    Key         string `json:"key"`
    Name        string `json:"name"`
    Extension   string `json:"extension"`
    Platform    string `json:"platform"`
}

var ScriptTypes = map[string]ScriptType{
    "sh": {
        Key:       "sh",
        Name:      "POSIX Shell (bash/zsh)",
        Extension: ".sh",
        Platform:  "unix",
    },
    "ps": {
        Key:       "ps", 
        Name:      "PowerShell",
        Extension: ".ps1",
        Platform:  "windows",
    },
}
```

### Progress Tracking

```go
type StepTracker struct {
    Title       string          `json:"title"`
    Steps       []Step          `json:"steps"`
    StatusOrder map[string]int  `json:"-"`
    refreshCb   func()          `json:"-"`
    mu          sync.RWMutex    `json:"-"`
}

type Step struct {
    Key     string `json:"key"`
    Label   string `json:"label"`
    Status  Status `json:"status"`
    Detail  string `json:"detail"`
    Started time.Time `json:"started,omitempty"`
    Ended   time.Time `json:"ended,omitempty"`
}

type Status string

const (
    StatusPending Status = "pending"
    StatusRunning Status = "running"
    StatusDone    Status = "done"
    StatusError   Status = "error"
    StatusSkipped Status = "skipped"
)
```

## Embedded Assets Strategy

### Template Embedding

```go
//go:embed assets/templates/*
var templateFS embed.FS

//go:embed assets/scripts/*
var scriptFS embed.FS

type EmbeddedAssets struct {
    Templates map[string][]byte
    Scripts   map[string][]byte
}

func LoadEmbeddedAssets() (*EmbeddedAssets, error) {
    assets := &EmbeddedAssets{
        Templates: make(map[string][]byte),
        Scripts:   make(map[string][]byte),
    }
    
    // Load all template files
    err := fs.WalkDir(templateFS, "assets/templates", func(path string, d fs.DirEntry, err error) error {
        if err != nil || d.IsDir() {
            return err
        }
        
        content, err := templateFS.ReadFile(path)
        if err != nil {
            return err
        }
        
        relativePath := strings.TrimPrefix(path, "assets/templates/")
        assets.Templates[relativePath] = content
        return nil
    })
    
    // Similar for scripts...
    return assets, err
}
```

### Dynamic Script Generation

```go
type ScriptGenerator struct {
    assets     *EmbeddedAssets
    projectCfg *ProjectConfig
}

func (sg *ScriptGenerator) GenerateScript(scriptName string, agent AIAssistant) ([]byte, error) {
    var templatePath string
    switch sg.projectCfg.ScriptType {
    case "sh":
        templatePath = fmt.Sprintf("bash/%s.sh", scriptName)
    case "ps":
        templatePath = fmt.Sprintf("powershell/%s.ps1", scriptName)
    default:
        return nil, fmt.Errorf("unsupported script type: %s", sg.projectCfg.ScriptType)
    }
    
    template, exists := sg.assets.Scripts[templatePath]
    if !exists {
        return nil, fmt.Errorf("script template not found: %s", templatePath)
    }
    
    // Perform placeholder replacement
    script := string(template)
    script = strings.ReplaceAll(script, "__AGENT__", agent.Key)
    script = strings.ReplaceAll(script, "{SCRIPT}", sg.getScriptPath(scriptName))
    
    return []byte(script), nil
}
```

## CLI Interface Implementation

### Root Command

```go
func NewRootCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "gospecify",
        Short: "Setup tool for Specify spec-driven development projects",
        Long:  getBanner() + "\n\nGitHub Spec Kit - Spec-Driven Development Toolkit",
        PersistentPreRun: func(cmd *cobra.Command, args []string) {
            // Global setup
        },
        Run: func(cmd *cobra.Command, args []string) {
            showBanner()
            fmt.Println("Run 'gospecify --help' for usage information")
        },
    }
    
    // Add subcommands
    cmd.AddCommand(NewInitCmd())
    cmd.AddCommand(NewCheckCmd())
    cmd.AddCommand(NewVersionCmd())
    
    return cmd
}

func getBanner() string {
    return `
███████╗██████╗ ███████╗ ██████╗██╗███████╗██╗   ██╗
██╔════╝██╔══██╗██╔════╝██╔════╝██║██╔════╝╚██╗ ██╔╝
███████╗██████╔╝█████╗  ██║     ██║█████╗   ╚████╔╝ 
╚════██║██╔═══╝ ██╔══╝  ██║     ██║██╔══╝    ╚██╔╝  
███████║██║     ███████╗╚██████╗██║██║        ██║   
╚══════╝╚═╝     ╚══════╝ ╚═════╝╚═╝╚═╝        ╚═╝   `
}
```

### Init Command

```go
func NewInitCmd() *cobra.Command {
    var cfg ProjectConfig
    
    cmd := &cobra.Command{
        Use:   "init [project-name]",
        Short: "Initialize a new Specify project from the latest template",
        Long: `Initialize a new Specify project from the latest template.

This command will:
1. Check that required tools are installed (git is optional)
2. Let you choose your AI assistant (Claude Code, Gemini CLI, GitHub Copilot, etc.)
3. Download the appropriate template from GitHub
4. Extract the template to a new project directory or current directory  
5. Initialize a fresh git repository (if not --no-git and no existing repo)
6. Optionally set up AI assistant commands

Examples:
  gospecify init my-project
  gospecify init my-project --ai claude
  gospecify init --here --ai claude
  gospecify init --here --force`,
        Args: func(cmd *cobra.Command, args []string) error {
            if cfg.Here && len(args) > 0 {
                return fmt.Errorf("cannot specify both project name and --here flag")
            }
            if !cfg.Here && len(args) == 0 {
                return fmt.Errorf("must specify either a project name or use --here flag")
            }
            return nil
        },
        RunE: func(cmd *cobra.Command, args []string) error {
            if len(args) > 0 {
                cfg.Name = args[0]
            }
            return runInit(&cfg)
        },
    }
    
    // Flags
    cmd.Flags().StringVar(&cfg.AIAssistant, "ai", "", 
        "AI assistant to use: claude, gemini, copilot, cursor, qwen, opencode, codex, windsurf, kilocode, or auggie")
    cmd.Flags().StringVar(&cfg.ScriptType, "script", "", 
        "Script type to use: sh or ps")
    cmd.Flags().BoolVar(&cfg.IgnoreTools, "ignore-agent-tools", false, 
        "Skip checks for AI agent tools like Claude Code")
    cmd.Flags().BoolVar(&cfg.NoGit, "no-git", false, 
        "Skip git repository initialization")
    cmd.Flags().BoolVar(&cfg.Here, "here", false, 
        "Initialize project in the current directory instead of creating a new one")
    cmd.Flags().BoolVar(&cfg.Force, "force", false, 
        "Force merge/overwrite when using --here (skip confirmation)")
    cmd.Flags().BoolVar(&cfg.SkipTLS, "skip-tls", false, 
        "Skip SSL/TLS verification (not recommended)")
    cmd.Flags().BoolVar(&cfg.Debug, "debug", false, 
        "Show verbose diagnostic output for network and extraction failures")
    cmd.Flags().StringVar(&cfg.GitHubToken, "github-token", "", 
        "GitHub token to use for API requests (or set GH_TOKEN or GITHUB_TOKEN environment variable)")
    
    return cmd
}
```

## Interactive UI Implementation

### Selection Interface

```go
type Selector struct {
    prompt   string
    options  map[string]string
    default  string
    selected string
    index    int
}

func NewSelector(prompt string, options map[string]string, defaultKey string) *Selector {
    return &Selector{
        prompt:  prompt,
        options: options,
        default: defaultKey,
        index:   0,
    }
}

func (s *Selector) Run() (string, error) {
    keys := make([]string, 0, len(s.options))
    for k := range s.options {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    
    // Set default index
    if s.default != "" {
        for i, key := range keys {
            if key == s.default {
                s.index = i
                break
            }
        }
    }
    
    // Use bubbletea for interactive selection
    program := tea.NewProgram(s)
    result, err := program.Run()
    if err != nil {
        return "", err
    }
    
    finalModel := result.(Selector)
    return finalModel.selected, nil
}

func (s Selector) Init() tea.Cmd {
    return nil
}

func (s Selector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        keys := make([]string, 0, len(s.options))
        for k := range s.options {
            keys = append(keys, k)
        }
        sort.Strings(keys)
        
        switch msg.Type {
        case tea.KeyUp:
            s.index = (s.index - 1 + len(keys)) % len(keys)
        case tea.KeyDown:
            s.index = (s.index + 1) % len(keys)
        case tea.KeyEnter:
            s.selected = keys[s.index]
            return s, tea.Quit
        case tea.KeyEsc, tea.KeyCtrlC:
            return s, tea.Quit
        }
    }
    return s, nil
}

func (s Selector) View() string {
    keys := make([]string, 0, len(s.options))
    for k := range s.options {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    
    var output strings.Builder
    
    output.WriteString(lipgloss.NewStyle().Bold(true).Render(s.prompt) + "\n\n")
    
    for i, key := range keys {
        prefix := "  "
        style := lipgloss.NewStyle()
        
        if i == s.index {
            prefix = "▶ "
            style = style.Foreground(lipgloss.Color("cyan"))
        }
        
        line := fmt.Sprintf("%s%s %s", prefix, 
            style.Render(key), 
            lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(fmt.Sprintf("(%s)", s.options[key])))
        output.WriteString(line + "\n")
    }
    
    output.WriteString("\nUse ↑/↓ to navigate, Enter to select, Esc to cancel")
    
    return lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("cyan")).
        Padding(1, 2).
        Render(output.String())
}
```

### Progress Tracking UI

```go
func (st *StepTracker) Render() string {
    st.mu.RLock()
    defer st.mu.RUnlock()
    
    tree := treeprint.New()
    tree.SetValue(lipgloss.NewStyle().Foreground(lipgloss.Color("cyan")).Render(st.Title))
    
    for _, step := range st.Steps {
        var symbol, style string
        
        switch step.Status {
        case StatusDone:
            symbol = "●"
            style = "green"
        case StatusPending:
            symbol = "○"  
            style = "240"
        case StatusRunning:
            symbol = "○"
            style = "cyan"
        case StatusError:
            symbol = "●"
            style = "red"
        case StatusSkipped:
            symbol = "○"
            style = "yellow"
        }
        
        label := step.Label
        if step.Detail != "" {
            if step.Status == StatusPending {
                label = fmt.Sprintf("%s (%s)", label, step.Detail)
                style = "240"
            } else {
                label = fmt.Sprintf("%s (%s)", 
                    lipgloss.NewStyle().Render(label),
                    lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(step.Detail))
            }
        }
        
        symbolStyled := lipgloss.NewStyle().Foreground(lipgloss.Color(style)).Render(symbol)
        line := fmt.Sprintf("%s %s", symbolStyled, label)
        
        if step.Status == StatusPending {
            line = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(line)
        }
        
        tree.AddNode(line)
    }
    
    return tree.String()
}
```

## GitHub Integration

### Release Asset Download

```go
type GitHubClient struct {
    httpClient *http.Client
    token      string
    baseURL    string
}

func NewGitHubClient(token string, skipTLS bool) *GitHubClient {
    client := &http.Client{
        Timeout: 30 * time.Second,
    }
    
    if skipTLS {
        tr := &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }
        client.Transport = tr
    }
    
    return &GitHubClient{
        httpClient: client,
        token:      token,
        baseURL:    "https://api.github.com",
    }
}

func (gc *GitHubClient) GetLatestRelease(owner, repo string) (*Release, error) {
    url := fmt.Sprintf("%s/repos/%s/%s/releases/latest", gc.baseURL, owner, repo)
    
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    if gc.token != "" {
        req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", gc.token))
    }
    req.Header.Set("Accept", "application/vnd.github.v3+json")
    
    resp, err := gc.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("GitHub API returned %d", resp.StatusCode)
    }
    
    var release Release
    if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
        return nil, err
    }
    
    return &release, nil
}

func (gc *GitHubClient) DownloadAsset(asset ReleaseAsset, destPath string, progress *progressbar.ProgressBar) error {
    req, err := http.NewRequest("GET", asset.BrowserDownloadURL, nil)
    if err != nil {
        return err
    }
    
    if gc.token != "" {
        req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", gc.token))
    }
    
    resp, err := gc.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("download failed with status %d", resp.StatusCode)
    }
    
    file, err := os.Create(destPath)
    if err != nil {
        return err
    }
    defer file.Close()
    
    if progress != nil {
        progress.ChangeMax64(resp.ContentLength)
        _, err = io.Copy(io.MultiWriter(file, progress), resp.Body)
    } else {
        _, err = io.Copy(file, resp.Body)
    }
    
    return err
}

type Release struct {
    TagName string         `json:"tag_name"`
    Assets  []ReleaseAsset `json:"assets"`
}

type ReleaseAsset struct {
    Name               string `json:"name"`
    Size               int64  `json:"size"`
    BrowserDownloadURL string `json:"browser_download_url"`
}
```

## Cross-Platform Script Execution

### Script Execution Engine

```go
type ScriptExecutor struct {
    projectPath string
    scriptType  string
    assets      *EmbeddedAssets
}

func NewScriptExecutor(projectPath, scriptType string, assets *EmbeddedAssets) *ScriptExecutor {
    return &ScriptExecutor{
        projectPath: projectPath,
        scriptType:  scriptType,
        assets:      assets,
    }
}

func (se *ScriptExecutor) Execute(scriptName string, args ...string) error {
    script, err := se.getScript(scriptName)
    if err != nil {
        return err
    }
    
    // Create temporary script file
    tmpFile, err := se.createTempScript(script)
    if err != nil {
        return err
    }
    defer os.Remove(tmpFile)
    
    // Execute script
    return se.executeScript(tmpFile, args...)
}

func (se *ScriptExecutor) getScript(scriptName string) ([]byte, error) {
    var scriptPath string
    switch se.scriptType {
    case "sh":
        scriptPath = fmt.Sprintf("bash/%s.sh", scriptName)
    case "ps":
        scriptPath = fmt.Sprintf("powershell/%s.ps1", scriptName)
    default:
        return nil, fmt.Errorf("unsupported script type: %s", se.scriptType)
    }
    
    script, exists := se.assets.Scripts[scriptPath]
    if !exists {
        return nil, fmt.Errorf("script not found: %s", scriptPath)
    }
    
    return script, nil
}

func (se *ScriptExecutor) createTempScript(content []byte) (string, error) {
    var pattern string
    switch se.scriptType {
    case "sh":
        pattern = "specify-script-*.sh"
    case "ps":
        pattern = "specify-script-*.ps1"
    }
    
    tmpFile, err := os.CreateTemp("", pattern)
    if err != nil {
        return "", err
    }
    defer tmpFile.Close()
    
    if _, err := tmpFile.Write(content); err != nil {
        os.Remove(tmpFile.Name())
        return "", err
    }
    
    // Make executable on Unix systems
    if se.scriptType == "sh" && runtime.GOOS != "windows" {
        if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
            os.Remove(tmpFile.Name())
            return "", err
        }
    }
    
    return tmpFile.Name(), nil
}

func (se *ScriptExecutor) executeScript(scriptPath string, args ...string) error {
    var cmd *exec.Cmd
    
    switch se.scriptType {
    case "sh":
        if runtime.GOOS == "windows" {
            // Use Git Bash or WSL if available
            if bash := findBashOnWindows(); bash != "" {
                cmd = exec.Command(bash, scriptPath)
            } else {
                return fmt.Errorf("bash not found on Windows")
            }
        } else {
            cmd = exec.Command("/bin/bash", scriptPath)
        }
    case "ps":
        if runtime.GOOS == "windows" {
            cmd = exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", scriptPath)
        } else {
            cmd = exec.Command("pwsh", "-File", scriptPath)
        }
    }
    
    cmd.Args = append(cmd.Args, args...)
    cmd.Dir = se.projectPath
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    
    return cmd.Run()
}

func findBashOnWindows() string {
    // Common locations for bash on Windows
    locations := []string{
        "C:\\Program Files\\Git\\bin\\bash.exe",
        "C:\\Windows\\System32\\bash.exe", // WSL
        "bash", // In PATH
    }
    
    for _, location := range locations {
        if _, err := exec.LookPath(location); err == nil {
            return location
        }
    }
    
    return ""
}
```

## Template Processing System

### Template Engine

```go
type TemplateProcessor struct {
    assets      *EmbeddedAssets
    config      *ProjectConfig
    assistant   *AIAssistant
    replacements map[string]string
}

func NewTemplateProcessor(assets *EmbeddedAssets, config *ProjectConfig, assistant *AIAssistant) *TemplateProcessor {
    return &TemplateProcessor{
        assets:    assets,
        config:    config,
        assistant: assistant,
        replacements: map[string]string{
            "__AGENT__":    assistant.Key,
            "{SCRIPT}":     "", // Will be set per template
            "$ARGUMENTS":   assistant.ArgFormat,
            "{{args}}":     assistant.ArgFormat,
        },
    }
}

func (tp *TemplateProcessor) ProcessTemplate(templateName string) ([]byte, error) {
    template, exists := tp.assets.Templates[templateName]
    if !exists {
        return nil, fmt.Errorf("template not found: %s", templateName)
    }
    
    content := string(template)
    
    // Apply all replacements
    for placeholder, replacement := range tp.replacements {
        content = strings.ReplaceAll(content, placeholder, replacement)
    }
    
    // Handle format-specific processing
    switch tp.assistant.Format {
    case FormatMarkdown:
        return tp.processMarkdownTemplate(content)
    case FormatTOML:
        return tp.processTOMLTemplate(content)
    case FormatPrompt:
        return tp.processPromptTemplate(content)
    }
    
    return []byte(content), nil
}

func (tp *TemplateProcessor) processMarkdownTemplate(content string) ([]byte, error) {
    // Extract frontmatter if present
    parts := strings.SplitN(content, "---", 3)
    if len(parts) == 3 {
        frontmatter := parts[1]
        body := parts[2]
        
        // Process frontmatter
        var fm map[string]interface{}
        if err := yaml.Unmarshal([]byte(frontmatter), &fm); err != nil {
            return nil, fmt.Errorf("invalid frontmatter: %w", err)
        }
        
        // Apply script replacement in frontmatter
        if scripts, ok := fm["scripts"].(map[string]interface{}); ok {
            if scriptCmd, ok := scripts[tp.config.ScriptType].(string); ok {
                scriptPath := tp.resolveScriptPath(scriptCmd)
                body = strings.ReplaceAll(body, "{SCRIPT}", scriptPath)
            }
        }
        
        // Reconstruct
        processedFM, _ := yaml.Marshal(fm)
        return []byte(fmt.Sprintf("---\n%s---\n%s", processedFM, body)), nil
    }
    
    return []byte(content), nil
}

func (tp *TemplateProcessor) processTOMLTemplate(content string) ([]byte, error) {
    // TOML format processing for Gemini/Qwen
    lines := strings.Split(content, "\n")
    var processed []string
    
    inPrompt := false
    for _, line := range lines {
        if strings.HasPrefix(line, "prompt = \"\"\"") {
            inPrompt = true
        } else if inPrompt && strings.HasPrefix(line, "\"\"\"") {
            inPrompt = false
        }
        
        if inPrompt {
            // Apply replacements to prompt content
            for placeholder, replacement := range tp.replacements {
                line = strings.ReplaceAll(line, placeholder, replacement)
            }
        }
        
        processed = append(processed, line)
    }
    
    return []byte(strings.Join(processed, "\n")), nil
}

func (tp *TemplateProcessor) resolveScriptPath(scriptCmd string) string {
    // Convert relative script path to absolute based on project structure
    return fmt.Sprintf(".specify/scripts/%s", scriptCmd)
}
```

## Build and Deployment System

### Build Script

```bash
#!/bin/bash
# scripts/build.sh

set -euo pipefail

VERSION=${VERSION:-"dev"}
COMMIT=$(git rev-parse --short HEAD)
DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS="-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}"

echo "Building GoSpecify ${VERSION} (${COMMIT})"

# Build for multiple platforms
PLATFORMS=(
    "linux/amd64"
    "linux/arm64" 
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

mkdir -p dist

for platform in "${PLATFORMS[@]}"; do
    GOOS=${platform%/*}
    GOARCH=${platform#*/}
    
    output="dist/gospecify-${GOOS}-${GOARCH}"
    if [ "$GOOS" = "windows" ]; then
        output="${output}.exe"
    fi
    
    echo "Building for ${GOOS}/${GOARCH}..."
    
    GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags "$LDFLAGS" \
        -o "$output" \
        ./cmd/gospecify
done

echo "Build complete. Binaries in dist/"
```

### Asset Embedding Script

```bash
#!/bin/bash
# scripts/embed.sh

set -euo pipefail

echo "Embedding assets..."

# Ensure assets directory exists
mkdir -p assets/templates assets/scripts

# Copy templates from source
cp -r templates/* assets/templates/

# Copy scripts from source  
cp -r scripts/bash assets/scripts/
cp -r scripts/powershell assets/scripts/

# Generate embedded assets file
cat > internal/templates/embedded.go << 'EOF'
// Code generated by embed.sh; DO NOT EDIT.

package templates

import "embed"

//go:embed ../../assets/templates/*
var TemplateFS embed.FS

//go:embed ../../assets/scripts/*  
var ScriptFS embed.FS
EOF

echo "Assets embedded successfully"
```

## Migration Strategy

### Compatibility Matrix

| Feature | Python Version | Go Version | Status |
|---------|---------------|------------|---------|
| CLI Interface | ✅ | ✅ | 1:1 Compatible |
| AI Assistant Support | ✅ | ✅ | All 11 assistants |
| Template System | ✅ | ✅ | Embedded + GitHub |
| Progress Tracking | ✅ | ✅ | Enhanced UI |
| Cross-platform Scripts | ✅ | ✅ | Embedded execution |
| GitHub Integration | ✅ | ✅ | Native HTTP client |
| Interactive Selection | ✅ | ✅ | Bubbletea UI |
| Git Operations | ✅ | ✅ | Native exec |
| SSL/TLS Support | ✅ | ✅ | Standard library |
| Configuration | ✅ | ✅ | Viper + JSON |

### Migration Path

1. **Phase 1**: Core CLI structure and basic commands
2. **Phase 2**: Template system and GitHub integration  
3. **Phase 3**: Interactive UI and progress tracking
4. **Phase 4**: Script execution engine
5. **Phase 5**: Full feature parity testing
6. **Phase 6**: Performance optimization and release

### Performance Targets

| Metric | Python Version | Go Version Target |
|--------|---------------|-------------------|
| Binary Size | N/A (interpreter) | <50MB (with assets) |
| Cold Start | ~500ms | <100ms |
| Template Download | ~2-3s | ~1-2s |
| Template Extraction | ~1-2s | <500ms |
| Memory Usage | ~50-100MB | <20MB |

## Testing Strategy

### Test Structure

```
test/
├── unit/                   # Unit tests for all packages
├── integration/            # Integration tests  
├── e2e/                   # End-to-end CLI tests
├── compatibility/         # Python vs Go comparison tests
└── performance/           # Benchmark tests

benchmarks/
├── startup_test.go        # Startup time benchmarks
├── download_test.go       # Download performance
└── memory_test.go         # Memory usage tests
```

### Test Categories

1. **Unit Tests**: Individual function/method testing
2. **Integration Tests**: Component interaction testing
3. **CLI Tests**: Command-line interface testing
4. **Compatibility Tests**: Output comparison with Python version
5. **Performance Tests**: Speed and memory benchmarks

## Deployment and Distribution

### Release Process

1. **Automated Builds**: GitHub Actions for all platforms
2. **Asset Preparation**: Embedded templates and scripts
3. **Version Management**: Semantic versioning with Git tags
4. **Distribution Channels**: 
   - GitHub Releases (primary)
   - Homebrew (macOS)
   - Chocolatey (Windows)
   - Package managers (Linux)

### Installation Methods

```bash
# Direct download
curl -L https://github.com/github/gospecify/releases/latest/download/gospecify-linux-amd64 -o gospecify
chmod +x gospecify

# Homebrew (macOS)
brew install github/tap/gospecify

# Chocolatey (Windows)
choco install gospecify

# Go install
go install github.com/github/gospecify/cmd/gospecify@latest
```

## Success Metrics

### Compatibility Goals

- [ ] 100% CLI argument compatibility
- [ ] 100% output format compatibility  
- [ ] 100% template compatibility
- [ ] 100% workflow compatibility

### Performance Goals

- [ ] 5x faster startup time
- [ ] 3x faster template processing
- [ ] 50% smaller resource usage
- [ ] Single binary deployment

### User Experience Goals

- [ ] Identical user workflows
- [ ] Enhanced visual feedback
- [ ] Better error messages
- [ ] Faster operations

This specification provides a complete roadmap for implementing GoSpecify as a high-performance, single-binary replacement for the Python Specify CLI while maintaining full compatibility and enhancing the user experience.