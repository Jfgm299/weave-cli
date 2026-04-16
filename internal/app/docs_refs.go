package app

const (
	DocsBaseURL          = "https://github.com/Jfgm299/weave-cli/blob/main/docs"
	DocsPathDoctor       = "docs/reference/doctor.md"
	DocsPathConfig       = "docs/reference/config.md"
	DocsPathTransactions = "docs/reference/transactions.md"
	DocsPathProviders    = "docs/reference/providers.md"
	DocsPathMigration    = "docs/reference/migration.md"
)

func DocsURL(path string) string {
	if path == "" {
		return DocsBaseURL
	}
	if len(path) >= 5 && path[:5] == "docs/" {
		return DocsBaseURL + "/" + path[len("docs/"):]
	}
	return DocsBaseURL + "/" + path
}
