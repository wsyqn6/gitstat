package scanner

import (
	"bufio"
	"io"
	"path/filepath"
	"sort"
	"strings"

	"gitstat/internal/model"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var extLangMap = map[string]string{
	".rs":      "Rust",
	".go":      "Go",
	".py":      "Python",
	".js":      "JavaScript",
	".jsx":     "React (JSX)",
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
}

var exactNameMap = map[string]string{
	"Dockerfile": "Dockerfile",
	"Makefile":   "Makefile",
	"CMakeLists.txt": "CMake",
}

func GetRepoMeta(repoPath string) (model.RepoInfo, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return model.RepoInfo{}, err
	}

	head, err := r.Head()
	if err != nil {
		return model.RepoInfo{}, err
	}
	currentBranch := head.Name().Short()

	branchCount := 0
	bIter, err := r.Branches()
	if err == nil {
		bIter.ForEach(func(ref *plumbing.Reference) error {
			branchCount++
			return nil
		})
	}

	fileCount := 0

	var lastCommitTime string
	headCommit, err := r.CommitObject(head.Hash())
	if err == nil {
		lastCommitTime = headCommit.Committer.When.Format("2006-01-02 15:04:05")

		tree, err := headCommit.Tree()
		if err == nil {
			tree.Files().ForEach(func(f *object.File) error {
				fileCount++
				return nil
			})
		}
	}

	return model.RepoInfo{
		Path:           repoPath,
		Name:           filepath.Base(repoPath),
		CurrentBranch:  currentBranch,
		BranchCount:    branchCount,
		FileCount:      fileCount,
		LastCommitTime: lastCommitTime,
	}, nil
}

func AnalyzeRepoDeep(repoPath string) (model.AnalyzeResult, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return model.AnalyzeResult{}, err
	}

	head, err := r.Head()
	if err != nil {
		return model.AnalyzeResult{}, err
	}
	currentBranch := head.Name().Short()

	branchCount := 0
	var branchNames []string
	bIter, err := r.Branches()
	if err == nil {
		bIter.ForEach(func(ref *plumbing.Reference) error {
			branchCount++
			name := ref.Name().Short()
			if name != currentBranch {
				branchNames = append(branchNames, name)
			}
			return nil
		})
	}
	sort.Strings(branchNames)
	branchNames = append([]string{currentBranch + " (current)"}, branchNames...)

	headCommit, err := r.CommitObject(head.Hash())
	if err != nil {
		return model.AnalyzeResult{}, err
	}

	tree, err := headCommit.Tree()
	if err != nil {
		return model.AnalyzeResult{}, err
	}

	type fileInfo struct {
		path  string
		lines int
		lang  string
	}
	var files []fileInfo

	tree.Files().ForEach(func(f *object.File) error {
		ext := strings.ToLower(filepath.Ext(f.Name))
		base := filepath.Base(f.Name)

		var lang string
		if v, ok := exactNameMap[base]; ok {
			lang = v
		} else if v, ok := extLangMap[ext]; ok {
			lang = v
		} else {
			lang = "Other"
		}

		blob, err := r.BlobObject(f.Hash)
		if err != nil {
			files = append(files, fileInfo{f.Name, 0, lang})
			return nil
		}
		reader, err := blob.Reader()
		if err != nil {
			files = append(files, fileInfo{f.Name, 0, lang})
			return nil
		}
		lines := countLines(reader)
		reader.Close()

		files = append(files, fileInfo{f.Name, lines, lang})
		return nil
	})

	langMap := make(map[string]*model.LanguageStat)
	totalLines := 0
	for _, f := range files {
		totalLines += f.lines
		if _, ok := langMap[f.lang]; !ok {
			langMap[f.lang] = &model.LanguageStat{Name: f.lang}
		}
		langMap[f.lang].FileCount++
		langMap[f.lang].Lines += f.lines
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

	return model.AnalyzeResult{
		Name:        filepath.Base(repoPath),
		Path:        repoPath,
		BranchCount: branchCount,
		Branches:    branchNames,
		FileCount:   len(files),
		TotalLines:  totalLines,
		Languages:   languages,
	}, nil
}

func GetRemoteUrl(repoPath string) string {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return ""
	}
	remote, err := r.Remote("origin")
	if err != nil {
		return ""
	}
	if len(remote.Config().URLs) > 0 {
		return remote.Config().URLs[0]
	}
	return ""
}

func GetRepoSize(repoPath string) int64 {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return 0
	}
	head, err := r.Head()
	if err != nil {
		return 0
	}
	commit, err := r.CommitObject(head.Hash())
	if err != nil {
		return 0
	}
	tree, err := commit.Tree()
	if err != nil {
		return 0
	}
	var total int64
	tree.Files().ForEach(func(f *object.File) error {
		blob, err := r.BlobObject(f.Hash)
		if err != nil {
			return nil
		}
		total += blob.Size
		return nil
	})
	return total
}

func countLines(reader io.Reader) int {
	scanner := bufio.NewScanner(reader)
	lines := 0
	for scanner.Scan() {
		lines++
	}
	return lines
}
