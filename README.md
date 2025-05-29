# goeditor

A lightweight single-document text editor written in Go using the [Fyne](https://fyne.io/) GUI toolkit.

## Features

- Simple and clean interface
- Open, edit, save, and save as files
- Keyboard shortcuts for common actions:
    - New (Ctrl+N)
    - Open (Ctrl+O)
    - Save (Ctrl+S)
    - Save As (Ctrl+Shift+S)
    - Quit (Ctrl+Q)

## Installation

### 1. Clone the repository
Make sure you have Go 1.18+ installed and [Fyne](https://developer.fyne.io/started/) set up.

```bash
git clone https://github.com/kenelite/goeditor.git
cd goeditor
```

### 2. Install Fyne dependencies

Fyne requires some native dependencies depending on your OS.

Linux (Debian/Ubuntu example):

```bash
sudo apt install libgl1-mesa-dev x11proto-core-dev libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev
```

macOS

Install Xcode Command Line Tools:

```bash
xcode-select --install
```

Windows

No special native dependencies needed, but you need to have a working Go environment.

### 3. Download Go modules
```bash
go mod tidy
```


## Build & Run

Run directly
```bash
go run main.go
```

Build binary
```bash
go build -o goedit main.go
```

##  Usage
Use the File menu or keyboard shortcuts to create, open, save, or quit.

The editor supports a single file at a time.

## Folder Structure
ui/ — user interface code, windows, menus, editor widget

backend/ — file handling, state management, and helper utilities

## Contributions
Contributions and suggestions are welcome! Feel free to open issues or submit pull requests.

## License
MIT License © 2025 Kenelite