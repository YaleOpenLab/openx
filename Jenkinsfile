node('docker') {
    checkout scm
    stage('Build') {
        docker.image('go').inside {
            sh 'go version'
        }
    }
}
