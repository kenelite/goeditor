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

Make sure you have Go 1.18+ installed and [Fyne](https://developer.fyne.io/started/) set up.

```bash
git clone https://github.com/kenelite/goeditor.git
cd goedit
go mod tidy
go run main.go
```
Or build the binary:

```bash
go build -o goeditor main.go
./goeditor
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