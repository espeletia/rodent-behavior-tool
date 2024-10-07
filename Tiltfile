load('ext://restart_process', 'docker_build_with_restart')


load_dynamic('./ci/minio/minio.Tiltfile')
load_dynamic('./ci/nats/nats.Tiltfile')

k8s_yaml("ci/kube_ratt.yaml")
k8s_yaml("ci/kube_echoes.yaml")

local_resource(
      'compile echoes',
      'cd echoes && bash ./ci/build.sh',
      deps=[
      './echoes/',
      './ghiaccio/'
      ],
      ignore=[
      'tilt_modules',
      'Tiltfile',
      'graph/schema.graphqls',
      'echoes/build',
      'dep',
      'ci/docker-compose.yaml',
      'swagger.yaml',
      '**/testdata'
      ],
      labels=["compile"],
  )
  
docker_build_with_restart('echoes',
    '.',
    dockerfile='./echoes/ci/Dockerfile',
    entrypoint='/app/start_server',
    only=[
        './echoes/build',
        './echoes/configurations',
        './echoes/certs',
        './echoes/videos',
        './echoes/migrations',
        './echoes/cmd/migrations'
    ],
    live_update=[
        sync('./configurations', '/app/configurations'),
        sync('./build', '/app')
    ]
)




name = "echoes"

local_resource(
        'compile-%s-migrations' % name,
        'cd %s && ./ci/build_migrations.sh' % name,
        deps=[
         '%s/cmd/migrations/' % name,
         '%s/cmd/dataInit/' % name,
         './%s/migrations' % name,
        ],
        resource_deps=['nats', 'minio', 'postgresql']
    )


docker_build_with_restart('echoes-migrations',
           '.',
           dockerfile='echoes/ci/Dockerfile',
           entrypoint='/app/run_migrations',
           only=[
               './%s/build' % name,
               './%s/ci' % name,
               './%s/configurations' % name,
               './%s/certs' % name,
               './echoes/videos',
               './%s/migrations' % name
           ],
           live_update=[
               sync('./%s/build' % name , '/app'),
               sync('./%s/configurations' % name , '/app/configurations')
           ],
           build_args={"app": name})



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

k8s_resource('ratt', port_forwards=["0.0.0.0:8080:8081"], labels=["AI"], resource_deps=['nats','minio'])
k8s_resource('echoes', labels=["ENCODING"])

