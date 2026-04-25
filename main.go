package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"docgen/parser"
	"docgen/renderer"
	"docgen/utils"
)

// Defaults assume you run from this module directory with a sibling layout:
//   dsa-doc/dsa, dsa-doc/docgen, dsa-doc/dsa-pavilion
const (
	categoryOrderFile = "category_order.txt" // relative to cwd (run from docgen/)
)

func main() {
	var (
		readmeOnly       = flag.Bool("readme-only", false, "Only write readme table; run from docgen cwd (category_order.txt). Does not touch docs/site.")
		codeDir          = flag.String("code", envOr("DOCGEN_CODE", "../dsa"), "Path to the DSA Go module root (category subdirs, e.g. linkedlist/)")
		docsDir          = flag.String("docs", envOr("DOCGEN_DOCS", "../dsa-pavilion/docs"), "Output directory for generated .mdx")
		sidebar          = flag.String("sidebar", envOr("DOCGEN_SIDEBAR", "../dsa-pavilion/sidebars.js"), "Output path for sidebars.js")
		readmeOut        = flag.String("readme", envOr("DOCGEN_README", "../dsa/readme.md"), "Output path for generated readme (code repo index table)")
		sitePrefix       = flag.String("base", envOr("DOCGEN_BASE", "/dsa-pavilion/"), "Site base path for internal doc links (must end with /)")
		githubRepo       = flag.String("github-repo", envOr("DOCGEN_GITHUB_REPO", "https://github.com/ghrushneshr25/dsa"), "DSA repo URL (no trailing slash) for navbar source links")
		githubBranch     = flag.String("github-branch", envOr("DOCGEN_GITHUB_BRANCH", "master"), "Branch name for GitHub blob/tree links")
		githubRoutesPath = flag.String("github-routes", envOr("DOCGEN_GITHUB_ROUTES", "../dsa-pavilion/src/data/dsaGithubRoutes.json"), "Output JSON: doc id -> GitHub URL")
	)
	flag.Parse()

	if *readmeOnly {
		codeDirAbs, err := filepath.Abs(*codeDir)
		if err != nil {
			fmt.Fprintln(os.Stderr, "code path:", err)
			os.Exit(1)
		}
		readmeAbs, err := filepath.Abs(*readmeOut)
		if err != nil {
			fmt.Fprintln(os.Stderr, "readme path:", err)
			os.Exit(1)
		}
		if err := renderer.RenderReadme(codeDirAbs, categoryOrderFile, readmeAbs); err != nil {
			fmt.Fprintln(os.Stderr, "readme:", err)
			os.Exit(1)
		}
		fmt.Println("✅ readme.md generated")
		return
	}

	if *sitePrefix == "" {
		*sitePrefix = "/"
	} else if !strings.HasSuffix(*sitePrefix, "/") {
		*sitePrefix = *sitePrefix + "/"
	}

	codeDirAbs, err := filepath.Abs(*codeDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "code path:", err)
		os.Exit(1)
	}
	docsDirAbs, err := filepath.Abs(*docsDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "docs path:", err)
		os.Exit(1)
	}
	sidebarAbs, err := filepath.Abs(*sidebar)
	if err != nil {
		fmt.Fprintln(os.Stderr, "sidebar path:", err)
		os.Exit(1)
	}
	readmeAbs, err := filepath.Abs(*readmeOut)
	if err != nil {
		fmt.Fprintln(os.Stderr, "readme path:", err)
		os.Exit(1)
	}

	if err := os.RemoveAll(docsDirAbs); err != nil {
		fmt.Fprintln(os.Stderr, "remove docs dir:", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(docsDirAbs, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "mkdir docs dir:", err)
		os.Exit(1)
	}

	var sidebarData []renderer.SidebarCategory

	repoRoot := strings.TrimSuffix(*githubRepo, "/")
	branch := *githubBranch
	githubRoutes := map[string]string{
		"index": repoRoot,
	}

	categories, err := os.ReadDir(codeDirAbs)
	if err != nil {
		fmt.Fprintln(os.Stderr, "read code dir:", err)
		os.Exit(1)
	}
	categories = utils.OrderedDirEntries(categories, categoryOrderFile)

	for _, category := range categories {
		if !category.IsDir() {
			continue
		}

		catName := category.Name()
		if strings.HasPrefix(catName, ".") {
			continue
		}
		catPath := filepath.Join(codeDirAbs, catName)
		docPath := filepath.Join(docsDirAbs, catName)

		files, err := filepath.Glob(filepath.Join(catPath, "*.go"))
		if err != nil {
			fmt.Fprintln(os.Stderr, "glob:", err)
			os.Exit(1)
		}
		sort.Strings(files)
		// Skip repo folders that are not doc categories (e.g. scripts, .github)
		if len(files) == 0 {
			continue
		}

		if err := os.MkdirAll(docPath, 0o755); err != nil {
			fmt.Fprintln(os.Stderr, "mkdir category:", err)
			os.Exit(1)
		}

		var conceptItems []renderer.IndexItem
		var problemItems []renderer.IndexItem
		cat := renderer.SidebarCategory{Name: utils.FormatTitle(catName)}

		for _, file := range files {
			if utils.IsTestFile(file) {
				continue
			}

			meta := parser.ParseMetadata(file)
			docType := meta["type"]
			if docType == "" {
				docType = "problem"
			}

			title := parser.ResolveTitle(meta, docType)
			if title == "" {
				fmt.Fprintln(os.Stderr, "skip (add @problem: or @title:):", file)
				continue
			}

			slug := utils.Slugify(title)

			var sections string
			var code string
			var structs []parser.StructInfo
			var conceptDesc, conceptStructIntro, conceptOpsMD, operationsCode string

			if docType == "concept" {
				cdoc, _ := parser.ParseConceptDocBlock(file)
				structs, _ = parser.ExtractStructs(file)
				if cdoc != nil {
					conceptDesc = cdoc.Description
					conceptStructIntro = cdoc.StructureIntro
					conceptOpsMD = cdoc.Operations
					for i := range structs {
						for k, v := range cdoc.StructureSubsections {
							if strings.EqualFold(strings.TrimSpace(k), structs[i].Name) {
								structs[i].Doc = v
								break
							}
						}
					}
				}
				funcs, _ := parser.ExtractFunctions(file)
				operationsCode = strings.Join(funcs, "\n\n")
			} else {
				sections = parser.ParseSections(file)
				code = parser.ExtractCode(file)
			}

			var subtests []parser.SubtestInfo
			if docType == "problem" {
				if tp, ok := parser.TestFilePath(file); ok {
					subtests, _ = parser.ParseSubtests(tp)
				}
			}
			hasTests := len(subtests) > 0

			outFile := filepath.Join(docPath, slug+".mdx")

			renderer.RenderDoc(renderer.Doc{
				Title:                 title,
				Type:                  docType,
				Sections:              sections,
				Code:                  code,
				Subtests:              subtests,
				HasTests:              hasTests,
				Meta:                  meta,
				Structs:               structs,
				ConceptDescription:    conceptDesc,
				ConceptStructureIntro: conceptStructIntro,
				ConceptOperationsMD:   conceptOpsMD,
				OperationsCode:        operationsCode,
				Output:                outFile,
			})

			docID := catName + "/" + slug
			relSource := filepath.ToSlash(filepath.Join(catName, filepath.Base(file)))
			githubRoutes[docID] = repoRoot + "/blob/" + branch + "/" + relSource

			entry := renderer.IndexItem{
				Title:      title,
				Link:       *sitePrefix + docID,
				Difficulty: meta["difficulty"],
				Tags:       meta["tags"],
			}
			if docType == "concept" {
				cat.Concepts = append(cat.Concepts, docID)
				conceptItems = append(conceptItems, entry)
			} else {
				cat.Problems = append(cat.Problems, docID)
				problemItems = append(problemItems, entry)
			}
		}

		renderer.RenderIndex(renderer.Index{
			Title:        cat.Name,
			Slug:         catName,
			ConceptItems: conceptItems,
			ProblemItems: problemItems,
			Path:         filepath.Join(docPath, "index.mdx"),
		})

		githubRoutes[catName+"/index"] = repoRoot + "/tree/" + branch + "/" + catName

		sidebarData = append(sidebarData, renderer.SidebarCategory{
			Name:     cat.Name,
			Slug:     catName,
			Concepts: cat.Concepts,
			Problems: cat.Problems,
		})
	}

	renderer.RenderSidebar(sidebarData, sidebarAbs)
	renderer.RenderHome(sidebarData, filepath.Join(docsDirAbs, "index.mdx"))

	if err := renderer.RenderReadme(codeDirAbs, categoryOrderFile, readmeAbs); err != nil {
		fmt.Fprintln(os.Stderr, "readme:", err)
		os.Exit(1)
	}

	routesAbs, err := filepath.Abs(*githubRoutesPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "github-routes path:", err)
		os.Exit(1)
	}
	if err := renderer.WriteGithubRoutes(routesAbs, githubRoutes); err != nil {
		fmt.Fprintln(os.Stderr, "github routes:", err)
		os.Exit(1)
	}

	fmt.Println("✅ Docs generated")
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
