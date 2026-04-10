import 'dart:async';
import 'dart:typed_data';

import 'package:flutter/material.dart';
import 'package:frontend/core/blob_colors.dart';
import 'package:frontend/models/config.dart';
import 'package:frontend/models/game_world_entity.dart';

class GameWorld extends ChangeNotifier {
  Map<String, GameWorldEntity> entities = <String, GameWorldEntity>{};
  Size? screenSize;

  double zoom = 12.0;
  double targetZoom = 12.0;
  bool directingPlayer = false;
  Offset playerDirection = Offset(0, 0);

  Timer? _inputTimer;

  Function(Offset target)? sendPlayerDirection;

  void provideFrame(ByteData data) {
    int entityCount = data.getUint16(2, Endian.big);

    for (var entity in entities.values) {
      entity.updatedThisFrame = false;
    }

    int playerId = data.getUint32(4);

    for (int el = 0; el < entityCount; el++) {
      int offset = 8 + el * 16;
      final entityType = data.getUint8(offset + 0);
      final entityColorIndex = data.getUint8(offset + 1);
      final entitySize = data.getUint16(offset + 2, Endian.big);
      final entityID = data.getUint32(offset + 4, Endian.big);
      final entityRelPosX =
          data.getInt16(offset + 8, Endian.big).toDouble() / precision;
      final entityRelPosY =
          data.getInt16(offset + 10, Endian.big).toDouble() / precision;

      final existing = entities[entityID.toString()];
      if (existing != null) {
        existing.provideUpdate(
          Offset(entityRelPosX, entityRelPosY),
          entitySize,
        );
      } else {
        final newEntity = GameWorldEntity(
          id: entityID,
          type: (entityType == 1)
              ? GameEntityType.player
              : GameEntityType.crumb,
          colorIndex: entityColorIndex.abs() % blobColors.length,
          size: entitySize,
          relPos: Offset(entityRelPosX, entityRelPosY),
          relPosTarget: Offset(entityRelPosX, entityRelPosY),
        );
        entities[entityID.toString()] = newEntity;
      }

      if (entityID == playerId) {
        // this is our player
        targetZoom = 120.0 / (entitySize + 110) * 20.0;
      }
    }

    entities.removeWhere((id, entity) => !entity.updatedThisFrame);

    notifyListeners();
  }

  void startDirectingPlayer(Offset screenPos) {
    directingPlayer = true;
    playerDirection = _convertScreenToWorld(screenPos);
    sendPlayerDirection?.call(playerDirection);
    _startInputTimer();
  }

  void updatePlayerDirection(Offset screenPos) {
    playerDirection = _convertScreenToWorld(screenPos);
  }

  void stopDirectingPlayer() {
    directingPlayer = false;
    _stopInputTimer();
  }

  void reset() {}

  void update(double delta) {
    zoom = zoom + (targetZoom - zoom) * delta;
    for (var entity in entities.values) {
      entity.update(delta);
    }
  }

  void updateScreenSize(Size newSize) {
    if (screenSize != newSize) {
      screenSize = newSize;
    }
  }

  void _startInputTimer() {
    _stopInputTimer();
    _inputTimer = Timer.periodic(const Duration(milliseconds: 100), (timer) {
      sendPlayerDirection?.call(playerDirection);
    });
  }

  void _stopInputTimer() {
    _inputTimer?.cancel();
    _inputTimer = null;
  }

  Offset _convertScreenToWorld(Offset screenPos) {
    if (screenSize == null) return Offset.zero;

    final center = Offset(screenSize!.width / 2, screenSize!.height / 2);

    return Offset(
      (screenPos.dx - center.dx) / zoom,
      (screenPos.dy - center.dy) / zoom,
    );
  }

  @override
  void dispose() {
    _stopInputTimer();
    super.dispose();
  }
}
