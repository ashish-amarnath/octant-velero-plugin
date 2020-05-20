package plugin

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"log"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/ashish-amarnath/octant-velero-plugin/internal/backup"
	"github.com/ashish-amarnath/octant-velero-plugin/internal/restore"
	"github.com/ashish-amarnath/octant-velero-plugin/internal/schedule"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

type rowPrinter func(unstructured.Unstructured) (component.TableRow, error)

func getRowPrinterForAccessor(accessor string) (rowPrinter, error) {
	switch accessor {
	case "backup":
		return backup.RowPrinter, nil
	case "restore":
		return restore.RowPrinter, nil
	case "schedule":
		return schedule.RowPrinter, nil
	default:
		return nil, fmt.Errorf("bad accessor, no printer found")
	}
}

func itemToRow(accessor string, object unstructured.Unstructured) component.TableRow {
	printer, err := getRowPrinterForAccessor(accessor)
	if err != nil {
		log.Printf("loading printer error: %s", err)
		return nil
	}
	row, err := printer(object)
	if err != nil {
		log.Printf("printing error: %s", err)
	}
	return row
}

func getColumnsForTable(name string) []string {
	switch name {
	case "backup":
		return backup.Columns
	case "restore":
		return restore.Columns
	case "schedule":
		return schedule.Columns
	default:
		return []string{}
	}
}

func getGVKForAccessor(accessor string) (schema.GroupVersionKind, error) {
	switch accessor {
	case "backup":
		return backup.GVK, nil
	case "restore":
		return restore.GVK, nil
	case "schedule":
		return schedule.GVK, nil
	default:
		return schema.GroupVersionKind{}, fmt.Errorf("bad accessor, no GVK found")
	}
}
