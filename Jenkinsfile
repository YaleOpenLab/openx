node {

	environment {
		WORKSPACE = '/home/jenkins'
	}

  docker.image('golang').inside {

		stage('Print go version') {
    	sh 'go version'
  	}

		stage('Get openx package') {
			sh 'echo "$GOPATH"'
			sh 'go get -v github.com/YaleOpenLab/openx'
		}
  }
}
