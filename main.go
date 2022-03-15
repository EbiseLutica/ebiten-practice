package main

import (
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// セルサイズ
const CELL_SIZE = 8

// 画面幅
const SCREEN_WIDTH = CELL_SIZE * 64

// 画面高さ
const SCREEN_HEIGHT = CELL_SIZE * 64

// セルテーブル列数
const TABLE_COLUMN = SCREEN_WIDTH / CELL_SIZE

// セルテーブル行数
const TABLE_ROW = SCREEN_WIDTH / CELL_SIZE

// テーブルの型定義
type Table [TABLE_COLUMN][TABLE_ROW]bool

// テーブルバッファA
var tableA Table

// テーブルバッファB
var tableB Table

// レンダリングポインタ - Drawで描画中のテーブルのポインタ
var tableR *Table

// 更新ポインタ - 次世代を計算するのに使うテーブルのポインタ
var tableU *Table

// セル画像
var cellImage = ebiten.NewImage(CELL_SIZE, CELL_SIZE)

// ライフゲーム再生中かどうか。falseであれば描画モードとなる
var isPlaying = false

// 現在のフレーム数
var frame = 0

// Boolean to Int
func btoi(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

// 周囲のセルの生死を取得する
func searchAlive(x int, y int) int {
	isTop := y == 0
	isBottom := y == TABLE_ROW-1
	isLeft := x == 0
	isRight := x == TABLE_COLUMN-1

	var count = 0

	if !isTop {
		if !isLeft {
			count += btoi(tableR[x-1][y-1])
		}
		count += btoi(tableR[x][y-1])
		if !isRight {
			count += btoi(tableR[x+1][y-1])
		}
	}
	if !isLeft {
		count += btoi(tableR[x-1][y])
	}
	if !isRight {
		count += btoi(tableR[x+1][y])
	}
	if !isBottom {
		if !isLeft {
			count += btoi(tableR[x-1][y+1])
		}
		count += btoi(tableR[x][y+1])
		if !isRight {
			count += btoi(tableR[x+1][y+1])
		}
	}

	return count
}

// 次世代を計算する
func calculateNext() {
	for x := 0; x < TABLE_COLUMN; x++ {
		for y := 0; y < TABLE_ROW; y++ {
			var isAlive = tableR[x][y]
			var aroundAlive = searchAlive(x, y)

			// カレントセルが死んでいる & 周囲に生存セル数が3であれば、カレントセルの位置に次世代が誕生する
			if !isAlive {
				tableU[x][y] = aroundAlive == 3
			}

			// カレントセルが生きている & 周囲の生存セル数が1以下または4以上であれば、カレントセルは死ぬ
			if isAlive {
				tableU[x][y] = (aroundAlive > 1 && aroundAlive < 4)
			}
		}
	}
	if tableR == &tableA {
		tableR = &tableB
		tableU = &tableA
	} else {
		tableR = &tableA
		tableU = &tableB
	}
}

// 入力を処理する
func processInput() {
	// 再生状態切替
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		isPlaying = !isPlaying
	}

	if isPlaying {
		return
	}

	// 編集モードであれば、セルの編集を受け付ける
	var mx, my = ebiten.CursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		var cx = int(mx / CELL_SIZE)
		var cy = int(my / CELL_SIZE)
		tableR[cx][cy] = !tableR[cx][cy]
	}
}

// Ebiten ゲームロジック
type Game struct{}

// Ebiten - アップデート
func (g *Game) Update() error {
	processInput()
	if isPlaying {
		ebiten.SetWindowTitle("[SPACE]: Stop")
		frame++
		if frame%4 == 0 {
			calculateNext()
		}
	} else {
		ebiten.SetWindowTitle("[SPACE]: Start")
		frame = 0
	}
	return nil
}

// Ebiten - レンダリング
func (g *Game) Draw(screen *ebiten.Image) {
	for x := 0; x < TABLE_COLUMN; x++ {
		for y := 0; y < TABLE_ROW; y++ {
			if tableR[x][y] {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x*CELL_SIZE), float64(y*CELL_SIZE))
				screen.DrawImage(cellImage, op)
			}
		}
	}
}

// Ebiten - 画面サイズ
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

// エントリ ポイント
func main() {
	cellImage.Fill(color.White)
	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)

	tableR = &tableA
	tableU = &tableB

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
