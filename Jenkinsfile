pipeline {

	agent {
		docker { image 'golang'}
	}

	environment {
		WORKSPACE = '/home/jenkins'
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
