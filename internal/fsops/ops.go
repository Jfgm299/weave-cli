package fsops

type OperationType string

const (
	OpEnsureDir  OperationType = "ensure_dir"
	OpWriteFile  OperationType = "write_file"
	OpCreateLink OperationType = "create_link"
	OpRemovePath OperationType = "remove_path"
)

type Operation struct {
	Type    OperationType
	Path    string
	Target  string
	Content []byte
}
