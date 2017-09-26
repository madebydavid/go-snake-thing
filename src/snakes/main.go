package main

import (
    tm "github.com/buger/goterm"
    "time"
    "math"
    "math/rand"
)

type World struct {
    Width int
    Height int
    Grid [][]rune 
}

type Point struct {
    X float64
    Y float64
}

type Snake struct {
    Point
    VelocityX float64
    VelocityY float64
    Length int
    Symbol rune
    Tail []Point
}

type Renderable interface {
    GetRow() int
    GetColumn() int
}

func (point Point) GetRow() int {
    return round((point).Y)
}

func (point Point) GetColumn() int {
    return round((point).X)
}

func (snake Snake) GetRow() int {
    return round((snake).Y)
}

func (snake Snake) GetColumn() int {
    return round((snake).X)
}

func round(num float64) int {
    return int(num + math.Copysign(0.5, num))
}

func newWorld(width, height int) (World) {

    grid := make([][]rune, height)

    for i := range grid {
        grid[i] = make([]rune, width)
    }

    return World{
        Width: width,
        Height: height,
        Grid: grid,
    }
}

func renderWorld(world World) {
    tm.MoveCursor(1,1)
    
    for row := range world.Grid {
        rowString := ""
        for col := range world.Grid[row] {
            if world.Grid[row][col] == 0 {
                rowString += " "
            } else {
                rowString += string(world.Grid[row][col])
            }
        }
        tm.Println(rowString)
    }

    tm.Flush()
}

func putSymbolAtPoint(world *World, symbol rune, point Renderable) {
    (*world).Grid[point.GetRow()][point.GetColumn()] = symbol;
}

func moveSnakes(snakes *[]Snake, world *World) {

    maxX := float64((*world).Width) - 1
    maxY := float64((*world).Height) - 1

    for i := 0; i < len(*snakes); i++ {

        snake := (*snakes)[i]

        for _, tailPoint := range snake.Tail {
            putSymbolAtPoint(world, 0, tailPoint);
        }

        putSymbolAtPoint(world, 0, snake);

        (*snakes)[i].X += snake.VelocityX

        if (*snakes)[i].X >= maxX || (*snakes)[i].X < 0 {
            (*snakes)[i].VelocityX *= -1
            (*snakes)[i].X += (*snakes)[i].VelocityX
        }

        (*snakes)[i].Y += snake.VelocityY

        if (*snakes)[i].Y >= maxY || (*snakes)[i].Y < 0 {
            (*snakes)[i].VelocityY *= -1
            (*snakes)[i].Y += (*snakes)[i].VelocityY
        }

        newTailPoint := Point{
            X: (*snakes)[i].X,
            Y: (*snakes)[i].Y,
        }

        (*snakes)[i].Tail = append([]Point{newTailPoint} ,(*snakes)[i].Tail...)

        if len((*snakes)[i].Tail) > (*snakes)[i].Length { 
            (*snakes)[i].Tail = (*snakes)[i].Tail[1:(*snakes)[i].Length]
        }

        putSymbolAtPoint(world, snake.Symbol, (*snakes)[i]);

        for _, tailPoint := range (*snakes)[i].Tail {
            putSymbolAtPoint(world, snake.Symbol, tailPoint);
        }

    }
}


func main() {

    rand.Seed(time.Now().UTC().UnixNano())

    world := newWorld(tm.Width() - 2 , tm.Height() - 2)

    snakes := []Snake{}

    const snakeCount = 200
    const firstSymbol rune = 'A'

    for i := 0; i < snakeCount; i++ {
        newSnake := Snake{
            Point: Point{
                X: rand.Float64() * (float64(world.Width) - 1),
                Y: rand.Float64() * (float64(world.Height) - 1),
            },
            VelocityX: rand.Float64() / 5,
            VelocityY: rand.Float64() / 5,
            Length: rand.Intn(20) + 1,
            Symbol: firstSymbol + rune(i % 50),
        }
        snakes = append(snakes, newSnake)
    }

    tm.Clear()

    for {
        moveSnakes(&snakes, &world)
        renderWorld(world)
        time.Sleep(time.Second / 100)
    }

}