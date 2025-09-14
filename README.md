# gedit 🖋️

A simple, modern, and lightweight terminal-based text editor built in Go.

[![Go Version](https://img.shields.io/badge/Go-1.22+-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/v/release/arit-pal/gedit)](https://github.com/arit-pal/gedit/releases)

`gedit` was born as a learning project and has evolved into a fully functional text editor for the terminal, designed with a clean, modular architecture. It's fast, easy to use, and perfect for quick edits without leaving your command line.
***

### Features

* **Full Text Editing:** Type, delete, and insert text with ease.
* **File Management:** Open, edit, and save files.
* **Smooth Scrolling:** Handles files larger than the screen.
* **Search:** Find any text with `Ctrl+F`.
* **Find Next:** Cycle through search results with `Tab`.
* **Dynamic Status Bar:** Real-time feedback on filename, modified status, line count, and cursor position.
* **Smart Indentation:** Use `Tab` for 4-space indents and `Backspace` to intelligently remove them.
* **Keyboard-Driven:** Designed for efficient, mouse-free operation.

***

### Installation & Building from Source

To use `gedit`, you need to have Go installed (version 1.20+ is recommended).

1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/arit-pal/gedit.git](https://github.com/arit-pal/gedit.git)
    cd gedit
    ```

2.  **Build the binary:**
    ```bash
    go build -o gedit ./cmd/gedit/
    ```

3.  **(Optional) Move the binary to your PATH:**
    For easy access from anywhere in your terminal, move the compiled `gedit` binary to a directory in your system's PATH.
    ```bash
    # For Linux/macOS
    sudo mv gedit /usr/local/bin/
    ```

***

### Usage

Run `gedit` from your terminal, passing a filename as an argument.

```bash
# Edit an existing file
gedit my_file.txt

# Create a new file
gedit new_document.md
````

-----

### Key Bindings

| Key | Action |
| :--- | :--- |
| `Ctrl` + `S` | Save the current file. |
| `Ctrl` + `X` | Quit the editor. |
| `Ctrl` + `F` | Open the search prompt. |
| `Arrow Keys` | Move the cursor. |
| `Enter` | Create a new line. |
| `Backspace` | Delete character to the left / join lines / smart delete tab. |
| `Delete` | Delete character to the right / join lines. |
| `Tab` | Insert a 4-space indent **or** find the next search result. |
| `Escape` | Cancel the search prompt. |

-----

### Code Architecture

`gedit` is built with a clean, modular architecture to make it easy for developers to understand and contribute. The logic is separated into distinct packages:

  * **`app`**: The main application orchestrator, containing the event loop.
  * **`editor`**: Defines the core `State` struct, which holds all the data about the editor's current state.
  * **`view`**: Handles all rendering and drawing logic to the terminal screen.
  * **`input`**: Manages all keyboard input and executes commands that modify the editor state.
  * **`file`**: Contains the logic for loading from and saving to the file system.

-----

### Contributing

Contributions are welcome\! If you find a bug or have a feature idea, please open an issue or submit a pull request.

-----

### License

This project is licensed under the MIT License. See the `LICENSE` file for details.