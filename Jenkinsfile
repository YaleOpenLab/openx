node {
  try {
      githubNotify(status: 'INPROGRESS')
      ws("${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/") {
        withEnv(["GOPATH=${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"]) {
          env.PATH="${GOPATH}/bin:$PATH"

          stage('Pre Test') {
            echo 'Printing go version'
            sh 'go version'
          }

          stage('Build') {
            //List all our project files with 'go list ./... | grep -v /vendor/ | grep -v github.com | grep -v golang.org'
            //Push our project files relative to ./src
            sh 'cd $GOPATH && go list ./...'

            echo 'Building'
            sh """cd $GOPATH && go build -v ./..."""
          }
        }
      }
  } catch (e) {
    // If there was an exception thrown, the build FAILURE
    currentBuild.result = "FAILURE"
    githubNotify(status: 'FAILURE')
  }
}
