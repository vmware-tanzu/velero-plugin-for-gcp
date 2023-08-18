## All changes

  * Add Volume Project setting to support the scenario describe in issue #92 (#149, @blackpiglet)
  * Add support for CSI driver gcp.csi.confidential.cloud (#146, @ps-occrp)
  * Check the "SNAPSHOTS" quota on Google Cloud Platform and do not attempt to create snapshots if the quota is reached. (#144, @0x113)
  * Add cross-project backup and restore functionality in the Velero GCP Plugin. (#143, @Savostov-Arseny)
  * Disable blob signing initialization for non service account file based credentials (#142, @kaovilai)
  * Bump Golang version to v1.20 and add push image to gcr.io in push action. (#138, @blackpiglet)
  * Replace busybox with internal copy binary and fix CVEs. (#137, @blackpiglet)