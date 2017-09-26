package main

import (
    tm "github.com/buger/goterm"
    "github.com/go-redis/redis"
    "fmt"
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
    VelocityX float64
    VelocityY float64
    Length int
    Symbol rune
    Tail []Point
}

func (point Point) GetRow() int {
    return round((point).Y)
}

func (point Point) GetColumn() int {
    return round((point).X)
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

func putSymbolAtPoint(world *World, symbol rune, point Point) {
    (*world).Grid[point.GetRow()][point.GetColumn()] = symbol;
}

func moveSnakes(snakes *[]Snake, world *World) {

    maxX := float64((*world).Width) - 1
    maxY := float64((*world).Height) - 1

    for i := 0; i < len(*snakes); i++ {

        snake := (*snakes)[i]

        // Clear the rendered tail
        for _, tailPoint := range snake.Tail {
            putSymbolAtPoint(world, 0, tailPoint);
        }

        // Get the most recent point
        previousTailPoint := (*snakes)[i].Tail[0]

        // Apply the velocity and create new point
        newTailPoint := Point{
            X: previousTailPoint.X + snake.VelocityX,
            Y: previousTailPoint.Y + snake.VelocityY,
        }

        // If on edge bounce
        if newTailPoint.X >= maxX || newTailPoint.X < 0 {
            (*snakes)[i].VelocityX *= -1
            newTailPoint.X += (*snakes)[i].VelocityX
        }

        if newTailPoint.Y >= maxY || newTailPoint.Y < 0 {
            (*snakes)[i].VelocityY *= -1
            newTailPoint.Y += (*snakes)[i].VelocityY
        }

        // Append to begining of tail
        (*snakes)[i].Tail = append([]Point{newTailPoint} ,(*snakes)[i].Tail...)

        // Trim tail to tail length
        if len((*snakes)[i].Tail) > (*snakes)[i].Length { 
            (*snakes)[i].Tail = (*snakes)[i].Tail[0:(*snakes)[i].Length]
        }
        
        // Render the tail
        for _, tailPoint := range (*snakes)[i].Tail {
            putSymbolAtPoint(world, snake.Symbol, tailPoint);
        }

    }
}


func main() {

    client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    pong, err := client.Ping().Result()
    fmt.Println(pong, err)

    rand.Seed(time.Now().UTC().UnixNano())

    world := newWorld(tm.Width() - 2 , tm.Height() - 2)

    snakes := []Snake{}

    const snakeCount = 200
    const firstSymbol rune = 'A'

    for i := 0; i < snakeCount; i++ {

        initialTailPoint := Point{
            X: rand.Float64() * (float64(world.Width) - 1),
            Y: rand.Float64() * (float64(world.Height) - 1),
        }

        newSnake := Snake{
            VelocityX: rand.Float64() / 5,
            VelocityY: rand.Float64() / 5,
            Length: rand.Intn(20) + 1,
            Symbol: firstSymbol + rune(i % 50),
            Tail: []Point{initialTailPoint},
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