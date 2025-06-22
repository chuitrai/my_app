// Jenkinsfile - Sử dụng Docker Agent để đảm bảo môi trường Linux

pipeline {
    // ---- THAY ĐỔI QUAN TRỌNG NHẤT ----
    // Yêu cầu Jenkins chạy pipeline này bên trong một container Docker
    agent {
        docker {
            // Sử dụng một image có sẵn Docker client và các công cụ cơ bản
            image 'docker:20.10.16' 
            // Cung cấp các cờ bổ sung để agent có thể giao tiếp với Docker daemon của host
            args '-v /var/run/docker.sock:/var/run/docker.sock' 
        }
    }

    // Các biến môi trường (giữ nguyên)
    environment {
        DOCKER_USERNAME       = 'chuitrai2901'
        BACKEND_IMAGE_NAME    = "${DOCKER_USERNAME}/my-go-backend"
        CONFIG_REPO_URL_HTTPS = 'https://github.com/chuitrai/my_app_config.git'
        CONFIG_REPO_DIR       = 'my_app_config_clone'
        DOCKER_CREDENTIALS_ID = 'dock-cre'
        GIT_CREDENTIALS_ID    = 'github-pat'
    }

    stages {
        // Giai đoạn 1: Lấy source code ứng dụng
        stage('Checkout Source') {
            steps {
                echo 'Checking out application source code...'
                checkout scm 
            }
        }

        // --- GIAI ĐOẠN MỚI: CÀI ĐẶT CÔNG CỤ CẦN THIẾT ---
        stage('Install Dependencies') {
            steps {
                script {
                    echo "Installing required tools inside the agent container..."
                    // Image docker:20.10.16 dùng Alpine Linux, sử dụng apk để cài đặt
                    sh 'apk add --no-cache git openssh-client sed'
                }
            }
        }

        // Giai đoạn 2: Build Docker Image
        stage('Build Image') {
            steps {
                script {
                    env.NEW_IMAGE_TAG = "v1.0.${env.BUILD_NUMBER}"
                    echo "Building image with tag: ${env.NEW_IMAGE_TAG}"
                    
                    // Lệnh docker.build bây giờ sẽ được thực thi bởi Docker client
                    // bên trong container agent, và nó sẽ ra lệnh cho Docker daemon
                    // bên ngoài (của máy host) để thực hiện việc build.
                    docker.build("${BACKEND_IMAGE_NAME}:${env.NEW_IMAGE_TAG}", "./backend")
                }
            }
        }

        // Các stage 'Push Image' và 'Update Configuration' giữ nguyên vì chúng
        // đã sử dụng các lệnh sh, và bây giờ sẽ được chạy trong môi trường Linux.
        stage('Push Image') {
            steps {
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