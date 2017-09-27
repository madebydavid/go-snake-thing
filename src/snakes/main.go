package main

import (
    tm "github.com/buger/goterm"
    "github.com/go-redis/redis"
    "time"
    "math"
    "math/rand"
    "github.com/golang/protobuf/proto"
    "snakedata"
)

type World struct {
    Width int
    Height int
    Grid [][]rune 
}

func newWorld(width, height int) World {

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

func putSymbolAtPoint(world *World, symbol rune, point snakedata.Snake_Point) {
    column := round(point.X);
    row := round(point.Y);
    (*world).Grid[row][column] = symbol;
}

func round(num float32) int {
    return int(float64(num) + math.Copysign(0.5, float64(num)))
}

func moveSnakes(snakes *[]snakedata.Snake, world *World, client *redis.Client) {

    maxX := float32((*world).Width) - 1
    maxY := float32((*world).Height) - 1

    for i := 0; i < len(*snakes); i++ {

        snake := (*snakes)[i]

        // Clear the rendered tail
        for _, tailPoint := range snake.Tail {
            putSymbolAtPoint(world, 0, *tailPoint);
        }

        // Get the most recent point
        previousTailPoint := *(*snakes)[i].Tail[0]

        // Apply the velocity and create new point
        
        newTailPoint := &snakedata.Snake_Point{
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
        (*snakes)[i].Tail = append([]*snakedata.Snake_Point{newTailPoint} ,(*snakes)[i].Tail...)

        // Trim tail to tail length
        if int32(len((*snakes)[i].Tail)) > (*snakes)[i].Length { 
            (*snakes)[i].Tail = (*snakes)[i].Tail[0:(*snakes)[i].Length]
        }
        
        // Render the tail
        for _, tailPoint := range (*snakes)[i].Tail {
            putSymbolAtPoint(world, snake.Symbol, *tailPoint);
        }

        // Serialise the snake
        serialisedSnake, _ := proto.Marshal(&(*snakes)[i])
        // Persist to redis
        client.LSet("snakes", int64(i), serialisedSnake)


    }
}


func main() {

    client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })

    client.LTrim("snakes", 0, 0)

    rand.Seed(time.Now().UTC().UnixNano())

    world := newWorld(tm.Width() - 2 , tm.Height() - 2)

    snakes := []snakedata.Snake{}

    const snakeCount = 200
    const firstSymbol rune = 'A'

    for i := 0; i < snakeCount; i++ {

        initialTailPoint := &snakedata.Snake_Point{
            X: rand.Float32() * (float32(world.Width) - 1),
            Y: rand.Float32() * (float32(world.Height) - 1),
        }

        newSnake := &snakedata.Snake{
            VelocityX: rand.Float32() / 5,
            VelocityY: rand.Float32() / 5,
            Length: int32(rand.Intn(20) + 1),
            Symbol: firstSymbol + rune(i % 50),
            Tail: []*snakedata.Snake_Point{initialTailPoint},
        }

        snakes = append(snakes, *newSnake)

        // Serialise the snake
        serialisedSnake, _ := proto.Marshal(newSnake)
        // Persist
        client.LPush("snakes", serialisedSnake)
    }

    tm.Clear()

    for {
        moveSnakes(&snakes, &world, client)
        renderWorld(world)
        time.Sleep(time.Second / 100)
    }

}