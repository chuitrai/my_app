// Jenkinsfile - Thêm các dòng echo để tăng khả năng quan sát

pipeline {
    agent {
        kubernetes {
            yaml """
            apiVersion: v1
            kind: Pod
            spec:
              containers:
              - name: jnlp
                image: jenkins/inbound-agent:latest
                args: ['\$(JENKINS_SECRET)', '\$(JENKINS_NAME)']
                workingDir: /home/jenkins/agent
              - name: docker
                image: docker:20.10.16
                command: ['sleep']
                args: ['infinity']
                volumeMounts:
                - name: docker-socket
                  mountPath: /var/run/docker.sock
              volumes:
              - name: docker-socket
                hostPath:
                  path: /var/run/docker.sock
            """
            label 'k8s-agent-with-docker'
        }
    }

    environment {
        // ... (Giữ nguyên các biến môi trường của bạn) ...
        DOCKER_USERNAME       = 'chuitrai2901'
        BACKEND_IMAGE_NAME    = "${DOCKER_USERNAME}/my-go-backend"
        CONFIG_REPO_URL_HTTPS = 'https://github.com/chuitrai/my_app_config.git'
        CONFIG_REPO_DIR       = 'my_app_config_clone'
        DOCKER_CREDENTIALS_ID = 'dock-cre'
        GIT_CREDENTIALS_ID    = 'git-pat'
    }

    stages {
        stage('Setup and Build') {
            steps {
                container('docker') {
                    script {
                        // --- Setup ---
                        echo '1. Checking out source code and installing dependencies...'
                        checkout scm
                        sh 'apk add --no-cache git sed'

                        // --- Define and Write Tag to a file ---
                        // **ECHO 1: In ra các biến môi trường để debug**
                        echo "=========================================="
                        echo "GIT TAG DETECTED (env.TAG_NAME): ${env.TAG_NAME}"
                        echo "JENKINS BUILD NUMBER (env.BUILD_NUMBER): ${env.BUILD_NUMBER}"
                        echo "=========================================="
                        
                        def imageTag = env.TAG_NAME ?: "dev-${env.BUILD_NUMBER}"
                        
                        // **ECHO 2: In ra tag cuối cùng đã được xác định**
                        echo "==> Final image tag for this build is: ${imageTag}"
                        
                        sh "echo ${imageTag} > image.tag"

                        // --- Build Image ---
                        echo "2. Building image: ${BACKEND_IMAGE_NAME}:${imageTag}"
                        docker.build("${BACKEND_IMAGE_NAME}:${imageTag}", "./backend")
                    }
                }
            }
        }

        stage('Publish and Deploy Release') {
            when {
                tag pattern: "*", comparator: "REGEXP"
            }
            steps {
                container('docker') {
                    script {
                        // **ECHO 3: Đọc lại tag từ file để xác nhận**
                        def releaseTag = readFile('image.tag').trim()
                        echo "=========================================="
                        echo "ENTERING DEPLOYMENT STAGE"
                        echo "Tag read from file for deployment: ${releaseTag}"
                        echo "=========================================="
                        
                        // --- Push Release Image ---
                        echo "3. Publishing release image: ${BACKEND_IMAGE_NAME}:${releaseTag}"
                        docker.withRegistry("https://index.docker.io/v1/", DOCKER_CREDENTIALS_ID) {
                            docker.image("${BACKEND_IMAGE_NAME}:${releaseTag}").push()
                        }

                        // --- Update Config Repo ---
                        echo "4. Updating config repo to release version: ${releaseTag}"
                        withCredentials([string(credentialsId: GIT_CREDENTIALS_ID, variable: 'GIT_TOKEN')]) {
                            sh "rm -rf ${CONFIG_REPO_DIR}"
                            sh "git clone https://${GIT_TOKEN}@github.com/chuitrai/my_app_config.git ${CONFIG_REPO_DIR}"
                            
                            dir(CONFIG_REPO_DIR) {
                                sh "git config user.email 'jenkins-bot@example.com'"
                                sh "git config user.name 'Jenkins Bot'"
                                sh "sed -i 's|^    tag: .*#backend-tag|    tag: ${releaseTag} #backend-tag|' values.yaml"
                                
                                def changes = sh(script: "git status --porcelain", returnStdout: true).trim()
                                if (changes) {
                                    sh """
                                        git add .
                                        git commit -m 'CI: Release backend version ${env.TAG_NAME}'
                                        git push origin main
                                    """
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
    }
}