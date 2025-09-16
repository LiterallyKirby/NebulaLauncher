package main

import (
	"embed"
  "github.com/shirou/gopsutil/process"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"fmt"

)

//go:embed all:frontend/dist
var assets embed.FS

// Your Inject function
func (a *App) Inject(pid int) {
    runtime.LogInfo(a.ctx, fmt.Sprintf("Injecting PID %d", pid))
    // call your injector logic here
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
