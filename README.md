# GitStat

> Git 仓库提交统计与可视化分析平台 / Git Repository Commit Statistics & Visualization Platform

---

## 目录 / Table of Contents

- [中文文档](#中文文档)
  - [项目简介](#项目简介)
  - [功能特性](#功能特性)
  - [技术栈](#技术栈)
  - [快速开始](#快速开始)
  - [使用方式](#使用方式)
  - [API 概览](#api-概览)
  - [从源码构建](#从源码构建)
  - [许可证](#许可证)
- [English Documentation](#english-documentation)
  - [Introduction](#introduction)
  - [Features](#features)
  - [Tech Stack](#tech-stack)
  - [Quick Start](#quick-start)
  - [Usage](#usage)
  - [API Overview](#api-overview)
  - [Build from Source](#build-from-source)
  - [License](#license)

---

# 中文文档

## 项目简介

**GitStat** 是一个轻量级的 Git 提交统计与可视化分析平台。它能够扫描指定目录下的所有 Git 仓库，从提交历史中提取多维度的活动数据，并通过一个赛博朋克风格的 Web 界面呈现直观的交互式图表。

支持两种运行模式：

- **CLI 模式** — 终端直接输出统计信息
- **Web UI 模式** — 启动内置 HTTP 服务器，通过浏览器访问可视化面板

整个项目打包为一个独立的可执行文件，前端资源内嵌在 Go 二进制中，无需额外部署。

---

## 功能特性

- **多仓库发现与扫描** — 自动扫描指定目录下的所有 Git 仓库，支持增量扫描与懒加载
- **今日看板** — 实时展示今日提交数、增删行数、活跃作者等概览数据
- **多维分析中心** — 支持时间范围筛选（本周/月/年/自定义）、仓库筛选
  - 每日提交趋势图（多用户折线图）
  - 代码变更分布（堆叠柱状图）
  - 开发者排行榜（横向柱状图）
  - 活动热力图（按星期×小时）
  - 仓库对比雷达图（多维度对比）
- **作者洞察** — 自动识别"本周之星"，展示平均提交规模、净代码增长等指标
- **数据导出** — 一键导出全部数据为 JSON 文件
- **国际化** — 内置中文 / 英文界面切换
- **独立二进制** — 前端嵌入 Go 二进制，单文件分发，零依赖

---

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端语言 | Go 1.26 |
| Web 框架 | go-chi/chi v5 |
| Git 操作库 | go-git/go-git v5 |
| 前端框架 | Vue 3 (Composition API, `<script setup>`) |
| 构建工具 | Vite 8 |
| 图表库 | ECharts 6 |
| 包管理器 | pnpm |
| 静态资源嵌入 | Go `embed` 包 |

---

## 快速开始

### 下载预编译二进制

从 [Releases](https://github.com/anomalyco/gitstat/releases) 页面下载对应平台的最新版本，解压后即可使用。

### 启动 Web UI

```bash
gitstat serve [扫描目录] --port 12580
```

例如：

```bash
gitstat serve D:/work --port 12580
```

启动后浏览器将自动打开，访问 `http://localhost:12580` 即可进入仪表盘。

### CLI 模式查看统计

```bash
gitstat stats [扫描目录]
```

---

## 使用方式

### 命令说明

| 命令 | 说明 |
|------|------|
| `gitstat serve [dir]` | 启动 Web 服务器，`dir` 为扫描目录，默认 `D:/work` |
| `gitstat stats [dir]` | 终端输出统计信息 |
| `gitstat --version` / `-v` | 显示版本号 |
| `gitstat --help` / `-h` | 显示帮助信息 |

### Web UI 页面

- **仪表盘 (Dashboard)** — 今日概览卡片 + 各仓库各作者逐日明细
- **分析中心 (Analytics)** — 时间范围/仓库筛选下的多图表分析
- **设置 (Settings)** — 扫描路径配置、数据导出、版本信息

### 时间范围

支持 `1d`、`7d`、`30d`、`90d`、`all` 等快捷时间范围，以及在 Web UI 中自定义起止日期。

---

## API 概览

| 方法 | 端点 | 说明 |
|------|------|------|
| `GET` | `/api/version` | 服务端版本 |
| `POST` | `/api/scan` | 扫描目录（指定时间范围） |
| `POST` | `/api/scan/path` | 设置扫描路径 |
| `GET` | `/api/scan/path` | 获取当前扫描路径 |
| `GET` | `/api/repositories` | 列出已发现的仓库 |
| `GET` | `/api/stats/overview` | 概览统计数据 |
| `GET` | `/api/stats/daily` | 逐日统计（按仓库/作者） |
| `GET` | `/api/stats/authors` | 作者排行榜 |
| `GET` | `/api/stats/activity-heatmap` | 活动热力图数据 |
| `GET` | `/api/stats/repo-comparison` | 仓库多维度对比 |
| `POST` | `/api/export/json` | 导出全部数据为 JSON |
| `GET` | `/health` | 健康检查 |

---

## 从源码构建

### 前置要求

- Go 1.26+
- Node.js 22+
- pnpm

### 构建步骤

```bash
# 1. 构建前端
cd frontend
pnpm install
pnpm run build

# 2. 复制前端产物到后端嵌入目录
cp -r frontend/dist/* backend/web/dist/

# 3. 构建 Go 二进制
cd backend
go build -ldflags="-s -w" -o gitstat.exe .
```

或使用项目自带的构建脚本：

```powershell
# PowerShell
.\build.ps1
```

构建产物为单文件 `gitstat.exe`（Windows）或 `gitstat`（Linux/macOS），可直接分发运行。

---

## 许可证

[MIT](LICENSE)

---

# English Documentation

## Introduction

**GitStat** is a lightweight Git commit statistics and visualization platform. It scans all Git repositories under a specified directory, extracts multi-dimensional activity data from commit history, and presents intuitive interactive charts through a cyberpunk-themed web interface.

Two operation modes are available:

- **CLI mode** — Output statistics directly in the terminal
- **Web UI mode** — Start a built-in HTTP server and access the dashboard through a browser

The entire project is packaged as a single standalone executable, with frontend assets embedded in the Go binary — no additional deployment required.

---

## Features

- **Multi-Repository Discovery & Scanning** — Automatically discovers Git repositories in a directory tree; supports incremental scanning and lazy loading
- **Today's Dashboard** — Real-time overview cards showing today's commits, additions/deletions, active authors, and per-repo/per-author daily breakdown
- **Multi-Dimensional Analytics** — Filter by time range (this week/month/year, custom date range) and repositories
  - Daily commit trend (multi-user line chart)
  - Code change distribution (stacked bar chart)
  - Developer leaderboard (horizontal bar chart)
  - Activity heatmap (by day-of-week × hour)
  - Repository comparison radar chart
- **Author Insights** — Automatically identifies "Star of the Week", displays average commit size, net code growth, etc.
- **Data Export** — One-click export of all data as JSON
- **Internationalization** — Built-in Chinese / English UI toggle
- **Single Binary** — Frontend embedded in the Go binary, zero dependencies

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend Language | Go 1.26 |
| Web Framework | go-chi/chi v5 |
| Git Library | go-git/go-git v5 |
| Frontend Framework | Vue 3 (Composition API, `<script setup>`) |
| Build Tool | Vite 8 |
| Charting Library | ECharts 6 |
| Package Manager | pnpm |
| Static Assets | Go `embed` package |

---

## Quick Start

### Download Pre-built Binary

Download the latest release for your platform from the [Releases](https://github.com/anomalyco/gitstat/releases) page, extract and run.

### Start Web UI

```bash
gitstat serve [scan-directory] --port 12580
```

Example:

```bash
gitstat serve D:/work --port 12580
```

The browser will open automatically at `http://localhost:12580`.

### CLI Mode

```bash
gitstat stats [scan-directory]
```

---

## Usage

### Commands

| Command | Description |
|---------|-------------|
| `gitstat serve [dir]` | Start web server; `dir` defaults to `D:/work` |
| `gitstat stats [dir]` | Output statistics to terminal |
| `gitstat --version` / `-v` | Show version |
| `gitstat --help` / `-h` | Show help |

### Web UI Pages

- **Dashboard** — Today's overview cards + per-repo/per-author daily breakdown
- **Analytics** — Multi-chart analysis with time range and repository filters
- **Settings** — Scan path configuration, data export, version info

### Time Ranges

Quick ranges: `1d`, `7d`, `30d`, `90d`, `all`. Custom date ranges are supported in the Web UI.

---

## API Overview

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/version` | Server version |
| `POST` | `/api/scan` | Scan directory with time range |
| `POST` | `/api/scan/path` | Set scan path |
| `GET` | `/api/scan/path` | Get current scan path |
| `GET` | `/api/repositories` | List discovered repositories |
| `GET` | `/api/stats/overview` | Overview statistics |
| `GET` | `/api/stats/daily` | Daily statistics (by repo/author) |
| `GET` | `/api/stats/authors` | Author leaderboard |
| `GET` | `/api/stats/activity-heatmap` | Activity heatmap data |
| `GET` | `/api/stats/repo-comparison` | Multi-repo comparison metrics |
| `POST` | `/api/export/json` | Export all data as JSON |
| `GET` | `/health` | Health check |

---

## Build from Source

### Prerequisites

- Go 1.26+
- Node.js 22+
- pnpm

### Build Steps

```bash
# 1. Build frontend
cd frontend
pnpm install
pnpm run build

# 2. Copy frontend output to backend embed directory
cp -r frontend/dist/* backend/web/dist/

# 3. Build Go binary
cd backend
go build -ldflags="-s -w" -o gitstat.exe .
```

Or use the provided build script:

```powershell
# PowerShell
.\build.ps1
```

The output is a single file `gitstat.exe` (Windows) or `gitstat` (Linux/macOS), ready to distribute and run.

---

## License

[MIT](LICENSE)
