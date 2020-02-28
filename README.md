# Plugins for Google Cloud Platform (GCP)

## Overview

This repository contains these plugins to support running Velero on GCP:

- An object store plugin for persisting and retrieving backups on Google Cloud Storage. Content of backup is log files, warning/error files, restore logs.

- A volume snapshotter plugin for creating snapshots from volumes (during a backup) and volumes from snapshots (during a restore) on Google Compute Engine Disks.

You can run Kubernetes on Google Cloud Platform in either:

* Kubernetes on Google Compute Engine virtual machines
* Google Kubernetes Engine

For common use-cases, take a look at the [Examples][10] page.

## Compatibility

Below is a listing of plugin versions and respective Velero versions that are compatible.

| Plugin Version  | Velero Version |
|-----------------|----------------|
| v1.0.x          | v1.2.0         |

## Setup

To set up Velero on GCP, you:

* [Create an GCS bucket][1]
* [Set permissions for Velero][2]
* [Install and start Velero][3]

If you do not have the `gcloud` and `gsutil` CLIs locally installed, follow the [user guide][5] to set them up.

## Create an GCS bucket

Velero requires an object storage bucket in which to store backups, preferably unique to a single Kubernetes cluster (see the [FAQ][11] for more details). Create a GCS bucket, replacing the <YOUR_BUCKET> placeholder with the name of your bucket:

```bash
BUCKET=<YOUR_BUCKET>

gsutil mb gs://$BUCKET/
```

## Set permissions for Velero

If you run Google Kubernetes Engine (GKE), make sure that your current IAM user is a cluster-admin. This role is required to create RBAC objects.
See [the GKE documentation][22] for more information.

### Option 1: Set permissions with a Service Account

To integrate Velero with GCP, create a Velero-specific [Service Account][15]:

1. View your current config settings:

    ```bash
    gcloud config list
    ```

    Store the `project` value from the results in the environment variable `$PROJECT_ID`.

    ```bash
    PROJECT_ID=$(gcloud config get-value project)
    ```

2. Create a service account:

    ```bash
    gcloud iam service-accounts create velero \
        --display-name "Velero service account"
    ```

    > If you'll be using Velero to backup multiple clusters with multiple GCS buckets, it may be desirable to create a unique username per cluster rather than the default `velero`.

    Then list all accounts and find the `velero` account you just created:

    ```bash
    gcloud iam service-accounts list
    ```

    Set the `$SERVICE_ACCOUNT_EMAIL` variable to match its `email` value.

    ```bash
    SERVICE_ACCOUNT_EMAIL=$(gcloud iam service-accounts list \
      --filter="displayName:Velero service account" \
      --format 'value(email)')
    ```

3. Attach policies to give `velero` the necessary permissions to function:

    ```bash
    ROLE_PERMISSIONS=(
        compute.disks.get
        compute.disks.create
        compute.disks.createSnapshot
        compute.snapshots.get
        compute.snapshots.create
        compute.snapshots.useReadOnly
        compute.snapshots.delete
        compute.zones.get
    )

    gcloud iam roles create velero.server \
        --project $PROJECT_ID \
        --title "Velero Server" \
        --permissions "$(IFS=","; echo "${ROLE_PERMISSIONS[*]}")"

    gcloud projects add-iam-policy-binding $PROJECT_ID \
        --member serviceAccount:$SERVICE_ACCOUNT_EMAIL \
        --role projects/$PROJECT_ID/roles/velero.server

    gsutil iam ch serviceAccount:$SERVICE_ACCOUNT_EMAIL:objectAdmin gs://${BUCKET}
    ```

4. Create a service account key, specifying an output file (`credentials-velero`) in your local directory:

    ```bash
    gcloud iam service-accounts keys create credentials-velero \
        --iam-account $SERVICE_ACCOUNT_EMAIL
    ```

### Option 2: Set permissions with using Workload Identity (Optional)

If you are running Velero on a GKE cluster with workload identity enabled, you may want to bind Velero's Kubernetes service account to a GCP service account with the appropriate permissions instead of providing the key file during installation.

To do this, you must grant the GCP service account(the one you created in Step 3) the 'iam.serviceAccounts.signBlob' role. This is so that Velero's Kubernetes service account can create signed urls for the GCP bucket.

Next, add an IAM policy binding to grant Velero's Kubernetes service account access to your created GCP service account.

```bash
gcloud iam service-accounts add-iam-policy-binding \
    --role roles/iam.workloadIdentityUser \
    --member serviceAccount:[PROJECT_ID].svc.id.goog[velero/velero] \
    [GSA_NAME]@[PROJECT_ID].iam.gserviceaccount.com
```

For more information on configuring workload identity on GKE, look at the [official GCP documentation][24] for more details.

## Install and start Velero

[Download][4] Velero

Install Velero, including all prerequisites, into the cluster and start the deployment. This will create a namespace called `velero`, and place a deployment named `velero` in it.

**If using a Service Account**:

```bash
velero install \
    --provider gcp \
    --plugins velero/velero-plugin-for-gcp:v1.0.1 \
    --bucket $BUCKET \
    --secret-file ./credentials-velero
```

**If using Workload Identity**:

You must add a service account annotation to the Kubernetes service account so that it will know which GCP service account to use. You can do this during installation with `--sa-annotations`. Use the flag `--no-secret` so that Velero will know not to look for a key file. You must also add the GCP service account name in `--backup-location-config`.

```bash
velero install \
    --provider gcp \
    --plugins velero/velero-plugin-for-gcp:v1.0.1 \
    --no-secret \
    --sa-annotations iam.gke.io/gcp-service-account=[GSA_NAME]@[PROJECT_ID].iam.gserviceaccount.com \
    --backup-location-config serviceAccount=[GSA_NAME]@[PROJECT_ID].iam.gserviceaccount.com \
```

Additionally, you can specify `--use-restic` to enable restic support, and `--wait` to wait for the deployment to be ready.

(Optional) Specify [additional configurable parameters][7] for the `--backup-location-config` flag.

(Optional) Specify [additional configurable parameters][8] for the `--snapshot-location-config` flag.

(Optional) Specify [CPU and memory resource requests and limits][9] for the Velero/restic pods.

For more complex installation needs, use either the Helm chart, or add `--dry-run -o yaml` options for generating the YAML representation for the installation.

[1]: #Create-an-GCS-bucket
[2]: #Set-permissions-for-Velero
[3]: #Install-and-start-Velero
[4]: https://velero.io/docs/master/install-overview/#install-the-cli
[5]: https://cloud.google.com/sdk/docs/
[7]: backupstoragelocation.md
[8]: volumesnapshotlocation.md
[9]: https://velero.io/docs/master/install-requirement
[10]: ./examples
[11]: https://velero.io/docs/master/faq/
[15]: https://cloud.google.com/compute/docs/access/service-accounts
[22]: https://cloud.google.com/kubernetes-engine/docs/how-to/role-based-access-control#iam-rolebinding-bootstrap
[24]: https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity
