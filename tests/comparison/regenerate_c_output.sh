#!/bin/bash
# C版/Go版EESLISMで出力を再生成するスクリプト
#
# 使用法:
#   ./regenerate_c_output.sh              # 全テストのc_outputを再生成
#   ./regenerate_c_output.sh --both       # 全テストのc_output + go_outputを再生成
#   ./regenerate_c_output.sh --both --all-variants
#                                         # 全テスト + 全バリアントを再生成
#   ./regenerate_c_output.sh L2_equipment/vav        # 特定テストのc_outputのみ
#   ./regenerate_c_output.sh --both L2_equipment/vav # 特定テストの両方を再生成
#   ./regenerate_c_output.sh --both --file pipevptr_test.txt L2_equipment/pump_pipe
#                                         # 特定ファイルをvariants/に再生成
#   ./regenerate_c_output.sh --both --all-variants L2_equipment/pump_pipe
#                                         # ディレクトリ内の全テストファイルを再生成

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TESTDATA_DIR="$SCRIPT_DIR/testdata"
PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"
BASE_DIR="$PROJECT_ROOT/Base"
C_EESLISM="/home/uda/ws/open_eeslism/build/eeslism"
GO_EESLISM="$PROJECT_ROOT/eeslism-go"

REGENERATE_GO=false
SPECIFIC_FILE=""
ALL_VARIANTS=false

# オプション解析
while [[ $# -gt 0 ]]; do
    case $1 in
        --both)
            REGENERATE_GO=true
            shift
            ;;
        --file)
            SPECIFIC_FILE="$2"
            shift 2
            ;;
        --all-variants)
            ALL_VARIANTS=true
            shift
            ;;
        *)
            TARGET="$1"
            shift
            ;;
    esac
done

# C版実行ファイルの確認
if [ ! -f "$C_EESLISM" ]; then
    echo "ERROR: C executable not found at $C_EESLISM"
    echo "Build it with:"
    echo "  cd /home/uda/ws/open_eeslism"
    echo "  cmake -B build -DCMAKE_C_FLAGS=\"-fcommon\""
    echo "  cd build && make"
    exit 1
fi

# Go版のビルド（--both指定時）
if [ "$REGENERATE_GO" = true ]; then
    echo "Building Go executable..."
    cd "$PROJECT_ROOT"
    go build -o "$GO_EESLISM" .
    echo "Done."
    echo ""
fi

# 単一テストファイルを再生成
# 引数: test_dir, txt_file, output_subdir (variants/xxx または 空)
regenerate_single_test() {
    local test_dir="$1"
    local txt_file="$2"
    local output_subdir="$3"
    local name=$(basename "$test_dir")
    local parent=$(basename "$(dirname "$test_dir")")

    cd "$test_dir"

    # 出力ディレクトリを決定
    local c_out="c_output"
    local go_out="go_output"
    if [ -n "$output_subdir" ]; then
        c_out="$output_subdir/c_output"
        go_out="$output_subdir/go_output"
    fi

    echo "Regenerating: $parent/$name ($txt_file) -> $c_out"

    # === C版の実行 ===
    # Baseファイルをコピー（C版は--eflオプションがないため）
    cp "$BASE_DIR"/* . 2>/dev/null || true

    # 既存の.esファイルを削除
    rm -f *.es *.gchi 2>/dev/null || true

    # C版を実行
    echo "" | timeout 120 "$C_EESLISM" "$txt_file" > /dev/null 2>&1 || true
    local exit_code=$?

    if [ $exit_code -eq 139 ] || [ $exit_code -eq 11 ]; then
        echo "  WARNING: C version crashed (segfault)"
    elif [ $exit_code -eq 124 ]; then
        echo "  WARNING: C version timed out"
    fi

    # c_outputディレクトリを作成/クリア
    mkdir -p "$c_out"
    rm -f "$c_out"/*.es 2>/dev/null || true

    # 生成された.esファイルをc_outputに移動
    mv *.es "$c_out/" 2>/dev/null || true

    local c_count=$(ls "$c_out"/*.es 2>/dev/null | wc -l)
    echo "  C version: $c_count files"

    # === Go版の実行（--both指定時） ===
    if [ "$REGENERATE_GO" = true ]; then
        # 既存の.esファイルを削除
        rm -f *.es *.gchi 2>/dev/null || true

        # Go版を実行（--eflオプションでBaseディレクトリを指定）
        "$GO_EESLISM" "$txt_file" --efl "$BASE_DIR" > /dev/null 2>&1 || true

        # go_outputディレクトリを作成/クリア
        mkdir -p "$go_out"
        rm -f "$go_out"/*.es 2>/dev/null || true

        # 生成された.esファイルをgo_outputに移動
        mv *.es "$go_out/" 2>/dev/null || true

        local go_count=$(ls "$go_out"/*.es 2>/dev/null | wc -l)
        echo "  Go version: $go_count files"
    fi

    # Baseファイルと一時ファイルを削除
    rm -f supw.efl wbmlist.efl pumpfanlst.efl reflist.efl windowcat.efl dayweek.efl *.has *.gchi 2>/dev/null || true
    rm -f sample.txt radiant_floor_heating.txt 2>/dev/null || true
    rm -f *.log *bdata.ewk *bdata0.ewk *schnma.ewk *schtba.ewk *week.ewk 2>/dev/null || true
}

# ディレクトリのメインテストを再生成
regenerate_main_test() {
    local test_dir="$1"
    local name=$(basename "$test_dir")
    local parent=$(basename "$(dirname "$test_dir")")

    cd "$test_dir"

    # 最初にBaseからコピーされた残骸を削除
    rm -f supw.efl wbmlist.efl pumpfanlst.efl reflist.efl windowcat.efl dayweek.efl *.has 2>/dev/null || true
    rm -f sample.txt radiant_floor_heating.txt 2>/dev/null || true
    rm -f *.log *.gchi 2>/dev/null || true
    rm -f *bdata.ewk *bdata0.ewk *schnma.ewk *schtba.ewk *week.ewk 2>/dev/null || true

    # テスト入力ファイルを探す（優先順位: {name}_test.txt > {name}.txt > 最初の.txt）
    local txt_file=""
    if [ -f "${name}_test.txt" ]; then
        txt_file="${name}_test.txt"
    elif [ -f "${name}.txt" ]; then
        txt_file="${name}.txt"
    else
        txt_file=$(ls *.txt 2>/dev/null | head -1)
    fi

    if [ -z "$txt_file" ]; then
        echo "SKIP: $parent/$name (no .txt file)"
        return
    fi

    regenerate_single_test "$test_dir" "$txt_file" ""
}

# ディレクトリの全バリアントテストを再生成
regenerate_all_variants() {
    local test_dir="$1"
    local name=$(basename "$test_dir")
    local parent=$(basename "$(dirname "$test_dir")")

    cd "$test_dir"

    # 最初にBaseからコピーされた残骸を削除
    rm -f supw.efl wbmlist.efl pumpfanlst.efl reflist.efl windowcat.efl dayweek.efl *.has 2>/dev/null || true
    rm -f sample.txt radiant_floor_heating.txt 2>/dev/null || true
    rm -f *.log *.gchi 2>/dev/null || true

    # メインテストファイルを特定
    local main_file=""
    if [ -f "${name}_test.txt" ]; then
        main_file="${name}_test.txt"
    elif [ -f "${name}.txt" ]; then
        main_file="${name}.txt"
    fi

    # 全.txtファイルを処理
    for txt_file in *.txt; do
        [ -f "$txt_file" ] || continue

        if [ "$txt_file" = "$main_file" ]; then
            # メインテスト -> c_output/go_output
            regenerate_single_test "$test_dir" "$txt_file" ""
        else
            # バリアント -> variants/{basename}/c_output
            local variant_name="${txt_file%.txt}"
            regenerate_single_test "$test_dir" "$txt_file" "variants/$variant_name"
        fi
    done
}

# メイン処理
if [ -n "$TARGET" ]; then
    target_dir="$TESTDATA_DIR/$TARGET"
    if [ ! -d "$target_dir" ]; then
        echo "ERROR: Test directory not found: $target_dir"
        exit 1
    fi

    if [ -n "$SPECIFIC_FILE" ]; then
        # 特定ファイルを指定 -> variants/に出力
        if [ ! -f "$target_dir/$SPECIFIC_FILE" ]; then
            echo "ERROR: Test file not found: $target_dir/$SPECIFIC_FILE"
            exit 1
        fi
        variant_name="${SPECIFIC_FILE%.txt}"
        regenerate_single_test "$target_dir" "$SPECIFIC_FILE" "variants/$variant_name"
    elif [ "$ALL_VARIANTS" = true ]; then
        # 全バリアントを再生成
        regenerate_all_variants "$target_dir"
    else
        # メインテストのみ
        regenerate_main_test "$target_dir"
    fi
else
    # 全テストを再生成
    if [ "$ALL_VARIANTS" = true ]; then
        echo "=== Regenerating all tests with variants ==="
    elif [ "$REGENERATE_GO" = true ]; then
        echo "=== Regenerating all c_output and go_output directories ==="
    else
        echo "=== Regenerating all c_output directories ==="
    fi
    echo ""

    for level_dir in "$TESTDATA_DIR"/L*; do
        if [ -d "$level_dir" ]; then
            echo "--- $(basename "$level_dir") ---"
            for test_dir in "$level_dir"/*/; do
                if [ -d "$test_dir" ]; then
                    if [ "$ALL_VARIANTS" = true ]; then
                        regenerate_all_variants "$test_dir"
                    else
                        regenerate_main_test "$test_dir"
                    fi
                fi
            done
            echo ""
        fi
    done

    echo "=== Done ==="
fi
