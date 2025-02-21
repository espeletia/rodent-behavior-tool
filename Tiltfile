load('ext://restart_process', 'docker_build_with_restart')


load_dynamic('./ci/minio/minio.Tiltfile')
load_dynamic('./ci/nats/nats.Tiltfile')
load_dynamic('./ci/postgres/postgres.Tiltfile')

k8s_yaml("ci/kube_ratt.yaml")
k8s_yaml("ci/kube_echoes.yaml")
k8s_yaml("ci/kube_tusk.yaml")
k8s_yaml("ci/kube_valentine.yaml")

def app(name, special_dockerfile = False, migrations = True):
    local_resource(
      'compile-%s' % name,
      'cd %s && ./ci/build.sh' % name,
      deps=[
      './%s/' % name,
      './ghiaccio/'
      ],
      ignore=[
      'tilt_modules',
      'Tiltfile',
      './iris/proto/gen',
      '%s/graph/schema.graphqls' % name,
      '%s/build' % name,
      '%s/dep' % name,
      '%s/ci/docker-compose.yaml' % name,
      '%s/swagger.yaml' % name,
      '%s/internal/handlers/swagger.yaml' % name,
      '%s/internal/handlers/generated.go' % name,
      '%s/migrations/*' % name,
      '%s/cmd/migrations/*' % name,
      '%s/cmd/dataInit/*' % name,
      '%s/**/testdata' % name
      ],
      resource_deps=['nats', 'minio', 'postgresql']
    )
    if migrations:
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

    if special_dockerfile:
        docker_build_with_restart('%s-migrations' % name,
            '.',
            dockerfile='./%s/ci/Dockerfile' % name,
            entrypoint='/app/run_migrations',
            only=[
                './%s/build' % name,
                './%s/ci' % name,
                './%s/configurations' % name,
                './%s/certs' % name,
                './%s/migrations' % name,
                './%s/videos' % name
            ],
            live_update=[
                sync('./%s/build' % name , '/app'),
                sync('./%s/configurations' % name , '/app/configurations')
            ],
            build_args={"app": name})

        docker_build_with_restart('%s' % name,
            '.',
            dockerfile='./%s/ci/Dockerfile' % name,
            entrypoint='/app/start_server',
            only=[
                './%s/build' % name,
                './%s/dep' % name,
                './%s/files' % name,
                './%s/ci' % name,
                './%s/configurations' % name,
                './%s/certs' % name,
                './%s/migrations' % name,
                './%s/videos' % name
            ],
            live_update=[
                sync('./%s/build' % name , '/app'),
                sync('./%s/configurations' % name , '/app/configurations')
            ],
            build_args={"app": name})
    else:
        if migrations:
            docker_build_with_restart('%s-migrations' % name,
                '.',
                dockerfile='./ci/Dockerfile',
                entrypoint='/app/run_migrations',
                only=[
                    './%s/build' % name,
                    './%s/ci' % name,
                    './%s/configurations' % name,
                    './%s/certs' % name,
                    './%s/migrations' % name,
                    './%s/videos' % name
                ],
                live_update=[
                    sync('./%s/build' % name , '/app'),
                    sync('./%s/configurations' % name , '/app/configurations')
                ],
                build_args={"app": name})

        docker_build_with_restart('%s' % name,
            '.',
            dockerfile='ci/Dockerfile',
            entrypoint='/app/start_server',
            only=[
                './%s/build' % name,
                './%s/configurations' % name,
                './%s/certs' % name,
                './%s/migrations' % name,
                './%s/videos' % name
            ],
            live_update=[
                sync('./%s/build' % name , '/app'),
                sync('./%s/configurations' % name , '/app/configurations')
            ],
            build_args={"app": name})



app('echoes', True)
app('tusk')
app('valentine', False, False)

local_resource(
    'reg-tusk',
    'cd tusk && go generate cmd/main.go',
    deps=[
    './tusk/graph/',
    ],
    ignore=[
    './tusk/graph/**/*.go*',
    './tusk/graph/generated',
    './tusk/graph/model',
    './tusk/migrations/migrations.xml',
    ],
    resource_deps=['postgresql']
)

local_resource(
    'reg-valentine',
    'cd valentine && go generate cmd/main.go',
    deps=[
    './valentine/view/',
    ],
    ignore=[
    './valentine/view/*.go',
    './valentine/view/**/*.go*',
    ],
    resource_deps=['postgresql']
)

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

k8s_resource('ratt', port_forwards=["0.0.0.0:8083:8081"], labels=["AI"], resource_deps=['nats','minio'])
k8s_resource('echoes', labels=["ENCODING"], resource_deps=['nats', 'minio'])
k8s_resource('tusk', labels=["BE"], port_forwards=["0.0.0.0:8081:8080"], resource_deps=['nats', 'minio'])
k8s_resource('valentine', labels=["FE"], port_forwards=["0.0.0.0:3000:3000"], resource_deps=['tusk'])
k8s_resource('tusk-queue', labels=["BE"], resource_deps=['nats', 'minio', 'tusk'])

