package main

import (
    "fmt"
    "log"
    "flag"
    "os"
    "os/exec"
    "encoding/json"
    "github.com/aws/aws-sdk-go/service/sfn"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws"
)

var (
    argProfile = flag.String("profile", "", "Specify a Profile Name to the AWS Shared Credential.")
    argRegion = flag.String("region", "ap-northeast-1", "Specify an AWS Region.")
    argActivityArn = flag.String("arn", "", "Specify an Activity ARN.")
)

type OutputMessage struct {
    Message string    `json:"message"`
}

func aws_sfn_client(profile string, region string) *sfn.SFN {
    var config aws.Config
    if profile != "" {
        creds := credentials.NewSharedCredentials("", profile)
        config = aws.Config{Region: aws.String(region), Credentials: creds}
    } else {
        config = aws.Config{Region: aws.String(region)}
    }
    sess := session.New(&config)
    sfn_client := sfn.New(sess)
    return sfn_client
}

func generageOutputMessage(message string) []byte {

    outputMessage := OutputMessage{Message: string(message)}
    jsonBytes, err := json.Marshal(outputMessage)
    if err != nil {
        log.Printf(err.Error())
    }
    return jsonBytes
}

func runTask(sfnSession *sfn.SFN , activity *sfn.GetActivityTaskOutput, command string) {
    cmd := exec.Command("sh", "-c", command)
    stdout, err := cmd.CombinedOutput()
    if err != nil {

        errMessage := err.Error() + ": " + string(stdout)
        log.Printf(string(errMessage))

        jsonBytes := generageOutputMessage(string(stdout))
        params := &sfn.SendTaskFailureInput{
            Error: aws.String(string(jsonBytes)),
            TaskToken: activity.TaskToken,
        }
        sfnSession.SendTaskFailure(params)
    }
    fmt.Println(string(stdout))
    jsonBytes := generageOutputMessage(string(stdout))

    params := &sfn.SendTaskSuccessInput{
        Output: aws.String(string(jsonBytes)),
        TaskToken: activity.TaskToken,
    }

    sfnSession.SendTaskSuccess(params)
}

func main() {
     flag.Parse()

     if *argActivityArn == ""  {
         fmt.Println("Please specify `-arn` option.")
         os.Exit(1)
     }

    sfn_client := aws_sfn_client(*argProfile, *argRegion)

    for {
        params := &sfn.GetActivityTaskInput{
            ActivityArn: aws.String(*argActivityArn),
        }
        activity, err := sfn_client.GetActivityTask(params)
        log.Printf(*activity.Input)

        if err != nil {
            log.Printf(err.Error())
            return
        } else if activity.TaskToken != nil {

            var input Input
            err := json.Unmarshal([]byte(*activity.Input), &input)
            if err != nil {
                log.Printf(err.Error())
                return
            }

            go runTask(sfn_client, activity, input.Command)
        }
    }
}

type Input struct {
    Command string `json:"Command"`
}
