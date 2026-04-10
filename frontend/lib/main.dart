import 'package:flutter/material.dart';
import 'package:frontend/core/layout/night_mode.dart';
import 'package:frontend/core/theme.dart';
import 'package:frontend/models/game_controller.dart';
import 'package:frontend/models/game_model.dart';
import 'package:frontend/screens/auth_screen.dart';
import 'package:frontend/screens/game_screen.dart';
import 'package:frontend/screens/loading_screen.dart';

void main() {
  runApp(GameApp());
}

class GameApp extends StatefulWidget {
  const GameApp({super.key});

  @override
  State<GameApp> createState() => _GameState();
}

class _GameState extends State<GameApp> with SingleTickerProviderStateMixin {
  late GameController controller;
  late GameModel gameModel;

  @override
  void initState() {
    super.initState();

    gameModel = GameModel();
    controller = GameController(gameModel, this);
    controller.start();
    gameModel.initGame();
  }

  @override
  void dispose() {
    controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return GameModelProvider(
      model: gameModel,
      child: ListenableBuilder(
        listenable: gameModel.themeModel,
        builder: (context, _) {
          return MaterialApp(
            theme: lightTheme,
            darkTheme: darkTheme,
            themeMode: gameModel.themeModel.themeMode,
            home: Stack(
              children: [
                Positioned.fill(child: ScreenSwitcher()),
                Align(
                  alignment: AlignmentGeometry.topRight,
                  child: NightModeSwitcher(),
                ),
              ],
            ),
          );
        },
      ),
    );
  }
}

class ScreenSwitcher extends StatelessWidget {
  const ScreenSwitcher({super.key});

  @override
  Widget build(BuildContext context) {
    final gameModel = GameModelProvider.of(context);

    return ListenableBuilder(
      listenable: gameModel,
      builder: (context, _) {
        return switch (gameModel.getState()) {
          GameState.loading => LoadingScreen(),
          GameState.auth => AuthScreen(),
          GameState.game => GameScreen(),
        };
      },
    );
  }
}
