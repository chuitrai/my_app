// Jenkinsfile - Sử dụng HTTPS và Personal Access Token trên agent mặc định

pipeline {
    // Chạy trên bất kỳ agent nào có sẵn của Jenkins trên K8s.
    // Agent này sẽ là một môi trường Linux.
    agent any

    // Các biến môi trường để quản lý tập trung
    environment {
        // --- Cấu hình Docker & Image ---
        DOCKER_USERNAME       = 'chuitrai2901'
        BACKEND_IMAGE_NAME    = "${DOCKER_USERNAME}/my-go-backend"
        DOCKER_CREDENTIALS_ID = 'dock-cre'      // ID của credential Docker Hub

        // --- Cấu hình Git & Repo ---
        // Sử dụng URL HTTPS cho repo cấu hình
        CONFIG_REPO_URL_HTTPS = 'https://github.com/chuitrai/my_app_config.git' 
        CONFIG_REPO_DIR       = 'my_app_config_clone' 
        GIT_CREDENTIALS_ID    = 'github-pat'          // ID của credential chứa PAT
    }

    stages {
        // Giai đoạn 1: Checkout code (chạy trong container mặc định 'jnlp')
        stage('Checkout Source') {
            steps {
                echo 'Checking out application source code...'
                checkout scm 
            }
        }

        // Giai đoạn 2: Cài đặt công cụ (chạy trong container 'docker')
        stage('Install Dependencies') {
            steps {
                // Chuyển sang container 'docker' để thực thi
                container('docker') {
                    script {
                        echo "Installing required tools inside the 'docker' container..."
                        sh 'apk add --no-cache git sed'
                    }
                }
            }
        }

        // Giai đoạn 3: Đẩy Image lên Docker Hub
        stage('Push Image') {
            steps {
                container('docker') {
                    script {
                        // Sử dụng credential đã lưu để đăng nhập và push
                        docker.withRegistry("https://index.docker.io/v1/", DOCKER_CREDENTIALS_ID) {
                            docker.image("${BACKEND_IMAGE_NAME}:${env.NEW_IMAGE_TAG}").push()
                            echo "Successfully pushed ${BACKEND_IMAGE_NAME}:${env.NEW_IMAGE_TAG}"
                        }
                    }
                }
            }
        }

        // Giai đoạn 4: Cập nhật Repo Cấu hình
        stage('Update Configuration') {
            steps {
                container('docker') {
                    script {
                        echo "Updating config repo with new image tag: ${env.NEW_IMAGE_TAG}"
                        
                        // Sử dụng credential 'github-pat' (Loại: Secret text)
                        withCredentials([string(credentialsId: GIT_CREDENTIALS_ID, variable: 'GIT_TOKEN')]) {
                            
                            // Xóa thư mục clone cũ nếu tồn tại để đảm bảo sạch sẽ
                            sh "rm -rf ${CONFIG_REPO_DIR}"

                            // Clone repo config bằng URL HTTPS có chèn token để xác thực
                            sh "git clone https://${GIT_TOKEN}@github.com/chuitrai/my_app_config.git ${CONFIG_REPO_DIR}"
                            
                            // Di chuyển vào thư mục repo vừa clone
                            dir(CONFIG_REPO_DIR) {
                                // Cấu hình thông tin người commit (Jenkins Bot)
                                sh "git config user.email 'jenkins-bot@example.com'"
                                sh "git config user.name 'Jenkins Bot'"

                                // Dùng sed để tìm và thay thế dòng tag của backend
                                // Yêu cầu: trong values.yaml, dòng tag của backend phải có comment #backend-tag
                                sh "sed -i 's|^    tag: .*#backend-tag|    tag: ${env.NEW_IMAGE_TAG} #backend-tag|' values.yaml"
                                
                                // Commit và push thay đổi lên nhánh main
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
}