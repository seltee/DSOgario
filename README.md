# DSOgar.io

## Overview

DSOgar.io is a real-time multiplayer game inspired by Agar.io and Ogario-style mechanics.  
The project is built with a custom high-performance backend in Go and a cross-platform frontend in Flutter.

The backend uses a single game loop and communicates with clients over WebSockets using a compact binary protocol. The server is designed to avoid blocking operations, keeping gameplay smooth even with multiple players.

The game is fully playable with a working client-server architecture.

---

## Why This Project

This project focuses on real-time multiplayer architecture rather than gameplay itself.

Key areas:

- Designing a non-blocking game loop in Go
- Efficient binary communication over WebSockets
- Handling multiple players without blocking the main loop
- Minimizing bandwidth via compact data structures

---

## Tech Stack

### Backend

- Go (Golang)
- net/http
- Chi router
- Gorilla WebSocket

### Frontend

- Dart
- Flutter (Windows / Linux / Web)

---

## Key Features

- **Authoritative game server architecture**
  - Single game loop running in its own goroutine
  - Communication via channels (`join`, `input`, etc.) to avoid blocking

- **Efficient networking**
  - Custom **binary protocol** for minimal bandwidth usage
  - Compact entity packets

- **Token-based authentication**
  - Random unique tokens

- **Optimized entity system**
  - Relative positioning for rendering
  - Visibility filtering per player
  - Fast removal using swap-and-pop

- **Scalable WebSocket handling**
  - Dedicated `readPump` / `writePump` per player
  - Non-blocking send channels

---

## Architecture Highlights

- Game state updates are fully decoupled from network I/O
- All gameplay logic (movement, collision, entity updates) runs inside the main loop
- Clients receive only **visible entities**, reducing bandwidth usage
- Designed with performance and scalability in mind

---

## Getting Started

### Prerequisites

- Go
- Dart
- Flutter

---

### Run Backend

```bash
cd backend
go run main.go
```

### Run Frontend Windows

```bash
cd frontend
flutter run -d windows
```

or

### Run Frontend Linux

```bash
cd frontend
flutter run -d linux
```

---

## License

MIT License – feel free to explore and learn from the code.

---

## Project Status

- Backend: implemented and stable
- Frontend: implemented and fully connected to backend
- Core gameplay: working (movement, eating, growth)
- WebSocket communication: fully functional
- Binary protocol: implemented and in use
- UI/UX improvements (player names, score display, HUD)
- Game is playable end-to-end.

---

## Next Steps

- Complete day/night mode
- General polish and quality-of-life improvements
- Additional gameplay features (splitting, merging, etc.)
- Performance tuning and optimization
- Add sound effects and background music
