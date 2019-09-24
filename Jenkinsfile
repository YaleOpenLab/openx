node {
  docker.image('golang').inside {

		stage('Print go version') {
    	sh 'go version'
  	}

		environment {
			WORKSPACE = '/home/jenkins'
		}

		stage('Get openx package') {
			script {
				withEnv(["GOPATH=${env.WORKSPACE}/go", "GOROOT=${root}", "GOBIN=${root}/bin", "PATH+GO=${root}/bin"]) {
					sh 'echo "$GOPATH"'
					sh 'go get -v github.com/YaleOpenLab/openx'
				}
			}
		}
  }
}
