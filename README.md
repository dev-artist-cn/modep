# modep - Go项目模块依赖关系可视化

[![Go Report Card](https://goreportcard.com/badge/github.com/dev-artist-cn/modep)](https://goreportcard.com/report/github.com/dev-artist-cn/modep)
[![GoDoc](https://godoc.org/github.com/dev-artist-cn/modep?status.svg)](https://pkg.go.dev/github.com/dev-artist-cn/modep)

把 Go 项目中模块之间依赖的层级关系，用树形结构打印出来。

## Features

[<img src="docs/screen-shot.png?raw=true">](https://www.bilibili.com/video/BV1SN6JYkETW)

## Installation

### Option 1
To install modep, use the following command:

```bash
go install github.com/dev-artist-cn/modep
```
Ensure that your Go bin directory is in your system's PATH.

### Option 2

```bash
git clone https://github.com/dev-artist-cn/modep.git
cd modep
go install
```

## Usage

Navigate to your Go project's root directory and run:
```
modep
```
This will print the module dependencies in tree hierachy.

## Video Course

[Go项目模块依赖关系可视化](https://www.bilibili.com/video/BV1Cu6BYQEdx)