package parser

import (
	"errors"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"text/template"
)

var ostype = runtime.GOOS

func GenForJava(schema, tableName, projectPkg string) error {
	var err error
	var t *table
	t, err = parseTable(schema, tableName, projectPkg, "model")
	if err != nil {
		log.Println(err)
		return err
	}
	tpls := map[string]string{
		"controller":   "controller.template.tpl",
		"service":      "service.template.tpl",
		"service.impl": "service.impl.template.tpl",
		"dao":          "dao.template.tpl",
		"dao.impl":     "dao.impl.template.tpl",
		"mapper":       "mapper.template.tpl",
		"model":        "model.template.tpl",
	}

	fileSuffix := map[string]string{
		"controller":   "Controller",
		"service":      "Service",
		"service.impl": "ServiceImpl",
		"dao":          "DAO",
		"dao.impl":     "DAOImpl",
		"mapper":       "Mapper",
		"model":        "",
	}

	for k, v := range tpls {
		err = GenJavaFile(t, projectPkg, k, v, fileSuffix[k])

		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func GenJavaFile(t *table, projectPkg, filePkg, tplFile, fileSuffix string) (err error) {
	rootDir := ParentDir(CurrentDir())
	var tmpl *template.Template
	tmpl, err = parseTemplate(tplFile, path.Join(rootDir, "tpl/java"))
	if err != nil {
		log.Println(err)
		return err
	}

	pkgPath := path.Join(rootDir, "gen", strings.ReplaceAll(projectPkg, ".", string(os.PathSeparator)))
	if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
		os.MkdirAll(pkgPath, os.ModePerm)
	}

	modelPath := path.Join(pkgPath, strings.ReplaceAll(filePkg, ".", string(os.PathSeparator)))
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		os.MkdirAll(modelPath, os.ModePerm)
	}

	outPath := path.Join(modelPath, t.ModelName+fileSuffix+".java")
	var file *os.File
	file, err = os.Create(outPath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, t)
}

func CurrentFile() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic(errors.New("Can not get current file info"))
	}
	return file
}

func CurrentDir() string {
	fp := CurrentFile()

	if ostype == "windows" {
		return fp[:strings.LastIndex(fp, "/")]
	}
	//else if ostype == "linux"{
	//	path = path +"/" + "config/"
	//}
	return fp[:strings.LastIndex(fp, string(os.PathSeparator))]
}

func ParentDir(currentDir string) string {
	if ostype == "windows" {
		return currentDir[:strings.LastIndex(currentDir, "/")]
	}

	return currentDir[:strings.LastIndex(currentDir, string(os.PathSeparator))]
}
