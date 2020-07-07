pipeline {
    agent any

    environment {
        DOCKER_IMAGE_TAG = "jenkins-${env.BUILD_NUMBER}"
        STACK_NAME = "filer-gateway"
    }
    
    options {
        disableConcurrentBuilds()
    }

    stages {
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

        // The following health check requires default network to be set "attachable". 
        // stage('Health check') {
        //      when {
        //         expression {
        //             return !params.PRODUCTION
        //         }
        //     }
        //     agent {
        //         docker {
        //             image 'jwilder/dockerize'
        //             args '--network filer-gateway_default'
        //         }
        //     }
        //     steps {
        //         sh (
        //             label: 'Waiting for services to become available',
        //             script: 'dockerize \
        //                 -timeout 120s \
        //                 -wait http://filer-gateway:8080'
        //         )
        //     }
        // }
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

                // check whether the docker service reports running
                sh (
                    label: 'checking if filer-gateway_api-server in running state',
                    script: "docker service ps filer-gateway_api-server --format '{{.CurrentState}}' | grep '^Running'"
                )
            }
        }


        // stage('Integration test') {
        //     steps {
        //         echo 'hi'
        //     }
        // }

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
                        def result = sh(script: "git ls-remote https://${GITHUB_USERNAME}:${GITHUB_PASSWORD}@github.com/Donders-Institute/filer-gateway.git refs/tags/${params.PRODUCTION_GITHUB_TAG}", returnStdout: true).trim()
                        if (result != "") {
                            sh "git push --delete https://${GITHUB_USERNAME}:${GITHUB_PASSWORD}@github.com/Donders-Institute/filer-gateway.git ${params.PRODUCTION_GITHUB_TAG}"
                            echo "Removed existing remote tag ${params.PRODUCTION_GITHUB_TAG}"
                        }
                    }

                    // Create remote tag
                    sh "git push https://${GITHUB_USERNAME}:${GITHUB_PASSWORD}@github.com/Donders-Institute/filer-gateway.git ${params.PRODUCTION_GITHUB_TAG}"
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
            archiveArtifacts "docker-compose.yml, docker-compose.swarm.yml, env.sh"
        }
        always {
            echo 'cleaning'
            sh 'docker system prune -f'
        }
    }
}
