package restore

import "k8s.io/apimachinery/pkg/runtime/schema"

var (
	Columns = []string{"Name", "Backup", "Status", "Warnings", "Errors", "Created", "Selector"}
	GVK     = schema.GroupVersionKind{Group: "velero.io", Version: "v1", Kind: "Restore"}
)
