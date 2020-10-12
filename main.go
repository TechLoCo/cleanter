package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	targetDir := os.Args[1]

	forbiddenImports, err := readConfig()
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	targetFiles, err := dirwalk(targetDir)
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	cnt := 0

	for _, targetFilePath := range targetFiles {
		imports, err := extractImports(targetFilePath)
		if err != nil {
			fmt.Printf("%+v\n", err)
			os.Exit(1)
		}

		for _, ngLine := range ngLines(imports, forbiddenImports) {
			fmt.Printf("%s の %s は禁止されているimportです\n", targetFilePath, ngLine)
			cnt++
		}
	}

	if cnt > 0 {
		os.Exit(2)
	}
}

func ngLines(imports, forbiddenImports []string) []string {
	var res []string
	for _, target := range imports {
		for _, forbiddenImport := range forbiddenImports {
			if strings.Contains(target, forbiddenImport) {
				res = append(res, target)
			}
		}
	}
	return res
}

func readConfig() ([]string, error) {
	configFile, err := os.Open(".cleanter")
	if err != nil {
		return nil, fmt.Errorf(".cleanterファイルがありません: %w", err)
	}

	reader := bufio.NewReader(configFile)

	var res []string
	for {
		lineBytes, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf(".cleanterファイルの読み込みに失敗しました: %w", err)
		}

		line := string(lineBytes)
		res = append(res, line)
	}
	return res, nil
}

func extractImports(targetFilePath string) ([]string, error) {
	var res []string

	targetFile, err := os.Open(targetFilePath)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(targetFile)

	isStartedImport := false
	for {
		lineBytes, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		line := string(lineBytes)

		if strings.HasPrefix(line, "import (") {
			isStartedImport = true
			continue
		}

		if isStartedImport && strings.HasPrefix(line, ")") {
			isStartedImport = false
			continue
		}

		if isStartedImport {
			res = append(res, strings.TrimSpace(line))
			continue
		}

		if strings.HasPrefix(line, "import") {
			res = append(res, strings.TrimSpace(strings.TrimPrefix(line, "import")))
			continue
		}
	}

	return res, nil
}

func dirwalk(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("ディレクトリを開けませんでした: %w", err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			childFiles, err := dirwalk(filepath.Join(dir, file.Name()))
			if err != nil {
				return nil, fmt.Errorf("ディレクトリのwalkに失敗しました: %w", err)
			}

			paths = append(paths, childFiles...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths, nil
}
