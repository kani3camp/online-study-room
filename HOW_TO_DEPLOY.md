# 開発中
## デプロイ先
- AWS Lambda（API）
- AWS API Gateway（API）
- Cloud Functions（API）
- （WebアプリはローカルでOK）
- （モバイルアプリはローカルでOK）

## デプロイ方法
### AWS Lambda
開発用のLambda関数名にはprefixとして`test_`をつける。
1. コード内で、Firestoreの認証情報（common.goのSecretManagerSecretNameという定数）が開発用のプロジェクトのものか確認する
1. コード内で、ProjectIdが開発用のプロジェクトのものか確認する
1. AWS CLIで関数をデプロイする


### AWS API Gateway
websocketと普通のhttpの2種類。
開発用とリリース用で別々にAPIを用意する（ステージで分けてもいいけどその都度関数を書き換えていくのが面倒なので）。
websocketのほうは自動デプロイができないため、統合の切り替え・追加・削除といったAPI Gateway側での変更を行った場合には毎度手動でデプロイする必要がある。
httpのほうは$defaultというステージで自動デプロイをオンにしてあるため、手動デプロイは必要ない。
websocketもhttpも、統合するLambda関数を更新しただけであればデプロイは必要ない。
1. websocketは開発用のAPIで手動でデプロイする
1. クライアント側で正しいurlになっているか確認する


### Cloud Functions（API）
FirebaseAuthNewUserListener関数だけはAWSではなくCloud Functionsにデプロイする。
1. コンパイルが通るように、リポジトリのcloud_functionsディレクトリに行き、aws_lambdaディレクトリからコードを手動でマージする。package名はgo_apiのはず。LambdaだとFirestoreの認証をSecret Managerから取得する処理があるが、Cloud Functionsなら必要ないのでその部分のコードはマージしない。
1. gcloud CLIでprojectが開発用のプロジェクトになっていることを確認する
1. 関数をデプロイする




# リリース
- AWS Lambda（API）
- AWS API Gateway（API）
- Cloud Functions（API）
- Firebase Hosting（Webアプリ）
- Google Play Store（モバイルアプリ）

## デプロイ方法
### AWS Lambda
本番用のLambda関数には`test_`がついてない名前の関数を用意する。
1. コード内で、Firestoreの認証情報を本番用のプロジェクトのものに切り替える
1. コード内で、ProjectIdを本番用のプロジェクトのものに切り替える

### AWS API Gateway
1. websocketは本番用のAPIで手動デプロイする。
1. クライアント側でAPIのurlを本番用に書き換える


### Cloud Functions
FirebaseAuthNewUserListener関数だけはAWSではなくCloud Functionsにデプロイする。
1. コンパイルが通るように、リポジトリのcloud_functionsディレクトリに行き、aws_lambdaディレクトリからコードを手動でマージする。package名はgo_apiのはず。LambdaだとFirestoreの認証をSecret Managerから取得する処理があるが、Cloud Functionsなら必要ないのでその部分のコードはマージしない。
1. gcloud CLIでprojectが本番用のプロジェクトになっていることを確認する
1. 関数をデプロイする

### Firebase Hosting
1. firebase CLIでプロジェクトが本番用のものであることを確認する
1. `firebase deploy`でデプロイする


### Google Play Store
次のことに注意してリリースする。
- アプリバージョン
- Firebase Authの認証が有効なものか(google-service.jsonみたいなの。SHA-1、SHA-256を登録するやつ)

