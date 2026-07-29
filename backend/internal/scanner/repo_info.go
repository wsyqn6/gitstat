package scanner

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"gitstat/internal/model"
)

var headerPool = sync.Pool{New: func() any { return make([]byte, 8192) }}
var lineBufPool = sync.Pool{New: func() any { return make([]byte, 32*1024) }}

var extLangMap = map[string]string{
	".rs":      "Rust",
	".go":      "Go",
	".py":      "Python",
	".js":      "JavaScript",
	".jsx":     "React (JSX)",
	".mjs":     "JavaScript",
	".ts":      "TypeScript",
	".tsx":     "React (TSX)",
	".vue":     "Vue",
	".svelte":  "Svelte",
	".java":    "Java",
	".kt":      "Kotlin",
	".kts":     "Kotlin",
	".scala":   "Scala",
	".c":       "C",
	".h":       "C/C++ Header",
	".cpp":     "C++",
	".cxx":     "C++",
	".hpp":     "C++ Header",
	".cs":      "C#",
	".rb":      "Ruby",
	".php":     "PHP",
	".swift":   "Swift",
	".m":       "Objective-C",
	".mm":      "Objective-C++",
	".r":       "R",
	".dart":    "Dart",
	".lua":     "Lua",
	".elixir":  "Elixir",
	".ex":      "Elixir",
	".exs":     "Elixir",
	".hs":      "Haskell",
	".zig":     "Zig",
	".pl":      "Perl",
	".pm":      "Perl",
	".sql":     "SQL",
	".sh":      "Shell",
	".bash":    "Shell",
	".zsh":     "Shell",
	".fish":    "Shell",
	".ps1":     "PowerShell",
	".css":     "CSS",
	".scss":    "SCSS",
	".less":    "Less",
	".html":    "HTML",
	".htm":     "HTML",
	".xml":     "XML",
	".yaml":    "YAML",
	".yml":     "YAML",
	".toml":    "TOML",
	".json":    "JSON",
	".jsonc":   "JSON",
	".md":      "Markdown",
	".rst":     "reStructuredText",
	".tex":     "LaTeX",
	".dockerfile": "Dockerfile",
	".makefile":   "Makefile",
	".cmake":      "CMake",
	".gradle":     "Gradle",
	".proto":      "Protobuf",
	".graphql":    "GraphQL",
	".gql":        "GraphQL",
	".mod":     "Go Module",
	".sum":     "Go Sum",
	".lock":    "Lockfile",
	".svg":     "SVG",
	".env":     "Config",
	".npmrc":   "Config",
	".gitignore": "Git Config",
	".gitattributes": "Git Config",
}

var exactNameMap = map[string]string{
	"Dockerfile":      "Dockerfile",
	"Makefile":        "Makefile",
	"CMakeLists.txt":  "CMake",
	"go.mod":          "Go Module",
	"go.sum":          "Go Sum",
	"Cargo.toml":      "Rust",
	"Cargo.lock":      "Rust",
	"Gemfile":         "Ruby",
	"Gemfile.lock":    "Ruby",
	"package.json":    "JavaScript",
	"yarn.lock":       "JavaScript",
	"pnpm-lock.yaml":  "JavaScript",
	"requirements.txt":"Python",
	"Pipfile":         "Python",
	"Pipfile.lock":    "Python",
	"pyproject.toml":  "Python",
	"pom.xml":         "Java",
	"build.gradle":    "Gradle",
	".env.example":    "Config",
	".gitattributes":  "Git Config",
}

var binaryExts = map[string]bool{
	".exe": true, ".dll": true, ".so": true, ".dylib": true,
	".o": true, ".obj": true, ".a": true, ".lib": true,
	".class": true, ".pyc": true, ".pyo": true,
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true,
	".ico": true, ".bmp": true, ".webp": true,
	".ttf": true, ".otf": true, ".woff": true, ".woff2": true,
	".pdf": true,
	".zip": true, ".tar": true, ".gz": true, ".bz2": true, ".xz": true,
	".7z": true, ".rar": true,
	".mp3": true, ".mp4": true, ".avi": true, ".mov": true,
	".wav": true, ".flac": true,
	".jar": true, ".war": true,
	".wasm": true,
}

func GetRepoMeta(repoPath string) (model.RepoInfo, error) {
	currentBranch, _ := gitExec(repoPath, "rev-parse", "--abbrev-ref", "HEAD")

	var branches []string
	branchCount := 0
	out, err := gitExec(repoPath, "for-each-ref", "refs/heads", "--format=%(refname:short)")
	if err == nil && out != "" {
		parts := strings.Split(out, "\n")
		for _, b := range parts {
			b = strings.TrimSpace(b)
			if b == "" {
				continue
			}
			branches = append(branches, b)
			branchCount++
		}
	}

	fileCount := 0
	out, err = gitExec(repoPath, "ls-tree", "-r", "HEAD", "--name-only")
	if err == nil && out != "" {
		fileCount = len(strings.Split(out, "\n"))
	}

	var lastCommitTime string
	out, err = gitExec(repoPath, "log", "-1", "--format=%ci")
	if err == nil {
		if t, e := time.Parse("2006-01-02 15:04:05 -0700", out); e == nil {
			lastCommitTime = t.Format("2006-01-02 15:04:05")
		}
	}

	return model.RepoInfo{
		Path:           repoPath,
		Name:           filepath.Base(repoPath),
		CurrentBranch:  currentBranch,
		BranchCount:    branchCount,
		Branches:       branches,
		FileCount:      fileCount,
		LastCommitTime: lastCommitTime,
	}, nil
}

func AnalyzeRepoDeep(repoPath string) (model.AnalyzeResult, error) {
	currentBranch, _ := gitExec(repoPath, "rev-parse", "--abbrev-ref", "HEAD")

	branchCount := 0
	var branchNames []string
	out, err := gitExec(repoPath, "for-each-ref", "refs/heads", "--format=%(refname:short)")
	if err == nil && out != "" {
		branches := strings.Split(out, "\n")
		branchCount = len(branches)
		var others []string
		for _, b := range branches {
			b = strings.TrimSpace(b)
			if b == "" {
				continue
			}
			if b == currentBranch {
				branchNames = append(branchNames, b)
			} else {
				others = append(others, b)
			}
		}
		sort.Strings(others)
		branchNames = append(branchNames, others...)
	}
	if len(branchNames) == 0 {
		branchNames = []string{currentBranch}
	}

	out, err = gitExec(repoPath, "ls-tree", "-r", "HEAD", "--name-only", "-z")
	if err != nil {
		return model.AnalyzeResult{}, fmt.Errorf("ls-tree: %w", err)
	}
	filePaths := strings.Split(strings.TrimSuffix(out, "\x00"), "\x00")

	langMap := make(map[string]*model.LanguageStat)
	totalLines := 0
	analyzedCount := 0

	for _, rel := range filePaths {
		rel = strings.TrimSpace(rel)
		if rel == "" {
			continue
		}

		ext := strings.ToLower(filepath.Ext(rel))
		base := filepath.Base(rel)
		fullPath := filepath.Join(repoPath, rel)

		if binaryExts[ext] {
			continue
		}

		lang := detectLanguageByExt(base, ext)
		if lang == "" {
			lang = "Other"
		}

		if isNonCode(ext) {
			analyzedCount++
			continue
		}

		f, err := os.Open(fullPath)
		if err != nil {
			continue
		}

		header := headerPool.Get().([]byte)
		defer headerPool.Put(header)
		n, err := f.Read(header)
		if err != nil && err != io.EOF {
			f.Close()
			continue
		}
		header = header[:n]

		if isBinary(header) {
			f.Close()
			continue
		}

		if isGenerated(header) {
			f.Close()
			continue
		}

		if l := detectByShebang(header); l != "" {
			lang = l
		}

		lines := bytes.Count(header, []byte{'\n'})
		more, err := countLinesStream(f)
		f.Close()
		if err != nil {
			continue
		}
		lines += more

		analyzedCount++
		totalLines += lines

		ls, ok := langMap[lang]
		if !ok {
			ls = &model.LanguageStat{Name: lang}
			langMap[lang] = ls
		}
		ls.FileCount++
		ls.Lines += lines
	}

	var languages []model.LanguageStat
	for _, ls := range langMap {
		if totalLines > 0 {
			ls.Percentage = float64(ls.Lines) / float64(totalLines) * 100
		}
		languages = append(languages, *ls)
	}
	sort.Slice(languages, func(i, j int) bool {
		return languages[i].Lines > languages[j].Lines
	})

	remoteBranches, _ := GetRemoteBranches(repoPath)
	tags := GetTags(repoPath)

	return model.AnalyzeResult{
		Name:           filepath.Base(repoPath),
		Path:           repoPath,
		BranchCount:    branchCount,
		Branches:       branchNames,
		RemoteBranches: remoteBranches,
		FileCount:      analyzedCount,
		TotalLines:     totalLines,
		Languages:      languages,
		Tags:           tags,
	}, nil
}

func detectLanguageByExt(base, ext string) string {
	if v, ok := exactNameMap[base]; ok {
		return v
	}
	if v, ok := extLangMap[ext]; ok {
		return v
	}
	if strings.HasPrefix(base, "Dockerfile.") {
		return "Dockerfile"
	}
	if strings.HasPrefix(base, "Makefile.") {
		return "Makefile"
	}
	return ""
}

func countLinesStream(r io.Reader) (int, error) {
	buf := lineBufPool.Get().([]byte)
	defer lineBufPool.Put(buf)
	count := 0
	for {
		n, err := r.Read(buf)
		if n > 0 {
			count += bytes.Count(buf[:n], []byte{'\n'})
		}
		if err != nil {
			if err == io.EOF {
				return count, nil
			}
			return 0, err
		}
	}
}

// ── binary detection ──

func isBinary(b []byte) bool {
	n := len(b)
	if n > 8000 {
		n = 8000
	}
	return bytes.IndexByte(b[:n], 0) >= 0
}

// ── generated file detection ──

var genHeaders = []string{
	"DO NOT EDIT",
	"@generated",
	"auto-generated",
	"autogenerated",
	"Generated by",
	"this file was generated",
	"THIS FILE IS GENERATED",
}

func isGenerated(b []byte) bool {
	// check first 10 lines
	lines := 0
	for _, line := range bytes.Split(b, []byte{'\n'}) {
		if lines >= 10 {
			break
		}
		trimmed := bytes.TrimSpace(line)
		if len(trimmed) == 0 {
			lines++
			continue
		}
		for _, h := range genHeaders {
			if bytes.Contains(trimmed, []byte(h)) {
				return true
			}
		}
		lines++
	}
	return false
}

// ── shebang detection ──

func detectByShebang(b []byte) string {
	if len(b) < 2 || b[0] != '#' || b[1] != '!' {
		return ""
	}
	nl := bytes.IndexByte(b, '\n')
	if nl < 0 {
		nl = len(b)
	}
	line := strings.TrimSpace(string(b[2:nl]))
	if line == "" {
		return ""
	}

	// take last path component, strip args
	parts := strings.Fields(line)
	interp := parts[0]
	if strings.Contains(interp, "/") {
		interp = interp[strings.LastIndex(interp, "/")+1:]
	}
	// /usr/bin/env node → take next arg
	if interp == "env" && len(parts) > 1 {
		interp = parts[1]
	}

	switch interp {
	case "node", "nodejs":
		return "JavaScript"
	case "python", "python2", "python3":
		return "Python"
	case "bash", "sh", "zsh", "fish", "dash":
		return "Shell"
	case "ruby":
		return "Ruby"
	case "perl":
		return "Perl"
	case "php":
		return "PHP"
	case "lua":
		return "Lua"
	case "rustc", "cargo":
		return "Rust"
	case "go", "golang":
		return "Go"
	case "java":
		return "Java"
	case "kotlin":
		return "Kotlin"
	case "swift":
		return "Swift"
	case "deno":
		return "TypeScript"
	case "ts-node":
		return "TypeScript"
	case "racket":
		return "Scheme"
	case "guile":
		return "Scheme"
	case "lisp", "sbcl", "clisp":
		return "Lisp"
	case "elixir":
		return "Elixir"
	case "haskell", "runhaskell":
		return "Haskell"
	case "bashdb":
		return "Shell"
	case "awk":
		return "AWK"
	case "sed":
		return "Shell"
	case "make", "gmake":
		return "Makefile"
	}
	return ""
}

// ── non-code filter (excluded from language stats) ──

func isNonCode(ext string) bool {
	switch ext {
	case ".md", ".rst", ".tex", ".txt":
		return true
	case ".yaml", ".yml", ".svg", ".lock":
		return true
	case ".sum", ".mod":
		return true
	}
	return false
}

func GetTags(repoPath string) []string {
	out, err := gitExec(repoPath, "for-each-ref", "refs/tags", "--sort=-creatordate", "--format=%(refname:short)")
	if err != nil || out == "" {
		return nil
	}
	tags := strings.Split(out, "\n")
	var result []string
	for _, t := range tags {
		t = strings.TrimSpace(t)
		if t != "" {
			result = append(result, t)
		}
	}
	return result
}

func GetTagsCount(repoPath string) int {
	out, err := gitExec(repoPath, "tag", "--list")
	if err != nil || out == "" {
		return 0
	}
	count := 0
	for _, line := range strings.Split(out, "\n") {
		if strings.TrimSpace(line) != "" {
			count++
		}
	}
	return count
}

func GetRemoteBranches(repoPath string) ([]string, int) {
	out, err := gitExec(repoPath, "for-each-ref", "refs/remotes", "--format=%(refname:short)")
	if err != nil || out == "" {
		return nil, 0
	}
	branches := strings.Split(out, "\n")
	var result []string
	for _, b := range branches {
		b = strings.TrimSpace(b)
		if b == "" || strings.HasSuffix(b, "/HEAD") || !strings.Contains(b, "/") {
			continue
		}
		result = append(result, b)
	}
	return result, len(result)
}

func GetRemoteUrl(repoPath string) string {
	out, err := gitExec(repoPath, "remote", "get-url", "origin")
	if err != nil {
		return ""
	}
	return out
}

func GetRepoSize(repoPath string) int64 {
	out, err := gitExec(repoPath, "ls-tree", "-r", "-l", "HEAD")
	if err != nil || out == "" {
		return 0
	}
	var total int64
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// 100644 blob <hash> <size>\t<path>
		tabIdx := strings.IndexByte(line, '\t')
		if tabIdx < 0 {
			continue
		}
		meta := strings.Fields(line[:tabIdx])
		if len(meta) < 4 || meta[1] != "blob" {
			continue
		}
		if n, e := strconv.ParseInt(meta[3], 10, 64); e == nil {
			total += n
		}
	}
	return total
}


