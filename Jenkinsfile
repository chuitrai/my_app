// Jenkinsfile - Đã sửa lỗi cú pháp Declarative Pipeline

pipeline {
    agent any

    environment {
        DOCKER_USERNAME       = 'chuitrai2901'
        BACKEND_IMAGE_NAME    = "${DOCKER_USERNAME}/my-go-backend"
        CONFIG_REPO_URL_HTTPS = 'https://github.com/chuitrai/my_app_config.git'
        CONFIG_REPO_DIR       = 'my_app_config_clone'
        DOCKER_CREDENTIALS_ID = 'dock-cre'
        GIT_CREDENTIALS_ID    = 'github-pat'
    }

    stages {
        stage('Checkout Source') {
            steps {
                echo 'Checking out application source code...'
                checkout scm 
            }
        }

        stage('Build Image') {
            steps {
                // ---- BỌC TOÀN BỘ LOGIC VÀO KHỐI SCRIPT ----
                script {
                    env.NEW_IMAGE_TAG = "v1.0.${env.BUILD_NUMBER}"
                    echo "Building image with tag: ${env.NEW_IMAGE_TAG}"
                    docker.build("${BACKEND_IMAGE_NAME}:${env.NEW_IMAGE_TAG}", "./backend")
                }
            }
        }

        stage('Push Image') {
            steps {
                // ---- BỌC LỆNH PHỨC TẠP VÀO KHỐI SCRIPT ----
                script {
                    docker.withRegistry("https://index.docker.io/v1/", DOCKER_CREDENTIALS_ID) {
                        docker.image("${BACKEND_IMAGE_NAME}:${env.NEW_IMAGE_TAG}").push()
                        echo "Successfully pushed ${BACKEND_IMAGE_NAME}:${env.NEW_IMAGE_TAG}"
                    }
                }
            }
        }

        stage('Update Configuration') {
            steps {
                // ---- BỌC TOÀN BỘ LOGIC VÀO KHỐI SCRIPT ----
                script {
                    echo "Updating config repo with new image tag: ${env.NEW_IMAGE_TAG}"
                    withCredentials([string(credentialsId: GIT_CREDENTIALS_ID, variable: 'GIT_TOKEN')]) {
                        sh "rm -rf ${CONFIG_REPO_DIR}"
                        sh "git clone https://${GIT_TOKEN}@github.com/chuitrai/my_app_config.git ${CONFIG_REPO_DIR}"
                        
                        dir(CONFIG_REPO_DIR) {
                            sh "git config user.email 'jenkins-bot@example.com'"
                            sh "git config user.name 'Jenkins Bot'"
                            sh "sed -i 's|^    tag: .*#backend-tag|    tag: ${env.NEW_IMAGE_TAG} #backend-tag|' values.yaml"
                            sh "git add values.yaml"
                            sh "git commit -m 'CI: Bump backend image to ${env.NEW_IMAGE_TAG}'"
                            sh "git push origin main"
                            echo "Successfully pushed configuration update."
                        }
                    }
                }
            }
        }
    }
}