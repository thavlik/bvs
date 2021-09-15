local kp = (import 'kube-prometheus/kube-prometheus.libsonnet') +
           (import 'kube-prometheus/kube-prometheus-anti-affinity.libsonnet') +
           (import 'kube-prometheus/kube-prometheus-all-namespaces.libsonnet') + {
  _config+:: {
    namespace: 'monitoring',
    grafana+: {
      folderDashboards+:: {
        BVS: {
          //'bvs.json': (import 'dashboards/bvs.json'),
        },
      },
    },
    prometheus+:: {
      namespaces: [],
    },
  },
  prometheus+:: {
    prometheus+: {
      spec+: {
        retention: '120d',
        //storage: {  // https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#storagespec
        //  volumeClaimTemplate: {  // https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.11/#persistentvolumeclaim-v1-core (defines variable named 'spec' of type 'PersistentVolumeClaimSpec')
        //    apiVersion: 'v1',
        //    kind: 'PersistentVolumeClaim',
        //    spec: {
        //      accessModes: ['ReadWriteOnce'],
        //      // https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.11/#resourcerequirements-v1-core (defines 'requests'),
        //      // and https://kubernetes.io/docs/concepts/policy/resource-quotas/#storage-resource-quota (defines 'requests.storage')
        //      resources: { requests: { storage: '80Gi' } },
        //      // A StorageClass of the following name (which can be seen via `kubectl get storageclass` from a node in the given K8s cluster) must exist prior to kube-prometheus being deployed.
        //      storageClassName: 'do-block-storage',
        //      // The following 'selector' is only needed if you're using manual storage provisioning (https://github.com/coreos/prometheus-operator/blob/master/Documentation/user-guides/storage.md#manual-storage-provisioning).
        //      // And note that this is not supported/allowed by AWS - uncommenting the following 'selector' line (when deploying kube-prometheus to a K8s cluster in AWS) will cause the pvc to be stuck in the Pending status and have the following error:
        //      //  * 'Failed to provision volume with StorageClass "ssd": claim.Spec.Selector is not supported for dynamic provisioning on AWS'
        //      // selector: { matchLabels: {} },
        //    },
        //  },
        //}, // storage
      },
    }
  },
};

{ ['setup/0namespace-' + name]: kp.kubePrometheus[name] for name in std.objectFields(kp.kubePrometheus) } +
{
  ['setup/prometheus-operator-' + name]: kp.prometheusOperator[name]
  for name in std.filter((function(name) name != 'serviceMonitor'), std.objectFields(kp.prometheusOperator))
} +
// serviceMonitor is separated so that it can be created after the CRDs are ready
{ 'prometheus-operator-serviceMonitor': kp.prometheusOperator.serviceMonitor } +
{ ['node-exporter-' + name]: kp.nodeExporter[name] for name in std.objectFields(kp.nodeExporter) } +
{ ['kube-state-metrics-' + name]: kp.kubeStateMetrics[name] for name in std.objectFields(kp.kubeStateMetrics) } +
{ ['alertmanager-' + name]: kp.alertmanager[name] for name in std.objectFields(kp.alertmanager) } +
{ ['prometheus-' + name]: kp.prometheus[name] for name in std.objectFields(kp.prometheus) } +
{ ['prometheus-adapter-' + name]: kp.prometheusAdapter[name] for name in std.objectFields(kp.prometheusAdapter) } +
{ ['grafana-' + name]: kp.grafana[name] for name in std.objectFields(kp.grafana) }