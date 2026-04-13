import 'package:flutter/material.dart';
import 'package:frontend/core/blob_colors.dart';
import 'package:frontend/core/layout/game_over.dart';
import 'package:frontend/core/layout/score.dart';
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
        return Stack(
          children: [
            Positioned.fill(
              child: Container(
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

                  child: CustomPaint(
                    painter: GamePainter(gameWorld: gameWorld),
                  ),
                ),
              ),
            ),
            Score(),
            Center(child: NewGame()),
          ],
        );
      },
    );
  }
}

class GamePainter extends CustomPainter {
  final GameWorld gameWorld;
  GamePainter({required this.gameWorld});

  double getVisibleRadius(double size) {
    return (size + 10) * 0.02;
  }

  @override
  void paint(Canvas canvas, Size size) {
    gameWorld.updateScreenSize(size);

    final center = Offset(size.width / 2, size.height / 2);
    final crumbColor = Color.fromARGB(255, 120, 120, 120);
    double zoom = gameWorld.zoom;

    // Draw all entities at screen center
    for (var entity in gameWorld.entities.values) {
      double radius = getVisibleRadius(entity.size.toDouble());

      final screenPos = center.translate(
        entity.relPos.dx * zoom,
        entity.relPos.dy * zoom,
      );

      canvas.drawCircle(
        screenPos,
        radius * zoom,
        Paint()
          ..color = entity.type == GameEntityType.player
              ? Color(blobColors[entity.colorIndex])
              : crumbColor, // use your entity.color
      );
    }

    // Draw entity names
    for (var entity in gameWorld.entities.values) {
      var item = gameWorld.gameScore.getScoreItemById(entity.id);
      if (item == null) continue;

      double radius = getVisibleRadius(entity.size.toDouble());

      final textStyle = TextStyle(
        color: Color(blobColors[item.colorIndex]),
        fontSize: 14,
        fontWeight: FontWeight.bold,
        shadows: const [
          Shadow(blurRadius: 3, color: Colors.white, offset: Offset(0, 0)),
        ],
      );

      final textPainter = TextPainter(
        text: TextSpan(text: item.name, style: textStyle),
        textDirection: TextDirection.ltr,
        textAlign: TextAlign.center,
      );

      textPainter.layout();

      final double textWidth = textPainter.width;

      final screenPos = center.translate(
        entity.relPos.dx * zoom - textWidth / 2,
        entity.relPos.dy * zoom - (radius * 1.5 + 1.5) * zoom,
      );

      textPainter.paint(canvas, screenPos);
    }
  }

  @override
  bool shouldRepaint(covariant GamePainter oldDelegate) => true; // or compare data
}
