package schedule

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/pkg/view/component"
	velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
)

func RowPrinter(object unstructured.Unstructured) (component.TableRow, error) {
	schedule := &velerov1.Schedule{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(object.Object, schedule)
	if err != nil {
		return component.TableRow{}, err
	}
	row := component.TableRow{}
	//TODO: Properly format these entries
	row["Name"] = component.NewText(fmt.Sprintf("%s", schedule.Name))
	row["Status"] = component.NewText(fmt.Sprintf("%s", schedule.Status.Phase))
	row["Created"] = component.NewText(fmt.Sprintf("%s", schedule.CreationTimestamp))
	row["Backup TTL"] = component.NewText(fmt.Sprintf("%s", schedule.Spec.Template.TTL))
	row["Last Backup"] = component.NewText(fmt.Sprintf("%s", schedule.Status.LastBackup))
	row["Selector"] = component.NewLabels(schedule.Spec.Template.LabelSelector.MatchLabels)
	return row, nil
}
