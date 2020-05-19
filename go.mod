module github.com/ashish-amarnath/octant-velero-plugin

go 1.13

require (
	github.com/joho/godotenv v1.3.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/vmware-tanzu/octant v0.12.1
	github.com/vmware-tanzu/velero v1.3.2
	k8s.io/apimachinery v0.0.0-20191016225534-b1267f8c42b4
	k8s.io/client-go v0.18.2 // indirect
)

replace k8s.io/client-go => k8s.io/client-go v0.0.0-20190620085101-78d2af792bab
