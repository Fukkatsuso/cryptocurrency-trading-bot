# trading

リソース制約と価格変動から取引戦略を考える

- 現物の成行注文
- ETH/JPYの現物なら0.01が最小(https://bitflyer.com/ja-jp/faq/4-27)
- 0.01ずつ，売り買いを交互に行う
- 元手は5000円
- 直近30日の取引量が10万円未満なら手数料0.15%．そんなに痛くない
- 1日1〜2回の取引
- チャートや取引履歴をチェックしやすい時間帯が良い
- 9:00/21:00の12時間周期か，どちらかの時間で1日周期で取引を行うことにする
- とりあえず9:00，1日1回取引する
