package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

const (
	version      string = "1.1"
	robtopServer string = "http://www.boomlings.com/database"
)

func scanData() string {
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	return input.Text()
}

func main() {
	fmt.Printf("[ GDPathsRipper v%s created by inREZy. ]\n\n", version)
	for {
		fmt.Print("Enter the path to the .exe/libcocos2dcpp.so file of Geometry Dash (or drag & drop the file): ")
		filePath := scanData()
		filePath = strings.ReplaceAll(filePath, "\"", "")
		if strings.HasSuffix(filePath, ".exe") || strings.HasSuffix(filePath, ".so") {
			fileArray := strings.Split(filePath, "\\")
			fileName := fileArray[len(fileArray)-1]

			fileData, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("[ERROR] %v\n", err)
				continue
			}

			fmt.Print("Enter your server address (33 characters): ")
			serverAddress := scanData()
			if len(serverAddress) != 33 {
				fmt.Printf("[ERROR] Your characters count is %d\n", len(serverAddress))
				continue
			}

			data := strings.ReplaceAll(string(fileData), robtopServer, serverAddress)
			data = strings.ReplaceAll(data, base64.StdEncoding.EncodeToString([]byte(robtopServer)), base64.StdEncoding.EncodeToString([]byte(serverAddress)))

			os.Mkdir("dist", 0755)

			newFile, err := os.OpenFile("dist/"+fileName, os.O_CREATE|os.O_WRONLY, 0755)
			if err != nil {
				fmt.Printf("[ERROR] %v\n", err)
				continue
			}
			newFile.Write([]byte(data))
			newFile.Close()

			fmt.Printf("%s is successfully modified!\n", fileName)
		} else {
			fmt.Print("[ERROR] Your file should have .exe/.so extension!\n")
			continue
		}
	}
}
