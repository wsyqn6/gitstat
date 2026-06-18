package model

import "time"

type Repository struct {
	Path           string   `json:"path"`
	Name           string   `json:"name"`
	CurrentBranch  string   `json:"currentBranch"`
	LastCommitTime string   `json:"lastCommitTime"`
	UserEmail      string   `json:"userEmail"`
	Commits        []Commit `json:"commits"`
}

type Commit struct {
	Hash      string    `json:"hash"`
	Author    string    `json:"author"`
	Email     string    `json:"-"`
	Date      time.Time `json:"date"`
	Message   string    `json:"message"`
	Additions int       `json:"additions"`
	Deletions int       `json:"deletions"`
}

type OverviewStats struct {
	TotalCommits    int              `json:"totalCommits"`
	TotalAdditions  int              `json:"totalAdditions"`
	TotalDeletions  int              `json:"totalDeletions"`
	ActiveAuthors   int              `json:"activeAuthors"`
	RepositoryCount int              `json:"repositoryCount"`
	Authors         []AuthorRankItem `json:"authors,omitempty"`
}

type AuthorDailyStats struct {
	Author    string          `json:"author"`
	Email     string          `json:"email"`
	Commits   int             `json:"commits"`
	Additions int             `json:"additions"`
	Deletions int             `json:"deletions"`
	IsMe      bool            `json:"isMe"`
	DailyData []DayCommitData `json:"dailyData,omitempty"`
}

type DayCommitData struct {
	Date      string `json:"date"`
	Commits   int    `json:"commits"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
}

type RepositoryDailyStats struct {
	RepoName       string             `json:"repoName"`
	RepoPath       string             `json:"repoPath"`
	CurrentBranch  string             `json:"currentBranch"`
	LastCommitTime string             `json:"lastCommitTime"`
	Authors        []AuthorDailyStats `json:"authors"`
}

// 周期聚合数据结构 (周/月/年)
type PeriodCommitData struct {
	Period    string `json:"date"`
	Commits   int    `json:"commits"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
}

type AuthorPeriodStats struct {
	Author     string            `json:"author"`
	Email      string            `json:"email"`
	Commits    int               `json:"commits"`
	Additions  int               `json:"additions"`
	Deletions  int               `json:"deletions"`
	IsMe       bool              `json:"-"`
	PeriodData []PeriodCommitData `json:"dailyData,omitempty"`
}

type RepositoryPeriodStats struct {
	RepoName       string             `json:"repoName"`
	RepoPath       string             `json:"repoPath"`
	CurrentBranch  string              `json:"-"`
	LastCommitTime string              `json:"-"`
	Authors        []AuthorPeriodStats `json:"authors"`
}

type ApiResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// 开发者排行榜项
type AuthorRankItem struct {
	Author        string  `json:"author"`
	Email         string  `json:"email"`
	Commits       int     `json:"commits"`
	Additions     int     `json:"additions"`
	Deletions     int     `json:"deletions"`
	NetChange     int     `json:"netChange"` // 新增-删除
	IsMe          bool    `json:"isMe"`
	AvgCommitSize float64 `json:"-"` // 平均提交大小
	LastCommitDate string `json:"-"` // 最后提交时间
}

// 活动热力图数据点
type ActivityHeatmapPoint struct {
	DayOfWeek   int `json:"dayOfWeek"` // 0=周日, 1=周一...6=周六
	Hour        int `json:"hour"`      // 0-23
	CommitCount int `json:"commitCount"`
}

// 仓库列表项（快数据）
type RepoInfo struct {
	Path              string   `json:"path"`
	Name              string   `json:"name"`
	CurrentBranch     string   `json:"currentBranch"`
	BranchCount       int      `json:"branchCount"`
	RemoteBranchCount int      `json:"-"`
	FileCount         int      `json:"fileCount"`
	LastCommitTime    string   `json:"lastCommitTime"`
	RemoteUrl         string   `json:"-"`
	Branches          []string `json:"branches"`
	Tags              []string `json:"tags"`
	RemoteBranches    []string `json:"remoteBranches"`
}

// 语言统计
type LanguageStat struct {
	Name       string  `json:"name"`
	FileCount  int     `json:"fileCount"`
	Lines      int     `json:"lines"`
	Percentage float64 `json:"percentage"`
}

// 深度分析结果
type AnalyzeResult struct {
	Name            string        `json:"-"`
	Path            string        `json:"-"`
	BranchCount     int           `json:"-"`
	Branches        []string      `json:"-"`
	RemoteBranches  []string      `json:"-"`
	FileCount       int           `json:"fileCount"`
	TotalLines      int           `json:"totalLines"`
	Languages       []LanguageStat `json:"languages"`
	Tags            []string      `json:"-"`
}

// 深度分析请求
type AnalyzeRequest struct {
	Path string `json:"path"`
}

// 贡献者统计
type ContributorStat struct {
	Author         string `json:"author"`
	Email          string `json:"email"`
	CommitCount    int    `json:"commitCount"`
	Additions      int    `json:"additions"`
	Deletions      int    `json:"deletions"`
	LastCommitDate string `json:"lastCommitDate"`
}

// 仓库统计（计算字段）
type RepoStats struct {
	Path                 string            `json:"path"`
	Name                 string            `json:"name"`
	CurrentBranch        string            `json:"currentBranch"`
	LastCommitTime       string            `json:"lastCommitTime"`
	EarliestDate         string            `json:"earliestDate"`
	EarliestCommitAuthor string            `json:"earliestCommitAuthor"`
	RepoSize             int64             `json:"repoSize"`
	RecentCommits        []Commit          `json:"recentCommits"`
	Contributors         []ContributorStat `json:"contributors"`
	Tags                 []string          `json:"tags"`
	RemoteBranches       []string          `json:"remoteBranches"`
	Analysis             *AnalyzeResult    `json:"analysis,omitempty"`
}

// 仓库对比数据
type RepoComparison struct {
	RepoName         string  `json:"repoName"`
	RepoPath         string  `json:"repoPath"`
	Commits          int     `json:"commits"`
	Authors          int     `json:"authors"`
	Additions        int     `json:"additions"`
	Deletions        int     `json:"deletions"`
	LastCommitTime   string  `json:"-"`
	ActiveDays       int     `json:"activeDays"`       // 活跃天数
	AvgCommitsPerDay float64 `json:"avgCommitsPerDay"` // 日均提交数
}

// 提交日历数据点
type CalendarPoint struct {
	Date  string `json:"date"`  // "2006-01-02"
	Count int    `json:"count"` // 当日提交数
}

// 累计提交数据点
type CumulativePoint struct {
	Date  string `json:"date"`  // "2006-01-02"
	Total int    `json:"total"` // 累计提交数(含当天)
}

// 提交时段数据点
type HourlyPoint struct {
	Hour  int `json:"hour"`  // 0-23
	Count int `json:"count"` // 该时段提交数
}

// Repo 图表数据集
type RepoChartData struct {
	Calendar   []CalendarPoint   `json:"calendar"`
	Cumulative []CumulativePoint `json:"cumulative"`
	Hourly     []HourlyPoint     `json:"hourly"`
}
