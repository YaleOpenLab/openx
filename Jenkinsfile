pipeline {

	environment {
		WORKSPACE = '/home/jenkins'
	}

	agent {
		docker { image 'golang'}
	}

	stages {
		stage('Print go version') {
			steps {
				sh 'go version'
			}
		}

		stage('Get openx package') {
			steps {
				sh 'echo "$GOPATH"'
				sh 'go get -v github.com/YaleOpenLab/openx'
			}
		}
	}
}
