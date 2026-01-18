---
title: Zeroconf の思想で考える LAN 内デバイス Discovery
labels:
    - tech
    - 雑記
blogger_id: ""
status: publish
---
![](/images/zeroconf_mdns_discovery.jpg)

LAN 内に存在するデバイスを見つけたい、という要求はよくある。
一方で、実際にそれを「ちゃんと」実装しようとすると、想像以上に設計要素が多いことに気づく。

- 同一ネットワークにあるはずなのに見つからない
- レスポンスは返るが、それが本当に対象か分からない
- 複数台存在した場合の扱い
- discovery が不安定なときの UX

本記事では、ある LAN 内デバイス制御ツールを作る過程で実装した
Zeroconf の思想に基づく mDNS discovery について整理する。  

なお、デバイスの操作・制御部分には踏み込まない。
あくまで「どう見つけるか」という設計に焦点を当てる。

---

## Zeroconf と mDNS の関係を整理する

まず用語を整理しておく。

- Zeroconf (Zero Configuration Networking)
  概念・思想。
  人手による設定なしにネットワーク上でサービスを発見するための考え方。

- mDNS (Multicast DNS)
  Zeroconf を構成する技術要素のひとつ。
  UDP/5353 を用いて、デバイスが自身の情報を LAN 内に広告する仕組み。

つまり、

> Zeroconf の思想を、具体的に実現するための中核技術が mDNS

という関係になる。
![](/images/zeroconf_mdns_discovery-1.jpg)

今回の discovery は、Zeroconf の思想に従い、mDNS を使って実装したものと位置づけるのが一番正確だ。

---

## なぜ Zeroconf / mDNS を選んだのか

LAN 内 discovery の手法はいくつかある。

- IP レンジスキャン
- SSDP / UPnP
- mDNS

今回 mDNS を選んだ理由はシンプルで、

- こちらから能動的にスキャンしない
- デバイスが「公開している情報」だけを扱える
- 仕組みとして正規で、IoT 機器に広く使われている
- 無闇にポートを叩かない

という点が大きい。

特に、

> 探しに行くのではなく、広告されているものを受け取る

という点は、Zeroconf の思想そのものでもある。

---

## mDNS は「見つかる」だけでは足りない

mDNS を使うと、確かに色々なレスポンスが返ってくる。
しかし、実装して最初に直面するのは次の問題だ。

> 見つかるが、それが使いたいデバイスとは限らない

実際には、次のようなものが平然と混ざる。

- 同一ベンダの別機器
- 無関係な IoT デバイス
- 名前が似ているだけの別物
- キャッシュされた古い情報

![](/images/zeroconf_mdns_discovery-2.jpg)

mDNS では各レコードに TTL が設定されており、またデバイスがネットワークから離脱する際にはGoodbye packet が送信されることもある。  
しかし実装や環境によっては、これらが必ずしも期待通りに機能するとは限らない。  

そのため、discovery（発見）と selection（選別）を明確に分けて設計した。

---

## Discovery の全体像（抽象）

![](/images/zeroconf_mdns_discovery-3.jpg)

実装上の流れは大まかに以下。

1. mDNS query を送信
2. 一定時間レスポンスを収集
3. raw なレスポンスを内部モデルに変換
4. 条件に合うものだけを候補として残す
5. 最低限の情報に正規化

ここで重要なのは、mDNS レスポンスを信頼しすぎないこと。

---

## service type と選別の考え方

mDNS では、サービスは `_http._tcp.local` のような service type として広告される。

これにより、

- どのプロトコルを前提としているか
- TCP / UDP の別
- 利用されるポート

といった情報を、問い合わせ段階である程度絞り込める。

実装では service type をそのまま信用するのではなく、「候補を狭めるためのヒント」として扱った。

TXT レコードも同様に、識別のヒントとして使えることは多いが、それ単体で確定判断を下すことは避けている。

---

## IPv4 / IPv6 の扱い

mDNS レスポンスには IPv4 と IPv6 の両方が含まれることがある。
どちらを使うかは、実装上の判断が必要になる。

今回は以下の優先順位で選択した。

1. IPv4 アドレスを優先
2. IPv6 の場合はリンクローカル（`fe80::`）を除外

```rust
let ip = addresses
    .iter()
    .find(|a| a.is_ipv4())
    .or_else(|| addresses.iter().find(|a| {
        if let IpAddr::V6(v6) = **a {
            // 簡易的な判定。厳密には is_unicast_link_local() を使うべきだが、
            // 実用上 fe80:: 以外のリンクローカルはほぼ存在しない
            !v6.to_string().starts_with("fe80")
        } else {
            true
        }
    }));
```

リンクローカル IPv6 は同一セグメント内でしか到達できず、ルーティングされない。
そのため、LAN 内であっても扱いが面倒になることが多い。

> IPv4 があれば IPv4、なければグローバル IPv6

という単純なルールで十分だった。

---

## TXT レコードの活用

mDNS の TXT レコードには、デバイス固有の情報が含まれることが多い。
ただし、これらのキーは標準化されていない。

よく見られるキーの例：

| キー | 意味 |
|------|------|
| `fn`, `name` | フレンドリー名 |
| `md`, `model` | モデル名 |
| `id`, `uuid` | 一意識別子 |

```rust
let friendly_name = properties
    .iter()
    .find(|p| p.key() == "fn" || p.key() == "name")
    .map(|p| p.val_str().to_string());
```

実装では、複数の候補キーを試す形にしている。
存在しなければ `None` のまま持ち、無理に補完しない。

> TXT レコードはあくまで「ヒント」

と考えて扱うのが安全だ。

---

## 複合条件による判定

単一条件での判定は、ほぼ確実に誤検出を生む。
実装では以下を複合的に評価した。

```rust
let is_target =
    name.to_lowercase().contains("keyword")
    || hostname.to_lowercase().contains("keyword")
    || properties.iter().any(|p| {
        p.key().to_lowercase().contains("keyword")
        || p.val_str().to_lowercase().contains("keyword")
    });

// さらに service type も確認
if is_target || info.get_type() == SERVICE_TYPE {
    // 候補として採用
}
```

フルネーム、ホスト名、TXT レコード、service type のいずれかにマッチすれば候補とする。

一見すると緩い条件に見えるが、service type で browse している時点である程度絞られているため、これで十分機能する。

> 厳しすぎる条件は、正当なデバイスを弾くリスクがある

このバランスは、実際に動かしながら調整した。

---

## 実装の一部（簡略化）

以下は、mDNS レコードを内部表現に変換する部分の一例。
再現性を下げるため、意図的に簡略化している。

```rust
// mDNS レスポンスを抽象化した構造体
struct MdnsRecord {
    addr: Option<IpAddr>,
    hostname: Option<String>,
    txt: Vec<(String, String)>,
}

// discovery 段階で扱う最小限の内部表現
struct DiscoveredDevice {
    addr: IpAddr,
    hostname: Option<String>,
    model: Option<String>,
}

fn normalize(record: MdnsRecord) -> Option<DiscoveredDevice> {
    // IP アドレスが取得できないものは候補から除外する
    let addr = record.addr?;

    let hostname = record.hostname;

    // TXT レコードなどからモデル情報らしきものを抽出する
    let model = extract_model_hint(&record.txt);

    Some(DiscoveredDevice {
        addr,
        hostname,
        model,
    })
}
```

![](/images/zeroconf_mdns_discovery-4.jpg)

ここで意識している点は：

- 情報が欠けていれば `None` のまま持つ
- 無理に補完しない
- この段階では「確定」させない

Zeroconf 的には、確実でないものを確定情報にしないことが重要だ。

---

## 重複レスポンスの扱い

mDNS は UDP マルチキャストのため、同一デバイスから複数回レスポンスが届くことがある。  
これは異常ではなく、仕様上想定される動作だ。

実装では IP アドレスをキーとした HashMap で管理し、後から届いた情報で上書きする形で対応した。

```rust
let devices: HashMap<String, Device> = HashMap::new();

// イベント受信ループ内
devices.insert(ip, device);  // 同じ IP なら上書き
```

これにより、重複は自然に排除され、常に最新の情報が保持される。

---

## 非同期での discovery ループ

mDNS の受信は継続的に行いつつ、一定時間で打ち切る必要がある。

Rust + tokio の場合、`tokio::select!` で実現できる。

```rust
let deadline = Instant::now() + timeout;

loop {
    tokio::select! {
        _ = sleep_until(deadline) => {
            break;  // タイムアウト
        }
        event = receive_mdns_event() => {
            // イベント処理
        }
    }
}
```

`select!` は複数の Future を競合させ、最初に完了したものを処理する。

これにより、

- イベントが来れば即座に処理
- タイムアウトに達すればループを抜ける

という動作が自然に書ける。

なお、mDNS ライブラリが同期的な場合は、`spawn_blocking` でブロッキング操作を別スレッドに逃がす必要がある。

```rust
let event = tokio::task::spawn_blocking(move || {
    receiver.recv_timeout(Duration::from_millis(100))
}).await;
```

非同期と同期の境界は、discovery のような I/O 処理では特に意識する必要がある。

---

## タイムアウト設計は UX の話

discovery はネットワーク I/O なので、必ず「待ち」が発生する。

- 長く待てば見つかる可能性は上がる
- しかし CLI としての体験は悪くなる

最終的には、

- 数秒で打ち切る
- 見つからなければ「存在しない」と判断する

という割り切りにした。

> 必ず見つけるより、すぐ終わる

という判断も、Zeroconf を扱う上では重要だと感じている。

---

## おわりに

Zeroconf / mDNS を使った LAN 内 discovery は、「とりあえず動いた」で終わらせるには惜しい題材だ。

- どこまで信じるか
- どこで諦めるか
- 何を確定情報として扱うか

これらを考える過程そのものが、ネットワークソフトウェア設計の練習になる。

もし mDNS を使う機会があれば、「見つかった」その先を一段考えてみると、意外と奥が深くて面白いかもしれない。
