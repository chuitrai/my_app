// Jenkinsfile - Phiên bản cuối cùng, định nghĩa Pod Agent dùng được Docker

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
    command:
    - sleep
    args:
    - infinity
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
        DOCKER_USERNAME       = 'chuitrai2901'
        BACKEND_IMAGE_NAME    = "${DOCKER_USERNAME}/my-go-backend"
        CONFIG_REPO_URL_HTTPS = 'https://github.com/chuitrai/my_app_config.git'
        CONFIG_REPO_DIR       = 'my_app_config_clone'
        DOCKER_CREDENTIALS_ID = 'dock-cre'
        GIT_CREDENTIALS_ID    = 'github-pat'
    }

    stages {
        stage('CI/CD Pipeline') {
            steps {
                container('docker') {
                    script {
                        // --- Stage: Setup ---
                        echo '🔧 Checking out source code and installing dependencies...'
                        checkout scm
                        sh 'apk add --no-cache git sed'

                        // --- Stage: Build & Push ---
                        def newTag = "v1.0.${env.BUILD_NUMBER}"
                        echo "🐳 Building and pushing image: ${BACKEND_IMAGE_NAME}:${newTag}"

                        docker.withRegistry("https://index.docker.io/v1/", DOCKER_CREDENTIALS_ID) {
                            def builtImage = docker.build("${BACKEND_IMAGE_NAME}:${newTag}", "./backend")
                            builtImage.push()
                        }

                        // --- Stage: Update Config ---
                        echo "🛠 Updating config repo with new image tag: ${newTag}"
                        withCredentials([string(credentialsId: GIT_CREDENTIALS_ID, variable: 'GIT_TOKEN')]) {
                            sh "rm -rf ${CONFIG_REPO_DIR}"
                            sh "git clone https://${GIT_TOKEN}@github.com/chuitrai/my_app_config.git ${CONFIG_REPO_DIR}"

                            dir(CONFIG_REPO_DIR) {
                                sh "git config user.email 'jenkins-bot@example.com'"
                                sh "git config user.name 'Jenkins Bot'"

                                // Cập nhật image tag bằng sed
                                sh "sed -i 's|tag:.*#backend-tag|tag: ${newTag} #backend-tag|' values.yaml"

                                // Commit nếu có thay đổi
                                sh """
                                    if ! git diff --quiet; then
                                        git add values.yaml
                                        git commit -m 'CI: Bump backend image to ${newTag}'
                                        git push origin main
                                        echo "✅ Successfully pushed configuration update."
                                    else
                                        echo "ℹ️ No changes to commit."
                                    fi
                                """
                            }
                        }
                    }
                }
            }
        }
    }
}
