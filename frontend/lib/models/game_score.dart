import 'dart:convert';
import 'dart:typed_data';

import 'package:flutter/material.dart';

class GameScore extends ChangeNotifier {
  List<GameScoreListItem> scoreList = [];

  void parseBinary(ByteData data) {
    int entityCount = data.getUint16(2, Endian.big);

    prepareForNewList();

    int offset = 4;
    for (int el = 0; el < entityCount; el++) {
      final entityId = data.getUint32(offset + 0);
      final entitySize = data.getUint16(offset + 4);
      final entityColorIndex = data.getUint8(offset + 6);
      final entityNameLength = data.getUint8(offset + 7);

      // Create a view of just the name bytes (zero-copy)
      final nameBytes = data.buffer.asUint8List(
        data.offsetInBytes + offset + 8,
        entityNameLength,
      );

      // Decode UTF-8 bytes to String
      final entityName = utf8.decode(nameBytes);

      updatePlayer(entityId, entitySize, entityColorIndex, entityName);

      offset += 8 + entityNameLength;
    }

    removeUnupdated();
    sortPlayers();

    notifyListeners();
  }

  void prepareForNewList() {
    for (var entity in scoreList) {
      entity.updatedThisFrame = false;
    }
  }

  void updatePlayer(int playerID, int size, int colorIndex, String name) {
    final index = scoreList.indexWhere((listItem) => listItem.id == playerID);
    if (index != -1) {
      scoreList[index].update(size);
    } else {
      scoreList.add(
        GameScoreListItem(
          id: playerID,
          size: size,
          colorIndex: colorIndex,
          name: name,
        ),
      );
    }
    notifyListeners();
  }

  void removeUnupdated() {
    scoreList.removeWhere((listItem) => !listItem.updatedThisFrame);
    notifyListeners();
  }

  void sortPlayers() {
    scoreList.sort((a, b) => b.size.compareTo(a.size));
    notifyListeners();
  }

  GameScoreListItem? getScoreItemById(int playerID) {
    final index = scoreList.indexWhere((listItem) => listItem.id == playerID);
    if (index != -1) {
      return scoreList[index];
    } else {
      return null;
    }
  }
}

class GameScoreListItem {
  bool updatedThisFrame = true;
  int id;
  int size;
  int colorIndex = 0;
  String name;

  GameScoreListItem({
    required this.id,
    required this.size,
    required this.name,
    required this.colorIndex,
  });

  void update(int size) {
    this.size = size;
    updatedThisFrame = true;
  }
}
