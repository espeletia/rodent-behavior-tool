load('ext://helm_resource', 'helm_resource')
load('ext://helm_resource', 'helm_repo')

helm_repo('minio-helm', 'https://charts.bitnami.com/bitnami', labels=["helm"])

helm_resource('minio',
            chart='minio-helm/minio',
            release_name='minio',
            resource_deps=['minio-helm'],
            labels=["S3"],
            flags=[
            '--set', 'statefulset.replicaCount=1',
            '--set', 'auth.rootUser=minio123',
            '--set', 'auth.rootPassword=minio123',
            '--set', 'defaultBuckets=test:public'
            ]
)

k8s_resource('minio', port_forwards=["0.0.0.0:9000:9000", "0.0.0.0:9001:9001"], labels=["S3"])
