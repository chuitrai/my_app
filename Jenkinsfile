// Jenkinsfile - Phiên bản nâng cấp cho Backend & Frontend, sử dụng Kaniko

pipeline {
    // Không định nghĩa agent ở cấp cao nhất, sẽ định nghĩa cho từng stage
    agent none

    environment {
        // --- Repo và Credentials ---
        DOCKER_REGISTRY_URL   = 'https://index.docker.io/v1/'
        DOCKER_CREDENTIALS_ID = 'dock-cre'      // Jenkins credential ID cho Docker Hub
        CONFIG_REPO_URL       = 'https://github.com/chuitrai/my_app_config.git'
        CONFIG_REPO_DIR       = 'my_app_config_clone'
        GIT_CREDENTIALS_ID    = 'git-pat'       // Jenkins credential ID cho GitHub PAT

        // --- Tên Image ---
        DOCKER_USERNAME       = 'chuitrai2901'
        BACKEND_IMAGE_REPO    = "${DOCKER_USERNAME}/my-go-backend"
        FRONTEND_IMAGE_REPO   = "${DOCKER_USERNAME}/my-react-frontend"
    }

    stages {
        // ======================================================================
        // STAGE 1: Chuẩn bị và xác định phiên bản
        // ======================================================================
        stage('1. Preparation & Versioning') {
            agent any // Chạy trên một agent bất kỳ để thực hiện các tác vụ nhẹ
            steps {
                script {
                    echo "=========================================="
                    echo "Triggered by: ${currentBuild.fullDisplayName}"
                    // Điều kiện: Chỉ chạy khi được trigger bởi một tag
                    if (!env.TAG_NAME) {
                        error "BUILD ABORTED: This pipeline is designed to run only on git tags."
                    }
                    echo "VERSION TO BUILD: ${env.TAG_NAME}"
                    echo "=========================================="

                    // Lưu tag vào workspace để các stage sau có thể dùng
                    writeFile file: 'version.txt', text: env.TAG_NAME
                }
            }
        }

        // ======================================================================
        // STAGE 2: Build Images Song Song
        // ======================================================================
        stage('2. Build Application Images') {
            // Chạy hai stage con này song song
            parallel {
                // --- Build Backend ---
                stage('Build Backend') {
                    // Sử dụng agent Kaniko, không cần Docker-in-Docker
                    agent {
                        kubernetes {
                            cloud 'kubernetes'
                            label 'kaniko-agent' // Label cho pod template, cần định nghĩa trong Jenkins config
                            yamlFile 'kaniko-pod-template.yaml' // Sử dụng file template cho sạch sẽ
                        }
                    }
                    steps {
                        container(name: 'kaniko') {
                            script {
                                def imageTag = readFile('version.txt').trim()
                                def finalImageName = "${env.BACKEND_IMAGE_REPO}:${imageTag}"
                                echo "Building and pushing Backend image: ${finalImageName}"

                                // Checkout SCM bên trong workspace của container
                                checkout scm

                                // Lệnh thực thi Kaniko
                                sh """
                                /kaniko/executor --context=dir://\$(pwd)/backend \
                                                 --dockerfile=\`pwd\`/backend/Dockerfile \
                                                 --destination=${finalImageName} \
                                                 --build-arg version=${imageTag}
                                """
                            }
                        }
                    }
                }

                // --- Build Frontend ---
                stage('Build Frontend') {
                    agent {
                        kubernetes {
                            cloud 'kubernetes'
                            label 'kaniko-agent'
                            yamlFile 'kaniko-pod-template.yaml'
                        }
                    }
                    steps {
                        container(name: 'kaniko') {
                            script {
                                def imageTag = readFile('version.txt').trim()
                                def finalImageName = "${env.FRONTEND_IMAGE_REPO}:${imageTag}"
                                echo "Building and pushing Frontend image: ${finalImageName}"

                                // Checkout SCM
                                checkout scm

                                // Lệnh thực thi Kaniko
                                sh """
                                /kaniko/executor --context=dir://\$(pwd)/frontend \
                                                 --dockerfile=\`pwd\`/frontend/Dockerfile \
                                                 --destination=${finalImageName} \
                                                 --build-arg version=${imageTag}
                                """
                            }
                        }
                    }
                }
            }
        }

        // ======================================================================
        // STAGE 3: Cập nhật repo cấu hình (GitOps)
        // ======================================================================
        stage('3. Update Deployment Configuration') {
            agent any // Agent cần có git và sed
            steps {
                script {
                    def releaseTag = readFile('version.txt').trim()
                    echo "Updating config repo to version: ${releaseTag}"

                    // Cần cài đặt git trên agent này nếu chưa có
                    sh 'git --version'
                    sh 'sed --version'

                    withCredentials([string(credentialsId: GIT_CREDENTIALS_ID, variable: 'GIT_TOKEN')]) {
                        sh "rm -rf ${CONFIG_REPO_DIR}"
                        // Clone repo config
                        sh "git clone https://${GIT_TOKEN}@github.com/chuitrai/my_app_config.git ${CONFIG_REPO_DIR}"

                        dir(CONFIG_REPO_DIR) {
                            sh "git config user.email 'jenkins-ci-bot@noreply.com'"
                            sh "git config user.name 'Jenkins CI Bot'"

                            // Cập nhật tag cho backend
                            // Sử dụng #comment làm anchor để sed không bị nhầm lẫn
                            sh "sed -i 's|^    tag: .*#backend-tag|    tag: \"${releaseTag}\" #backend-tag|' my-go-app/values.yaml"

                            // Cập nhật tag cho frontend
                            sh "sed -i 's|^    tag: .*#frontend-tag|    tag: \"${releaseTag}\" #frontend-tag|' my-go-app/values.yaml"

                            // Cập nhật appVersion trong Chart.yaml
                            sh "sed -i 's|^appVersion: .*|appVersion: \"${releaseTag}\"|' my-go-app/Chart.yaml"
                            
                            def changes = sh(script: "git status --porcelain", returnStdout: true).trim()
                            if (changes) {
                                echo "Config changes detected. Committing and pushing..."
                                sh "git add ."
                                sh "git commit -m 'CI: Release backend and frontend to version ${releaseTag}'"
                                sh "git push origin main"
                                echo "Successfully pushed configuration update."
                            } else {
                                echo "No changes detected in config repo. Skipping commit and push."
                            }
                        }
                    }
                }
            }
        }
    }
    post {
        always {
            // Dọn dẹp workspace
            cleanWs()
        }
    }
}