![Go](https://github.com/ashish-amarnath/octant-velero-plugin/workflows/Go/badge.svg?branch=master)

# Octant Velero Plugin

_This plugin is still under development and is not ready for use. PRs are welcome_ 😀

This is an This is an [Octant](https://octant.dev/) plugin for [Velero](https://velero.io) which provides dashboadrd for Velero running in your Kubernetes cluster.

## Prerequisites

The following are the prerequisites for using this plugin:
1. Octant should be installed. Installation instructions available [here](https://github.com/vmware-tanzu/octant#installation)
2. [Go 1.13 or above](https://golang.org/dl/)

## Install

### Build from source

This plugin may be installed using the below command:

```bash
make install
```

This command will build the code into an executable `octant-velero-plugin` and places this under `$(HOME)/.config/.octant/plugins`. This is the directory that Octant looks in to discover plugins. 


### Download a releases

There are no releases available at this time. Please check back in later.

## Discussion

Feature requests, bug reports and enhancements are appreciated.
Communication Channes:
1. [Kubernetes Slack](http://slack.k8s.io/) in the [#velero](https://kubernetes.slack.com/app_redirect?channel=C6VCGP4MT) or [#octant](https://kubernetes.slack.com/app_redirect?channel=CM37M9FCG) channels
2. Twitter [@ProjectVelero](https://twitter.com/projectvelero) [@ProjectOctant](https://twitter.com/projectoctant)
3. [Github Issues](https://github.com/ashish-amarnath/octant-velero-plugin/issues/new)

