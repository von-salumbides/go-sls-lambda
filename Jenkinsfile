pipeline {
    agent any
    environment {
        GO111MODULE = 'on'
        GOOS = 'linux'
        GOARCH = 'amd64'
        CGO_ENABLED = '0'
    }
    stages {
        stage('Switch Environment') {
            steps {
                script {
                    switch(DEPLOY_ENV) {
                        case "dev":
                            AWS_ROLE = ""
                        break
                            error("Build Failed for ${DEPLOY_ENV}. No match found.")
                    }
                }
            }
        }
        stage('Deploy') {
            steps {
                script {
                    currentBuild.displayName = "${SLS_ACTION}-${FUNCTION_NAME}"
                    if ( SLS_ACTION == 'deploy') {
                        if ( FUNCTION_NAME == 'all' ) {
                            sh "make deploy"
                        } else {
                            sh "deployFunc"
                        }
                    } else if ( SLS_ACTION == 'remove' ) {
                        sh "make remove"
                    } else {
                        error("Build Failed, ${SLS_ACTION} is not defined")
                    }
                }
            }
        }
    }
    post {
        always {
            cleanWs()
        }
    }
}