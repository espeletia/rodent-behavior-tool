load('ext://restart_process', 'docker_build_with_restart')


load_dynamic('./ci/minio/minio.Tiltfile')
load_dynamic('./ci/nats/nats.Tiltfile')

k8s_yaml("ci/kube_ratt.yaml")

docker_build_with_restart('ratt',
      '.',
      dockerfile='./ratt/ci/Dockerfile',
      entrypoint='python main.py',
      only=[
          './ratt',
      ],
      live_update=[
          sync('./ratt/' , '/app/'),
      ])

k8s_resource('ratt', port_forwards=["0.0.0.0:8080:8081"], labels=["AI"])

