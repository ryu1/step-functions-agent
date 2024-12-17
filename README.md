# _Step Functions Agent_

Description:
_サーバ上に常駐し、AWS Step Functionsのタスクに定義された任意のコマンドをサーバ上で実行します。_



## Development

### Requirement

- Go v1.12.9
- Go Modules

### Project Structure

```
step-functions-agent
├── .go-version ... goenv用（goenv使用しない場合は意識する必要はない）
├── go.mod      ... ライブラリ管理ファイル
├── go.sum      ... ライブラリ管理ファイル
├── main.go     ... 本プログラムのソースファイル
└── Makefile    ... makeコマンドでビルドを自動化しています
```

## Building

```
$ make
```

`dist`配下に、ビルド成果物として、`step-functions-agent`ファイルが生成されます。

## Deployment

1. `Building`の手順により生成されたstep-functions-agent`ファイルをサーバ上に配置します

   ```
   $ sudo su -
   $ mkdir -p /tool/step-functions-agent/bin
   $ mv ~/step-functions-agent /tool/step-functions-agent/bin/
   $ chmod +x /tool/step-functions-agent/bin/step-functions-agent
   ```

2. サービスとして起動し、プロセスの死活管理を行うためimmortalを導入します
   1. immortalをインストールします

      ```
      $ curl -s https://packagecloud.io/install/repositories/immortal/immortal/script.rpm.sh | sudo bash
      $ yum install immortal
      ```

   2. 設定ファイルを作成し配置します

      ```
      $ cd /etc/immortal
      $ touch step-functions-agent.yml
      ```

       `step-functions-agent.yml`の内容

      `{ActivityのARN}`には、`Using`の項目で、AWS Step Functionsに作成したアクティビティのARNを指定します。

      https://ap-northeast-1.console.aws.amazon.com/states/home?region=ap-northeast-1#/activities

      ```
      cmd: /tool/step-functions-agent/bin/step-functions-agent -arn {ActivityのARN}
      cwd: /root
      env:
          HOME: /root
      log:
          file: /var/log/step-functions-agent.log
          age: 86400
          num: 7
          size: 1
          timestamp: true
      user: root
      ```

   3. サービスとして起動する

      1. ステータスを確認します

      ```
      $ immortalctl status
      ```

      2. サービスが起動してなかったら再起動します

      ```
      $ immortalctl restart step-functions-agent
      ```

## Using

1. AWS Step Functionsにアクティビティを作成する

2. AWS Step Functionsのステートマシンに、以下のようにタスクを定義して使用します

```
 {
    "Type": "Task",
    "Parameters": {
      "Command": "/PAHT/TO/COMMAND" # ここに実行したいコマンドを指定する
    },
    "TimeoutSeconds": 300,
    "HeartbeatSeconds": 60,
    "Resource": "1で作成したAcitvityのARNを指定",
}
```
