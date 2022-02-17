package app

import (
	"fmt"
	"log"
	"os"
)

// Config contains global parameters for wrappers
type Config struct {
	TmuxPrefix        string `yaml:"tmux-prefix"`
	StartTmuxAttached bool   `yaml:"start-tmux-attached"`

	// Global default parameters
	RAMMin    int    `yaml:"ram-min"`
	RAMMax    int    `yaml:"ram-max"`
	RootDir   string `yaml:"root-dir"`
	ServerDir string `yaml:"server-dir"`
	JarName   string `yaml:"jar-name"`
	JavaFlags string `yaml:"java-params"`
}

// NewConfig New returns a new Config initialized to default values
func NewConfig() *Config {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("can't get current dir '%v'", err)
	}

	return &Config{
		RAMMin:    4096,
		RAMMax:    4096,
		RootDir:   pwd,
		ServerDir: "servers",
		JarName:   "server.jar",
		JavaFlags: "-XX:+UseG1GC -XX:+ParallelRefProcEnabled -XX:MaxGCPauseMillis=200 -XX:+UnlockExperimentalVMOptions -XX:+DisableExplicitGC -XX:+AlwaysPreTouch -XX:G1NewSizePercent=30 -XX:G1MaxNewSizePercent=40 -XX:G1HeapRegionSize=8M -XX:G1ReservePercent=20 -XX:G1HeapWastePercent=5 -XX:G1MixedGCCountTarget=4 -XX:InitiatingHeapOccupancyPercent=15 -XX:G1MixedGCLiveThresholdPercent=90 -XX:G1RSetUpdatingPauseTimePercent=5 -XX:SurvivorRatio=32 -XX:+PerfDisableSharedMem -XX:MaxTenuringThreshold=1",
	}
}

func (c Config) GetPath() string {
	return fmt.Sprintf("%s/%s", c.RootDir, c.ServerDir)
}
