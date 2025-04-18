package pkg

import (
	"fmt"

	"github.com/vmware-tanzu/sonobuoy/pkg/buildinfo"
)

const (
	PrivilegedClusterRole          = "opct-scc-privileged"
	PrivilegedClusterRoleBinding   = "opct-scc-privileged"
	CertificationNamespace         = "opct"
	VersionInfoConfigMapName       = "opct-version"
	PluginsVarsConfigMapName       = "plugins-config"
	DedicatedNodeRoleLabel         = "node-role.kubernetes.io/tests"
	DedicatedNodeRoleLabelSelector = "node-role.kubernetes.io/tests="
	SonobuoyServiceAccountName     = "sonobuoy-serviceaccount"
	SonobuoyLabelNamespaceName     = "namespace"
	SonobuoyLabelComponentName     = "component"
	SonobuoyLabelComponentValue    = "sonobuoy"
	DefaultToolsRepository         = "quay.io/opct"
	ControllerImage                = "quay.io/opct/opct:v0.6.0-rc.0"
	PluginsImage                   = "plugin-openshift-tests:v0.6.0-rc.0"
	CollectorImage                 = "plugin-artifacts-collector:v0.6.0-rc.0"
	MustGatherMonitoringImage      = "must-gather-monitoring:v0.6.0-rc.0"
	OpenShiftTestsImage            = "image-registry.openshift-image-registry.svc:5000/openshift/tests"
	DedicatedControllerName        = "opct-dedicated-e2e-controller"
)

var (
	SonobuoyImage = fmt.Sprintf("sonobuoy:%s", buildinfo.Version)
)

var (
	SonobuoyDefaultLabels = map[string]string{
		SonobuoyLabelComponentName: SonobuoyLabelComponentValue,
		SonobuoyLabelNamespaceName: CertificationNamespace,
		// Enforcing privileged mode for PSA on Conformance/Sonobuoy environment.
		// https://issues.redhat.com/browse/OPCT-11
		// https://issues.redhat.com/browse/OPCT-31
		"pod-security.kubernetes.io/enforce": "privileged",
		"pod-security.kubernetes.io/audit":   "privileged",
		"pod-security.kubernetes.io/warn":    "privileged",
	}
)

func GetSonobuoyImage() string {
	return fmt.Sprintf("%s/%s", DefaultToolsRepository, SonobuoyImage)
}

func GetPluginsImage() string {
	return fmt.Sprintf("%s/%s", DefaultToolsRepository, PluginsImage)
}

func GetCollectorImage() string {
	return fmt.Sprintf("%s/%s", DefaultToolsRepository, CollectorImage)
}

func GetMustGatherMonitoring() string {
	return fmt.Sprintf("%s/%s", DefaultToolsRepository, MustGatherMonitoringImage)
}
