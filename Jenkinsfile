// Jenkinsfile - Hỗ trợ cả CI (cho nhánh) và CD (cho tag)

pipeline {
    // ---- ĐỊNH NGHĨA AGENT TƯỜNG MINH ----
    agent {
        kubernetes {
            // Định nghĩa Pod Template bằng YAML
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
        DOCKER_USERNAME       = 'chuitrai2901'
        BACKEND_IMAGE_NAME    = "${DOCKER_USERNAME}/my-go-backend"
        CONFIG_REPO_URL_HTTPS = 'https://github.com/chuitrai/my_app_config.git'
        CONFIG_REPO_DIR       = 'my_app_config_clone'
        DOCKER_CREDENTIALS_ID = 'dock-cre'
        GIT_CREDENTIALS_ID    = 'github-pat'
    }

    stages {
        // Giai đoạn 1: Luôn chạy để chuẩn bị môi trường
        stage('Setup Environment') {
            steps {
                container('docker') {
                    script {
                        echo "Checking out code..."
                        checkout scm
                        echo "Installing dependencies..."
                        sh 'apk add --no-cache git sed'
                    }
                }
            }
        }
        
        // Giai đoạn 2: Luôn chạy để build và xác nhận code
        stage('Build Image') {
            steps {
                container('docker') {
                    script {
                        // Xác định tag: nếu là build từ Git tag, dùng tên tag. Nếu không, dùng dev-<số-build>.
                        env.IMAGE_TAG = env.TAG_NAME ?: "dev-${env.BUILD_NUMBER}"
                        echo "Building image with tag: ${env.IMAGE_TAG}"
                        docker.build("${BACKEND_IMAGE_NAME}:${env.IMAGE_TAG}", "./backend")
                    }
                }
            }
        }

        // Giai đoạn 3: Chỉ chạy khi được trigger bởi một Git Tag
        stage('Publish and Deploy Release') {
            // ---- ĐIỀU KIỆN QUAN TRỌNG NHẤT ----
            when {
                tag pattern: "v.*"
            }
            steps {
                container('docker') {
                    script {
                        // --- Push Release Image ---
                        echo "Publishing release image: ${BACKEND_IMAGE_NAME}:${env.TAG_NAME}"
                        docker.withRegistry("https://index.docker.io/v1/", DOCKER_CREDENTIALS_ID) {
                            docker.image("${BACKEND_IMAGE_NAME}:${env.TAG_NAME}").push()
                        }

                        // --- Update Config Repo ---
                        echo "Updating config repo to release version: ${env.TAG_NAME}"
                        withCredentials([string(credentialsId: GIT_CREDENTIALS_ID, variable: 'GIT_TOKEN')]) {
                            sh "rm -rf ${CONFIG_REPO_DIR}"
                            sh "git clone https://${GIT_TOKEN}@github.com/chuitrai/my_app_config.git ${CONFIG_REPO_DIR}"
                            
                            dir(CONFIG_REPO_DIR) {
                                sh "git config user.email 'jenkins-bot@example.com'"
                                sh "git config user.name 'Jenkins Bot'"
                                sh "sed -i 's|^    tag: .*#backend-tag|    tag: ${env.TAG_NAME} #backend-tag|' values.yaml"
                                sh "git add . ; git commit -m 'CI: Release backend version ${env.TAG_NAME}' ; git push origin main"
                                echo "Successfully pushed configuration update."
                            }
                        }
                    }
                }
            }
        }
    }
}