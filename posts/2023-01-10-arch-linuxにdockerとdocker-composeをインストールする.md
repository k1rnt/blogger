---
title: Arch LinuxにDockerとDocker Composeをインストールする
labels:
    - archlinux
    - Docker
    - Linux
    - 雑記
blogger_id: "461080664898657547"
slug: arch-linuxにdockerとdocker-composeをインストールする
status: publish
published: "2023-01-10T08:30:00+09:00"
---
[![](images/5e969ec228d1.jpg)](https://blogger.googleusercontent.com/img/a/AVvXsEgr655lNcyk7ZDTydorIr_WZo2htMEP2_3k3R2YNchWW4oyUbei1pYvU_r1fwOthrkpj_8H0PjRQ-KDpUrP1dpcenKrQtGcSUtlz5cnrwRvaZhpNweSshZrULXsuEC5GPU4CRPnrRirh8Hc9q1NUBhCGR3OvVjgEU3xL04RNshMYpoDo9ughmPEo-NQ)

何故か苦戦したのでメモ代わりに

## 何に苦戦したのか

linux版のdocker-desktopとdockerの情報が混在していたのでインストールしたファイルが壊れていた

## インストール方法

まずdocker-desktopを過去に何らかの方法でインストールしていた場合下記の動作が重要

[https://docs.docker.jp/desktop/install/archlinux.html#desktop-archlinux-uninstall-docker-desktop](https://docs.docker.jp/desktop/install/archlinux.html#desktop-archlinux-uninstall-docker-desktop)

こちらにある完全アンインストールを行う

```
$ rm -r $HOME/.docker/desktop
$ sudo rm /usr/local/bin/com.docker.cli
$ sudo pacman -Rns docker-desktop
```

その後

```
$ vim $HOME/.docker/config.json
```

でファイルを開きcredsStoreとcurrentContextのキーとバリューを削除する

ここまできたら後は普通にインストール

```
$ sudo pacman -S docker docker-compose
```

起動

```
$ systemctl start docker.service
```

自動起動

```
$ systemctl enable docker.service
```

以上となります
