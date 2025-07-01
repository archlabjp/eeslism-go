# 8. トラブルシューティング

## 8.1 はじめに

このガイドでは、EESLISM Go版を使用する際によく遭遇する問題とその解決方法を説明します。エラーメッセージの読み方から、計算が正常に実行されない場合の対処法まで、実践的な解決策を提供します。

## 8.2 インストール・実行時の問題

### 8.2.1 プログラムが起動しない

#### 問題: 実行ファイルが見つからない
```bash
$ ./eeslism
bash: ./eeslism: No such file or directory
```

**原因と解決方法:**
1. **ビルドされていない**
   ```bash
   # Go言語でビルド
   go build
   ```

2. **実行権限がない**
   ```bash
   # 実行権限を付与
   chmod +x eeslism
   ```

3. **パスが間違っている**
   ```bash
   # 現在のディレクトリを確認
   ls -la eeslism
   # フルパスで実行
   /path/to/eeslism input.txt
   ```

#### 問題: Go言語がインストールされていない
```bash
$ go build
bash: go: command not found
```

**解決方法:**
1. **Go言語のインストール**
   - [公式サイト](https://go.dev/doc/install)からダウンロード
   - パッケージマネージャーを使用:
     ```bash
     # Ubuntu/Debian
     sudo apt install golang
     
     # macOS (Homebrew)
     brew install go
     
     # Windows (Chocolatey)
     choco install golang
     ```

2. **バージョン確認**
   ```bash
   go version
   # go version go1.23.0 linux/amd64 (推奨: 1.18以上)
   ```

### 8.2.2 依存関係の問題

#### 問題: モジュールが見つからない
```bash
$ go build
go: module github.com/archlabjp/eeslism-go: cannot find module providing package
```

**解決方法:**
```bash
# 依存関係を整理
go mod tidy

# 依存関係をダウンロード
go mod download

# 再ビルド
go build
```

## 8.3 入力ファイルの問題

### 8.3.1 ファイルが見つからない

#### 問題: 入力ファイルが見つからない
```
Error: Cannot open input file: input.txt
```

**解決方法:**
1. **ファイル存在確認**
   ```bash
   ls -la input.txt
   ```

2. **パス指定**
   ```bash
   # 相対パス
   ./eeslism samples/simple_room.txt
   
   # 絶対パス
   ./eeslism /full/path/to/input.txt
   ```

#### 問題: 気象データファイルが見つからない
```
Error: Cannot open weather file: tokyo_3column_SI.has
```

**解決方法:**
1. **EFLディレクトリの指定**
   ```bash
   ./eeslism input.txt --efl Base
   ```

2. **ファイル存在確認**
   ```bash
   ls Base/tokyo_3column_SI.has
   ```

3. **入力ファイルでのパス修正**
   ```
   GDAT
       FILE    w=Base/tokyo_3column_SI.has ;
   ```

### 8.3.2 構文エラー

#### 問題: セクション終了記号の不備
```
Error: Syntax error in line 25: Expected '*' to close section
```

**解決方法:**
各セクションの最後に `*` を追加:
```
WALL
    exterior_wall
    3
    concrete    0.15
    insulation  0.10
    gypsum      0.012
    ;
*  ← この行が必要
```

#### 問題: パラメータの不備
```
Error: Invalid parameter in ROOM definition: line 45
```

**解決方法:**
1. **必須パラメータの確認**
   ```
   ROOM
       RoomName
       Rg=1.0 Cg=1000000 Ag=50.0  ← 必須パラメータ
       ...
   ```

2. **数値フォーマットの確認**
   ```
   # 正しい例
   Rg=1.0
   
   # 間違った例
   Rg=1,0  ← カンマではなくピリオド
   ```

### 8.3.3 データ定義の問題

#### 問題: 未定義の材料・機器
```
Error: Undefined wall material: unknown_material
```

**解決方法:**
1. **材料定義の確認**
   ```
   WALL
       wall_name
       2
       concrete    0.15  ← wbmlist.eflで定義済みの材料名
       insulation  0.10
       ;
   ```

2. **カスタム材料の定義**
   ```
   WALL
       wall_name
       2
       concrete    0.15
       my_material 0.10 cond=0.04 dens=50 spht=1000
       ;
   ```

## 8.4 計算実行時の問題

### 8.4.1 収束しない

#### 問題: 計算が収束しない
```
Warning: Calculation did not converge after 100 iterations
```

**解決方法:**
1. **最大反復回数の増加**
   ```
   GDAT
       RUN 1/1-12/31 MaxIterate=200 dTime=1800 ;
   ```

2. **時間間隔の短縮**
   ```
   GDAT
       RUN 1/1-12/31 dTime=900 ;  ← 30分→15分
   ```

3. **PCM収束計算の調整**
   ```
   PCM
       PCM_name
       ...
       nWeight=0.5     ← 重み係数を大きく（0.1→0.5）
       -iterate
       ;
   ```

#### 問題: VAVシステムで収束しない
```
Warning: VAV system did not converge
```

**解決方法:**
1. **風量範囲の調整**
   ```
   EQPCAT
       VAV
       VAV_name A 0.5 0.05  ← 最小風量を大きく
   ```

2. **制御設定の見直し**
   ```
   CONTL
       VAV_name
       %v -t SetTemp Summer 26.0 ;  ← 設定温度を現実的に
   ```

### 8.4.2 異常な計算結果

#### 問題: 室温が異常に高い/低い
```
Room temperature: 150.5°C  ← 異常値
```

**原因と解決方法:**
1. **熱容量の確認**
   ```
   ROOM
       RoomName
       Rg=100.0 Cg=10000000 Ag=50.0  ← Cgが小さすぎる場合は増加
   ```

2. **壁体構成の確認**
   ```
   WALL
       wall_name
       3
       concrete    0.15  ← 厚さが現実的か確認
       insulation  0.10
       gypsum      0.012
       ;
   ```

3. **境界条件の確認**
   ```
   VCFILE
       Schedule
       %v -t RoomTemp AllDay 20.0 ;  ← 設定温度が現実的か
   ```

#### 問題: 負荷が異常に大きい
```
Heating load: 50000W  ← 異常に大きい
```

**原因と解決方法:**
1. **断熱性能の確認**
   ```
   # wbmlist.eflで断熱材の熱伝導率を確認
   insulation  dens=50  cond=0.04  ← 適切な値か確認
   ```

2. **窓面積の確認**
   ```
   ROOM
       RoomName
       ...
       south exterior_wall 10.0 single_glass 5.0 south
       #                                    ↑ 窓面積が適切か
   ```

## 8.5 出力ファイルの問題

### 8.5.1 出力ファイルが生成されない

#### 問題: 出力ファイルが作成されない
**解決方法:**
1. **出力設定の確認**
   ```
   GDAT
       PRINT 1/1-12/31 *wd ;  ← 出力期間と項目を確認
   ```

2. **ディスク容量の確認**
   ```bash
   df -h .  # 現在のディスクの空き容量を確認
   ```

3. **書き込み権限の確認**
   ```bash
   ls -la .  # 現在のディレクトリの権限を確認
   ```

### 8.5.2 出力データが異常

#### 問題: 出力データが空または異常
**解決方法:**
1. **計算期間の確認**
   ```
   GDAT
       RUN 1/1-1/7 dTime=3600 ;     ← 実行期間
       PRINT 1/1-1/7 *wd ;          ← 出力期間（一致させる）
   ```

2. **出力項目の確認**
   ```
   GDAT
       PRINT 1/1-1/7 rm sf wd ;  ← 必要な項目のみ指定
   ```

## 8.6 性能・メモリの問題

### 8.6.1 実行速度が遅い

#### 問題: 計算に時間がかかりすぎる
**解決方法:**
1. **時間間隔の調整**
   ```
   GDAT
       RUN 1/1-12/31 dTime=3600 ;  ← 30分→1時間
   ```

2. **出力項目の削減**
   ```
   GDAT
       PRINT 1/1-12/31 rm ;  ← 必要最小限の出力
   ```

3. **計算期間の短縮**
   ```
   GDAT
       RUN 1/1-1/31 dTime=1800 ;  ← テスト用に短期間
   ```

### 8.6.2 メモリ不足

#### 問題: メモリ不足エラー
```
Error: Out of memory
```

**解決方法:**
1. **計算条件の調整**
   ```
   GDAT
       RUN 1/1-3/31 dTime=3600 ;  ← 期間短縮
   ```

2. **出力の最適化**
   ```
   GDAT
       PRINT 1/1-3/31 rm ;  ← 出力項目削減
   ```

3. **システムメモリの確認**
   ```bash
   free -h  # Linux
   # または
   top      # 実行中のメモリ使用量確認
   ```

## 8.7 WebAssembly版特有の問題

### 8.7.1 ブラウザでの実行問題

#### 問題: WebAssembly版が動作しない
**解決方法:**
1. **ブラウザの対応確認**
   - Chrome 57+, Firefox 52+, Safari 11+, Edge 16+

2. **ファイルアクセスの制限**
   - ローカルファイルの直接アクセスは制限される
   - Webサーバー経由でのアクセスが必要

3. **メモリ制限**
   - ブラウザのメモリ制限により大規模計算は困難
   - 計算規模を小さくする

## 8.8 デバッグのヒント

### 8.8.1 段階的なデバッグ

1. **最小構成から開始**
   ```
   # 1室、1日、1時間間隔から開始
   GDAT
       RUN 1/1-1/1 dTime=3600 ;
   ```

2. **徐々に複雑化**
   ```
   # 期間延長 → 時間間隔短縮 → 室数増加 → 設備追加
   ```

3. **ログの活用**
   ```bash
   # 詳細ログの出力
   ./eeslism input.txt 2>&1 | tee debug.log
   ```

### 8.8.2 入力ファイルの検証

1. **構文チェック**
   ```bash
   # セクション終了記号の確認
   grep -n "\*" input.txt
   ```

2. **文字コードの確認**
   ```bash
   # UTF-8であることを確認
   file input.txt
   ```

3. **改行コードの確認**
   ```bash
   # Unix形式（LF）であることを確認
   od -c input.txt | head
   ```

## 8.9 よくある質問（FAQ）

### Q1: 計算結果がC言語版と異なる
**A:** Go言語版はC言語版と同じアルゴリズムを使用していますが、浮動小数点演算の微小な差により、わずかな違いが生じる場合があります。実用上問題のない範囲です。

### Q2: 大規模モデルの計算時間を短縮したい
**A:** 以下の方法を試してください：
- 時間間隔を大きくする（30分→1時間）
- 出力項目を必要最小限にする
- 並列計算（将来実装予定）

### Q3: 独自の機器モデルを追加したい
**A:** 現在はソースコード修正が必要です。将来的にはプラグイン機能の実装を検討しています。

### Q4: 商用利用は可能か
**A:** GPL-2.0ライセンスに従って利用可能です。詳細はLICENSEファイルを確認してください。

## 8.10 サポート・コミュニティ

### 8.10.1 問題報告

1. **GitHub Issues**
   - https://github.com/archlabjp/eeslism-go/issues
   - バグレポート、機能要求

2. **情報提供時の注意**
   - OS・Go言語バージョン
   - 入力ファイル（可能な範囲で）
   - エラーメッセージ
   - 実行環境

### 8.10.2 コミュニティ

1. **GitHub Discussions**
   - 使用方法の質問
   - 情報交換

2. **ドキュメント改善**
   - Pull Requestによる文書改善
   - 翻訳・多言語化

このトラブルシューティングガイドを参考に、問題の解決を図ってください。解決しない場合は、コミュニティに相談することをお勧めします。