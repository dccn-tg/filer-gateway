def docker_ps_health_check(list) {
    for (int i = 0; i < list.size(); i++) {
        statusCode = sh(
            returnStatus: true,
            label: "checking if ${list[i]} in running state",
            script: "docker ps --format {{.Names}} --filter status=running --filter name=${list[i]} | grep ${list[i]}"
        )

        if ( statusCode != 0 ) {
            return statusCode
        }
    }
}

pipeline {
    agent any

    environment {
        DOCKER_IMAGE_TAG = "jenkins-${env.BUILD_NUMBER}"
    }
    
    options {
        disableConcurrentBuilds()
    }

    stages {
        // build containers.
        stage('Build') {
            when {
                expression {
                    return !params.PRODUCTION
                }
            }
            agent {
                label 'swarm-manager'
            }
            steps {
                sh 'docker-compose build --parallel'
            }
        }

        // build contains with tags ready for pushing to production registry.
        stage('Build (PRODUCTION)') {
            when {
                expression {
                    return params.PRODUCTION
                }
            }
            agent {
                label 'swarm-manager'
            }
            steps {
                withEnv([
                    "DOCKER_REGISTRY=${params.PRODUCTION_DOCKER_REGISTRY}",
                    "DOCKER_IMAGE_TAG=${params.PRODUCTION_GITHUB_TAG}"
                ]) {
                    sh 'docker-compose build --parallel'
                }
            }
        }

        // perform unit test, and push the containers to the test environment registry on success.
        stage('Unit tests') {
            when {
                expression {
                    return !params.PRODUCTION
                }
            }
            agent {
                label 'swarm-manager'
            }
            steps {
                echo 'Unit tests go here'
            }
            post {
                success {
                    script {
                        if (env.DOCKER_REGISTRY) {
                            sh (
                                label: "Pushing images to repository '${env.DOCKER_REGISTRY}'",
                                script: 'docker-compose push'
                            )
                        }
                    }
                }
            }
        }

        // deploying containers on the test environment.
        stage('Staging') {
             when {
                expression {
                    return !params.PRODUCTION
                }
            }
            agent {
                label 'swarm-manager'
            }
            steps {

                script {
                    def statusCode = sh(script: "bash ./stop.sh", returnStatus: true)
                    echo "statusCode: ${statusCode}"
                }
                
                sleep(10)

                withCredentials([
                    // Fill in runtime credentials
                    usernamePassword (
                        credentialsId: params.PLAYGROUND_DOCKER_REGISTRY_CREDENTIALS,
                        usernameVariable: 'DOCKER_REGISTRY_USER',
                        passwordVariable: 'DOCKER_REGISTRY_PASSWORD'
                    )
                ]) {
                    // Overwrite the env.sh file to be stored later as an artifact
                    script {
                        def statusCode = sh(script: "bash ./print_env.sh > ${WORKSPACE}/env.sh", returnStatus: true)
                        echo "statusCode: ${statusCode}"
                    }

                    // Use the same approach as for production
                    script {
                        def statusCode = sh(script: "bash ./start.sh", returnStatus: true)
                        echo "statusCode: ${statusCode}"
                    }
                }
            }
        }

        // check if containers are properly started on the test environment.
        stage('Health check') {
            when {
                expression {
                    return !params.PRODUCTION
                }
            }
            agent {
                label 'swarm-manager'
            }
            steps {
                // wait for 10 seconds
                sleep 10

                // check if containers are in running state
                docker_ps_health_check(['filer-gateway_api-server','filer-gateway_worker','filer-gateway_db'])
            }
        }

        // making release tag and push containers to production registry.
        stage('Tag and push (PRODUCTION)') {
            when {
                expression {
                    return params.PRODUCTION
                }
            }
            agent {
                label 'swarm-manager'
            }
            steps {
                echo "production: true"
                echo "production github tag: ${params.PRODUCTION_GITHUB_TAG}"

                // Handle Github tags
                withCredentials([
                    usernamePassword (
                        credentialsId: params.GITHUB_CREDENTIALS,
                        usernameVariable: 'GITHUB_USERNAME',
                        passwordVariable: 'GITHUB_PASSWORD'
                    )
                ]) {
                    // Remove local tag (if any)
                    script {
                        def statusCode = sh(script: "git tag --list | grep ${params.PRODUCTION_GITHUB_TAG}", returnStatus: true)
                        if(statusCode == 0) {
                            sh "git tag -d ${params.PRODUCTION_GITHUB_TAG}"
                            echo "Removed existing local tag ${params.PRODUCTION_GITHUB_TAG}"
                        }
                    }
                    
                    // Create local tag
                    sh "git tag -a ${params.PRODUCTION_GITHUB_TAG} -m 'jenkins'"
                    echo "Created local tag ${params.PRODUCTION_GITHUB_TAG}"

                    // Remove remote tag (if any)
                    script {
                        def result = sh(script: "git ls-remote https://${GITHUB_USERNAME}:${GITHUB_PASSWORD}@github.com/dccn-tg/filer-gateway.git refs/tags/${params.PRODUCTION_GITHUB_TAG}", returnStdout: true).trim()
                        if (result != "") {
                            sh "git push --delete https://${GITHUB_USERNAME}:${GITHUB_PASSWORD}@github.com/dccn-tg/filer-gateway.git ${params.PRODUCTION_GITHUB_TAG}"
                            echo "Removed existing remote tag ${params.PRODUCTION_GITHUB_TAG}"
                        }
                    }

                    // Create remote tag
                    sh "git push https://${GITHUB_USERNAME}:${GITHUB_PASSWORD}@github.com/dccn-tg/filer-gateway.git ${params.PRODUCTION_GITHUB_TAG}"
                    echo "Created remote tag ${params.PRODUCTION_GITHUB_TAG}"
                }

                // Override the env variables and 
                // push the Docker images to the production Docker registry
                withEnv([
                    "DOCKER_REGISTRY=${params.PRODUCTION_DOCKER_REGISTRY}",
                    "DOCKER_IMAGE_TAG=${params.PRODUCTION_GITHUB_TAG}"
                ]) {
                    withCredentials([
                        usernamePassword (
                            credentialsId: params.PRODUCTION_DOCKER_REGISTRY_CREDENTIALS,
                            usernameVariable: 'DOCKER_USERNAME',
                            passwordVariable: 'DOCKER_PASSWORD'
                        )
                    ]) {
                        sh "docker login -u ${DOCKER_USERNAME} -p ${DOCKER_PASSWORD} ${params.PRODUCTION_DOCKER_REGISTRY}"
                        sh 'docker-compose push'
                        echo "Pushed images to ${DOCKER_REGISTRY}"
                    }
                } 
            }
        }
    }

    post {
        success {
            script {
                // regenerate env.sh; but strip out the username/password
                def statusCode = sh(returnStatus:true, script: "bash ./print_env.sh | sed 's/^DOCKER_REGISTRY_USER=.*/DOCKER_REGISTRY_USER=/' | sed 's/^DOCKER_REGISTRY_PASSWORD=.*/DOCKER_REGISTRY_PASSWORD=/' > env.sh")
                if ( statusCode != 0 ) {
                    echo "unable to generate env.sh file, check it manually."
                }
            }
            archiveArtifacts "docker-compose.yml, start.sh, stop.sh, env.sh"
        }
        always {
            echo 'cleaning'
            sh 'docker system prune -f'
        }
    }
}
