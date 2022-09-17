package main

import (
	"fmt"
	"os/exec"
)

func ZipWriter(folderToZip string, outputZipName string) {
	out, err := exec.Command("zip", "-r", outputZipName+".zip", folderToZip).Output()
	if err != nil {
		fmt.Println(err)
	}
	println(string(out))
}
