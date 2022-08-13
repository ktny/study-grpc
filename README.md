# grpc学習

## grpcとは

- Googleで開発された多言語間のRPC(Remote Procedure Call)を実現するプロトコル
- ネットワークを介して異なる言語のプログラムを関数呼び出しすることができる
- 異なる言語間での呼び出しができるようにIDL(Interface Description Language)と呼ばれるインターフェースを定義する

### grpcの長所

- HTTP/2による高速な通信
    - テキストではなくバイナリにシリアライズされて小さな容量で転送される
    - ひとつのコネクションで複数のリクエスト/レスポンスをやり取りできる
    - リクエストの度に接続/切断を行う必要がなく、ヘッダーを都度送る必要もない
- IDL は Protocol Buffers
    - スキーマファーストの開発になり生産性が高くなる
    - コンパイルすると任意の言語のサーバ/クライアント用コードを自動生成できる
    - カスタマイズすることでJSONのレスポンスに変えられるなど拡張性が高い
- 柔軟な4つのストリーム対応
    - シンプルなRPC
        - クライアントから送られた1リクエストに対して、レスポンスを1度返して終了する
    - サーバーストリーミングRPC
        - レスポンスを複数回に分けて返すことで、時間のかかる処理を非同期的に処理できる
        - クライアントからポーリング通信を行うのではなく、サーバからpushできるので無駄な通信を減らせる
        - サーバ側から状態変化時にクライアントに直接伝えられるのでネットワーク負荷を最小化できる
    - クライアントストリーミングRPC
        - クライアントからリクエストを分割して贈り、サーバは全てのリクエストを受ける前に処理を逐次開始する
        - サーバは全てのリクエストを受け取ってからレスポンスを返す
        - 大きなデータを分割してアップロードしたい場合などに有用
    - 双方向ストリーミングRPC
        - クライアントから初めのリクエストが送られたあと、サーバ・クライアントどちらも任意のタイミングでリクエスト/レスポンスを送ることができる
        - チャットやゲームなどで有用
        - REST APIではストリーミングが行えないので、WebSocketサーバを別途立てる必要があった。gRPCの場合、単一のサーバで双方向通信が必要なAPIとしても対応可能
- 強力なタイムアウト、キャンセル
    - リクエストチェーン全体に渡るタイムアウトやキャンセルをプロトコルレベルでサポート

### grpcの短所

- HTTP/2非対応のアーキテクチャもある
- ブラウザの対応状況が不十分
    - gRPC-webが2018年にリリースされたが、EnvoyでProxyサーバを立てる必要があったり、双方向やクライアントストリーミングRPCに非対応など、フロントエンドとのアプリケーション通信には難がある
- 言語によって実装状況にばらつきがある
    - Go, Javaなどはクライアントロードバランシングが実装されているが、Ruby, PHPなどはなかったりする
- バイナリにシリアライズすると人間が読めない
    - RESTでは出力を目視できるが、gRPCは専用のクライアントをインストールする必要がある
- RESTも十分速い
    - gRPCはRESTに比べてたしかに速いが何倍も改善できるほど速いわけではない

### 他のプロトコルとの比較

| 機能                           | gRPC                         | REST           | GraphQL                    |
| ------------------------------ | ---------------------------- | -------------- | -------------------------- |
| スキーマ言語                   | Protocol Buffers             | OpenAPI        | GraphQL                    |
| クエリ言語                     | なし                         | なし           | あり                       |
| IDL コンパイラ                 | 公式が多数提供               | 公式提供なし   | 公式提供なし               |
| ストリーム処理                 | 双方向可能                   | サーバーサイド | [ドラフト][graphql-stream] |
| クライアント指定のタイムアウト | あり                         | なし           | なし                       |
| クライアントからのキャンセル   | あり                         | なし           | なし                       |
| バイナリデータ                 | 扱える                       | 扱える         | [実装次第][graphql-binary] |
| 最大メッセージサイズ           | デフォルト 4 MiB, 最大 4 GiB | 仕様上はなし   | 仕様上はなし               |

## Protocol Buffers

- service
    - APIにおけるサービスを定義する。サービスには複数のRPCメソッドを定義できる
- message
    - プログラミング言語の構造体やクラスに変換される概念
    - フィールドにはスカラ型とmessage型を扱うことができる
- タグナンバー
	- フィールド識別のため同じメッセージ内で一意である必要がある
	- 一度使ったタグナンバーは再利用せず廃盤にする必要がある
-  repeated
	- 配列を表現できる。スカラ型、メッセージ型どちらでも使える
- enum
	- 列挙型を定義できる
	- 要素の先頭は必ず`UNKNOWN = 0` である必要がある
	- `option allow_alias = true`で同じ値を異なるラベルに割り当てられる
- マップ
	- `map<key_type, value_type> map_field = N`
	- キーにできるのは整数値、文字列、真偽値のみ
- oneof
	- フィールドの先頭にoneofと付与することで「複数の型の中からどれかひとつ」という定義をできる

```proto
message GreetingCard {
int32 id = 1;
oneof message {
    string text = 2;
    Image image = 3;
    Video video = 4;
}
message Image {...}
message Video {...}
}
```

- Well Known Types
	- Googleが定義した便利なメッセージ型
	    - `google.protobuf.Timestamp`: 日時を表す型
	    - `google.protobuf.Duration`: 期間を表す型
	    - `google.protobuf.Empty`: 特に値を返す必要がない場合
	    - `google.protobuf.Any`: 明示的に定義されていない型

## 本リポジトリでの使い方

### 準備

[Quick start | Go | gRPC](https://grpc.io/docs/languages/go/quickstart/)

```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
export PATH="$PATH:$(go env GOPATH)/bin"
```

### protoファイルをコンパイル

```sh
protoc -I. -Iinclude --go_out=module=github.com/ktny/study-grpc:. deepthought.proto
protoc -I. -Iinclude --go-grpc_out=module=github.com/ktny/study-grpc:. deepthought.proto
```
