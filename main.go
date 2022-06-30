package main

import (
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

func xlsx2json(path string, filename string) {
	f, err := excelize.OpenFile(path + "/" + filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheetmap := f.GetSheetMap()
	for _, sheetname := range sheetmap {
		fmt.Println("\n开始操作文件：", sheetname)
		// Get all the rows in the Sheet1.
		rows, err := f.GetRows(sheetname)
		if err != nil {
			fmt.Println(err)
			return
		}

		bigmap := make(map[string]map[string]interface{})

		typeArray := make([]string, len(rows[3]))
		for i := 0; i < len(rows[3]); i++ {
			typeArray[i] = rows[3][i]
		}

		for i := 4; i < len(rows); i++ {
			//先遍历一遍，排除空行
			emptyRow := true
			for j := 0; j < len(rows[i]); j++ {
				if rows[i][j] != "" {
					emptyRow = false
					break
				}
			}
			if emptyRow {
				continue
			}

			rowMap := make(map[string]interface{})
			for j := 0; j < len(rows[i]); j++ {
				rowMap[rows[2][j]] = rows[i][j]
				if typeArray[j] == "int" {
					rowMap[rows[2][j]], _ = strconv.Atoi(rows[i][j])
				}
			}

			//默认第一列为id字段
			bigmap[rows[i][0]] = rowMap
		}

		filepath := fmt.Sprintf("%s/%s.json", path, sheetname)
		f, err1 := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
		if err1 != nil {
			panic(any(err1))
		}
		data, _ := json.Marshal(&bigmap)
		n, err1 := io.WriteString(f, string(data))
		fmt.Println("写入文件：", sheetname, " 字节数：", n, err1)
	}
}

func FileForEachComplete(fileFullPath string) {
	files, err := ioutil.ReadDir(fileFullPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.Name()[0] == '.' {
			continue // 排除 目录文件和隐藏文件
		}
		if file.IsDir() {
			subPath := strings.TrimSuffix(fileFullPath, "/") + "/" + file.Name()
			FileForEachComplete(subPath) // 递归子目录
		} else {
			if path.Ext(file.Name()) == ".xlsx" {
				xlsx2json(fileFullPath, file.Name())
			}
		}
	}
}

func main() {
	FileForEachComplete(".")
}
