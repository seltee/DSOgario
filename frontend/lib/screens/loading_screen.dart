import 'package:flutter/material.dart';
import 'package:loading_animation_widget/loading_animation_widget.dart';

class LoadingScreen extends StatelessWidget {
  const LoadingScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(color: Theme.of(context).canvasColor),
      child: Center(
        child: LoadingAnimationWidget.twoRotatingArc(
          color: Theme.of(context).hintColor,
          size: 160,
        ),
      ),
    );
  }
}
