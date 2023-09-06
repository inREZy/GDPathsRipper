package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

const (
	robtopServer    string = "http://www.boomlings.com/database"
	robtopServerB64 string = "aHR0cDovL3d3dy5ib29tbGluZ3MuY29tL2RhdGFiYXNl"
)

func main() {
	var (
		filePath      string
		serverAddress string
	)

	for {
		fmt.Print("Enter the path to the .exe/libcocos2dcpp.so file of Geometry Dash (or drag & drop the file): ")
		fmt.Scan(&filePath)
		if strings.HasSuffix(filePath, ".exe") || strings.HasSuffix(filePath, ".so") {
			fileArray := strings.Split(filePath, "\\")
			fileFullName := fileArray[len(fileArray)-1]
			file := strings.Split(fileFullName, ".")

			fileData, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Println("The file doesn't exist!")
				continue
			}

			fmt.Print("Enter your server address (33 characters): ")
			fmt.Scan(&serverAddress)
			if len(serverAddress) < 33 || len(serverAddress) > 33 {
				fmt.Printf("Invalid characters count! Your characters count is %d\n", len(serverAddress))
				continue
			}

			data := strings.ReplaceAll(string(fileData), robtopServer, serverAddress)
			dataB64 := strings.ReplaceAll(data, robtopServerB64, base64.StdEncoding.EncodeToString([]byte(serverAddress)))

			os.Mkdir("dist", 0755)

			newFile, _ := os.OpenFile("dist/"+file[0]+"_modified."+file[1], os.O_CREATE|os.O_WRONLY, 0755)
			newFile.Write([]byte(dataB64))

			fmt.Printf("%s file is modified now!\n", fileFullName)
		} else {
			fmt.Println("Your file must have a .exe/.so extension!")
		}
	}
}
