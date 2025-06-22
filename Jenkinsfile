// Jenkinsfile - Sửa lỗi cú pháp 'when' và lỗi scope của biến

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
        // Giữ các biến tĩnh ở đây
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
                        echo 'Checking out source code and installing dependencies...'
                        checkout scm
                        sh 'apk add --no-cache git sed'

                        // --- Define and Write Tag to a file ---
                        def imageTag = env.TAG_NAME ?: "dev-${env.BUILD_NUMBER}"
                        echo "Determined image tag: ${imageTag}"
                        // Ghi tag vào một file trong workspace để các stage sau có thể đọc
                        sh "echo ${imageTag} > image.tag"

                        // --- Build Image ---
                        echo "Building image: ${BACKEND_IMAGE_NAME}:${imageTag}"
                        docker.build("${BACKEND_IMAGE_NAME}:${imageTag}", "./backend")
                    }
                }
            }
        }

        stage('Publish and Deploy Release') {
            // Sửa lại cú pháp 'when'
            when {
                tag pattern: "*", comparator: "REGEXP"
            }
            steps {
                container('docker') {
                    script {
                        // Đọc tag từ file đã lưu
                        def releaseTag = readFile('image.tag').trim()
                        
                        // --- Push Release Image ---
                        echo "Publishing release image: ${BACKEND_IMAGE_NAME}:${releaseTag}"
                        docker.withRegistry("https://index.docker.io/v1/", DOCKER_CREDENTIALS_ID) {
                            docker.image("${BACKEND_IMAGE_NAME}:${releaseTag}").push()
                        }

                        // --- Update Config Repo ---
                        echo "Updating config repo to release version: ${releaseTag}"
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
                                        git push origin my_app
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