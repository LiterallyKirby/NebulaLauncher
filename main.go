package main

import (
	"embed"
  "github.com/shirou/gopsutil/process"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"fmt"

	"io"
	"net/http"
	"os"
	"os/exec"

)

//go:embed all:frontend/dist
var assets embed.FS

// Your Inject function
func (a *App) Inject(pid int) {
	// Log start
	runtime.LogInfo(a.ctx, fmt.Sprintf("Injecting PID %d...", pid))
	runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("Injecting PID %d...", pid))

	// Download the library from GitHub
	libURL := "https://github.com/your/repo/raw/main/libmod.so"
	tmpPath := "/tmp/libmod.so"

	resp, err := http.Get(libURL)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Failed to download library: %s", err))
		runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("Failed to download library: %s", err))
		return
	}
	defer resp.Body.Close()

	out, err := os.Create(tmpPath)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Failed to create file: %s", err))
		runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("Failed to create file: %s", err))
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Failed to write library: %s", err))
		runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("Failed to write library: %s", err))
		return
	}

	runtime.EventsEmit(a.ctx, "log", "Library downloaded successfully!")

	// Call nebula_injector
	cmd := exec.Command("nebula_injector", tmpPath, fmt.Sprintf("%d", pid), "init")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Injection failed: %s", err))
		runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("Injection failed: %s", err))
		return
	}

	runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("PID %d injected successfully!", pid))
	runtime.LogInfo(a.ctx, fmt.Sprintf("PID %d injected successfully!", pid))
}

func (a *App) GetProcesses() ([]ProcessInfo, error) {
    procs, err := process.Processes()
    if err != nil {
        return nil, err
    }

    var list []ProcessInfo
    for _, p := range procs {
        name, err := p.Name()
        if err != nil {
            continue
        }
        list = append(list, ProcessInfo{
            PID:  p.Pid,
            Name: name,
        })
    }
    return list, nil
}

type ProcessInfo struct {
    PID  int32  `json:"pid"`
    Name string `json:"name"`
}

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "NebulaLauncher",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
