// Jenkinsfile - Phiên bản cuối cùng, tối ưu và sửa lỗi

pipeline {
    // Không định nghĩa agent ở cấp cao nhất, sẽ định nghĩa cho từng stage
    agent none

    environment {
        // --- Repo và Credentials ---
        DOCKER_REGISTRY_URL   = 'https://index.docker.io/v1/'
        DOCKER_CREDENTIALS_ID = 'dock-cre'
        CONFIG_REPO_URL       = 'https://github.com/chuitrai/my_app_config.git'
        CONFIG_REPO_DIR       = 'my_app_config_clone'
        GIT_CREDENTIALS_ID    = 'git-pat'

        // --- Tên Image ---
        DOCKER_USERNAME       = 'chuitrai2901'
        BACKEND_IMAGE_REPO    = "${DOCKER_USERNAME}/my-go-backend"
        FRONTEND_IMAGE_REPO   = "${DOCKER_USERNAME}/my-react-frontend"
    }

    stages {
        // ======================================================================
        // STAGE 1: Checkout, xác định phiên bản và lưu trữ
        // ======================================================================
        stage('1. Checkout & Versioning') {
            agent any // Chạy trên một agent bất kỳ
            steps {
                script {
                    // --- Checkout code một lần duy nhất ---
                    echo "Checking out source code..."
                    checkout scm

                    // --- Xác định phiên bản ---
                    echo "=========================================="
                    echo "Triggered by: ${currentBuild.fullDisplayName}"
                    if (!env.TAG_NAME) {
                        error "BUILD ABORTED: This pipeline is designed to run only on git tags."
                    }
                    def imageTag = env.TAG_NAME
                    echo "VERSION TO BUILD: ${imageTag}"
                    echo "=========================================="

                    // --- Lưu trữ (stash) workspace và version cho các stage sau ---
                    // 'includes' giúp stash nhẹ hơn, chỉ chứa những gì cần thiết.
                    stash name: 'source', includes: 'backend/**, frontend/**, kaniko-pod-template.yaml'
                    writeFile file: 'version.txt', text: imageTag
                    stash name: 'version', includes: 'version.txt'
                }
            }
        }

        // ======================================================================
        // STAGE 2: Build Images Song Song
        // ======================================================================
        stage('2. Build Application Images') {
            parallel {
                // --- Build Backend ---
                stage('Build Backend') {
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
                                // Lấy lại source code và file version
                                unstash 'source'
                                unstash 'version'

                                def imageTag = readFile('version.txt').trim()
                                def finalImageName = "${env.BACKEND_IMAGE_REPO}:${imageTag}"
                                echo "Building and pushing Backend image: ${finalImageName}"

                                // Lệnh Kaniko đã được sửa cú pháp
                                sh (
                                    script: '/kaniko/executor ' +
                                            '--context=dir://$(pwd)/backend ' +
                                            '--dockerfile=$(pwd)/backend/Dockerfile ' +
                                            '--destination=' + finalImageName + ' ' +
                                            '--destination=' + "${env.BACKEND_IMAGE_REPO}:latest" + ' ' + // Push thêm tag latest
                                            '--build-arg version=' + imageTag
                                )
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
                                // Lấy lại source code và file version
                                unstash 'source'
                                unstash 'version'

                                def imageTag = readFile('version.txt').trim()
                                def finalImageName = "${env.FRONTEND_IMAGE_REPO}:${imageTag}"
                                echo "Building and pushing Frontend image: ${finalImageName}"

                                // Lệnh Kaniko đã được sửa cú pháp
                                sh (
                                    script: '/kaniko/executor ' +
                                            '--context=dir://$(pwd)/frontend ' +
                                            '--dockerfile=$(pwd)/frontend/Dockerfile ' +
                                            '--destination=' + finalImageName + ' ' +
                                            '--destination=' + "${env.FRONTEND_IMAGE_REPO}:latest" + ' ' + // Push thêm tag latest
                                            '--build-arg version=' + imageTag
                                )
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
                    // Lấy lại file version
                    unstash 'version'
                    def releaseTag = readFile('version.txt').trim()
                    echo "Updating config repo to version: ${releaseTag}"

                    // Cần cài đặt git trên agent này nếu chưa có
                    sh 'apk add --no-cache git'

                    withCredentials([string(credentialsId: GIT_CREDENTIALS_ID, variable: 'GIT_TOKEN')]) {
                        sh "rm -rf ${CONFIG_REPO_DIR}"
                        sh "git clone https://x-access-token:${GIT_TOKEN}@github.com/chuitrai/my_app_config.git ${CONFIG_REPO_DIR}"

                        dir(CONFIG_REPO_DIR) {
                            sh "git config user.email 'jenkins-ci-bot@noreply.com'"
                            sh "git config user.name 'Jenkins CI Bot'"

                            // Lệnh sed đã được làm cho an toàn hơn
                            sh "sed -i 's|^    tag:.*#backend-tag|    tag: \"${releaseTag}\" #backend-tag|' my-go-app/values.yaml"
                            sh "sed -i 's|^    tag:.*#frontend-tag|    tag: \"${releaseTag}\" #frontend-tag|' my-go-app/values.yaml"
                            sh "sed -i 's|^appVersion:.*|appVersion: \"${releaseTag}\"|' my-go-app/Chart.yaml"
                            
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
    
    // Khối post đã được sửa lỗi
    post {
        always {
            // Cần một agent để thực hiện việc dọn dẹp
            agent any
            steps {
                echo "Cleaning up workspace."
                cleanWs()
            }
        }
    }
}