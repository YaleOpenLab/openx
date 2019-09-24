node {

	environment {
		WORKSPACE = '/home/jenkins'
	}

	agent {
		docker { image: 'golang'}
	}

	stages {
		stage('Print go version') {
			sh 'go version'
		}

		stage('Get openx package') {
			sh 'echo "$GOPATH"'
			sh 'go get -v github.com/YaleOpenLab/openx'
		}
	}
}
