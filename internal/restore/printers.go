package restore

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/pkg/view/component"
	velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
)

func RowPrinter(object unstructured.Unstructured) (component.TableRow, error) {
	restore := &velerov1.Restore{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(object.Object, restore)
	if err != nil {
		return component.TableRow{}, err
	}
	row := component.TableRow{}
	//TODO: Properly format these entries
	row["Name"] = component.NewText(fmt.Sprintf("%s", restore.Name))
	row["Backup"] = component.NewText(fmt.Sprintf("%s", restore.Spec.BackupName))
	row["Status"] = component.NewText(fmt.Sprintf("%s", restore.Status.Phase))
	row["Warnings"] = component.NewText(fmt.Sprintf("%s", restore.Status.Warnings))
	row["Errors"] = component.NewText(fmt.Sprintf("%s", restore.Status.Errors))
	row["Created"] = component.NewText(fmt.Sprintf("%s", restore.CreationTimestamp))
	row["Selector"] = component.NewLabels(restore.Spec.LabelSelector.MatchLabels)
	return row, nil
}
