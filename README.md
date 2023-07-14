# Go-Academy-Presentation

## 実装メモ

### golang-migration

主な参考

https://kakakakakku.hatenablog.com/entry/2022/09/26/131311

マイグレーションファイル作成
```sh
migrate create --ext sql --dir db/migrate --seq users
```

マイグレーション実行
```sh
migrate -path db/migrate -database 'mysql://root:password@tcp(localhost:3306)/academy15' up
```

初期化
```sh
migrate -path db/migrate -database 'mysql://root:password@tcp(localhost:3306)/academy15' drop
```