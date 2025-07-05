---
title: Google OSSにマージされた話
labels:
    - GitHub
    - Google
    - OSS
    - アドベントカレンダー
blogger_id: "5805870221516918520"
slug: google-ossにマーシ-された話
status: publish
published: "2022-12-07T00:00:00+09:00"
---
この投稿は私の現在所属しているスターフェスティバル株式会社の [スターフェスティバル Advent Calendar 2022](https://qiita.com/advent-calendar/2022/stafes) の7日目記事になります。

昨日は [yui\_tang](https://twitter.com/yui_tang) さんによる機械学習でメタルかそうでないかを聞き分ける自作アプリの紹介でした！面白かったので是非読んでみてください！

[https://zenn.dev/stafes\_blog/articles/89a29ade69ec6d](https://zenn.dev/stafes_blog/articles/89a29ade69ec6d)

こんにちは！ [k1rnt](https://twitter.com/k1rnt) です。

今回はGoogle OSSにissueを出しPRの作成、マージまでの経緯をお話しできればと思います。

## なぜやろうと思ったのか

つい先日GitHub Actionsの _save-state_ と _set-output_ コマンドが廃止になるという速報が社内Slackに投稿されました

[https://github.blog/changelog/2022-10-11-github-actions-deprecating-save-state-and-set-output-commands/](https://github.blog/changelog/2022-10-11-github-actions-deprecating-save-state-and-set-output-commands/)

こちらですね。

社内のリポジトリのGithub Actionsでも結構 _set-output_ は使われていることもあり自分も書いていました。

そんな中弊社、 [ikkitang](https://stafes.notion.site/Ikki-Takahashi-948c3d875cab4ba9a06382aaa81d4585) からこんな一言が！

[![](images/4031f7c38750.jpg)](https://blogger.googleusercontent.com/img/a/AVvXsEgz3IskaLlbIDe6_2jY1FbW-a-0Gkm1MlfjHSedkFnD82mVyOhuLatRlLdn8GvwNdkta-y0Rvz_WjHLBnjIzq1g-TFWAydL06ME8UXazD0kKlNj4ZySLP9Yz4PdcmPxEMxujNRuooDPSiPFRhPlj8B9UmpQ7YpGMY9l7Tm6nnZRTyGjMISMvucU6tkk)

簡単な内容な割にはクリティカルで影響範囲も大きいはず...！

そう感じた私はひとまずGoogle OSSのGitHubで _set-output_ が使われていないか検索をしてみることにしました。

## OSSコミットチャンス

Google OSSを検索するとset-outputを使用しているリポジトリが多々ありました。

その中でも一度使ったことある [go-github](https://github.com/google/go-github) も対象に含まれていたのでこちらに貢献することを決めました。

## CONTRIBUTING.mdを読む

まずどのように貢献したら良いかは大抵CONTRIBUTING.mdといったファイルがあるはずなので読みます。

今回の場合 [こちら](https://github.com/google/go-github/blob/master/CONTRIBUTING.md)

これを読むとissueにて問題の報告PRにてパッチの提出を行えば良いことがわかります。

このドキュメントにはパッチの提出の際 go fmt の実行やテストを書く必要があると書かれていますが今回修正するファイルはGitHub Actionsなのでこれらは無視します。

## issueを書く

OSSではissueだけでなくPRなども全て英語で書く必要がありますがdeepl等で翻訳したテキストを投下しても問題ないと思われます。

今回書いたissueはこちらになります。

[https://github.com/google/go-github/issues/2491](https://github.com/google/go-github/issues/2491)

[![](images/c1ff7cf6b203.jpg)](https://blogger.googleusercontent.com/img/a/AVvXsEgvhFJC8oyVI2hcCyNqoUIDyPzF7aIXsUgOFZn02kbBYjIW9tIN2QQrUjXCioE4JE01zdivH6P2hl4PgfAEpjPYXP54gvMlOU4_rivFZMfYYzGwZ0Iz33EwfMDeDUkB-BCehEMJjP4y98mKL9LTx3g9xC3gzd0yX165txHbKWRn8emNHMB_tqJac8fd)

今回はフォーマットなどが無かったのでなるべく簡潔に書きました。

リアクションがついて嬉しかったので皆さんも是非issueにもリアクションしていきましょう。

## リポジトリをForkして作業

他社のリポジトリに貢献する際、通常はそのリポジトリをForkして作業を行います。

今回行った修正はこちらになります。

[https://github.com/google/go-github/pull/2492/commits/fa30a07caf75a5a8354e93221f82d5701f991d4a](https://github.com/google/go-github/pull/2492/commits/fa30a07caf75a5a8354e93221f82d5701f991d4a)

## PRを提出

少しドキドキしながらPRの提出を行いました。

[https://github.com/google/go-github/pull/2492](https://github.com/google/go-github/pull/2492)

[![](images/0657e58de088.jpg)](https://blogger.googleusercontent.com/img/a/AVvXsEhz3KZE-pgEE2dNgWX1BJAbI-JS7dA6oeeYliMWiJwiUvY97kdqew31E5oO6S25F_80O2iVVfU4VRjnwYafRDcitqwQteThyMgvi4eNvjsqXmCdWhu6LxuYzN9pLTc34Ufy1gyB8VCFz5cgjxUsRq3vmKox5YHhTxuV0_gVvOj8BuuQX3dX9q5ZJgla)

私はGoogle OSSに貢献するのが初めてだったのでCLAに署名する必要がありました。

こちらはリンクから簡単に行えました。

CLAとは？

_CLA は Contributor License Agreement (コントリビューター ライセンス契約) の頭文字であり、オープンソース プロジェクトの主体と、コードを提供する個人開発者あるいは企業との間で締結される契約を意味します。_

[_https://prtimes.jp/main/html/rd/p/000000052.000042042.html_](https://prtimes.jp/main/html/rd/p/000000052.000042042.html)

## 変更リクエストがありました

[![](images/931c0d9305c1.jpg)](https://blogger.googleusercontent.com/img/a/AVvXsEj4RVU-dBdeAVdGZ-ylVtiFaF1hpaJpAE4sJV7rfNvHKNwZouzi3CPriEy7UoX-zqmfwWixT3H4XW93RTlcwIX2uiR1xDptxDD9GXmxg5IoRn37BLpVzzAAbKGsvIX1cB_gdX-hflPCIQii1tCB3gHTY68p0WaaZHDIoc5EFsH3mHNol5MHBY9a-NCj)

私の作業のダブルクォーテーションの位置の間違いです。

丁寧にsuggestion機能でレビューを頂きました。

恥ずかしい\> <

上記を治すと、

他にもWindowsでのJobでtestが落ちる、などが発生しました。

[![](images/e214f741137a.jpg)](https://blogger.googleusercontent.com/img/a/AVvXsEgmTL-xEtSjFkh6B1wn_evDmY5zFK37EblSL7pZD9HivyDr78wQfvO4GuB4atPVo1VYrNY_cBb-iwhiCXd6CXCjWkYa5wGpXan6ZCQ_Up5Gmbp_yJZ8m1yC3PLJ4tliL0CckfWyVcMRGYdRawFTkLl7crMgEmilK1nExr9x69DV5IBS2i82yLwMZbzg)

こちらについてもしばらく時間がかかったもののWindowsのJobに対しては正しくshellの種類をセットしなければいけないと言ったものでした。

## マージされました！

こうしてやっとPRがマージされることになりました。

[![](images/603c068e35e3.jpg)](images/603c068e35e3.jpg)

初めてGoogle OSSに貢献できた...！！

## 終わりに

Google OSSへの貢献は今回初めてだったのですが、 [gmlewis](https://github.com/gmlewis) さんが非常に優しくレビューをして頂いたりリアクションがあったりと、とても勉強になり非常に良い体験ができました。

OSS貢献は敷居が高く感じられることも多々ありますが意外とやればすんなり出来るものです。

是非ここまで読んでいただいた皆さんも気軽にOSS貢献してみてください。

ありがとうございました！

明日は@delucciさんの記事になります！お楽しみに！