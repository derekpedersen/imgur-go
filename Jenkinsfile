pipeline {
    agent {
        label 'build-golang-stable'
    }
    options {
        skipDefaultCheckout true
    }
    stages {
        stage('Checkout') {
            steps {
                dir('/root/workspace/go/src/github.com/derekpedersen/imgur-go') {
                    checkout scm
                }
            }
        }
        stage('dependencies') {
            steps {
                dir('/root/workspace/go/src/github.com/derekpedersen/imgur-go') {
                    sh 'make dependencies'
                }
            }
        }
        stage('build') {
            steps {
                dir('/root/workspace/go/src/github.com/derekpedersen/imgur-go') {
                    sh 'make build'
                }
            }
        }
        stage('test') {
            steps {
                dir('/root/workspace/go/src/github.com/derekpedersen/imgur-go') {
                    sh 'make test'
                }
            }
        }
    }
}
