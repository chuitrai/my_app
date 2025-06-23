// Jenkinsfile - Phiên bản cuối cùng với Kaniko
pipeline {
    agent {
        kubernetes {
            // Định nghĩa Pod Template với container JNLP và KANIKO
            yaml """
            apiVersion: v1
            kind: Pod
            spec:
              containers:
              - name: jnlp
                image: jenkins/inbound-agent:3107.v665000b_51092-5
                args: ['\$(JENKINS_SECRET)', '\$(JENKINS_NAME)']
                workingDir: /home/jenkins/agent
                command:
                - sleep
                args:
                - infinity
              - name: kaniko
                image: gcr.io/kaniko-project/executor:v1.9.0-debug
                command:
                - sleep
                args:
                - infinity
                volumeMounts:
                - name: kaniko-secret
                  mountPath: /kaniko/.docker
              volumes:
              - name: kaniko-secret
                secret:
                  secretName: regcred # Tên secret ta đã tạo bằng kubectl
                  items:
                    - key: .dockerconfigjson
                      path: config.json
            """
            label 'k8s-agent-with-kaniko'
        }
    }

    environment {
        DOCKER_USERNAME       = 'chuitrai2901'
        BACKEND_IMAGE_NAME    = "${DOCKER_USERNAME}/my-go-backend"
        CONFIG_REPO_URL_HTTPS = 'https://github.com/chuitrai/my_app_config.git'
        CONFIG_REPO_DIR       = 'my_app_config_clone'
        GIT_CREDENTIALS_ID    = 'github-pat'
    }

    stages {
        stage('Checkout & Setup') {
            steps {
                // Chạy trong container jnlp mặc định
                container('jnlp') {
                    echo 'Checking out source code...'
                    checkout scm
                    echo 'Installing dependencies...'
                    sh 'apk add --no-cache git sed'
                }
            }
        }

        stage('Build & Push with Kaniko') {
            steps {
                // Chuyển sang container kaniko để thực hiện build
                container('kaniko') {
                    script {
                        def newTag = "v1.0.${env.BUILD_NUMBER}"
                        echo "Building and pushing image with Kaniko: ${BACKEND_IMAGE_NAME}:${newTag}"

                        // Chạy lệnh Kaniko executor
                        sh """
                        /kaniko/executor \
                          --context="dir:///home/jenkins/agent/backend" \
                          --dockerfile="dir:///home/jenkins/agent/backend/Dockerfile" \
                          --destination="${BACKEND_IMAGE_NAME}:${newTag}" \
                          --cache=true
                        """
                    }
                }
            }
        }

        stage('Update Config Repository') {
            steps {
                // Quay lại container jnlp để chạy các lệnh git
                container('jnlp') {
                    script {
                        def newTag = "v1.0.${env.BUILD_NUMBER}"
                        echo "Updating config repo with new image tag: ${newTag}"
                        withCredentials([string(credentialsId: GIT_CREDENTIALS_ID, variable: 'GIT_TOKEN')]) {
                            sh "rm -rf ${CONFIG_REPO_DIR}"
                            sh "git clone https://${GIT_TOKEN}@github.com/chuitrai/my_app_config.git ${CONFIG_REPO_DIR}"
                            dir(CONFIG_REPO_DIR) {
                                sh "git config user.email 'jenkins-bot@example.com'"
                                sh "git config user.name 'Jenkins Bot'"

                                sh "sed -i 's|tag:.*#backend-tag|tag: ${newTag} #backend-tag|' values.yaml"

                                sh """
                                    if ! git diff --quiet; then
                                        git add values.yaml
                                        git commit -m 'CI: Bump backend image to ${newTag}'
                                        git push origin main
                                        echo "Successfully pushed configuration update."
                                    else
                                        echo "No changes to commit."
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