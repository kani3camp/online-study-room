# youtube_displayer

## デプロイ方法
このプロジェクトで生成した静的ファイルを `https://github.com/sorarideblog/youtube-displayer-hosting` へデプロイする。
具体的には、
 - `npm run generate` でできたdistディレクトリの中身を全てhttps://github.com/sorarideblog/youtube-displayer-hostingのローカルフォルダへコピペする
 - 注意：コピー先の.gitのフォルダと.gitignoreを削除しないように、消してからペーストじゃなくてそのまま上書きすると良い。
 - 注意：また、.nojekyllとかいうファイルはペーストできないらしいのでコピーしないこと。

## Build Setup

```bash
# install dependencies
$ npm install

# serve with hot reload at localhost:3000
$ npm run dev

# build for production and launch server
$ npm run build
$ npm run start

# generate static project
$ npm run generate
```
