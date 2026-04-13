import 'package:flutter/material.dart';
import 'package:frontend/core/blob_colors.dart';
import 'package:frontend/models/game_model.dart';
import 'package:loading_animation_widget/loading_animation_widget.dart';

class AuthScreen extends StatelessWidget {
  const AuthScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final gameModel = GameModelProvider.of(context);

    return ListenableBuilder(
      listenable: gameModel,
      builder: (context, _) {
        return Container(
          decoration: BoxDecoration(color: Theme.of(context).canvasColor),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Text("Dogario", style: Theme.of(context).textTheme.displayLarge),
              Padding(padding: EdgeInsetsGeometry.symmetric(vertical: 32)),
              Text(
                gameModel.name,
                style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                  color: Color(blobColors[gameModel.colorIndex]),
                ),
              ),
              Padding(padding: EdgeInsetsGeometry.symmetric(vertical: 4)),
              ElevatedButton(
                onPressed: gameModel.isLoggingIn
                    ? null
                    : () => gameModel.suggestNewName(),
                child: const Text('Change name'),
              ),
              Padding(padding: EdgeInsetsGeometry.symmetric(vertical: 8)),
              LoginButton(),
            ],
          ),
        );
      },
    );
  }
}

class LoginButton extends StatelessWidget {
  const LoginButton({super.key});

  @override
  Widget build(BuildContext context) {
    final gameModel = GameModelProvider.of(context);

    return SizedBox(
      width: 200,
      height: 32,
      child: Center(
        child: ListenableBuilder(
          listenable: gameModel,
          builder: (context, _) {
            if (gameModel.isLoggingIn) {
              return LoadingAnimationWidget.twoRotatingArc(
                color: Theme.of(context).hintColor,
                size: 24,
              );
            } else {
              return ElevatedButton(
                onPressed: () => gameModel.enterTheGame(),
                child: const Text('Enter'),
              );
            }
          },
        ),
      ),
    );
  }
}
