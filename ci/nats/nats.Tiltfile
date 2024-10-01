load('ext://helm_resource', 'helm_resource')
load('ext://helm_resource', 'helm_repo')
load('ext://namespace', 'namespace_create')

helm_repo('nats-helm', 'https://nats-io.github.io/k8s/helm/charts', labels=["helm"])

helm_resource('nats',
            chart='nats-helm/nats',
            release_name='nats',
            resource_deps=['nats-helm'],
            labels=["DB"],
            flags=[
            '--set', 'config.jetstream.enabled=true',
            '--set', 'natsbox.enabled=false',
            '--set', 'cluster.enabled=false'
            ]
)

k8s_resource('nats', port_forwards=["0.0.0.0:4222:4222", "0.0.0.0:8222:8222"], labels=["DB"])

