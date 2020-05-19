package backup

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/pkg/view/component"
	velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
)

func RowPrinter(object unstructured.Unstructured) (component.TableRow, error) {
	backup := &velerov1.Backup{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(object.Object, backup)
	if err != nil {
		return component.TableRow{}, err
	}

	row := component.TableRow{}
	//TODO: Properly format these entries
	row["Name"] = component.NewText(fmt.Sprintf("%s", backup.Name))
	row["Status"] = component.NewText(fmt.Sprintf("%s", backup.Status.Phase))
	row["Created"] = component.NewText(fmt.Sprintf("%s", backup.CreationTimestamp))
	row["Expires"] = component.NewText(fmt.Sprintf("%s", backup.Spec.TTL))
	row["Storage Location"] = component.NewText(fmt.Sprintf("%s", backup.Spec.StorageLocation))
	row["Selector"] = component.NewLabels(backup.Spec.LabelSelector.MatchLabels)
	return row, nil
}
