import 'package:flutter/material.dart';

enum GameEntityType { unknown, player, crumb }

class GameWorldEntity {
  final int id;
  final GameEntityType type;
  final int colorIndex;

  int size;
  Offset relPos;
  Offset relPosTarget;

  bool updatedThisFrame = true;

  GameWorldEntity({
    required this.id,
    required this.type,
    required this.colorIndex,
    required this.size,
    required this.relPos,
    required this.relPosTarget,
  });

  void provideUpdate(Offset newRelPosTarget, int newSize) {
    relPosTarget = newRelPosTarget;
    size = newSize;
    updatedThisFrame = true;
  }

  void update(double delta) {
    relPos = Offset(
      relPos.dx + (relPosTarget.dx - relPos.dx) * delta * 20.0,
      relPos.dy + (relPosTarget.dy - relPos.dy) * delta * 20.0,
    );
  }
}
