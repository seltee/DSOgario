import 'dart:ui';

import 'package:flutter/material.dart';

class ThemeModel extends ChangeNotifier {
  ThemeMode _themeMode = ThemeMode.system;

  ThemeMode get themeMode => _themeMode;

  bool get isSystemDark =>
      PlatformDispatcher.instance.platformBrightness == Brightness.dark;
  bool get isSystemLight => !isSystemDark;

  void toggleTheme() {
    if (_themeMode == ThemeMode.light) {
      _themeMode = ThemeMode.dark;
    } else if (_themeMode == ThemeMode.dark) {
      _themeMode = ThemeMode.light;
    } else if (_themeMode == ThemeMode.system) {
      _themeMode = isSystemDark ? ThemeMode.light : ThemeMode.dark;
    }
    notifyListeners();
  }
}
