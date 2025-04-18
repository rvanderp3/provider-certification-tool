# Installation Review

> Note: This document is constantly updated and provides **guidance** to review the installed environment. It's always encouraged to review the product documentation first: [docs.openshift.com][ocp-docs].

This document complements the [official page of "Installing a cluster on any platform"][ocp-install-agnostic] to review specific configurations and components after the cluster has been installed.

This document is also a helper for the ["OPCT - Installation Checklist"][opct-install-checklist] user document.

- [Compute](#compute)
- [Load Balancers](#loadbalancers)
    - [Review the Load Balancer Size](#loadbalancers-size)
    - [Review Health Check configurations](#loadbalancers-healthcheck)
    - [Review Hairpin Traffic](#loadbalancers-hairpin)
- [Components](#components)
    - [etcd](#components-etcd)
        - [Review disk performance with etcd-fio](#components-etcd-ocp-fio)
        - [Review etcd logs: etcd slow requests](#components-etcd-logs-slow-req)
        - [Alternative: Mount /var/lib/etcd in separate disk](#components-etcd-mount)
    - [Image Registry](#components-imageregistry)

## Compute <a name="compute"></a>

- Minimal requirements for Compute nodes: [User Documentation -> Prerequisites][opct-user-guide#prerequisites]

## Load Balancers <a name="loadbalancers"></a>

Review the Load Balancer requirements: [Load balancing requirements for user-provisioned infrastructure][ocp-upi-req-lb]

### Review the Load Balancer size <a name="loadbalancers-size"></a>

The Load Balancer used by the API must support a throughput higher than 100Mbps.

<!-- We haven't this information in the Product Documentation, this minimal was based on the utilization, mainly when installing the cluster (higher than 10Mbps on AWS), and on integrated providers: AWS (NLB) and AlibabaCloud (SLB flavor used by IPI). -->

Reference:

* [AWS](https://github.com/openshift/installer/blob/master/data/data/aws/cluster/vpc/master-elb.tf#L3): NLB (Network Load Balancer)
* [Alibaba](https://github.com/openshift/installer/blob/master/data/data/alibabacloud/cluster/vpc/slb.tf#L49): `slb.s2.small`
* [Azure](https://github.com/openshift/installer/blob/master/data/data/azure/vnet/internal-lb.tf#L7): Standard

### Review the private Load Balancer <a name="loadbalancers"></a>

The basic OpenShift Installations with support of external Load Balancers deploy 3 Load Balancers: public and private for control plane services (Kubernetes API and Machine Config Server), and one public for the ingress.

The DNS or IP address for the private Load Balancer must point to the DNS record `api-int.<cluster>.<domain>`, which will be accessed for internal services.

Reference: [User-provisioned DNS requirements][ocp-upi-req-dns].

### Review Health Check configurations <a name="loadbalancers-healthcheck"></a>

The kube-apiserver has a graceful termination engine that requires the Load Balancer health check probe to the HTTP path.

| Service | Protocol | Port | Path | Threshold | Interval | Timeout |
| -- | -- | -- | -- | -- | -- | -- |
| Kubernetes API Server | HTTPS* | 6443 | /readyz | 2  | 10 | 10 |
| Machine Config Server | HTTPS* | 22623 | /healthz | 2  | 10 | 10 |
| Ingress | TCP | 80 | - | 2  | 10 | 10 |
| Ingress | TCP | 443 | - | 2  | 10 | 10 |

<!-- > Note/Question: Not sure if we need to keep the HTTP (non-SSL on the doc). In the past, I talked with the KAS team and he had plans to remove that option, but due to the limitation of a few cloud providers, it will not. Some providers that still use this: [Alibaba](https://github.com/openshift/installer/blob/master/data/data/alibabacloud/cluster/vpc/slb.tf#L31), [GCP Public](https://github.com/openshift/installer/blob/master/data/data/gcp/cluster/network/lb-public.tf#L20-L21)
*It's required to health check support HTTP protocol. If the Load Balancer used does not support SSL, alternatively and not preferably you can use HTTP - but never TCP:

| Service | Protocol | Port | Path | Threshold | Interval | Timeout |
| -- | -- | -- | -- | -- | -- | -- |
| Kubernetes API Server | HTTP* | 6080 | /readyz | 2  | 10 | 10 |
| Machine Config Server | HTTP* | 22624 | /healthz | 2  | 10 | 10 |

-->


Reminder for the API Load Balancer Health Check:

*"The load balancer must be configured to take a maximum of 30 seconds from the time the API server turns off the /readyz endpoint to the removal of the API server instance from the pool. Within the time frame after /readyz returns an error or becomes healthy, the endpoint must have been removed or added. Probing every 5 or 10 seconds, with two successful requests to become healthy and three to become unhealthy, are well-tested values."* [Load balancing requirements for user-provisioned infrastructure][ocp-upi-req-lb-agnostic].


### Review Hairpin Traffic <a name="loadbalancers-hairpin"></a>

Hairpin traffic is when a backend node's traffic is load-balanced to itself. If this type of network traffic is dropped because your load balancer does not allow hairpin traffic, you need to provide a solution.

On the integrated clouds that do not support hairpin traffic, OpenShift provides a static pod to redirect traffic destined for the load balancer VIP back to the node on the kube-apiserver.

For Reference:

> This is not a recommendation, any solution provided by you will not be supported by Red Hat.

- [Static pods to redirect hairpin traffic for Azure][ocp-src-mco-haiirpin-az]
- [Static pods to redirect hairpin traffic for AlibabaCloud][ocp-src-mco-haiirpin-alc]

Steps to reproduce the Hairpin traffic to a node:

- Deploy one sample pod
- Add one service with a node port
- Create the load balancer with the listener in any port. Example 80
- Create the backend/target group pointing to the node port
- Add the node which the pod is running to the LB/target group/backend nodes
- Try to reach the load balancer IP/DNS through the pod

## Components <a name="components"></a>

### etcd <a name="components-etcd"></a>

Review etcd's disk speed requirements:

- [etcd: Hardware recommendations][etcd-hw-rec]
- [OpenShift Docs: Planning your environment according to object maximums][ocp-perf-obj]
- [OpenShift KCS: Backend Performance Requirements for OpenShift etcd][ocp-kcs-etcd-perf]
- [IBM: Using Fio to Tell Whether Your Storage is Fast Enough for Etcd][ibm-etcd-fio]

#### Review disk performance with etcd-fio <a name="components-etcd-ocp-fio"></a>

The [KCS "How to Use 'fio' to Check Etcd Disk Performance in OCP"][ocp-kcs-fio-etcd] is a guide to check if the disk used by etcd has the expected performance on OpenShift.

<!-- #### Run dense FIO tests

> Note: Keep this section commented as we don't have a strong need to implement or share this broadly.

This section documents how to run dense disk tests using `fio`.

> References:
- https://fio.readthedocs.io/en/latest/fio_doc.html
- https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/benchmark_procedures.html
- https://cloud.google.com/compute/docs/disks/benchmarking-pd-performance
-->

#### Review etcd logs: etcd slow requests <a name="components-etcd-logs-slow-req"></a>

This section provides a guide to check the etcd slow requests from the logs on the etcd pods to understand how the etcd is performing while running the e2e tests.

The command `opct adm parse-etcd-logs` reads the logs, aggregates the requests and displays results in buckets of 100ms increments up to 1s.

`opct adm parse-etcd-logs` is the utility to help troubleshoot the slow requests in the cluster, and help make decisions like changing the flavor of the block device used by the control plane, increasing IOPS, changing the flavor of the instances, etc.

See the command [`opct adm parse-etcd-logs`](../../opct/adm/parse-etcd-logs.md) for more information.

#### Mount /var/lib/etcd in separate disk <a name="components-etcd-mount"></a>

One way to improve the performance on etcd is to use a dedicated block device.

You can mount `/var/lib/etcd` by following the documentation:

- [OpenShift Docs: Disk partitioning][ocp-etcd-isolate]
- [KCS: Mounting separate disk for OpenShift 4 etcd][ocp-kcs]

### Image Registry <a name="components-imageregistry"></a>

You should be able to access the registry and make sure you can push and pull images on it, otherwise, the e2e tests will be reported as failed.

Please check the OpenShift documentation to validate it:

- [Accessing the registry][ocp-registry]
- [Installing a cluster on any platform > Image registry storage configuration][ocp-registry-agnostic]

You can also explore the OpenShift sample projects that create PVC and BuildConfigs (which result in images being built and pushed to image registry). For example:

```bash
oc new-app nodejs-postgresql-persistent
```


[ocp-docs]: https://docs.openshift.com/
[ocp-install-agnostic]: https://docs.openshift.com/container-platform/4.11/installing/installing_platform_agnostic/installing-platform-agnostic.html

[ocp-upi-req-lb]: https://docs.openshift.com/container-platform/4.11/installing/installing_platform_agnostic/installing-platform-agnostic.html#installation-load-balancing-user-infra_installing-platform-agnostic
[ocp-upi-req-lb-agnostic]: https://docs.openshift.com/container-platform/4.11/installing/installing_platform_agnostic/installing-platform-agnostic.html#installation-load-balancing-user-infra_installing-platform-agnostic
[ocp-upi-req-dns]: https://docs.openshift.com/container-platform/4.11/installing/installing_platform_agnostic/installing-platform-agnostic.html#installation-dns-user-infra_installing-platform-agnostic

[ocp-src-mco-haiirpin-az]: https://github.com/openshift/machine-config-operator/blob/master/templates/master/00-master/azure/files/opt-libexec-openshift-azure-routes-sh.yaml
[ocp-src-mco-haiirpin-alc]:https://github.com/openshift/machine-config-operator/tree/master/templates/master/00-master/alibabacloud
[ocp-etcd-isolate]: https://docs.openshift.com/container-platform/4.11/installing/installing_bare_metal/installing-bare-metal.html#installation-user-infra-machines-advanced_disk_installing-bare-metal
[ocp-registry]: https://docs.openshift.com/container-platform/4.11/registry/accessing-the-registry.html
[ocp-registry-agnostic]: https://docs.openshift.com/container-platform/4.11/installing/installing_platform_agnostic/installing-platform-agnostic.html#installation-registry-storage-config_installing-platform-agnostic

[ocp-kcs]: https://access.redhat.com/solutions/5840061

[etcd-hw-rec]: https://etcd.io/docs/v3.5/op-guide/hardware/
[ocp-perf-obj]: https://docs.openshift.com/container-platform/4.11/scalability_and_performance/planning-your-environment-according-to-object-maximums.html
[ocp-kcs-etcd-perf]: https://access.redhat.com/solutions/4770281
[ibm-etcd-fio]: https://www.ibm.com/cloud/blog/using-fio-to-tell-whether-your-storage-is-fast-enough-for-etcd
[ocp-kcs-fio-etcd]: https://access.redhat.com/solutions/4885641

[opct-install-checklist]: ./installation-checklist.md
[opct-user-guide]: ./index.md
