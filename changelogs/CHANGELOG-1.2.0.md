## All changes

- Upgrade cloud.google.com/go/storage to v1.8.0 (#32, @ashish-amarnath)
- Edit readme velero installation w/ workload identity to include --bucket parameter (#38, @bryanro92)
- Add support for the new `credentialsFile` config key which enables per-BSL credentials. If set, the plugin will use this path as the credentials file for authentication rather than the credentials file path in the environment (#52, @zubron)