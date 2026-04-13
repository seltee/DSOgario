import 'package:flutter/material.dart';
import 'package:frontend/models/game_model.dart';

class NewGame extends StatelessWidget {
  const NewGame({super.key});

  @override
  Widget build(BuildContext context) {
    final gameModel = GameModelProvider.of(context);
    final gameWorld = gameModel.gameWorld;

    return ListenableBuilder(
      listenable: gameWorld,
      builder: (context, widget) {
        if (gameWorld.playerIsEaten) {
          return SizedBox(
            height: 120,
            child: Column(
              children: [
                Padding(padding: EdgeInsetsGeometry.symmetric(vertical: 16)),
                ElevatedButton(
                  onPressed: () => gameModel.restart(),
                  child: const Text('Restart'),
                ),
                Padding(padding: EdgeInsetsGeometry.symmetric(vertical: 8)),
                ElevatedButton(
                  onPressed: () => gameModel.goToMainMenu(),
                  child: const Text('Main menu'),
                ),
              ],
            ),
          );
        } else {
          return Container();
        }
      },
    );
  }
}
