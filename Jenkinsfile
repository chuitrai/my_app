// file: my_app/Jenkinsfile

pipeline {
    // Chạy trên bất kỳ agent nào có sẵn
    agent any 

    // Các biến môi trường để dễ quản lý
    environment {
        DOCKER_USERNAME       = 'chuitrai2901' // <-- Thay bằng Docker ID của bạn
        BACKEND_IMAGE_NAME    = "${DOCKER_USERNAME}/my-go-backend"
        CONFIG_REPO_URL_SSH   = 'git@github.com:chuitrai/my_app_config.git' // Dùng SSH để push
        CONFIG_REPO_DIR       = 'my_app_config'
        // Sử dụng ID của credential bạn đã tạo
        DOCKER_CREDENTIALS_ID = 'dock-cre' 
        // ID của credential chứa SSH key cho repo config
        GIT_CREDENTIALS_ID    = 'github-pat' // Chúng ta sẽ tạo credential này
    }

    stages {
        // Giai đoạn 1: Lấy source code
        stage('Checkout Source') {
            steps {
                echo 'Checking out application source code...'
                // Tự động checkout nhánh đã trigger pipeline
                checkout scm 
            }   
        }

        // Giai đoạn 2: Build Docker Image
        stage('Build Image') {
            steps {
                script {
                    // Tạo một tag duy nhất dựa trên số lần build
                    def newTag = "v1.0.${env.BUILD_NUMBER}"
                    echo "Building image with tag: ${newTag}"
                    // Build image từ Dockerfile trong thư mục ./backend
                    docker.build("${BACKEND_IMAGE_NAME}:${newTag}", "./backend")
                }
            }
        }

        // Giai đoạn 3: Đẩy Image lên Docker Hub
        stage('Push Image') {
            steps {
                script {
                    def newTag = "v1.0.${env.BUILD_NUMBER}"
                    // Sử dụng credential đã lưu để đăng nhập và push
                    docker.withRegistry("https://index.docker.io/v1/", DOCKER_CREDENTIALS_ID) {
                        docker.image("${BACKEND_IMAGE_NAME}:${newTag}").push()
                        echo "Successfully pushed ${BACKEND_IMAGE_NAME}:${newTag}"
                    }
                }
            }
        }

        // Giai đoạn 4: Cập nhật Repo Cấu hình
        stage('Update Configuration') {
            steps {
                script {
                    def newTag = "v1.0.${env.BUILD_NUMBER}"
                    echo "Updating config repo with new image tag: ${newTag}"
                    // Sử dụng SSH key để checkout và push vào repo config
                    withCredentials([sshUserPrivateKey(credentialsId: GIT_CREDENTIALS_ID, keyFileVariable: 'SSH_KEY')]) {
                        // Cần cài đặt ssh-agent trên Jenkins agent
                        sh '''
                            # Bắt đầu ssh-agent và thêm key vào
                            eval $(ssh-agent -s)
                            ssh-add $SSH_KEY
                            
                            # Tắt kiểm tra host key nghiêm ngặt
                            mkdir -p ~/.ssh
                            echo "Host github.com\\n\\tStrictHostKeyChecking no\\n" >> ~/.ssh/config

                            # Clone, sửa, commit, và push
                            git clone ${CONFIG_REPO_URL_SSH} ${CONFIG_REPO_DIR}
                            cd ${CONFIG_REPO_DIR}
                            
                            # Dùng yq (cần cài trên agent) hoặc sed để sửa file
                            # Ví dụ dùng sed
                            sed -i "s|tag:.*#backend|tag: ${newTag} #backend|" values.yaml

                            git config user.email "jenkins-bot@example.com"
                            git config user.name "Jenkins Bot"
                            git add values.yaml
                            git commit -m "CI: Bump backend image to ${newTag}"
                            git push origin main
                        '''
                    }
                }
            }
        }
    }
}