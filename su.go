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
	wossfile := filepath.Join(fldPath, "SYSGO", "woss.su")
	if _, err := os.Stat(wossfile); os.IsNotExist(err) {
		f, err := os.OpenFile(wossfile, os.O_CREATE|os.O_WRONLY, 0666)
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
	osInstance.setconf["go"] = "you give it a folder name and it goes there, if no argument it goes to the root folder, and .. goes to the last folder you were in"
	osInstance.setconf["cat"] = "you give it a file name and it gives you the contents of the file"
	osInstance.setconf["rm"] = "you give it a file name and it deletes the file, you can give it a -f flag to force delete, and -a flag to delete all files in the folder you are in"
	osInstance.setconf["hold"] = "you give it a folder name and it creates the folder"
	osInstance.setconf["touch"] = "you give it a file name and it creates the file"
	osInstance.setconf["lf"] = "it lists the files in the folder you are in"
	osInstance.setconf["pwf"] = "it prints the current working folder"
	osInstance.setconf["tx"] = "it will use your text editor to edit the file you give the name as an argument"
	osInstance.setconf["date"] = "it will print the date and the time"
	osInstance.setconf["clear"] = "it clears the screen"
	osInstance.setconf["cls"] = "it clears the screen"
	osInstance.setconf["clrscr"] = "it clears the screen"
	osInstance.setconf["clr"] = "it clears the screen"
	osInstance.setconf["rn"] = "it renames the file you give the name as an argument to the second name you give it"
	osInstance.setconf["tell"] = "it will print the argument you give it"
	osInstance.setconf["mnt"] = "it open a the shell on your OS and if you type exit it will close the shell and also if you give the argument -kf it will open in the folder the executable is in"
	osInstance.setconf["nefech"] = "it will create ascii art from text you give it"
	osInstance.setconf["snowflake"] = "if you use -i it installes the next argument you gave it from github, if you use -c it will read the description of the next argument you gave it from github, if you use -d it will delete the next argument you gave it from github, if you use -r it will run the next argument you gave it, if it is installed"
	osInstance.setconf["./"] = "it will run the next argument you gave it, if it is a su file from the Home folder"
	osInstance.setconf["././"] = "it will run the next argument you gave it, if it is a su file from the folder you are in"
	osInstance.setconf["cap"] = "it lists all the paths in the folder you are in"
	osInstance.setconf["sf"] = "it searches for the next argument you gave it in the folder you are in"
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
				if strings.Contains(argv[1], "..") {
					parentDir := filepath.Dir(osInstance.workDir)
					if osInstance.workDir == "fld\\Home" || osInstance.workDir == "Home" {
						fmt.Println("\033[31msu: go: You are already at the root directory\033[0m")
					} else {
						parentDirParts := strings.Split(parentDir, "\\")
						if parentDirParts[len(parentDirParts)-1] != "fld" {
							osInstance.workDir = parentDir
						} else {
							fmt.Println("\033[31msu: go: You are already at the root directory\033[0m")
						}
					}
				} else {
					newDir := filepath.Join(osInstance.workDir, argv[1])
					if _, err := os.Stat(newDir); os.IsNotExist(err) {
						fmt.Printf("\033[31msu: go: %s: No such file or directory\n\033[0m", argv[1])
					} else if fi, err := os.Stat(newDir); err == nil && fi.Mode().IsDir() {
						osInstance.workDir = newDir
					}
				}
			} else if strings.ToLower(argv[0]) == "go" && len(argv) == 1 {
				osInstance.workDir = filepath.Join("fld", "Home")
			} else if strings.ToLower(argv[0]) == "go.." {
				parentDir := filepath.Dir(osInstance.workDir)
				parentDirParts := strings.Split(parentDir, "\\")
				if parentDirParts[len(parentDirParts)-1] != "fld" {
					osInstance.workDir = parentDir
				} else {
					fmt.Println("\033[31msu: go: You are already at the root directory\033[0m")
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
				all := false
				for _, arg := range argv {
					if strings.ToLower(arg) == "-f" {
						forceDelete = true
					}
					if strings.ToLower(arg) == "-a" {
						all = true
					}
				}
				if all {
					files, err := filepath.Glob(filepath.Join(osInstance.workDir, "*"))
					if err != nil {
						fmt.Printf("\033[31msu: rm: error getting files in folder\n\033[0m")
					} else {
						for _, file := range files {
							err := os.RemoveAll(file)
							if err != nil {
								fmt.Printf("\033[31msu: rm: %s: Could not remove file/folder\n\033[0m", filepath.Base(file))
							}
						}
					}
				} else {
					if forceDelete {
						err := os.RemoveAll(filepath.Join(osInstance.workDir, argv[2]))
						if err != nil {
							fmt.Printf("\033[31msu: rm: %s: Could not remove file/folder\n\033[0m", argv[2])
						}
					} else {
						err := os.Remove(filepath.Join(osInstance.workDir, argv[1]))
						if err != nil {
							if os.IsNotExist(err) {
								fmt.Printf("\033[31msu: rm: %s: No such file or directory\n\033[0m", argv[1])
							} else {
								fmt.Print("\033[31msu: rm: Folder is not empty, use the -f to force delete\n\033[0m")
							}
						}

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
							last_modified := fileInfo.ModTime().Format("02 Jan 2006 15:04:05")
							var folderSize int64
							err := filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
								if err != nil {
									return err
								}
								if !info.IsDir() {
									folderSize += info.Size()
								}
								return nil
							})
							if err != nil {
								fmt.Printf("\033[31msu: lf: error getting folder size for %s\n\033[0m", name)
								continue
							}
							space := fmt.Sprintf("%d bytes", folderSize)
							fmt.Printf("\033[32m[ %s ] \033[36m%s \033[34m{ %s } \033[35m( %s )\n\033[0m", typef, name, last_modified, space)
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
					fmt.Printf("\033[36m%s\033[0m\n", fullPath)
				}
			} else if strings.ToLower(argv[0]) == "tx" && len(argv) > 1 {
				if argv[1] == "./woss.su" {
					if runtime.GOOS == "windows" {
						cmd := exec.Command("notepad.exe", "fld\\SYSGO\\woss.su")
						cmd.Stdin = os.Stdin
						cmd.Stdout = os.Stdout
						cmd.Stderr = os.Stderr
						err := cmd.Run()
						if err != nil {
							fmt.Printf("\033[31msu: tx: error executing notepad.exe: %s\n\033[0m", err)
						}
					} else if runtime.GOOS == "linux" {
						filePath := filepath.Join("fld\\SYSGO\\woss.su")
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
						fmt.Println("\033[31mUnsupported OS\033[0m")
					}
				} else if runtime.GOOS == "windows" {
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
					fmt.Println("\033[31mUnsupported OS\033[0m")
				}
			} else if argv[0] == "date" {
				now := time.Now()
				fmt.Printf("\033[34m%02d/%02d/%04d\n%02d:%02d:%02d\n\033[0m", now.Day(), now.Month(), now.Year(), now.Hour(), now.Minute(), now.Second())
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
				fmt.Print(strings.ReplaceAll(strings.Join(argv[1:], " "), "\\n", "\n"))
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
					fmt.Println("\033[31mUnsupported OS\033[0m")
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
			} else if strings.HasPrefix(argv[0], "./") && strings.HasSuffix(argv[0], ".su") {
				if argv[0] == "./woss.su" {
					data, err := ioutil.ReadFile(filepath.Join("fld", "SYSGO", "woss.su"))
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
				var filePath string
				var err error
				if strings.HasPrefix(argv[0], "././") {
					filePath, err = filepath.Abs(filepath.Join(filepath.Dir(osInstance.workDir), filepath.Base(osInstance.workDir), strings.TrimPrefix(argv[0], "././")))
				} else {
					filePath, err = filepath.Abs(filepath.Join("fld/Home", strings.TrimPrefix(argv[0], "./")))
				}
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
				fmt.Print("\033[31m♥\nurt\n\033[0m")
			} else if strings.ToLower(argv[0]) == "cap" {
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
						var folderSize int64
						err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
							if err != nil {
								return err
							}
							if !info.IsDir() {
								folderSize += info.Size()
							}
							return nil
						})
						if err != nil {
							fmt.Printf("\033[31msu: lf: error getting folder size for %s: %v\n\033[0m", relativePath, err)
						}
						space := fmt.Sprintf("%d bytes", folderSize)
						last_modified := info.ModTime().Format("02 Jan 2006 15:04:05")
						fmt.Printf("\033[32m[ folder ] \033[36m%s \033[34m{ %s } \033[35m( %s )\n\033[0m", relativePath, last_modified, space)
					} else if strings.HasSuffix(relativePath, ".su") {
						var folderSize int64
						err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
							if err != nil {
								return err
							}
							if !info.IsDir() {
								folderSize += info.Size()
							}
							return nil
						})
						if err != nil {
							fmt.Printf("\033[31msu: lf: error getting folder size for %s: %v\n\033[0m", relativePath, err)
						}
						space := fmt.Sprintf("%d bytes", folderSize)
						last_modified := info.ModTime().Format("02 Jan 2006 15:04:05")
						fmt.Printf("\033[32m[ su file ] \033[36m%s \033[34m{ %s } \033[35m( %s )\n\033[0m", relativePath, last_modified, space)
					} else if relativePath != "." {
						var folderSize int64
						err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
							if err != nil {
								return err
							}
							if !info.IsDir() {
								folderSize += info.Size()
							}
							return nil
						})
						if err != nil {
							fmt.Printf("\033[31msu: lf: error getting folder size for %s: %v\n\033[0m", relativePath, err)
						}
						space := fmt.Sprintf("%d bytes", folderSize)
						last_modified := info.ModTime().Format("02 Jan 2006 15:04:05")
						fmt.Printf("\033[32m[ file ] \033[36m%s \033[34m{ %s } \033[35m( %s )\n\033[0m", relativePath, last_modified, space)
					}
					return nil
				})
			} else if strings.ToLower(argv[0]) == "sf" && len(argv) > 1 {
				nameToFind := argv[1]
				err := filepath.Walk(osInstance.workDir, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						fmt.Printf("\033[31mError accessing path %s: %v\n\033[0m", path, err)
						return err
					}
					if strings.EqualFold(strings.TrimSuffix(info.Name(), filepath.Ext(info.Name())), strings.TrimSuffix(nameToFind, filepath.Ext(nameToFind))) {
						relativePath, err := filepath.Rel(osInstance.workDir, path)
						if err != nil {
							fmt.Printf("\033[31mError calculating relative path: %v\n\033[0m", err)
							return err
						}
						relativePath = strings.TrimPrefix(relativePath, "Home"+string(os.PathSeparator))
						relativePath = strings.ReplaceAll(relativePath, "\\", "/")
						if info.IsDir() {
							var folderSize int64
							err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
								if err != nil {
									return err
								}
								if !info.IsDir() {
									folderSize += info.Size()
								}
								return nil
							})
							if err != nil {
								fmt.Printf("\033[31msu: lf: error getting folder size for %s: %v\n\033[0m", relativePath, err)
							}
							space := fmt.Sprintf("%d bytes", folderSize)
							last_modified := info.ModTime().Format("02 Jan 2006 15:04:05")
							fmt.Printf("\033[32m[ folder ] \033[36m/%s \033[34m{ %s } \033[35m( %s )\n\033[0m", relativePath, last_modified, space)
						} else if strings.HasSuffix(relativePath, ".su") {
							var folderSize int64
							err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
								if err != nil {
									return err
								}
								if !info.IsDir() {
									folderSize += info.Size()
								}
								return nil
							})
							if err != nil {
								fmt.Printf("\033[31msu: lf: error getting folder size for %s: %v\n\033[0m", relativePath, err)
							}
							space := fmt.Sprintf("%d bytes", folderSize)
							last_modified := info.ModTime().Format("02 Jan 2006 15:04:05")
							fmt.Printf("\033[32m[ su file ] \033[36m%s \033[34m{ %s } \033[35m( %s )\n\033[0m", relativePath, last_modified, space)
						} else {
							var folderSize int64
							err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
								if err != nil {
									return err
								}
								if !info.IsDir() {
									folderSize += info.Size()
								}
								return nil
							})
							if err != nil {
								fmt.Printf("\033[31msu: lf: error getting folder size for %s: %v\n\033[0m", relativePath, err)
							}
							space := fmt.Sprintf("%d bytes", folderSize)
							last_modified := info.ModTime().Format("02 Jan 2006 15:04:05")
							fmt.Printf("\033[32m[ file ] \033[36m%s \033[34m{ %s } \033[35m( %s )\n\033[0m", relativePath, last_modified, space)
						}
					}
					if !info.IsDir() {
						fileContent, err := ioutil.ReadFile(path)
						if err != nil {
							fmt.Printf("\033[31mError reading file: %v\n\033[0m", err)
							return err
						}
						if strings.Contains(string(fileContent), nameToFind) {
							relativePath, err := filepath.Rel(osInstance.workDir, path)
							if err != nil {
								fmt.Printf("\033[31mError calculating relative path: %v\n\033[0m", err)
								return err
							}
							relativePath = strings.TrimPrefix(relativePath, "Home"+string(os.PathSeparator))
							relativePath = strings.ReplaceAll(relativePath, "\\", "/")
							fmt.Printf("\033[32m[ file containing ] \033[36mHome/%s\n\033[0m", relativePath)
						}
					}
					return nil
				})
				if err != nil {
					fmt.Printf("\033[31mError walking the directory: %v\n\033[0m", err)
				}
			} else if strings.ToLower(argv[0]) == "help" {
				if len(argv) > 1 {
					if value, ok := osInstance.setconf[strings.ToLower(argv[1])]; ok {
						fmt.Printf("%s\n", value)
					} else {
						fmt.Printf("\033[31mError: %s is not a valid option\n\033[0m", argv[1])
					}
				} else if len(argv) == 1 {
					for key, value := range osInstance.setconf {
						fmt.Printf("%s: %s\n", key, value)
					}
				}
			} else {
				fmt.Printf("\033[31msu: %s: command not found\n\033[0m", argv[0])
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
	osInstance.loop("./woss.su")
}