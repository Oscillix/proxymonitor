// This project was made for abz.agency as a test project.
package main

import (
	"fmt"           // We'll need this to actually log to the console and file.
	"os"            // We'll need this to get the environment variables, create directories and open files.
	"path/filepath" // We'll need this to join the path of the log file (line 34).
	"strconv"       // We'll need this to convert the string PID from the lock file to an integer. (line 153 only).
	"syscall"       // We'll need this to kill the process by PID. (killProcess() function).
	"time"          // Time is needed to get the current time. (Obviously....).

	"golang.org/x/sys/windows/registry" // We'll need this to access the registry.

	"github.com/allan-simon/go-singleinstance"   // We'll need this to check if the program is already running via .lock file.
	"github.com/getlantern/systray"              // Third party library, has multi os support (Windows, Linux & MacOS), easy to use but for this case we're using only windows specific calls.
	"github.com/getlantern/systray/example/icon" // We'll use the example icon from systray, the test didn't say WHICH icon to use. :)
)

// We need this to actually stop logging since you can't stop a goroutine "directly".
var stopLoggingChan bool = false

func startLogging() {
	// Open the registry key where proxy settings are stored, personally I think this is the easiest way to get the proxy settings.
	// This will most likely fail IF the registry key has never been created.
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.QUERY_VALUE)
	if err != nil {
		fmt.Println("Error opening registry key:", err)
		return
	}
	defer key.Close()

	// Create a log file in %APPDATA%\proxymonitor folder named proxymonitor.log
	// We'll need to create the folder if it doesn't exist.
	logDir := filepath.Join(os.Getenv("APPDATA"), "proxymonitor")
	dirErr := os.MkdirAll(logDir, os.ModePerm)
	if dirErr != nil {
		fmt.Println("Error creating log directory: ", dirErr)
		return
	}

	// Open the log file in append mode, create it if it doesn't exist
	logFilePath := filepath.Join(logDir, "proxymonitor.log")
	logFile, fileErr := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // 0644 is the permission for the file
	if fileErr != nil {
		fmt.Println("Error opening log file: ", fileErr)
		return
	}
	defer logFile.Close() // Close the log file

	// Start logging to the log file
	fmt.Println("Logging to: ", logFile.Name())

	var lastProxyEnable uint64
	var lastProxyServer string

	for {
		if stopLoggingChan {
			fmt.Println("Proxy stopped.")
			return // This will kill the for loop.
		} else {
			now := time.Now() // Current local time

			// Read the proxy settings, enabled
			proxyEnable, _, enabledErr := key.GetIntegerValue("ProxyEnable")
			if enabledErr != nil {
				fmt.Println("Error reading ProxyEnable value:", enabledErr)
				return
			}

			// Read the proxy settings, server / IP address
			proxyServer, _, serverErr := key.GetStringValue("ProxyServer")
			if serverErr != nil {
				fmt.Println("Error reading ProxyServer value:", serverErr)
				return
			}

			// Check if the proxy settings have changed, to avoid spamming the log file with the same information.
			if proxyEnable != lastProxyEnable || proxyServer != lastProxyServer {
				var status string

				// Since golang doesn't have match statements, we'll have to do a if else statement and since proxyEnable is an integer, we'll have to do "!= 0"...
				if proxyEnable != 0 {
					status = "on"
				} else {
					status = "off"
				}

				// Log the proxy status to console. DEBUG PURPOSES ONLY!
				// fmt.Printf("%s	proxy %s, %s\n", now.Format(time.ANSIC), status, proxyServer)

				// Log the proxy status to file.
				fmt.Fprintf(logFile, "%s	proxy %s, %s\n", now.Format(time.ANSIC), status, proxyServer)

				// Update the last known proxy settings
				lastProxyEnable = proxyEnable
				lastProxyServer = proxyServer
			}

			time.Sleep(1 * time.Second) // Sleep for a second before checking again
		}
	}
}

func startProxy() {
	stopLoggingChan = false
	go startLogging() // This starts a goroutine to start checking and logging.
}

func stopProxy() {
	if !stopLoggingChan {
		stopLoggingChan = true
	}
}

// This killProcess code is not made by me, I've found it on stackoverflow.
func killProcess(pid int) error {
	// Find the process by PID
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process: %w", err)
	}

	// Send a SIGKILL signal to the process
	err = process.Signal(syscall.SIGKILL)
	if err != nil {
		return fmt.Errorf("failed to kill process: %w", err)
	}

	return nil
}

func main() {
	fmt.Println("This project was made for abz.agency as a test project.")

	// This is used to check if the program is already running, if it is, it will not start another instance. This is done via a .lock file.
	lockFile, err := singleinstance.CreateLockFile("proxy.lock")

	// This is used for the command line arguments.
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-start":
			startProxy()
		case "-stop":
			stopProxy()
		case "-quit":
			// Now, I'm not sure if this is what you wanted, but I've decided to kill the first instance if the -quit argument is passed.
			// This is done by reading the PID from the lock file (that singleinstance created) and killing the process.
			if err != nil {
				data, readErr := os.ReadFile("proxy.lock")
				if readErr != nil {
					fmt.Println("A lock file already exists, but failed to read PID:", readErr)
				} else {
					pid, atoiErr := strconv.Atoi(string(data))
					if atoiErr != nil {
						fmt.Println("Failed to convert PID to integer:", atoiErr)
					} else {
						fmt.Printf("An instance already exists with PID, killing the first instance: %d\n", pid)
						killProcess(pid)
					}
				}
				defer lockFile.Close() // Close the lock file
				return
			}
			// This will quit the current instance if the "-quit" argument is passed.
			return
		default:
			fmt.Print("Invalid argument!")
			return
		}
	}

	systray.Run(onReady, onExit) // Start the systray
}

func onReady() {
	systray.SetIcon(icon.Data)          // Sets the icon
	systray.SetTitle("ABZProxyMonitor") // Sets the title of the tray

	startItem := systray.AddMenuItem("Start", "Start") // Creates the menu item with the text "Start"
	stopItem := systray.AddMenuItem("Stop", "Stop")    // Same as above
	quitItem := systray.AddMenuItem("Quit", "Quit")    // Same as above

	// Infinite loop to check if the menu items are clicked
	go func() {
		for {
			select {
			case <-startItem.ClickedCh:
				startProxy() // Start logging

			case <-stopItem.ClickedCh:
				stopProxy() // Stop logging

			case <-quitItem.ClickedCh:
				systray.Quit() // This also works as a global "quit", surprisingly.
				return
			}
		}
	}()
}

func onExit() {
	stopProxy() // Stop logging if the program is closed, eventhough it seems a bit reduntant but why not make use of the onExit() func anyway...
}
