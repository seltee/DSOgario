import 'package:flutter/material.dart';
import 'package:frontend/core/blob_colors.dart';
import 'package:frontend/models/game_model.dart';

// ignore: use_key_in_widget_constructors
class Score extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    final gameWorld = GameModelProvider.of(context).gameWorld;

    return ListenableBuilder(
      listenable: gameWorld.gameScore,
      builder: (context, widget) {
        return SizedBox(
          height: 106,
          child: ShaderMask(
            shaderCallback: (Rect bounds) {
              return LinearGradient(
                begin: Alignment.topCenter,
                end: Alignment.bottomCenter,
                colors: [Colors.black, Colors.black, Colors.transparent],
                stops: [0.0, 0.92, 1.0],
              ).createShader(bounds);
            },
            blendMode: BlendMode.dstIn,
            child: Stack(
              children: [
                for (final (index, scoreItem)
                    in gameWorld.gameScore.scoreList.take(8).indexed)
                  AnimatedPositioned(
                    key: ValueKey(scoreItem.id),
                    duration: const Duration(milliseconds: 400),
                    curve: Curves.easeInOut,
                    left: 4.0,
                    top: 4.0 + index.toDouble() * 20.0,
                    child: Text(
                      '${scoreItem.name} (${scoreItem.size})',
                      style: Theme.of(context).textTheme.displayLarge?.merge(
                        TextStyle(
                          fontSize: 14,
                          color: Color(blobColors[scoreItem.colorIndex]),
                        ),
                      ),
                    ),
                  ),
              ],
            ),
          ),
        );
      },
    );
  }
}
