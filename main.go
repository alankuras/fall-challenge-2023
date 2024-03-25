package main

import (
	"fmt"
	"math"
	"os"
)

var DEBUG, TRACE bool

func debug(args ...interface{}) {
	if DEBUG {
		fmt.Fprintln(os.Stderr, concatStringsAndInt(args))
	}
}

func concatStringsAndInt(args ...interface{}) string {
	return fmt.Sprint(args...)
}

func trace(d string) {
	if TRACE {
		fmt.Fprintln(os.Stderr, d)
	}
}

func wait(i int) {
	fmt.Println(fmt.Sprintf("WAIT %d", i))
}

func move(x, y, light float64, monster *Fish, drone Drone) {
	collision := monster.Id
	if monster.Type == MONSTER && monster.Pos.DistanceTo(Vector{drone.Pos.x, drone.Pos.y}) <= 810.0 {
		collision = int(monster.Pos.DistanceTo(Vector{drone.Pos.x, drone.Pos.y}))
		fmt.Println(fmt.Sprintf("MOVE %d %d %d %d", int(x)+int(monster.Speed.x*-1), int(y), 0, collision))
		return
	}
	fmt.Println(fmt.Sprintf("MOVE %d %d %d %d", int(x), int(y), int(light), collision))
}

type FishType int

const (
	JELLY FishType = iota
	FISH
	CRAB
	MONSTER
)

func (f FishType) String() string {
	return [...]string{"JELLY", "FISH", "CRAB", "MONSTER"}[f]
}

func (f FishType) EnumIndex() int {
	return int(f)
}

type Color int

const (
	PINK Color = iota
	YELLOW
	GREEN
	BLUE
	RED
)

func (c Color) String() string {
	return [...]string{"PINK", "YELLOW", "GREEN", "BLUE", "RED"}[c]
}

func (c Color) EnumIndex() int {
	return int(c)
}

type Vector struct {
	x float64
	y float64
}

func (v Vector) DistanceTo(v2 Vector) float64 {
	return math.Sqrt(math.Pow(v.x-v2.x, 2) + math.Pow(v.y-v2.y, 2))
	//deltaX := v.x - v2.x
	//deltaY := v.y - v2.y
	//return math.Sqrt(deltaX*deltaX + deltaY*deltaY)
}

func NewVector(x, y float64) Vector {
	return Vector{
		x: x,
		y: y,
	}
}

type DroneLocation struct {
	droneID  int
	location string
}
type Fish struct {
	Type              FishType
	Pos               Vector
	Color             Color
	StartY            int
	Speed             Vector
	Id                int
	LowY, HighY       int
	IsFleeing         bool
	FleeingFromPlayer int
	Saved             bool
	SavedByOpponent   bool
	Scanned           bool
	ScannedByOpponent bool
	ScannedByDroneId  int
	LocationToDrones  []DroneLocation
	BeingTracked      bool
}

func (f *Fish) OutOfBound() bool {
	return f.Pos.x < 0 || f.Pos.x > 9999 || f.Pos.y < 0 || f.Pos.y > 9999
}
func (f *Fish) PrintFields() string {
	return fmt.Sprint(
		"Type: ", f.Type,
		", BeingTracked: ", f.BeingTracked,
		", Pos: ", f.Pos,
		", Color: ", f.Color,
		", Speed: ", f.Speed,
		", ID: ", f.Id,
		", Saved: ", f.Saved,
		", Scanned: ", f.Scanned,
		", Saved: ", f.Saved,
		", ScannedByOpponent: ", f.ScannedByOpponent,
		", SavedByOpponent: ", f.SavedByOpponent,
		", ScannedByDroneId: ", f.ScannedByDroneId,
		", Locations: ", f.LocationToDrones,
	)
}

func NewFish(x, y, fishType, color, id int) *Fish {
	if fishType == -1 {
		fishType = 3
		color = 4
	}
	return &Fish{
		Type:              FishType(fishType),
		Color:             Color(color),
		Pos:               NewVector(0, 0),
		Id:                id,
		Speed:             NewVector(0, 0),
		Saved:             false,
		SavedByOpponent:   false,
		Scanned:           false,
		ScannedByOpponent: false,
		ScannedByDroneId:  0,
	}
}

type Drone struct {
	Id      int
	Pos     Vector
	Battery int
}

func (d *Drone) PrintFields() string {
	return fmt.Sprint(
		"Battery: ", d.Battery,
		", Pos: ", d.Pos,
		", ID: ", d.Id,
	)
}

var creatureCount int

func getCreatures() []*Fish {
	fmt.Scan(&creatureCount)

	var creatures = make([]*Fish, 0)
	for i := 0; i < creatureCount; i++ {

		var creatureId, color, _type int
		fmt.Scan(&creatureId, &color, &_type)
		creatures = append(creatures, NewFish(0, 0, _type, color, creatureId))
	}
	return creatures
}

var myScore, opponentScore, myScanCount, opponentScanCount int
var scannedCreatures []int
var myDroneCount, myOpponentDroneCount int
var myDroneId, opponentDroneId int
var myDroneX, myDroneY, opponentDroneX, opponentDroneY float64
var battery, emergency int

var myDrones []Drone
var opponentDrones []Drone

func getMyDrones(droneCount int) []Drone {
	var drones = make([]Drone, 0)

	for i := 0; i < droneCount; i++ {
		fmt.Scan(&myDroneId, &myDroneX, &myDroneY, &emergency, &battery)
		myDrone := Drone{
			Id:      myDroneId,
			Pos:     Vector{myDroneX, myDroneY},
			Battery: battery,
		}
		drones = append(drones, myDrone)
	}
	return drones
}

func getOpponentDrones(droneCount int) []Drone {
	var drones = make([]Drone, 0)
	for i := 0; i < droneCount; i++ {
		fmt.Scan(&opponentDroneId, &opponentDroneX, &opponentDroneY, &emergency, &battery)
		opponentDrone := Drone{
			Id:      myDroneId,
			Pos:     Vector{opponentDroneX, opponentDroneY},
			Battery: battery,
		}
		drones = append(drones, opponentDrone)
	}

	return drones
}

func FindClosest(creatures []*Fish, d Drone) *Fish {
	var closest *Fish
	for _, creature := range creatures {
		if creature.Saved || creature.Scanned || (creature.Pos.x == 0 && creature.Pos.y == 0) || creature.OutOfBound() || creature.Type == MONSTER {
			continue
		}
		if closest == nil || d.Pos.DistanceTo(creature.Pos) < d.Pos.DistanceTo(closest.Pos) {
			closest = creature
		}
	}

	return closest
}

func FindClosestMonster(creatures []*Fish, d Drone) *Fish {
	var closest *Fish
	for _, creature := range creatures {
		if creature.Type != MONSTER {
			continue
		}
		if closest == nil || d.Pos.DistanceTo(creature.Pos) < d.Pos.DistanceTo(closest.Pos) {
			closest = creature
		}
	}

	return closest
}

type Direction struct {
	x int
	y int
}

var directions = map[string]Direction{
	"TL": {x: -1, y: -1},
	"TR": {x: 1, y: -1},
	"BL": {x: -1, y: 1},
	"BR": {x: 1, y: 1},
}

func getGameTurn(creatures []*Fish) {
	fmt.Scan(&myScore)
	fmt.Scan(&opponentScore)
	fmt.Scan(&myScanCount)

	for i := 0; i < myScanCount; i++ {
		var creatureId int
		fmt.Scan(&creatureId)
		for _, fish := range creatures {
			if fish.Id == creatureId {
				fish.Saved = true
				break
			}
		}
	}
	fmt.Scan(&opponentScanCount)

	for i := 0; i < opponentScanCount; i++ {
		var creatureId int
		fmt.Scan(&creatureId)
		for _, fish := range creatures {
			if fish.Id == creatureId {
				fish.SavedByOpponent = true
				break
			}
		}
	}
	fmt.Scan(&myDroneCount)

	myDrones := getMyDrones(myDroneCount)
	_ = myDrones

	fmt.Scan(&myOpponentDroneCount)

	opponentDrones := getOpponentDrones(myOpponentDroneCount)
	_ = opponentDrones

	var droneScanCount int
	fmt.Scan(&droneScanCount)

	for i := 0; i < droneScanCount; i++ {
		var droneId, creatureId int
		fmt.Scan(&droneId, &creatureId)
		for _, fish := range creatures {
			if fish.Id == creatureId {
				fish.ScannedByOpponent = true
				fish.ScannedByDroneId = droneId
				break
			}
		}
		for _, drone := range myDrones {
			if drone.Id == droneId {
				for _, fish := range creatures {
					if fish.Id == creatureId {
						fish.Scanned = true
						fish.ScannedByOpponent = false
						fish.ScannedByDroneId = droneId
						break
					}
				}
			}
		}
	}
	var visibleCreatureCount int
	fmt.Scan(&visibleCreatureCount)

	for _, fish := range creatures {
		fish.Pos = Vector{0, 0}
	}
	for i := 0; i < visibleCreatureCount; i++ {
		var scannedCreatureId int
		var x, y, Vx, Vy float64
		fmt.Scan(&scannedCreatureId, &x, &y, &Vx, &Vy)
		for _, fish := range creatures {
			if fish.Id == scannedCreatureId {
				fish.Pos = Vector{x, y}
				fish.Speed = Vector{Vx, Vy}
				break
			}
		}
	}

	var radarBlipCount int
	fmt.Scan(&radarBlipCount)

	// reset locations
	for _, fish := range creatures {
		fish.LocationToDrones = []DroneLocation{}
		fish.BeingTracked = false
	}

	for i := 0; i < radarBlipCount; i++ {
		var droneId, creatureId int
		var radar string
		fmt.Scan(&droneId, &creatureId, &radar)
		for _, creature := range creatures {
			if creature.Id == creatureId {
				creature.LocationToDrones = append(creature.LocationToDrones, DroneLocation{
					droneID:  droneId,
					location: radar,
				})
			}
		}
	}

	for _, creature := range creatures {
		trace(creature.PrintFields())
	}
	for _, drone := range myDrones {
		debug(drone.PrintFields())
		monster := FindClosestMonster(creatures, drone)

		//
		// Check if something need to be saved
		// *****************************************************************************************************
		savedCounter := 0
		for _, fish := range creatures {
			//if !fish.Saved && fish.Scanned && !fish.SavedByOpponent && fish.ScannedByDroneId == drone.Id {
			if !fish.Saved && fish.Scanned && fish.ScannedByDroneId == drone.Id {
				savedCounter++
			}
		}
		if savedCounter >= 3 {
			debug("Going to save some animals by drone id ", drone.Id)
			move(drone.Pos.x, 500, 1, monster, drone)
			continue
		}

		// if there is something visible close
		closest := FindClosest(creatures, drone)
		if closest != nil && !closest.Scanned && !closest.Saved {
			debug(closest.PrintFields())
			//move(closest.Pos.x+closest.Speed.x, closest.Pos.y+closest.Speed.y, 1, *monster, drone)
			move(closest.Pos.x, closest.Pos.y, 1, monster, drone)
			continue
		}

		speed := float64(600)
		droneFoundSomething := false
		// Hunting
		// *****************************************************************************************************
		for _, fish := range creatures {
			if fish.OutOfBound() || fish.Type == MONSTER {
				continue
			}
			// ignore the fact that its scanned by opponent, we are not doing race, we are collecting everything
			//if !fish.Scanned && !fish.Saved && !fish.ScannedByOpponent {
			if !fish.Scanned && !fish.Saved && !fish.BeingTracked {
				location := ""
				for _, fishLocation := range fish.LocationToDrones {
					if drone.Id == fishLocation.droneID {
						location = fishLocation.location
					}
				}

				if direction, ok := directions[location]; ok {
					debug("Tracking fish ", fish.Id, " by drone ", drone.Id, " going ", location, direction.x, direction.y)
					if drone.Pos.y < 2000 {
						move(drone.Pos.x, drone.Pos.y+(speed*float64(direction.y)), 0, monster, drone)
					} else {
						move(drone.Pos.x+(speed*float64(direction.x)), drone.Pos.y+(speed*float64(direction.y)), 1, monster, drone)
					}
					fish.BeingTracked = true
					droneFoundSomething = true
					break
				}
			}
		}

		if droneFoundSomething {
			continue
		}

		//for _, fish := range creatures {
		//	if !fish.Scanned && !fish.Saved {
		//		move(fish.Pos.x+fish.Speed.x, fish.Pos.y+fish.Speed.y, 1)
		//		droneFoundSomething = true
		//		break
		//	}
		//}

		if !droneFoundSomething {
			move(5000, 0, 0, monster, drone)
		}

	}

}

func main() {
	TRACE = true
	DEBUG = true

	creatures := getCreatures()

	for {
		getGameTurn(creatures)
	}
}
