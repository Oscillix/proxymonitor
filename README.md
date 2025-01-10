# ProxyMonitor

This project was created as a test project for abz.agency and is my first project written in [Golang](https://golang.org/).

---

## How It Works

The application uses the `golang.org/x/sys/windows/registry` package to monitor the Windows Registry for changes to proxy settings.

- **Registry Key Monitored:**  
  `HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Internet Settings`

- **Logging Location:**  
  Logs are stored at: `%APPDATA%\proxymonitor\proxymonitor.log`

---

## Dependencies

This project uses the following dependencies:

1. **[`golang.org/x/sys/windows/registry`](https://pkg.go.dev/golang.org/x/sys/windows/registry)**  
   Provides access to Windows Registry APIs.

2. **[`github.com/allan-simon/go-singleinstance`](https://github.com/allan-simon/go-singleinstance)**  
   A cross-platform library to ensure only one instance of the software runs.

3. **[`github.com/getlantern/systray`](https://github.com/getlantern/systray)**  
   A cross-platform Go library for placing an icon and menu in the notification area.

---

## Build Instructions

Before getting started, ensure you have [Golang](https://golang.org/) installed on your system. To build and run the project, follow these steps:

1. Clone the repository:
   ```bash
   git clone https://github.com/Oscillix/proxymonitor.git
   ```

2. Navigate to the project directory:
   ```bash
   cd proxymonitor
   ```

3. Download the project dependencies:
   ```bash
   go mod download
   ```

4. Build the project:
   ```bash
   go build
   ```

5. Run the executable:
   ```bash
   proxymonitor.exe
   ```

---

## CLI Commands

- **Start the service:**
  ```bash
  proxymonitor.exe -start
  ```

- **Stop the service:**
  ```bash
  proxymonitor.exe -stop
  ```

- **Quit the service entirely:**
  ```bash
  proxymonitor.exe -quit
  ```

---

## Troubleshooting

The program is designed to handle most issues and will report errors to the console. 

If you see `Error opening registry key`, it's likely because the registry key doesn't exist. To fix this:  
1. Go to **Settings > Network & Internet > Proxy**.  
2. Under **Manual proxy setup**, toggle **Use a proxy server** on and save the setting.  
3. Then, toggle it off again (if needed).  