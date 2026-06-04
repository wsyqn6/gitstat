# GitStat

![Go version](https://img.shields.io/github/go-mod/go-version/wsyqn6/gitstat?filename=backend%2Fgo.mod)
![License](https://img.shields.io/github/license/wsyqn6/gitstat)
![Release](https://img.shields.io/github/v/release/wsyqn6/gitstat)
![CI](https://img.shields.io/github/actions/workflow/status/wsyqn6/gitstat/cd.yml?branch=main)

> Git 仓库提交统计与可视化分析平台 / Git Repository Commit Statistics & Visualization Platform

---

# 中文文档

## 特点

- **单文件分发** — Go 单二进制，前端内嵌，解压即用，零依赖
- **全离线** — 不调用 GitHub API，纯本地扫描，无需联网
- **惰性加载** — 按需拉取 git log，支持增量补数据，启动飞快
- **赛博朋克 UI** — 暗色霓虹风格，ECharts 交互图表
- **开箱即用** — serve 后自动打开浏览器，中英文双语

## 快速开始

从 [Releases](https://github.com/wsyqn6/gitstat/releases) 下载预编译二进制，解压后：

```bash
gitstat serve                  # 扫描当前目录，启动 Web UI
gitstat serve D:/work          # 扫描指定目录
gitstat serve D:/work --port 8080  # 自定义端口
```

启动后自动打开浏览器，访问 `http://localhost:12580` 查看仪表盘。

## 从源码构建

前置要求：Go 1.26+、Node.js 22+、pnpm

```bash
cd frontend && pnpm install && pnpm build
mkdir -p backend/web/dist && cp -r frontend/dist/* backend/web/dist/
cd backend && go build -ldflags="-s -w" -o gitstat.exe .
```

## 技术栈

后端：Go + chi（调 git log 命令解析提交数据）
前端：Vue 3 + Vite + ECharts

## 许可证

[MIT](LICENSE)

---

# English Documentation

## Features

- **Single binary** — Go binary with embedded frontend, extract and run, zero dependencies
- **Fully offline** — No GitHub API calls, scans local repos only, no internet needed
- **Lazy loading** — Fetches git log on demand with incremental updates, fast startup
- **Cyberpunk UI** — Dark neon theme, interactive ECharts visualizations
- **Ready to go** — Auto-opens browser on serve, bilingual (CN/EN)

## Quick Start

Download the pre-built binary from [Releases](https://github.com/wsyqn6/gitstat/releases), extract and run:

```bash
gitstat serve                     # Scan current dir, start Web UI
gitstat serve D:/work             # Scan specified directory
gitstat serve D:/work --port 8080 # Custom port
```

Browser opens automatically. Visit `http://localhost:12580` to view the dashboard.

## Build from Source

Prerequisites: Go 1.26+, Node.js 22+, pnpm

```bash
cd frontend && pnpm install && pnpm build
mkdir -p backend/web/dist && cp -r frontend/dist/* backend/web/dist/
cd backend && go build -ldflags="-s -w" -o gitstat.exe .
```

## Tech Stack

Backend: Go + chi (parses git log output directly)
Frontend: Vue 3 + Vite + ECharts

## License

[MIT](LICENSE)
