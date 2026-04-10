import 'package:flutter/material.dart';
import 'package:frontend/core/blob_colors.dart';
import 'package:frontend/models/game_model.dart';
import 'package:frontend/models/game_world.dart';
import 'package:frontend/models/game_world_entity.dart';

class GameScreen extends StatelessWidget {
  const GameScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final gameWorld = GameModelProvider.of(context).gameWorld;

    return ListenableBuilder(
      listenable: gameWorld,
      builder: (context, widget) {
        return Container(
          decoration: BoxDecoration(color: Theme.of(context).canvasColor),
          child: GestureDetector(
            // Tap / Move / Release callbacks
            onPanStart: (DragStartDetails details) {
              final Offset localPos = details.localPosition;
              gameWorld.startDirectingPlayer(localPos);
            },

            onPanUpdate: (DragUpdateDetails details) {
              final Offset localPos = details.localPosition;
              gameWorld.updatePlayerDirection(localPos);
            },

            onPanEnd: (DragEndDetails details) {
              gameWorld.stopDirectingPlayer();
            },

            child: CustomPaint(painter: GamePainter(gameWorld: gameWorld)),
          ),
        );
      },
    );
  }
}

class GamePainter extends CustomPainter {
  final GameWorld gameWorld;
  GamePainter({required this.gameWorld});

  @override
  void paint(Canvas canvas, Size size) {
    gameWorld.updateScreenSize(size);

    final center = Offset(size.width / 2, size.height / 2);
    final crumbColor = Color.fromARGB(255, 120, 120, 120);
    double zoom = gameWorld.zoom;

    // Draw all entities at screen center
    for (var entity in gameWorld.entities.values) {
      double radius = (entity.size.toDouble() + 10) * 0.02 * zoom;

      final screenPos = center.translate(
        entity.relPos.dx * zoom,
        entity.relPos.dy * zoom,
      );

      canvas.drawCircle(
        screenPos,
        radius,
        Paint()
          ..color = entity.type == GameEntityType.player
              ? Color(blobColors[entity.colorIndex])
              : crumbColor, // use your entity.color
      );
    }
  }

  @override
  bool shouldRepaint(covariant GamePainter oldDelegate) => true; // or compare data
}
