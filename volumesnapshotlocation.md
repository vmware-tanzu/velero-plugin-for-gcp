# Volume Snapshot Location

The following sample GCP `VolumeSnapshotLocation` YAML shows all of the configurable parameters. The items under `spec.config` can be provided as key-value pairs to the `velero install` command's `--snapshot-location-config` flag -- for example, `snapshotLocation=us-central1,project=my-project,...`.

```yaml
apiVersion: velero.io/v1
kind: VolumeSnapshotLocation
metadata:
  name: gcp-default
  namespace: velero
spec:
  # Name of the volume snapshotter plugin to use to connect to this location.
  #
  # Required.
  provider: velero.io/gcp
  
  config:
    # The GCP location where snapshots should be stored. See the GCP documentation
    # (https://cloud.google.com/storage/docs/locations#available_locations) for the
    # full list. If not specified, snapshots are stored in the default location
    # (https://cloud.google.com/compute/docs/disks/create-snapshots#default_location).
    #
    # Optional.
    snapshotLocation: us-central1

    # The project ID where snapshots should be stored.
    #
    # Optional (defaults to the project that the GCP IAM account is in).
    snapshotProject: my-project

    # The project ID where existing snapshots should be retrieved from during restores, if 
    # different than the project that your IAM account is in. This field has no effect on 
    # where new snapshots are created; it is only useful for restoring existing snapshots 
    # from a different project.
    # 
    # Optional (defaults to the project that the GCP IAM account is in).
    restoreFromProject: my-alternate-project
```
