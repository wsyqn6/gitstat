# GitStat

![Go version](https://img.shields.io/github/go-mod/go-version/wsyqn6/gitstat?filename=backend%2Fgo.mod)
![License](https://img.shields.io/github/license/wsyqn6/gitstat)
![Release](https://img.shields.io/github/v/release/wsyqn6/gitstat)
![CI](https://img.shields.io/github/actions/workflow/status/wsyqn6/gitstat/cd.yml?branch=main)

> Git 仓库提交统计与可视化分析平台 / Git Repository Commit Statistics & Visualization Platform

---

# 中文文档

## 简介

**GitStat** 扫描指定目录下所有 Git 仓库，提取提交活动数据，通过赛博朋克风格 Web 界面呈现交互式图表。支持 CLI 终端输出和 Web UI 仪表盘两种模式。单文件分发，零依赖。

## 快速开始

从 [Releases](https://github.com/wsyqn6/gitstat/releases) 下载预编译二进制，解压后即可使用。

```bash
gitstat serve D:/work --port 12580   # 启动 Web UI
gitstat stats D:/work                 # CLI 统计
```

## 从源码构建

前置要求：Go 1.26+、Node.js 22+、pnpm

```bash
cd frontend && pnpm install && pnpm build
mkdir -p backend/web/dist && cp -r frontend/dist/* backend/web/dist/
cd backend && go build -ldflags="-s -w" -o gitstat.exe .
```

或使用 `.\build.ps1`

## 技术栈

后端：Go + chi + go-git | 前端：Vue 3 + Vite + ECharts

## 许可证

[MIT](LICENSE)

---

# English Documentation

## Introduction

**GitStat** scans Git repositories under a directory, extracts commit activity data, and displays interactive charts through a cyberpunk-themed web UI. Supports CLI terminal output and Web UI dashboard modes. Single binary, zero dependencies.

## Quick Start

Download the pre-built binary from [Releases](https://github.com/wsyqn6/gitstat/releases), extract and run.

```bash
gitstat serve D:/work --port 12580   # Start Web UI
gitstat stats D:/work                 # CLI stats
```

## Build from Source

Prerequisites: Go 1.26+, Node.js 22+, pnpm

```bash
cd frontend && pnpm install && pnpm build
mkdir -p backend/web/dist && cp -r frontend/dist/* backend/web/dist/
cd backend && go build -ldflags="-s -w" -o gitstat.exe .
```

Or use `.\build.ps1`

## Tech Stack

Backend: Go + chi + go-git | Frontend: Vue 3 + Vite + ECharts

## License

[MIT](LICENSE)
