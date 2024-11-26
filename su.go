package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func init() {
	var Path, fldPath string
	Path, _ = os.Getwd()
	fldPath = filepath.Join(Path, "fld")
	if _, err := os.Stat(fldPath); os.IsNotExist(err) {
		os.Mkdir(fldPath, os.ModePerm)
	}
	for _, dir := range []string{"SYSGO", "Home"} {
		completePath := filepath.Join(fldPath, dir)
		if _, err := os.Stat(completePath); os.IsNotExist(err) {
			os.Mkdir(completePath, os.ModePerm)
		}
	}
	cmdsfile := filepath.Join(fldPath, "SYSGO", "cmds")
	if _, err := os.Stat(cmdsfile); os.IsNotExist(err) {
		f, err := os.OpenFile(cmdsfile, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Println(err)
		}
		f.Close()
	}
	libsfolder := filepath.Join(fldPath, "SYSGO", "libs")
	if _, err := os.Stat(libsfolder); os.IsNotExist(err) {
		os.Mkdir(libsfolder, os.ModePerm)
	}
}

type OS struct {
	setconf map[string]string
	workDir string
}

func NewOS() *OS {
	osInstance := &OS{setconf: make(map[string]string), workDir: ""}
	osInstance.workDir, _ = os.Getwd()
	osInstance.workDir = filepath.Join(osInstance.workDir, "fld", "Home")
	if _, err := os.Stat(osInstance.workDir); os.IsNotExist(err) {
		os.MkdirAll(osInstance.workDir, os.ModePerm)
	}
	return osInstance
}

func (osInstance *OS) loop(line string) {
	reader := bufio.NewReader(os.Stdin)
	strconf := true
	for {
		if line == "" {
			fmt.Printf("\033[36m%s\033[0m\033[31m:\033[0m\033[32m$\033[0m ", filepath.Base(osInstance.workDir))
			var err error
			fmt.Print("\033[32m")
			line, err = reader.ReadString('\n')
			fmt.Print("\033[0m")
			if err != nil {
				fmt.Println()
				break
			}
		} else if strconf == false {
			line = ""
		}
		strconf = false
		commandLines := strings.Split(line, ";")
		for _, commandLine := range commandLines {
			argv := strings.Fields(strings.TrimSpace(commandLine))
			if len(argv) == 0 {
				continue
			}
			if strings.ToLower(argv[0]) == "go" && len(argv) > 1 {
				if argv[1] == ".." {
					parentDir := filepath.Dir(osInstance.workDir)
					if osInstance.workDir == "fld\\Home" || osInstance.workDir == "Home" {
						fmt.Println("\033[31msu: go: You are already at the root directory\n\033[0m")
					} else {
						parentDirParts := strings.Split(parentDir, "\\")
						if parentDirParts[len(parentDirParts)-1] != "fld" {
							osInstance.workDir = parentDir
						} else {
							fmt.Println("\033[31msu: go: You are already at the root directory\n\033[0m")
						}
					}
				} else {
					newDir := filepath.Join(osInstance.workDir, argv[1])
					if _, err := os.Stat(newDir); os.IsNotExist(err) {
						fmt.Printf("\033[31msu: go: %s: No such file or directory\n\033[0m", argv[1])
					} else {
						osInstance.workDir = newDir
					}
				}
			} else if strings.ToLower(argv[0]) == "go.." {
				parentDir := filepath.Dir(osInstance.workDir)
				parentDirParts := strings.Split(parentDir, "\\")
				if parentDirParts[len(parentDirParts)-1] != "fld" {
					osInstance.workDir = parentDir
				} else {
					fmt.Println("\033[31msu: go: You are already at the root directory\n\033[0m")
				}
			} else if strings.ToLower(argv[0]) == "cat" && len(argv) > 1 {
				if runtime.GOOS == "windows" {
					osInstance.workDir = strings.Replace(osInstance.workDir, "/", "\\", -1)
				} else if runtime.GOOS == "linux" {
					osInstance.workDir = strings.Replace(osInstance.workDir, "\\", "/", -1)
				}
				filePath := filepath.Join(osInstance.workDir, argv[1])
				data, err := ioutil.ReadFile(filePath)
				if err != nil {
					fmt.Printf("\033[31msu: cat: %s: No such file or directory\n\033[0m", argv[1])
				} else {
					fmt.Println(string(data))
				}
			} else if strings.ToLower(argv[0]) == "rm" && len(argv) > 1 {
				if runtime.GOOS == "windows" {
					osInstance.workDir = strings.Replace(osInstance.workDir, "/", "\\", -1)
				} else if runtime.GOOS == "linux" {
					osInstance.workDir = strings.Replace(osInstance.workDir, "\\", "/", -1)
				}
				forceDelete := false
				if len(argv) > 2 && argv[1] == "-f" {
					forceDelete = true
					argv = argv[2:]
				} else {
					argv = argv[1:]
				}
				filePath := filepath.Join(osInstance.workDir, argv[0])
				info, err := os.Stat(filePath)
				if err != nil {
					fmt.Printf("\033[31msu: rm: %s: No such file or folder\n\033[0m", argv[0])
				} else if info.IsDir() {
					if forceDelete {
						err := os.RemoveAll(filePath)
						if err != nil {
							fmt.Printf("\033[31msu: rm: %s: Could not remove folder\n\033[0m", argv[0])
						}
					} else {
						files, err := filepath.Glob(filepath.Join(filePath, "*"))
						if err != nil || len(files) > 0 {
							fmt.Printf("\033[31msu: rm: %s: folder is not empty\n\033[0m", argv[0])
						} else {
							err := os.Remove(filePath)
							if err != nil {
								fmt.Printf("\033[31msu: rm: %s: Could not remove folder\n\033[0m", argv[0])
							}
						}
					}
				} else {
					if forceDelete {
						err := os.Remove(filePath)
						if err != nil {
							fmt.Printf("\033[31msu: rm: %s: Could not remove file\n\033[0m", argv[0])
						}
					} else {
						fmt.Printf("\033[31msu: rm: %s: use -f to force delete\n\033[0m", argv[0])
					}
				}
			} else if strings.ToLower(argv[0]) == "hold" && len(argv) > 1 {
				newDir := filepath.Join(osInstance.workDir, argv[1])
				if _, err := os.Stat(newDir); !os.IsNotExist(err) {
					fmt.Printf("\033[31msu: hold: %s: File exists\n\033[0m", argv[1])
				} else {
					err := os.Mkdir(newDir, os.ModePerm)
					if err != nil {
						fmt.Printf("\033[31msu: hold: %s: Could not create directory\n\033[0m", argv[1])
					}
				}
			} else if strings.ToLower(argv[0]) == "touch" && len(argv) > 1 {
				if runtime.GOOS == "windows" {
					osInstance.workDir = strings.Replace(osInstance.workDir, "/", "\\", -1)
				} else if runtime.GOOS == "linux" {
					osInstance.workDir = strings.Replace(osInstance.workDir, "\\", "/", -1)
				}
				filePath := filepath.Join(osInstance.workDir, argv[1])
				f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
				if err != nil {
					fmt.Printf("\033[31msu: touch: %s: Could not create file\n\033[0m", filePath)
				} else {
					f.Close()
				}
			} else if strings.ToLower(argv[0]) == "lf" {
				files, err := filepath.Glob(filepath.Join(osInstance.workDir, "*"))
				if err != nil {
					fmt.Printf("\033[31msu: lf: error listing files or folders\n\033[0m")
				} else {
					for _, file := range files {
						name := filepath.Base(file)
						fileInfo, err := os.Stat(file)
						if err != nil {
							fmt.Printf("\033[31msu: lf: error getting file info for %s\n\033[0m", name)
							continue
						}
						typef := "file"
						if fileInfo.Mode().IsDir() {
							typef = "folder"
						} else if strings.HasSuffix(name, ".su") {
							typef = "su file"
						}
						if _, err := os.Stat(file); err == nil {
							fmt.Println(name, ":", typef)
						}
					}
				}
			} else if strings.ToLower(argv[0]) == "pwf" {
				relativePath, err := filepath.Rel(filepath.Join(osInstance.workDir, "..", "Home"), osInstance.workDir)
				if err != nil {
					fmt.Printf("\033[31msu: pwf: error determining relative path\n\033[0m")
				} else {
					fullPath := strings.Replace(filepath.Join(osInstance.workDir, relativePath), "\\", "/", -1)
					firstHomeIndex := strings.Index(fullPath, "Home")
					if firstHomeIndex != -1 {
						fullPath = fullPath[firstHomeIndex:]
					}
					fmt.Println(fullPath)
				}
			} else if strings.ToLower(argv[0]) == "tx" && len(argv) > 1 {
				if runtime.GOOS == "windows" {
					cmd := exec.Command("notepad.exe", fmt.Sprintf("%s\\%s", strings.ReplaceAll(osInstance.workDir, "\\", "/"), argv[1]))
					cmd.Stdin = os.Stdin
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					err := cmd.Run()
					if err != nil {
						fmt.Printf("\033[31msu: tx: error executing notepad.exe: %s\n\033[0m", err)
					}
				} else if runtime.GOOS == "linux" {
					filePath := filepath.Join(osInstance.workDir, argv[1])
					if _, err := os.Stat(filePath); os.IsNotExist(err) {
						fmt.Printf("\033[31msu: tx: %s: No such file\n\033[0m", argv[1])
					} else {
						cmd := exec.Command("nano", filePath)
						cmd.Stdin = os.Stdin
						cmd.Stdout = os.Stdout
						cmd.Stderr = os.Stderr
						err := cmd.Run()
						if err != nil {
							fmt.Printf("\033[31msu: tx: error executing nano: %s\n\033[0m", err)
						}
					}
				} else {
					fmt.Println("\033[31mUnsupported OS\n\033[0m")
				}
			} else if strings.ToLower(argv[0]) == "time" {
				ticker := time.NewTicker(1 * time.Second)
				defer ticker.Stop()
				for range ticker.C {
					if runtime.GOOS == "windows" {
						cmd := exec.Command("cmd", "/c", "cls")
						cmd.Stdout = os.Stdout
						cmd.Run()
					} else {
						cmd := exec.Command("clear")
						cmd.Stdout = os.Stdout
						cmd.Run()
					}
					now := time.Now()
					fmt.Printf("\x1b[48;5;0m\x1b[38;5;15m %02d:%02d:%02d \x1b[0m", now.Hour(), now.Minute(), now.Second())
					buf := bufio.NewReader(os.Stdin)
					go func() {
						_, _ = buf.ReadString('\n')
						main()
					}()
				}
			} else if argv[0] == "date" {
				now := time.Now()
				fmt.Printf("%02d/%02d/%04d\n", now.Day(), now.Month(), now.Year())
			} else if strings.ToLower(argv[0]) == "clrscr" || strings.ToLower(argv[0]) == "cls" || strings.ToLower(argv[0]) == "clear" || strings.ToLower(argv[0]) == "clr" {
				if runtime.GOOS == "windows" {
					cmd := exec.Command("cmd", "/c", "cls")
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					err := cmd.Run()
					if err != nil {
						fmt.Printf("\033[31mError clearing screen: %s\n\033[0m", err)
					}
				} else {
					cmd := exec.Command("clear")
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					err := cmd.Run()
					if err != nil {
						fmt.Printf("\033[31mError clearing screen: %s\n\033[0m", err)
					}
				}
			} else if strings.ToLower(argv[0]) == "rn" && len(argv) > 2 {
				err := os.Rename(filepath.Join(osInstance.workDir, argv[1]), filepath.Join(osInstance.workDir, argv[2]))
				if err != nil {
					fmt.Printf("\033[31mError renaming: %s\n\033[0m", err)
				}
			} else if strings.ToLower(argv[0]) == "tell" && len(argv) > 1 {
				fmt.Println(strings.Join(argv[1:], " "))
			} else if strings.ToLower(argv[0]) == "mnt" {
				if runtime.GOOS == "windows" {
					cmd := exec.Command("cmd", "/K", "cd", "/D", "c:\\")
					if len(argv) > 1 {
						if argv[1] == "-kf" {
							cmd = exec.Command("cmd", "/K")
						}
					}
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					cmd.Stdin = os.Stdin
					err := cmd.Run()
					if err != nil && err.Error() != "exit status 9009" {
						fmt.Printf("\033[31mError opening cmd: %s\n\033[0m", err)
					}
					clr := exec.Command("cmd", "/c", "cls")
					clr.Stdout = os.Stdout
					clr.Stderr = os.Stderr
					er := clr.Run()
					if er != nil {
						fmt.Printf("\033[31mError clearing screen: %s\n\033[0m", er)
					}
				} else if runtime.GOOS == "linux" {
					cmd := exec.Command("bash", "-c", "cd / && bash")
					if len(argv) > 1 {
						if argv[1] == "-kf" {
							cmd = exec.Command("bash")
						}
					}
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					cmd.Stdin = os.Stdin
					err := cmd.Run()
					if err != nil && err.Error() != "exit status 127" && err.Error() != "exit status 130" {
						fmt.Printf("\033[31mError opening bash: %s\n\033[0m", err)
					}
					clr := exec.Command("clear")
					clr.Stdout = os.Stdout
					clr.Stderr = os.Stderr
					er := clr.Run()
					if er != nil {
						fmt.Printf("\033[31mError clearing screen: %s\n\033[0m", er)
					}
				} else {
					fmt.Println("\033[31mUnsupported OS\n\033[0m")
				}
			} else if strings.ToLower(argv[0]) == "neofech" && len(argv) > 1 {
				fmt.Print("\n")
				neofech := strings.Join(argv[1:], " ")
				neofech = strings.TrimSpace(neofech)
				font := map[string]string{
					"a": " _\n _|\n|_|\n",
					"b": "|\n|_\n|_|\n",
					"c": " _\n|\n|_\n",
					"d": "  |\n _|\n|_|\n",
					"e": " _\n|_\n|_\n",
					"f": " _\n|_\n|\n",
					"g": " _\n|_|\n|_\n",
					"h": "|\n|_\n| |\n",
					"i": " *\n |\n |\n",
					"j": " *\n |\n_|\n",
					"k": "| /\n|<\n| \\ \n",
					"l": "|\n|\n|__\n",
					"m": "!_!_\n| | |\n",
					"n": "!_\n| |\n",
					"o": " _\n| |\n|_|\n",
					"p": " _\n|_|\n|\n",
					"q": " _\n|_|\n  |\n",
					"r": " _\n|_|\n|\\ \n",
					"s": " _\n|_\n _|\n",
					"t": "___\n |\n |\n",
					"u": "| |\n|_|\n",
					"v": "\\   /\n \\_/\n",
					"w": "\\   /\\   /\n \\_/  \\_/\n",
					"x": "\\_/\n/ \\\n",
					"y": "\\_/\n /\n",
					"z": "__\n /\n/_\n",
					" ": "\n\n\n\n",
					"0": " _\n| |\n|_|\n",
					"1": "|\n|\n",
					"2": " _\n _|\n|_\n",
					"3": " _\n _|\n _|\n",
					"4": "|_|\n  |\n",
					"5": " _\n|_\n _|\n",
					"6": " _\n|_\n|_|\n",
					"7": "__\n /\n|\n",
					"8": " _\n|_|\n|_|\n",
					"9": " _\n|_|\n _|\n",
				}
				for _, r := range neofech {
					fmt.Println(font[strings.ToLower(string(r))])
				}
			} else if strings.ToLower(argv[0]) == "snowflake" && len(argv) == 3 && argv[1] == "-i" {
				githubURL := fmt.Sprintf("https://api.github.com/repos/ArtikLamartik/compax-snowflake-lib/contents/%s?ref=main", argv[2])
				resp, err := http.Get(githubURL)
				if err != nil {
					fmt.Printf("\033[31mError getting contents: %s\n\033[0m", err)
					return
				}
				defer resp.Body.Close()
				if resp.StatusCode != http.StatusOK {
					fmt.Printf("\033[31mError getting contents: received status code %d\n\033[0m", resp.StatusCode)
					return
				}

				var contents []map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&contents)
				if err != nil {
					fmt.Printf("\033[31mError decoding JSON: %s\n\033[0m", err)
					return
				}

				destDir := filepath.Join("fld", "SYSGO", "libs", argv[2])
				err = os.MkdirAll(destDir, os.ModePerm)
				if err != nil {
					fmt.Printf("\033[31mError creating directory: %s\n\033[0m", err)
					return
				}

				for _, file := range contents {
					if file["type"] != nil && file["type"] == "file" {
						fileURL := file["download_url"].(string)
						resp, err = http.Get(fileURL)
						if err != nil {
							fmt.Printf("\033[31mError downloading file: %s\n\033[0m", err)
							return
						}
						defer resp.Body.Close()
						if resp.StatusCode != http.StatusOK {
							fmt.Printf("\033[31mError downloading file: received status code %d\n\033[0m", resp.StatusCode)
							return
						}

						filePath := filepath.Join(destDir, file["name"].(string))
						f, err := os.Create(filePath)
						if err != nil {
							fmt.Printf("\033[31mError creating file: %s\n\033[0m", err)
							return
						}
						defer f.Close()
						_, err = io.Copy(f, resp.Body)
						if err != nil {
							fmt.Printf("\033[31mError writing to file: %s\n\033[0m", err)
							return
						}
					}
				}
				fmt.Printf("\033[32mFolder installed successfully.\n\033[0m")
			} else if strings.ToLower(argv[0]) == "snowflake" && len(argv) == 3 && argv[1] == "-c" {
				githubURL := fmt.Sprintf("https://raw.githubusercontent.com/ArtikLamartik/compax-snowflake-lib/main/%s/description.txt", argv[2])
				resp, err := http.Get(githubURL)
				if err != nil {
					fmt.Printf("\033[31mError fetching description: %s\n\033[0m", err)
				}
				defer resp.Body.Close()
				if resp.StatusCode != http.StatusOK {
					fmt.Printf("\033[31mError fetching description: received status code %d\n\033[0m", resp.StatusCode)
				}
				description, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Printf("\033[31mError reading description: %s\n\033[0m", err)
				}
				fmt.Printf("Description to %s:", argv[2])
				fmt.Println()
				descriptionLines := strings.Split(string(description), "\n")
				for _, line := range descriptionLines {
					fmt.Println(line)
				}
			} else if strings.ToLower(argv[0]) == "snowflake" && len(argv) == 3 && argv[1] == "-r" {
				if runtime.GOOS == "windows" {
					executablePath := filepath.Join("fld", "SYSGO", "libs", argv[2], argv[2]+".exe")
					cmd := exec.Command(executablePath)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					err := cmd.Run()
					if err != nil {
						fmt.Printf("\033[31mError running executable, try to restart the compax system: %s\n\033[0m", err)
					}
				} else if runtime.GOOS == "linux" {
					executablePath := filepath.Join("fld", "SYSGO", "libs", argv[2], argv[2])
					cmd := exec.Command(executablePath)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					err := cmd.Run()
					if err != nil {
						fmt.Printf("\033[31mError running executable, try to restart the compax system: %s\n\033[0m", err)
					}
				}
			} else if strings.ToLower(argv[0]) == "snowflake" && len(argv) == 2 && argv[1] == "-l" {
				libDir := filepath.Join("fld", "SYSGO", "libs")
				files, err := os.ReadDir(libDir)
				if err != nil {
					fmt.Printf("\033[31mError reading directory: %s\n\033[0m", err)
				}
				fmt.Println("Available libraries:")
				for _, file := range files {
					if file.IsDir() {
						fmt.Println(file.Name())
					}
				}
			} else if strings.ToLower(argv[0]) == "snowflake" && len(argv) == 3 && argv[1] == "-d" {
				libDir := filepath.Join("fld", "SYSGO", "libs", argv[2])
				err := os.RemoveAll(libDir)
				if err != nil {
					fmt.Printf("\033[31mError deleting library: %s\n\033[0m", err)
				} else {
					fmt.Printf("\033[32mLibrary %s deleted successfully.\n\033[0m", argv[2])
				}
			} else if strings.ToLower(argv[0]) == "exit" {
				os.Exit(0)
			} else if strings.ToLower(argv[0]) == "adcom" && len(argv) >= 3 {
				if argv[1] == "-a" {
					if len(argv) < 4 {
						fmt.Println("\033[31mError: please provide the command and the corresponding action.\n\033[0m")
					} else {
						cmdsPath := filepath.Join("fld", "SYSGO", "cmds")
						f, err := os.OpenFile(cmdsPath, os.O_RDWR|os.O_CREATE, 0666)
						if err != nil {
							fmt.Printf("\033[31mError opening cmds file: %s\n\033[0m", err)
						} else {
							defer f.Close()
							data, err := ioutil.ReadFile(cmdsPath)
							if err != nil {
								fmt.Printf("\033[31mError reading cmds file: %s\n\033[0m", err)
							} else {
								commands := strings.Split(string(data), "\n")
								command := fmt.Sprintf("%s: %s\n", argv[2], strings.Join(argv[3:], " "))
								var commandExists bool
								for _, cmd := range commands {
									cmdParts := strings.SplitN(cmd, ":", 2)
									if len(cmdParts) > 0 && cmdParts[0] == argv[2] {
										commandExists = true
										break
									}
								}
								if !commandExists {
									if _, err := f.WriteString(command); err != nil {
										fmt.Printf("\033[31mError writing to cmds file: %s\n\033[0m", err)
									} else {
										fmt.Printf("\033[32mCommand %s added successfully.\n\033[0m", argv[2])
									}
								} else {
									fmt.Printf("\033[31mCommand with name %s already exists.\n\033[0m", argv[2])
								}
							}
						}
					}
				} else if argv[1] == "-d" {
					cmdsPath := filepath.Join("fld", "SYSGO", "cmds")
					f, err := os.OpenFile(cmdsPath, os.O_WRONLY|os.O_CREATE, 0666)
					if err != nil {
						fmt.Printf("\033[31mError opening cmds file: %s\n\033[0m", err)
					} else {
						defer f.Close()
						data, err := ioutil.ReadFile(cmdsPath)
						if err != nil {
							fmt.Printf("\033[31mError reading cmds file: %s\n\033[0m", err)
						} else {
							lines := strings.Split(string(data), "\n")
							var updatedLines []string
							command := fmt.Sprintf("%s:", argv[2])
							for _, line := range lines {
								if !strings.HasPrefix(line, command) {
									updatedLines = append(updatedLines, line)
								}
							}
							err = ioutil.WriteFile(cmdsPath, []byte(strings.Join(updatedLines, "\n")), 0666)
							if err != nil {
								fmt.Printf("\033[31mError writing to cmds file: %s\n\033[0m", err)
							} else {
								fmt.Printf("\033[32mCommand %s deleted successfully.\n\033[0m", argv[2])
							}
						}
					}
				}
			} else if strings.HasPrefix(argv[0], "./") && strings.HasSuffix(argv[0], ".su") {
				filePath, err := filepath.Abs(filepath.Join(osInstance.workDir, strings.TrimPrefix(argv[0], "./")))
				if err != nil {
					fmt.Printf("\033[31mError obtaining absolute path: %s\n\033[0m", err)
				} else {
					data, err := ioutil.ReadFile(filePath)
					if err != nil {
						fmt.Printf("\033[31mError reading file: %s\n\033[0m", err)
					} else {
						var updatedLines []string
						for _, line := range strings.Split(string(data), "\n") {
							if !strings.HasSuffix(line, ";") {
								updatedLines = append(updatedLines, fmt.Sprintf("%s;", line))
							} else {
								updatedLines = append(updatedLines, line)
							}
						}
						data = []byte(strings.Join(updatedLines, "\n"))
						osInstance.loop(string(data))
					}
				}
			} else if strings.ToLower(argv[0]) == "tsa" {
				fmt.Println("\033[31m♥\n\033[0m")
			} else if strings.ToLower(argv[0]) == "fs" && len(argv) > 1 {
				nameToFind := argv[1]
				filepath.Walk(osInstance.workDir, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if strings.TrimSuffix(info.Name(), filepath.Ext(info.Name())) == nameToFind {
						relativePath, err := filepath.Rel(osInstance.workDir, path)
						relativePath = strings.TrimPrefix(relativePath, "Home" + string(os.PathSeparator))
						relativePath = strings.ReplaceAll(relativePath, "\\", "/")
						if err != nil {
							return err
						}
						if info.IsDir() {
							fmt.Printf("\033[32mHome/%s : folder\n\033[0m", relativePath)
						} else if strings.HasSuffix(relativePath, ".su") {
							fmt.Printf("\033[32mHome/%s : su file\n\033[0m", relativePath)
						} else {
							fmt.Printf("\033[32mHome/%s : file\n\033[0m", relativePath)
						}
					}
					return nil
				})
			} else if strings.ToLower(argv[0]) == "caf" {
				filepath.Walk(osInstance.workDir, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					relativePath, err := filepath.Rel(osInstance.workDir, path)
					relativePath = strings.TrimPrefix(relativePath, "Home" + string(os.PathSeparator))
					relativePath = strings.ReplaceAll(relativePath, "\\", "/")
					if err != nil {
						return err
					}
					if info.IsDir() && relativePath != "." {
						fmt.Printf("Home/%s : folder\n", relativePath)
					} else if strings.HasSuffix(relativePath, ".su") {
						fmt.Printf("Home/%s : su file\n", relativePath)
					} else if relativePath != "." {
						fmt.Printf("Home/%s : file\n", relativePath)
					}
					return nil
				})
				
			} else {
				cmdsPath := filepath.Join("fld", "SYSGO", "cmds")
				data, err := ioutil.ReadFile(cmdsPath)
				if err != nil {
					fmt.Printf("\033[31mError reading cmds file: %s\n\033[0m", err)
				} else {
					lines := strings.Split(string(data), "\n")
					for _, line := range lines {
						if strings.HasPrefix(line, fmt.Sprintf("%s:", strings.ToLower(argv[0]))) {
							inputParts := strings.SplitN(line, ": ", 2)
							if len(inputParts) > 1 {
								osInstance.loop(inputParts[1])
							}
						}
					}
					fmt.Printf("\033[31msu: %s: command not found\n\033[0m", argv[0])
				}
			}
		}
	}
}

func main() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("\033[31mError clearing screen: %s\n\033[0m", err)
		}
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("\033[31mError clearing screen: %s\n\033[0m", err)
		}
	}
	osInstance := NewOS()
	osInstance.loop("")
}