# Velero at project A, backup and restore at other projects

This scenario is introduced in [issue 4806](https://github.com/vmware-tanzu/velero/issues/4806).

Assume the following...

- Project A [project-a]: The project where the Velero's service account is located, and the Velero service account is granted to have enough permission to do backup and restore in the other projects.
- Project B [project-b]: The GCP project we want to restore TO.
- Project C [project-c]: The GCP project we want to restore FROM.

## Set up Velero with permission in projects
* In **project-a**
  * Create "Velero Server" IAM role **role-a** with required role permissions.
  * Create ServiceAccount **sa-a**.
    * Assign **sa-a** with **role-a**.
    * Assign **sa-a** with **role-b**(need to run after role-b created in project-b).
    * Assign **sa-a** with **role-c**(need to run after role-c created in project-c).
  * Create a bucket **bucket-a**.
    * Assign [sa-a] "Storage Object Admin" permissions to [bucket-a]
    * Assign [sa-b] "Storage Object Admin" permissions to [bucket-a](need to run after sa-b created in project-b)
    * Assign [sa-c] "Storage Object Admin" permissions to [bucket-a](need to run after sa-c created in project-c)


* In **project-b**
  * Add the ServiceAccount **sa-a** into project **project-b** according to [Granting service accounts access to your projects](https://cloud.google.com/marketplace/docs/grant-service-account-access).
  * Create ServiceAccount **sa-b**.
  * Create "Velero Server" IAM role **role-b** with required role permissions.
  * Assign **sa-b** with **role-b**.

* In **project-c**
  * Add the ServiceAccount **sa-a** into project **project-c** according to [Granting service accounts access to your projects](https://cloud.google.com/marketplace/docs/grant-service-account-access).
  * Create ServiceAccount **sa-c**.
  * Create "Velero Server" IAM role **role-c** with required role permissions.
  * Assign **sa-c** with **role-c**.

## Backup at project C
* In **project-c**
  * Install Velero on the k8s cluster in this project with configurations:
    * SecretFile: **sa-a**
    * SnapshotLocation: project=**project-a** and volumeProject=**project-c**
    * Bucket: **bucket-a**
  * Create Velero backup **backup-c** with the PVC snapshots desired.

## Restore at project B
* In **project-b**
  * NOTE: Make sure to disable any scheduled backups.
  * Install Velero on the k8s cluster in this project with configurations
    * SecretFile: **sa-a**
    * SnapshotLocation: project=**project-a** and volumeProject=**project-b**
    * Bucket: **bucket-a**
  * Create Velero restore **restore-b** from backup **backup-c**