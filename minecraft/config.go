package minecraft

import (
	"os"
)

// Config contains relevent parameters for ServerControllers
type Config struct {
	Username          string `yaml:"username"`
	RAMMin            int    `yaml:"ram-min"`
	RAMMax            int    `yaml:"ram-max"`
	RootDir           string `yaml:"root-dir"`
	ServerDir         string `yaml:"server-dir"`
	JarName           string `yaml:"jar-name"`
	JavaFlags         string `yaml:"java-params"`
	TmuxPrefix        string `yaml:"tmux-prefix"`
	StartTmuxAttached bool   `yaml:"start-tmux-attached"`
}

// New returns a new Config initialized to default values
func newConfig() (*Config, error) {
	pwd, err := os.Getwd()
	return &Config{
		Username:  "minecraft",
		RAMMin:    4096,
		RAMMax:    4096,
		RootDir:   pwd,
		ServerDir: "servers",
		JarName:   "server.jar",
		JavaFlags: "-XX:+UseG1GC -XX:+ParallelRefProcEnabled -XX:MaxGCPauseMillis=200 -XX:+UnlockExperimentalVMOptions -XX:+DisableExplicitGC -XX:+AlwaysPreTouch -XX:G1NewSizePercent=30 -XX:G1MaxNewSizePercent=40 -XX:G1HeapRegionSize=8M -XX:G1ReservePercent=20 -XX:G1HeapWastePercent=5 -XX:G1MixedGCCountTarget=4 -XX:InitiatingHeapOccupancyPercent=15 -XX:G1MixedGCLiveThresholdPercent=90 -XX:G1RSetUpdatingPauseTimePercent=5 -XX:SurvivorRatio=32 -XX:+PerfDisableSharedMem -XX:MaxTenuringThreshold=1",
	}, err
}
