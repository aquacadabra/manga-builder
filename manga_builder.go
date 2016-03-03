package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: builder [data_path]")
		return
	}

	baseFolder := os.Args[1]

	if !exists(baseFolder) {
		fmt.Println("Data folder doesn't exist")
		return
	}

	if exists(path.Join(baseFolder, "pages")) {
		fmt.Println("pages folder already exists")
		return
	}

	if exists(path.Join(baseFolder, "index.html")) {
		fmt.Println("index.html file already exists")
		return
	}

	chapterList := getChapterList(baseFolder)

	err := os.Mkdir(path.Join(baseFolder, "pages"), os.ModeDir)
	if err != nil {
		fmt.Println("Create pages folder error")
		return
	}

	generateIndexFile(baseFolder, chapterList)
	generateChapterFiles(baseFolder, chapterList)

}

func getChapterList(baseFolder string) []string {
	chapterList := []string{}
	files, _ := ioutil.ReadDir(baseFolder)
	for _, f := range files {
		if f.IsDir() {
			chapterList = append(chapterList, f.Name())
		}
	}
	bubleSortStringSlice(chapterList)
	//chapterList = specialEditForPhongVan(chapterList)
	return chapterList
}

func exists(mpath string) bool {
	_, err := os.Stat(mpath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func generateIndexFile(baseFolder string, chapterList []string) {
	indexPageData := struct {
		Title    string
		Chapters []string
	}{
		Title:    "Phong Van",
		Chapters: chapterList,
	}

	indexTemplate, err := template.ParseFiles("template/index.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	indexFile, err := os.Create(path.Join(baseFolder, "index.html"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer indexFile.Close()
	indexTemplate.Execute(indexFile, indexPageData)
}

func getImagesOfChapter(baseFolder, chapterName string) []string {
	imageList := []string{}
	files, _ := ioutil.ReadDir(path.Join(baseFolder, chapterName))
	for _, f := range files {
		if !f.IsDir() && (strings.HasSuffix(f.Name(), ".png") || strings.HasSuffix(f.Name(), ".jpg")) {
			imageList = append(imageList, f.Name())
		}
	}
	sort.StringSlice(imageList).Sort()
	return imageList
}

func generateChapterFiles(baseFolder string, chapterList []string) {
	for index, chapter := range chapterList {
		previousChapter := chapter
		nextChapter := chapter
		if index > 0 {
			previousChapter = chapterList[index-1]
		}
		if index < (len(chapterList) - 1) {
			nextChapter = chapterList[index+1]
		}
		images := getImagesOfChapter(baseFolder, chapter)
		generateChapterFile(baseFolder, images, chapter, previousChapter, nextChapter)
	}
}

func generateChapterFile(baseFolder string, images []string, chapterName string, previousChapter string, nextChapter string) {
	indexPageData := struct {
		Title           string
		PreviousChapter string
		NextChapter     string
		Images          []string
	}{
		Title:           chapterName,
		PreviousChapter: previousChapter,
		NextChapter:     nextChapter,
		Images:          images,
	}

	indexTemplate, err := template.ParseFiles("template/page.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	indexFile, err := os.Create(path.Join(baseFolder, "pages", chapterName+".html"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer indexFile.Close()
	indexTemplate.Execute(indexFile, indexPageData)
}

func bubleSortStringSlice(slice []string) {
	for i := 0; i < len(slice)-1; i++ {
		for j := 1 + i; j < len(slice); j++ {
			if stringCompare(slice[i], slice[j]) > 0 {
				str1 := slice[i]
				slice[i] = slice[j]
				slice[j] = str1
			}
		}
	}
}

func stringCompare(str1, str2 string) int {
	if len(str1) == len(str2) {
		return strings.Compare(str1, str2)
	}
	if len(str1) > len(str2) {
		return 1
	}
	return -1
}

//For special case only
// func specialEditForPhongVan(slice []string) []string {
// 	newSlice := []string{}
// 	newSlice = append(newSlice, slice[0:6]...)
// 	newSlice = append(newSlice, "6.5")
// 	newSlice = append(newSlice, slice[6:len(slice)-1]...)
// 	return newSlice
// }
