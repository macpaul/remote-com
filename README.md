# Remote-COM

Remote-COM is a cross-platform application designed to securely redirect local
serial (COM) ports over the network using the SSH protocol. It provides a
robust infrastructure for remote access to embedded systems, industrial
controllers, or any serial-based hardware without exposing raw, unencrypted
TCP ports.

## Key Features

### 1. Core Redirection Engine
- **Serial Port Interface**: Robust, cross-platform enumeration and management
  of physical serial hardware using `go.bug.st/serial`.
- **Embedded SSH Server**: A standalone SSH server powered by
  `golang.org/x/crypto/ssh`, supporting secure, password-authenticated tunnels.
- **Bi-directional Pumping**: A reliable bridge that seamlessly streams data
  between SSH channels and serial ports.

### 2. Management Layer
- **Binding Lifecycle**: A central manager that handles the creation,
  activation, and persistence of port-to-TCP mappings.
- **Configuration Persistence**: Saves all settings to a local `config.json`
  file.
- **Automated Security**: Out-of-the-box encryption with automated SSH host key
  generation.

### 3. Graphical User Interface
- **Modern Dashboard**: A React-based frontend built with Wails for managing
  active redirectors and discovering system ports.
- **Integrated Bridge**: Real-time interaction between the Go engine and the
  web-based UI.

---

# Development README

## About

This is the official Wails React template.

You can configure the project by editing `wails.json`. More information about
the project settings can be found here:
https://wails.io/docs/reference/project-config

## Live Development

To run in live development mode, run `wails dev` in the project directory. This
will run a Vite development server that will provide very fast hot reload of
your frontend changes. If you want to develop in a browser and have access to
your Go methods, there is also a dev server that runs on
http://localhost:34115. Connect to this in your browser, and you can call your
Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.
