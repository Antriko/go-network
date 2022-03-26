# Project to learn GoLang

Please keep in mind this is a personal project in which I am trying to learn on my own for fun, everything is most likely inefficient and not the best.

## Run application
`go run . server` to host the server

`go run . client` to run a client application

`go run . map` for map experimentation

## Network
The backend uses DTLS and TLC.

The players are able to communicate to one another via the chat, as well as update the players locations for all users that are connected to the server and show off what user models they've chosen for their player pre-joining the game.


## Game Engine
I've chosen to use [raylib-go](https://github.com/gen2brain/raylib-go) which essentially is go bindings of the [raylib](https://www.raylib.com/) graphics library. My main focus was on the language more so than the engine.


## Modeling
I use [blockbench](https://www.blockbench.net/) to create my models which I then export as GLTF which raylib supports. Animations within raylib-go are not supported however I have been able to make do with the limitations.

## Procedurally generated maps
With the use of noise, I've been able to create unique islands that are able to be rendered into the game.

## "Gameplay"
Keep in mind I'm not very good at design to begin with, was more so focused on the functionality of the "game" rather than design.


TODO - Add images