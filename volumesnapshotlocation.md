# Velero Volume Snapshot Location

## Volume Snapshot Location

A volume snapshot location is the location in which to store the volume snapshots created for a backup.

Velero can be configured to take snapshots of volumes from multiple providers. Velero also allows you to configure multiple possible `VolumeSnapshotLocation` per provider, although you can only select one location per provider at backup time.

Each VolumeSnapshotLocation describes a provider + location. These are represented in the cluster via the `VolumeSnapshotLocation` CRD. Velero must have at least one `VolumeSnapshotLocation` per cloud provider.

A sample YAML `VolumeSnapshotLocation` looks like the following:

```yaml
apiVersion: velero.io/v1
kind: VolumeSnapshotLocation
metadata:
  name: gcp-default
  namespace: velero
spec:
  provider: gcp
```

### Parameter Reference

The configurable parameters are as follows:

#### Main config parameters

| Key | Type | Default | Meaning |
| --- | --- | --- | --- |
| `provider` | String `gcp` | Required Field | The name of the cloud provider which will be used to actually store the volume |