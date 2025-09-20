package main

import (
	"fmt"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	// Screen width
	screenWidth = 1600
	// Screen height
	screenHeight = 1200

	// Ship Dimensions In Pixels
	shipHeight = 96
	shipWidth  = 64

	maxAsteroids = 50
	maxBullets   = 100
)

type AsteroidSize int32

const (
	Big    AsteroidSize = 0
	Medium AsteroidSize = 1
	Small  AsteroidSize = 2
)

type Player struct {
	ShipRec           rl.Rectangle
	ShipDirection     rl.Vector2
	MovementDirection rl.Vector2
	Position          rl.Vector2
	Velocity          float32
	Tex               rl.Texture2D
	Hit               bool
}

type Asteroid struct {
	Rec              rl.Rectangle
	StartPos         rl.Vector2
	Pos              rl.Vector2
	Velocity         rl.Vector2
	Rotation         float32
	RotationVelocity float32
	Size             AsteroidSize
	Tex              rl.Texture2D
	Active           bool
}

type Bullet struct {
	Pos       rl.Vector2
	Direction rl.Vector2
	Velocity  float32
	Active    bool
}

type Particle struct {
	Pos       rl.Vector2
	Direction rl.Vector2
	Velocity  float32
	Active    bool
}

// Generates a single asteroid.
func (g *Game) Asteroid() (a Asteroid) {
	// Give Asteroid a random size
	a.Size = AsteroidSize(rand.Intn(3))

	// Choose Random texture based on size.
	choice := rand.Intn(3)
	switch a.Size {
	case Big:
		switch choice {
		case 0:
			a.Tex = g.Textures.BigAsteroidTex1
		case 1:
			a.Tex = g.Textures.BigAsteroidTex2
		case 2:
			a.Tex = g.Textures.BigAsteroidTex3
		default:
			break
		}
		a.Rec.Height = 160
		a.Rec.Width = 160

	case Medium:
		switch choice {
		case 0:
			a.Tex = g.Textures.MediumAsteroidTex1
		case 1:
			a.Tex = g.Textures.MediumAsteroidTex2
		case 2:
			a.Tex = g.Textures.MediumAsteroidTex3
		default:
			break
		}
		a.Rec.Height = 96
		a.Rec.Width = 96

	case Small:
		switch choice {
		case 0:
			a.Tex = g.Textures.SmallAsteroidTex1
		case 1:
			a.Tex = g.Textures.SmallAsteroidTex2
		case 2:
			a.Tex = g.Textures.SmallAsteroidTex3
		default:
			break
		}
		a.Rec.Height = 64
		a.Rec.Width = 64
	}

	// Generate random posistion off screen. I.e. x < 0 or x > 1600 && y < 0 or y > 1000
	x := rand.Intn(screenWidth+1000) - 500
	var y int
	if x >= -100 && x <= 1700 { // y must be out of screen pixel range. y < 0 and y > screenHeight + buffer
		choice1 := rand.Intn(500) * -1           // Spawns above screen
		choice2 := rand.Intn(500) + screenHeight // Spawns below screen
		choices := [2]int{choice1, choice2}
		choice := rand.Intn(2) // 0 or 1 for a choice
		y = choices[choice]

	} else { // y can be anywhere
		y = rand.Intn(screenWidth+1000) - 500
	}

	a.StartPos.X = float32(x)
	a.StartPos.Y = float32(y)
	a.Pos.X = float32(x)
	a.Pos.Y = float32(y)
	a.Rec.X = 0
	a.Rec.Y = 0

	// Generate random Velocity that passes through Screen Space
	p := rl.NewVector2(float32(rand.Intn(screenWidth)), float32(rand.Intn(screenHeight))) // random point in screenspace
	a.Velocity = rl.Vector2Scale(rl.Vector2Normalize(rl.Vector2Subtract(p, a.Pos)), .5+rand.Float32()*2.0)

	// Generate Random Rotation and Rotation Velocity
	a.Rotation = float32(rand.Intn(361))       // random degree
	a.RotationVelocity = float32(rand.Intn(5)) // random rotational speed

	a.Active = true

	return a
}

// Creates two Asteroids in place of a given Asteroid
func (g *Game) BreakAsteroid(a *Asteroid) (a1 Asteroid, a2 Asteroid) {
	choiceA1 := rand.Intn(3)
	choiceA2 := rand.Intn(3)
	switch a.Size {
	case Big:
		rl.PlaySound(g.Sounds.LargeBang)
		switch choiceA1 {
		case 0:
			a1.Tex = g.Textures.MediumAsteroidTex1
		case 1:
			a1.Tex = g.Textures.MediumAsteroidTex2
		case 2:
			a1.Tex = g.Textures.MediumAsteroidTex3
		default:
			break
		}
		switch choiceA2 {
		case 0:
			a2.Tex = g.Textures.MediumAsteroidTex1
		case 1:
			a2.Tex = g.Textures.MediumAsteroidTex2
		case 2:
			a2.Tex = g.Textures.MediumAsteroidTex3
		default:
			break
		}
		a1.Size = Medium
		a1.Rec.Height = 96
		a1.Rec.Width = 96
		a2.Size = Medium
		a2.Rec.Height = 96
		a2.Rec.Width = 96

	case Medium:
		rl.PlaySound(g.Sounds.MediumBang)
		switch choiceA1 {
		case 0:
			a1.Tex = g.Textures.SmallAsteroidTex1
		case 1:
			a1.Tex = g.Textures.SmallAsteroidTex2
		case 2:
			a1.Tex = g.Textures.SmallAsteroidTex3
		default:
			break
		}
		switch choiceA2 {
		case 0:
			a2.Tex = g.Textures.SmallAsteroidTex1
		case 1:
			a2.Tex = g.Textures.SmallAsteroidTex2
		case 2:
			a2.Tex = g.Textures.SmallAsteroidTex3
		default:
			break
		}
		a1.Size = Small
		a1.Rec.Height = 64
		a1.Rec.Width = 64
		a2.Size = Small
		a2.Rec.Height = 64
		a2.Rec.Width = 64
	}

	a1.Pos = a.Pos
	a2.Pos = a.Pos
	a1.StartPos = a.Pos
	a2.StartPos = a.Pos
	a1.Rec.X = 0
	a1.Rec.Y = 0
	a2.Rec.X = 0
	a2.Rec.Y = 0

	a1.Rotation = float32(rand.Intn(361))       // random degree
	a1.RotationVelocity = float32(rand.Intn(5)) // random rotational speed
	a2.Rotation = float32(rand.Intn(361))       // random degree
	a2.RotationVelocity = float32(rand.Intn(5)) // random rotational speed

	a1.Velocity = rl.Vector2Scale(rl.Vector2Normalize(rl.Vector2Add(a.Velocity, rl.NewVector2(rand.Float32(), rand.Float32()))), .5+rand.Float32()*2.0)
	a2.Velocity = rl.Vector2Negate(a1.Velocity)

	a1.Active = true
	a2.Active = true

	return a1, a2
}

// Generate a bullet
func (g *Game) Bullet() (b Bullet) {
	b.Direction = g.Player.ShipDirection
	angle := math.Atan2(float64(g.Player.ShipDirection.Y), float64(g.Player.ShipDirection.X))
	b.Pos.X = g.Player.Position.X + (float32(math.Cos(angle)) * (float32(g.Textures.ShipTex1.Width) / 2))
	b.Pos.Y = g.Player.Position.Y + (float32(math.Sin(angle))*(float32(g.Textures.ShipTex1.Height)/2) + 3)
	b.Velocity = g.Player.Velocity + 7
	b.Active = true
	rl.PlaySound(g.Sounds.Fire)
	return b
}

func (g *Game) Particle(pos rl.Vector2) (p Particle) {
	p.Direction = rl.Vector2Normalize(rl.Vector2Rotate(rl.NewVector2(0, 1), rand.Float32()*361.0))
	p.Velocity = 0.5 + rand.Float32()*5
	p.Pos = pos
	p.Active = true
	return p
}

// Update Asteroid Positions
func (g *Game) updateAsteroids() {
	for i := 0; i < g.maxAsteroids; i++ {
		asteroid := &g.Asteroids[i]
		if asteroid.Active && rl.Vector2Distance(asteroid.Pos, asteroid.StartPos) < 5000 {
			asteroid.Pos = rl.Vector2Add(asteroid.Pos, asteroid.Velocity)
			asteroid.Rotation = asteroid.Rotation + asteroid.RotationVelocity
		} else {
			g.Asteroids[i] = g.Asteroid()
		}
	}
}

func (g *Game) updatePlayer() {
	if !g.Player.Hit {
		g.Player.Tex = g.Textures.ShipTex1

		if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
			g.Player.Velocity = g.Player.Velocity + (1 / 8.0)
			if g.FramesCounter%15 == 0 {
				rl.PlaySound(g.Sounds.Thrust)
			}
			if g.Player.Velocity > 8 {
				g.Player.Velocity = 8
			}
			g.Player.Tex = g.Textures.ShipTex2
			g.Player.MovementDirection = rl.Vector2Normalize(rl.Vector2Add(g.Player.ShipDirection, g.Player.MovementDirection))
		} else {
			g.Player.Velocity = g.Player.Velocity - (1.0 / 14.0)
			if g.Player.Velocity < 0 {
				g.Player.Velocity = 0
			}
		}
		if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
			g.Player.ShipDirection = rl.Vector2Normalize(rl.Vector2Rotate(g.Player.ShipDirection, -0.09))
		}
		if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
			g.Player.ShipDirection = rl.Vector2Normalize(rl.Vector2Rotate(g.Player.ShipDirection, 0.09))
		}
		if rl.IsKeyPressed(rl.KeySpace) {
			g.Bullets[g.numBullets] = g.Bullet()
			g.numBullets = (g.numBullets + 1) % g.maxBullets
		}
		g.Player.Position = rl.Vector2Add(g.Player.Position, rl.Vector2Scale(g.Player.MovementDirection, g.Player.Velocity))

		if g.Player.Position.X < -400 {
			g.Player.Position.X = 1800
		} else if g.Player.Position.X > 2000 {
			g.Player.Position.X = -200
		}
		if g.Player.Position.Y < -200 {
			g.Player.Position.Y = 1200
		} else if g.Player.Position.Y > 1200 {
			g.Player.Position.Y = -200
		}

	} else {
		g.Player.Position = rl.NewVector2(screenWidth/2, screenHeight/2)
		g.Lives = g.Lives - 1
		g.Player.Hit = false
	}

}

func (g *Game) updateBullets() {
	for i := 0; i < g.maxBullets; i++ {
		b := &g.Bullets[i]
		if b.Active {
			g.Bullets[i].Pos = rl.Vector2Add(b.Pos, rl.Vector2Scale(b.Direction, b.Velocity))
		}
	}

}

func (g *Game) updateParticles() {
	for i := 0; i < g.maxParticles; i++ {
		p := &g.Particles[i]
		if p.Active {
			g.Particles[i].Pos = rl.Vector2Add(p.Pos, rl.Vector2Scale(p.Direction, p.Velocity))
			p.Velocity = p.Velocity - .1
			if p.Velocity < 0 {
				p.Active = false
			}
		}
	}
}

func (g *Game) BulletCollisions() {
	for i := 0; i < g.maxBullets; i++ {
		b := &g.Bullets[i]
		if !b.Active {
			continue
		}
		for j := 0; j < g.maxAsteroids; j++ {
			a := &g.Asteroids[j]
			if !a.Active {
				continue
			}
			//rl.DrawCircleLines(int32(b.Pos.X), int32(b.Pos.Y), 5, rl.Blue)
			//rl.DrawCircleLines(int32(a.Pos.X), int32(a.Pos.Y), a.Rec.Width/3, rl.Blue)
			if rl.CheckCollisionCircles(b.Pos, 5, a.Pos, a.Rec.Width/3) {
				g.Score += 10
				b.Active = false
				numNewParticles := 1 + rand.Intn(6)
				for nextParticle := 0; nextParticle < numNewParticles; nextParticle++ {
					g.Particles[(g.numParticles+nextParticle)%g.maxParticles] = g.Particle(a.Pos)
				}
				g.numParticles = (g.numParticles + numNewParticles) % g.maxParticles
				if a.Size == Small {
					a.Active = false
					rl.PlaySound(g.Sounds.SmallBang)

				} else {
					a1, a2 := g.BreakAsteroid(a)
					g.Asteroids[j] = a1
					g.Asteroids[(j+1)%g.maxAsteroids] = a2
				}
				break
			}
		}
	}
}

func (g *Game) PlayerCollision() {
	for i := 0; i < g.maxAsteroids; i++ {
		asteroid := &g.Asteroids[i]
		if asteroid.Active {
			//rl.DrawCircleLines(int32(g.Player.Position.X), int32(g.Player.Position.Y), 20, rl.Blue)
			//rl.DrawCircleLines(int32(asteroid.Pos.X), int32(asteroid.Pos.Y), asteroid.Rec.Width/4, rl.Blue)
			if rl.CheckCollisionCircles(g.Player.Position, 25, asteroid.Pos, asteroid.Rec.Width/4) {
				g.Player.Hit = true
				asteroid.Active = false
				rl.PlaySound(g.Sounds.LargeBang)
				numNewParticles := 5 + rand.Intn(15)
				for nextParticle := 0; nextParticle < numNewParticles; nextParticle++ {
					g.Particles[(g.numParticles+nextParticle)%g.maxParticles] = g.Particle(g.Player.Position)
				}
				g.numParticles = (g.numParticles + numNewParticles) % g.maxParticles
			}
		}
	}
}

type Textures struct {
	BigAsteroidTex1    rl.Texture2D
	BigAsteroidTex2    rl.Texture2D
	BigAsteroidTex3    rl.Texture2D
	MediumAsteroidTex1 rl.Texture2D
	MediumAsteroidTex2 rl.Texture2D
	MediumAsteroidTex3 rl.Texture2D
	SmallAsteroidTex1  rl.Texture2D
	SmallAsteroidTex2  rl.Texture2D
	SmallAsteroidTex3  rl.Texture2D
	ShipTex1           rl.Texture2D
	ShipTex2           rl.Texture2D
	BulletTex1         rl.Texture2D
	BulletTex2         rl.Texture2D
	UfoTex             rl.Texture2D
}

type Sounds struct {
	Fire        rl.Sound
	LargeBang   rl.Sound
	MediumBang  rl.Sound
	SmallBang   rl.Sound
	Thrust      rl.Sound
	BigSaucer   rl.Sound
	SmallSaucer rl.Sound
	ExtraShip   rl.Sound
	Beat1       rl.Sound
	Beat2       rl.Sound
}

type Game struct {
	GameOver bool
	Dead     bool
	Pause    bool
	SuperFX  bool

	Lives         int
	Score         int
	HiScore       int
	FramesCounter int32

	WindowShouldClose bool

	Player Player

	maxAsteroids int
	maxBullets   int
	maxParticles int
	numAsteroids int
	numBullets   int
	numParticles int
	Asteroids    []Asteroid
	Bullets      []Bullet
	Particles    []Particle

	Textures Textures
	Sounds   Sounds
}

func NewGame() (g Game) {
	g.Init()
	return
}

func (g *Game) Init() {
	// Player Ship
	g.Player = Player{rl.NewRectangle(0, 0, shipHeight, shipWidth), rl.NewVector2(1.0, 0), rl.NewVector2(1.0, 0), rl.NewVector2(screenWidth/2, screenHeight/2), 0, g.Textures.ShipTex1, false}

	g.Asteroids = make([]Asteroid, maxAsteroids)
	g.Bullets = make([]Bullet, maxBullets)
	g.Particles = make([]Particle, maxBullets)
	g.maxAsteroids = maxAsteroids
	g.maxBullets = maxBullets
	g.maxParticles = g.maxBullets
	g.numAsteroids = 0
	g.numBullets = 0
	g.numParticles = 0
	g.Lives = 3
	g.Score = 0
	g.FramesCounter = 0

	g.WindowShouldClose = false

	for i := 0; i < g.maxBullets; i++ {
		g.Bullets[i] = g.Bullet()
		g.Bullets[i].Active = false
	}

	for i := 0; i < g.maxParticles; i++ {
		g.Particles[i] = g.Particle(rl.NewVector2(0, 0))
		g.Particles[i].Active = false
	}

	g.GameOver = false
	g.Dead = false
	g.Pause = false
	g.SuperFX = false
}

// Load Assets
func (g *Game) Load() {
	g.Textures.BigAsteroidTex1 = rl.LoadTexture("images/BigAst1.png")
	g.Textures.BigAsteroidTex2 = rl.LoadTexture("images/BigAst2.png")
	g.Textures.BigAsteroidTex3 = rl.LoadTexture("images/BigAst3.png")
	g.Textures.MediumAsteroidTex1 = rl.LoadTexture("images/MedAst1.png")
	g.Textures.MediumAsteroidTex2 = rl.LoadTexture("images/MedAst2.png")
	g.Textures.MediumAsteroidTex3 = rl.LoadTexture("images/MedAst3.png")
	g.Textures.SmallAsteroidTex1 = rl.LoadTexture("images/SmaAst1.png")
	g.Textures.SmallAsteroidTex2 = rl.LoadTexture("images/SmaAst2.png")
	g.Textures.SmallAsteroidTex3 = rl.LoadTexture("images/SmaAst3.png")
	g.Textures.ShipTex1 = rl.LoadTexture("images/Ship1.png")
	g.Textures.ShipTex2 = rl.LoadTexture("images/Ship2.png")
	g.Textures.BulletTex1 = rl.LoadTexture("images/bullet1.png")
	g.Textures.BulletTex2 = rl.LoadTexture("images/bullet2.png")
	g.Textures.UfoTex = rl.LoadTexture("images/Ufo.png")

	g.Sounds.Fire = rl.LoadSound("sounds/fire.wav")
	g.Sounds.LargeBang = rl.LoadSound("sounds/bangLarge.wav")
	g.Sounds.MediumBang = rl.LoadSound("sounds/bangMedium.wav")
	g.Sounds.SmallBang = rl.LoadSound("sounds/bangSmall.wav")
	g.Sounds.Thrust = rl.LoadSound("sounds/thrust.wav")
	g.Sounds.BigSaucer = rl.LoadSound("sounds/saucerBig.wav")
	g.Sounds.SmallSaucer = rl.LoadSound("sounds/saucerSmall.wav")
	g.Sounds.ExtraShip = rl.LoadSound("sounds/extraShip.wav")
	g.Sounds.Beat1 = rl.LoadSound("sounds/beat1.wav")
	g.Sounds.Beat2 = rl.LoadSound("sounds/beat2.wav")
}

// UnLoad Assets
func (g *Game) Unload() {
	rl.UnloadTexture(g.Textures.BigAsteroidTex1)
	rl.UnloadTexture(g.Textures.BigAsteroidTex2)
	rl.UnloadTexture(g.Textures.BigAsteroidTex3)
	rl.UnloadTexture(g.Textures.MediumAsteroidTex1)
	rl.UnloadTexture(g.Textures.MediumAsteroidTex2)
	rl.UnloadTexture(g.Textures.MediumAsteroidTex3)
	rl.UnloadTexture(g.Textures.SmallAsteroidTex1)
	rl.UnloadTexture(g.Textures.SmallAsteroidTex2)
	rl.UnloadTexture(g.Textures.SmallAsteroidTex3)
	rl.UnloadTexture(g.Textures.ShipTex1)
	rl.UnloadTexture(g.Textures.ShipTex2)
	rl.UnloadTexture(g.Textures.BulletTex1)
	rl.UnloadTexture(g.Textures.BulletTex2)
	rl.UnloadTexture(g.Textures.UfoTex)

	rl.UnloadSound(g.Sounds.Fire)
	rl.UnloadSound(g.Sounds.LargeBang)
	rl.UnloadSound(g.Sounds.MediumBang)
	rl.UnloadSound(g.Sounds.SmallBang)
	rl.UnloadSound(g.Sounds.Thrust)
	rl.UnloadSound(g.Sounds.BigSaucer)
	rl.UnloadSound(g.Sounds.SmallSaucer)
	rl.UnloadSound(g.Sounds.ExtraShip)
	rl.UnloadSound(g.Sounds.Beat1)
	rl.UnloadSound(g.Sounds.Beat2)
}

func main() {
	game := NewGame()
	game.GameOver = false
	rl.InitWindow(1600, 1000, "Go Asteroids!")
	rl.InitAudioDevice() // Initialize audio device
	defer rl.CloseWindow()

	game.Load()

	rl.SetTargetFPS(60)

	// Main Loop
	for !game.WindowShouldClose {
		game.FramesCounter++
		flag := true
		game.Update()
		if game.FramesCounter%60 == 0 {
			if flag {
				rl.PlaySound(game.Sounds.Beat1)
				flag = false
			} else {
				rl.PlaySound(game.Sounds.Beat2)
				flag = true
			}

		}
		game.Draw()

	}

	game.Unload()

	rl.CloseAudioDevice()
	rl.CloseWindow()

}

func (g *Game) Update() {
	g.updatePlayer()

	g.updateAsteroids()

	g.updateBullets()

	g.BulletCollisions()

	g.PlayerCollision()

	g.updateParticles()
}

func (g *Game) DrawShip() {
	rl.DrawTexturePro(g.Player.Tex, g.Player.ShipRec,
		rl.NewRectangle(g.Player.Position.X, g.Player.Position.Y, float32(g.Player.Tex.Width), float32(g.Player.Tex.Height)),
		rl.NewVector2(float32(g.Player.Tex.Width/2), float32(g.Player.Tex.Height/2)), float32(math.Atan2(float64(g.Player.ShipDirection.Y), float64(g.Player.ShipDirection.X))*180/math.Pi), rl.Pink)
}

func (g *Game) DrawAsteroids() {
	for i := 0; i < g.maxAsteroids; i++ {
		asteroid := &g.Asteroids[i]
		if asteroid.Active {
			rl.DrawTexturePro(asteroid.Tex, asteroid.Rec, rl.NewRectangle(asteroid.Pos.X, asteroid.Pos.Y, float32(asteroid.Tex.Height), float32(asteroid.Tex.Width)),
				rl.NewVector2(float32(asteroid.Tex.Width/2), float32(asteroid.Tex.Height/2)), asteroid.Rotation, rl.Color{180, 100, 30, 255})
		}
	}
}

func (g *Game) DrawBullets() {
	for i := 0; i < g.maxBullets; i++ {
		b := &g.Bullets[i]
		if b.Active {
			rl.DrawTexturePro(g.Textures.BulletTex2, rl.NewRectangle(0, 0, 32, 32), rl.NewRectangle(b.Pos.X, b.Pos.Y, 32, 32), rl.NewVector2(16, 16), 0, rl.RayWhite)
		}
	}
}

func (g *Game) DrawParticles() {
	for i := 0; i < g.maxParticles; i++ {
		p := &g.Particles[i]
		if p.Active {
			rl.DrawTexturePro(g.Textures.BulletTex1, rl.NewRectangle(0, 0, 32, 32), rl.NewRectangle(p.Pos.X, p.Pos.Y, 32, 32), rl.NewVector2(16, 16), 0, rl.RayWhite)
		}
	}
}

func (g *Game) Draw() {
	rl.BeginDrawing()

	rl.ClearBackground(rl.Black)

	if !g.GameOver {
		rl.DrawText(fmt.Sprintf("SCORE: %d", g.Score), 10, 10, 40, rl.RayWhite)

		for i := 0; i < g.Lives; i++ {
			rl.DrawTextureEx(g.Textures.ShipTex1, rl.NewVector2(10+30*float32(i), 100), -90, .5, rl.Pink)
		}

		g.DrawShip()
		g.DrawAsteroids()
		g.DrawBullets()
		g.DrawParticles()
	}

	rl.EndDrawing()
}
