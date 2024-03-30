#!/bin/bash

echo "Misskeyのホスト名を入力してください（例：example.com）:"
read host
echo "MisskeyのAPIトークンを入力してください:"
read token
echo "時報の文章を入力してください:"
read content
echo "時報の時刻を24時間形式で入力してください（例：23:59）:"
read post_time

echo "MISSKEY_HOST=$host" > .env
echo "MISSKEY_TOKEN=$token" >> .env
echo "MISSKEY_CONTENT=\"$content\"" >> .env

post_hour=$(echo $post_time | cut -d: -f1)
post_minute=$(echo $post_time | cut -d: -f2)

current_dir=$(pwd)

echo "以下のコマンドをcrontabに追加してください:"
echo "$post_minute $post_hour * * * cd $current_dir && ./misskeyJihou"

echo "設定が完了しました。"
