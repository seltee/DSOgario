import 'package:flutter/material.dart';
import 'package:frontend/models/game_model.dart';

class NightModeSwitcher extends StatelessWidget {
  const NightModeSwitcher({super.key});

  @override
  Widget build(BuildContext context) {
    final themeModel = GameModelProvider.of(context).themeModel;

    return ListenableBuilder(
      listenable: themeModel,
      builder: (context, _) {
        return Padding(
          padding: EdgeInsetsGeometry.symmetric(vertical: 8, horizontal: 8),
          child: IconButton(
            icon: Icon(
              // Show correct icon based on current mode
              Theme.of(context).brightness == Brightness.light
                  ? Icons.dark_mode_outlined
                  : Icons.light_mode_outlined,
            ),
            tooltip: 'Toggle theme',
            onPressed: () {
              // Call your toggle function from provider / notifier / setState
              themeModel.toggleTheme();
            },
          ),
        );
      },
    );
  }
}
