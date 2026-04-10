import 'package:flutter/widgets.dart';

abstract class Provider<T> extends InheritedWidget {
  final T model;

  const Provider({super.key, required this.model, required super.child});

  @override
  bool updateShouldNotify(covariant Provider<T> oldWidget) {
    // Return true if the reference to the model changes
    return oldWidget.model != model;
  }
}
