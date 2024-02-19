package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v2"
)

// 定义YAML文件中相应的Go结构
type ReadsEntry struct {
	Class    string `yaml:"class"`
	Location string `yaml:"location"`
}

type MixcrJson struct {
	Class    string `yaml:"class"`
	Location string `yaml:"location"`
}

type MixcrPath struct {
	Class    string `yaml:"class"`
	Location string `yaml:"location"`
}

type Config struct {
	Reads     []ReadsEntry `yaml:"Reads"`
	MixcrPath MixcrPath    `yaml:"mixcr_path"`
	MixcrJson MixcrJson    `yaml:"mixcr_json"`
	Version   string       `yaml:"version"`
	Threads   int          `yaml:"threads"`
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func main() {
	// 定义命令行参数
	configFlag := flag.String("c", "", "Path to config file")
	workingDirFlag := flag.String("w", "", "Working directory")
	flag.Parse()

	// 检查命令行参数
	if *configFlag == "" || *workingDirFlag == "" {
		fmt.Println("Usage: ./program -c <config_file> -w <working_directory>")
		os.Exit(1)
	}

	// 读取配置文件内容
	data, err := ioutil.ReadFile(*configFlag)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// 将YAML内容解析到Config结构体中
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	// 打印解析结果
	fmt.Printf("--- config:\n%v\n\n", config)
	// 检查文件是否存在
	for _, read := range config.Reads {
		var fileExistsStat bool
		fileExistsStat = fileExists(read.Location)
		if !fileExistsStat {
			fmt.Printf("Error Reads: %s not in path\n", read.Location)
			os.Exit(1)
		}
	}

	// 检查mixcr存在与否
	var mixcrPathFileExists bool
	mixcrPathFileExists = fileExists(config.MixcrPath.Location)
	if !mixcrPathFileExists {
		fmt.Printf("Error Mixcr: %s not in path\n", config.MixcrPath.Location)
		os.Exit(1)
	}

	// 检查json文件是否存在
	var mixcrJsonFileExists bool
	mixcrJsonFileExists = fileExists(config.MixcrJson.Location)
	if !mixcrJsonFileExists {
		fmt.Printf("Error Mixcr reference file: %s not in path\n", config.MixcrJson.Location)
		os.Exit(1)
	}

	// 获取当前用户的 home 目录
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("获取用户 home 目录出错：%v", err)
	}
	// 构建 rabbit mixcr presets 的命令
	err = createNestedFolder(filepath.Join(usr.HomeDir, ".mixcr/presets/"))
	if err != nil {
		fmt.Printf("Create Dir error: %s", err)
	}

	presetCmd := exec.Command(config.MixcrPath.Location,
		"exportPreset", "--preset-name", "bd-sc-xcr-rhapsody-full-length", "-f",
		"-s", "rabbit", filepath.Join(usr.HomeDir, ".mixcr/presets/bd_rabbit_bcr.yaml"))

	// 执行命令并获取输出
	output, err := presetCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("build presets error: %s", output)
		return
	}

	// 输出预设命令的执行结果
	fmt.Printf("Output of build presets command:\n%s\n", output)
	//准备json文件
	err = createNestedFolder(filepath.Join(usr.HomeDir, ".mixcr/libraries/"))
	if err != nil {
		fmt.Printf("Create Dir error: %s", err)
	}

	err = copyFile(config.MixcrJson.Location, filepath.Join(usr.HomeDir, ".mixcr/libraries/imgt.202312-3.sv8_rabbit.json"))
	if err != nil {
		fmt.Printf("Copy json file error: %s", err)
	}

	// 构建 mixcr 命令
	cmdArgs := []string{
		"-Xmx100g",
		"analyze",
		"-t",
		strconv.Itoa(config.Threads),
		"local:bd_rabbit_bcr",
		"--library", "imgt.202312-3.sv8_rabbit.json",
		config.Reads[0].Location,
		config.Reads[1].Location,
		filepath.Join(*workingDirFlag, "Out"),
	}
	err = createNestedFolder(*workingDirFlag)
	// 构建 mixcr 命令
	cmdLine := exec.Command(config.MixcrPath.Location, cmdArgs...)

	// 执行命令并获取输出
	output, err = cmdLine.CombinedOutput()
	if err != nil {
		fmt.Printf("build presets error: %s", output)
		return
	}

	// 输出 mixcr 命令的执行结果
	fmt.Printf("Output of run bcr command:\n%s\n", output)
}

// executeCommand 执行指定的命令和参数
func executeCommand(command string, args []string) {
	cmd := exec.Command(command, args...)
	_, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Error executing command '%s': %s\n", command, err)
		return
	}

	// fmt.Printf("Output of command '%s':\n%s\n", command, output)
}

func copyFile(source, destination string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	return err
}

func createNestedFolder(dirPath string) error {
	err := os.MkdirAll(dirPath, os.ModePerm)
	return err
}
