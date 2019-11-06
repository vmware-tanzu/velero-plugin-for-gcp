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

    # The project ID where snapshots should be stored, if different than the project 
    # that your IAM account is in.
    # 
    # Optional (defaults to the project that the GCP IAM account is in).
    project: my-alternate-project
```
