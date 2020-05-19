package schedule

import "k8s.io/apimachinery/pkg/runtime/schema"

var (
	Columns = []string{"Name", "Status", "Created", "Schedule", "Backup TTL", "Last Backup", "Selector"}
	GVK     = schema.GroupVersionKind{Group: "velero.io", Version: "v1", Kind: "Schedule"}
)
