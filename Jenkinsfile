pipeline {
    agent { docker { image 'golang' } }
    stages {
        stage('build') {
            steps {
                sh 'go version'
								sh 'go env gocache'
								sh 'go get -v github.com/YaleOpenLab/openx'
            }
        }
    }
}
