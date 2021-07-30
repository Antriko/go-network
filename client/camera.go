package client

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// CustomCamera struct
type CustomCamera struct {
	Camera        rl.Camera
	Angle         rl.Vector2
	Zoom          float64
	RotationSpeed float32
	ZoomSpeed     float32
}

// NewCustomCamera requires zoom and rotationSpeed, returns a pointer to a
// CustomCamera object
func NewCustomCamera(zoom float64, rotationSpeed, zoomSpeed float32) *CustomCamera {
	c := &CustomCamera{}

	c.Zoom = zoom
	c.RotationSpeed = rotationSpeed
	c.ZoomSpeed = zoomSpeed
	c.Angle = rl.NewVector2(0.0, 0.0)
	c.Camera = rl.Camera{} // Set the internal raylib CustomCamera
	c.Camera.Position = rl.NewVector3(0.0, 32.0, 32.0)
	c.Camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	c.Camera.Fovy = 100.0
	//c.Camera.Type = rl.CameraTypeOrthographic
	c.Camera.Projection = rl.CameraPerspective

	// Set CustomCamera default values
	v1 := c.Camera.Position
	v2 := c.Camera.Target // Just use whatever the zero value is, update later
	dx := float64(v2.X - v1.X)
	dy := float64(v2.Y - v1.Y)
	dz := float64(v2.Z - v1.Z)
	c.Angle.X = float32(math.Atan2(dx, dz)) + math.Pi
	c.Angle.Y = float32(math.Atan2(dy, math.Sqrt(dx*dx+dz*dz)))

	return c
}

// Update update camera position and rotation using the keyboard
func (c *CustomCamera) Update(dt float32) {
	c.Angle.X = float32(math.Remainder(float64(c.Angle.X), math.Pi*2))

	c.Camera.Position.X = float32(math.Sin(float64(c.Angle.X))*c.Zoom*math.Cos(float64(c.Angle.Y))) + c.Camera.Target.X
	c.Camera.Position.Y = float32(math.Sin(float64(c.Angle.Y))*c.Zoom*math.Sin(float64(c.Angle.Y))) + c.Camera.Target.Y
	if c.Angle.Y > 0.0 {
		c.Camera.Position.Y *= -1.0
	}
	c.Camera.Position.Z = float32(math.Cos(float64(c.Angle.X))*c.Zoom*math.Cos(float64(c.Angle.Y))) + c.Camera.Target.Z
}

// SetTarget set the target for the camera
func (c *CustomCamera) SetTarget(target rl.Vector3) {
	c.Camera.Target = target
}

// SetPosition set the position of the camera
func (c *CustomCamera) SetPosition(x, y, z float32) {
	c.Camera.Position.X = x
	c.Camera.Position.Y = y
	c.Camera.Position.Z = z
}
