# modep
打印 Go 模块的依赖树

## Features


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
This will generate a file named `dependency_tree.html` in the current directory. Open this file in a web browser to view
your module's dependency graph.
