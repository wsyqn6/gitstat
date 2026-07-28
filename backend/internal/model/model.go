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

type FileStat struct {
	Path      string `json:"path"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
}

type FileRankItem struct {
	FilePath  string `json:"filePath"`
	Commits   int    `json:"commits"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
	NetChange int    `json:"netChange"`
}

type Commit struct {
	Hash      string    `json:"hash"`
	Author    string    `json:"author"`
	Email     string    `json:"-"`
	Date      time.Time `json:"date"`
	Message   string    `json:"message"`
	Additions int       `json:"additions"`
	Deletions int       `json:"deletions"`
	Files     []FileStat `json:"-"`
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
	Author       string          `json:"author"`
	Email        string          `json:"email"`
	Commits      int             `json:"commits"`
	Additions    int             `json:"additions"`
	Deletions    int             `json:"deletions"`
	FilesChanged int             `json:"filesChanged"`
	IsMe         bool            `json:"isMe"`
	DailyData    []DayCommitData `json:"dailyData,omitempty"`
}

type DayCommitData struct {
	Date      string `json:"date"`
	Commits   int    `json:"commits"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
}

type DayCommitCount struct {
	Date    string `json:"date"`
	Commits int    `json:"commits"`
}

type DailyTrendItem struct {
	RepoName     string           `json:"repoName"`
	DailyCommits []DayCommitCount `json:"dailyCommits,omitempty"`
}

type RepositoryDailyStats struct {
	RepoName       string             `json:"repoName"`
	RepoPath       string             `json:"repoPath"`
	CurrentBranch  string             `json:"currentBranch"`
	LastCommitTime string             `json:"lastCommitTime"`
	Authors        []AuthorDailyStats `json:"authors"`
	DailyCommits   []DayCommitCount   `json:"dailyCommits,omitempty"`
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
	IsMe       bool              `json:"isMe"`
	PeriodData []PeriodCommitData `json:"periodData,omitempty"`
}

type RepositoryPeriodStats struct {
	RepoName       string             `json:"repoName"`
	RepoPath       string             `json:"repoPath"`
	CurrentBranch  string             `json:"currentBranch"`
	LastCommitTime string             `json:"lastCommitTime"`
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
	IsMe bool `json:"isMe"`
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
	RemoteBranchCount int      `json:"remoteBranchCount"`
	FileCount         int      `json:"fileCount"`
	LastCommitTime    string   `json:"lastCommitTime"`
	RemoteUrl         string   `json:"remoteUrl"`
	Branches          []string `json:"branches"`
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
	Name            string        `json:"name"`
	Path            string        `json:"path"`
	BranchCount     int           `json:"branchCount"`
	Branches        []string      `json:"branches"`
	RemoteBranches  []string      `json:"remoteBranches"`
	FileCount       int           `json:"fileCount"`
	TotalLines      int           `json:"totalLines"`
	Languages       []LanguageStat `json:"languages"`
	Tags            []string      `json:"tags"`
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

// 提交分页响应
type CommitPage struct {
	Commits []Commit `json:"commits"`
	HasMore bool     `json:"hasMore"`
}

// Tag 分页响应
type TagPage struct {
	Tags    []string `json:"tags"`
	Total   int      `json:"total"`
	HasMore bool     `json:"hasMore"`
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
	Contributors         []ContributorStat `json:"contributors"`
	Analysis             *AnalyzeResult    `json:"analysis,omitempty"`
}

// 月历聚合项（前端直接使用，无需计算）
type MonthlyCalendarItem struct {
	Month     string `json:"month"`     // "2006-01"
	Commits   int    `json:"commits"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
}

// Dashboard 聚合响应（6 卡片 + 今日详情）
type DashboardSummary struct {
	TotalCommits    int `json:"totalCommits"`
	TotalAdditions  int `json:"totalAdditions"`
	TotalDeletions  int `json:"totalDeletions"`
	ActiveAuthors   int `json:"activeAuthors"`
	RepositoryCount int `json:"repositoryCount"`
}

type DashboardData struct {
	Summary    DashboardSummary       `json:"summary"`
	DailyRepos []RepositoryDailyStats `json:"dailyRepos"`
	Comparison OverviewComparison     `json:"comparison"`
}

// 同比/环比单个指标变化
type MetricChange struct {
	Current  int     `json:"current"`
	Previous int     `json:"previous"`
	Abs      int     `json:"abs"`
	Pct      float64 `json:"pct"`
}

// 同比/环比单个作者各指标变化
type AuthorMetricChange struct {
	Commits   MetricChange `json:"commits"`
	Additions MetricChange `json:"additions"`
	Deletions MetricChange `json:"deletions"`
	NetChange int          `json:"netChange"`
}

// 同比/环比作者对比项
type AuthorComparisonItem struct {
	Author string             `json:"author"`
	Email  string             `json:"email"`
	IsMe   bool               `json:"isMe"`
	Change AuthorMetricChange `json:"change"`
}

// 同比/环比概览对比
type OverviewComparison struct {
	TotalCommits   MetricChange           `json:"totalCommits"`
	TotalAdditions MetricChange           `json:"totalAdditions"`
	TotalDeletions MetricChange           `json:"totalDeletions"`
	ActiveAuthors  MetricChange           `json:"activeAuthors"`
	Authors        []AuthorComparisonItem `json:"authors,omitempty"`
}

// 仓库对比数据
type RepoComparison struct {
	RepoName         string  `json:"repoName"`
	RepoPath         string  `json:"repoPath"`
	Commits          int     `json:"commits"`
	Authors          int     `json:"authors"`
	Additions        int     `json:"additions"`
	Deletions        int     `json:"deletions"`
	FilesChanged     int     `json:"filesChanged"`
	LastCommitTime   string  `json:"-"`
	ActiveDays       int     `json:"activeDays"`       // 活跃天数
	AvgCommitsPerDay float64 `json:"avgCommitsPerDay"` // 日均提交数
}

// 提交时段数据点
type HourlyPoint struct {
	Hour  int `json:"hour"`  // 0-23
	Count int `json:"count"` // 该时段提交数
}

// 每日聚合数据点（合并日历/累计/代码行）
type DailyAggPoint struct {
	Date     string `json:"date"`     // "2006-01-02"
	Commits  int    `json:"commits"`  // 当日提交数
	Total    int    `json:"total"`    // 累计提交数
	NetLines int    `json:"netLines"` // 累计净增行数
}

// Repo 图表数据集
type RepoChartData struct {
	DailyAgg []DailyAggPoint `json:"dailyAgg"`
	Hourly   []HourlyPoint   `json:"hourly"`
}
