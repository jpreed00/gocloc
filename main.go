package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

const FILE_HEADER string = "File"
const LANG_HEADER string = "Language"
const COMMON_HEADER string = "files          blank        comment           code"
const ROW string = "-------------------------------------------------------------------------" +
	"-------------------------------------------------------------------------" +
	"-------------------------------------------------------------------------"

var rowLen = 79

func main() {
	// parse command line options
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = "gocloc"
	parser.Usage = "[OPTIONS] PATH[...]"

	paths, err := flags.Parse(&opts)
	if err != nil {
		parser.WriteHelp(os.Stdout)
		return
	}
	if len(paths) <= 0 {
		parser.WriteHelp(os.Stdout)
		return
	}

	// setup option for exclude extensions
	ExcludeExts = make(map[string]struct{})
	for _, ext := range strings.Split(opts.ExcludeExt, ",") {
		e, ok := Exts[ext]
		if ok {
			ExcludeExts[e] = struct{}{}
		} else {
			ExcludeExts[ext] = struct{}{}
		}
	}

	// setup option for not match directory
	if opts.NotMatchDir != "" {
		reNotMatchDir = regexp.MustCompile(opts.NotMatchDir)
	}

	// define languages
	action_script := NewLanguage("ActionScript", "//", "/*", "*/")
	asm := NewLanguage("Assembly", "//,;,#,@,|,!", "/*", "*/")
	awk := NewLanguage("Awk", "#", "", "")
	bash := NewLanguage("BASH", "#", "", "")
	batch := NewLanguage("Batch", "REM,rem", "", "")
	c := NewLanguage("C", "//", "/*", "*/")
	c_header := NewLanguage("C Header", "//", "/*", "*/")
	c_sharp := NewLanguage("C#", "//", "/*", "*/")
	c_shell := NewLanguage("C Shell", "#", "", "")
	clojure := NewLanguage("Clojure", ",#,#_", "", "")
	coffee_script := NewLanguage("CoffeeScript", "#", "###", "###")
	cold_fusion := NewLanguage("ColdFusion", "", "<!---", "--->")
	cf_script := NewLanguage("ColdFusion CFScript", "//", "/*", "*/")
	cmake := NewLanguage("CMake", "#", "", "")
	cpp := NewLanguage("C++", "//", "/*", "*/")
	cpp_header := NewLanguage("C++ Header", "//", "/*", "*/")
	css := NewLanguage("CSS", "//", "/*", "*/")
	cuda := NewLanguage("CUDA", "//", "/*", "*/")
	d := NewLanguage("D", "//", "/*", "*/")
	dart := NewLanguage("Dart", "//", "/*", "*/")
	dtrace := NewLanguage("DTrace", "#", "/*", "*/")
	device_tree := NewLanguage("Device Tree", "//", "/*", "*/")
	elixir := NewLanguage("Elixir", "#", "", "")
	erlang := NewLanguage("Erlang", "%", "", "")
	expect := NewLanguage("Expect", "#", "", "")
	lisp := NewLanguage("LISP", ";;", "#|", "|#")
	fortran_legacy := NewLanguage("FORTRAN Legacy", "c,C,!,*", "", "")
	fortran_modern := NewLanguage("FORTRAN Modern", "!", "", "")
	golang := NewLanguage("Go", "//", "/*", "*/")
	haskell := NewLanguage("Haskell", "--", "", "")
	html := NewLanguage("HTML", "<!--", "<!--", "-->")
	jai := NewLanguage("JAI", "//", "/*", "*/")
	java := NewLanguage("Java", "//", "/*", "*/")
	java_script := NewLanguage("JavaScript", "//", "/*", "*/")
	julia := NewLanguage("Julia", "#", "#:=", ":=#")
	json := NewLanguage("JSON", "", "", "")
	jsx := NewLanguage("JSX", "//", "/*", "*/")
	less := NewLanguage("LESS", "//", "/*", "*/")
	linker_script := NewLanguage("LD Script", "//", "/*", "*/")
	lua := NewLanguage("Lua", "--", "--[[", "]]")
	makefile := NewLanguage("make", "#", "", "")
	markdown := NewLanguage("Markdown", "", "", "")
	matlab := NewLanguage("MATLAB", "%", "%{", "}%")
	maven := NewLanguage("Maven", "<!--", "<!--", "-->")
	mustache := NewLanguage("Mustache", "", "{{!", "))")
	m4 := NewLanguage("M4", "#", "", "")
	lex := NewLanguage("lex", "", "/*", "*/")
	objective_c := NewLanguage("Objective C", "//", "/*", "*/")
	objective_cpp := NewLanguage("Objective C++", "//", "/*", "*/")
	ocaml := NewLanguage("OCaml", "", "(*", "*)")
	php := NewLanguage("PHP", "#,//", "/*", "*/")
	pascal := NewLanguage("Pascal", "//,(*", "{", ")")
	plan9_shell := NewLanguage("Plan9 Shell", "#", "", "")
	polly := NewLanguage("Polly", "<!--", "<!--", "-->")
	perl := NewLanguage("Perl", "#", ":=", ":=cut")
	powershell := NewLanguage("PowerShell", "#", "<#", "#>")
	protobuf := NewLanguage("Protocol Buffers", "//", "", "")
	python := NewLanguage("Python", "#", "\"\"\"", "\"\"\"")
	r := NewLanguage("R", "#", "", "")
	racket := NewLanguage("Racket", ";", "#|", "|#")
	ruby := NewLanguage("Ruby", "#", ":=begin", ":=end")
	ruby_html := NewLanguage("Ruby HTML", "<!--", "<!--", "-->")
	rust := NewLanguage("Rust", "//,///,//!", "/*", "*/")
	sass := NewLanguage("Sass", "//", "/*", "*/")
	scala := NewLanguage("Scala", "//", "/*", "*/")
	sed := NewLanguage("sed", "#", "", "")
	sh := NewLanguage("Bourne Shell", "#", "", "")
	scheme := NewLanguage("Scheme", ";", "#|", "|#")
	skill := NewLanguage("SKILL", ";", "/*", "*/")
	sml := NewLanguage("Standard ML", "", "(*", "*)")
	sql := NewLanguage("SQL", "--", "/*", "*/")
	swift := NewLanguage("Swift", "//", "/*", "*/")
	tex := NewLanguage("TeX", "%", "", "")
	text := NewLanguage("Plain Text", "", "", "")
	toml := NewLanguage("TOML", "#", "", "")
	type_script := NewLanguage("TypeScript", "//", "/*", "*/")
	vcproj := NewLanguage("MSBuild script", "<!--", "<!--", "-->")
	vim_script := NewLanguage("Vim script", "\"", "", "")
	wix := NewLanguage("WiX", "<!--", "<!--", "-->")
	xml := NewLanguage("XML", "<!--", "<!--", "-->")
	xsl := NewLanguage("XSLT", "<!--", "<!--", "-->")
	yaml := NewLanguage("YAML", "#", "", "")
	yacc := NewLanguage("Yacc", "//", "/*", "*/")
	zsh := NewLanguage("Zsh", "#", "", "")

	// value for language result
	languages := map[string]*Language{
		"as":          action_script,
		"s":           asm,
		"awk":         awk,
		"bat":         batch,
		"bash":        bash,
		"c":           c,
		"csh":         c_shell,
		"cs":          c_sharp,
		"clj":         clojure,
		"coffee":      coffee_script,
		"cfm":         cold_fusion,
		"cfc":         cf_script,
		"cmake":       cmake,
		"cpp":         cpp,
		"css":         css,
		"cu":          cuda,
		"d":           d,
		"dart":        dart,
		"dtrace":      dtrace,
		"dts":         device_tree,
		"exs":         elixir,
		"erl":         erlang,
		"exp":         expect,
		"lua":         lua,
		"lisp":        lisp,
		"f77":         fortran_legacy,
		"f90":         fortran_modern,
		"go":          golang,
		"h":           c_header,
		"hs":          haskell,
		"hpp":         cpp_header,
		"html":        html,
		"il":          skill,
		"jai":         jai,
		"java":        java,
		"js":          java_script,
		"jl":          julia,
		"json":        json,
		"jsx":         jsx,
		"lds":         linker_script,
		"less":        less,
		"Objective-C": objective_c,
		"md":          markdown,
		"ml":          ocaml,
		"mm":          objective_cpp,
		"makefile":    makefile,
		"Matlab":      matlab,
		"maven":       maven,
		"mustache":    mustache,
		"m4":          m4,
		"lex":         lex,
		"php":         php,
		"pas":         pascal,
		"pl":          perl,
		"text":        text,
		"plan9sh":     plan9_shell,
		"ps1":         powershell,
		"polly":       polly,
		"proto":       protobuf,
		"py":          python,
		"r":           r,
		"rkt":         racket,
		"rb":          ruby,
		"rhtml":       ruby_html,
		"rs":          rust,
		"scala":       scala,
		"scss":        sass,
		"scm":         scheme,
		"sed":         sed,
		"sh":          sh,
		"sml":         sml,
		"sql":         sql,
		"swift":       swift,
		"tex":         tex,
		"sty":         tex,
		"toml":        toml,
		"ts":          type_script,
		"vcproj":      vcproj,
		"vim":         vim_script,
		"wxs":         wix,
		"xml":         xml,
		"xsl":         xsl,
		"yaml":        yaml,
		"y":           yacc,
		"zsh":         zsh,
	}
	fileCache = make(map[string]struct{})

	total := NewLanguage("TOTAL", "", "", "")
	num, maxPathLen := getAllFiles(paths, languages)
	headerLen := 28
	header := LANG_HEADER
	clocFiles := make(map[string]*ClocFile, num)

	// write header
	if opts.Byfile {
		headerLen = maxPathLen + 1
		rowLen = maxPathLen + len(COMMON_HEADER) + 2
		header = FILE_HEADER
	}
	if opts.OutputType == OutputTypeDefault {
		fmt.Printf("%.[2]*[1]s\n", ROW, rowLen)
		fmt.Printf("%-[2]*[1]s %[3]s\n", header, headerLen, COMMON_HEADER)
		fmt.Printf("%.[2]*[1]s\n", ROW, rowLen)
	}

	for _, language := range languages {
		if language.printed {
			continue
		}

		for _, file := range language.files {
			clocFiles[file] = analyzeFile(file, language)

			language.code += clocFiles[file].Code
			language.comments += clocFiles[file].Comments
			language.blanks += clocFiles[file].Blanks
			clocFiles[file].Lang = language.name
		}

		files := int32(len(language.files))
		if len(language.files) <= 0 {
			continue
		}

		language.printed = true

		total.total += files
		total.blanks += language.blanks
		total.comments += language.comments
		total.code += language.code
	}

	// write result
	if opts.Byfile {
		var sortedFiles ClocFiles
		for _, file := range clocFiles {
			sortedFiles = append(sortedFiles, *file)
		}
		sort.Sort(sortedFiles)

		switch opts.OutputType {
		case OutputTypeClocXml:
			t := XmlTotal{
				Code:    total.code,
				Comment: total.comments,
				Blank:   total.blanks,
			}
			f := XmlResultFiles{
				Files: sortedFiles,
				Total: t,
			}
			xmlResult := XmlResult{
				XmlFiles: f,
			}
			xmlResult.Encode()
		case OutputTypeSloccount:
			for _, file := range sortedFiles {
				p := ""
				if strings.HasPrefix(file.Name, "./") || string(file.Name[0]) == "/" {
					splitPaths := strings.Split(file.Name, string(os.PathSeparator))
					if len(splitPaths) >= 3 {
						p = splitPaths[1]
					}
				}
				fmt.Printf("%v\t%v\t%v\t%v\n",
					file.Code, file.Lang, p, file.Name)
			}
		default:
			for _, file := range sortedFiles {
				clocFile := file
				fmt.Printf("%-[1]*[2]s %21[3]v %14[4]v %14[5]v\n",
					maxPathLen, file.Name, clocFile.Blanks, clocFile.Comments, clocFile.Code)
			}
		}
	} else {
		var sortedLanguages Languages
		for _, language := range languages {
			if len(language.files) != 0 && language.printed {
				sortedLanguages = append(sortedLanguages, *language)
			}
		}
		sort.Sort(sortedLanguages)

		for _, language := range sortedLanguages {
			fmt.Printf("%-27v %6v %14v %14v %14v\n",
				language.name, len(language.files), language.blanks, language.comments, language.code)
		}
	}

	// write footer
	if opts.OutputType == OutputTypeDefault {
		fmt.Printf("%.[2]*[1]s\n", ROW, rowLen)
		if opts.Byfile {
			fmt.Printf("%-[1]*[2]v %6[3]v %14[4]v %14[5]v %14[6]v\n",
				maxPathLen, "TOTAL", total.total, total.blanks, total.comments, total.code)
		} else {
			fmt.Printf("%-27v %6v %14v %14v %14v\n",
				"TOTAL", total.total, total.blanks, total.comments, total.code)
		}
		fmt.Printf("%.[2]*[1]s\n", ROW, rowLen)
	}
}
