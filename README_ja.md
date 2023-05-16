# Cusi

* [English](README.md)
* 日本語

CusiはM5Stack MicroPython環境(UIFlow)向けのコマンドラインツールです。
M5Stackデバイスに対するファイルの読み書きが行えます。

# 特徴
* マルチプラットフォーム
* 単一のバイナリファイルで動作
* シリアルポートの一覧表示
* デバイス上のファイル/ディレクトリ一覧表示
* デバイスへのローカルファイルのアップロード
* デバイス上のファイルダウンロード又は表示
* デバイス上のファイル削除

# 使用方法
## 準備
1. [releaseページ](https://github.com/zuku/cusi/releases/latest)から `cusi-vN.N.N.zip` という形式の名前のファイルをダウンロードします。 (`N` 部分は数字に置き換えます 例: `cusi-v1.0.0.zip`)
2. ダウンロードしたZIPファイルを展開します。
3. 展開した中から使用している環境向けのコマンドファイル(`cusi` 又は `cusi.exe`)を見つけます。

|Directory      |Platform         |
|---------------|-----------------|
|`darwin_amd64` |Intel Mac        |
|`darwin_arm64` |Apple Silicon Mac|
|`linux_amd64`  |Linux (x86_64)   |
|`windows_amd64`|Windows (x86_64) |

ZIPファイルのダウンロード時、「危険なファイル」とWebブラウザが警告する場合があります。
安全性が気になる場合はGitHubからソースをダウンロードして(内容をレビューしてから)自身のGo環境でビルドしてください。

### macOS
macOSではGatekeeperが開発元が未確認のアプリケーション実行をブロックします。
`cusi` コマンドを実行したい場合は以下の手順を行なってください。

1. Finderでコマンドファイルのあるフォルダを開きます。
2. controlキーを押しながらコマンドファイルをクリック(又は右ボタンクリック)して、表示されたショートカットメニューから _開く_ を選択します。
3. ダイアログボックスの _開く_ をクリックします。
4. ターミナルが起動してコマンドが実行されます。実行後ターミナルの画面を閉じます。
5. 一度コマンドが実行された後はお好きなターミナルアプリを使用してコマンドを実行できます。

## コマンド詳細

### デバイスへの接続

最初にM5StackデバイスをUSBモードに設定し、コンピュータに接続します。

#### macOS
シリアルポートの一覧表示

```
$ cusi -l
/dev/cu.Bluetooth-Incoming-Port
/dev/cu.usbserial-XXXXXXXXXX
/dev/cu.wlan-debug
/dev/tty.Bluetooth-Incoming-Port
/dev/tty.usbserial-XXXXXXXXXX
/dev/tty.wlan-debug
```
接続には `/dev/tty.usbserial-XXXXXXXXXX` を使用します。
`XXXXXXXXXX` 部分はデバイスごとに異なる16進数文字列です。

```
$ cusi /dev/tty.usbserial-XXXXXXXXXX
```

#### Windows
シリアルポートの一覧表示

```
$ cusi.exe -l
COMX
```
`COMX` の `X` 部分は数字です。例えば `COM1`, `COM3` といった名前になります。

```
$ cusi.exe COMX
```

### プロンプト
デバイスに接続されると入力を待つプロンプトが表示されます。

```
>
```
#### ファイル/ディレクトリ一覧
```
> ls
apps
blocks
boot.py
emojiImg
img
main.py
res
temp.py
test.py
update
```

#### ファイルアップロード
```
> put /path/to/my_app.py apps/my_app.py
Uploading...
1234 / 1234 bytes
```

#### ディレクトリ内の一覧
```
> ls apps
my_app.py
```

#### 終了
```
> exit
```

#### ヘルプ
`help` と入力するとより詳しい情報が表示されます。
```
> help
```

# ライセンス
CusiはMITライセンスでリリースされています。[LICENSE](./LICENSE)をご確認ください。
使用しているライブラリ等のライセンスは[THIRD-PARTY-NOTICES.txt](./THIRD-PARTY-NOTICES.txt)をご確認ください。
