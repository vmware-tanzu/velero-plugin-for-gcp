# Velero Backup Storage Locations

## Backup Storage Location

Velero can store backups in a number of locations. These are represented in the cluster via the `BackupStorageLocation` CRD.

Velero must have at least one `BackupStorageLocation`. By default, this is expected to be named `default`, however the name can be changed by specifying `--default-backup-storage-location` on `velero server`.  Backups that do not explicitly specify a storage location will be saved to this `BackupStorageLocation`.

A sample YAML `BackupStorageLocation` looks like the following:

```yaml
apiVersion: velero.io/v1
kind: BackupStorageLocation
metadata:
  name: default
  namespace: velero
spec:
  provider: gcp
  objectStorage:
    bucket: myBucket
```

### Parameter Reference

The configurable parameters are as follows:

#### Main config parameters

| Key | Type | Default | Meaning |
| --- | --- | --- | --- |
| `provider` | String `gcp` | Required Field | The name for the cloud provider which will be used to actually store the backups. |
| `objectStorage` | ObjectStorageLocation | Specification of the object storage for the given provider. |
| `objectStorage/bucket` | String | Required Field | The storage bucket where backups are to be uploaded. |
| `objectStorage/prefix` | String | Optional Field | The directory inside a storage bucket where backups are to be uploaded. |
| `accessMode` | String | `ReadWrite` | How Velero can access the backup storage location. Valid values are `ReadWrite`, `ReadOnly`. |
| `config` | map[string]string<br><br> (See the corresponding GCP specific configs below) | None (Optional) | Configuration keys/values to be passed to the cloud provider for backup storage. |

#### GCP specific

##### config

| Key | Type | Default | Meaning |
| --- | --- | --- | --- |
| `kmsKeyName` | string | Empty | Name of the Cloud KMS key to use to encrypt backups stored in this location, in the form `projects/P/locations/L/keyRings/R/cryptoKeys/K`. See [customer-managed Cloud KMS keys](https://cloud.google.com/storage/docs/encryption/using-customer-managed-keys) for details.  |
| `serviceAccount` | string | Empty | Name of the GCP service account to use for this backup storage location. Specify the service account here if you want to use workload identity instead of providing the key file.
