pipeline {
    agent {
        label 'build-golang-stable'
    }
    options {
        skipDefaultCheckout true
    }
    stages {
        stage('Checkout') {
            steps{
                dir('/root/workspace/go/src/github.com/derekpedersen/imgur-go') {
                    checkout scm
                    sh 'make dependencies'
                }
            }
        }
        stage('Build') {
            steps{
                dir('/root/workspace/go/src/github.com/derekpedersen/imgur-go') {
                    sh 'make build'
                }
            }
        }
        stage('Test') {
            steps {
                withCredentials([[$class: 'StringBinding', credentialsId: 'IMGUR_API_KEY', variable: 'IMGUR_API_KEY']]) {
                    dir('/root/workspace/go/src/github.com/derekpedersen/imgur-go') {
                        sh 'make test'
                    }
                }
            }
        }
    }
    // post {
    //     always {
    //         withCredentials([[$class: 'StringBinding', credentialsId: 'IMGUR_GO_COVERALLS_TOKEN', variable: 'COVERALLS_TOKEN']]) {
    //             dir('/root/workspace/go/src/github.com/derekpedersen/imgur-go') {
    //                 step([$class: 'CoberturaPublisher', autoUpdateHealth: false, autoUpdateStability: false, coberturaReportFile: '**/cp.xml', failUnhealthy: false, failUnstable: false, maxNumberOfBuilds: 0, onlyStable: false, sourceEncoding: 'ASCII', zoomCoverageChart: false]) 
    //                 sh 'go get github.com/derekpedersen/goveralls'
    //                 sh 'goveralls -coverprofile=cp.out'
    //             }
    //         }
    //     }
    // }
}