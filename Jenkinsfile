pipeline {

	environment{
		WORKSPACE = '/home/jenkins'
	}

	agent {
		docker { image 'golang'}
	}

	stages {
		stage ('Test') {
			steps {
				sh 'echo $GOPATH'
				sh 'echo $HOME'
			}
		}
	}
}
