package utils

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

// LoadTmpl 加载指定目录解析为模板 仅限 .html 文件
// eg: views/index.html ,views/subtmpl/subtmpl.html
// tmpl:=LoadTmpl("./views",funcs)
// tmpl.ExecuteTemplate(yourWriter, "index.html", yourData)、tmpl.ExecuteTemplate(yourWriter, "subtmpl/subtmpl.html", yourData)
// 模板名是相对 root 的路径，统一用 / 分隔
func LoadTmpl(root string, funcs template.FuncMap) (*template.Template, error) {
	tmpl := template.New("LoadTmpl")
	if funcs != nil {
		tmpl.Funcs(funcs)
	}
	err := filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 是文件 并且是 .html 结尾
		if fi.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}
		buf, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		// 相对 root 的路径作为模板名，避免带上调用方传入的前缀
		name, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		name = filepath.ToSlash(name)
		if _, err = tmpl.New(name).Parse(string(buf)); err != nil {
			return err
		}
		return nil
	})
	return tmpl, err
}
