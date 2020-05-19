package backup

import "k8s.io/apimachinery/pkg/runtime/schema"

var (
	Columns = []string{"Name", "Status", "Created", "Expires", "Storage Location", "Selector"}
	GVK     = schema.GroupVersionKind{Group: "velero.io", Version: "v1", Kind: "Backup"}
)
