import 'package:flutter/material.dart';
import 'package:frontend/models/game_model.dart';
import 'package:loading_animation_widget/loading_animation_widget.dart';

class LoadingScreen extends StatelessWidget {
  const LoadingScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final gameModel = GameModelProvider.of(context);

    return ListenableBuilder(
      listenable: gameModel,
      builder: (context, _) {
        return Container(
          decoration: BoxDecoration(color: Theme.of(context).canvasColor),
          child: Center(
            child: gameModel.serverIsUnavailable
                ? Text(
                    "Server is unavailalbe :c",
                    style: Theme.of(context).textTheme.displaySmall?.merge(
                      TextStyle(color: Color(0xFFAA0000), fontSize: 20),
                    ),
                  )
                : LoadingAnimationWidget.twoRotatingArc(
                    color: Theme.of(context).hintColor,
                    size: 160,
                  ),
          ),
        );
      },
    );
  }
}
